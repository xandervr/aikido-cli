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
	if err := c.Save("xyz"); err != nil {
		t.Fatal(err)
	}
	got, err := c.Load()
	if err != nil {
		t.Fatal(err)
	}
	if got != "xyz" {
		t.Fatalf("got %q", got)
	}
	if err := c.Delete(); err != nil {
		t.Fatal(err)
	}
	if _, err := c.Load(); !errors.Is(err, ErrNoCredential) {
		t.Fatalf("expected ErrNoCredential, got %v", err)
	}
}
