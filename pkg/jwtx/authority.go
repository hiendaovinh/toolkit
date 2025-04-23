package jwtx

import (
	"context"
	"crypto/ed25519"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"errors"
	"log"
	"time"

	"github.com/MicahParks/jwkset"
	"github.com/MicahParks/keyfunc/v3"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

type JWK struct {
	Alg   string `json:"alg"`
	Type  string `json:"kty"`
	Curve string `json:"crv"`
	ID    string `json:"kid"`
	X     string `json:"x"`
}

type JWTClaims struct {
	jwt.RegisteredClaims
	Metadata map[string]any `json:"metadata,omitempty"`
}

var (
	ErrUnableToParse = errors.New("unable to parse")
	ErrInvalidClaims = errors.New("invalid token claims")
)

type Authority struct {
	issuer     string
	expiration time.Duration

	pub  ed25519.PublicKey
	priv ed25519.PrivateKey
}

func (g *Authority) kid() string {
	hash := sha256.Sum256(g.pub)
	return base64.RawURLEncoding.EncodeToString(hash[:])
}

func (g *Authority) IssueToken(ctx context.Context, subject string, audience []string, metadata map[string]any) (*JWTClaims, string, error) {
	id, err := uuid.NewV7()
	if err != nil {
		return nil, "", err
	}

	claims := &JWTClaims{
		Metadata: metadata,
		RegisteredClaims: jwt.RegisteredClaims{
			ID:        id.String(),
			Issuer:    g.issuer,
			Subject:   subject,
			Audience:  audience,
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(g.expiration)),
			NotBefore: jwt.NewNumericDate(time.Now()),
		},
	}

	x, err := json.Marshal(&claims)
	log.Println(string(x), err)

	token := jwt.NewWithClaims(jwt.SigningMethodEdDSA, claims)
	token.Header["kid"] = g.kid()

	tokenSigned, err := token.SignedString(g.priv)
	return claims, tokenSigned, err
}

func (g *Authority) PublicJWK(ctx context.Context) (jwkset.JWK, error) {
	id := g.kid()

	metadata := jwkset.JWKMetadataOptions{
		ALG: jwkset.AlgEdDSA,
		KID: id,
		USE: jwkset.UseSig,
	}
	options := jwkset.JWKOptions{
		Metadata: metadata,
	}

	return jwkset.NewJWKFromKey(g.pub, options)
}

func (g *Authority) GenerateKeySet(ctx context.Context) (*jwkset.MemoryJWKSet, keyfunc.Keyfunc, error) {
	set := jwkset.NewMemoryStorage()

	jwk, err := g.PublicJWK(ctx)
	if err != nil {
		return nil, nil, err
	}

	err = set.KeyWrite(ctx, jwk)
	if err != nil {
		return nil, nil, err
	}

	jwks, err := keyfunc.New(keyfunc.Options{
		Storage: set,
	})
	if err != nil {
		return nil, nil, err
	}

	return set, jwks, nil
}

func (g *Authority) PublicJWKS(ctx context.Context) ([]byte, error) {
	set, _, err := g.GenerateKeySet(ctx)
	if err != nil {
		return nil, err
	}

	return set.JSONPublic(ctx)
}

func NewAuthority(issuer string, expiration time.Duration, pub ed25519.PublicKey, priv ed25519.PrivateKey) (*Authority, error) {
	// TODO: do some defensive checking here
	return &Authority{issuer, expiration, pub, priv}, nil
}
