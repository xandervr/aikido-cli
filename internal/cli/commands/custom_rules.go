package commands

import (
	"net/http"

	"github.com/spf13/cobra"
	"github.com/xandervr/aikido-cli/internal/cli"
)

func NewCustomRules(g *cli.Globals) *cobra.Command {
	cmd := &cobra.Command{Use: "custom-rules", Short: "Custom SAST rules"}
	cmd.AddCommand(
		simpleList(g, "list", "List custom rules", "/repositories/sast/custom-rules"),
		endpointCommand(g, endpointCommandConfig{Use: "create", Short: "Create custom rule", Method: http.MethodPost, Path: staticPath("/repositories/sast/custom-rules")}),
		simpleGet(g, "get <id>", "Get a custom rule", "/repositories/sast/custom-rules"),
		endpointCommand(g, endpointCommandConfig{Use: "update <id>", Short: "Edit custom rule", Method: http.MethodPut, Args: cobra.ExactArgs(1), Path: oneArgPath("/repositories/sast/custom-rules/%s")}),
		endpointCommand(g, endpointCommandConfig{Use: "delete <id>", Short: "Remove custom rule", Method: http.MethodDelete, Args: cobra.ExactArgs(1), Path: oneArgPath("/repositories/sast/custom-rules/%s"), Confirm: true}),
	)
	return cmd
}
