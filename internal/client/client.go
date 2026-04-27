package client

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"
)

const DefaultBaseURL = "https://app.aikido.dev/api/public/v1"

type Config struct {
	BaseURL  string
	APIKey   string
	Debug    bool
	DebugOut io.Writer
	HTTP     *http.Client
}

type Client struct {
	baseURL  string
	apiKey   string
	debug    bool
	debugOut io.Writer
	http     *http.Client
}

func New(cfg Config) *Client {
	if cfg.BaseURL == "" {
		cfg.BaseURL = DefaultBaseURL
	}
	if cfg.HTTP == nil {
		cfg.HTTP = &http.Client{Timeout: 30 * time.Second}
	}
	if cfg.DebugOut == nil {
		cfg.DebugOut = io.Discard
	}
	return &Client{
		baseURL:  strings.TrimRight(cfg.BaseURL, "/"),
		apiKey:   cfg.APIKey,
		debug:    cfg.Debug,
		debugOut: cfg.DebugOut,
		http:     cfg.HTTP,
	}
}

func (c *Client) Get(ctx context.Context, path string, query map[string]string, out any) error {
	return c.do(ctx, http.MethodGet, path, query, nil, out)
}

func (c *Client) Post(ctx context.Context, path string, body, out any) error {
	return c.do(ctx, http.MethodPost, path, nil, body, out)
}

func (c *Client) Put(ctx context.Context, path string, body, out any) error {
	return c.do(ctx, http.MethodPut, path, nil, body, out)
}

func (c *Client) Delete(ctx context.Context, path string, out any) error {
	return c.do(ctx, http.MethodDelete, path, nil, nil, out)
}

func (c *Client) GetRaw(ctx context.Context, path string, query map[string]string) ([]byte, string, error) {
	req, err := c.newRequest(ctx, http.MethodGet, path, query, nil)
	if err != nil {
		return nil, "", err
	}
	resp, err := c.http.Do(req)
	if err != nil {
		return nil, "", err
	}
	defer resp.Body.Close()
	body, _ := io.ReadAll(resp.Body)
	if resp.StatusCode/100 != 2 {
		return nil, "", parseAPIError(resp.StatusCode, body)
	}
	return body, resp.Header.Get("Content-Type"), nil
}

func (c *Client) do(ctx context.Context, method, path string, query map[string]string, body, out any) error {
	req, err := c.newRequest(ctx, method, path, query, body)
	if err != nil {
		return err
	}
	if c.debug {
		fmt.Fprintf(c.debugOut, "%s %s\n", req.Method, req.URL.String())
	}
	resp, err := c.http.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	respBody, _ := io.ReadAll(resp.Body)
	if c.debug {
		fmt.Fprintf(c.debugOut, "<- %d %d bytes\n", resp.StatusCode, len(respBody))
	}
	if resp.StatusCode/100 != 2 {
		return parseAPIError(resp.StatusCode, respBody)
	}
	if out == nil || len(respBody) == 0 {
		return nil
	}
	if err := json.Unmarshal(respBody, out); err != nil {
		return fmt.Errorf("decode response: %w (body=%s)", err, truncate(respBody, 200))
	}
	return nil
}

func (c *Client) newRequest(ctx context.Context, method, path string, query map[string]string, body any) (*http.Request, error) {
	full := c.baseURL + ensureLeadingSlash(path)
	if len(query) > 0 {
		u, err := url.Parse(full)
		if err != nil {
			return nil, err
		}
		q := u.Query()
		for k, v := range query {
			if v != "" {
				q.Set(k, v)
			}
		}
		u.RawQuery = q.Encode()
		full = u.String()
	}
	var rdr io.Reader
	if body != nil {
		buf, err := json.Marshal(body)
		if err != nil {
			return nil, err
		}
		rdr = bytes.NewReader(buf)
	}
	req, err := http.NewRequestWithContext(ctx, method, full, rdr)
	if err != nil {
		return nil, err
	}
	if c.apiKey != "" {
		req.Header.Set("Authorization", "Bearer "+c.apiKey)
	}
	if body != nil {
		req.Header.Set("Content-Type", "application/json")
	}
	req.Header.Set("Accept", "application/json")
	req.Header.Set("User-Agent", "aikido-cli")
	return req, nil
}

func ensureLeadingSlash(p string) string {
	if strings.HasPrefix(p, "/") {
		return p
	}
	return "/" + p
}

func truncate(b []byte, n int) string {
	if len(b) <= n {
		return string(b)
	}
	return string(b[:n]) + "..."
}
