package id3

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"os"
	"testing"

	v1 "github.com/mikkyang/id3-go/v1"
)

func TestV1_SetComment(t *testing.T) {
	tag := &v1.Tag{}
	tag.SetComment("My Comment")
	comments := tag.Comments()
	if len(comments) != 1 || comments[0] != "My Comment" {
		t.Errorf("expected comment 'My Comment', got %v", comments)
	}
}

func TestJSONConvert_ErrorHandling(t *testing.T) {
	err := ConvertToJSON("non_existent_file_path.mp3", "output.json")
	if err == nil {
		t.Error("expected error for non-existent file in ConvertToJSON")
	}

	err = ConvertToID3("non_existent_file_path.json", "output.mp3")
	if err == nil {
		t.Error("expected error for non-existent JSON path in ConvertToID3")
	}
}

func TestJSONConvert_AdvancedFramesRoundtrip(t *testing.T) {
	tempFile, err := ioutil.TempFile("", "json_adv")
	if err != nil {
		t.Fatal(err)
	}
	defer os.Remove(tempFile.Name())
	defer tempFile.Close()

	dummyAudio := make([]byte, 200)
	if _, err := tempFile.Write(dummyAudio); err != nil {
		t.Fatal(err)
	}

	jsonFile, err := ioutil.TempFile("", "adv_json")
	if err != nil {
		t.Fatal(err)
	}
	jsonPath := jsonFile.Name()
	jsonFile.Close()
	defer os.Remove(jsonPath)

	meta := Metadata{
		ID3v2: &V2Metadata{
			Frames: []V2Frame{
				{ID: "TCON", Type: "text", Value: "Alternative"},
				{ID: "TXXX", Type: "desc_text", Description: "Description", Value: "Value"},
				{ID: "COMM", Type: "unsynch_text", Language: "deu", Description: "MyDesc", Value: "MyComment"},
				{ID: "PRIV", Type: "data", Data: []byte{1, 2, 3, 4}},
			},
		},
	}

	importData, err := json.Marshal(meta)
	if err != nil {
		t.Fatal(err)
	}
	if err := ioutil.WriteFile(jsonPath, importData, 0644); err != nil {
		t.Fatal(err)
	}

	// Convert JSON -> MP3
	if err := ConvertToID3(jsonPath, tempFile.Name()); err != nil {
		t.Fatalf("failed to convert to ID3: %v", err)
	}

	// Convert MP3 -> JSON
	os.Remove(jsonPath)
	if err := ConvertToJSON(tempFile.Name(), jsonPath); err != nil {
		t.Fatalf("failed to convert to JSON: %v", err)
	}

	// Read and verify
	exportedData, err := ioutil.ReadFile(jsonPath)
	if err != nil {
		t.Fatal(err)
	}

	var exported Metadata
	if err := json.Unmarshal(exportedData, &exported); err != nil {
		t.Fatal(err)
	}

	if exported.ID3v2 == nil {
		t.Fatal("expected exported ID3v2 metadata")
	}

	// Map frames by ID for easy checks
	framesMap := make(map[string]V2Frame)
	for _, f := range exported.ID3v2.Frames {
		framesMap[f.ID] = f
	}

	// Verify TCON
	if f, ok := framesMap["TCON"]; !ok || f.Value != "Alternative" {
		t.Errorf("TCON failed: %+v", f)
	}

	// Verify TXXX
	if f, ok := framesMap["TXXX"]; !ok || f.Description != "Description" || f.Value != "Value" {
		t.Errorf("TXXX failed: %+v", f)
	}

	// Verify COMM
	if f, ok := framesMap["COMM"]; !ok || f.Language != "deu" || f.Description != "MyDesc" || f.Value != "MyComment" {
		t.Errorf("COMM failed: %+v", f)
	}

	// Verify PRIV
	if f, ok := framesMap["PRIV"]; !ok || !bytes.Equal(f.Data, []byte{1, 2, 3, 4}) {
		t.Errorf("PRIV failed: %+v", f)
	}
}
