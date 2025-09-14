package cmd

import (
	"fmt"
	"os"

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
		// Decrypt the file
		helpers.LogVerbose("Getting Arguments from command line")
		inputFile := args[0]
		helpers.LogVerbose("Getting Arguments from command line")
		rawKey, err := vault.GetEnv("CLOAK_KEY")
		if err != nil {
			fmt.Println("Error getting key from environment:", err)
			return
		}

		helpers.LogVerbose("Checking if file exists and have correct suffix")
		_, err = helpers.ReadFile(inputFile)
		if err != nil {
			fmt.Println("Error reading input file:", err)
			return
		}
		if !helpers.CheckSuffix(inputFile, ".vault") {
			fmt.Println("Input file does not have .vault suffix")
			return
		}

		plainText, err := vault.ReadEncryptedFile(rawKey, inputFile)
		if err != nil {
			fmt.Println("Error decrypting file:", err)
			return
		}

		tempFile, err := vault.CreateTempFile()
		if err != nil {
			fmt.Println("Error creating temporary file:", err)
			return
		}

		defer func() {
			if err := tempFile.Close(); err != nil {
				fmt.Println("Error closing temporary file:", err)
			}
		}()
		defer func() {
			if err := os.Remove(tempFile.Name()); err != nil {
				fmt.Println("Error removing temporary file:", err)
			}
		}()

		if _, err := tempFile.Write(plainText); err != nil {
			fmt.Println("Error writing to temporary file:", err)
			return
		}

		editor := vault.GetEditor()
		if err := vault.LaunchEditor(editor, tempFile.Name()); err != nil {
			fmt.Println("Error launching editor:", err)
			return
		}

		_, err = helpers.ReadFile(tempFile.Name())
		if err != nil {
			fmt.Println("Error reading edited content:", err)
		}

		if err := vault.EncryptFile(rawKey, tempFile.Name(), inputFile); err != nil {
			fmt.Println("Error re-encrypting file:", err)
			return
		}

		fmt.Println("edit called")
	},
}

func init() {
	rootCmd.AddCommand(editCmd)
}
