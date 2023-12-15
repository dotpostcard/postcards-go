package compile

import (
	"bytes"
	_ "embed"
	"image"
	"image/draw"
	"io"
	"math/big"
	"text/template"

	"github.com/dotpostcard/postcards-go/internal/types"
)

//go:embed web.md.tmpl
var webTemplateStr string
var webTemplate = template.Must(template.New("web").Parse(webTemplateStr))

func CompileWeb(front, back io.Reader, mp MetadataProvider) (imgBytes []byte, mdBytes []byte, err error) {
	frontImg, backImg, frontDims, _, meta, err := processImages(front, back, mp)
	if err != nil {
		return nil, nil, err
	}

	combinedDims := types.Size{
		PxWidth:  frontDims.PxWidth,
		PxHeight: frontDims.PxHeight * 2,
		CmWidth:  frontDims.CmWidth,
		CmHeight: frontDims.CmHeight.Mul(frontDims.CmHeight, big.NewRat(2, 1)),
	}

	frontBounds := frontImg.Bounds()
	backBounds := image.Rectangle{
		Min: image.Point{0, frontDims.PxHeight},
		Max: image.Point{frontDims.PxWidth, combinedDims.PxHeight},
	}

	combinedImg := image.NewRGBA(image.Rect(0, 0, combinedDims.PxWidth, combinedDims.PxHeight))
	draw.Draw(combinedImg, frontBounds, frontImg, image.Point{}, draw.Src)

	if meta.Flip == types.FlipLeftHand || meta.Flip == types.FlipRightHand {
		backImg = rotateImage(backImg, meta.Flip)
	}
	draw.Draw(combinedImg, backBounds, backImg, image.Point{}, draw.Src)

	combined, err := encodeWebp(combinedImg, combinedDims)
	if err != nil {
		return nil, nil, err
	}

	buf := new(bytes.Buffer)
	if err := webTemplate.Execute(buf, meta); err != nil {
		return nil, nil, err
	}

	return combined, buf.Bytes(), nil
}

func rotateImage(img image.Image, flip types.Flip) image.Image {
	bounds := img.Bounds()
	rotated := image.NewRGBA(image.Rect(0, 0, bounds.Dy(), bounds.Dx()))

	switch flip {
	case types.FlipLeftHand:
		// Top left of the source should be bottom left of the output
		for x := bounds.Min.X; x < bounds.Max.X; x++ {
			for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
				rotated.Set(y, bounds.Max.X-x, img.At(x, y))
			}
		}
	case types.FlipRightHand:
		// Top left of the source should be top right of the output
		for x := bounds.Min.X; x < bounds.Max.X; x++ {
			for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
				rotated.Set(bounds.Max.Y-y, x, img.At(x, y))
			}
		}
	}

	return rotated
}
