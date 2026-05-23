package id3

import (
	"bytes"
	"io/ioutil"
	"os"
	"strings"
	"testing"

	"github.com/unxed/localecp"
	v1 "github.com/unxed/id3-go/v1"
	v2 "github.com/unxed/id3-go/v2"
	"golang.org/x/text/encoding/charmap"
)

func TestV1_Detailed(t *testing.T) {
	// Parse non-existent / invalid V1 tag
	invalidData := []byte("Not a tag at all")
	v1Tag := v1.ParseTag(bytes.NewReader(invalidData))
	if v1Tag != nil {
		t.Error("expected ParseTag to return nil for invalid V1 tag")
	}

	// Create and modify V1 Tag
	tag := &v1.Tag{}
	tag.SetTitle("Title Test")
	tag.SetArtist("Artist Test")
	tag.SetAlbum("Album Test")
	tag.SetYear("2026")
	tag.SetGenre("Blues") // Valid genre

	if tag.Title() != "Title Test" {
		t.Errorf("expected Title Test, got %s", tag.Title())
	}
	if tag.Genre() != "Blues" {
		t.Errorf("expected Blues, got %s", tag.Genre())
	}

	tag.SetGenre("Non-existent Genre")
	if tag.Genre() != "" {
		t.Errorf("expected empty string for invalid genre, got %q", tag.Genre())
	}

	if tag.Version() != "1.0" {
		t.Errorf("expected version 1.0, got %s", tag.Version())
	}

	if tag.Padding() != 0 {
		t.Errorf("expected padding 0, got %d", tag.Padding())
	}

	if len(tag.AllFrames()) != 0 {
		t.Error("expected no frames in V1")
	}

	bytesData := tag.Bytes()
	if len(bytesData) != v1.TagSize {
		t.Errorf("expected size %d, got %d", v1.TagSize, len(bytesData))
	}
}

func TestV2_FrameCreationAndModification(t *testing.T) {
	tag := v2.NewTag(3)

	// Add dynamic Frame
	ft := v2.V23FrameTypeMap["TIT2"]
	frame := v2.NewTextFrame(ft, "My Title")
	tag.AddFrames(frame)

	if tag.Title() != "My Title" {
		t.Errorf("expected Title 'My Title', got %q", tag.Title())
	}

	// Change Title Text
	tag.SetTitle("New Title")
	if tag.Title() != "New Title" {
		t.Errorf("expected Title 'New Title', got %q", tag.Title())
	}

	// Non-existent Frame Check
	nonExistent := tag.Frame("ABCD")
	if nonExistent != nil {
		t.Error("expected non-existent frame to be nil")
	}

	// Check Comments
	tag.SetGenre("Rock")
	if tag.Genre() != "Rock" {
		t.Errorf("expected Genre 'Rock', got %q", tag.Genre())
	}

	// Check V2 Bytes and dirty flag
	if !tag.Dirty() {
		t.Error("tag should be marked as dirty after modifications")
	}
}

func TestDataFrame_Detailed(t *testing.T) {
	ft := v2.V23FrameTypeMap["PRIV"]
	dataFrame := v2.NewDataFrame(ft, []byte{1, 2, 3, 4})

	if dataFrame.Id() != "PRIV" {
		t.Errorf("expected PRIV, got %s", dataFrame.Id())
	}

	if !bytes.Equal(dataFrame.Data(), []byte{1, 2, 3, 4}) {
		t.Errorf("unexpected data: %v", dataFrame.Data())
	}

	dataFrame.SetData([]byte{5, 6})
	if !bytes.Equal(dataFrame.Data(), []byte{5, 6}) {
		t.Errorf("unexpected data after set: %v", dataFrame.Data())
	}

	if dataFrame.String() != "<binary data>" {
		t.Errorf("unexpected string representation: %s", dataFrame.String())
	}
}

func TestIdFrame_Detailed(t *testing.T) {
	ft := v2.V23FrameTypeMap["UFID"]
	idFrame := v2.NewIdFrame(ft, "http://example.com", []byte{0xDE, 0xAD})

	if idFrame.OwnerIdentifier() != "http://example.com" {
		t.Errorf("expected owner identifier, got %s", idFrame.OwnerIdentifier())
	}

	if !bytes.Equal(idFrame.Identifier(), []byte{0xDE, 0xAD}) {
		t.Errorf("expected identifier, got %v", idFrame.Identifier())
	}

	err := idFrame.SetIdentifier(make([]byte, 100))
	if err == nil {
		t.Error("expected error when setting identifier exceeding 64 bytes")
	}

	idFrame.SetOwnerIdentifier("http://new.com")
	if idFrame.OwnerIdentifier() != "http://new.com" {
		t.Errorf("expected http://new.com, got %s", idFrame.OwnerIdentifier())
	}
}

func TestImageFrame_Detailed(t *testing.T) {
	head := v2.FrameHead{} // mock / blank head
	invalidImageFrame := v2.ParseImageFrame(head, []byte{0x00})
	if invalidImageFrame != nil {
		t.Error("expected ParseImageFrame to fail on truncated data")
	}
}

func TestParse_CorruptedOrTruncatedV2(t *testing.T) {
	// Truncated header
	truncatedHeader := []byte("ID3\x03\x00\x00")
	r := bytes.NewReader(truncatedHeader)
	tag := v2.ParseTag(r)
	if tag != nil {
		t.Error("expected ParseTag to return nil on truncated header")
	}
}

func TestFileOpenSave_Edges(t *testing.T) {
	tempFile, err := ioutil.TempFile("", "id3go_edge")
	if err != nil {
		t.Fatal(err)
	}
	defer os.Remove(tempFile.Name())
	defer tempFile.Close()

	// Parse on an empty file (this creates a new V2 tag)
	f, err := Open(tempFile.Name())
	if err != nil {
		t.Fatalf("unexpected open error: %v", err)
	}

	// Modify tag data
	f.SetArtist("Test Artist")
	f.SetTitle("Test Title")

	// Save/Close
	err = f.Close()
	if err != nil {
		t.Fatalf("unexpected close error: %v", err)
	}

	// Reopen and check
	f2, err := Open(tempFile.Name())
	if err != nil {
		t.Fatalf("unexpected reopen error: %v", err)
	}
	defer f2.Close()

	if f2.Artist() != "Test Artist" {
		t.Errorf("expected Test Artist, got %s", f2.Artist())
	}
	if f2.Title() != "Test Title" {
		t.Errorf("expected Test Title, got %s", f2.Title())
	}
}

func TestShiftBytesBack_Direct(t *testing.T) {
	tempFile, err := ioutil.TempFile("", "shift_bytes")
	if err != nil {
		t.Fatal(err)
	}
	defer os.Remove(tempFile.Name())
	defer tempFile.Close()

	initialData := []byte("ABCDEFGHIJ")
	if _, err := tempFile.Write(initialData); err != nil {
		t.Fatal(err)
	}

	// Shift back by 3 bytes from position 4
	err = shiftBytesBack(tempFile, 4, 3)
	if err != nil {
		t.Fatalf("shiftBytesBack failed: %v", err)
	}

	result, err := ioutil.ReadFile(tempFile.Name())
	if err != nil {
		t.Fatal(err)
	}

	// Expected total length should be 10 + 3 = 13
	if len(result) != 13 {
		t.Errorf("expected file length 13, got %d", len(result))
	}

	// First 4 bytes must be unaffected: "ABCD"
	if string(result[:4]) != "ABCD" {
		t.Errorf("expected prefix ABCD, got %q", string(result[:4]))
	}

	// Last 6 bytes should be shifted: "EFGHIJ"
	if string(result[7:]) != "EFGHIJ" {
		t.Errorf("expected shifted suffix EFGHIJ, got %q", string(result[7:]))
	}
}

func TestV2_FrameDeletionsAndAllFrames(t *testing.T) {
	tag := v2.NewTag(3)
	ft1 := v2.V23FrameTypeMap["TIT2"]
	ft2 := v2.V23FrameTypeMap["TPE1"]

	f1 := v2.NewTextFrame(ft1, "Title 1")
	f2 := v2.NewTextFrame(ft2, "Artist 1")

	tag.AddFrames(f1, f2)

	all := tag.AllFrames()
	if len(all) != 2 {
		t.Errorf("expected 2 frames, got %d", len(all))
	}

	deleted := tag.DeleteFrames("TIT2")
	if len(deleted) != 1 {
		t.Errorf("expected deleted 1 frame, got %d", len(deleted))
	}

	allAfter := tag.AllFrames()
	if len(allAfter) != 1 {
		t.Errorf("expected 1 frame after deletion, got %d", len(allAfter))
	}

	if tag.Frame("TIT2") != nil {
		t.Error("expected TIT2 frame to be deleted")
	}

	// Delete non-existent frame
	deletedNone := tag.DeleteFrames("ABCD")
	if len(deletedNone) != 0 {
		t.Errorf("expected 0 deleted frames, got %d", len(deletedNone))
	}
}

func TestDescTextFrame_Direct(t *testing.T) {
	ft := v2.V23FrameTypeMap["TXXX"]
	f := v2.NewDescTextFrame(ft, "Description", "Value")

	if f.Description() != "Description" {
		t.Errorf("expected Description, got %s", f.Description())
	}
	if f.Text() != "Value" {
		t.Errorf("expected Value, got %s", f.Text())
	}

	err := f.SetDescription("New Description")
	if err != nil {
		t.Fatal(err)
	}
	if f.Description() != "New Description" {
		t.Errorf("expected New Description, got %s", f.Description())
	}

	err = f.SetEncoding("UTF-16")
	if err != nil {
		t.Fatal(err)
	}

	if f.Encoding() != "UTF-16" {
		t.Errorf("expected UTF-16, got %s", f.Encoding())
	}
}

func TestUnsynchTextFrame_Direct(t *testing.T) {
	ft := v2.V23FrameTypeMap["COMM"]
	f := v2.NewUnsynchTextFrame(ft, "Comment Description", "Comment Content")

	if f.Language() != "eng" {
		t.Errorf("expected eng, got %s", f.Language())
	}

	err := f.SetLanguage("rus")
	if err != nil {
		t.Fatal(err)
	}
	if f.Language() != "rus" {
		t.Errorf("expected rus, got %s", f.Language())
	}

	err = f.SetLanguage("invalid")
	if err == nil {
		t.Error("expected error for invalid language code length")
	}
}

func TestV1_CloseAndSave(t *testing.T) {
	tempFile, err := ioutil.TempFile("", "v1_save")
	if err != nil {
		t.Fatal(err)
	}
	defer os.Remove(tempFile.Name())
	defer tempFile.Close()

	// Write 500 bytes of dummy audio data
	dummyAudio := make([]byte, 500)
	for i := range dummyAudio {
		dummyAudio[i] = byte(i % 256)
	}
	if _, err := tempFile.Write(dummyAudio); err != nil {
		t.Fatal(err)
	}

	// Now append a valid ID3v1 tag to the end of the file
	v1Tag := &v1.Tag{}
	v1Tag.SetTitle("V1 Title")
	v1Tag.SetArtist("V1 Artist")
	if _, err := tempFile.Write(v1Tag.Bytes()); err != nil {
		t.Fatal(err)
	}

	// Open the file with id3-go
	f, err := Open(tempFile.Name())
	if err != nil {
		t.Fatalf("failed to open file: %v", err)
	}

	// Verify it has been parsed as ID3v1
	if f.Version() != "1.0" {
		t.Errorf("expected version 1.0, got %s", f.Version())
	}

	title := strings.TrimRight(f.Title(), "\x00")
	if title != "V1 Title" {
		t.Errorf("expected Title 'V1 Title', got %q", title)
	}

	// Modify tag details
	f.SetTitle("New V1 Title")

	// Close the file
	err = f.Close()
	if err != nil {
		t.Fatalf("failed to close and save file: %v", err)
	}

	// Re-verify after saving
	f2, err := Open(tempFile.Name())
	if err != nil {
		t.Fatalf("failed to reopen file: %v", err)
	}
	defer f2.Close()

	title2 := strings.TrimRight(f2.Title(), "\x00")
	if title2 != "New V1 Title" {
		t.Errorf("expected Title 'New V1 Title', got %q", title2)
	}
}

type dummyTagger struct {
	*v1.Tag
}

func TestClose_UnknownVersion(t *testing.T) {
	tempFile, err := ioutil.TempFile("", "dummy")
	if err != nil {
		t.Fatal(err)
	}
	defer os.Remove(tempFile.Name())

	f := &File{
		Tagger: dummyTagger{&v1.Tag{}},
		file:   tempFile,
	}
	f.SetTitle("A") // makes dirty = true

	err = f.Close()
	if err == nil {
		t.Error("expected error for unknown tag version")
	} else if err.Error() != "Close: unknown tag version" {
		t.Errorf("expected 'Close: unknown tag version', got %q", err.Error())
	}
}

func TestV1_LocaleDecoding(t *testing.T) {
	// Backup original decoder
	origDecoder := localecp.ANSIDecoder
	defer func() {
		localecp.ANSIDecoder = origDecoder
	}()

	// Set specifically to Windows-1251 for non-ASCII tests
	localecp.ANSIDecoder = charmap.Windows1251.NewDecoder()

	tempFile, err := ioutil.TempFile("", "v1_locale")
	if err != nil {
		t.Fatal(err)
	}
	defer os.Remove(tempFile.Name())
	defer tempFile.Close()

	// Write 128 bytes of ID3v1 tag directly
	tagData := make([]byte, 128)
	copy(tagData[:3], "TAG")
	// "Привет" in CP1251: 0xcf, 0xf0, 0xe8, 0xe2, 0xe5, 0xf2
	copy(tagData[3:3+6], []byte{0xcf, 0xf0, 0xe8, 0xe2, 0xe5, 0xf2})

	if _, err := tempFile.Write(tagData); err != nil {
		t.Fatal(err)
	}

	// Open with id3-go
	f, err := Open(tempFile.Name())
	if err != nil {
		t.Fatal(err)
	}
	defer f.Close()

	title := strings.TrimRight(f.Title(), "\x00")
	if title != "Привет" {
		t.Errorf("expected 'Привет' after legacy decoding, got %q", title)
	}
}
func TestV1_LocaleEncoding(t *testing.T) {
	// Backup original encoder
	origEncoder := localecp.ANSIEncoder
	defer func() {
		localecp.ANSIEncoder = origEncoder
	}()

	// Set specifically to Windows-1251
	localecp.ANSIEncoder = charmap.Windows1251.NewEncoder()

	tag := &v1.Tag{}
	tag.SetTitle("Привет")

	bytesData := tag.Bytes()
	// "Привет" в кодировке CP1251: 0xcf, 0xf0, 0xe8, 0xe2, 0xe5, 0xf2
	expectedBytes := []byte{0xcf, 0xf0, 0xe8, 0xe2, 0xe5, 0xf2}

	if !bytes.Equal(bytesData[3:3+len(expectedBytes)], expectedBytes) {
		t.Errorf("expected CP1251 bytes %v, got %v", expectedBytes, bytesData[3:3+len(expectedBytes)])
	}
}
func TestV1_LocaleEncoding_Fallback(t *testing.T) {
	// Backup original encoder
	origEncoder := localecp.ANSIEncoder
	defer func() {
		localecp.ANSIEncoder = origEncoder
	}()

	// Set specifically to Windows-1251 (Cyrillic)
	localecp.ANSIEncoder = charmap.Windows1251.NewEncoder()

	tag := &v1.Tag{}
	// Japanese text which cannot be encoded into CP1251
	japaneseText := "日本"
	tag.SetTitle(japaneseText)

	bytesData := tag.Bytes()

	// Since Japanese cannot be represented in CP1251, it must fall back to raw UTF-8 bytes
	expectedUTF8 := []byte(japaneseText)
	if !bytes.Equal(bytesData[3:3+len(expectedUTF8)], expectedUTF8) {
		t.Errorf("expected fallback to UTF-8 bytes %v, got %v", expectedUTF8, bytesData[3:3+len(expectedUTF8)])
	}
}
func TestV1_DummyMethods(t *testing.T) {
	tag := &v1.Tag{}
	tag.SetArtist("Artist")

	// Comments
	comments := tag.Comments()
	if len(comments) != 1 {
		t.Errorf("expected 1 comment, got %d", len(comments))
	}

	// Padding
	if tag.Padding() != 0 {
		t.Errorf("expected padding 0, got %d", tag.Padding())
	}

	// AllFrames
	if len(tag.AllFrames()) != 0 {
		t.Error("expected 0 frames")
	}

	// Frame / Frames / DeleteFrames / AddFrames
	if tag.Frame("TIT2") != nil {
		t.Error("expected nil frame")
	}
	if len(tag.Frames("TIT2")) != 0 {
		t.Error("expected 0 frames slice")
	}
	if len(tag.DeleteFrames("TIT2")) != 0 {
		t.Error("expected 0 deleted frames")
	}
	// AddFrames should not panic
	tag.AddFrames(nil)
}

func TestOpen_NonExistentFile(t *testing.T) {
	_, err := Open("non_existent_file_path_12345.mp3")
	if err == nil {
		t.Error("expected error when opening non-existent file")
	}
}

func TestV2_RealSize(t *testing.T) {
	tag := v2.NewTag(3)
	// Add some frame to change size
	ft := v2.V23FrameTypeMap["TIT2"]
	frame := v2.NewTextFrame(ft, "My Title")
	tag.AddFrames(frame)

	realSize := tag.RealSize()
	// RealSize = size - padding.
	// Since we created Tag and added frame, size should match and we can check it's non-negative.
	if realSize < 0 {
		t.Errorf("expected non-negative RealSize, got %d", realSize)
	}
}

func TestParseHeader_Validation(t *testing.T) {
	// Invalid signature
	invalidSig := []byte("BAD\x03\x00\x00\x00\x00\x00\x00")
	h1 := v2.ParseHeader(bytes.NewReader(invalidSig))
	if h1 != nil {
		t.Error("expected nil header for invalid signature")
	}

	// Invalid synchsafe size (bit 7 set in size bytes)
	invalidSize := []byte("ID3\x03\x00\x00\x80\x00\x00\x00")
	h2 := v2.ParseHeader(bytes.NewReader(invalidSize))
	if h2 != nil {
		t.Error("expected nil header for invalid synchsafe size")
	}
}

func TestShiftBytesBack_ClosedFile(t *testing.T) {
	tempFile, err := ioutil.TempFile("", "closed_file")
	if err != nil {
		t.Fatal(err)
	}
	os.Remove(tempFile.Name())
	tempFile.Close()

	// Calling on closed file should fail
	err = shiftBytesBack(tempFile, 0, 10)
	if err == nil {
		t.Error("expected error when calling shiftBytesBack on closed file")
	}
}
