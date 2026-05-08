package commands

import (
	"net/http"

	"github.com/spf13/cobra"
	"github.com/xandervr/aikido-cli/internal/cli"
)

func NewCodeQuality(g *cli.Globals) *cobra.Command {
	cmd := &cobra.Command{Use: "code-quality", Short: "Code quality findings"}
	cmd.AddCommand(endpointCommand(g, endpointCommandConfig{
		Use:    "findings",
		Short:  "List code quality findings for a pull request",
		Method: http.MethodGet,
		Path:   staticPath("/code-quality/findings"),
	}))
	return cmd
}
