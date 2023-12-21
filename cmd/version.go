package cmd

import (
	"fmt"

	"github.com/dtrejod/goexif/version"
	"github.com/spf13/cobra"
)

var versionCmd = &cobra.Command{
	Use: "version",
	Run: versionRun,
}

func versionRun(_ *cobra.Command, _ []string) {
	fmt.Printf("version: %s\n", version.Version())
	return
}

func init() {
	rootCmd.AddCommand(versionCmd)
}
