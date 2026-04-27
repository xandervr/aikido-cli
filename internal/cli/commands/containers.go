package commands

import (
	"context"
	"fmt"

	"github.com/spf13/cobra"
	"github.com/xandervr/aikido-cli/internal/cli"
)

func NewContainers(g *cli.Globals) *cobra.Command {
	cmd := &cobra.Command{Use: "containers", Short: "Container repositories"}
	cmd.AddCommand(
		simpleList(g, "list", "List container repositories", "/repositories/container"),
		simpleGet(g, "get <id>", "Get a container repo", "/repositories/container"),
		containersSBOM(g),
	)
	return cmd
}

func containersSBOM(g *cli.Globals) *cobra.Command {
	return &cobra.Command{
		Use:   "sbom <id>",
		Short: "Export the SBOM for a container",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			c, err := g.Client()
			if err != nil {
				return err
			}
			body, _, err := c.GetRaw(context.Background(), "/repositories/container/"+args[0]+"/licenses", nil)
			if err != nil {
				return err
			}
			fmt.Fprint(g.Renderer().Out, string(body))
			return nil
		},
	}
}
