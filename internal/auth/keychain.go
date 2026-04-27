package auth

import (
	"errors"

	"github.com/zalando/go-keyring"
)

const (
	keyringService = "aikido-cli"
	keyringAccount = "default"
)

var ErrNoCredential = errors.New("no credential stored")

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

func (c *CredentialStore) Save(token string) error {
	return c.store.Set(keyringService, keyringAccount, token)
}

func (c *CredentialStore) Load() (string, error) {
	return c.store.Get(keyringService, keyringAccount)
}

func (c *CredentialStore) Delete() error {
	err := c.store.Delete(keyringService, keyringAccount)
	if errors.Is(err, ErrNoCredential) {
		return nil
	}
	return err
}
