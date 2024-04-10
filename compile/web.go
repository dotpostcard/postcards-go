package compile

import (
	"bytes"
	_ "embed"
	"encoding/json"
	"image"
	"io"
	"math/big"

	"github.com/dotpostcard/postcards-go/internal/types"
	"golang.org/x/image/draw"
)

const maxSize = 2048

func CompileWeb(front, back io.Reader, mp MetadataProvider) (imgBytes []byte, jsonBytes []byte, err error) {
	frontImg, backImg, frontDims, _, meta, err := processImages(front, back, mp)
	if err != nil {
		return nil, nil, err
	}

	scaleW := float64(frontDims.PxWidth) / float64(maxSize)
	scaleH := float64(frontDims.PxHeight) / float64(maxSize)
	scale := float64(1)
	if scaleW > scale {
		scale = scaleW
	}
	if scaleH > scale {
		scale = scaleH
	}

	newW := int(float64(frontDims.PxWidth) / scale)
	newH := int(float64(frontDims.PxHeight) / scale)

	combinedDims := types.Size{
		PxWidth:  newW,
		PxHeight: newH * 2,
		CmWidth:  frontDims.CmWidth,
		CmHeight: frontDims.CmHeight.Mul(frontDims.CmHeight, big.NewRat(2, 1)),
	}

	frontBounds := image.Rectangle{
		Min: image.Point{0, 0},
		Max: image.Point{newW, newH},
	}
	backBounds := image.Rectangle{
		Min: image.Point{0, newH},
		Max: image.Point{newW, newH * 2},
	}

	combinedImg := image.NewRGBA(image.Rect(0, 0, combinedDims.PxWidth, combinedDims.PxHeight))
	draw.CatmullRom.Scale(combinedImg, frontBounds, frontImg, frontImg.Bounds(), draw.Src, nil)

	if meta.Flip == types.FlipLeftHand || meta.Flip == types.FlipRightHand {
		backImg = rotateImage(backImg, meta.Flip)
	}
	draw.CatmullRom.Scale(combinedImg, backBounds, backImg, backImg.Bounds(), draw.Src, nil)

	combined, err := encodeWebp(combinedImg, combinedDims)
	if err != nil {
		return nil, nil, err
	}

	buf := new(bytes.Buffer)
	if err := json.NewEncoder(buf).Encode(meta); err != nil {
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
