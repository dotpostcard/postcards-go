package cmd

import (
	"fmt"

	"github.com/jphastings/postcarder"
	"github.com/spf13/cobra"
)

var infoCmd = &cobra.Command{
	Use:   "info",
	Short: "prints info about the specified postcard",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		pc, err := postcarder.ReadFile(args[0], true)
		if err != nil {
			return fmt.Errorf("unable to read postcard file: %w", err)
		}

		fmt.Println(pc.Meta)
		return nil
	},
}

func init() {
	rootCmd.AddCommand(infoCmd)
}
