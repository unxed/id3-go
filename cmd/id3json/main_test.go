package main

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"testing"

	"github.com/unxed/id3-go"
)

func TestJSONRoundTrip(t *testing.T) {
	// Create a temp MP3 file
	mp3File, err := ioutil.TempFile("", "test_mp3")
	if err != nil {
		t.Fatal(err)
	}
	defer os.Remove(mp3File.Name())
	defer mp3File.Close()

	// Write dummy data so we can open it
	dummyAudio := make([]byte, 200)
	if _, err := mp3File.Write(dummyAudio); err != nil {
		t.Fatal(err)
	}

	// Set up initial V1 and V2 tags manually using ConvertToID3 and a mock JSON
	jsonFile, err := ioutil.TempFile("", "test_json")
	if err != nil {
		t.Fatal(err)
	}
	jsonPath := jsonFile.Name()
	jsonFile.Close()
	defer os.Remove(jsonPath)

	initialMeta := id3.Metadata{
		ID3v1: &id3.V1Metadata{
			Title:   "V1 Initial Title",
			Artist:  "V1 Initial Artist",
			Album:   "V1 Initial Album",
			Year:    "2010",
			Comment: "V1 Comment",
		},
		ID3v2: &id3.V2Metadata{
			Frames: []id3.V2Frame{
				{ID: "TIT2", Type: "text", Value: "V2 Initial Title"},
				{ID: "TPE1", Type: "text", Value: "V2 Initial Artist"},
			},
		},
	}

	initData, err := json.Marshal(initialMeta)
	if err != nil {
		t.Fatal(err)
	}
	if err := ioutil.WriteFile(jsonPath, initData, 0644); err != nil {
		t.Fatal(err)
	}

	// Convert JSON -> MP3 (installs tags)
	if err := id3.ConvertToID3(jsonPath, mp3File.Name()); err != nil {
		t.Fatalf("failed to install initial tags: %v", err)
	}

	// Clear JSON and convert MP3 -> JSON (dumps tags)
	os.Remove(jsonPath)
	if err := id3.ConvertToJSON(mp3File.Name(), jsonPath); err != nil {
		t.Fatalf("failed to dump tags: %v", err)
	}

	// Read and verify dumped JSON
	jsonData, err := ioutil.ReadFile(jsonPath)
	if err != nil {
		t.Fatal(err)
	}

	var dumped id3.Metadata
	if err := json.Unmarshal(jsonData, &dumped); err != nil {
		t.Fatal(err)
	}

	if dumped.ID3v1 == nil || dumped.ID3v1.Title != "V1 Initial Title" {
		t.Errorf("ID3v1 failed, got: %+v", dumped.ID3v1)
	}
	if dumped.ID3v2 == nil || len(dumped.ID3v2.Frames) < 2 {
		t.Errorf("ID3v2 failed, got: %+v", dumped.ID3v2)
	}
}