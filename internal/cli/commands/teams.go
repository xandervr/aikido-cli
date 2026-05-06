package commands

import (
	"encoding/json"
	"errors"
	"fmt"
	"strconv"

	"github.com/spf13/cobra"
	"github.com/xandervr/aikido-cli/internal/cli"
)

type Team struct {
	ID     int64  `aikido:"column,header=ID"`
	Name   string `aikido:"column,header=Name"`
	Source string `aikido:"column,header=Source"`
	Active string `aikido:"column,header=Active"`

	raw json.RawMessage
}

func (t *Team) UnmarshalJSON(b []byte) error {
	t.raw = append(t.raw[:0], b...)
	var m map[string]any
	if err := json.Unmarshal(b, &m); err != nil {
		return err
	}
	t.ID = pickInt(m, "id", "team_id")
	t.Name = pickStr(m, "name", "team_name")
	t.Source = pickStr(m, "external_source", "source")
	if t.Source == "" {
		if imported, ok := m["is_imported"].(bool); ok && imported {
			t.Source = "imported"
		}
	}
	if active, ok := m["active"].(bool); ok {
		t.Active = strconv.FormatBool(active)
	}
	return nil
}

func (t Team) MarshalJSON() ([]byte, error) {
	if len(t.raw) > 0 {
		return t.raw, nil
	}
	m := map[string]any{
		"id":   t.ID,
		"name": t.Name,
	}
	if t.Source != "" {
		m["external_source"] = t.Source
	}
	if t.Active != "" {
		if active, err := strconv.ParseBool(t.Active); err == nil {
			m["active"] = active
		}
	}
	return json.Marshal(m)
}

func NewTeams(g *cli.Globals) *cobra.Command {
	cmd := &cobra.Command{Use: "teams", Short: "Team management"}
	cmd.AddCommand(
		teamsList(g),
		teamsCreate(g),
		teamsUpdate(g),
		teamsDelete(g),
		teamsLink(g),
		teamsUnlink(g),
		teamsRemoveUser(g),
	)
	return cmd
}

func teamsList(g *cli.Globals) *cobra.Command {
	var page int
	cmd := &cobra.Command{
		Use:   "list",
		Short: "List teams",
		RunE: func(cmd *cobra.Command, args []string) error {
			c, err := g.Client()
			if err != nil {
				return err
			}
			q := map[string]string{}
			if page > 0 {
				q["page"] = fmt.Sprintf("%d", page)
			}
			var teams []Team
			if err := c.Get(cmd.Context(), "/teams", q, &teams); err != nil {
				return err
			}
			return g.Renderer().Render(teams)
		},
	}
	cmd.Flags().IntVar(&page, "page", 0, "page (0-indexed)")
	return cmd
}

func teamsCreate(g *cli.Globals) *cobra.Command {
	var name string
	cmd := &cobra.Command{
		Use:   "create",
		Short: "Create a team",
		RunE: func(cmd *cobra.Command, args []string) error {
			if name == "" {
				return errors.New("--name is required")
			}
			c, err := g.Client()
			if err != nil {
				return err
			}
			var resp any
			if err := c.Post(cmd.Context(), "/teams", map[string]string{"name": name}, &resp); err != nil {
				return err
			}
			return g.Renderer().Render(resp)
		},
	}
	cmd.Flags().StringVar(&name, "name", "", "team name (required)")
	_ = cmd.MarkFlagRequired("name")
	return cmd
}

func teamsUpdate(g *cli.Globals) *cobra.Command {
	var name string
	cmd := &cobra.Command{
		Use:   "update <id>",
		Short: "Update a team (rename)",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			body := map[string]any{}
			if name != "" {
				body["name"] = name
			}
			if len(body) == 0 {
				return errors.New("nothing to update; provide --name")
			}
			c, err := g.Client()
			if err != nil {
				return err
			}
			var resp any
			if err := c.Put(cmd.Context(), "/teams/"+args[0], body, &resp); err != nil {
				return err
			}
			return g.Renderer().Render(resp)
		},
	}
	cmd.Flags().StringVar(&name, "name", "", "new team name")
	return cmd
}

func teamsDelete(g *cli.Globals) *cobra.Command {
	var confirm bool
	cmd := &cobra.Command{
		Use:   "delete <id>",
		Short: "Delete a non-imported team",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			if !confirm {
				return &cli.ExitError{Code: cli.ExitUsage, Err: errors.New("destructive: pass --confirm")}
			}
			c, err := g.Client()
			if err != nil {
				return err
			}
			if err := c.Delete(cmd.Context(), "/teams/"+args[0], nil); err != nil {
				return err
			}
			fmt.Fprintln(g.Renderer().Out, `{"deleted":true}`)
			return nil
		},
	}
	cmd.Flags().BoolVar(&confirm, "confirm", false, "required for destructive operation")
	return cmd
}

func teamsLink(g *cli.Globals) *cobra.Command {
	return &cobra.Command{
		Use:   "link <team-id> <resource-type> <resource-id>",
		Short: "Link a resource (repo|container|cloud|app|domain) to a team",
		Args:  cobra.ExactArgs(3),
		RunE: func(cmd *cobra.Command, args []string) error {
			c, err := g.Client()
			if err != nil {
				return err
			}
			body, err := teamResourceBody(args[1], args[2])
			if err != nil {
				return &cli.ExitError{Code: cli.ExitUsage, Err: err}
			}
			var resp any
			if err := c.Post(cmd.Context(), "/teams/"+args[0]+"/linkResource", body, &resp); err != nil {
				return err
			}
			return g.Renderer().Render(resp)
		},
	}
}

func teamsUnlink(g *cli.Globals) *cobra.Command {
	return &cobra.Command{
		Use:   "unlink <team-id> <resource-type> <resource-id>",
		Short: "Unlink a resource from a team",
		Args:  cobra.ExactArgs(3),
		RunE: func(cmd *cobra.Command, args []string) error {
			c, err := g.Client()
			if err != nil {
				return err
			}
			body, err := teamResourceBody(args[1], args[2])
			if err != nil {
				return &cli.ExitError{Code: cli.ExitUsage, Err: err}
			}
			var resp any
			path := fmt.Sprintf("/teams/%s/unlinkResource", args[0])
			if err := c.Post(cmd.Context(), path, body, &resp); err != nil {
				return err
			}
			return g.Renderer().Render(resp)
		},
	}
}

func teamsRemoveUser(g *cli.Globals) *cobra.Command {
	var confirm bool
	cmd := &cobra.Command{
		Use:   "remove-user <team-id> <user-id>",
		Short: "Remove a user from a team",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			if !confirm {
				return &cli.ExitError{Code: cli.ExitUsage, Err: errors.New("destructive: pass --confirm")}
			}
			c, err := g.Client()
			if err != nil {
				return err
			}
			userID, err := strconv.Atoi(args[1])
			if err != nil {
				return &cli.ExitError{Code: cli.ExitUsage, Err: fmt.Errorf("user-id must be an integer: %w", err)}
			}
			var resp any
			path := fmt.Sprintf("/teams/%s/removeUser", args[0])
			if err := c.Post(cmd.Context(), path, map[string]any{"user_id": userID}, &resp); err != nil {
				return err
			}
			return g.Renderer().Render(resp)
		},
	}
	cmd.Flags().BoolVar(&confirm, "confirm", false, "required for destructive operation")
	return cmd
}

func teamResourceBody(resourceType, resourceID string) (map[string]any, error) {
	fieldByType := map[string]string{
		"repo":      "repo_id",
		"cloud":     "cloud_id",
		"container": "image_id",
		"image":     "image_id",
		"domain":    "domain_id",
		"app":       "zen_app_id",
		"zen-app":   "zen_app_id",
	}
	field, ok := fieldByType[resourceType]
	if !ok {
		return nil, fmt.Errorf("unsupported resource type %q", resourceType)
	}
	id, err := strconv.Atoi(resourceID)
	if err != nil {
		return nil, fmt.Errorf("resource-id must be an integer: %w", err)
	}
	return map[string]any{field: id}, nil
}
