package commands

import (
	"net/http"

	"github.com/spf13/cobra"
	"github.com/xandervr/aikido-cli/internal/cli"
)

func NewLicenses(g *cli.Globals) *cobra.Command {
	cmd := &cobra.Command{Use: "licenses", Short: "License inventory"}
	cmd.AddCommand(
		simpleList(g, "list", "List licenses across the workspace", "/licenses"),
		endpointCommand(g, endpointCommandConfig{Use: "overwrite", Short: "Overwrite license metadata", Method: http.MethodPost, Path: staticPath("/licenses/overwrite")}),
	)
	return cmd
}
