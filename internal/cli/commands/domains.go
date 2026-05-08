package commands

import (
	"net/http"

	"github.com/spf13/cobra"
	"github.com/xandervr/aikido-cli/internal/cli"
)

func NewDomains(g *cli.Globals) *cobra.Command {
	cmd := &cobra.Command{Use: "domains", Short: "Domains and surface monitoring"}
	cmd.AddCommand(
		endpointCommand(g, endpointCommandConfig{
			Use:    "list",
			Short:  "List domains",
			Method: http.MethodGet,
			Path:   staticPath("/domains"),
		}),
		endpointCommand(g, endpointCommandConfig{
			Use:    "create",
			Short:  "Create a domain",
			Method: http.MethodPost,
			Path:   staticPath("/domains"),
		}),
		endpointCommand(g, endpointCommandConfig{
			Use:     "delete <domain-id>",
			Short:   "Remove a domain",
			Method:  http.MethodDelete,
			Args:    cobra.ExactArgs(1),
			Path:    oneArgPath("/domains/%s"),
			Confirm: true,
		}),
		endpointCommand(g, endpointCommandConfig{
			Use:    "scan",
			Short:  "Start a scan for a domain",
			Method: http.MethodPost,
			Path:   staticPath("/domains/scan"),
		}),
		endpointCommand(g, endpointCommandConfig{
			Use:    "headers <domain-id>",
			Short:  "Update auth headers",
			Method: http.MethodPost,
			Args:   cobra.ExactArgs(1),
			Path:   oneArgPath("/domains/%s/headers"),
		}),
		endpointCommand(g, endpointCommandConfig{
			Use:    "custom-headers <domain-id>",
			Short:  "Update custom scan headers",
			Method: http.MethodPost,
			Args:   cobra.ExactArgs(1),
			Path:   oneArgPath("/domains/%s/custom-headers"),
		}),
		endpointCommand(g, endpointCommandConfig{
			Use:    "update-openapi <domain-id>",
			Short:  "Update OpenAPI spec",
			Method: http.MethodPut,
			Args:   cobra.ExactArgs(1),
			Path:   oneArgPath("/domains/%s/update/openapi-spec"),
		}),
		domainSubdomains(g),
	)
	return cmd
}

func domainSubdomains(g *cli.Globals) *cobra.Command {
	cmd := &cobra.Command{Use: "subdomains", Short: "Domain subdomains"}
	cmd.AddCommand(
		endpointCommand(g, endpointCommandConfig{
			Use:    "list <domain-id>",
			Short:  "List subdomains",
			Method: http.MethodGet,
			Args:   cobra.ExactArgs(1),
			Path:   oneArgPath("/domains/%s/subdomains"),
		}),
		endpointCommand(g, endpointCommandConfig{
			Use:    "add <domain-id>",
			Short:  "Add a subdomain",
			Method: http.MethodPost,
			Args:   cobra.ExactArgs(1),
			Path:   oneArgPath("/domains/%s/subdomains"),
		}),
	)
	return cmd
}
