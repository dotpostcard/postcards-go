package cmd

import (
	"fmt"

	"github.com/dotpostcard/postcards-go/internal/cmdhelp"
	"github.com/dotpostcard/postcards-go/make3d"
	"github.com/spf13/cobra"
)

var make3DCmd = &cobra.Command{
	Use:   "3d",
	Short: "Creates 3D models of postcards",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		pc, info, err := cmdhelp.OpenPostcard(args[0], false)
		if err != nil {
			return err
		}

		f, out, err := cmdhelp.CreateFile(rootCmd, info, ".obj.zip")
		if err != nil {
			return err
		}

		if err := make3d.WriteObjZip(pc, f, nil); err != nil {
			return err
		}

		fmt.Printf("Writing 3D model to %s\n", out)
		return nil
	},
}

func init() {
	rootCmd.AddCommand(make3DCmd)
}
