package id3

import (
	"io/ioutil"
	"os"
	"testing"
)

func TestShiftBytesBack(t *testing.T) {
	tempFile, err := ioutil.TempFile("", "shift_bytes")
	if err != nil {
		t.Fatal(err)
	}
	defer os.Remove(tempFile.Name())
	defer tempFile.Close()

	if _, err := tempFile.Write([]byte("ABCDEFGHIJ")); err != nil {
		t.Fatal(err)
	}

	// Move everything from offset 4 on 3 bytes towards the end of the file.
	if err := shiftBytesBack(tempFile, 4, 3); err != nil {
		t.Fatalf("shiftBytesBack failed: %v", err)
	}

	result, err := ioutil.ReadFile(tempFile.Name())
	if err != nil {
		t.Fatal(err)
	}
	if len(result) != 13 {
		t.Fatalf("file length = %d, want 13", len(result))
	}
	if string(result[:4]) != "ABCD" {
		t.Errorf("prefix = %q, want %q", result[:4], "ABCD")
	}
	if string(result[7:]) != "EFGHIJ" {
		t.Errorf("shifted data = %q, want %q", result[7:], "EFGHIJ")
	}
}
