package media_test

import (
	"bytes"
	"encoding/hex"
	"image"
	_ "image/jpeg"
	"media"
	"os"
	"strings"
	"testing"
)

// TestFfmpeg tests that the function works on some level. It does not test
// that the output is correct, just that it produces somewhat valid output.
// TODO: write more comprehensive tests.
func TestFfmpeg(t *testing.T) {
	files, err := os.ReadDir("./test")
	if err != nil {
		t.Error(err)
	}

	for _, file := range files {
		fp := "./test/" + file.Name()
		r, err := os.Open(fp)
		if err != nil {
			t.Error(err)
		}

		w := &bytes.Buffer{}
		err = media.MakeThumbnail(r, w)
		if err != nil {
			t.Error(err)
		}
		img, _, err := image.Decode(w)
		if err != nil {
			t.Error(err)
		}

		hash := media.GeneratePhash(img)
		shouldFail := strings.Contains(file.Name(), "fail")
		if hex.EncodeToString(hash[:]) != expect && !shouldFail {
			t.Errorf("expected %s, got %s", expect, hex.EncodeToString(hash[:]))
		}
	}
}
