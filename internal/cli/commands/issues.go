package commands

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/xandervr/aikido-cli/internal/cli"
)

type IssueGroup struct {
	ID       int    `json:"id"        aikido:"column,header=ID"`
	Title    string `json:"title"     aikido:"column,header=Title"`
	Severity string `json:"severity"  aikido:"column,header=Severity"`
	Type     string `json:"type"      aikido:"column,header=Type"`
	RepoName string `json:"repo_name" aikido:"column,header=Repo"`
	Status   string `json:"status"    aikido:"column,header=Status"`
}

type issuesListOpts struct {
	Severity string
	Status   string
	Repo     int
	Team     int
	Page     int
	PerPage  int
}

func NewIssues(g *cli.Globals) *cobra.Command {
	cmd := &cobra.Command{Use: "issues", Short: "Open issues / vulnerabilities"}
	cmd.AddCommand(issuesList(g), issuesGet(g), issuesExport(g))
	return cmd
}

func issuesList(g *cli.Globals) *cobra.Command {
	var opts issuesListOpts
	cmd := &cobra.Command{
		Use:   "list",
		Short: "List open issue groups (vulnerabilities)",
		RunE: func(cmd *cobra.Command, args []string) error {
			c, err := g.Client()
			if err != nil {
				return err
			}
			q := map[string]string{}
			if opts.Severity != "" {
				q["filter_severity"] = opts.Severity
			}
			if opts.Status != "" {
				q["filter_status"] = opts.Status
			}
			if opts.Repo > 0 {
				q["filter_code_repo_id"] = fmt.Sprintf("%d", opts.Repo)
			}
			if opts.Team > 0 {
				q["filter_team_id"] = fmt.Sprintf("%d", opts.Team)
			}
			if opts.Page > 0 {
				q["page"] = fmt.Sprintf("%d", opts.Page)
			}
			if opts.PerPage > 0 {
				q["per_page"] = fmt.Sprintf("%d", opts.PerPage)
			}
			var groups []IssueGroup
			if err := c.Get(cmd.Context(), "/open-issue-groups", q, &groups); err != nil {
				return err
			}
			return g.Renderer().Render(groups)
		},
	}
	cmd.Flags().StringVar(&opts.Severity, "severity", "", "filter: critical|high|medium|low")
	cmd.Flags().StringVar(&opts.Status, "status", "", "filter: open|ignored|snoozed|closed")
	cmd.Flags().IntVar(&opts.Repo, "repo", 0, "filter by repo ID")
	cmd.Flags().IntVar(&opts.Team, "team", 0, "filter by team ID")
	cmd.Flags().IntVar(&opts.Page, "page", 0, "page (0-indexed)")
	cmd.Flags().IntVar(&opts.PerPage, "per-page", 0, "page size")
	return cmd
}

func issuesGet(g *cli.Globals) *cobra.Command {
	return &cobra.Command{
		Use:   "get <group-id>",
		Short: "Get details for an issue group",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			c, err := g.Client()
			if err != nil {
				return err
			}
			var raw any
			if err := c.Get(cmd.Context(), "/open-issue-groups/"+args[0], nil, &raw); err != nil {
				return err
			}
			return g.Renderer().Render(raw)
		},
	}
}

func issuesExport(g *cli.Globals) *cobra.Command {
	var format string
	cmd := &cobra.Command{
		Use:   "export",
		Short: "Export all issues (paginates server-side)",
		RunE: func(cmd *cobra.Command, args []string) error {
			c, err := g.Client()
			if err != nil {
				return err
			}
			q := map[string]string{}
			if format != "" {
				q["format"] = format
			}
			body, _, err := c.GetRaw(cmd.Context(), "/issues/export", q)
			if err != nil {
				return err
			}
			fmt.Fprint(g.Renderer().Out, string(body))
			return nil
		},
	}
	cmd.Flags().StringVar(&format, "format", "json", "json|csv")
	return cmd
}
