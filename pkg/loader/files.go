package loader

import (
	"fmt"
	"io"
	"os"
	"path"
	"strings"

	"github.com/jphastings/postcarder/pkg/postcards"
)

func QuickLoad(dir, prefix string) (*postcards.Postcard, error) {
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

	return Load(front, back, meta)
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
