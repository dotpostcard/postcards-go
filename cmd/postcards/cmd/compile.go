package cmd

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/dotpostcard/postcards-go/compile"
	"github.com/dotpostcard/postcards-go/internal/cmdhelp"
	"github.com/spf13/cobra"
)

// compileCmd represents the compile command
var compileCmd = &cobra.Command{
	Use:   "compile",
	Short: "Compiles images & metadata into a postcard file, or web-compatible equivalent",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		path, err := filepath.Abs(args[0])
		if err != nil {
			return fmt.Errorf("unknown file path: %w", err)
		}

		outdir, err := cmdhelp.Outdir(cmd, path)
		if err != nil {
			return err
		}

		override, err := cmd.Flags().GetBool("override")
		if err != nil {
			override = false
		}

		webFormat, err := cmd.Flags().GetBool("web")
		if err != nil {
			webFormat = false
		}

		filenames, datas, err := compile.Files(path, !override, webFormat)
		if err != nil {
			return err
		}

		if datas == nil {
			fmt.Printf("Postcard files already exist, skipping: %s\n", strings.Join(filenames, ", "))
		}

		fmt.Printf("Writing postcard files to %s\n", outdir)
		for i, filename := range filenames {
			if err := os.WriteFile(filepath.Join(outdir, filename), datas[i], 0600); err != nil {
				return fmt.Errorf("unable to write file %s: %w", filename, err)
			}
			fmt.Printf("↪ Wrote postcard file to %s\n", filename)
		}
		return nil
	},
}

func init() {
	compileCmd.Flags().Bool("override", false, "overrides output files")
	compileCmd.Flags().Bool("web", false, "make web-compatible postcard file")
	rootCmd.AddCommand(compileCmd)
}
