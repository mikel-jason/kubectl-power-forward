package main

import (
	"fmt"

	"github.com/spf13/cobra"
)

var version, commit, date string // filled by goreleaser

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Show version details",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Printf("Version: %s\n", version)
		fmt.Printf("Build revision: %s\n", commit)
		fmt.Printf("Build date: %s\n", date)
	},
}
