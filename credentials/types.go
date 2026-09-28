package credentials

import (
	"gopkg.in/yaml.v3"

	"github.com/fjaeckel/open-aircrew-rules/engine"
)

// Credential is one file under credentials/<authority>/<kind dir>/.
type Credential struct {
	Name            string           `yaml:"credential"`
	ID              string           `yaml:"id"`
	Kind            string           `yaml:"kind"`
	Authority       string           `yaml:"authority"`
	HeldOn          []string         `yaml:"held_on"`
	Selects         Selects          `yaml:"selects"`
	Validity        *ValidityNote    `yaml:"validity"`
	Evaluations     []*Evaluation    `yaml:"evaluations"`
	Interpretations []Interpretation `yaml:"interpretations"`

	File string `yaml:"-"`
}

// Shared is one parameterised evaluation under credentials/<authority>/shared/.
type Shared struct {
	Name            string            `yaml:"shared"`
	ID              string            `yaml:"id"`
	Authority       string            `yaml:"authority"`
	Params          map[string]string `yaml:"params"`
	Evaluation      yaml.Node         `yaml:"evaluation"`
	Interpretations []Interpretation  `yaml:"interpretations"`

	File string `yaml:"-"`
}

// Selects says how a credential appears in the record.
type Selects struct {
	Licence    *Part `yaml:"licence"`
	Ratings    *Part `yaml:"ratings"`
	Privilege  *Part `yaml:"privilege"`
	Credential *Part `yaml:"credential"`
}

// Part selects record items; kinds are licence kinds, privilege kinds or credential types
// depending on the part.
type Part struct {
	Kinds           []string `yaml:"kinds"`
	Classes         []string `yaml:"classes"`
	Authorities     []string `yaml:"authorities"`
	NotAuthorities  []string `yaml:"not_authorities"`
	LicenceKinds    []string `yaml:"licence_kinds"`
	NotLicenceKinds []string `yaml:"not_licence_kinds"`
}

// ValidityNote is the credential's validity period as a reader sees it.
type ValidityNote struct {
	PeriodMonths int    `yaml:"period_months"`
	Ref          string `yaml:"ref"`
	Note         string `yaml:"note"`
}

// Interpretation is a reading of the text that the evaluations depend on.
type Interpretation struct {
	ID         string   `yaml:"id"`
	Reading    string   `yaml:"reading"`
	Ref        RefList  `yaml:"ref"`
	Affects    []string `yaml:"affects"`
	ApprovedBy *string  `yaml:"approved_by"`
	ApprovedOn *string  `yaml:"approved_on"`
}

// RefList is one reference or a list of them.
type RefList []string

// UnmarshalYAML accepts a string or a list of strings.
func (r *RefList) UnmarshalYAML(n *yaml.Node) error {
	if n.Kind == yaml.ScalarNode {
		*r = RefList{n.Value}
		return nil
	}
	var l []string
	if err := n.Decode(&l); err != nil {
		return err
	}
	*r = l
	return nil
}

// Evaluation is one thing a credential needs.
type Evaluation struct {
	ID        string    `yaml:"id"`
	Asks      string    `yaml:"asks"`
	Source    string    `yaml:"source"`
	AlsoCites []string  `yaml:"also_cites"`
	Uses      string    `yaml:"uses"`
	With      yaml.Node `yaml:"with"`
	// RequiresAll and RequiresAny name the credentials this one needs: every one of
	// RequiresAll and at least one of RequiresAny (see Composite).
	RequiresAll    []string  `yaml:"requires_all"`
	RequiresAny    []string  `yaml:"requires_any"`
	About          string    `yaml:"about"`
	OnlyFor        *Scope    `yaml:"only_for"`
	Scope          *Scope    `yaml:"scope"`
	RelevantClass  *PoolDef  `yaml:"relevant_class"`
	Counting       yaml.Node `yaml:"counting"`
	PassesIf       yaml.Node `yaml:"passes_if"`
	RestoredBy     yaml.Node `yaml:"restored_by"`
	ValidFor       *ValidFor `yaml:"valid_for"`
	Outcomes       yaml.Node `yaml:"outcomes"`
	DescriptionKey string    `yaml:"description_key"`
	OnFail         *OnFail   `yaml:"on_fail"`
	// EffectiveFrom and EffectiveTo bound the dates the evaluation applies to (both
	// inclusive); a regulation change ends one evaluation and starts its successor.
	EffectiveFrom *engine.Date `yaml:"effective_from"`
	EffectiveTo   *engine.Date `yaml:"effective_to"`
	Line          int          `yaml:"-"`
	From          string       `yaml:"-"`
	Shared        *Shared      `yaml:"-"`
	Owner         *Credential  `yaml:"-"`
}

// Scope narrows (only_for) or replaces (scope) the subjects an evaluation selects.
type Scope struct {
	Authorities     []string      `yaml:"authorities"`
	NotAuthorities  []string      `yaml:"not_authorities"`
	LicenceKinds    []string      `yaml:"licence_kinds"`
	NotLicenceKinds []string      `yaml:"not_licence_kinds"`
	Classes         []string      `yaml:"classes"`
	NotClasses      []string      `yaml:"not_classes"`
	CredentialTypes []string      `yaml:"credential_types"`
	PrivilegeKinds  []string      `yaml:"privilege_kinds"`
	LaunchMethods   []string      `yaml:"launch_methods"`
	TypeRated       *bool         `yaml:"type_rated"`
	Programme       string        `yaml:"programme"`
	WhenHolding     *engine.Holds `yaml:"when_holding"`
	Ref             RefList       `yaml:"ref"`
}

// PoolDef makes in_class count the classes of a pool the holder rates on the same licence.
type PoolDef struct {
	PooledWithHeld []string `yaml:"pooled_with_held"`
	Ref            RefList  `yaml:"ref"`
}

// ValidFor derives an expiry date.
type ValidFor struct {
	CountedFrom string       `yaml:"counted_from"`
	Periods     []PeriodSpec `yaml:"periods"`
	Ref         RefList      `yaml:"ref"`
}

// PeriodSpec is one validity period; the first whose ages hold applies.
type PeriodSpec struct {
	AgeUnder   *int    `yaml:"age_under"`
	AgeFrom    *int    `yaml:"age_from"`
	Months     int     `yaml:"months"`
	EndsAtAge  *int    `yaml:"ends_at_age"`
	EndOfMonth bool    `yaml:"end_of_month"`
	Ref        RefList `yaml:"ref"`
}

// OnFail says what the holder must do when the evaluation fails.
type OnFail struct {
	What string  `yaml:"what"`
	Ref  RefList `yaml:"ref"`
}
