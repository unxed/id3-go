package encodedbytes

import "testing"

func TestConverters_RoundTrip(t *testing.T) {
	for i, s := range []string{"café", "café ✓", "café ✓", "café ✓"} {
		enc, err := Encoders[i].ConvertString(s)
		if err != nil {
			t.Fatalf("%s: encode: %v", EncodingForIndex(byte(i)), err)
		}
		dec, err := Decoders[i].ConvertString(enc)
		if err != nil {
			t.Fatalf("%s: decode: %v", EncodingForIndex(byte(i)), err)
		}
		if dec != s {
			t.Errorf("%s: round trip = %q, want %q", EncodingForIndex(byte(i)), dec, s)
		}
	}
}

func TestConverters_Bytes(t *testing.T) {
	for _, tc := range []struct {
		encoding byte
		in, want string
	}{
		{0, "é", "\xe9"},
		{1, "a", "\xfe\xff\x00a"},
		{2, "a", "\x00a"},
		{3, "é", "é"},
	} {
		got, err := Encoders[tc.encoding].ConvertString(tc.in)
		if err != nil {
			t.Fatalf("%s: %v", EncodingForIndex(tc.encoding), err)
		}
		if got != tc.want {
			t.Errorf("%s: encoded %q as %q, want %q", EncodingForIndex(tc.encoding), tc.in, got, tc.want)
		}
	}

	// Little-endian UTF-16 with a BOM must decode as well.
	if got, err := Decoders[1].ConvertString("\xff\xfea\x00"); err != nil || got != "a" {
		t.Errorf("UTF-16LE with BOM decoded as %q, %v", got, err)
	}

	if _, err := Encoders[0].ConvertString("✓"); err == nil {
		t.Error("ISO-8859-1 encoder accepted a character outside Latin-1")
	}
}
