package cmd

import (
	"fmt"
	"os"

	"github.com/dotpostcard/postcards-go"
	"github.com/spf13/cobra"
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:     "postcards",
	Short:   "A tool to create postcard files from images & descriptions of the front and back",
	Version: postcards.Version.String(),
}

func checkErr(err error) {
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}

func Execute() {
	rootCmd.PersistentFlags().Bool("here", false, "Output files in the current working directory")
	rootCmd.PersistentFlags().Bool("there", true, "Output files in the same directory as the source data")
	rootCmd.PersistentFlags().String("outdir", "", "Output files to the given directory")
	rootCmd.MarkFlagsMutuallyExclusive("here", "there", "outdir")

	checkErr(rootCmd.Execute())
}
