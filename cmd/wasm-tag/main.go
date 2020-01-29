package main

import (
	"bytes"
	"fmt"
	"syscall/js"

	"github.com/dhowden/tag"
)

func main() {
	c := make(chan struct{}, 0)

	imgRead := js.FuncOf(func(this js.Value, args []js.Value) interface{} {
		var inBuf []uint8
		array := args[0]
		inBuf = make([]uint8, array.Get("byteLength").Int())
		js.CopyBytesToGo(inBuf, array)

		reader := bytes.NewReader(inBuf)
		m, err := tag.ReadFrom(reader)
		if err != nil {
			fmt.Printf("error reading file: %v\n", err)
			return nil
		}

		pic := m.Picture()
		if pic == nil {
			return nil
		}


		dst := js.Global().Get("Uint8Array").New(len(pic.Data))
		js.CopyBytesToJS(dst, pic.Data)

		out := map[string]interface{}{
			"ext": pic.Ext,
			"mime_type": pic.MIMEType,
			"type": pic.Type,
			"description": pic.Description,
			"data": dst,
		}

		return js.ValueOf(out)
	})

	val := js.FuncOf(func(this js.Value, args []js.Value) interface{} {
		var inBuf []uint8
		array := args[0]
		inBuf = make([]uint8, array.Get("byteLength").Int())
		js.CopyBytesToGo(inBuf, array)

		reader := bytes.NewReader(inBuf)
		m, err := tag.ReadFrom(reader)
		if err != nil {
			fmt.Printf("error reading file: %v\n", err)
			return nil
		}

		discNo, discTotal := m.Disc()
		trackNo, trackTotal := m.Track()
		out := map[string]interface{}{
			"format":       string(m.Format()),
			"file_type":    string(m.FileType()),
			"title":        m.Title(),
			"album":        m.Album(),
			"artist":       m.Artist(),
			"album_artist": m.AlbumArtist(),
			"composer":     m.Composer(),
			"genre":        m.Genre(),
			"year":         m.Year(),
			"track": map[string]interface{}{
				"number": trackNo,
				"total":  trackTotal,
			},
			"disc": map[string]interface{}{
				"number": discNo,
				"total":  discTotal,
			},
			"lyrics":  m.Lyrics(),
			"comment": m.Comment(),
		}

		return js.ValueOf(out)
	})

	js.Global().Set("loadTags", val)
	js.Global().Set("loadImage", imgRead)

	<-c
}
