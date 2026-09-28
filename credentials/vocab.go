// Package credentials loads the credential-centric catalogue (credentials/**), compiles
// each credential evaluation onto the engine's rule model, resolves source references and
// runs worked examples (DESIGN.md sections 3 and 4).
package credentials

import (
	"fmt"
	"os"

	"gopkg.in/yaml.v3"
)

// Vocab is the credential_vocabulary section of vocabulary.yaml.
type Vocab struct {
	Kinds          map[string]string          `yaml:"kinds"`
	Selects        map[string]string          `yaml:"selects"`
	About          map[string]AboutDef        `yaml:"about"`
	RefAuthorities map[string]RefAuthorityDef `yaml:"ref_authorities"`
	Counts         map[string]CountDef        `yaml:"counts"`
	Qualifiers     map[string]QualifierDef    `yaml:"qualifiers"`
	OutcomePresets map[string]PresetDef       `yaml:"outcome_presets"`
	// LimitationScopes are the scopes a requirement or evaluation may limit (limits:).
	LimitationScopes map[string]string `yaml:"limitation_scopes"`
	// AuthorityConventions are conventions every evaluation of an authority's credentials
	// follows, keyed by credential authority.
	AuthorityConventions map[string]AuthorityConvention `yaml:"authority_conventions"`
}

// AuthorityConvention is how one authority counts periods.
type AuthorityConvention struct {
	// ValidityEnds is same_day (default) or day_before: a derived validity period ends the
	// day before the day that corresponds to its start (§§ 187(2), 188(2) BGB).
	ValidityEnds string  `yaml:"validity_ends"`
	Ref          RefList `yaml:"ref"`
}

// PresetDef is one outcome preset: its description and, per stage id, the policy ref behind
// a status the preset sets by convention.
type PresetDef struct {
	Description string            `yaml:"description"`
	Conventions map[string]string `yaml:"conventions"`
}

// AboutDef is what an evaluation can be about.
type AboutDef struct {
	Subject     string `yaml:"subject"`
	Selects     string `yaml:"selects"`
	Description string `yaml:"description"`
}

// RefAuthorityDef maps a reference prefix to a sources/ directory.
type RefAuthorityDef struct {
	Directory string `yaml:"directory"`
	File      string `yaml:"file"`
	Cite      string `yaml:"cite"`
}

// CountDef is one count word of passes_if.
type CountDef struct {
	Metric string         `yaml:"metric"`
	Events []string       `yaml:"events"`
	Amount string         `yaml:"amount"`
	ID     string         `yaml:"id"`
	Name   string         `yaml:"name"`
	Unit   string         `yaml:"unit"`
	Remedy string         `yaml:"remedy"`
	Filter map[string]any `yaml:"filter"`
}

// QualifierDef is one qualifier word.
type QualifierDef struct {
	Window      string `yaml:"window"`
	Anchor      string `yaml:"anchor"`
	Filter      string `yaml:"filter"`
	Value       string `yaml:"value"`
	Document    bool   `yaml:"document"`
	Description string `yaml:"description"`
}

// LoadVocab reads credential_vocabulary from vocabulary.yaml.
func LoadVocab(path string) (*Vocab, error) {
	b, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	var doc struct {
		V *Vocab `yaml:"credential_vocabulary"`
	}
	if err := yaml.Unmarshal(b, &doc); err != nil {
		return nil, fmt.Errorf("%s: %w", path, err)
	}
	if doc.V == nil {
		return nil, fmt.Errorf("%s: no credential_vocabulary section", path)
	}
	return doc.V, nil
}
