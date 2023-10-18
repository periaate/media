package media

import (
	"image"
	"image/color"
	"io"

	"golang.org/x/image/draw"
)

type Phasher struct {
	Scale int
}

func SetPhashScale(scale int)  { pHasher.Scale = scale }
func DefaultPhasher() *Phasher { return pHasher }

var pHasher *Phasher

func Phash(r io.Reader) (b []byte, err error) {
	img, err := ImageFromReader(r)
	if err != nil {
		return b, err
	}

	return pHasher.GeneratePhash(img), nil
}

// GeneratePhash returns a basic perceptual hash of the given image.
func (p Phasher) GeneratePhash(img image.Image) []byte {
	resizedImg := ResizeImage(img, p.Scale, p.Scale)

	grayImg := toGrayscale(resizedImg)
	avgLight := averageBrightness(grayImg)

	// width * height
	hash := make([]byte, p.Scale*p.Scale/8)
	for i, pixel := range grayImg.Pix {
		if pixel >= avgLight {
			hash[i/8] |= 1 << uint(i%8)
		} else {
			hash[i/8] &= ^(1 << uint(i%8))
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

func MakePhash(filename string) (b []byte, err error) {
	img, err := OpenImage(filename)
	if err != nil {
		return b, err
	}

	return pHasher.GeneratePhash(img), nil
}

func init() {
	pHasher = &Phasher{
		Scale: 16,
	}
}
