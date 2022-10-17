package cmd

import (
	"fmt"
	"reflect"

	"github.com/dotpostcard/postcards-go/internal/cmdhelp"
	"github.com/spf13/cobra"
)

var infoCmd = &cobra.Command{
	Use:   "info",
	Short: "prints info about the specified postcard",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		pc, info, err := cmdhelp.OpenPostcard(args[0], true)
		if err != nil {
			return err
		}

		fmt.Printf("Postcard: %s.postcard\n", info.Base)
		fmt.Printf("Size:     %s\n", info.SizeHuman)
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
