package images

import (
	"github.com/Ernyoke/Imger/threshold"
	"image"
	"image/draw"
)

func RemoveBackground(scan image.Image) (image.Image, error) {
	gray := image.NewGray(image.Rect(0, 0, scan.Bounds().Dx(), scan.Bounds().Dy()))
	draw.Draw(gray, scan.Bounds(), scan, scan.Bounds().Min, draw.Src)

	thresh, err := threshold.Threshold(gray, 75, threshold.ThreshBinary)
	if err != nil {
		return nil, err
	}

	return thresh, nil
}

func Trim(front, back image.RGBA) (image.RGBA, image.RGBA, error) {
	panic("not implemented")
}
