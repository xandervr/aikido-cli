package commands

import (
	"github.com/spf13/cobra"
	"github.com/xandervr/aikido-cli/internal/cli"
)

func NewLicenses(g *cli.Globals) *cobra.Command {
	cmd := &cobra.Command{Use: "licenses", Short: "License inventory"}
	cmd.AddCommand(simpleList(g, "list", "List licenses across the workspace", "/licenses"))
	return cmd
}
