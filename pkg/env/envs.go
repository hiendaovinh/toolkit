package env

import (
	"fmt"
	"os"
)

func EnvsRequired(vs ...string) (map[string]string, error) {
	m := map[string]string{}

	for _, v := range vs {
		str := os.Getenv(v)
		if str == "" {
			return nil, fmt.Errorf("missing env: %s", v)
		}

		m[v] = str
	}

	return m, nil
}
