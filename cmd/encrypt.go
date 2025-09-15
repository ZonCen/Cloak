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

		helpers.LogVerbose("Checking if input file already is encrypted")
		encrypted, err := vault.IsEncryptedFile(inputFile)
		if err == nil && encrypted {
			fmt.Println("input file is already encrypted. Please choose a different input file.")
			return
		}

		rawKey, err := vault.GetEnv("CLOAK_KEY")
		if err != nil {
			fmt.Println("Error getting key from environment:", err)
			return
		}

		helpers.LogVerbose("Checking if file exists")
		_, err = helpers.ReadFile(inputFile)
		if err != nil {
			fmt.Println("Error reading input file:", err)
			return
		}

		helpers.LogVerbose("Confirming output file")
		if outputFile == "" {
			outputFile = inputFile + ".vault"
		}
		if helpers.CheckSuffix(inputFile, ".vault") {
			outputFile = inputFile
		}

		helpers.LogVerbose("Encrypting file")
		err = vault.EncryptFile(rawKey, inputFile, outputFile)
		if err != nil {
			fmt.Println("Error encrypting file:", err)
			return
		}

		fmt.Printf("File %s encrypted successfully to %s\n", inputFile, outputFile)
	},
}

func init() {
	rootCmd.AddCommand(encryptCmd)

	encryptCmd.Flags().StringVar(&outputFile, "output", "", "Path to the output file")
}
