package media

import (
	"image"
	"image/jpeg"
	"io"

	"golang.org/x/image/draw"
)

func MakeThumb(in io.Reader, out io.Writer, maxDimension int) error {
	img, err := ImageFromReader(in)
	if err != nil {
		return err
	}
	if 0 >= maxDimension {
		maxDimension = 64
	}
	return ReduceImageSize(img, out, maxDimension)
}

func ReduceImageSize(img image.Image, w io.Writer, maxDimension int) error {
	// Calculate new dimensions while maintaining aspect ratio
	srcWidth := img.Bounds().Dx()
	srcHeight := img.Bounds().Dy()
	newWidth, newHeight := srcWidth, srcHeight

	if srcWidth > maxDimension || srcHeight > maxDimension {
		if srcWidth > srcHeight {
			newWidth = maxDimension
			newHeight = (srcHeight * maxDimension) / srcWidth
		} else {
			newHeight = maxDimension
			newWidth = (srcWidth * maxDimension) / srcHeight
		}
	}

	// Create a new image with the calculated dimensions
	dst := image.NewRGBA(image.Rect(0, 0, newWidth, newHeight))

	// Scale the image using the ApproxBiLinear interpolation algorithm
	draw.ApproxBiLinear.Scale(dst, dst.Bounds(), img, img.Bounds(), draw.Src, nil)

	// Encode the image as JPEG and write to the writer
	return jpeg.Encode(w, dst, &jpeg.Options{Quality: 50})
}
