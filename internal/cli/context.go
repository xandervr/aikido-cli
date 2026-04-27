package cli

import (
	"context"
	"errors"
	"fmt"
	"io"
	"os"

	"github.com/xandervr/aikido-cli/internal/auth"
	"github.com/xandervr/aikido-cli/internal/client"
	"github.com/xandervr/aikido-cli/internal/output"
)

type Globals struct {
	JSON         bool
	Table        bool
	NoColor      bool
	Debug        bool
	BaseURL      string
	ClientID     string
	ClientSecret string
	AccessToken  string
	DebugOut     io.Writer
}

func (g *Globals) Renderer() *output.Renderer {
	r := output.New()
	r.ForceJSON = g.JSON
	r.ForceTable = g.Table
	r.NoColor = g.NoColor || os.Getenv("NO_COLOR") != ""
	return r
}

// APIBase returns the API base URL (flag → env → default).
func (g *Globals) APIBase() string {
	if g.BaseURL != "" {
		return g.BaseURL
	}
	if v := os.Getenv("AIKIDO_BASE_URL"); v != "" {
		return v
	}
	return client.DefaultBaseURL
}

// Client returns an HTTP client with a Bearer access token attached.
// It resolves auth in this order:
//  1. --access-token flag or AIKIDO_ACCESS_TOKEN env (used directly)
//  2. cached access token on disk (if not expired)
//  3. exchange client_id/client_secret for a fresh access token
func (g *Globals) Client() (*client.Client, error) {
	token, err := g.resolveAccessToken()
	if err != nil {
		return nil, err
	}
	cfg := client.Config{
		BaseURL:  g.APIBase(),
		APIKey:   token,
		Debug:    g.Debug,
		DebugOut: g.DebugOut,
	}
	if cfg.DebugOut == nil {
		cfg.DebugOut = os.Stderr
	}
	return client.New(cfg), nil
}

func (g *Globals) resolveAccessToken() (string, error) {
	if g.AccessToken != "" {
		return g.AccessToken, nil
	}
	if v := os.Getenv("AIKIDO_ACCESS_TOKEN"); v != "" {
		return v, nil
	}
	creds, err := g.ResolveCredentials()
	if err != nil {
		return "", err
	}
	if cached, cerr := auth.LoadCachedToken(); cerr == nil {
		return cached.Token, nil
	}
	oauthURL := auth.DeriveOAuthURL(g.APIBase())
	tok, err := auth.ExchangeClientCredentials(context.Background(), oauthURL, creds.ClientID, creds.ClientSecret, nil)
	if err != nil {
		return "", &ExitError{Code: ExitAuth, Err: fmt.Errorf("oauth: %w", err)}
	}
	_ = auth.SaveCachedToken(tok)
	return tok.Token, nil
}

// ResolveCredentials returns the client_id/client_secret pair (flags → env → keychain).
// The "source" string identifies which path was used (for diagnostics).
func (g *Globals) ResolveCredentials() (auth.ClientCredentials, error) {
	if g.ClientID != "" && g.ClientSecret != "" {
		return auth.ClientCredentials{ClientID: g.ClientID, ClientSecret: g.ClientSecret}, nil
	}
	id := os.Getenv("AIKIDO_CLIENT_ID")
	secret := os.Getenv("AIKIDO_CLIENT_SECRET")
	if id != "" && secret != "" {
		return auth.ClientCredentials{ClientID: id, ClientSecret: secret}, nil
	}
	creds, err := auth.NewCredentialStore().LoadCredentials()
	if err == nil {
		return creds, nil
	}
	if errors.Is(err, auth.ErrNoCredential) {
		return auth.ClientCredentials{}, &ExitError{Code: ExitAuth, Err: errors.New("not authenticated. Run 'aikido auth login' or set AIKIDO_CLIENT_ID + AIKIDO_CLIENT_SECRET")}
	}
	return auth.ClientCredentials{}, &ExitError{Code: ExitAuth, Err: fmt.Errorf("read keychain: %w", err)}
}

// CredentialSource returns a label for where the active credentials came from.
func (g *Globals) CredentialSource() string {
	if g.AccessToken != "" {
		return "flag-token"
	}
	if os.Getenv("AIKIDO_ACCESS_TOKEN") != "" {
		return "env-token"
	}
	if g.ClientID != "" && g.ClientSecret != "" {
		return "flag"
	}
	if os.Getenv("AIKIDO_CLIENT_ID") != "" && os.Getenv("AIKIDO_CLIENT_SECRET") != "" {
		return "env"
	}
	if _, err := auth.NewCredentialStore().LoadCredentials(); err == nil {
		return "keychain"
	}
	return "none"
}
