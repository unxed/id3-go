package id3

import (
	"os"
	"testing"
)

// The bug that corrupted the checked-in fixture: writing text into a v2.3
// file produced frames with encoding byte 3, which v2.3 does not have.
// After a save the tag must still be a valid v2.3 tag.
func TestSave_KeepsV23Encodings(t *testing.T) {
	path := testFixture(t)
	file, err := Open(path)
	if err != nil {
		t.Fatal(err)
	}
	file.SetArtist("Paloalto ✓")
	file.SetTitle("Nice Life")
	if err := file.Close(); err != nil {
		t.Fatal(err)
	}
	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	if data[3] != 3 {
		t.Fatalf("tag version byte = %d, want 3", data[3])
	}
	size := int(data[6])<<21 | int(data[7])<<14 | int(data[8])<<7 | int(data[9])
	for i := 10; i < 10+size; {
		id := string(data[i : i+4])
		if id == "\x00\x00\x00\x00" {
			break
		}
		fsize := int(data[i+4])<<24 | int(data[i+5])<<16 | int(data[i+6])<<8 | int(data[i+7])
		if id[0] == 'T' || id == "COMM" || id == "USLT" {
			if enc := data[i+10]; enc > 1 {
				t.Errorf("frame %s written with encoding byte %d, which v2.3 does not define", id, enc)
			}
		}
		i += 10 + fsize
	}
	reopened, err := Open(path)
	if err != nil {
		t.Fatal(err)
	}
	defer reopened.Close()
	if got := reopened.Artist(); got != "Paloalto ✓" {
		t.Errorf("artist after save = %q", got)
	}
}
