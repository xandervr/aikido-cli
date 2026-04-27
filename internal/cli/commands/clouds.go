package commands

import (
	"github.com/spf13/cobra"
	"github.com/xandervr/aikido-cli/internal/cli"
)

func NewClouds(g *cli.Globals) *cobra.Command {
	cmd := &cobra.Command{Use: "clouds", Short: "Connected cloud environments"}
	cmd.AddCommand(
		simpleList(g, "list", "List connected clouds", "/clouds"),
		cloudsAssets(g),
	)
	return cmd
}

func cloudsAssets(g *cli.Globals) *cobra.Command {
	return &cobra.Command{
		Use:   "assets",
		Short: "List cloud assets (POST under the hood)",
		RunE: func(cmd *cobra.Command, args []string) error {
			c, err := g.Client()
			if err != nil {
				return err
			}
			var raw any
			if err := c.Post(cmd.Context(), "/clouds/assets", map[string]any{}, &raw); err != nil {
				return err
			}
			return g.Renderer().Render(raw)
		},
	}
}
