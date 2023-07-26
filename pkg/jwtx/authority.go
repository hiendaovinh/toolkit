package jwtx

import (
	"context"
	"crypto/ed25519"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v4"
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
	Scopes []string `json:"scopes,omitempty"`
	jwt.RegisteredClaims
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

func (g *Authority) IssueToken(ctx context.Context, subject string, audience []string, scopes []string) (*JWTClaims, string, error) {
	claims := &JWTClaims{
		Scopes: scopes,
		RegisteredClaims: jwt.RegisteredClaims{
			ID:        uuid.New().String(),
			Issuer:    g.issuer,
			Subject:   subject,
			Audience:  audience,
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(g.expiration)),
			NotBefore: jwt.NewNumericDate(time.Now()),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodEdDSA, claims)
	token.Header["kid"] = g.kid()

	tokenSigned, err := token.SignedString(g.priv)
	return claims, tokenSigned, err
}

func (g *Authority) PublicJWK(ctx context.Context) (*JWK, error) {
	x := base64.RawURLEncoding.EncodeToString(g.pub)
	id := g.kid()

	return &JWK{
		Alg:   "EdDSA",
		Type:  "OKP",
		Curve: "Ed25519",
		ID:    id,
		X:     x,
	}, nil
}

func (g *Authority) PublicJWKS(ctx context.Context) ([]byte, error) {
	jwk, err := g.PublicJWK(ctx)
	if err != nil {
		return nil, err
	}

	output, err := json.Marshal(map[string][]any{
		"keys": {jwk},
	})
	return output, err
}

func NewAuthority(issuer string, expiration time.Duration, pub ed25519.PublicKey, priv ed25519.PrivateKey) (*Authority, error) {
	// TODO: do some defensive checking here
	return &Authority{issuer, expiration, pub, priv}, nil
}
