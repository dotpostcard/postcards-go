package images

import (
	"fmt"
	"image"
	"image/color"

	"github.com/disintegration/gift"
)

const (
	cutoff = 0.85
)

func thresh(val float32) float32 {
	if val < cutoff {
		return 0
	}

	return 1
}

func Edges(src *image.RGBA) (image.Image, error) {
	g := gift.New(
		gift.ColorFunc(
			func(r0, g0, b0, a0 float32) (r, g, b, a float32) {
				v := float32(0)
				tr := thresh(r0)
				tg := thresh(g0)
				tb := thresh(b0)

				if tr < cutoff || tg < cutoff || tb < cutoff {
					v = 1
				}

				r = v
				g = v
				b = v
				a = a0
				return r, g, b, a
			},
		),
		gift.Threshold(0.9),
	)
	threshImg := image.NewRGBA(g.Bounds(src.Bounds()))
	g.Draw(threshImg, src)

	dst := image.NewRGBA(g.Bounds(src.Bounds()))

	for y := src.Bounds().Min.Y; y < src.Bounds().Dy(); y++ {
		for x := src.Bounds().Min.X; x < src.Bounds().Dx(); x++ {
			c := shouldBe(src, threshImg, dst, x, y)
			dst.Set(x, y, c)
		}
	}

	return dst, nil
}

var bg = color.RGBA{255, 255, 255, 0}

func shouldBe(src, threshImg, dst image.Image, x, y int) color.Color {
	if x == 0 || y == 0 {
		return bg
	}

	if x == 1 && y == 1 {
		fmt.Println(threshImg.At(x, y).RGBA())
		fmt.Println(isWhite(threshImg.At(x, y)))
		fmt.Println(isBG(dst.At(x-1, y)))
	}

	if isWhite(threshImg.At(x, y)) {
		return src.At(x, y)
	}

	if isBG(dst.At(x-1, y)) || isBG(dst.At(x, y-1)) {
		return bg
	}

	return src.At(x, y)
}

func isWhite(c color.Color) bool {
	r, _, _, _ := c.RGBA()

	return r >= 65535
}

func isBG(c color.Color) bool {
	_, _, _, a := c.RGBA()

	return a == 0
}
