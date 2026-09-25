// Copyright 2013 Michael Yang. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.
package v1

import (
	"io"
	"os"

	v2 "github.com/mikkyang/id3-go/v2"
	"github.com/unxed/localecp"
	"golang.org/x/text/encoding"
)

const (
	TagSize = 128
)

var (
	Genres = []string{
		"Blues", "Classic Rock", "Country", "Dance",
		"Disco", "Funk", "Grunge", "Hip-Hop",
		"Jazz", "Metal", "New Age", "Oldies",
		"Other", "Pop", "R&B", "Rap",
		"Reggae", "Rock", "Techno", "Industrial",
		"Alternative", "Ska", "Death Metal", "Pranks",
		"Soundtrack", "Euro-Techno", "Ambient", "Trip-Hop",
		"Vocal", "Jazz+Funk", "Fusion", "Trance",
		"Classical", "Instrumental", "Acid", "House",
		"Game", "Sound Clip", "Gospel", "Noise",
		"AlternRock", "Bass", "Soul", "Punk",
		"Space", "Meditative", "Instrumental Pop", "Instrumental Rock",
		"Ethnic", "Gothic", "Darkwave", "Techno-Industrial",
		"Electronic", "Pop-Folk", "Eurodance", "Dream",
		"Southern Rock", "Comedy", "Cult", "Gangsta",
		"Top 40", "Christian Rap", "Pop/Funk", "Jungle",
		"Native American", "Cabaret", "New Wave", "Psychadelic",
		"Rave", "Showtunes", "Trailer", "Lo-Fi",
		"Tribal", "Acid Punk", "Acid Jazz", "Polka",
		"Retro", "Musical", "Rock & Roll", "Hard Rock",
	}
)

// Encoding is the character set used for the text fields of ID3v1 tags.
// ID3v1 does not define one; in practice taggers wrote the legacy code page
// of the system they ran on (Windows-1251 in Russia, Windows-1250 in Poland,
// and so on).
//
// When Encoding is nil (the default), fields are returned exactly as stored
// and strings are written as their UTF-8 bytes, as before. When it is set,
// fields are decoded from it on parse and encoded to it on write; a string
// that cannot be encoded is written as UTF-8 bytes.
//
// Set it once, before parsing or writing tags; it is not synchronized.
var Encoding encoding.Encoding

// UseLocaleEncoding sets Encoding to the legacy ANSI code page of the
// current system locale (for example Windows-1251 for ru_RU), as detected
// by github.com/unxed/localecp: GetACP on Windows, LC_ALL/LC_CTYPE/LANG
// elsewhere, Windows-1252 if nothing more specific is found.
func UseLocaleEncoding() {
	Encoding = localecp.ANSIEncoding
}

func decodeField(b []byte) string {
	if Encoding == nil {
		return string(b)
	}
	decoded, err := Encoding.NewDecoder().Bytes(b)
	if err != nil {
		return string(b)
	}
	return string(decoded)
}

func encodeField(dst []byte, s string) {
	if Encoding != nil {
		if encoded, err := Encoding.NewEncoder().Bytes([]byte(s)); err == nil {
			copy(dst, encoded)
			return
		}
	}
	copy(dst, []byte(s))
}

// Tag represents an ID3v1 tag
type Tag struct {
	title, artist, album, year, comment string
	genre                               byte
	dirty                               bool
}

func ParseTag(readSeeker io.ReadSeeker) *Tag {
	readSeeker.Seek(-TagSize, os.SEEK_END)

	data := make([]byte, TagSize)
	n, err := io.ReadFull(readSeeker, data)
	if n < TagSize || err != nil || string(data[:3]) != "TAG" {
		return nil
	}

	return &Tag{
		title:   decodeField(data[3:33]),
		artist:  decodeField(data[33:63]),
		album:   decodeField(data[63:93]),
		year:    decodeField(data[93:97]),
		comment: decodeField(data[97:127]),
		genre:   data[127],
		dirty:   false,
	}
}

func (t Tag) Dirty() bool {
	return t.dirty
}

func (t Tag) Title() string  { return t.title }
func (t Tag) Artist() string { return t.artist }
func (t Tag) Album() string  { return t.album }
func (t Tag) Year() string   { return t.year }

func (t Tag) Genre() string {
	if int(t.genre) < len(Genres) {
		return Genres[t.genre]
	}

	return ""
}

func (t Tag) Comments() []string {
	return []string{t.comment}
}

func (t *Tag) SetTitle(text string) {
	t.title = text
	t.dirty = true
}

func (t *Tag) SetArtist(text string) {
	t.artist = text
	t.dirty = true
}

func (t *Tag) SetAlbum(text string) {
	t.album = text
	t.dirty = true
}

func (t *Tag) SetYear(text string) {
	t.year = text
	t.dirty = true
}

func (t *Tag) SetGenre(text string) {
	t.genre = 255
	for i, genre := range Genres {
		if text == genre {
			t.genre = byte(i)
			break
		}
	}
	t.dirty = true
}

func (t Tag) Bytes() []byte {
	data := make([]byte, TagSize)

	copy(data[:3], []byte("TAG"))
	encodeField(data[3:33], t.title)
	encodeField(data[33:63], t.artist)
	encodeField(data[63:93], t.album)
	encodeField(data[93:97], t.year)
	encodeField(data[97:127], t.comment)
	data[127] = t.genre

	return data
}

func (t Tag) Size() int {
	return TagSize
}

func (t Tag) Version() string {
	return "1.0"
}

// Dummy methods to satisfy Tagger interface
func (t Tag) Padding() uint                      { return 0 }
func (t Tag) AllFrames() []v2.Framer             { return []v2.Framer{} }
func (t Tag) Frame(id string) v2.Framer          { return nil }
func (t Tag) Frames(id string) []v2.Framer       { return []v2.Framer{} }
func (t Tag) DeleteFrames(id string) []v2.Framer { return []v2.Framer{} }
func (t Tag) AddFrames(f ...v2.Framer)           {}
