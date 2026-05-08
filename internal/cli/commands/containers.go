package commands

import (
	"context"
	"fmt"
	"net/http"

	"github.com/spf13/cobra"
	"github.com/xandervr/aikido-cli/internal/cli"
)

func NewContainers(g *cli.Globals) *cobra.Command {
	cmd := &cobra.Command{Use: "containers", Short: "Container repositories"}
	cmd.AddCommand(
		simpleList(g, "list", "List container repositories", "/containers"),
		simpleGet(g, "get <id>", "Get a container repo", "/containers"),
		containersSBOM(g),
		endpointCommand(g, endpointCommandConfig{Use: "delete <id>", Short: "Delete container", Method: http.MethodDelete, Args: cobra.ExactArgs(1), Path: oneArgPath("/containers/%s"), Confirm: true}),
		endpointCommand(g, endpointCommandConfig{Use: "raw-sbom <id>", Short: "Export raw SBOM", Method: http.MethodGet, Args: cobra.ExactArgs(1), Path: oneArgPath("/containers/%s/sbom/exportRaw")}),
		endpointCommand(g, endpointCommandConfig{Use: "sensitivity <id>", Short: "Update sensitivity", Method: http.MethodPut, Args: cobra.ExactArgs(1), Path: oneArgPath("/containers/%s/sensitivity")}),
		endpointCommand(g, endpointCommandConfig{Use: "connectivity <id>", Short: "Update internet connectivity", Method: http.MethodPut, Args: cobra.ExactArgs(1), Path: oneArgPath("/containers/%s/internetConnection")}),
		endpointCommand(g, endpointCommandConfig{Use: "upload-sbom", Short: "Upload container SBOM", Method: http.MethodPost, Path: staticPath("/containers/sbom")}),
		endpointCommand(g, endpointCommandConfig{Use: "generate-sbom", Short: "Generate bulk SBOM", Method: http.MethodPost, Path: staticPath("/containers/sbom/generate")}),
		endpointCommand(g, endpointCommandConfig{Use: "activate", Short: "Activate container", Method: http.MethodPost, Path: staticPath("/containers/activate")}),
		endpointCommand(g, endpointCommandConfig{Use: "deactivate", Short: "Deactivate container", Method: http.MethodPost, Path: staticPath("/containers/deactivate")}),
		endpointCommand(g, endpointCommandConfig{Use: "link-code-repo", Short: "Link code repository to container", Method: http.MethodPost, Path: staticPath("/containers/linkCodeRepo")}),
		endpointCommand(g, endpointCommandConfig{Use: "unlink-code-repo", Short: "Unlink code repository from container", Method: http.MethodPost, Path: staticPath("/containers/unlinkCodeRepo")}),
		endpointCommand(g, endpointCommandConfig{Use: "update-tag-filter", Short: "Update container tag filter", Method: http.MethodPost, Path: staticPath("/containers/updateTagFilter")}),
		endpointCommand(g, endpointCommandConfig{Use: "public", Short: "Add public container", Method: http.MethodPost, Path: staticPath("/containers/public")}),
		endpointCommand(g, endpointCommandConfig{Use: "clone", Short: "Clone container", Method: http.MethodPost, Path: staticPath("/containers/clone")}),
		endpointCommand(g, endpointCommandConfig{Use: "scan <id>", Short: "Scan container", Method: http.MethodPost, Args: cobra.ExactArgs(1), Path: oneArgPath("/containers/%s/scan")}),
		endpointCommand(g, endpointCommandConfig{Use: "registry <id>", Short: "Get container registry", Method: http.MethodGet, Args: cobra.ExactArgs(1), Path: oneArgPath("/containers/registries/%s")}),
		endpointCommand(g, endpointCommandConfig{Use: "registry-acr", Short: "Add Azure container registry", Method: http.MethodPost, Path: staticPath("/containers/registries/acr")}),
		endpointCommand(g, endpointCommandConfig{Use: "registry-gcp", Short: "Add GCP Artifact Registry", Method: http.MethodPost, Path: staticPath("/containers/registries/gcp-artifact-registry")}),
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
