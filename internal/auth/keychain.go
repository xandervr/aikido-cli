package auth

import (
	"encoding/json"
	"errors"

	"github.com/zalando/go-keyring"
)

const (
	keyringService = "aikido-cli"
	keyringAccount = "default"
)

var ErrNoCredential = errors.New("no credential stored")

type ClientCredentials struct {
	ClientID     string `json:"client_id"`
	ClientSecret string `json:"client_secret"`
}

type secretStore interface {
	Set(service, account, secret string) error
	Get(service, account string) (string, error)
	Delete(service, account string) error
}

type keyringStore struct{}

func (keyringStore) Set(s, a, v string) error { return keyring.Set(s, a, v) }
func (keyringStore) Delete(s, a string) error { return keyring.Delete(s, a) }
func (keyringStore) Get(s, a string) (string, error) {
	v, err := keyring.Get(s, a)
	if errors.Is(err, keyring.ErrNotFound) {
		return "", ErrNoCredential
	}
	return v, err
}

type CredentialStore struct {
	store secretStore
}

func NewCredentialStore() *CredentialStore { return &CredentialStore{store: keyringStore{}} }

// SaveCredentials stores the client_id/client_secret pair as a JSON blob in
// the OS keychain.
func (c *CredentialStore) SaveCredentials(creds ClientCredentials) error {
	if creds.ClientID == "" || creds.ClientSecret == "" {
		return errors.New("client_id and client_secret must be non-empty")
	}
	buf, err := json.Marshal(creds)
	if err != nil {
		return err
	}
	return c.store.Set(keyringService, keyringAccount, string(buf))
}

// LoadCredentials retrieves the stored client_id/client_secret pair.
// Returns ErrNoCredential if nothing is stored.
func (c *CredentialStore) LoadCredentials() (ClientCredentials, error) {
	raw, err := c.store.Get(keyringService, keyringAccount)
	if err != nil {
		return ClientCredentials{}, err
	}
	var creds ClientCredentials
	if err := json.Unmarshal([]byte(raw), &creds); err != nil {
		return ClientCredentials{}, errors.New("stored credential is not in client_id/secret JSON format; run 'aikido auth login' to refresh")
	}
	if creds.ClientID == "" || creds.ClientSecret == "" {
		return ClientCredentials{}, errors.New("stored credential is missing client_id or client_secret; run 'aikido auth login' to refresh")
	}
	return creds, nil
}

func (c *CredentialStore) Delete() error {
	err := c.store.Delete(keyringService, keyringAccount)
	if errors.Is(err, ErrNoCredential) {
		return nil
	}
	return err
}
