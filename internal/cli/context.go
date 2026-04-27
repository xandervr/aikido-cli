package cli

import (
	"errors"
	"fmt"
	"io"
	"os"

	"github.com/xandervr/aikido-cli/internal/auth"
	"github.com/xandervr/aikido-cli/internal/client"
	"github.com/xandervr/aikido-cli/internal/output"
)

type Globals struct {
	JSON     bool
	Table    bool
	NoColor  bool
	Debug    bool
	BaseURL  string
	APIKey   string
	DebugOut io.Writer
}

func (g *Globals) Renderer() *output.Renderer {
	r := output.New()
	r.ForceJSON = g.JSON
	r.ForceTable = g.Table
	r.NoColor = g.NoColor || os.Getenv("NO_COLOR") != ""
	return r
}

func (g *Globals) Client() (*client.Client, error) {
	key, err := g.resolveKey()
	if err != nil {
		return nil, err
	}
	cfg := client.Config{
		BaseURL:  g.BaseURL,
		APIKey:   key,
		Debug:    g.Debug,
		DebugOut: g.DebugOut,
	}
	if cfg.BaseURL == "" {
		cfg.BaseURL = os.Getenv("AIKIDO_BASE_URL")
	}
	if cfg.DebugOut == nil {
		cfg.DebugOut = os.Stderr
	}
	return client.New(cfg), nil
}

func (g *Globals) resolveKey() (string, error) {
	if g.APIKey != "" {
		return g.APIKey, nil
	}
	if v := os.Getenv("AIKIDO_API_KEY"); v != "" {
		return v, nil
	}
	store := auth.NewCredentialStore()
	v, err := store.Load()
	if err == nil {
		return v, nil
	}
	if errors.Is(err, auth.ErrNoCredential) {
		return "", &ExitError{Code: ExitAuth, Err: errors.New("not authenticated. Run 'aikido auth login' or set AIKIDO_API_KEY")}
	}
	return "", &ExitError{Code: ExitAuth, Err: fmt.Errorf("read keychain: %w", err)}
}
