package cmd

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/jphastings/postcarder/pkg/compiler"
	"github.com/jphastings/postcarder/pkg/loader"
	"github.com/spf13/cobra"
)

// compileCmd represents the compile command
var compileCmd = &cobra.Command{
	Use:   "compile",
	Short: "Compiles images & metadata into a postcard file",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		path, err := filepath.Abs(args[0])
		if err != nil {
			return fmt.Errorf("unknown file path: %w", err)
		}

		if _, err := os.Stat(path); err != nil {
			return fmt.Errorf("file doesn't exist: %w", err)
		}

		dir := filepath.Dir(path)
		base := strings.SplitN(filepath.Base(path), "-", 2)[0]

		postcard, err := loader.QuickLoad(dir, base)
		if err != nil {
			return err
		}

		return compiler.WritePostcardFile(postcard, fmt.Sprintf("%s.postcard", base))
	},
}

func init() {
	rootCmd.AddCommand(compileCmd)
}
