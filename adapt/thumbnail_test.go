package adapt_test

import (
	"bytes"
	"crypto/md5"
	"fmt"
	"image/png"

	"github.com/dotpostcard/postcards-go/adapt"
)

func ExampleThumbnailFile() {
	thumb, err := adapt.ThumbnailFile("../fixtures/hello.postcard", 128, 0)
	if err != nil {
		panic(err)
	}

	buf := new(bytes.Buffer)
	if err := png.Encode(buf, thumb); err != nil {
		panic(err)
	}

	fmt.Printf("Thumbnail PNG (%dx%d) has checksum %x",
		thumb.Bounds().Dx(), thumb.Bounds().Dy(), md5.Sum(buf.Bytes()))
	// Output: Thumbnail PNG (128x157) has checksum eba1acbde60137aff808c9783c05706d
}
