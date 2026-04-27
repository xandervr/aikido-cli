package commands

import (
	"github.com/spf13/cobra"
	"github.com/xandervr/aikido-cli/internal/cli"
)

func NewCompliance(g *cli.Globals) *cobra.Command {
	cmd := &cobra.Command{Use: "compliance", Short: "Compliance overviews"}
	cmd.AddCommand(
		simpleList(g, "soc2", "SOC2 compliance overview", "/compliance/soc2"),
		simpleList(g, "nis2", "NIS2 compliance overview", "/compliance/nis2"),
		simpleList(g, "iso27001", "ISO 27001 compliance overview", "/compliance/iso27001"),
	)
	return cmd
}
