package jwtx_test

import (
	"context"
	"crypto/ed25519"
	"encoding/json"
	"math/rand"
	"testing"
	"time"

	"github.com/hiendaovinh/toolkit/pkg/jwtx"
	"github.com/stretchr/testify/assert"
)

func TestIssueToken(t *testing.T) {
	issuer := "issuer"
	expiration := time.Second * 10
	subject := "subject"
	audience := []string{"audience"}
	meta := map[string]any{
		"foo": 1,
		"bar": "baz",
		"qux": map[int]bool{
			1: true,
			2: false,
		},
	}

	seed := int64(42) // Deterministic seed for testing.
	randSource := rand.New(rand.NewSource(seed))
	randReader := rand.New(randSource)

	pub, priv, err := ed25519.GenerateKey(randReader)
	assert.Equal(t, err, nil, "should be successful")

	a, err := jwtx.NewAuthority(issuer, expiration, pub, priv)
	assert.Equal(t, err, nil, "should be successful")

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

	_, tokenStr, err := a.IssueToken(context.Background(), subject, audience, meta)
	assert.Equal(t, err, nil, "should be successful")

	_, claims, err := jwtx.ValidateToken(context.Background(), tokenStr, jwks.Keyfunc)
	assert.Equal(t, err, nil, "should be successful")

	assert.Equal(t, claims.Issuer, issuer, "mismatched issuers")
	assert.Equal(t, claims.Subject, subject, "mismatched subject")
	assert.Equal(t, []string(claims.Audience), audience, "mismatched audience")
	assert.NotEqual(t, meta, claims.Metadata, "mismatched metadata") // maps are always encoded as map[string]any, no matter what is the original values
	assert.Equal(t, map[string]any{
		"foo": float64(1),
		"bar": "baz",
		"qux": map[string]any{
			"1": true,
			"2": false,
		},
	}, claims.Metadata)

	assert.LessOrEqual(t, time.Until(claims.ExpiresAt.Time), expiration, "invalid expiration")
}
