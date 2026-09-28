package credentials

import (
	"testing"

	"gopkg.in/yaml.v3"
)

func load(t *testing.T) *Catalogue {
	t.Helper()
	cat, err := Load("..")
	if err != nil {
		t.Fatal(err)
	}
	for _, e := range cat.Errors {
		t.Error(e)
	}
	return cat
}

func pols(t *testing.T) map[string]*Policy {
	t.Helper()
	p, err := LoadPolicies("..", false)
	if err != nil {
		t.Fatal(err)
	}
	return p
}

func yamlDecode(s string, v any) error { return yaml.Unmarshal([]byte(s), v) }
