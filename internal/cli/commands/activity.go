package commands

import (
	"github.com/spf13/cobra"
	"github.com/xandervr/aikido-cli/internal/cli"
)

func NewActivity(g *cli.Globals) *cobra.Command {
	var from, to, user string
	cmd := &cobra.Command{
		Use:   "activity",
		Short: "Workspace activity log",
		RunE: func(cmd *cobra.Command, args []string) error {
			c, err := g.Client()
			if err != nil {
				return err
			}
			q := map[string]string{}
			if from != "" {
				q["start"] = from
			}
			if to != "" {
				q["end"] = to
			}
			if user != "" {
				q["user_type_filter"] = user
			}
			var raw any
			if err := c.Get(cmd.Context(), "/report/activityLog", q, &raw); err != nil {
				return err
			}
			return g.Renderer().Render(raw)
		},
	}
	cmd.Flags().StringVar(&from, "from", "", "ISO date (inclusive)")
	cmd.Flags().StringVar(&to, "to", "", "ISO date (inclusive)")
	cmd.Flags().StringVar(&user, "user", "", "filter by user type")
	return cmd
}
