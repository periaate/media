package media_test

import (
	"encoding/hex"
	"media"
	"os"
	"strings"
	"testing"
)

const (
	// Expected hash as hex string of all non-fail images in ./test
	expect = "3e7c3dbc3bdc37ec2ff41ff8c00340024002c0031ff82ff437ec3bdc3dbc3e7c"
)

// test hashes of all files in ./test
func TestPerceptualHash(t *testing.T) {
	files, err := os.ReadDir("./test")
	if err != nil {
		t.Error(err)
	}

	for _, file := range files {
		fp := "./test/" + file.Name()
		img, err := media.OpenImage(fp)
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

// test hashes of all files in ./test
func TestMakePhash(t *testing.T) {
	files, err := os.ReadDir("./test")
	if err != nil {
		t.Error(err)
	}

	for _, file := range files {
		fp := "./test/" + file.Name()

		hash := media.MakePhash(fp)
		shouldFail := strings.Contains(file.Name(), "fail")
		if hex.EncodeToString(hash[:]) != expect && !shouldFail {
			t.Errorf("expected %s, got %s", expect, hex.EncodeToString(hash[:]))
		}
	}
}
