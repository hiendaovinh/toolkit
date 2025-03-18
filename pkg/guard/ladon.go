package guard

import (
	"context"

	"github.com/ory/ladon"
	manager "github.com/ory/ladon/manager/memory"
)

func NewLadon(policies ladon.Policies) (*ladon.Ladon, error) {
	warden := &ladon.Ladon{
		Manager: manager.NewMemoryManager(),
	}

	for _, pol := range policies {
		err := warden.Manager.Create(context.TODO(), pol)
		if err != nil {
			return nil, err
		}
	}

	return warden, nil
}
