package compile

import (
	"bytes"
	"image"

	"github.com/dotpostcard/postcards-go/internal/types"
	"github.com/kolesa-team/go-webp/encoder"
	"github.com/kolesa-team/go-webp/webp"
	"github.com/kolesa-team/goexiv"
)

var webpEncoderOpts *encoder.Options

func init() {
	opts, err := encoder.NewLossyEncoderOptions(encoder.PresetDefault, 85)
	if err != nil {
		panic(err)
	}
	webpEncoderOpts = opts
}

// encodeWebp turns a the image.Image into bytes in Webp format. Currently does *not* write the resolution
// bytes into exif tags, as I can't find a good library for completing this (goexiv doesn't support writing
// rational numbers, which XResolution and YResolution are.)
func encodeWebp(img image.Image, size types.Size) ([]byte, error) {
	data := new(bytes.Buffer)
	if err := webp.Encode(data, img, webpEncoderOpts); err != nil {
		return nil, err
	}

	goIm, err := goexiv.OpenBytes(data.Bytes())
	if err != nil {
		return nil, err
	}

	return goIm.GetBytes(), nil
}
