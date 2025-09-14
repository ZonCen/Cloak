package cmd

import (
	"os"

	"github.com/spf13/cobra"

	"github.com/ZonCen/Cloak/internal/helpers"
)

var (
	outputFile string
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "Cloak",
	Short: "A brief description of your application",
}

func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	rootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")

	rootCmd.PersistentFlags().BoolVarP(&helpers.Verbose, "verbose", "v", false, "Show detailed output")
}
