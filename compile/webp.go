package compile

import (
	"bytes"
	"image"

	"github.com/chai2010/webp"
	"github.com/dotpostcard/postcards-go/internal/types"
)

// encodeWebp turns a the image.Image into bytes in Webp format. Currently does *not* write the resolution
// bytes into exif tags, as I can't find a good library for completing this (goexiv doesn't support writing
// rational numbers, which XResolution and YResolution are.)
func encodeWebp(img image.Image, size types.Size) ([]byte, error) {
	data := new(bytes.Buffer)
	if err := webp.Encode(data, img, &webp.Options{Lossless: false, Quality: 75}); err != nil {
		return nil, err
	}

	return data.Bytes(), nil
}
