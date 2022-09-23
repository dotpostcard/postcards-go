package compile

import (
	"image"

	"github.com/fogleman/gg"
	"github.com/jphastings/postcard-go/internal/types"
	"golang.org/x/image/draw"
)

var pixelatedSize = 0.5 // centimeters

func hideSecrets(img *image.NRGBA, dim *types.Dimensions, secrets []types.Polygon) error {
	if len(secrets) == 0 {
		return nil
	}

	obscured := makeObscuredImage(img, dim)
	mask := makeMask(obscured, secrets)

	draw.Copy(img, image.Point{}, mask, mask.Bounds(), draw.Over, nil)

	return nil
}

func makeObscuredImage(img *image.NRGBA, dim *types.Dimensions) *image.NRGBA {
	cmW, cmH := dim.AsFloats()

	micro := image.NewNRGBA(image.Rect(0, 0, int(cmW/pixelatedSize), int(cmH/pixelatedSize)))
	draw.NearestNeighbor.Scale(micro, micro.Rect, img, img.Bounds(), draw.Over, nil)

	obscured := image.NewNRGBA(img.Bounds())
	draw.CatmullRom.Scale(obscured, obscured.Rect, micro, micro.Bounds(), draw.Over, nil)

	return obscured
}

func makeMask(obscured *image.NRGBA, secrets []types.Polygon) image.Image {
	w, h := obscured.Bounds().Dx(), obscured.Bounds().Dy()

	dc := gg.NewContext(w, h)

	for _, pts := range secrets {
		x, y := pts[0].ToPixels(w, h)
		dc.MoveTo(x, y)
		for _, p := range pts[1:] {
			x, y := p.ToPixels(w, h)
			dc.LineTo(x, y)
		}

		dc.Clip()
	}
	dc.DrawImage(obscured, 0, 0)

	return dc.Image()
}
