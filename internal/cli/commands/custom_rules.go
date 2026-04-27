package commands

import (
	"github.com/spf13/cobra"
	"github.com/xandervr/aikido-cli/internal/cli"
)

func NewCustomRules(g *cli.Globals) *cobra.Command {
	cmd := &cobra.Command{Use: "custom-rules", Short: "Custom SAST rules"}
	cmd.AddCommand(
		simpleList(g, "list", "List custom rules", "/custom-rules"),
		simpleGet(g, "get <id>", "Get a custom rule", "/custom-rules"),
	)
	return cmd
}
