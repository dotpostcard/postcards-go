package cmd

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/jphastings/postcard-go"
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

		filename, data, err := postcard.CompileFiles(path)
		if err != nil {
			return err
		}

		fmt.Printf("Writing postcard file to %s\n", filename)
		return os.WriteFile(filename, data, 0600)
	},
}

func init() {
	rootCmd.AddCommand(compileCmd)
}
