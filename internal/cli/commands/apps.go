package commands

import (
	"github.com/spf13/cobra"
	"github.com/xandervr/aikido-cli/internal/cli"
)

func NewApps(g *cli.Globals) *cobra.Command {
	cmd := &cobra.Command{Use: "apps", Short: "Zen apps"}
	cmd.AddCommand(simpleList(g, "list", "List Zen apps", "/apps"))
	return cmd
}
