package guard

import (
	"context"
	"errors"

	"github.com/golang-jwt/jwt/v5"
	"github.com/hiendaovinh/toolkit/pkg/jwtx"
	"github.com/ory/ladon"
)

type Guard struct {
	authn         jwt.Keyfunc
	authz         AuthChecker
	claimsFactory func() jwtx.RegisteredClaims
}

type AuthChecker interface {
	IsAllowed(ctx context.Context, r *ladon.Request) error
}

func NewGuard(authn jwt.Keyfunc, authz AuthChecker, claimsFactory func() jwtx.RegisteredClaims) (*Guard, error) {
	if authn == nil || authz == nil {
		return nil, errors.New("invalid authn or authz")
	}

	return &Guard{authn, authz, claimsFactory}, nil
}

func (guard *Guard) Allow(sub string, resource string, action string, ctx map[string]any) error {
	r := &ladon.Request{
		Subject:  sub,
		Resource: resource,
		Action:   action,
		Context:  ctx,
	}

	return guard.authz.IsAllowed(context.TODO(), r)
}

func (guard *Guard) AuthenticateJWT(ctx context.Context, tokenStr string, claims jwtx.RegisteredClaims) (*jwt.Token, error) {
	return jwtx.ValidateToken(ctx, tokenStr, guard.authn, claims)
}

func (guard *Guard) NewClaims() jwtx.RegisteredClaims {
	return guard.claimsFactory()
}
