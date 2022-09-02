package main

import (
	"github.com/Ernyoke/Imger/imgio"
	"github.com/jphastings/postcarder/pkg/images"
)

func check(err error) {
	if err != nil {
		panic(err)
	}
}

func main() {
	img, err := imgio.ImreadRGBA("/Users/jp/Pictures/shutup-madrid-front.jpeg")
	check(err)
	outImg, err := images.RemoveBackground(img)
	check(err)

	check(imgio.Imwrite(outImg, "out.jpg"))
}
