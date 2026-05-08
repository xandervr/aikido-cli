package commands

import (
	"net/http"

	"github.com/spf13/cobra"
	"github.com/xandervr/aikido-cli/internal/cli"
)

func NewUsers(g *cli.Globals) *cobra.Command {
	cmd := &cobra.Command{Use: "users", Short: "Workspace users"}
	cmd.AddCommand(
		simpleList(g, "list", "List users", "/users"),
		simpleGet(g, "get <id>", "Get a user", "/users"),
		endpointCommand(g, endpointCommandConfig{Use: "rights <id>", Short: "Update user rights", Method: http.MethodPut, Args: cobra.ExactArgs(1), Path: oneArgPath("/users/%s/rights")}),
	)
	return cmd
}
