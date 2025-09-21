package cmd

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/ZonCen/Cloak/internal/helpers"
	"github.com/ZonCen/Cloak/internal/vault"
)

// editCmd represents the edit command
var editCmd = &cobra.Command{
	Use:   "edit",
	Short: "A brief description of your command",
	Args:  cobra.MaximumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		helpers.LogVerbose("Getting Arguments from command line")
		inputFile := args[0]

		rawKey, err := vault.GetKey()
		if err != nil {
			fmt.Println("Error getting key from environment:", err)
			return
		}

		err = vault.EditFile(inputFile, rawKey)
		if err != nil {
			fmt.Println("Error editing file:", err)
			return
		}

		fmt.Println("editing file finished successfully")
	},
}

func init() {
	rootCmd.AddCommand(editCmd)
}
