package auth

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"
)

type AccessToken struct {
	Token     string    `json:"access_token"`
	TokenType string    `json:"token_type"`
	ExpiresAt time.Time `json:"expires_at"`
}

const oauthSafetyMargin = 60 * time.Second

func ExchangeClientCredentials(ctx context.Context, oauthURL, clientID, clientSecret string, httpClient *http.Client) (*AccessToken, error) {
	if httpClient == nil {
		httpClient = &http.Client{Timeout: 30 * time.Second}
	}
	body := strings.NewReader(`{"grant_type":"client_credentials"}`)
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, oauthURL, body)
	if err != nil {
		return nil, err
	}
	authHdr := base64.StdEncoding.EncodeToString([]byte(clientID + ":" + clientSecret))
	req.Header.Set("Authorization", "Basic "+authHdr)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")
	req.Header.Set("User-Agent", "aikido-cli")

	resp, err := httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	respBody, _ := io.ReadAll(resp.Body)
	if resp.StatusCode/100 != 2 {
		return nil, fmt.Errorf("oauth: %d: %s", resp.StatusCode, string(respBody))
	}
	var raw struct {
		AccessToken string `json:"access_token"`
		TokenType   string `json:"token_type"`
		ExpiresIn   int    `json:"expires_in"`
	}
	if err := json.Unmarshal(respBody, &raw); err != nil {
		return nil, fmt.Errorf("decode oauth response: %w", err)
	}
	if raw.AccessToken == "" {
		return nil, fmt.Errorf("oauth response missing access_token (body=%s)", string(respBody))
	}
	exp := time.Now().Add(time.Duration(raw.ExpiresIn) * time.Second).Add(-oauthSafetyMargin)
	return &AccessToken{Token: raw.AccessToken, TokenType: raw.TokenType, ExpiresAt: exp}, nil
}

// DeriveOAuthURL returns the canonical OAuth endpoint for a given API base URL.
// The OAuth endpoint lives at /api/oauth/token on the same host as the public API.
func DeriveOAuthURL(apiBaseURL string) string {
	u, err := url.Parse(apiBaseURL)
	if err != nil || u.Host == "" {
		return "https://app.aikido.dev/api/oauth/token"
	}
	scheme := u.Scheme
	if scheme == "" {
		scheme = "https"
	}
	return scheme + "://" + u.Host + "/api/oauth/token"
}
