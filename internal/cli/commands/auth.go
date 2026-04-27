package commands

import (
	"bufio"
	"context"
	"errors"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/spf13/cobra"
	"golang.org/x/term"

	"github.com/xandervr/aikido-cli/internal/auth"
	"github.com/xandervr/aikido-cli/internal/cli"
	"github.com/xandervr/aikido-cli/internal/client"
)

func NewAuth(g *cli.Globals) *cobra.Command {
	cmd := &cobra.Command{Use: "auth", Short: "Manage Aikido API credentials"}
	cmd.AddCommand(authLogin(g), authLogout(), authStatus(g), authRefresh(g))
	return cmd
}

func authLogin(g *cli.Globals) *cobra.Command {
	var clientID, clientSecret string
	cmd := &cobra.Command{
		Use:   "login",
		Short: "Verify Aikido OAuth credentials and store them in the OS keychain",
		Long: `Exchanges client_id and client_secret for an access token via
POST /api/oauth/token (HTTP Basic auth, grant_type=client_credentials).
On success, stores both values in the OS keychain and caches the access
token on disk.`,
		RunE: func(cmd *cobra.Command, args []string) error {
			if clientID == "" {
				clientID = os.Getenv("AIKIDO_CLIENT_ID")
			}
			if clientID == "" {
				fmt.Fprint(os.Stderr, "Aikido client ID: ")
				reader := bufio.NewReader(os.Stdin)
				line, err := reader.ReadString('\n')
				if err != nil {
					return err
				}
				clientID = strings.TrimSpace(line)
			}
			if clientSecret == "" {
				clientSecret = os.Getenv("AIKIDO_CLIENT_SECRET")
			}
			if clientSecret == "" {
				fmt.Fprint(os.Stderr, "Aikido client secret: ")
				b, err := term.ReadPassword(int(os.Stdin.Fd()))
				fmt.Fprintln(os.Stderr)
				if err != nil {
					return err
				}
				clientSecret = strings.TrimSpace(string(b))
			}
			if clientID == "" || clientSecret == "" {
				return errors.New("client_id and client_secret are required")
			}

			apiBase := g.APIBase()
			oauthURL := auth.DeriveOAuthURL(apiBase)
			tok, err := auth.ExchangeClientCredentials(cmd.Context(), oauthURL, clientID, clientSecret, nil)
			if err != nil {
				return fmt.Errorf("verify failed: %w", err)
			}

			c := client.New(client.Config{BaseURL: apiBase, APIKey: tok.Token})
			var ws struct {
				Name string `json:"name"`
				ID   int    `json:"id"`
			}
			if err := c.Get(cmd.Context(), "/workspace", nil, &ws); err != nil {
				return fmt.Errorf("workspace probe: %w", err)
			}

			if err := auth.NewCredentialStore().SaveCredentials(auth.ClientCredentials{ClientID: clientID, ClientSecret: clientSecret}); err != nil {
				return fmt.Errorf("save to keychain: %w", err)
			}
			_ = auth.SaveCachedToken(tok)
			fmt.Fprintf(os.Stderr, "✓ Verified workspace %q (id %d)\n", ws.Name, ws.ID)
			fmt.Fprintln(os.Stderr, "✓ Stored client credentials in OS keychain")
			fmt.Fprintf(os.Stderr, "✓ Cached access token (valid for %s)\n", time.Until(tok.ExpiresAt).Round(time.Second))
			return nil
		},
	}
	cmd.Flags().StringVar(&clientID, "client-id", "", "OAuth client ID")
	cmd.Flags().StringVar(&clientSecret, "client-secret", "", "OAuth client secret")
	return cmd
}

func authLogout() *cobra.Command {
	return &cobra.Command{
		Use:   "logout",
		Short: "Delete stored credentials and cached access token",
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := auth.NewCredentialStore().Delete(); err != nil {
				return err
			}
			_ = auth.ClearCachedToken()
			fmt.Fprintln(os.Stderr, "✓ Removed credentials and cached token")
			return nil
		},
	}
}

func authStatus(g *cli.Globals) *cobra.Command {
	return &cobra.Command{
		Use:   "status",
		Short: "Show current authentication state",
		RunE: func(cmd *cobra.Command, args []string) error {
			source := g.CredentialSource()
			if source == "none" {
				fmt.Fprintln(os.Stderr, "not authenticated")
				return &cli.ExitError{Code: cli.ExitAuth, Err: errors.New("no credential")}
			}

			var clientIDMasked string
			var hasSecret bool
			if source == "flag-token" || source == "env-token" {
				clientIDMasked = "(direct token)"
				hasSecret = false
			} else {
				creds, err := g.ResolveCredentials()
				if err != nil {
					return err
				}
				clientIDMasked = mask(creds.ClientID)
				hasSecret = creds.ClientSecret != ""
			}

			tokenStatus := "no cached token"
			if t, err := auth.LoadCachedToken(); err == nil {
				tokenStatus = fmt.Sprintf("cached (expires in %s)", time.Until(t.ExpiresAt).Round(time.Second))
			}

			out := struct {
				Source    string `json:"source"        aikido:"column,header=Source"`
				ClientID  string `json:"client_id"     aikido:"column,header=ClientID"`
				HasSecret bool   `json:"has_secret"    aikido:"column,header=HasSecret"`
				Token     string `json:"token_status"  aikido:"column,header=Token"`
				APIBase   string `json:"api_base_url"  aikido:"column,header=API"`
			}{
				Source:    source,
				ClientID:  clientIDMasked,
				HasSecret: hasSecret,
				Token:     tokenStatus,
				APIBase:   g.APIBase(),
			}
			return g.Renderer().Render(out)
		},
	}
}

func authRefresh(g *cli.Globals) *cobra.Command {
	return &cobra.Command{
		Use:   "refresh",
		Short: "Force a fresh OAuth token exchange (clears the cache)",
		RunE: func(cmd *cobra.Command, args []string) error {
			creds, err := g.ResolveCredentials()
			if err != nil {
				return err
			}
			_ = auth.ClearCachedToken()
			oauthURL := auth.DeriveOAuthURL(g.APIBase())
			tok, err := auth.ExchangeClientCredentials(context.Background(), oauthURL, creds.ClientID, creds.ClientSecret, nil)
			if err != nil {
				return err
			}
			if err := auth.SaveCachedToken(tok); err != nil {
				return err
			}
			fmt.Fprintf(os.Stderr, "✓ New access token cached (valid for %s)\n", time.Until(tok.ExpiresAt).Round(time.Second))
			return nil
		},
	}
}

func mask(s string) string {
	if len(s) <= 8 {
		return "***"
	}
	return s[:4] + "..." + s[len(s)-4:]
}
