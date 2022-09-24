package postcard

import (
	"bytes"
	"image"

	"github.com/jphastings/postcard-go/internal/types"
	"golang.org/x/image/draw"
	_ "golang.org/x/image/webp"
)

// Thumbnail generates a thumbnail image representing this postcard's front. Providing a maximum width or height of zero will allow
// the thumbnail to grow up to the postcard's original size in that dimension.
func Thumbnail(pc *types.Postcard, maxWidth, maxHeight int) (image.Image, error) {
	img, _, err := image.Decode(bytes.NewReader(pc.Front))
	if err != nil {
		return nil, err
	}

	thumb := image.NewNRGBA(thumbDimensions(img.Bounds(), maxWidth, maxHeight))
	draw.BiLinear.Scale(thumb, thumb.Bounds(), img, img.Bounds(), draw.Over, nil)

	return thumb, nil
}

// ThumbnailFile is a convenience method for retrieving a thumbnail directly on a file using Thumbnail.
func ThumbnailFile(path string, maxWidth, maxHeight int) (image.Image, error) {
	pc, err := ReadFile(path, false)
	if err != nil {
		return nil, err
	}

	return Thumbnail(pc, maxWidth, maxHeight)
}

func thumbDimensions(b image.Rectangle, maxWidth, maxHeight int) image.Rectangle {
	if maxWidth == 0 {
		maxWidth = b.Dx()
	}
	if maxHeight == 0 {
		maxHeight = b.Dy()
	}

	newWidth := maxWidth
	newHeight := b.Dy() * maxWidth / b.Dx()

	if newHeight > maxHeight {
		newHeight = maxHeight
		newWidth = b.Dx() * maxHeight / b.Dy()
	}

	return image.Rect(0, 0, newWidth, newHeight)
}
