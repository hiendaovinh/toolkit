package guard

import (
	"context"
	"errors"

	"github.com/golang-jwt/jwt/v5"
	"github.com/hiendaovinh/toolkit/pkg/jwtx"
)

type GuardLite struct {
	authn         jwt.Keyfunc
	claimsFactory func() jwtx.RegisteredClaims
}

func (guard *GuardLite) AuthenticateJWT(ctx context.Context, tokenStr string, claims jwtx.RegisteredClaims) (*jwt.Token, error) {
	return jwtx.ValidateToken(ctx, tokenStr, guard.authn, claims)
}

func (guard *GuardLite) NewClaims() jwtx.RegisteredClaims {
	return guard.claimsFactory()
}

func NewGuardLite(authn jwt.Keyfunc, claimsFactory func() jwtx.RegisteredClaims) (*GuardLite, error) {
	if authn == nil {
		return nil, errors.New("missing authentication")
	}

	return &GuardLite{authn: authn, claimsFactory: claimsFactory}, nil
}
