package media

import (
	"bytes"
	"errors"
	"io"
	"os"
	"os/exec"
)

var ErrFFMPEGNotFound = errors.New("ffmpeg not found in PATH")
var ffmpegNotFound = false

// MakeThumbnail takes any files input stream and tries to make a thumbnail from it.
// It uses ffmpeg to do this, as it has the best support for all file types,
// including video files.
func MakeThumbnail(input io.Reader, output io.Writer) error {
	args := []string{"-i", "pipe:", "-f", "image2pipe", "-vframes", "1", "-vf", "scale=1024:1024:force_original_aspect_ratio=decrease", "-q:v", "10", "pipe:.jpg"}
	return RunFfmpeg(input, output, args...)
}

func MakeThumb(in io.Reader, outpath string) error {
	thumbBuf := &bytes.Buffer{}
	err := MakeThumbnail(in, thumbBuf)
	if err != nil {
		return err
	}

	// write to outpath

	return os.WriteFile(outpath, thumbBuf.Bytes(), 0777)
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

// Todo: Implement a function which uses image.Image instead of ffmpeg.
// This will require less dependencies, and will likely be faster.
