package auth

import (
	"errors"
	"testing"
)

type fakeStore struct {
	data map[string]string
	err  error
}

func newFakeStore() *fakeStore { return &fakeStore{data: map[string]string{}} }

func (f *fakeStore) Set(service, account, secret string) error {
	if f.err != nil {
		return f.err
	}
	f.data[service+"/"+account] = secret
	return nil
}
func (f *fakeStore) Get(service, account string) (string, error) {
	if f.err != nil {
		return "", f.err
	}
	v, ok := f.data[service+"/"+account]
	if !ok {
		return "", ErrNoCredential
	}
	return v, nil
}
func (f *fakeStore) Delete(service, account string) error {
	if f.err != nil {
		return f.err
	}
	delete(f.data, service+"/"+account)
	return nil
}

func TestCredentialStore_RoundTrip(t *testing.T) {
	c := &CredentialStore{store: newFakeStore()}
	creds := ClientCredentials{ClientID: "cid", ClientSecret: "csecret"}
	if err := c.SaveCredentials(creds); err != nil {
		t.Fatal(err)
	}
	got, err := c.LoadCredentials()
	if err != nil {
		t.Fatal(err)
	}
	if got.ClientID != "cid" || got.ClientSecret != "csecret" {
		t.Fatalf("got %+v", got)
	}
	if err := c.Delete(); err != nil {
		t.Fatal(err)
	}
	if _, err := c.LoadCredentials(); !errors.Is(err, ErrNoCredential) {
		t.Fatalf("expected ErrNoCredential, got %v", err)
	}
}

func TestCredentialStore_RejectsEmpty(t *testing.T) {
	c := &CredentialStore{store: newFakeStore()}
	if err := c.SaveCredentials(ClientCredentials{}); err == nil {
		t.Fatal("expected error for empty credentials")
	}
}

func TestCredentialStore_RejectsMalformed(t *testing.T) {
	store := newFakeStore()
	store.data[keyringService+"/"+keyringAccount] = "not-json"
	c := &CredentialStore{store: store}
	if _, err := c.LoadCredentials(); err == nil {
		t.Fatal("expected error for malformed stored value")
	}
}
