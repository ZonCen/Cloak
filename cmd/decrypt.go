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
		helpers.LogVerbose("Getting Arguments from command line")
		rawKey, err := vault.GetEnv("CLOAK_KEY")
		if err != nil {
			fmt.Println("Error getting key from environment:", err)
			return
		}

		helpers.LogVerbose("Checking if file exists and is a vault file")
		_, err = helpers.ReadFile(inputFile)
		if err != nil {
			fmt.Println("Error reading input file:", err)
			return
		}
		encrypted, err := vault.IsEncryptedFile(inputFile)
		if err != nil {
			fmt.Println("Error checking if file is encrypted:", err)
			return
		}
		if !encrypted {
			fmt.Println("Input file is not an encrypted vault file")
			return
		}

		helpers.LogVerbose("Confirming output file")
		if outputFile == "" {
			outputFile = inputFile
		}

		helpers.LogVerbose("Decrypting file")
		err = vault.DecryptFile(rawKey, inputFile, outputFile)
		if err != nil {
			fmt.Println("Error decrypting file:", err)
			return
		}

		fmt.Println("File has been decrypted successfully to", outputFile)
	},
}

func init() {
	rootCmd.AddCommand(decryptCmd)

	decryptCmd.Flags().StringVar(&outputFile, "output", "", "Path to the output file")
}
