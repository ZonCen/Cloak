package cmd

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/ZonCen/Cloak/internal/helpers"
	"github.com/ZonCen/Cloak/internal/vault"
)

// encryptCmd represents the encrypt command
var encryptCmd = &cobra.Command{
	Use:   "encrypt",
	Short: "Encrypt a file",
	Args:  cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		helpers.LogVerbose("Getting Arguments from command line")
		inputFile := args[0]

		rawKey, err := vault.GetKey()
		if err != nil {
			fmt.Println("Error getting key from environment:", err)
			return
		}

		err = vault.EncryptFile(inputFile, outputFile, rawKey)
		if err != nil {
			fmt.Println("Error encrypting file:", err)
			return
		}

		fmt.Printf("File %s encrypted successfully to %s\n", inputFile, outputFile)
	},
}

func init() {
	rootCmd.AddCommand(encryptCmd)

	encryptCmd.Flags().StringVar(&outputFile, "output-file", "", "Path to the output file")
}
