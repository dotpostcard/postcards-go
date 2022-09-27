package compile

import (
	"bytes"
	"image"
	_ "image/jpeg"
	_ "image/png"
	"io"

	"github.com/jphastings/postcard-go/internal/types"
	"github.com/jphastings/postcard-go/pkg/resolution"
)

func readerToImage(r io.Reader) (image.Image, types.Size, error) {
	buf := new(bytes.Buffer)
	imgR := io.TeeReader(r, buf)

	img, _, err := image.Decode(imgR)
	if err != nil {
		return nil, types.Size{}, err
	}

	res, err := resolution.Decode(buf.Bytes())
	if err != nil {
		return nil, types.Size{}, err
	}

	size := res.Size(img.Bounds().Dx(), img.Bounds().Dy())
	return img, size, nil
}
