package engine

import (
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"gopkg.in/yaml.v3"
)

// Keys is messages/keys.yaml.
type Keys struct {
	Keys []KeyDef `yaml:"keys"`
	byID map[string]*KeyDef
}

// KeyDef is one message key.
type KeyDef struct {
	Key    string   `yaml:"key"`
	Kind   string   `yaml:"kind"`
	Params []string `yaml:"params"`
	Notes  string   `yaml:"notes"`
}

// Get returns the key definition.
func (k *Keys) Get(key string) (*KeyDef, bool) {
	d, ok := k.byID[key]
	return d, ok
}

// LoadKeys reads messages/keys.yaml.
func LoadKeys(path string) (*Keys, error) {
	b, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	var k Keys
	if err := yaml.Unmarshal(b, &k); err != nil {
		return nil, fmt.Errorf("%s: %w", path, err)
	}
	k.byID = map[string]*KeyDef{}
	for i := range k.Keys {
		d := &k.Keys[i]
		if _, dup := k.byID[d.Key]; dup {
			return nil, fmt.Errorf("%s: duplicate key %s", path, d.Key)
		}
		k.byID[d.Key] = d
	}
	return &k, nil
}

// Catalogue is a set of rules with the vocabulary and message keys they use. The credential
// compiler (package credentials) builds one from the credential files.
type Catalogue struct {
	Root       string
	Vocabulary *Vocabulary
	Keys       *Keys
	Rules      []*Rule
	// Errors lists rules that could not be added (duplicate ids); those rules are skipped.
	Errors []error
	byID   map[string]*Rule
}

// Rule returns the rule with the id.
func (c *Catalogue) Rule(id string) (*Rule, bool) {
	r, ok := c.byID[id]
	return r, ok
}

// YAMLFiles lists *.yaml files under dir, sorted; a missing dir yields none.
func YAMLFiles(dir string) ([]string, error) {
	var out []string
	err := filepath.WalkDir(dir, func(p string, d fs.DirEntry, err error) error {
		if err != nil {
			if os.IsNotExist(err) {
				return filepath.SkipDir
			}
			return err
		}
		if !d.IsDir() && strings.HasSuffix(p, ".yaml") {
			out = append(out, p)
		}
		return nil
	})
	sort.Strings(out)
	return out, err
}

// NewCatalogue returns a catalogue of the given rules (compiled credential evaluations);
// duplicate ids are reported in Errors and skipped.
func NewCatalogue(root string, v *Vocabulary, k *Keys, rules []*Rule) *Catalogue {
	c := &Catalogue{Root: root, Vocabulary: v, Keys: k, byID: map[string]*Rule{}}
	for _, r := range rules {
		if _, dup := c.byID[r.ID]; dup {
			c.Errors = append(c.Errors, fmt.Errorf("%s: duplicate rule id %s", r.File, r.ID))
			continue
		}
		c.byID[r.ID] = r
		c.Rules = append(c.Rules, r)
	}
	return c
}
