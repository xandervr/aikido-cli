package commands

import (
	"strconv"

	"github.com/spf13/cobra"
	"github.com/xandervr/aikido-cli/internal/cli"
)

func NewCodeQuality(g *cli.Globals) *cobra.Command {
	cmd := &cobra.Command{Use: "code-quality", Short: "Code quality findings"}
	cmd.AddCommand(codeQualityFindings(g))
	return cmd
}

func codeQualityFindings(g *cli.Globals) *cobra.Command {
	var repoID int
	var pr string
	cmd := &cobra.Command{
		Use:   "findings",
		Short: "List code quality findings for a pull request",
		RunE: func(cmd *cobra.Command, args []string) error {
			c, err := g.Client()
			if err != nil {
				return err
			}
			q := map[string]string{
				"code_repo_id": strconv.Itoa(repoID),
				"pr_number":    pr,
			}
			var raw any
			if err := c.Get(cmd.Context(), "/code-quality/findings", q, &raw); err != nil {
				return err
			}
			return g.Renderer().Render(raw)
		},
	}
	cmd.Flags().IntVar(&repoID, "repo", 0, "code repository ID (required)")
	cmd.Flags().StringVar(&pr, "pr", "", "pull request number (required)")
	_ = cmd.MarkFlagRequired("repo")
	_ = cmd.MarkFlagRequired("pr")
	return cmd
}
