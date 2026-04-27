package auth

import (
	"testing"
	"time"
)

const fakeJWT = "eyJ0eXAiOiJKV1QiLCJhbGciOiJIUzI1NiJ9." +
	"eyJpc3MiOiJhaWtpZG8uZGV2IiwiYXVkIjoiaWRlLmFpa2lkbyIsImlhdCI6MTc3NTYzODg5MCwibmJmIjoxNzc1NjM4ODgwLCJleHAiOjI1NjQ1NTcyOTAsImlzX2lkZV90b2tlbiI6dHJ1ZSwidXNlcl9pZCI6MTQ5NDY5LCJ0b2tlbl9pZCI6MTk2ODksInJlZ2lvbiI6ImV1In0." +
	"signature"

func TestDecodeClaims_ParsesRegionAndUser(t *testing.T) {
	c, err := DecodeClaims(fakeJWT)
	if err != nil {
		t.Fatal(err)
	}
	if c.Region != "eu" {
		t.Errorf("region: %q", c.Region)
	}
	if c.UserID != 149469 {
		t.Errorf("user_id: %d", c.UserID)
	}
	if c.TokenID != 19689 {
		t.Errorf("token_id: %d", c.TokenID)
	}
	if c.Expiry().IsZero() {
		t.Errorf("expiry zero")
	}
	if !c.Expiry().After(time.Now()) {
		t.Errorf("expected exp in future")
	}
}

func TestDecodeClaims_RejectsMalformed(t *testing.T) {
	if _, err := DecodeClaims("not-a-jwt"); err == nil {
		t.Fatal("expected error")
	}
	if _, err := DecodeClaims(""); err == nil {
		t.Fatal("expected error")
	}
}
