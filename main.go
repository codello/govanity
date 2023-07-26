package main

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"codello.dev/govanity/cmd/server"
	"codello.dev/govanity/cmd/version"
)

// rootCmd is the starting point of the govanity program.
var rootCmd = &cobra.Command{
	Use:   "govanity",
	Short: "govanity is a simple vanity URL server for Go packages",
}

func init() {
	rootCmd.AddCommand(version.Command, server.Command)
}

func main() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
