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

// OpenImage opens a media file from disk and attempts to decode it into an image.Image.
// If the file is not an image, it will be passed to MakeThumbnail to generate a thumbnail. This function currently only returns thumbnails in case of non StdImage media.
// Either change the function name or make an ffmpeg command which returns first frame of video/decodes image to jpg.
//
// For this package, this functionality with thumbnails is preferred.
func OpenImage(filename string) (img image.Image, err error) {
	ext := strings.ToLower(filepath.Ext(filename))
	mt := MediaType(ext)
	if mt == Notmedia {
		return nil, errors.New("file not recognized as media: " + filename)
	}
	if mt == Audio {
		return nil, errors.New("audio files are not supported: " + filename)
	}

	f, err := os.Open(filename)
	if err != nil {
		panic(err)
	}
	defer f.Close()

	if mt == StdImage {
		img, _, err = image.Decode(f)
		if err != nil {
			return nil, err
		}
		return img, nil
	}

	w := &bytes.Buffer{}
	MakeThumbnail(f, w)

	img, _, err = image.Decode(w)
	if err != nil {
		return nil, err
	}

	return img, nil

}
