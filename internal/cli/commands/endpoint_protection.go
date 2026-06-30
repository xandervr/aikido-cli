package commands

import (
	"net/http"

	"github.com/spf13/cobra"
	"github.com/xandervr/aikido-cli/internal/cli"
)

func NewEndpointProtection(g *cli.Globals) *cobra.Command {
	cmd := &cobra.Command{Use: "endpoint-protection", Short: "Endpoint protection"}
	cmd.AddCommand(
		endpointCommand(g, endpointCommandConfig{
			Use:    "activity-logs",
			Short:  "List endpoint activity logs",
			Method: http.MethodGet,
			Path:   staticPath("/endpoint-protection/activityLogs"),
		}),
		endpointCommand(g, endpointCommandConfig{
			Use:    "devices",
			Short:  "List endpoint devices",
			Method: http.MethodGet,
			Path:   staticPath("/endpoint-protection/devices"),
		}),
		endpointCommand(g, endpointCommandConfig{
			Use:    "installed-packages",
			Short:  "List installed packages",
			Method: http.MethodGet,
			Path:   staticPath("/endpoint-protection/installed-packages"),
		}),
		endpointCommand(g, endpointCommandConfig{
			Use:    "permission-groups",
			Short:  "List endpoint permission groups",
			Method: http.MethodGet,
			Path:   staticPath("/endpoint-protection/permission-groups"),
		}),
		endpointCommand(g, endpointCommandConfig{
			Use:    "exceptions <ecosystem>",
			Short:  "List endpoint exceptions",
			Method: http.MethodGet,
			Args:   cobra.ExactArgs(1),
			Path:   oneArgPath("/endpoint-protection/%s/exceptions"),
		}),
		endpointCommand(g, endpointCommandConfig{
			Use:    "add-exception <ecosystem>",
			Short:  "Add endpoint exception",
			Method: http.MethodPost,
			Args:   cobra.ExactArgs(1),
			Path:   oneArgPath("/endpoint-protection/%s/exceptions"),
		}),
		endpointCommand(g, endpointCommandConfig{
			Use:     "remove-exception <package-exception-id>",
			Short:   "Remove endpoint exception",
			Method:  http.MethodDelete,
			Args:    cobra.ExactArgs(1),
			Path:    oneArgPath("/endpoint-protection/exceptions/%s"),
			Confirm: true,
		}),
	)
	return cmd
}
