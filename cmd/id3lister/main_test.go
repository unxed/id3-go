package main

import (
	"bytes"
	"io/ioutil"
	"os"
	"strings"
	"testing"

	v1 "github.com/unxed/id3-go/v1"
)

func TestListFile_Success(t *testing.T) {
	tempFile, err := ioutil.TempFile("", "lister_test")
	if err != nil {
		t.Fatal(err)
	}
	defer os.Remove(tempFile.Name())
	defer tempFile.Close()

	// Write dummy data with ID3v1 tag
	dummyAudio := make([]byte, 200)
	if _, err := tempFile.Write(dummyAudio); err != nil {
		t.Fatal(err)
	}

	tag := &v1.Tag{}
	tag.SetTitle("Lister Title")
	tag.SetArtist("Lister Artist")
	if _, err := tempFile.Write(tag.Bytes()); err != nil {
		t.Fatal(err)
	}

	var buf bytes.Buffer
	err = listFile(tempFile.Name(), &buf)
	if err != nil {
		t.Fatalf("unexpected listFile error: %v", err)
	}

	output := buf.String()
	if !strings.Contains(output, "Title: Lister Title") {
		t.Errorf("expected output to contain Title, got:\n%s", output)
	}
	if !strings.Contains(output, "Artist: Lister Artist") {
		t.Errorf("expected output to contain Artist, got:\n%s", output)
	}
}

func TestListFile_Error(t *testing.T) {
	var buf bytes.Buffer
	err := listFile("non_existent_file_path.mp3", &buf)
	if err == nil {
		t.Error("expected error for non-existent file")
	}
}