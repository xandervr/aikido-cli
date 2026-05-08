package commands

import (
	"fmt"
	"net/http"

	"github.com/spf13/cobra"
	"github.com/xandervr/aikido-cli/internal/cli"
)

func NewTasks(g *cli.Globals) *cobra.Command {
	cmd := &cobra.Command{Use: "tasks", Short: "Task tracker integrations"}
	cmd.AddCommand(
		simpleList(g, "projects", "List task tracking projects", "/task_tracking/projects"),
		tasksList(g),
		endpointCommand(g, endpointCommandConfig{Use: "project-mapping", Short: "Get project mapping", Method: http.MethodGet, Path: staticPath("/task_tracking/projectMapping")}),
		endpointCommand(g, endpointCommandConfig{Use: "map-repos", Short: "Map code repos to task tracking projects", Method: http.MethodPost, Path: staticPath("/task_tracking/mapCodeReposToProjects")}),
		endpointCommand(g, endpointCommandConfig{Use: "link-task", Short: "Link existing task to issue", Method: http.MethodPost, Path: staticPath("/task_tracking/linkTaskToIssueGroup")}),
	)
	return cmd
}

func tasksList(g *cli.Globals) *cobra.Command {
	var search string
	cmd := &cobra.Command{
		Use:   "list <project-id>",
		Short: "List tasks in a project",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			c, err := g.Client()
			if err != nil {
				return err
			}
			q := map[string]string{}
			if search != "" {
				q["search"] = search
			}
			var raw any
			path := fmt.Sprintf("/task_tracking/projects/%s/tasks", args[0])
			if err := c.Get(cmd.Context(), path, q, &raw); err != nil {
				return err
			}
			return g.Renderer().Render(raw)
		},
	}
	cmd.Flags().StringVar(&search, "search", "", "search tasks")
	return cmd
}
