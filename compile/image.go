package compile

import (
	"bytes"
	"image"
	_ "image/jpeg"
	_ "image/png"
	"io"

	"github.com/dotpostcard/postcards-go/internal/types"
	"github.com/dotpostcard/postcards-go/pkg/resolution"
)

func readerToImage(r io.Reader) (image.Image, types.Size, error) {
	buf := new(bytes.Buffer)
	imgR := io.TeeReader(r, buf)

	img, _, err := image.Decode(imgR)
	if err != nil {
		return nil, types.Size{}, err
	}

	size := types.Size{
		PxWidth:  img.Bounds().Dx(),
		PxHeight: img.Bounds().Dy(),
	}

	if xRes, yRes, err := resolution.Decode(buf.Bytes()); err == nil {
		size.SetResolution(xRes, yRes)
	}

	return img, size, nil
}
