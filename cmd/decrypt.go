package cmd

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/ZonCen/Cloak/internal/helpers"
	"github.com/ZonCen/Cloak/internal/vault"
)

// decryptCmd represents the decrypt command
var decryptCmd = &cobra.Command{
	Use:   "decrypt",
	Short: "Decrypt a file",
	Args:  cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		helpers.LogVerbose("Getting Arguments from command line")
		inputFile := args[0]

		rawKey, err := vault.GetKey()
		if err != nil {
			fmt.Println("Error getting key from environment:", err)
			return
		}

		err = vault.DecryptFile(inputFile, outputFile, rawKey)
		if err != nil {
			fmt.Println("Error decrypting file:", err)
			return
		}

		fmt.Println("File has been decrypted successfully")
	},
}

func init() {
	rootCmd.AddCommand(decryptCmd)

	decryptCmd.Flags().StringVar(&outputFile, "output-file", "", "Path to the output file")
}
