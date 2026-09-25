package v2

import (
	"testing"
)

func TestUnsynchTextFrameSetEncoding(t *testing.T) {
	f := NewUnsynchTextFrame(V23CommonFrame["Comments"], "Foo", "Bar")
	size := f.Size()
	expectedDiff := 11

	err := f.SetEncoding("UTF-16")
	if err != nil {
		t.Fatal(err)
	}
	newSize := f.Size()
	if int(newSize-size) != expectedDiff {
		t.Errorf("expected size to increase to %d, but it was %d", size+1, newSize)
	}

	size = newSize
	err = f.SetEncoding("ISO-8859-1")
	if err != nil {
		t.Fatal(err)
	}
	newSize = f.Size()
	if int(newSize-size) != -expectedDiff {
		t.Errorf("expected size to decrease to %d, but it was %d", size-1, newSize)
	}
}

func TestDescTextFrame_BytesRoundTrip(t *testing.T) {
	f := NewDescTextFrame(V23FrameTypeMap["TXXX"], "Description", "Value")
	data := f.Bytes()
	if uint(len(data)) != f.Size() {
		t.Fatalf("Bytes() length %d, Size() %d", len(data), f.Size())
	}

	parsed, ok := ParseDescTextFrame(f.FrameHead, data).(*DescTextFrame)
	if !ok {
		t.Fatal("ParseDescTextFrame failed")
	}
	if got := parsed.Description(); got != "Description" {
		t.Errorf("description = %q, want %q", got, "Description")
	}
	if got := parsed.Text(); got != "Value" {
		t.Errorf("text = %q, want %q", got, "Value")
	}
}
