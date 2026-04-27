package commands

import (
	"context"
	"fmt"

	"github.com/spf13/cobra"
	"github.com/xandervr/aikido-cli/internal/cli"
)

func NewVMs(g *cli.Globals) *cobra.Command {
	cmd := &cobra.Command{Use: "vms", Short: "Virtual machines"}
	cmd.AddCommand(
		simpleList(g, "list", "List virtual machines", "/virtual-machines"),
		vmsSBOM(g),
	)
	return cmd
}

func vmsSBOM(g *cli.Globals) *cobra.Command {
	return &cobra.Command{
		Use:   "sbom <id>",
		Short: "Export the SBOM for a virtual machine",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			c, err := g.Client()
			if err != nil {
				return err
			}
			body, _, err := c.GetRaw(context.Background(), "/virtual-machines/"+args[0]+"/sbom", nil)
			if err != nil {
				return err
			}
			fmt.Fprint(g.Renderer().Out, string(body))
			return nil
		},
	}
}
