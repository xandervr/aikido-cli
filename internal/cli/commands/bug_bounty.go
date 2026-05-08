package commands

import (
	"net/http"

	"github.com/spf13/cobra"
	"github.com/xandervr/aikido-cli/internal/cli"
)

func NewBugBounty(g *cli.Globals) *cobra.Command {
	cmd := &cobra.Command{Use: "bug-bounty", Short: "Bug bounty report validation"}
	cmd.AddCommand(endpointCommand(g, endpointCommandConfig{
		Use:    "validate-report <program-id>",
		Short:  "Validate a bug bounty report",
		Method: http.MethodPost,
		Args:   cobra.ExactArgs(1),
		Path:   oneArgPath("/bug_bounty/program/%s/report"),
	}))
	return cmd
}
