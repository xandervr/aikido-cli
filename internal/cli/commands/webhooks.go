package commands

import (
	"net/http"

	"github.com/spf13/cobra"
	"github.com/xandervr/aikido-cli/internal/cli"
)

func NewWebhooks(g *cli.Globals) *cobra.Command {
	cmd := &cobra.Command{Use: "webhooks", Short: "Configured webhooks"}
	cmd.AddCommand(
		simpleList(g, "list", "List webhooks", "/webhooks"),
		endpointCommand(g, endpointCommandConfig{Use: "add", Short: "Add webhook", Method: http.MethodPost, Path: staticPath("/webhooks")}),
		endpointCommand(g, endpointCommandConfig{Use: "delete <webhook-id>", Short: "Remove webhook", Method: http.MethodDelete, Args: cobra.ExactArgs(1), Path: oneArgPath("/webhooks/%s"), Confirm: true}),
	)
	return cmd
}
