package commands

import (
	"context"
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/xandervr/aikido-cli/internal/cli"
)

func NewWorkspace(g *cli.Globals) *cobra.Command {
	cmd := &cobra.Command{Use: "workspace", Short: "Workspace info and configuration"}
	cmd.AddCommand(workspaceInfo(g), workspaceConfigErrors(g), workspaceIntrospect(g))
	return cmd
}

func workspaceInfo(g *cli.Globals) *cobra.Command {
	return &cobra.Command{
		Use:   "info",
		Short: "Show workspace summary",
		RunE: func(cmd *cobra.Command, args []string) error {
			c, err := g.Client()
			if err != nil {
				return err
			}
			var raw any
			if err := c.Get(context.Background(), "/workspace", nil, &raw); err != nil {
				return err
			}
			return g.Renderer().Render(raw)
		},
	}
}

func workspaceConfigErrors(g *cli.Globals) *cobra.Command {
	return &cobra.Command{
		Use:   "config-errors",
		Short: "List workspace configuration errors",
		RunE: func(cmd *cobra.Command, args []string) error {
			c, err := g.Client()
			if err != nil {
				return err
			}
			var raw any
			if err := c.Get(context.Background(), "/workspace/configuration-errors", nil, &raw); err != nil {
				return err
			}
			return g.Renderer().Render(raw)
		},
	}
}

func workspaceIntrospect(g *cli.Globals) *cobra.Command {
	return &cobra.Command{
		Use:   "introspect",
		Short: "Dump the OpenAPI spec from the workspace",
		RunE: func(cmd *cobra.Command, args []string) error {
			c, err := g.Client()
			if err != nil {
				return err
			}
			body, _, err := c.GetRaw(context.Background(), "/openapi/spec", nil)
			if err != nil {
				return err
			}
			_, err = fmt.Fprint(os.Stdout, string(body))
			return err
		},
	}
}
