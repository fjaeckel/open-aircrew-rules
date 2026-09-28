package credentials

import (
	"fmt"
	"os"
	"path/filepath"

	"gopkg.in/yaml.v3"
)

// Policy is one convention with no legal text behind it (policies.yaml).
type Policy struct {
	ID        string `yaml:"id"`
	Statement string `yaml:"statement"`
	Rationale string `yaml:"rationale"`
	File      string `yaml:"-"`
}

// LoadPolicies reads policies.yaml and, with fragments, fragments/policies/*.yaml. A policy id
// declared twice is an error.
func LoadPolicies(root string, fragments bool) (map[string]*Policy, error) {
	files := []string{filepath.Join(root, "policies.yaml")}
	if fragments {
		more, err := fragmentFiles(root, "policies")
		if err != nil {
			return nil, err
		}
		files = append(files, more...)
	}
	out := map[string]*Policy{}
	for _, f := range files {
		b, err := os.ReadFile(f)
		if err != nil {
			return nil, err
		}
		var doc struct {
			Policies []*Policy `yaml:"policies"`
		}
		if err := yaml.Unmarshal(b, &doc); err != nil {
			return nil, fmt.Errorf("%s: %w", f, err)
		}
		for _, p := range doc.Policies {
			if _, dup := out[p.ID]; dup {
				return nil, fmt.Errorf("%s: policy %s declared twice", f, p.ID)
			}
			p.File = f
			out[p.ID] = p
		}
	}
	return out, nil
}
