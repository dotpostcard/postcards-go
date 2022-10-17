package cmd

import (
	"github.com/dotpostcard/postcards-go"
	"github.com/dotpostcard/postcards-go/internal/cmdhelp"
	"github.com/spf13/cobra"
)

// extractCmd represents the compile command
var extractCmd = &cobra.Command{
	Use:   "extract",
	Short: "Extracts images & metadata from a postcard file",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		pc, info, err := cmdhelp.OpenPostcard(args[0], false)
		if err != nil {
			return err
		}

		if err := cmdhelp.WriteFile(cmd, info, "-front.webp", pc.Front, "front image"); err != nil {
			return err
		}
		if err := cmdhelp.WriteFile(cmd, info, "-back.webp", pc.Back, "back image"); err != nil {
			return err
		}

		meta, err := postcards.MetadataBytes(pc.Meta, true)
		if err != nil {
			return err
		}
		if err := cmdhelp.WriteFile(cmd, info, "-meta.json", meta, "metadata"); err != nil {
			return err
		}

		return nil
	},
}

func init() {
	rootCmd.AddCommand(extractCmd)
}
