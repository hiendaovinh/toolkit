package jwtx_test

import (
	"context"
	"crypto/ed25519"
	"testing"
	"time"

	"github.com/MicahParks/keyfunc"
	"github.com/golang-jwt/jwt/v4"
	"github.com/hiendaovinh/toolkit/pkg/jwtx"
	"github.com/stretchr/testify/assert"
	"golang.org/x/exp/slices"
)

func TestIssueToken(t *testing.T) {
	issuer := "issuer"
	expiration := time.Second * 10
	subject := "subject"
	audience := []string{"audience"}
	scopes := []string{"foo", "bar"}

	pub, priv, err := ed25519.GenerateKey(nil)
	assert.Equal(t, err, nil, "should be successful")

	a, err := jwtx.NewAuthority(issuer, expiration, pub, priv)
	assert.Equal(t, err, nil, "should be successful")

	jwk, err := a.PublicJWK(context.Background())
	assert.Equal(t, err, nil, "should be successful")

	_, tokenStr, err := a.IssueToken(context.Background(), subject, audience, scopes)
	assert.Equal(t, err, nil, "should be successful")

	jwks := keyfunc.NewGiven(map[string]keyfunc.GivenKey{
		jwk.ID: keyfunc.NewGivenEdDSACustomWithOptions(pub, keyfunc.GivenKeyOptions{
			Algorithm: jwt.SigningMethodEdDSA.Alg(),
		}),
	})

	_, claims, err := jwtx.ValidateToken(context.Background(), tokenStr, jwks.Keyfunc)
	assert.Equal(t, err, nil, "should be successful")

	assert.Equal(t, claims.Issuer, issuer, "mismatched issuers")
	assert.Equal(t, claims.Subject, subject, "mismatched subject")
	assert.Equal(t, []string(claims.Audience), audience, "mismatched audience")
	assert.Equal(t, claims.Scopes, scopes, "mismatched scopes")
	assert.True(t, slices.Contains(claims.Scopes, "foo"))
	assert.True(t, slices.Contains(claims.Scopes, "bar"))
	assert.False(t, slices.Contains(claims.Scopes, "qux"))

	assert.LessOrEqual(t, time.Until(claims.ExpiresAt.Time), expiration, "invalid expiration")
}
