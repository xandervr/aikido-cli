package commands

import (
	"context"
	"errors"
	"fmt"
	"os"
	"strings"

	"github.com/spf13/cobra"
	"golang.org/x/term"

	"github.com/xandervr/aikido-cli/internal/auth"
	"github.com/xandervr/aikido-cli/internal/cli"
)

func NewAuth(g *cli.Globals) *cobra.Command {
	cmd := &cobra.Command{Use: "auth", Short: "Manage Aikido API credentials"}
	cmd.AddCommand(authLogin(g), authLogout(), authStatus(g))
	return cmd
}

func authLogin(g *cli.Globals) *cobra.Command {
	var keyFlag string
	cmd := &cobra.Command{
		Use:   "login",
		Short: "Verify an API key and store it in the OS keychain",
		RunE: func(cmd *cobra.Command, args []string) error {
			key := keyFlag
			if key == "" {
				key = os.Getenv("AIKIDO_API_KEY")
			}
			if key == "" {
				fmt.Fprint(os.Stderr, "Aikido API key: ")
				b, err := term.ReadPassword(int(os.Stdin.Fd()))
				fmt.Fprintln(os.Stderr)
				if err != nil {
					return err
				}
				key = strings.TrimSpace(string(b))
			}
			if key == "" {
				return errors.New("no API key provided")
			}
			g.APIKey = key
			c, err := g.Client()
			if err != nil {
				return err
			}
			var ws struct {
				Name string `json:"name"`
				ID   int    `json:"id"`
			}
			if err := c.Get(context.Background(), "/workspace", nil, &ws); err != nil {
				return fmt.Errorf("verify failed: %w", err)
			}
			if err := auth.NewCredentialStore().Save(key); err != nil {
				return fmt.Errorf("save to keychain: %w", err)
			}
			claims, _ := auth.DecodeClaims(key)
			fmt.Fprintf(os.Stderr, "✓ Verified workspace %q (user_id %d, region %s)\n", ws.Name, claims.UserID, defaultStr(claims.Region, "unknown"))
			fmt.Fprintln(os.Stderr, "✓ Stored in OS keychain")
			return nil
		},
	}
	cmd.Flags().StringVar(&keyFlag, "key", "", "API key (otherwise read from env or prompt)")
	return cmd
}

func authLogout() *cobra.Command {
	return &cobra.Command{
		Use:   "logout",
		Short: "Delete the stored API key",
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := auth.NewCredentialStore().Delete(); err != nil {
				return err
			}
			fmt.Fprintln(os.Stderr, "✓ Removed credential")
			return nil
		},
	}
}

func authStatus(g *cli.Globals) *cobra.Command {
	return &cobra.Command{
		Use:   "status",
		Short: "Show current authentication state",
		RunE: func(cmd *cobra.Command, args []string) error {
			source, key := resolveSource(g)
			if key == "" {
				fmt.Fprintln(os.Stderr, "not authenticated")
				return &cli.ExitError{Code: cli.ExitAuth, Err: errors.New("no credential")}
			}
			claims, _ := auth.DecodeClaims(key)
			out := struct {
				Source     string `json:"source"      aikido:"column,header=Source"`
				Masked     string `json:"masked_key"  aikido:"column,header=Key"`
				Region     string `json:"region"      aikido:"column,header=Region"`
				UserID     int    `json:"user_id"     aikido:"column,header=UserID"`
				ExpiryUnix int64  `json:"expires_at"  aikido:"column,header=Expires"`
			}{
				Source:     source,
				Masked:     mask(key),
				Region:     defaultStr(claims.Region, "unknown"),
				UserID:     claims.UserID,
				ExpiryUnix: claims.Exp,
			}
			return g.Renderer().Render(out)
		},
	}
}

func resolveSource(g *cli.Globals) (string, string) {
	if g.APIKey != "" {
		return "flag", g.APIKey
	}
	if v := os.Getenv("AIKIDO_API_KEY"); v != "" {
		return "env", v
	}
	if v, err := auth.NewCredentialStore().Load(); err == nil {
		return "keychain", v
	}
	return "none", ""
}

func mask(s string) string {
	if len(s) <= 12 {
		return "***"
	}
	return s[:8] + "..." + s[len(s)-4:]
}

func defaultStr(s, fallback string) string {
	if s == "" {
		return fallback
	}
	return s
}
