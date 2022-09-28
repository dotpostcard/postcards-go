package resolution

import (
	"fmt"

	"github.com/jphastings/postcards-go/internal/types"
)

var readers = []struct {
	magicBytes []byte
	fn         func([]byte) (types.Resolution, error)
}{
	{[]byte(pngHeader), decodePng},          // PNG
	{[]byte("\xff\xd8"), decodeExif},        // JPEG
	{[]byte("RIFF????WEBPVP8"), decodeExif}, // WebP
}

func Decode(data []byte) (types.Resolution, error) {
	for _, r := range readers {
		if isMagic(data[0:len(r.magicBytes)], r.magicBytes) {
			return r.fn(data)
		}
	}
	return types.Resolution{}, fmt.Errorf("unparseable image format")
}

func isMagic(data, magic []byte) bool {
	if len(magic) != len(data) {
		return false
	}

	for i, b := range data {
		if magic[i] != b && magic[i] != '?' {
			return false
		}
	}

	return true
}
