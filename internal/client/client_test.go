package client

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestClient_Get_SendsBearerAndDecodes(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if got := r.Header.Get("Authorization"); got != "Bearer test-key" {
			t.Errorf("bad auth header: %q", got)
		}
		if r.URL.Path != "/workspace" {
			t.Errorf("bad path: %q", r.URL.Path)
		}
		_ = json.NewEncoder(w).Encode(map[string]any{"name": "Focus"})
	}))
	defer srv.Close()

	c := New(Config{BaseURL: srv.URL, APIKey: "test-key"})
	var resp struct {
		Name string `json:"name"`
	}
	if err := c.Get(context.Background(), "/workspace", nil, &resp); err != nil {
		t.Fatal(err)
	}
	if resp.Name != "Focus" {
		t.Fatalf("bad decode: %+v", resp)
	}
}

func TestClient_Get_ReturnsAPIErrorOnNon2xx(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
		_, _ = w.Write([]byte(`{"error":"not_found","message":"missing"}`))
	}))
	defer srv.Close()

	c := New(Config{BaseURL: srv.URL, APIKey: "k"})
	err := c.Get(context.Background(), "/anything", nil, nil)
	var apiErr *APIError
	if !errors.As(err, &apiErr) {
		t.Fatalf("expected APIError, got %v", err)
	}
	if apiErr.Status != 404 || apiErr.Code != "not_found" {
		t.Fatalf("bad apiErr: %+v", apiErr)
	}
}

func TestClient_Get_AppendsQueryParams(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if got := r.URL.Query().Get("severity"); got != "high" {
			t.Errorf("bad query: %q", r.URL.RawQuery)
		}
		_, _ = w.Write([]byte(`[]`))
	}))
	defer srv.Close()

	c := New(Config{BaseURL: srv.URL, APIKey: "k"})
	var out []any
	if err := c.Get(context.Background(), "/x", map[string]string{"severity": "high"}, &out); err != nil {
		t.Fatal(err)
	}
}

func TestClient_DebugLogsToWriter(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_, _ = w.Write([]byte("ok"))
	}))
	defer srv.Close()

	var buf strings.Builder
	c := New(Config{BaseURL: srv.URL, APIKey: "secret-token", Debug: true, DebugOut: &buf})
	_ = c.Get(context.Background(), "/x", nil, nil)
	if !strings.Contains(buf.String(), "GET") {
		t.Fatalf("no debug output: %q", buf.String())
	}
	if strings.Contains(buf.String(), "secret-token") {
		t.Fatalf("debug leaked the token: %q", buf.String())
	}
}
