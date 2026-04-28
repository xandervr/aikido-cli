package commands

import (
	"context"
	"fmt"
	"os"
	"strconv"

	"github.com/spf13/cobra"
	"github.com/xandervr/aikido-cli/internal/cli"
)

func NewReport(g *cli.Globals) *cobra.Command {
	var out, sections string
	var teamID int
	cmd := &cobra.Command{Use: "report", Short: "Workspace reports"}
	pdf := &cobra.Command{
		Use:   "pdf",
		Short: "Export workspace report as PDF",
		RunE: func(cmd *cobra.Command, args []string) error {
			c, err := g.Client()
			if err != nil {
				return err
			}
			q := map[string]string{"included_sections": sections}
			if teamID > 0 {
				q["team_id"] = strconv.Itoa(teamID)
			}
			body, _, err := c.GetRaw(context.Background(), "/report/export/pdf", q)
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
	pdf.Flags().StringVar(&sections, "sections", "", "comma-separated sections to include (required)")
	pdf.Flags().IntVar(&teamID, "team", 0, "filter report by team ID")
	_ = pdf.MarkFlagRequired("sections")
	cmd.AddCommand(pdf)
	return cmd
}
