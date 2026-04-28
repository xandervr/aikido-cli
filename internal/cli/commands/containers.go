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
		simpleList(g, "list", "List container repositories", "/containers"),
		simpleGet(g, "get <id>", "Get a container repo", "/containers"),
		containersSBOM(g),
	)
	return cmd
}

func containersSBOM(g *cli.Globals) *cobra.Command {
	var format string
	cmd := &cobra.Command{
		Use:   "sbom <id>",
		Short: "Export the SBOM for a container",
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
			body, _, err := c.GetRaw(context.Background(), "/containers/"+args[0]+"/licenses/export", q)
			if err != nil {
				return err
			}
			fmt.Fprint(g.Renderer().Out, string(body))
			return nil
		},
	}
	cmd.Flags().StringVar(&format, "format", "", "format passthrough")
	return cmd
}
