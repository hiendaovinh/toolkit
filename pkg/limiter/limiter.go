package limiter

import (
	"context"
	"errors"
	"fmt"

	"github.com/go-redis/redis_rate/v10"
	"github.com/hiendaovinh/toolkit/pkg/auth"
	"github.com/redis/go-redis/v9"
)

type ctxKey string

const (
	ctxKeySkip ctxKey = "skip"
)

func Skip(ctx context.Context) context.Context {
	return context.WithValue(ctx, ctxKeySkip, true)
}

var ErrRateLimited = errors.New("rate limited")

type Limiter struct {
	limiter *redis_rate.Limiter
}

func NewLimiter(rdb *redis.Client) (*Limiter, error) {
	limiter := redis_rate.NewLimiter(rdb)
	return &Limiter{limiter}, nil
}

func (l *Limiter) Allow(ctx context.Context, key string, limit redis_rate.Limit) error {
	skip := ctx.Value(ctxKeySkip)
	if v, ok := skip.(bool); ok && v {
		return nil
	}

	sub := auth.ResolveSubject(ctx)
	res, err := l.limiter.Allow(ctx, fmt.Sprintf("%s:%s", sub, key), limit)
	if err != nil {
		return err
	}

	if res.Allowed <= 0 {
		return ErrRateLimited
	}

	return nil
}
