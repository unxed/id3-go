package v2

import "testing"

// A v2.3 tag has no UTF-8 encoding (that is v2.4's byte 3); writing text
// into one must use UTF-16, and only a v2.4 tag gets UTF-8.
func TestSetText_EncodingFollowsTagVersion(t *testing.T) {
	for _, tc := range []struct {
		version byte
		want    string
	}{
		{2, "UTF-16"},
		{3, "UTF-16"},
		{4, "UTF-8"},
	} {
		tag := NewTag(tc.version)
		tag.SetArtist("Paloalto ✓")
		frame := tag.textFrame(tag.commonMap["Artist"])
		if frame == nil {
			t.Fatalf("v2.%d: no artist frame written", tc.version)
		}
		if got := frame.Encoding(); got != tc.want {
			t.Errorf("v2.%d: artist written as %s, want %s", tc.version, got, tc.want)
		}
		if got := tag.Artist(); got != "Paloalto ✓" {
			t.Errorf("v2.%d: artist reads back as %q", tc.version, got)
		}
	}
}
