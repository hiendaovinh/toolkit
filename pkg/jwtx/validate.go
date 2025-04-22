package jwtx

import (
	"context"
	"fmt"

	"github.com/golang-jwt/jwt/v5"
)

func ValidateToken(ctx context.Context, tokenStr string, fnc jwt.Keyfunc, claims RegisteredClaims) (*jwt.Token, error) {
	token, err := jwt.ParseWithClaims(tokenStr, claims, fnc)
	if err != nil {
		return nil, fmt.Errorf("%w: %q", ErrInvalid, err)
	}

	return token, nil
}
