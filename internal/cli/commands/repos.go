package commands

import (
	"context"
	"fmt"

	"github.com/spf13/cobra"
	"github.com/xandervr/aikido-cli/internal/cli"
	"github.com/xandervr/aikido-cli/internal/client"
	"github.com/xandervr/aikido-cli/internal/output"
)

type Repo struct {
	ID       int    `json:"id"            aikido:"column,header=ID"`
	Name     string `json:"name"          aikido:"column,header=Name"`
	Provider string `json:"provider"      aikido:"column,header=Provider"`
	External string `json:"external_id"   aikido:"column,header=External"`
	IsActive bool   `json:"active"        aikido:"column,header=Active"`
	TeamID   int    `json:"team_id"       aikido:"column,header=Team"`
}

type reposListOpts struct {
	Team    int
	Search  string
	Page    int
	PerPage int
}

func NewRepos(g *cli.Globals) *cobra.Command {
	cmd := &cobra.Command{Use: "repos", Short: "Code repositories"}
	cmd.AddCommand(reposList(g), reposGet(g), reposSBOM(g))
	return cmd
}

func reposList(g *cli.Globals) *cobra.Command {
	var opts reposListOpts
	cmd := &cobra.Command{
		Use:   "list",
		Short: "List code repositories",
		RunE: func(cmd *cobra.Command, args []string) error {
			c, err := g.Client()
			if err != nil {
				return err
			}
			return runReposList(cmd.Context(), c, g.Renderer(), opts)
		},
	}
	cmd.Flags().IntVar(&opts.Team, "team", 0, "filter by team ID")
	cmd.Flags().StringVar(&opts.Search, "search", "", "name search")
	cmd.Flags().IntVar(&opts.Page, "page", 0, "page (0-indexed)")
	cmd.Flags().IntVar(&opts.PerPage, "per-page", 0, "page size")
	return cmd
}

func runReposList(ctx context.Context, c *client.Client, r *output.Renderer, opts reposListOpts) error {
	q := map[string]string{}
	if opts.Team > 0 {
		q["team_id"] = fmt.Sprintf("%d", opts.Team)
	}
	if opts.Search != "" {
		q["filter_name"] = opts.Search
	}
	if opts.Page > 0 {
		q["page"] = fmt.Sprintf("%d", opts.Page)
	}
	if opts.PerPage > 0 {
		q["per_page"] = fmt.Sprintf("%d", opts.PerPage)
	}
	var repos []Repo
	if err := c.Get(ctx, "/repositories/code", q, &repos); err != nil {
		return err
	}
	return r.Render(repos)
}

func reposGet(g *cli.Globals) *cobra.Command {
	return &cobra.Command{
		Use:   "get <id>",
		Short: "Get a single code repository",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			c, err := g.Client()
			if err != nil {
				return err
			}
			var raw any
			if err := c.Get(cmd.Context(), "/repositories/code/"+args[0], nil, &raw); err != nil {
				return err
			}
			return g.Renderer().Render(raw)
		},
	}
}

func reposSBOM(g *cli.Globals) *cobra.Command {
	var format string
	cmd := &cobra.Command{
		Use:   "sbom <id>",
		Short: "Export the SBOM (license overview) for a repo",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			c, err := g.Client()
			if err != nil {
				return err
			}
			q := map[string]string{}
			if format != "" {
				q["format"] = format
			}
			body, _, err := c.GetRaw(cmd.Context(), "/repositories/code/"+args[0]+"/licenses/export", q)
			if err != nil {
				return err
			}
			fmt.Fprint(g.Renderer().Out, string(body))
			return nil
		},
	}
	cmd.Flags().StringVar(&format, "format", "", "format passthrough (json|csv|spdx)")
	return cmd
}
