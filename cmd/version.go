package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(versionCmd)
}

// Version is the klein package version
// to be filled in at compile time using ldflags
var Version = "-dev"

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "print the version number of klein",
	Long:  "print the version number of klein",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Printf("klein v%s\n", Version)
	},
}
