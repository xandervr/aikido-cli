package commands

import (
	"github.com/spf13/cobra"
	"github.com/xandervr/aikido-cli/internal/cli"
)

func NewTasks(g *cli.Globals) *cobra.Command {
	cmd := &cobra.Command{Use: "tasks", Short: "Task tracker integrations"}
	cmd.AddCommand(
		simpleList(g, "projects", "List task tracking projects", "/tasks/projects"),
		simpleGet(g, "list <project-id>", "List tasks in a project", "/tasks/projects"),
	)
	return cmd
}
