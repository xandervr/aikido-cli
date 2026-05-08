package commands

import (
	"net/http"

	"github.com/spf13/cobra"
	"github.com/xandervr/aikido-cli/internal/cli"
)

func NewAccessTokens(g *cli.Globals) *cobra.Command {
	cmd := &cobra.Command{Use: "access-tokens", Short: "Aikido access token management"}
	cmd.AddCommand(endpointCommand(g, endpointCommandConfig{
		Use:    "code-scanning",
		Short:  "Update code scanning access token",
		Method: http.MethodPost,
		Path:   staticPath("/access-tokens/code-scanning"),
	}))
	return cmd
}
