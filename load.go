package media

import (
	"bytes"
	"errors"
	"fmt"
	"image"
	_ "image/gif"
	_ "image/jpeg"
	_ "image/png"
	"io"
	"os"
	"os/exec"

	"github.com/gabriel-vasile/mimetype"
)

var ffmpegNotFound = false

// supported is a map of supported extensions, those with value true are supported natively
// those with false values depend on ffmpeg.
var supported = map[string]bool{
	".jpg":  true,
	".jpeg": true,
	".jfif": true,
	".png":  true,
	".gif":  true,
	".bmp":  false,
	".jxl":  false,
	".tiff": false,
	".webp": false,
	".mp4":  false,
	".webm": false,
	".mkv":  false,
	".avi":  false,
	".mov":  false,
	".m4v":  false,
}

// OpenImage opens a media file from disk and attempts to decode it into an image.Image.
// If the file is not an image, it will be passed to MakeThumbnail to generate a thumbnail. This function currently only returns thumbnails in case of non StdImage media.
// Either change the function name or make an ffmpeg command which returns first frame of video/decodes image to jpg.
//
// For this package, this functionality with thumbnails is preferred.
func OpenImage(filename string) (img image.Image, err error) {
	f, err := os.Open(filename)
	if err != nil {
		panic(err)
	}
	defer f.Close()

	return ImageFromReader(f)
}

func ImageFromReader(r io.Reader) (img image.Image, err error) {
	b, err := io.ReadAll(r)
	if err != nil {
		return nil, err
	}
	return ImageFromBytes(b)
}

func ImageFromBytes(b []byte) (img image.Image, err error) {
	m := mimetype.Detect(b)
	if m.Extension() == "octet-stream" {
		return nil, ErrUnsupportedMIME
	}

	w := bytes.NewBuffer(b)

	if std, ok := supported[m.Extension()]; ok {
		if std {
			img, _, err = image.Decode(w)
			if err != nil {
				return nil, err
			}
			return img, nil
		}

		tw := &bytes.Buffer{}
		imageWithFFMPEG(w, tw, -1)

		img, _, err = image.Decode(tw)
		if err != nil {
			return nil, err
		}
		return img, nil
	}
	return nil, ErrUnsupportedMIME
}

// Scale is the maximum dimension in both width and height
func imageWithFFMPEG(in io.Reader, out io.Writer, scale int) error {
	args := []string{"-i", "pipe:", "-f", "image2pipe", "-vframes", "1", "-vf"}

	if scale != -1 {
		if scale <= 0 {
			scale = 512
		}
		md := fmt.Sprintf("scale=%v:%v:force_original_aspect_ratio=decrease", scale, scale)
		args = append(args, []string{md, "-q:v", "10", "pipe:.jpg"}...)
	} else {
		args = append(args, []string{"-q:v", "10", "pipe:.jpg"}...)
	}

	return RunFfmpeg(in, out, args...)
}

func RunFfmpeg(r io.Reader, w io.Writer, args ...string) error {
	if ffmpegNotFound {
		return ErrFFMPEGNotFound
	}
	cmd := exec.Command("ffmpeg", args...)
	if r == nil || w == nil {
		return errors.New("invalid reader or writer")
	}
	cmd.Stdin = r
	cmd.Stdout = w
	return cmd.Run()
}

func init() {
	_, err := exec.LookPath("ffmpeg")
	if err != nil {
		ffmpegNotFound = true
	}
}
