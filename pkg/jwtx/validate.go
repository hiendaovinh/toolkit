package jwtx

import (
	"context"
	"fmt"

	"github.com/golang-jwt/jwt/v4"
)

func ValidateToken(ctx context.Context, tokenStr string, fnc jwt.Keyfunc) (*jwt.Token, *JWTClaims, error) {
	token, err := jwt.ParseWithClaims(tokenStr, &JWTClaims{}, fnc)
	if err != nil {
		return nil, nil, fmt.Errorf("%w: %q", ErrUnableToParse, err)
	}

	if !token.Valid {
		return nil, nil, ErrInvalidClaims
	}

	claims, ok := token.Claims.(*JWTClaims)
	if !ok {
		return nil, nil, ErrInvalidClaims
	}

	return token, claims, nil
}
