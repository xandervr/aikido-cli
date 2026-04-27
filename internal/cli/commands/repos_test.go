package commands

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/xandervr/aikido-cli/internal/client"
	"github.com/xandervr/aikido-cli/internal/output"
)

func TestReposList_RendersJSON(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if !strings.HasSuffix(r.URL.Path, "/repositories/code") {
			t.Errorf("path: %q", r.URL.Path)
		}
		_ = json.NewEncoder(w).Encode([]map[string]any{
			{"id": 1, "name": "alpha"},
			{"id": 2, "name": "beta"},
		})
	}))
	defer srv.Close()

	buf := &bytes.Buffer{}
	r := &output.Renderer{Out: buf, ForceJSON: true}
	c := client.New(client.Config{BaseURL: srv.URL, APIKey: "k"})

	if err := runReposList(context.Background(), c, r, reposListOpts{}); err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(buf.String(), "alpha") || !strings.Contains(buf.String(), "beta") {
		t.Fatalf("missing rows: %s", buf.String())
	}
}
