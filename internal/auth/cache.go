package auth

import (
	"encoding/json"
	"errors"
	"os"
	"path/filepath"
	"time"
)

func tokenCachePath() (string, error) {
	dir, err := os.UserCacheDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(dir, "aikido-cli", "token.json"), nil
}

func LoadCachedToken() (*AccessToken, error) {
	p, err := tokenCachePath()
	if err != nil {
		return nil, err
	}
	buf, err := os.ReadFile(p)
	if err != nil {
		return nil, err
	}
	var t AccessToken
	if err := json.Unmarshal(buf, &t); err != nil {
		return nil, err
	}
	if time.Now().After(t.ExpiresAt) {
		return nil, errors.New("cached token expired")
	}
	return &t, nil
}

func SaveCachedToken(t *AccessToken) error {
	p, err := tokenCachePath()
	if err != nil {
		return err
	}
	if err := os.MkdirAll(filepath.Dir(p), 0o700); err != nil {
		return err
	}
	buf, err := json.Marshal(t)
	if err != nil {
		return err
	}
	return os.WriteFile(p, buf, 0o600)
}

func ClearCachedToken() error {
	p, err := tokenCachePath()
	if err != nil {
		return err
	}
	if err := os.Remove(p); err != nil && !os.IsNotExist(err) {
		return err
	}
	return nil
}
