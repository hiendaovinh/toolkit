package auth_test

import (
	"errors"
	"sync"
	"testing"

	"github.com/MicahParks/keyfunc"
	"github.com/hiendaovinh/toolkit/pkg/auth"
	"github.com/ory/ladon"
	"github.com/stretchr/testify/assert"
)

type authzFake struct {
	registry *sync.Map
}

func (authz *authzFake) IsAllowed(r *ladon.Request) error {
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
	authn *keyfunc.JWKS
	authz *authzFake
}

func beforeEach(t *testing.T) *test {
	var registry sync.Map
	registry.Store("foo", map[string][]string{
		"bar": {"qux"},
	})

	return &test{
		&keyfunc.JWKS{},
		&authzFake{&registry},
	}
}

func TestGuardAuthz(t *testing.T) {
	test := beforeEach(t)
	guard, err := auth.NewGuard(test.authn, test.authz)
	assert.NoError(t, err)

	sub := "foo"
	resource := "bar"
	ctx := map[string]any{}

	err = guard.Allow(sub, resource, "qux", ctx)
	assert.NoError(t, err)

	err = guard.Allow(sub, resource, "quxx", ctx)
	assert.Error(t, err)
}
