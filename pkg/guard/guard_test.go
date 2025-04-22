package guard_test

import (
	"context"
	"errors"
	"sync"
	"testing"

	"github.com/MicahParks/jwkset"
	"github.com/MicahParks/keyfunc/v3"
	"github.com/golang-jwt/jwt/v5"
	"github.com/hiendaovinh/toolkit/pkg/guard"
	"github.com/ory/ladon"
	"github.com/stretchr/testify/assert"
)

type authzFake struct {
	registry *sync.Map
}

func (authz *authzFake) IsAllowed(ctx context.Context, r *ladon.Request) error {
	x, ok := authz.registry.Load(r.Subject)
	if !ok {
		return errors.New("not allowed")
	}

	m, ok := x.(map[string][]string)
	if !ok {
		return errors.New("not allowed")
	}

	for k, v := range m {
		if k != r.Resource {
			continue
		}

		for _, a := range v {
			if a == r.Action {
				return nil
			}
		}

		return errors.New("not allowed")
	}

	return errors.New("not allowed")
}

type test struct {
	authn jwt.Keyfunc
	authz *authzFake
}

func beforeEach(t *testing.T) *test {
	var registry sync.Map
	registry.Store("foo", map[string][]string{
		"bar": {"qux"},
	})

	given := jwkset.NewMemoryStorage()
	x, err := keyfunc.New(keyfunc.Options{Storage: given})
	assert.NoError(t, err)

	return &test{
		x.Keyfunc,
		&authzFake{&registry},
	}
}

func TestGuardAuthz(t *testing.T) {
	test := beforeEach(t)
	guard, err := guard.NewGuard(test.authn, test.authz, nil)
	assert.NoError(t, err)

	sub := "foo"
	resource := "bar"
	ctx := map[string]any{}

	err = guard.Allow(sub, resource, "qux", ctx)
	assert.NoError(t, err)

	err = guard.Allow(sub, resource, "quxx", ctx)
	assert.Error(t, err)
}
