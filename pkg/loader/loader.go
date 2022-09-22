package loader

import (
	"io"

	"github.com/jphastings/postcarder/pkg/postcards"
)

func Load(frontReader, backReader, metaReader io.Reader) (*postcards.Postcard, error) {
	meta, err := readerToMeta(metaReader)
	if err != nil {
		return nil, err
	}
	frontImg, err := readerToImage(frontReader)
	if err != nil {
		return nil, err
	}
	backImg, err := readerToImage(backReader)
	if err != nil {
		return nil, err
	}

	if err := validateDimensions(frontImg, backImg); err != nil {
		return nil, err
	}

	return &postcards.Postcard{
		Front: frontImg,
		Back:  backImg,
		Meta:  meta,
	}, nil
}
