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

		pc, err := postcards.ReadFile(path, false)
		if err != nil {
			return err
		}

		base := filepath.Base(path)
		name := base[:len(base)-len(filepath.Ext(base))]

		frontName := fmt.Sprintf("%s-front.webp", name)
		if err := os.WriteFile(frontName, pc.Front, 0666); err != nil {
			return err
		}
		fmt.Printf("%s: Front image of postcard\n", frontName)

		backName := fmt.Sprintf("%s-back.webp", name)
		if err := os.WriteFile(backName, pc.Back, 0666); err != nil {
			return err
		}
		fmt.Printf("%s:  Back image of postcard\n", backName)

		metaName := fmt.Sprintf("%s-meta.json", name)
		meta, err := postcards.MetadataBytes(pc.Meta, true)
		if err != nil {
			return err
		}
		if err := os.WriteFile(metaName, meta, 0666); err != nil {
			return err
		}
		fmt.Printf("%s:  Metadata of postcard\n", metaName)

		return nil
	},
}

func init() {
	rootCmd.AddCommand(extractCmd)
}
