package auth

import (
	"context"
	"encoding/base64"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"
)

func TestDeriveOAuthURL(t *testing.T) {
	cases := map[string]string{
		"https://app.aikido.dev/api/public/v1": "https://app.aikido.dev/api/oauth/token",
		"http://localhost:9999/api/public/v1":  "http://localhost:9999/api/oauth/token",
		"":                                     "https://app.aikido.dev/api/oauth/token",
	}
	for in, want := range cases {
		if got := DeriveOAuthURL(in); got != want {
			t.Errorf("DeriveOAuthURL(%q) = %q, want %q", in, got, want)
		}
	}
}

func TestExchangeClientCredentials_HappyPath(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "POST" || r.URL.Path != "/api/oauth/token" {
			t.Errorf("bad request line: %s %s", r.Method, r.URL.Path)
		}
		got := r.Header.Get("Authorization")
		want := "Basic " + base64.StdEncoding.EncodeToString([]byte("id:secret"))
		if got != want {
			t.Errorf("auth header = %q, want %q", got, want)
		}
		buf := make([]byte, 256)
		n, _ := r.Body.Read(buf)
		if !strings.Contains(string(buf[:n]), `"grant_type":"client_credentials"`) {
			t.Errorf("body missing grant_type: %q", string(buf[:n]))
		}
		_, _ = w.Write([]byte(`{"access_token":"abc","token_type":"Bearer","expires_in":3600}`))
	}))
	defer srv.Close()

	tok, err := ExchangeClientCredentials(context.Background(), srv.URL+"/api/oauth/token", "id", "secret", nil)
	if err != nil {
		t.Fatal(err)
	}
	if tok.Token != "abc" {
		t.Errorf("token: %q", tok.Token)
	}
	if !tok.ExpiresAt.After(time.Now()) {
		t.Errorf("expires_at not in future: %v", tok.ExpiresAt)
	}
}

func TestExchangeClientCredentials_ErrorOnNon2xx(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusUnauthorized)
		_, _ = w.Write([]byte(`{"error":"invalid_client"}`))
	}))
	defer srv.Close()
	if _, err := ExchangeClientCredentials(context.Background(), srv.URL, "x", "y", nil); err == nil {
		t.Fatal("expected error")
	}
}
