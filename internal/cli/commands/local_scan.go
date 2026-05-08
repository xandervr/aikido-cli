package commands

import (
	"net/http"

	"github.com/spf13/cobra"
	"github.com/xandervr/aikido-cli/internal/cli"
)

func NewLocalScan(g *cli.Globals) *cobra.Command {
	cmd := &cobra.Command{Use: "local-scan", Short: "Local scanner metadata"}
	cmd.AddCommand(endpointCommand(g, endpointCommandConfig{
		Use:    "latest",
		Short:  "Get latest local scanner version",
		Method: http.MethodGet,
		Path:   staticPath("/localscan/latest"),
	}))
	return cmd
}
