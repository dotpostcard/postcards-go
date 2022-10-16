package cmd

import (
	"fmt"
	"os"
	"path"
	"reflect"

	"github.com/dotpostcard/postcards-go"
	"github.com/dotpostcard/postcards-go/internal/helpers"
	"github.com/spf13/cobra"
)

var infoCmd = &cobra.Command{
	Use:   "info",
	Short: "prints info about the specified postcard",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		file := args[0]

		f, err := os.Open(file)
		if err != nil {
			return fmt.Errorf("unable to open file: %w", err)
		}
		defer f.Close()

		pc, err := postcards.Read(f, true)
		if err != nil {
			return fmt.Errorf("unable to read postcard file: %w", err)
		}

		fmt.Printf("Postcard:   %s\n", path.Base(file))
		fmt.Printf("Size:       %s\n", helpers.SizeHuman(f))
		fmt.Println()

		printUnlessZero("From:       %s\n", pc.Meta.Sender)
		printUnlessZero("To:         %s\n", pc.Meta.Recipient)
		printUnlessZero("Sent on:    %s\n", pc.Meta.SentOn)
		printUnlessZero("Location:   %s\n", pc.Meta.Location.Name)
		printUnlessZero("Flips:      %s\n", pc.Meta.Flip)
		printUnlessZero("Dimensions: %s\n", pc.Meta.FrontDimensions)
		fmt.Println()

		printUnlessZero("Front transcription\n-------------------%s\n\n", pc.Meta.Front.Transcription)
		printUnlessZero("Front description\n-----------------\n%s\n\n", pc.Meta.Front.Description)
		printUnlessZero("Back transcription\n------------------\n%s\n\n", pc.Meta.Back.Transcription)
		printUnlessZero("Back description\n----------------\n%s\n\n", pc.Meta.Back.Description)

		return nil
	},
}

func printUnlessZero(format string, vals ...interface{}) {
	for _, val := range vals {
		if reflect.ValueOf(val).IsZero() {
			return
		}
	}

	fmt.Printf(format, vals...)
}

func init() {
	rootCmd.AddCommand(infoCmd)
}
