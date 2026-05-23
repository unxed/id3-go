package encodedbytes

import (
	"bytes"
	"testing"
)

func TestByteInt_Errors(t *testing.T) {
	// Length exceeds BytesPerInt
	_, err := ByteInt([]byte{1, 2, 3, 4, 5}, 8)
	if err == nil {
		t.Error("expected error for slice longer than BytesPerInt")
	}

	// Byte value exceeds max bit
	_, err = ByteInt([]byte{0x80}, 7)
	if err == nil {
		t.Error("expected error for byte exceeding max bit")
	}
}

func TestReader_Edges(t *testing.T) {
	data := []byte("Hello World\x00UTF-16\x00\x00")
	r := NewReader(data)

	// Read empty slice
	n, err := r.Read(nil)
	if n != 0 || err != nil {
		t.Errorf("expected (0, nil), got (%d, %v)", n, err)
	}

	// Read rest
	r2 := NewReader([]byte{1, 2, 3})
	rest, err := r2.ReadRest()
	if err != nil || !bytes.Equal(rest, []byte{1, 2, 3}) {
		t.Errorf("ReadRest failed, got %v, err %v", rest, err)
	}
	_, err = r2.ReadByte()
	if err == nil {
		t.Error("expected EOF on subsequent ReadByte")
	}
}

func TestWriter_Edges(t *testing.T) {
	buf := make([]byte, 5)
	w := NewWriter(buf)

	// Write empty slice
	n, err := w.Write(nil)
	if n != 0 || err != nil {
		t.Errorf("expected (0, nil), got (%d, %v)", n, err)
	}

	// Overwrite boundary (copies only up to destination length)
	n, err = w.Write([]byte{1, 2, 3, 4, 5, 6})
	if err != nil {
		t.Errorf("unexpected error on partial copy: %v", err)
	}
	if n != 5 {
		t.Errorf("expected to write 5 bytes, wrote %d", n)
	}

	// Next write should return EOF
	_, err = w.Write([]byte{7})
	if err == nil {
		t.Error("expected EOF on subsequent Write")
	}

	// WriteByte boundary
	w2 := NewWriter(make([]byte, 1))
	err = w2.WriteByte(1)
	if err != nil {
		t.Fatal(err)
	}
	err = w2.WriteByte(2)
	if err == nil {
		t.Error("expected error when writing byte beyond limit")
	}
}

func TestEncodedDiff(t *testing.T) {
	diff, err := EncodedDiff(3, "Hello", 0, "Hello")
	if err != nil {
		t.Fatal(err)
	}
	if diff != 0 {
		t.Errorf("expected diff 0, got %d", diff)
	}
}