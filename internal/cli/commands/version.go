package commands

import (
	"encoding/json"
	"fmt"

	"github.com/spf13/cobra"
	"github.com/xandervr/aikido-cli/internal/cli"
	"github.com/xandervr/aikido-cli/internal/version"
)

func NewVersion(g *cli.Globals) *cobra.Command {
	return &cobra.Command{
		Use:   "version",
		Short: "Show CLI version",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			info := version.Current()
			if g.JSON {
				enc := json.NewEncoder(cmd.OutOrStdout())
				enc.SetIndent("", "  ")
				return enc.Encode(info)
			}
			fmt.Fprintf(cmd.OutOrStdout(), "aikido version %s\n", info.Version)
			return nil
		},
	}
}
