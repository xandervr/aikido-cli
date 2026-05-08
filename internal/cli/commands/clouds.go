package commands

import (
	"net/http"

	"github.com/spf13/cobra"
	"github.com/xandervr/aikido-cli/internal/cli"
)

func NewClouds(g *cli.Globals) *cobra.Command {
	cmd := &cobra.Command{Use: "clouds", Short: "Connected cloud environments"}
	cmd.AddCommand(
		simpleList(g, "list", "List connected clouds", "/clouds"),
		endpointCommand(g, endpointCommandConfig{Use: "rules", Short: "List cloud rules", Method: http.MethodGet, Path: staticPath("/clouds/rules")}),
		cloudsAssets(g),
		endpointCommand(g, endpointCommandConfig{Use: "delete <cloud-id>", Short: "Remove cloud", Method: http.MethodDelete, Args: cobra.ExactArgs(1), Path: oneArgPath("/clouds/%s"), Confirm: true}),
		endpointCommand(g, endpointCommandConfig{Use: "aws", Short: "Connect AWS cloud", Method: http.MethodPost, Path: staticPath("/clouds/aws")}),
		endpointCommand(g, endpointCommandConfig{Use: "azure", Short: "Connect Azure cloud", Method: http.MethodPost, Path: staticPath("/clouds/azure")}),
		endpointCommand(g, endpointCommandConfig{Use: "azure-credentials <cloud-id>", Short: "Update Azure cloud credentials", Method: http.MethodPut, Args: cobra.ExactArgs(1), Path: oneArgPath("/clouds/azure/%s/credentials")}),
		endpointCommand(g, endpointCommandConfig{Use: "gcp", Short: "Connect GCP cloud", Method: http.MethodPost, Path: staticPath("/clouds/gcp")}),
		endpointCommand(g, endpointCommandConfig{Use: "kubernetes", Short: "Create Kubernetes cloud", Method: http.MethodPost, Path: staticPath("/clouds/kubernetes")}),
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
