package auth

import (
	"context"
	"errors"

	"github.com/MicahParks/keyfunc"
	"github.com/golang-jwt/jwt/v4"
	"github.com/hiendaovinh/toolkit/pkg/jwtx"
	"github.com/ory/ladon"
)

type Guard struct {
	authn *keyfunc.JWKS
	authz AuthChecker
}

type AuthChecker interface {
	IsAllowed(r *ladon.Request) error
}

func NewGuard(authn *keyfunc.JWKS, authz AuthChecker) (*Guard, error) {
	if authn == nil || authz == nil {
		return nil, errors.New("invalid authn or authz")
	}

	return &Guard{authn, authz}, nil
}

func (guard *Guard) Allow(sub string, resource string, action string, ctx map[string]any) error {
	r := &ladon.Request{
		Subject:  sub,
		Resource: resource,
		Action:   action,
		Context:  ctx,
	}

	return guard.authz.IsAllowed(r)
}

func (guard *Guard) AuthenticateJWT(ctx context.Context, tokenStr string) (*jwt.Token, *jwtx.JWTClaims, error) {
	return jwtx.ValidateToken(ctx, tokenStr, guard.authn.Keyfunc)
}
