package cmdhelp

import (
	"fmt"
	"os"
	"path"
	"path/filepath"
	"strings"

	"github.com/dotpostcard/postcards-go"
	"github.com/dotpostcard/postcards-go/internal/helpers"
	"github.com/dotpostcard/postcards-go/internal/types"
	"github.com/spf13/cobra"
)

type FSInfo struct {
	Dir       string
	Base      string
	SizeHuman string
}

func OpenPostcard(filename string, metadataOnly bool) (*types.Postcard, FSInfo, error) {
	abs, err := filepath.Abs(filename)
	if err != nil {
		return nil, FSInfo{}, err
	}

	f, err := os.Open(abs)
	if err != nil {
		return nil, FSInfo{}, fmt.Errorf("unable to open file: %w", err)
	}
	defer f.Close()

	pc, err := postcards.Read(f, metadataOnly)
	if err != nil {
		return nil, FSInfo{}, fmt.Errorf("unable to read postcard file: %w", err)
	}

	return pc, FSInfo{
		Dir:       path.Dir(abs),
		Base:      strings.TrimSuffix(path.Base(abs), ".postcard"),
		SizeHuman: helpers.SizeHuman(f),
	}, nil
}

func CreateFile(cmd *cobra.Command, info FSInfo, suffix string) (*os.File, string, error) {
	out := path.Join(OutputDir(cmd, info.Dir), info.Base+suffix)
	f, err := os.OpenFile(out, os.O_CREATE|os.O_WRONLY, 0666)
	return f, out, err
}

func WriteFile(cmd *cobra.Command, info FSInfo, suffix string, data []byte, displayName string) error {
	f, out, err := CreateFile(cmd, info, suffix)
	if err != nil {
		return err
	}
	defer f.Close()

	if _, err := f.Write(data); err != nil {
		return err
	}

	if displayName != "" {
		fmt.Printf("Writing %s to %s\n", displayName, out)
	}
	return nil
}

// OutputDir is closely tied to
func OutputDir(cmd *cobra.Command, there string) string {
	if out := must(cmd.Flags().GetString("outdir")); out != "" {
		return out
	}

	if must(cmd.Flags().GetBool("here")) {
		return "."
	}
	if must(cmd.Flags().GetBool("there")) {
		return there
	}

	fmt.Println("Warning: No outdir specified, writing to current working directory")
	return "."
}

func must[T any](val T, err error) T {
	if err != nil {
		panic(err)
	}
	return val
}
