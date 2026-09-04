// Copyright 2013 Michael Yang. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.
package v2

import (
	"github.com/unxed/id3-go/encodedbytes"
	"io"
)

var (
	// Common frame IDs
	V23CommonFrame = map[string]FrameType{
		"Title":    V23FrameTypeMap["TIT2"],
		"Artist":   V23FrameTypeMap["TPE1"],
		"Album":    V23FrameTypeMap["TALB"],
		"Year":     V23FrameTypeMap["TYER"],
		"Genre":    V23FrameTypeMap["TCON"],
		"Comments": V23FrameTypeMap["COMM"],
	}

	// V23DeprecatedTypeMap contains deprecated frame IDs from ID3v2.2
	V23DeprecatedTypeMap = map[string]string{
		"BUF": "RBUF", "COM": "COMM", "CRA": "AENC", "EQU": "EQUA",
		"ETC": "ETCO", "GEO": "GEOB", "MCI": "MCDI", "MLL": "MLLT",
		"PIC": "APIC", "POP": "POPM", "REV": "RVRB", "RVA": "RVAD",
		"SLT": "SYLT", "STC": "SYTC", "TAL": "TALB", "TBP": "TBPM",
		"TCM": "TCOM", "TCO": "TCON", "TCR": "TCOP", "TDA": "TDAT",
		"TDY": "TDLY", "TEN": "TENC", "TFT": "TFLT", "TIM": "TIME",
		"TKE": "TKEY", "TLA": "TLAN", "TLE": "TLEN", "TMT": "TMED",
		"TOA": "TOPE", "TOF": "TOFN", "TOL": "TOLY", "TOR": "TORY",
		"TOT": "TOAL", "TP1": "TPE1", "TP2": "TPE2", "TP3": "TPE3",
		"TP4": "TPE4", "TPA": "TPOS", "TPB": "TPUB", "TRC": "TSRC",
		"TRD": "TRDA", "TRK": "TRCK", "TSI": "TSIZ", "TSS": "TSSE",
		"TT1": "TIT1", "TT2": "TIT2", "TT3": "TIT3", "TXT": "TEXT",
		"TXX": "TXXX", "TYE": "TYER", "UFI": "UFID", "ULT": "USLT",
		"WAF": "WOAF", "WAR": "WOAR", "WAS": "WOAS", "WCM": "WCOM",
		"WCP": "WCOP", "WPB": "WPB", "WXX": "WXXX",
	}

	// V23FrameTypeMap specifies the frame IDs and constructors allowed in ID3v2.3
	V23FrameTypeMap = map[string]FrameType{
		"AENC": {id: "AENC", description: "Audio encryption", constructor: ParseDataFrame},
		"APIC": {id: "APIC", description: "Attached picture", constructor: ParseImageFrame},
		"COMM": {id: "COMM", description: "Comments", constructor: ParseUnsynchTextFrame},
		"COMR": {id: "COMR", description: "Commercial frame", constructor: ParseDataFrame},
		"ENCR": {id: "ENCR", description: "Encryption method registration", constructor: ParseDataFrame},
		"EQUA": {id: "EQUA", description: "Equalization", constructor: ParseDataFrame},
		"ETCO": {id: "ETCO", description: "Event timing codes", constructor: ParseDataFrame},
		"GEOB": {id: "GEOB", description: "General encapsulated object", constructor: ParseDataFrame},
		"GRID": {id: "GRID", description: "Group identification registration", constructor: ParseDataFrame},
		"IPLS": {id: "IPLS", description: "Involved people list", constructor: ParseDataFrame},
		"LINK": {id: "LINK", description: "Linked information", constructor: ParseDataFrame},
		"MCDI": {id: "MCDI", description: "Music CD identifier", constructor: ParseDataFrame},
		"MLLT": {id: "MLLT", description: "MPEG location lookup table", constructor: ParseDataFrame},
		"OWNE": {id: "OWNE", description: "Ownership frame", constructor: ParseDataFrame},
		"PRIV": {id: "PRIV", description: "Private frame", constructor: ParseDataFrame},
		"PCNT": {id: "PCNT", description: "Play counter", constructor: ParseDataFrame},
		"POPM": {id: "POPM", description: "Popularimeter", constructor: ParseDataFrame},
		"POSS": {id: "POSS", description: "Position synchronisation frame", constructor: ParseDataFrame},
		"RBUF": {id: "RBUF", description: "Recommended buffer size", constructor: ParseDataFrame},
		"RVAD": {id: "RVAD", description: "Relative volume adjustment", constructor: ParseDataFrame},
		"RVRB": {id: "RVRB", description: "Reverb", constructor: ParseDataFrame},
		"SYLT": {id: "SYLT", description: "Synchronized lyric/text", constructor: ParseDataFrame},
		"SYTC": {id: "SYTC", description: "Synchronized tempo codes", constructor: ParseDataFrame},
		"TALB": {id: "TALB", description: "Album/Movie/Show title", constructor: ParseTextFrame},
		"TBPM": {id: "TBPM", description: "BPM (beats per minute)", constructor: ParseTextFrame},
		"TCOM": {id: "TCOM", description: "Composer", constructor: ParseTextFrame},
		"TCON": {id: "TCON", description: "Content type", constructor: ParseTextFrame},
		"TCOP": {id: "TCOP", description: "Copyright message", constructor: ParseTextFrame},
		"TDAT": {id: "TDAT", description: "Date", constructor: ParseTextFrame},
		"TDLY": {id: "TDLY", description: "Playlist delay", constructor: ParseTextFrame},
		"TENC": {id: "TENC", description: "Encoded by", constructor: ParseTextFrame},
		"TEXT": {id: "TEXT", description: "Lyricist/Text writer", constructor: ParseTextFrame},
		"TFLT": {id: "TFLT", description: "File type", constructor: ParseTextFrame},
		"TIME": {id: "TIME", description: "Time", constructor: ParseTextFrame},
		"TIT1": {id: "TIT1", description: "Content group description", constructor: ParseTextFrame},
		"TIT2": {id: "TIT2", description: "Title/songname/content description", constructor: ParseTextFrame},
		"TIT3": {id: "TIT3", description: "Subtitle/Description refinement", constructor: ParseTextFrame},
		"TKEY": {id: "TKEY", description: "Initial key", constructor: ParseTextFrame},
		"TLAN": {id: "TLAN", description: "Language(s)", constructor: ParseTextFrame},
		"TLEN": {id: "TLEN", description: "Length", constructor: ParseTextFrame},
		"TMED": {id: "TMED", description: "Media type", constructor: ParseTextFrame},
		"TOAL": {id: "TOAL", description: "Original album/movie/show title", constructor: ParseTextFrame},
		"TOFN": {id: "TOFN", description: "Original filename", constructor: ParseTextFrame},
		"TOLY": {id: "TOLY", description: "Original lyricist(s)/text writer(s)", constructor: ParseTextFrame},
		"TOPE": {id: "TOPE", description: "Original artist(s)/performer(s)", constructor: ParseTextFrame},
		"TORY": {id: "TORY", description: "Original release year", constructor: ParseTextFrame},
		"TOWN": {id: "TOWN", description: "File owner/licensee", constructor: ParseTextFrame},
		"TPE1": {id: "TPE1", description: "Lead performer(s)/Soloist(s)", constructor: ParseTextFrame},
		"TPE2": {id: "TPE2", description: "Band/orchestra/accompaniment", constructor: ParseTextFrame},
		"TPE3": {id: "TPE3", description: "Conductor/performer refinement", constructor: ParseTextFrame},
		"TPE4": {id: "TPE4", description: "Interpreted, remixed, or otherwise modified by", constructor: ParseTextFrame},
		"TPOS": {id: "TPOS", description: "Part of a set", constructor: ParseTextFrame},
		"TPUB": {id: "TPUB", description: "Publisher", constructor: ParseTextFrame},
		"TRCK": {id: "TRCK", description: "Track number/Position in set", constructor: ParseTextFrame},
		"TRDA": {id: "TRDA", description: "Recording dates", constructor: ParseTextFrame},
		"TRSN": {id: "TRSN", description: "Internet radio station name", constructor: ParseTextFrame},
		"TRSO": {id: "TRSO", description: "Internet radio station owner", constructor: ParseTextFrame},
		"TSIZ": {id: "TSIZ", description: "Size", constructor: ParseTextFrame},
		"TSRC": {id: "TSRC", description: "ISRC (international standard recording code)", constructor: ParseTextFrame},
		"TSSE": {id: "TSSE", description: "Software/Hardware and settings used for encoding", constructor: ParseTextFrame},
		"TYER": {id: "TYER", description: "Year", constructor: ParseTextFrame},
		"TXXX": {id: "TXXX", description: "User defined text information frame", constructor: ParseDescTextFrame},
		"UFID": {id: "UFID", description: "Unique file identifier", constructor: ParseIdFrame},
		"USER": {id: "USER", description: "Terms of use", constructor: ParseDataFrame},
		"TCMP": {id: "TCMP", description: "Part of a compilation (iTunes extension)", constructor: ParseTextFrame},
		"USLT": {id: "USLT", description: "Unsychronized lyric/text transcription", constructor: ParseUnsynchTextFrame},
		"WCOM": {id: "WCOM", description: "Commercial information", constructor: ParseDataFrame},
		"WCOP": {id: "WCOP", description: "Copyright/Legal information", constructor: ParseDataFrame},
		"WOAF": {id: "WOAF", description: "Official audio file webpage", constructor: ParseDataFrame},
		"WOAR": {id: "WOAR", description: "Official artist/performer webpage", constructor: ParseDataFrame},
		"WOAS": {id: "WOAS", description: "Official audio source webpage", constructor: ParseDataFrame},
		"WORS": {id: "WORS", description: "Official internet radio station homepage", constructor: ParseDataFrame},
		"WPAY": {id: "WPAY", description: "Payment", constructor: ParseDataFrame},
		"WPUB": {id: "WPUB", description: "Publishers official webpage", constructor: ParseDataFrame},
		"WXXX": {id: "WXXX", description: "User defined URL link frame", constructor: ParseDataFrame},
	}
)

func ParseV23Frame(reader io.Reader, unsynchronization bool) Framer {
	data := make([]byte, FrameHeaderSize)

	if unsynchronization {
		nextOk := false
		for i := 0; i < FrameHeaderSize; i++ {
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
		if n, err := io.ReadFull(reader, data); n < FrameHeaderSize || err != nil {
			return nil
		}
	}

	id := string(data[:4])
	t, ok := V23FrameTypeMap[id]
	size, err := encodedbytes.NormInt(data[4:8])

	if !ok {
		return nil
	}

	if err != nil {
		return nil
	}

	h := FrameHead{
		FrameType:   t,
		statusFlags: data[8],
		formatFlags: data[9],
		size:        size,
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

func V23Bytes(f Framer) []byte {
	headBytes := make([]byte, 0, FrameHeaderSize)

	headBytes = append(headBytes, f.Id()...)
	headBytes = append(headBytes, encodedbytes.NormBytes(uint32(f.Size()))...)
	headBytes = append(headBytes, f.StatusFlags(), f.FormatFlags())

	return append(headBytes, f.Bytes()...)
}
