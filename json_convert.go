package id3

import (
	"encoding/json"
	"io/ioutil"
	"os"

	v1 "github.com/unxed/id3-go/v1"
	v2 "github.com/unxed/id3-go/v2"
)

type V1Metadata struct {
	Title   string `json:"title"`
	Artist  string `json:"artist"`
	Album   string `json:"album"`
	Year    string `json:"year"`
	Comment string `json:"comment"`
}

type V2Frame struct {
	ID          string `json:"id"`
	Type        string `json:"type,omitempty"`
	Encoding    string `json:"encoding,omitempty"`
	Language    string `json:"language,omitempty"`
	Description string `json:"description,omitempty"`
	Value       string `json:"value,omitempty"`
	Data        []byte `json:"data,omitempty"`
}

type V2Metadata struct {
	Version string    `json:"version"`
	Padding uint      `json:"padding"`
	Frames  []V2Frame `json:"frames"`
}

type Metadata struct {
	ID3v1 *V1Metadata `json:"id3v1,omitempty"`
	ID3v2 *V2Metadata `json:"id3v2,omitempty"`
}

func ConvertToJSON(mp3Path, jsonPath string) error {
	fi, err := os.Open(mp3Path)
	if err != nil {
		return err
	}
	defer fi.Close()

	var meta Metadata

	// Parse V2
	if v2Tag := v2.ParseTag(fi); v2Tag != nil {
		v2m := &V2Metadata{
			Version: v2Tag.Version(),
			Padding: v2Tag.Padding(),
			Frames:  []V2Frame{},
		}

		for _, frame := range v2Tag.AllFrames() {
			var f V2Frame
			f.ID = frame.Id()
			f.Value = frame.String()

			switch tf := frame.(type) {
			case v2.TextFramer:
				f.Type = "text"
				f.Encoding = tf.Encoding()
				f.Value = tf.Text()
			}

			switch dtf := frame.(type) {
			case *v2.DescTextFrame:
				f.Type = "desc_text"
				f.Description = dtf.Description()
			case *v2.UnsynchTextFrame:
				f.Type = "unsynch_text"
				f.Language = dtf.Language()
			case *v2.IdFrame:
				f.Type = "id"
				f.Description = dtf.OwnerIdentifier()
				f.Data = dtf.Identifier()
			case *v2.DataFrame:
				f.Type = "data"
				f.Data = dtf.Data()
			}

			v2m.Frames = append(v2m.Frames, f)
		}
		meta.ID3v2 = v2m
	}

	// Parse V1
	if _, err := fi.Seek(0, os.SEEK_SET); err == nil {
		if v1Tag := v1.ParseTag(fi); v1Tag != nil {
			meta.ID3v1 = &V1Metadata{
				Title:   v1Tag.Title(),
				Artist:  v1Tag.Artist(),
				Album:   v1Tag.Album(),
				Year:    v1Tag.Year(),
				Comment: v1Tag.Comments()[0],
			}
		}
	}

	data, err := json.MarshalIndent(meta, "", "  ")
	if err != nil {
		return err
	}

	return ioutil.WriteFile(jsonPath, data, 0644)
}

func ConvertToID3(jsonPath, mp3Path string) error {
	data, err := ioutil.ReadFile(jsonPath)
	if err != nil {
		return err
	}

	var meta Metadata
	if err := json.Unmarshal(data, &meta); err != nil {
		return err
	}

	// 1. Process V2
	if meta.ID3v2 != nil {
		fi, err := os.OpenFile(mp3Path, os.O_RDWR, 0666)
		if err != nil {
			return err
		}

		var originalSize int
		if existingV2 := v2.ParseTag(fi); existingV2 != nil {
			originalSize = existingV2.Size()
		}

		tag := v2.NewTag(3)

		for _, jf := range meta.ID3v2.Frames {
			var ft v2.FrameType
			var ok bool
			ft, ok = v2.V23FrameTypeMap[jf.ID]
			if !ok {
				ft = v2.V23FrameTypeMap["TXXX"]
			}

			var framer v2.Framer
			switch jf.Type {
			case "desc_text":
				framer = v2.NewDescTextFrame(ft, jf.Description, jf.Value)
			case "unsynch_text":
				utf := v2.NewUnsynchTextFrame(ft, jf.Description, jf.Value)
				if jf.Language != "" {
					utf.SetLanguage(jf.Language)
				}
				framer = utf
			case "id":
				framer = v2.NewIdFrame(ft, jf.Description, jf.Data)
			case "data":
				framer = v2.NewDataFrame(ft, jf.Data)
			default:
				framer = v2.NewTextFrame(ft, jf.Value)
			}

			if framer != nil {
				if jf.Encoding != "" {
					if tf, ok := framer.(v2.TextFramer); ok {
						tf.SetEncoding(jf.Encoding)
					}
				}
				tag.AddFrames(framer)
			}
		}

		if _, err := fi.Seek(0, os.SEEK_SET); err != nil {
			fi.Close()
			return err
		}

		wrapFile := &File{
			Tagger:       tag,
			originalSize: originalSize,
			file:         fi,
		}

		if err := wrapFile.Close(); err != nil {
			return err
		}
	}

	// 2. Process V1
	if meta.ID3v1 != nil {
		fi, err := os.OpenFile(mp3Path, os.O_RDWR, 0666)
		if err != nil {
			return err
		}
		defer fi.Close()

		tag := &v1.Tag{}
		tag.SetTitle(meta.ID3v1.Title)
		tag.SetArtist(meta.ID3v1.Artist)
		tag.SetAlbum(meta.ID3v1.Album)
		tag.SetYear(meta.ID3v1.Year)
		tag.SetComment(meta.ID3v1.Comment)

		var hasV1 bool
		if _, err := fi.Seek(-v1.TagSize, os.SEEK_END); err == nil {
			sig := make([]byte, 3)
			if _, err := fi.Read(sig); err == nil && string(sig) == "TAG" {
				hasV1 = true
			}
		}

		if hasV1 {
			if _, err := fi.Seek(-v1.TagSize, os.SEEK_END); err != nil {
				return err
			}
		} else {
			if _, err := fi.Seek(0, os.SEEK_END); err != nil {
				return err
			}
		}

		if _, err := fi.Write(tag.Bytes()); err != nil {
			return err
		}
	}

	return nil
}