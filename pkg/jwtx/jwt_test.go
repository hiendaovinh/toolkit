package jwtx_test

import (
	"context"
	"crypto/ed25519"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"math/rand"
	"strings"
	"testing"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/hiendaovinh/toolkit/pkg/jwtx"
	"github.com/stretchr/testify/assert"
)

func TestIssueToken(t *testing.T) {
	issuer := "issuer"
	expiration := time.Second * 10
	subject := "subject"
	audience := []string{"audience"}

	seed := int64(42) // Deterministic seed for testing.
	randSource := rand.New(rand.NewSource(seed))
	randReader := rand.New(randSource)

	pub, priv, err := ed25519.GenerateKey(randReader)
	assert.NoError(t, err)

	a, err := jwtx.NewAuthority(issuer, expiration, pub, priv)
	assert.NoError(t, err)

	jwk, err := a.PublicJWK(context.Background())
	assert.NoError(t, err)

	x, err := json.Marshal(jwk.Marshal())
	assert.NoError(t, err)

	var m1, m2 map[string]interface{}
	err = json.Unmarshal(x, &m1)
	assert.NoError(t, err)

	err = json.Unmarshal([]byte(`{"alg":"EdDSA","kty":"OKP","crv":"Ed25519","kid":"03AveCUC5gHj1iNiADoa5PnOWfP2NH-oWSoQIMFySjM","x":"z-xJSnzIQtg0Dywzvl7VUzTLpj5VFajA8wZ3IRz6ztQ", "use":"sig"}`), &m2)
	assert.NoError(t, err)
	assert.Equal(t, m1, m2)

	_, jwks, err := a.GenerateKeySet(context.Background())
	assert.NoError(t, err)

	claimsIn := jwtx.DefaultJWTClaims{}
	claimsOut := jwtx.DefaultJWTClaims{}

	tokenStr, err := a.IssueToken(context.Background(), subject, audience, &claimsIn)
	assert.NoError(t, err)

	token, err := jwtx.ValidateToken(context.Background(), tokenStr, jwks.Keyfunc, &claimsOut)
	assert.NoError(t, err)
	assert.Equal(t, claimsIn, claimsOut)
	assert.Equal(t, &claimsOut, token.Claims)
	assert.Equal(t, issuer, claimsOut.Issuer, "mismatched issuers")
	assert.Equal(t, subject, claimsOut.Subject, "mismatched subject")
	assert.Equal(t, audience, []string(claimsOut.Audience), "mismatched audience")
	assert.LessOrEqual(t, time.Until(claimsOut.ExpiresAt.Time), expiration, "invalid expiration")

	parts := strings.Split(tokenStr, ".")
	data, err := base64.RawURLEncoding.DecodeString(parts[1])
	assert.NoError(t, err)

	var m3 map[string]any
	err = json.Unmarshal(data, &m3)
	assert.NoError(t, err)

	m3["subject"] = "malicious"
	data, err = json.Marshal(&m3)
	assert.NoError(t, err)
	forgedBody := base64.RawURLEncoding.EncodeToString(data)

	_, err = jwtx.ValidateToken(context.Background(), fmt.Sprintf("%s.%s.%s", parts[0], forgedBody, parts[2]), jwks.Keyfunc, &claimsOut)
	assert.ErrorIs(t, err, jwtx.ErrInvalid)
}

type CustomJWTClaims struct {
	jwt.RegisteredClaims
	Name   string   `json:"string"`
	Scopes []string `json:"scopes"`
}

func (c *CustomJWTClaims) Register(cs jwt.RegisteredClaims) jwtx.RegisteredClaims {
	c.RegisteredClaims = cs
	return c
}

func TestCustomClaims(t *testing.T) {
	issuer := "issuer"
	expiration := time.Second * 10
	subject := "subject"
	audience := []string{"audience"}

	seed := int64(42) // Deterministic seed for testing.
	randSource := rand.New(rand.NewSource(seed))
	randReader := rand.New(randSource)

	pub, priv, err := ed25519.GenerateKey(randReader)
	assert.NoError(t, err)

	a, err := jwtx.NewAuthority(issuer, expiration, pub, priv)
	assert.NoError(t, err)

	jwk, err := a.PublicJWK(context.Background())
	assert.NoError(t, err)

	x, err := json.Marshal(jwk.Marshal())
	assert.NoError(t, err)

	var m1, m2 map[string]interface{}
	err = json.Unmarshal(x, &m1)
	assert.NoError(t, err)

	err = json.Unmarshal([]byte(`{"alg":"EdDSA","kty":"OKP","crv":"Ed25519","kid":"03AveCUC5gHj1iNiADoa5PnOWfP2NH-oWSoQIMFySjM","x":"z-xJSnzIQtg0Dywzvl7VUzTLpj5VFajA8wZ3IRz6ztQ", "use":"sig"}`), &m2)
	assert.NoError(t, err)
	assert.Equal(t, m1, m2)

	_, jwks, err := a.GenerateKeySet(context.Background())
	assert.NoError(t, err)

	claimsIn := CustomJWTClaims{
		Name:   "foo",
		Scopes: []string{"bar", "qux"},
	}
	claimsOut := CustomJWTClaims{
		Name:   "foo",
		Scopes: []string{"bar", "qux"},
	}

	tokenStr, err := a.IssueToken(context.Background(), subject, audience, &claimsIn)
	assert.NoError(t, err)

	_, err = jwtx.ValidateToken(context.Background(), tokenStr, jwks.Keyfunc, &claimsOut)
	assert.NoError(t, err)
	assert.Equal(t, claimsIn, claimsOut)
	assert.Equal(t, issuer, claimsOut.Issuer, "mismatched issuers")
	assert.Equal(t, subject, claimsOut.Subject, "mismatched subject")
	assert.Equal(t, audience, []string(claimsOut.Audience), "mismatched audience")
	assert.LessOrEqual(t, time.Until(claimsOut.ExpiresAt.Time), expiration, "invalid expiration")
}
