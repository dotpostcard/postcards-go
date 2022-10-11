package cmd

import (
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

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}
