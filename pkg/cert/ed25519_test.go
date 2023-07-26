package cert_test

import (
	"crypto/ed25519"
	"crypto/rand"
	"encoding/base64"
	"encoding/json"
	"testing"

	"github.com/hiendaovinh/toolkit/pkg/cert"
	hydra "github.com/ory/hydra-client-go"
	"github.com/stretchr/testify/assert"
)

func TestMarshalED25519Pair(t *testing.T) {
	password := []byte(`123456`)

	pub, priv, err := ed25519.GenerateKey(rand.Reader)
	assert.Nil(t, err, "should be successful")

	_, pemPRIV, err := cert.MarshalED25519PKCS8(priv, password)
	assert.Nil(t, err, "should be successful")

	pubR, privR, err := cert.UnmarshalED25519PKCS8(pemPRIV, password)
	assert.Nil(t, err, "should be successful")

	assert.Equal(t, pubR, pub)
	assert.Equal(t, privR, priv)
}

func TestED25519FromJWK(t *testing.T) {
	// https://www.rfc-editor.org/rfc/rfc8037

	var jwkNullable hydra.NullableJSONWebKey
	err := json.Unmarshal([]byte(`{"kty":"OKP","crv":"Ed25519","d":"nWGxne_9WmC6hEr0kuwsxERJxWl7MmkZcDusAxyuf2A","x":"11qYAYKxCrfVS_7TyWQHOg7hcvPapiMlrwIaaPcHURo"}`), &jwkNullable)
	assert.Nil(t, err, "should be successful")

	jwk := jwkNullable.Get()
	assert.NotNil(t, jwk)

	pub, priv, err := cert.ED25519FromJWK(jwk)
	assert.Nil(t, err, "should be successful")
	seed, _ := base64.RawURLEncoding.DecodeString(*jwkNullable.Get().D)
	assert.Equal(t, seed, priv.Seed())

	pubString := base64.RawURLEncoding.EncodeToString(pub)
	assert.Equal(t, "11qYAYKxCrfVS_7TyWQHOg7hcvPapiMlrwIaaPcHURo", pubString)
}
