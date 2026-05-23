package v2

import (
	"bytes"
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
func TestImageFrame_Success(t *testing.T) {
	// Let's create an ImageFrame manually via parsing
	// Encoding = 0 (ISO-8859-1)
	// MIME Type = "image/png" + null (10 bytes)
	// Picture Type = 3
	// Description = "Cover" + null (6 bytes)
	// Data = 0x11, 0x22, 0x33
	frameData := []byte{0, 'i', 'm', 'a', 'g', 'e', '/', 'p', 'n', 'g', 0, 3, 'C', 'o', 'v', 'e', 'r', 0, 0x11, 0x22, 0x33}

	ft := V23FrameTypeMap["APIC"]
	head := FrameHead{
		FrameType: ft,
		size:      uint32(len(frameData)),
	}

	framer := ParseImageFrame(head, frameData)
	if framer == nil {
		t.Fatal("expected ParseImageFrame to succeed")
	}

	img, ok := framer.(*ImageFrame)
	if !ok {
		t.Fatal("expected *ImageFrame type")
	}

	if img.MIMEType() != "image/png\x00" && img.MIMEType() != "image/png" {
		t.Errorf("expected image/png, got %q", img.MIMEType())
	}

	if img.Encoding() != "ISO-8859-1" {
		t.Errorf("expected ISO-8859-1, got %s", img.Encoding())
	}

	img.SetMIMEType("image/jpeg")

	err := img.SetEncoding("UTF-8")
	if err != nil {
		t.Fatal(err)
	}

	if img.String() == "" {
		t.Error("expected non-empty string representation")
	}

	imgBytes := img.Bytes()
	if len(imgBytes) == 0 {
		t.Error("expected non-empty byte slice")
	}
}
func TestParseV23Frame_Unknown(t *testing.T) {
	// Let's pass a header with unknown ID "ZZZZ"
	headerData := []byte("ZZZZ\x00\x00\x00\x00\x00\x00")
	framer := ParseV23Frame(bytes.NewReader(headerData), false)
	if framer != nil {
		t.Error("expected ParseV23Frame to return nil for unknown frame ID")
	}
}

func TestParseV22Frame_Unknown(t *testing.T) {
	// V2.2 frame header size is 6. Let's pass "ZZZ"
	headerData := []byte("ZZZ\x00\x00\x00")
	framer := ParseV22Frame(bytes.NewReader(headerData), false)
	if framer != nil {
		t.Error("expected ParseV22Frame to return nil for unknown frame ID")
	}
}
