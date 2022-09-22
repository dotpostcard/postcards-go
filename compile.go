package postcarder

import (
	"fmt"
	"io"
	"os"
	"path"
	"strings"

	"github.com/jphastings/postcarder/internal/compile"
	"gopkg.in/yaml.v3"
)

func CompileFiles(dir, prefix string) (*Postcard, error) {
	meta, err := tryLoad(dir, prefix, "meta", "yml", "yaml")
	if err != nil {
		return nil, fmt.Errorf("couldn't load metadata: %w", err)
	}
	front, err := tryLoad(dir, prefix, "front", "png", "jpg", "tif", "tiff")
	if err != nil {
		return nil, fmt.Errorf("couldn't load postcard front: %w", err)
	}
	back, err := tryLoad(dir, prefix, "back", "png", "jpg", "tif", "tiff")
	if err != nil {
		return nil, fmt.Errorf("couldn't load postcard back: %w", err)
	}

	return Compile(front, back, meta)
}

func Compile(frontReader, backReader, metaReader io.Reader) (*Postcard, error) {
	var meta PostcardMetadata
	if err := yaml.NewDecoder(metaReader).Decode(&meta); err != nil {
		return nil, err
	}

	frontImg, err := compile.ReaderToImage(frontReader)
	if err != nil {
		return nil, err
	}
	backImg, err := compile.ReaderToImage(backReader)
	if err != nil {
		return nil, err
	}

	if err := compile.ValidateDimensions(frontImg, backImg); err != nil {
		return nil, err
	}

	return &Postcard{
		Front: frontImg,
		Back:  backImg,
		Meta:  meta,
	}, nil
}

func tryLoad(dir, prefix, suffix string, extensions ...string) (io.Reader, error) {
	for _, ext := range extensions {
		r, err := os.Open(path.Join(dir, fmt.Sprintf("%s-%s.%s", prefix, suffix, ext)))
		if err == nil {
			return r, nil
		}
	}
	return nil, fmt.Errorf("no file '%s-%s.{%s}' in %s", prefix, suffix, strings.Join(extensions, ","), dir)
}
