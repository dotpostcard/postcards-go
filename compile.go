package postcard

import (
	"fmt"
	"io"
	"os"
	"path"
	"path/filepath"
	"strings"

	"github.com/jphastings/postcard-go/internal/compile"
	"github.com/jphastings/postcard-go/internal/types"
	"gopkg.in/yaml.v3"
)

// CompileFiles accepts a path to one of the three needed files, attempts to find the others, and writes a compiled postcard file using conventional naming.
func CompileFiles(oneFile string) error {
	dir := filepath.Dir(oneFile)
	prefix := strings.SplitN(filepath.Base(oneFile), "-", 2)[0]

	meta, err := tryLoad(dir, prefix, "meta", "yml", "yaml")
	if err != nil {
		return fmt.Errorf("couldn't load metadata: %w", err)
	}
	front, err := tryLoad(dir, prefix, "front", "png", "jpg", "tif", "tiff")
	if err != nil {
		return fmt.Errorf("couldn't load postcard front: %w", err)
	}
	back, err := tryLoad(dir, prefix, "back", "png", "jpg", "tif", "tiff")
	if err != nil {
		return fmt.Errorf("couldn't load postcard back: %w", err)
	}

	pc, err := Compile(front, back, meta)
	if err != nil {
		return err
	}

	return WriteFile(pc, fmt.Sprintf("%s.postcard", prefix))
}

// Compile accepts reader objects for each of the components of a postcard file, and creates an in-memory Postcard object.
func Compile(frontReader, backReader, metaReader io.Reader) (*types.Postcard, error) {
	var meta types.Metadata
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

	return &types.Postcard{
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
