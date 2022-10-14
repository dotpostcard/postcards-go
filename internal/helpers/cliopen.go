package helpers

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"path"

	"github.com/dotpostcard/postcards-go"
	"github.com/dotpostcard/postcards-go/internal/types"
)

func OpenFromFilename(
	filename, fallbackName string,
	metaOnly bool,
) (pc *types.Postcard, name string, size string, err error) {
	var src io.Reader

	if filename == "-" {
		name = fallbackName
		src = bufio.NewReader(os.Stdin)
	} else {
		name = path.Base(filename)

		f, err := os.Open(filename)
		if err != nil {
			return nil, "", "", fmt.Errorf("unable to open file: %w", err)
		}
		defer f.Close()

		src = f
		size = FileSizeHuman(f)
	}

	pc, err = postcards.Read(src, metaOnly)
	if err != nil {
		return nil, "", "", fmt.Errorf("unable to read postcard file: %w", err)
	}

	return pc, name, size, nil
}
