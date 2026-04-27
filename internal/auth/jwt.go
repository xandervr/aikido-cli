package auth

import (
	"encoding/base64"
	"encoding/json"
	"errors"
	"strings"
	"time"
)

type Claims struct {
	Issuer   string `json:"iss"`
	Audience string `json:"aud"`
	Region   string `json:"region"`
	UserID   int    `json:"user_id"`
	TokenID  int    `json:"token_id"`
	Exp      int64  `json:"exp"`
	IssuedAt int64  `json:"iat"`
}

func (c Claims) Expiry() time.Time {
	if c.Exp == 0 {
		return time.Time{}
	}
	return time.Unix(c.Exp, 0)
}

func DecodeClaims(token string) (Claims, error) {
	parts := strings.Split(token, ".")
	if len(parts) < 2 {
		return Claims{}, errors.New("not a jwt")
	}
	raw, err := base64.RawURLEncoding.DecodeString(parts[1])
	if err != nil {
		raw, err = base64.StdEncoding.DecodeString(parts[1])
		if err != nil {
			return Claims{}, err
		}
	}
	var c Claims
	if err := json.Unmarshal(raw, &c); err != nil {
		return Claims{}, err
	}
	return c, nil
}
