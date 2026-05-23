# id3-go (Fork)

[![build status](https://travis-ci.org/unxed/id3-go.svg)](https://travis-ci.org/unxed/id3-go)

ID3 library for Go.

This is a modernized and maintained fork of the original `github.com/mikkyang/id3-go` library, which seems to be no longer developed.

## Key Improvements in this Fork

* **Locale-Aware ID3v1 Support**: Integrated with `github.com/unxed/localecp` to dynamically deduce the host system's legacy ANSI codepage (such as CP1251 on Cyrillic systems, CP932 on Japanese systems). This ensures that legacy ID3v1 tags are decoded and encoded correctly, matching the behavior of classic desktop media players.
* **Deterministic Tag Serialization**: Refactored internal frame storage to preserve the exact insertion order. Frame sequences remain stable and deterministic when writing tags back to disk, rather than being serialized in random map order.
* **ID3v2 Desynchronization**: Added desynchronization decoding for ID3v2.2 and ID3v2.3 frame contents to safely handle streams where synchsafe integrity is enforced.
* **Critical Bug Fixes**:
  * Corrected file boundary calculations in the `shiftBytesBack` utility to prevent stream corruption when enlarging tags.
  * Added missing null-terminators in `IdFrame` (`UFID`) byte formatting.
  * Stripped trailing padding null-bytes from ID3v1 string parses.
* **Thorough Test Suite**: Added a comprehensive unit and integration testing suite covering edge cases, corrupted/truncated headers, specific frames, and encoding fallbacks.
* **CLI Test Utility (`id3lister`)**: Added a convenient command-line tool `id3lister` to display all tag information (including properly decoded ID3v1 tag values according to the currently active locale).
* **JSON Converter Utility (`id3json`)**: Implemented a bidirectional command-line tool `id3json` to export audio tag metadata to structured JSON and apply it back into files.

Supported formats:

* ID3v1 (with locale-aware ANSI codepage detection)
* ID3v2.2 (with desynchronization support)
* ID3v2.3 (with desynchronization support)

# Install

The platform ($GOROOT/bin) "go get" tool is the best method to install.

    go get github.com/unxed/id3-go

This downloads and installs the package into your $GOPATH. If you only want to
recompile, use "go install".

    go install github.com/unxed/id3-go

# Usage

An import allows access to the package.

    import (
        id3 "github.com/unxed/id3-go"
    )

Version specific details can be accessed through the subpackages.

    import (
        "github.com/unxed/id3-go/v1"
        "github.com/unxed/id3-go/v2"
    )

# Quick Start

To access the tag of a file, first open the file using the package's `Open`
function.

    mp3File, err := id3.Open("All-In.mp3")

It's also a good idea to ensure that the file is closed using `defer`.

    defer mp3File.Close()

## Accessing Information

Some commonly used data have methods in the tag for easier access. These
methods are for `Title`, `Artist`, `Album`, `Year`, `Genre`, and `Comments`.

    mp3File.SetArtist("Okasian")
    fmt.Println(mp3File.Artist())

# ID3v2 Frames

v2 Frames can be accessed directly by using the `Frame` or `Frames` method
of the file, which return the first frame or a slice of frames as `Framer`
interfaces. These interfaces allow read access to general details of the file.

    lyricsFrame := mp3File.Frame("USLT")
    lyrics := lyricsFrame.String()

If more specific information is needed, or frame-specific write access is
needed, then the interface must be cast into the appropriate underlying type.
The example provided does not check for errors, but it is recommended to do
so.

    lyricsFrame := mp3File.Frame("USLT").(*v2.UnsynchTextFrame)

## Adding Frames

For common fields, a frame will automatically be created with the `Set` method.
For other frames or more fine-grained control, frames can be created with the
corresponding constructor, usually prefixed by `New`. These constructors require
the first argument to be a FrameType struct, which are global variables named by
version.

    ft := V23FrameTypeMap["TIT2"]
    text := "Hello"
    textFrame := NewTextFrame(ft, text)
    mp3File.AddFrames(textFrame)

# CLI Utilities

We include CLI utilities to inspect or modify audio files metadata directly from your terminal.

## ID3 Lister (`id3lister`)

Inspects and outputs all core metadata (Title, Artist, Album, etc.) from an audio file.

### Installation

    go install github.com/unxed/id3-go/cmd/id3lister@latest

### Usage

    id3lister <path_to_audio_file.mp3>

## JSON Converter (`id3json`)

Allows you to export metadata to a structured JSON file, or apply metadata from a JSON file back into an audio file.

### Installation

    go install github.com/unxed/id3-go/cmd/id3json@latest

### Export to JSON

    id3json tojson <path_to_audio_file.mp3> <output.json>

### Import from JSON

    id3json toid3 <input.json> <path_to_audio_file.mp3>
