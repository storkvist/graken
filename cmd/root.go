package cmd

import (
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use: "graken",
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		panic(err)
	}
}
