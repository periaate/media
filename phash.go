package media

import (
	"image"
	"image/color"
	"io"
	"path/filepath"
	"strings"

	"golang.org/x/image/draw"
)

const (
	resizeWidth  = 16
	resizeHeight = 16
)

// GeneratePhash returns a basic perceptual hash of the given image.
func GeneratePhash(img image.Image) [32]byte {
	resizedImg := ResizeImage(img, resizeWidth, resizeHeight)

	grayImg := toGrayscale(resizedImg)
	avgLight := averageBrightness(grayImg)

	hash := [32]byte{}
	for pixel := range grayImg.Pix {
		if grayImg.Pix[pixel] >= avgLight {
			hash[pixel/8] |= 1 << uint(pixel%8)
		} else {
			hash[pixel/8] &= ^(1 << uint(pixel%8))
		}
	}

	return hash
}

// ResizeImage resizes the given image to the specified width and height.
func ResizeImage(img image.Image, width, height int) image.Image {
	dst := image.NewRGBA(image.Rect(0, 0, width, height))
	draw.NearestNeighbor.Scale(dst, dst.Rect, img, img.Bounds(), draw.Over, nil)
	return dst
}

// toGrayscale converts the given image to grayscale.
func toGrayscale(img image.Image) *image.Gray {
	bounds := img.Bounds()
	grayImg := image.NewGray(bounds)

	for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
		for x := bounds.Min.X; x < bounds.Max.X; x++ {
			grayImg.Set(x, y, color.GrayModel.Convert(img.At(x, y)))
		}
	}

	return grayImg
}

// averageBrightness returns the average brightness of the given grayscale image.
func averageBrightness(img *image.Gray) uint8 {
	var total float64
	bounds := img.Bounds()
	area := bounds.Dx() * bounds.Dy()
	for pixel := range img.Pix {
		total += float64(pixel) / float64(area)
	}

	return uint8(total)
}

// LoadImage loads an image from disk.
func DecodeToImage(r io.Reader) image.Image {
	img, _, err := image.Decode(r)
	if err != nil {
		panic(err)
	}

	return img
}

func IsHashable(path string) bool {
	ext := strings.ToLower(filepath.Ext(path))
	switch ext {
	case ".jpg", ".jpeg", ".png", ".gif":
		return true
	}
	return false
}

func MakePhash(filename string) []byte {
	img, err := OpenImage(filename)
	if err != nil {
		panic(err)
	}

	hash := GeneratePhash(img)
	return hash[:]
}
