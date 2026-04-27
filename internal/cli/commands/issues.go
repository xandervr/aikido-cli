package commands

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/spf13/cobra"
	"github.com/xandervr/aikido-cli/internal/cli"
)

// IssueGroup is a schema-tolerant view of an Aikido issue group.
//
// The Aikido response shape varies a bit by issue type (SAST vs open-source
// vs leaked-secret all carry slightly different fields). UnmarshalJSON probes
// a list of known field-name variants so the table renderer fills consistently.
// MarshalJSON returns the raw response so `--json` output preserves every
// field the API returned, untouched.
type IssueGroup struct {
	ID       int64  `aikido:"column,header=ID"`
	Title    string `aikido:"column,header=Title"`
	Severity string `aikido:"column,header=Severity"`
	Type     string `aikido:"column,header=Type"`
	Repo     string `aikido:"column,header=Repo"`
	Status   string `aikido:"column,header=Status"`

	raw json.RawMessage
}

func (g *IssueGroup) UnmarshalJSON(b []byte) error {
	g.raw = append(g.raw[:0], b...)
	var m map[string]any
	if err := json.Unmarshal(b, &m); err != nil {
		return err
	}
	g.ID = pickInt(m, "id", "issue_group_id", "group_id")
	g.Title = pickStr(m, "title", "name", "summary", "rule_name")
	g.Severity = pickStr(m, "severity", "severity_score", "priority")
	g.Type = pickStr(m, "type", "issue_type", "category")
	g.Repo = pickStr(m, "repo_name", "code_repo_name", "repository_name", "repo", "container_repo_name")
	if g.Repo == "" {
		if v, ok := m["code_repo"].(map[string]any); ok {
			g.Repo = pickStr(v, "name", "external_id")
		}
	}
	if g.Repo == "" {
		if v, ok := m["repository"].(map[string]any); ok {
			g.Repo = pickStr(v, "name", "external_id")
		}
	}
	if g.Repo == "" {
		if locs, ok := m["locations"].([]any); ok && len(locs) > 0 {
			// locations[] is one entry per *place* the issue was found
			// (lockfile, branch, monorepo subpath), not per repo.
			// Dedupe by repo name so the "(+N)" annotation reflects
			// distinct repos, not raw location count.
			seen := map[string]struct{}{}
			order := []string{}
			for _, raw := range locs {
				loc, ok := raw.(map[string]any)
				if !ok {
					continue
				}
				name := pickStr(loc,
					"code_repo_name", "repo_name", "repository_name",
					"container_repo_name", "name", "external_id",
				)
				if name == "" {
					if cr, ok := loc["code_repo"].(map[string]any); ok {
						name = pickStr(cr, "name", "external_id")
					}
				}
				if name == "" {
					if cr, ok := loc["repository"].(map[string]any); ok {
						name = pickStr(cr, "name", "external_id")
					}
				}
				if name == "" {
					continue
				}
				if _, dup := seen[name]; dup {
					continue
				}
				seen[name] = struct{}{}
				order = append(order, name)
			}
			if len(order) > 0 {
				g.Repo = order[0]
				if len(order) > 1 {
					g.Repo = fmt.Sprintf("%s (+%d)", g.Repo, len(order)-1)
				}
			}
		}
	}
	g.Status = pickStr(m, "group_status", "status", "state", "issue_status")
	if g.Status == "" {
		if b, ok := m["is_open"].(bool); ok {
			if b {
				g.Status = "open"
			} else {
				g.Status = "closed"
			}
		}
	}
	if g.Status == "" {
		if b, ok := m["ignored"].(bool); ok && b {
			g.Status = "ignored"
		}
	}
	if g.Status == "" {
		if b, ok := m["snoozed"].(bool); ok && b {
			g.Status = "snoozed"
		}
	}
	return nil
}

// MarshalJSON returns the raw response unchanged so `--json` output never
// loses any field.
func (g IssueGroup) MarshalJSON() ([]byte, error) {
	if len(g.raw) > 0 {
		return g.raw, nil
	}
	type alias struct {
		ID       int64  `json:"id"`
		Title    string `json:"title"`
		Severity string `json:"severity"`
		Type     string `json:"type"`
		Repo     string `json:"repo"`
		Status   string `json:"status"`
	}
	return json.Marshal(alias{g.ID, g.Title, g.Severity, g.Type, g.Repo, g.Status})
}

type issuesListOpts struct {
	Severity string
	Status   string
	Type     string
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
			if opts.Type != "" {
				q["filter_issue_type"] = opts.Type
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
			// The Aikido API does not support filter_status server-side
			// (only severity, type, repo, team). Filter client-side so the
			// flag has the effect users expect.
			if opts.Status != "" {
				want := strings.ToLower(opts.Status)
				kept := groups[:0]
				for _, gr := range groups {
					if strings.ToLower(gr.Status) == want {
						kept = append(kept, gr)
					}
				}
				groups = kept
			}
			return g.Renderer().Render(groups)
		},
	}
	cmd.Flags().StringVar(&opts.Severity, "severity", "", "filter: critical|high|medium|low")
	cmd.Flags().StringVar(&opts.Status, "status", "", "filter (CLIENT-side; API has no status filter): task_open|task_in_progress|task_closed|task_done|todo|new|ignored|snoozed")
	cmd.Flags().StringVar(&opts.Type, "type", "", "filter: open_source|leaked_secret|sast|iac|cloud|docker_container|cloud_instance|surface_monitoring|malware|eol|mobile|scm_security|ai_pentest|license")
	cmd.Flags().IntVar(&opts.Repo, "repo", 0, "filter by code repo ID")
	cmd.Flags().IntVar(&opts.Team, "team", 0, "filter by team ID")
	cmd.Flags().IntVar(&opts.Page, "page", 0, "page (0-indexed)")
	cmd.Flags().IntVar(&opts.PerPage, "per-page", 0, "page size (server caps issues at 20)")
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
			if err := c.Get(cmd.Context(), "/issues/groups/"+args[0], nil, &raw); err != nil {
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
