package postcard

import (
	"bytes"
	"image"

	"github.com/jphastings/postcard-go/internal/types"
	"golang.org/x/image/draw"
)

func Thumbnail(pc *types.Postcard, maxWidth, maxHeight int) (image.Image, error) {
	img, _, err := image.Decode(bytes.NewReader(pc.Front))
	if err != nil {
		return nil, err
	}

	ratio := maxWidth / img.Bounds().Dx()
	newWidth := maxWidth
	newHeight := img.Bounds().Dy() * ratio

	if newHeight > maxHeight {
		newHeight = maxHeight
		newWidth = img.Bounds().Dx() * maxHeight / img.Bounds().Dy()
	}

	thumb := image.NewNRGBA(image.Rect(0, 0, newWidth, newHeight))
	draw.BiLinear.Scale(thumb, thumb.Bounds(), img, img.Bounds(), draw.Over, nil)

	return thumb, nil
}
