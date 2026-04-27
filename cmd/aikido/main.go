package main

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/xandervr/aikido-cli/internal/version"
)

func main() {
	root := &cobra.Command{
		Use:           "aikido",
		Short:         "Aikido Security CLI",
		Long:          "Command-line client for the Aikido Security public REST API.",
		SilenceUsage:  true,
		SilenceErrors: true,
		Version:       version.Version,
	}
	if err := root.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, "error:", err)
		os.Exit(1)
	}
}
