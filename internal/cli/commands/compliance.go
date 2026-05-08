package commands

import (
	"net/http"

	"github.com/spf13/cobra"
	"github.com/xandervr/aikido-cli/internal/cli"
)

func NewCompliance(g *cli.Globals) *cobra.Command {
	cmd := &cobra.Command{Use: "compliance", Short: "Compliance overviews"}
	cmd.AddCommand(
		simpleList(g, "soc2", "SOC2 compliance overview", "/report/soc2/overview"),
		simpleList(g, "nis2", "NIS2 compliance overview", "/report/nis2/overview"),
		simpleList(g, "iso27001", "ISO 27001 compliance overview", "/report/iso/overview"),
		endpointCommand(g, endpointCommandConfig{Use: "cis", Short: "CIS compliance overview", Method: http.MethodGet, Path: staticPath("/report/cis/overview")}),
		endpointCommand(g, endpointCommandConfig{Use: "cis-aws", Short: "CIS AWS compliance overview", Method: http.MethodGet, Path: staticPath("/report/cis_aws/overview")}),
	)
	return cmd
}
