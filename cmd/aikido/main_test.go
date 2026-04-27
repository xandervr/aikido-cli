package main

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/xandervr/aikido-cli/internal/cli"
	"github.com/xandervr/aikido-cli/internal/cli/commands"
)

func TestCLI_ReposListEndToEnd(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Header.Get("Authorization") != "Bearer test" {
			t.Errorf("auth header: %q", r.Header.Get("Authorization"))
		}
		_ = json.NewEncoder(w).Encode([]map[string]any{
			{"id": 1, "name": "alpha"},
		})
	}))
	defer srv.Close()

	t.Setenv("AIKIDO_ACCESS_TOKEN", "test")
	t.Setenv("AIKIDO_BASE_URL", srv.URL)

	root, g := cli.NewRoot()
	root.AddCommand(commands.NewRepos(g))
	root.SetArgs([]string{"repos", "list", "--json"})
	if err := root.Execute(); err != nil {
		t.Fatal(err)
	}
}

func TestCLI_TeamsDeleteRequiresConfirm(t *testing.T) {
	t.Setenv("AIKIDO_ACCESS_TOKEN", "test")
	t.Setenv("AIKIDO_BASE_URL", "http://invalid.local")

	root, g := cli.NewRoot()
	root.AddCommand(commands.NewTeams(g))
	root.SetArgs([]string{"teams", "delete", "1"})
	err := root.Execute()
	if err == nil {
		t.Fatal("expected error when --confirm is missing")
	}
}
