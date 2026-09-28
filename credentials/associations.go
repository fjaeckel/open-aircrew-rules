package credentials

import (
	"errors"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"

	"gopkg.in/yaml.v3"
)

// Association is one rule document an association sets under a statutory delegation
// (associations.yaml). It is cited by id and title only; its text is never stored.
type Association struct {
	ID          string `yaml:"id"`
	Publisher   string `yaml:"publisher"`
	Title       string `yaml:"title"`
	DelegatedBy string `yaml:"delegated_by"`
	Verified    bool   `yaml:"verified"`
	File        string `yaml:"-"`
}

// LoadAssociations reads associations.yaml; a missing file declares none. An id declared
// twice is an error.
func LoadAssociations(root string) (map[string]*Association, error) {
	path := filepath.Join(root, "associations.yaml")
	out := map[string]*Association{}
	b, err := os.ReadFile(path)
	if errors.Is(err, fs.ErrNotExist) {
		return out, nil
	}
	if err != nil {
		return nil, err
	}
	var doc struct {
		Associations []*Association `yaml:"associations"`
	}
	if err := yaml.Unmarshal(b, &doc); err != nil {
		return nil, fmt.Errorf("%s: %w", path, err)
	}
	for _, a := range doc.Associations {
		if _, dup := out[a.ID]; dup {
			return nil, fmt.Errorf("%s: association %s declared twice", path, a.ID)
		}
		a.File = path
		out[a.ID] = a
	}
	return out, nil
}
