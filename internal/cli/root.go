package cli

import (
	"github.com/spf13/cobra"
	"github.com/xandervr/aikido-cli/internal/version"
)

func NewRoot() (*cobra.Command, *Globals) {
	g := &Globals{}
	root := &cobra.Command{
		Use:           "aikido",
		Short:         "Aikido Security CLI",
		Long:          "Command-line client for the Aikido Security public REST API.",
		SilenceUsage:  true,
		SilenceErrors: true,
		Version:       version.Version,
	}
	root.PersistentFlags().BoolVar(&g.JSON, "json", false, "force JSON output")
	root.PersistentFlags().BoolVar(&g.Table, "table", false, "force table output")
	root.PersistentFlags().BoolVar(&g.NoColor, "no-color", false, "disable ANSI colors")
	root.PersistentFlags().BoolVar(&g.Debug, "debug", false, "log HTTP requests to stderr")
	root.PersistentFlags().StringVar(&g.BaseURL, "base-url", "", "override base URL (also: AIKIDO_BASE_URL)")
	root.PersistentFlags().StringVar(&g.APIKey, "api-key", "", "override API key (also: AIKIDO_API_KEY)")
	return root, g
}
