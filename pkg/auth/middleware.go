package auth

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"github.com/hiendaovinh/toolkit/pkg/jwtx"
)

type ctxKey string

var (
	ctxKeyAuthClaims  ctxKey = "AUTH_CLAIMS"
	ctxKeyAuthJWT     ctxKey = "AUTH_JWT"
	ctxKeyAuthSubject ctxKey = "AUTH_SUBJECT"
)

var ErrInvalidSession = errors.New("invalid session")

// context public setters are not recommended but used here to reuse the logic among 2 packages of middleware
func WithAuthClaims(ctx context.Context, v jwtx.RegisteredClaims) context.Context {
	subject, err := v.GetSubject()
	if err != nil {
		return context.WithValue(ctx, ctxKeyAuthClaims, v)
	}

	id, err := uuid.Parse(subject)
	if err != nil {
		return context.WithValue(ctx, ctxKeyAuthClaims, v)
	}

	return context.WithValue(context.WithValue(ctx, ctxKeyAuthSubject, id), ctxKeyAuthClaims, v)
}

func WithAuthJWT(ctx context.Context, v string) context.Context {
	return context.WithValue(ctx, ctxKeyAuthJWT, v)
}

func ResolveJWT(ctx context.Context) string {
	jwt, ok := ctx.Value(ctxKeyAuthJWT).(string)
	if !ok {
		return ""
	}

	return jwt
}

func ResolveSubject(ctx context.Context) string {
	claims, ok := ctx.Value(ctxKeyAuthClaims).(jwtx.RegisteredClaims)
	if !ok {
		return ""
	}

	sub, err := claims.GetSubject()
	if err != nil {
		return ""
	}

	return sub
}

func ResolveSubjectUUID(ctx context.Context) (uuid.UUID, bool) {
	id, ok := ctx.Value(ctxKeyAuthSubject).(uuid.UUID)
	if !ok {
		return uuid.UUID{}, false
	}

	return id, true
}

func ResolveValidSubject(ctx context.Context) (string, error) {
	sub := ResolveSubject(ctx)
	if sub == "" {
		return "", ErrInvalidSession
	}

	return sub, nil
}

func ResolveValidSubjectUUID(ctx context.Context) (uuid.UUID, error) {
	id, ok := ResolveSubjectUUID(ctx)
	if !ok {
		return id, ErrInvalidSession
	}

	return id, nil
}
