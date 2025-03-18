package jwtx_test

import (
	"context"
	"crypto/ed25519"
	"encoding/json"
	"math/rand"
	"testing"
	"time"

	"github.com/MicahParks/jwkset"
	"github.com/MicahParks/keyfunc/v3"
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

	seed := int64(42)
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

	store := jwkset.NewMemoryStorage()
	err = store.KeyWrite(context.Background(), jwk)
	assert.NoError(t, err)

	_, tokenStr, err := a.IssueToken(context.Background(), subject, audience, scopes)
	assert.Equal(t, err, nil, "should be successful")

	jwks, err := keyfunc.New(keyfunc.Options{
		Storage: store,
	})
	assert.NoError(t, err)

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
