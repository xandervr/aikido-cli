package commands

import (
	"context"
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/xandervr/aikido-cli/internal/cli"
)

func NewReport(g *cli.Globals) *cobra.Command {
	var out string
	cmd := &cobra.Command{Use: "report", Short: "Workspace reports"}
	pdf := &cobra.Command{
		Use:   "pdf",
		Short: "Export workspace report as PDF",
		RunE: func(cmd *cobra.Command, args []string) error {
			c, err := g.Client()
			if err != nil {
				return err
			}
			body, _, err := c.GetRaw(context.Background(), "/report/export/pdf", nil)
			if err != nil {
				return err
			}
			if out != "" {
				return os.WriteFile(out, body, 0o644)
			}
			_, err = fmt.Fprint(os.Stdout, string(body))
			return err
		},
	}
	pdf.Flags().StringVar(&out, "out", "", "write PDF to this path instead of stdout")
	cmd.AddCommand(pdf)
	return cmd
}
