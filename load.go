package media

import (
	"bytes"
	"errors"
	"image"
	_ "image/gif"
	_ "image/jpeg"
	_ "image/png"
	"os"
	"path/filepath"
	"strings"
)

// OpenImage loads an image from disk.
func OpenImage(filename string) (img image.Image, err error) {
	ext := strings.ToLower(filepath.Ext(filename))
	if !IsMedia(ext) {
		return nil, errors.New("file not recognized as media: " + filename)
	}

	f, err := os.Open(filename)
	if err != nil {
		panic(err)
	}
	defer f.Close()

	switch ext {
	case ".jpg", ".jpeg", ".png", ".gif":
		img, _, err = image.Decode(f)
		if err != nil {
			return nil, err
		}
		return img, nil
	default:
		w := &bytes.Buffer{}

		MakeThumbnail(f, w)

		img, _, err = image.Decode(w)
		if err != nil {
			return nil, err
		}

		return img, nil
	}
}
