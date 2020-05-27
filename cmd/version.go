package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"runtime"
)

const (
	version string = "0.0.1"
)

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Show graken version",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Printf("Graken — auto-fetch multiple Git repositories, v%s, %s/%s\n", version, runtime.GOOS, runtime.GOARCH)
	},
}

func init() {
	rootCmd.AddCommand(versionCmd)
}
