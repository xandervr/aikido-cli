package commands

import (
	"fmt"
	"strconv"
	"time"

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
				ts, err := activityTimestamp(from, false)
				if err != nil {
					return &cli.ExitError{Code: cli.ExitUsage, Err: fmt.Errorf("--from: %w", err)}
				}
				q["start"] = ts
			}
			if to != "" {
				ts, err := activityTimestamp(to, true)
				if err != nil {
					return &cli.ExitError{Code: cli.ExitUsage, Err: fmt.Errorf("--to: %w", err)}
				}
				q["end"] = ts
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
	cmd.Flags().StringVar(&from, "from", "", "Unix timestamp, RFC3339 timestamp, or YYYY-MM-DD start date")
	cmd.Flags().StringVar(&to, "to", "", "Unix timestamp, RFC3339 timestamp, or YYYY-MM-DD end date")
	cmd.Flags().StringVar(&user, "user", "", "filter by user type")
	return cmd
}

func activityTimestamp(value string, endOfDay bool) (string, error) {
	if _, err := strconv.ParseInt(value, 10, 64); err == nil {
		return value, nil
	}
	if ts, err := time.Parse(time.RFC3339, value); err == nil {
		return strconv.FormatInt(ts.Unix(), 10), nil
	}
	day, err := time.Parse("2006-01-02", value)
	if err != nil {
		return "", fmt.Errorf("expected Unix timestamp, RFC3339 timestamp, or YYYY-MM-DD date")
	}
	if endOfDay {
		day = day.Add(24*time.Hour - time.Second)
	}
	return strconv.FormatInt(day.Unix(), 10), nil
}
