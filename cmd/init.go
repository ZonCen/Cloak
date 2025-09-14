package cmd

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/ZonCen/Cloak/internal/vault"
)

// initCmd represents the init command
var initCmd = &cobra.Command{
	Use:   "init",
	Short: "A brief description of your command",
	Run: func(cmd *cobra.Command, args []string) {
		key, err := vault.GetEnv("CLOAK_KEY")
		if err != nil || string(key) != "" {
			fmt.Println("CLOAK_KEY environment variable is already set.")
			return
		}

		key, err = vault.GenerateRandomByteKey(32)
		if err != nil {
			fmt.Println("Error generating key:", err)
			return
		}
		encodedKey := vault.EncodeKey(key)

		fmt.Println("Your master key (store this securely!):", encodedKey)
		fmt.Println("To run cloak please set the CLOAK_KEY environment variable")
		fmt.Printf("export CLOAK_KEY=%s\n", encodedKey)
	},
}

func init() {
	rootCmd.AddCommand(initCmd)
}
