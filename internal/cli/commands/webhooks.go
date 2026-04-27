package commands

import (
	"github.com/spf13/cobra"
	"github.com/xandervr/aikido-cli/internal/cli"
)

func NewWebhooks(g *cli.Globals) *cobra.Command {
	cmd := &cobra.Command{Use: "webhooks", Short: "Configured webhooks"}
	cmd.AddCommand(simpleList(g, "list", "List webhooks", "/webhooks"))
	return cmd
}
