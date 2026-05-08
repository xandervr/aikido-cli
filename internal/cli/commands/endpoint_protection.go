package commands

import (
	"net/http"

	"github.com/spf13/cobra"
	"github.com/xandervr/aikido-cli/internal/cli"
)

func NewEndpointProtection(g *cli.Globals) *cobra.Command {
	cmd := &cobra.Command{Use: "endpoint-protection", Short: "Endpoint protection activity"}
	cmd.AddCommand(endpointCommand(g, endpointCommandConfig{
		Use:    "activity-logs",
		Short:  "List endpoint protection activity logs",
		Method: http.MethodGet,
		Path:   staticPath("/endpoint-protection/activityLogs"),
	}))
	return cmd
}
