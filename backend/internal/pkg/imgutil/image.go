package imgutil

import (
	"bytes"
	"errors"
	"fmt"
	"image"
	"image/color"
	"image/jpeg"
	_ "image/png"
	"io"
	"math"
)

var (
	ErrUnsupportedFormat = errors.New("unsupported image format: only jpeg and png are supported")
	ErrInvalidDimensions = errors.New("image has invalid or empty dimensions")
)

const (
	DefaultTargetDim = 400
	DefaultQuality   = 85
	ContentTypeJPEG  = "image/jpeg"
)

// ProcessProfilePicture reads an image from reader, validates format (JPEG/PNG),
// crops to a 1:1 square, resizes to targetDim x targetDim, and encodes to clean JPEG with quality.
func ProcessProfilePicture(r io.Reader, targetDim int, quality int) (io.Reader, int64, string, error) {
	const op = "pkg.imgutil.ProcessProfilePicture"

	if targetDim <= 0 {
		targetDim = DefaultTargetDim
	}
	if quality <= 0 || quality > 100 {
		quality = DefaultQuality
	}

	img, format, err := image.Decode(r)
	if err != nil {
		return nil, 0, "", fmt.Errorf("%s: failed to decode image: %w", op, err)
	}

	if format != "jpeg" && format != "png" {
		return nil, 0, "", fmt.Errorf("%s: %w (detected: %s)", op, ErrUnsupportedFormat, format)
	}

	bounds := img.Bounds()
	w, h := bounds.Dx(), bounds.Dy()
	if w <= 0 || h <= 0 {
		return nil, 0, "", fmt.Errorf("%s: %w", op, ErrInvalidDimensions)
	}

	// 1. Calculate center square crop
	minDim := w
	if h < minDim {
		minDim = h
	}

	cropX := bounds.Min.X + (w-minDim)/2
	cropY := bounds.Min.Y + (h-minDim)/2

	// 2. Bilinear resize into target square canvas
	dst := image.NewRGBA(image.Rect(0, 0, targetDim, targetDim))
	scaleFactor := float64(minDim) / float64(targetDim)

	for y := 0; y < targetDim; y++ {
		srcY := float64(cropY) + (float64(y)+0.5)*scaleFactor - 0.5
		y0 := int(math.Floor(srcY))
		y1 := y0 + 1
		yWeight := srcY - float64(y0)
		if y0 < cropY {
			y0 = cropY
		}
		if y1 >= cropY+minDim {
			y1 = cropY + minDim - 1
		}

		for x := 0; x < targetDim; x++ {
			srcX := float64(cropX) + (float64(x)+0.5)*scaleFactor - 0.5
			x0 := int(math.Floor(srcX))
			x1 := x0 + 1
			xWeight := srcX - float64(x0)
			if x0 < cropX {
				x0 = cropX
			}
			if x1 >= cropX+minDim {
				x1 = cropX + minDim - 1
			}

			r00, g00, b00, a00 := img.At(x0, y0).RGBA()
			r10, g10, b10, a10 := img.At(x1, y0).RGBA()
			r01, g01, b01, a01 := img.At(x0, y1).RGBA()
			r11, g11, b11, a11 := img.At(x1, y1).RGBA()

			topR := float64(r00)*(1-xWeight) + float64(r10)*xWeight
			topG := float64(g00)*(1-xWeight) + float64(g10)*xWeight
			topB := float64(b00)*(1-xWeight) + float64(b10)*xWeight
			topA := float64(a00)*(1-xWeight) + float64(a10)*xWeight

			botR := float64(r01)*(1-xWeight) + float64(r11)*xWeight
			botG := float64(g01)*(1-xWeight) + float64(g11)*xWeight
			botB := float64(b01)*(1-xWeight) + float64(b11)*xWeight
			botA := float64(a01)*(1-xWeight) + float64(a11)*xWeight

			rVal := uint8((topR*(1-yWeight) + botR*yWeight) / 257)
			gVal := uint8((topG*(1-yWeight) + botG*yWeight) / 257)
			bVal := uint8((topB*(1-yWeight) + botB*yWeight) / 257)
			aVal := uint8((topA*(1-yWeight) + botA*yWeight) / 257)

			dst.SetRGBA(x, y, color.RGBA{R: rVal, G: gVal, B: bVal, A: aVal})
		}
	}

	// 3. Re-encode as clean sanitized JPEG
	var buf bytes.Buffer
	if err := jpeg.Encode(&buf, dst, &jpeg.Options{Quality: quality}); err != nil {
		return nil, 0, "", fmt.Errorf("%s: failed to encode jpeg: %w", op, err)
	}

	return bytes.NewReader(buf.Bytes()), int64(buf.Len()), ContentTypeJPEG, nil
}
