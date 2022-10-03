package cmd

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/dotpostcard/postcards-go"
	"github.com/spf13/cobra"
)

// extractCmd represents the compile command
var extractCmd = &cobra.Command{
	Use:   "extract",
	Short: "Extracts images & metadata from a postcard file",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		path, err := filepath.Abs(args[0])
		if err != nil {
			return fmt.Errorf("unknown file path: %w", err)
		}

		if _, err := os.Stat(path); err != nil {
			return fmt.Errorf("file doesn't exist: %w", err)
		}

		pc, err := postcards.ReadFile(path, false)
		if err != nil {
			return err
		}

		fmt.Println(pc.Meta)

		return fmt.Errorf("not implemented")
	},
}

func init() {
	rootCmd.AddCommand(extractCmd)
}
