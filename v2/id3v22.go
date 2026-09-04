// Copyright 2013 Michael Yang. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.
package v2

import (
	"github.com/unxed/id3-go/encodedbytes"
	"io"
)

const (
	V22FrameHeaderSize = 6
)

var (
	// Common frame IDs
	V22CommonFrame = map[string]FrameType{
		"Title":    V22FrameTypeMap["TT2"],
		"Artist":   V22FrameTypeMap["TP1"],
		"Album":    V22FrameTypeMap["TAL"],
		"Year":     V22FrameTypeMap["TYE"],
		"Genre":    V22FrameTypeMap["TCO"],
		"Comments": V22FrameTypeMap["COM"],
	}

	// V22FrameTypeMap specifies the frame IDs and constructors allowed in ID3v2.2
	V22FrameTypeMap = map[string]FrameType{
		"BUF": {id: "BUF", description: "Recommended buffer size", constructor: ParseDataFrame},
		"CNT": {id: "CNT", description: "Play counter", constructor: ParseDataFrame},
		"COM": {id: "COM", description: "Comments", constructor: ParseUnsynchTextFrame},
		"CRA": {id: "CRA", description: "Audio encryption", constructor: ParseDataFrame},
		"CRM": {id: "CRM", description: "Encrypted meta frame", constructor: ParseDataFrame},
		"ETC": {id: "ETC", description: "Event timing codes", constructor: ParseDataFrame},
		"EQU": {id: "EQU", description: "Equalization", constructor: ParseDataFrame},
		"GEO": {id: "GEO", description: "General encapsulated object", constructor: ParseDataFrame},
		"IPL": {id: "IPL", description: "Involved people list", constructor: ParseDataFrame},
		"LNK": {id: "LNK", description: "Linked information", constructor: ParseDataFrame},
		"MCI": {id: "MCI", description: "Music CD Identifier", constructor: ParseDataFrame},
		"MLL": {id: "MLL", description: "MPEG location lookup table", constructor: ParseDataFrame},
		"PIC": {id: "PIC", description: "Attached picture", constructor: ParseDataFrame},
		"POP": {id: "POP", description: "Popularimeter", constructor: ParseDataFrame},
		"REV": {id: "REV", description: "Reverb", constructor: ParseDataFrame},
		"RVA": {id: "RVA", description: "Relative volume adjustment", constructor: ParseDataFrame},
		"SLT": {id: "SLT", description: "Synchronized lyric/text", constructor: ParseDataFrame},
		"STC": {id: "STC", description: "Synced tempo codes", constructor: ParseDataFrame},
		"TAL": {id: "TAL", description: "Album/Movie/Show title", constructor: ParseTextFrame},
		"TBP": {id: "TBP", description: "BPM (Beats Per Minute)", constructor: ParseTextFrame},
		"TCM": {id: "TCM", description: "Composer", constructor: ParseTextFrame},
		"TCO": {id: "TCO", description: "Content type", constructor: ParseTextFrame},
		"TCR": {id: "TCR", description: "Copyright message", constructor: ParseTextFrame},
		"TDA": {id: "TDA", description: "Date", constructor: ParseTextFrame},
		"TDY": {id: "TDY", description: "Playlist delay", constructor: ParseTextFrame},
		"TEN": {id: "TEN", description: "Encoded by", constructor: ParseTextFrame},
		"TFT": {id: "TFT", description: "File type", constructor: ParseTextFrame},
		"TIM": {id: "TIM", description: "Time", constructor: ParseTextFrame},
		"TKE": {id: "TKE", description: "Initial key", constructor: ParseTextFrame},
		"TLA": {id: "TLA", description: "Language(s)", constructor: ParseTextFrame},
		"TLE": {id: "TLE", description: "Length", constructor: ParseTextFrame},
		"TMT": {id: "TMT", description: "Media type", constructor: ParseTextFrame},
		"TOA": {id: "TOA", description: "Original artist(s)/performer(s)", constructor: ParseTextFrame},
		"TOF": {id: "TOF", description: "Original filename", constructor: ParseTextFrame},
		"TOL": {id: "TOL", description: "Original Lyricist(s)/text writer(s)", constructor: ParseTextFrame},
		"TOR": {id: "TOR", description: "Original release year", constructor: ParseTextFrame},
		"TOT": {id: "TOT", description: "Original album/Movie/Show title", constructor: ParseTextFrame},
		"TP1": {id: "TP1", description: "Lead artist(s)/Lead performer(s)/Soloist(s)/Performing group", constructor: ParseTextFrame},
		"TP2": {id: "TP2", description: "Band/Orchestra/Accompaniment", constructor: ParseTextFrame},
		"TP3": {id: "TP3", description: "Conductor/Performer refinement", constructor: ParseTextFrame},
		"TP4": {id: "TP4", description: "Interpreted, remixed, or otherwise modified by", constructor: ParseTextFrame},
		"TPA": {id: "TPA", description: "Part of a set", constructor: ParseTextFrame},
		"TPB": {id: "TPB", description: "Publisher", constructor: ParseTextFrame},
		"TRC": {id: "TRC", description: "ISRC (International Standard Recording Code)", constructor: ParseTextFrame},
		"TRD": {id: "TRD", description: "Recording dates", constructor: ParseTextFrame},
		"TRK": {id: "TRK", description: "Track number/Position in set", constructor: ParseTextFrame},
		"TSI": {id: "TSI", description: "Size", constructor: ParseTextFrame},
		"TSS": {id: "TSS", description: "Software/hardware and settings used for encoding", constructor: ParseTextFrame},
		"TT1": {id: "TT1", description: "Content group description", constructor: ParseTextFrame},
		"TT2": {id: "TT2", description: "Title/Songname/Content description", constructor: ParseTextFrame},
		"TT3": {id: "TT3", description: "Subtitle/Description refinement", constructor: ParseTextFrame},
		"TXT": {id: "TXT", description: "Lyricist/text writer", constructor: ParseTextFrame},
		"TXX": {id: "TXX", description: "User defined text information frame", constructor: ParseDescTextFrame},
		"TYE": {id: "TYE", description: "Year", constructor: ParseTextFrame},
		"UFI": {id: "UFI", description: "Unique file identifier", constructor: ParseDataFrame},
		"ULT": {id: "ULT", description: "Unsychronized lyric/text transcription", constructor: ParseDataFrame},
		"WAF": {id: "WAF", description: "Official audio file webpage", constructor: ParseDataFrame},
		"WAR": {id: "WAR", description: "Official artist/performer webpage", constructor: ParseDataFrame},
		"WAS": {id: "WAS", description: "Official audio source webpage", constructor: ParseDataFrame},
		"WCM": {id: "WCM", description: "Commercial information", constructor: ParseDataFrame},
		"WCP": {id: "WCP", description: "Copyright/Legal information", constructor: ParseDataFrame},
		"WPB": {id: "WPB", description: "Publishers official webpage", constructor: ParseDataFrame},
		"WXX": {id: "WXX", description: "User defined URL link frame", constructor: ParseDataFrame},
	}
)

func ParseV22Frame(reader io.Reader, unsynchronization bool) Framer {
	data := make([]byte, V22FrameHeaderSize)

	if unsynchronization {
		nextOk := false
		for i := 0; i < V22FrameHeaderSize; i++ {
			if n, err := io.ReadFull(reader, data[i:i+1]); n == 0 || err != nil {
				return nil
			}
			if i >= 1 && data[i-1] == 255 && data[i] == 0 && !nextOk {
				// we must skip this 00
				i--
				// but not the next one
				nextOk = true
			} else {
				nextOk = false
			}
		}
	} else {
		if n, err := io.ReadFull(reader, data); n < V22FrameHeaderSize || err != nil {
			return nil
		}
	}

	id := string(data[:3])
	t, ok := V22FrameTypeMap[id]
	if !ok {
		return nil
	}

	size, err := encodedbytes.NormInt(data[3:6])
	if err != nil {
		return nil
	}

	h := FrameHead{
		FrameType: t,
		size:      size,
	}

	frameData := make([]byte, size)

	if unsynchronization {
		nextOk := false
		for i := 0; i < int(size); i++ {
			if n, err := io.ReadFull(reader, frameData[i:i+1]); n == 0 || err != nil {
				return nil
			}
			if i >= 1 && frameData[i-1] == 255 && frameData[i] == 0 && !nextOk {
				// we must skip this 00
				i--
				// but not the next one
				nextOk = true
			} else {
				nextOk = false
			}
		}
	} else {
		if n, err := io.ReadFull(reader, frameData); n < int(size) || err != nil {
			return nil
		}
	}

	return t.constructor(h, frameData)
}

func V22Bytes(f Framer) []byte {
	headBytes := make([]byte, 0, V22FrameHeaderSize)

	headBytes = append(headBytes, f.Id()...)
	headBytes = append(headBytes, encodedbytes.NormBytes(uint32(f.Size()))[1:]...)

	return append(headBytes, f.Bytes()...)
}
