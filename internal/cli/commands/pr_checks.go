package commands

import (
	"net/http"

	"github.com/spf13/cobra"
	"github.com/xandervr/aikido-cli/internal/cli"
)

func NewPRChecks(g *cli.Globals) *cobra.Command {
	var repo string
	cmd := &cobra.Command{Use: "pr-checks", Short: "CI/PR scans"}
	listCmd := &cobra.Command{
		Use:   "list",
		Short: "List PR checks",
		RunE: func(cmd *cobra.Command, args []string) error {
			c, err := g.Client()
			if err != nil {
				return err
			}
			q := map[string]string{}
			if repo != "" {
				q["filter_code_repo_id"] = repo
			}
			var raw any
			if err := c.Get(cmd.Context(), "/report/ciScans", q, &raw); err != nil {
				return err
			}
			return g.Renderer().Render(raw)
		},
	}
	listCmd.Flags().StringVar(&repo, "repo", "", "filter by repo ID")
	cmd.AddCommand(
		listCmd,
		endpointCommand(g, endpointCommandConfig{
			Use:    "issue-actions",
			Short:  "List PR check manual actions",
			Method: http.MethodGet,
			Path:   staticPath("/report/ciScans/issueActions"),
		}),
	)
	return cmd
}
