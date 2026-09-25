package v1

import (
	"bytes"
	"testing"

	"github.com/unxed/localecp"
	"golang.org/x/text/encoding"
	"golang.org/x/text/encoding/charmap"
)

// "Athens" in Greek, as UTF-8 and as Windows-1253.
const (
	greekText  = "Αθήνα"
	greekBytes = "\xc1\xe8\xde\xed\xe1"
)

func withEncoding(t *testing.T, e encoding.Encoding) {
	t.Helper()
	orig := Encoding
	t.Cleanup(func() { Encoding = orig })
	Encoding = e
}

func rawTag(title string) []byte {
	data := make([]byte, TagSize)
	copy(data, "TAG")
	copy(data[3:33], title)
	return data
}

func TestEncoding_DefaultKeepsBytes(t *testing.T) {
	withEncoding(t, nil)

	tag := ParseTag(bytes.NewReader(rawTag(greekBytes)))
	if got := tag.Title()[:len(greekBytes)]; got != greekBytes {
		t.Errorf("title = %q, want the stored bytes %q", got, greekBytes)
	}

	tag = new(Tag)
	tag.SetTitle(greekText)
	if got := string(tag.Bytes()[3 : 3+len(greekText)]); got != greekText {
		t.Errorf("wrote %q, want UTF-8 %q", got, greekText)
	}
}

func TestEncoding_Decode(t *testing.T) {
	withEncoding(t, charmap.Windows1253)

	tag := ParseTag(bytes.NewReader(rawTag(greekBytes)))
	if got := tag.Title()[:len(greekText)]; got != greekText {
		t.Errorf("title = %q, want %q", got, greekText)
	}
}

func TestEncoding_Encode(t *testing.T) {
	withEncoding(t, charmap.Windows1253)

	tag := new(Tag)
	tag.SetTitle(greekText)
	if got := string(tag.Bytes()[3 : 3+len(greekBytes)]); got != greekBytes {
		t.Errorf("wrote %q, want %q", got, greekBytes)
	}
}

func TestEncoding_EncodeFallback(t *testing.T) {
	withEncoding(t, charmap.Windows1253)

	// Not representable in Windows-1253: written as UTF-8.
	const text = "日本"
	tag := new(Tag)
	tag.SetTitle(text)
	if got := string(tag.Bytes()[3 : 3+len(text)]); got != text {
		t.Errorf("wrote %q, want UTF-8 %q", got, text)
	}
}

func TestUseLocaleEncoding(t *testing.T) {
	withEncoding(t, nil)

	UseLocaleEncoding()
	if Encoding != localecp.ANSIEncoding {
		t.Errorf("Encoding = %v, want localecp.ANSIEncoding %v", Encoding, localecp.ANSIEncoding)
	}
}
