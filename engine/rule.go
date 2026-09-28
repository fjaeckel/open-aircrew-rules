package engine

import (
	"fmt"
	"reflect"
	"sort"
	"strings"

	"gopkg.in/yaml.v3"
)

// Rule is one rule of the closed vocabulary: what the compiler makes of one credential
// evaluation (DESIGN.md section 4).
type Rule struct {
	ID                 string      `yaml:"id"`
	Title              string      `yaml:"title"`
	Authority          string      `yaml:"authority"`
	EffectiveFrom      *Date       `yaml:"effective_from"`
	EffectiveTo        *Date       `yaml:"effective_to"`
	Citations          []string    `yaml:"citations"`
	AppliesTo          AppliesTo   `yaml:"applies_to"`
	Window             *Window     `yaml:"window"`
	Filter             *Filter     `yaml:"filter"`
	Validity           *Validity   `yaml:"validity"`
	Requirements       *Node       `yaml:"requirements"`
	Stages             []Stage     `yaml:"stages"`
	RestoredBy         []EventHook `yaml:"restored_by"`
	RuleDescriptionKey string      `yaml:"ruleDescriptionKey"`
	Supersedes         []string    `yaml:"supersedes"`

	// File is the path the rule was loaded from.
	File string `yaml:"-"`
}

// deref follows a YAML alias to its anchored node.
func deref(n *yaml.Node) *yaml.Node {
	for n.Kind == yaml.AliasNode {
		n = n.Alias
	}
	return n
}

// AppliesTo selects the subjects a rule evaluates.
type AppliesTo struct {
	Subject             string   `yaml:"subject"`
	Authorities         []string `yaml:"authorities"`
	ExcludeAuthorities  []string `yaml:"excludeAuthorities"`
	LicenceKinds        []string `yaml:"licenceKinds"`
	ExcludeLicenceKinds []string `yaml:"excludeLicenceKinds"`
	Classes             []string `yaml:"classes"`
	ExcludeClasses      []string `yaml:"excludeClasses"`
	ULKinds             []string `yaml:"ulKinds"`
	PrivilegeKinds      []string `yaml:"privilegeKinds"`
	CredentialTypes     []string `yaml:"credentialTypes"`
	LaunchMethods       []string `yaml:"launchMethods"`
	TypeRated           *bool    `yaml:"typeRated"`
	Programme           string   `yaml:"programme"`
	Holds               *Holds   `yaml:"holds"`
	// Categories narrows to subjects whose aircraft category (subjectCategory) is listed.
	Categories []string `yaml:"categories"`
	// DifferentEngineType narrows variants to those recorded as differing (or not) by engine type.
	DifferentEngineType *bool `yaml:"differentEngineType"`
	// ClassGroup makes passengers subjects one per group of vocabulary class_groups.
	ClassGroup string `yaml:"classGroup"`
	// UnknownWhenMissing lists the criteria (typeRated, ulKinds, categories,
	// differentEngineType) for which a selected item that lacks the data is evaluated as
	// unknown (selection.input_missing) instead of being skipped.
	UnknownWhenMissing []string `yaml:"unknownWhenMissing"`
}

// Window is one window kind with its parameter.
type Window struct {
	Kind   string
	N      int
	Anchor string
}

// UnmarshalYAML reads { <kind>: <param> }.
func (w *Window) UnmarshalYAML(n *yaml.Node) error {
	n = deref(n)
	if n.Kind != yaml.MappingNode || len(n.Content) != 2 {
		return fmt.Errorf("line %d: window must be a single-key map", n.Line)
	}
	w.Kind = n.Content[0].Value
	v := n.Content[1]
	switch w.Kind {
	case "since_issue":
		w.Anchor = v.Value
	case "lifetime", "validity_period":
	default:
		if err := v.Decode(&w.N); err != nil {
			return fmt.Errorf("line %d: window %s needs an integer", v.Line, w.Kind)
		}
	}
	return nil
}

// String returns a canonical form, e.g. "calendar_months=6".
func (w Window) String() string {
	switch w.Kind {
	case "since_issue":
		return w.Kind + "=" + w.Anchor
	case "lifetime", "validity_period":
		return w.Kind
	}
	return fmt.Sprintf("%s=%d", w.Kind, w.N)
}

// Filter restricts the items a metric counts; the keys set are remembered for merging.
type Filter struct {
	Classes         []string        `yaml:"classes"`
	ExcludeClasses  []string        `yaml:"excludeClasses"`
	HeldClassPools  [][]string      `yaml:"heldClassPools"`
	Categories      []string        `yaml:"categories"`
	ULKinds         []string        `yaml:"ulKinds"`
	ULCredit        []ULCredit      `yaml:"ulCredit"`
	LaunchMethods   []string        `yaml:"launchMethods"`
	TypeDesignators []string        `yaml:"typeDesignators"`
	Variants        []string        `yaml:"variants"`
	Roles           []string        `yaml:"roles"`
	WithMinutes     []string        `yaml:"withMinutes"`
	WithoutMinutes  []string        `yaml:"withoutMinutes"`
	Simulator       string          `yaml:"simulator"`
	FSTDTypes       []string        `yaml:"fstdTypes"`
	Flags           map[string]bool `yaml:"flags"`
	MinLandings     *int            `yaml:"minLandings"`
	MinDistanceKm   *float64        `yaml:"minDistanceKm"`
	MaxEngines      *int            `yaml:"maxEngines"`
	MaxMTOMKg       *int            `yaml:"maxMtomKg"`
	Tailwheel       *bool           `yaml:"tailwheel"`
	SoleManipulator *bool           `yaml:"soleManipulator"`
	PilotFlying     *bool           `yaml:"pilotFlying"`
	TowKinds        []string        `yaml:"towKinds"`
	TowTakeUps      []string        `yaml:"towTakeUps"`
	ClassGroup      string          `yaml:"classGroup"`
	EventKinds      []string        `yaml:"eventKinds"`
	EventRatings    []string        `yaml:"eventRatings"`
	// ExcludeEventRatings drops events recorded only for these ratings; an event recorded
	// without a rating, or also for another rating, stays.
	ExcludeEventRatings []string  `yaml:"excludeEventRatings"`
	EventAuthorities    []string  `yaml:"eventAuthorities"`
	Any                 []*Filter `yaml:"any"`

	keys map[string]bool
}

// ULCredit adds ultralight flights of some kinds when a class is counted.
type ULCredit struct {
	Class     string   `yaml:"class"`
	ULKinds   []string `yaml:"ulKinds"`
	MinMTOMKg *int     `yaml:"minMtomKg"`
	// FixedEngine credits only flights recorded with a fixed engine and propeller.
	FixedEngine bool `yaml:"fixedEngine"`
}

type filterAlias Filter

// UnmarshalYAML decodes the filter and records which keys were given.
func (f *Filter) UnmarshalYAML(n *yaml.Node) error {
	n = deref(n)
	var a filterAlias
	if err := n.Decode(&a); err != nil {
		return err
	}
	*f = Filter(a)
	f.keys = map[string]bool{}
	for i := 0; i+1 < len(n.Content); i += 2 {
		f.keys[n.Content[i].Value] = true
	}
	return nil
}

// Keys returns the filter keys given, sorted.
func (f *Filter) Keys() []string {
	if f == nil {
		return nil
	}
	out := make([]string, 0, len(f.keys))
	for k := range f.keys {
		out = append(out, k)
	}
	sort.Strings(out)
	return out
}

// Has reports whether key k was given.
func (f *Filter) Has(k string) bool { return f != nil && f.keys[k] }

// MergeFilter returns base overlaid with the keys given in over.
func MergeFilter(base, over *Filter) *Filter {
	if over == nil {
		return base
	}
	if base == nil {
		return over
	}
	out := *base
	out.keys = map[string]bool{}
	for k := range base.keys {
		out.keys[k] = true
	}
	dst, src := reflect.ValueOf(&out).Elem(), reflect.ValueOf(over).Elem()
	for i := range dst.NumField() {
		tag := strings.Split(dst.Type().Field(i).Tag.Get("yaml"), ",")[0]
		if tag != "" && over.keys[tag] {
			dst.Field(i).Set(src.Field(i))
			out.keys[tag] = true
		}
	}
	return &out
}

// Node is a requirement: a leaf (metric) or a combinator.
type Node struct {
	ID            string       `yaml:"id"`
	Metric        string       `yaml:"metric"`
	Min           *float64     `yaml:"min"`
	Unit          string       `yaml:"unit"`
	NameKey       string       `yaml:"nameKey"`
	RemedyKey     string       `yaml:"remedyKey"`
	Messages      *ReqMessages `yaml:"messages"`
	Informational bool         `yaml:"informational"`
	// Max caps a leaf's value (the credit it contributes to a sum_of).
	Max *float64 `yaml:"max"`
	// UnknownIfNone makes an unmet leaf with nothing counted untracked rather than unmet.
	UnknownIfNone bool       `yaml:"unknownIfNone"`
	EscapeHatch   string     `yaml:"escape_hatch"`
	Window        *Window    `yaml:"window"`
	Filter        *Filter    `yaml:"filter"`
	When          *Condition `yaml:"when"`
	AllOf         []*Node    `yaml:"all_of"`
	AnyOf         []*Node    `yaml:"any_of"`
	NOf           *NOf       `yaml:"n_of"`
	// SumOf adds its children's values; the node is met when the sum reaches Min.
	SumOf []*Node `yaml:"sum_of"`
}

// NOf is "at least n of".
type NOf struct {
	N  int     `yaml:"n"`
	Of []*Node `yaml:"of"`
}

// ReqMessages are per-row message keys.
type ReqMessages struct {
	Met       string `yaml:"met"`
	Unmet     string `yaml:"unmet"`
	Untracked string `yaml:"untracked"`
}

// IsLeaf reports whether n is a metric row.
func (n *Node) IsLeaf() bool {
	return n.AllOf == nil && n.AnyOf == nil && n.NOf == nil && n.SumOf == nil
}

// Combinator returns all_of, any_of, n_of or "".
func (n *Node) Combinator() string {
	switch {
	case n.AllOf != nil:
		return "all_of"
	case n.AnyOf != nil:
		return "any_of"
	case n.NOf != nil:
		return "n_of"
	case n.SumOf != nil:
		return "sum_of"
	}
	return ""
}

// Children returns a combinator's children.
func (n *Node) Children() []*Node {
	switch {
	case n.AllOf != nil:
		return n.AllOf
	case n.AnyOf != nil:
		return n.AnyOf
	case n.NOf != nil:
		return n.NOf.Of
	case n.SumOf != nil:
		return n.SumOf
	}
	return nil
}

// Walk visits n and its descendants depth first.
func (n *Node) Walk(fn func(*Node)) {
	if n == nil {
		return
	}
	fn(n)
	for _, c := range n.Children() {
		c.Walk(fn)
	}
}

// Stage sets status and message when its condition holds.
type Stage struct {
	ID         string                 `yaml:"id"`
	When       Condition              `yaml:"when"`
	Status     string                 `yaml:"status"`
	MessageKey string                 `yaml:"messageKey"`
	Params     map[string]ParamSource `yaml:"params"`
}

// Tag returns the stage's coverage name.
func (s Stage) Tag() string {
	if s.ID != "" {
		return s.ID
	}
	return s.Status
}

// ParamSource names where a message param comes from.
type ParamSource struct {
	Source string
	Ref    string
}

// UnmarshalYAML accepts a name or { needed|last_date: <requirement id> }.
func (p *ParamSource) UnmarshalYAML(n *yaml.Node) error {
	n = deref(n)
	if n.Kind == yaml.ScalarNode {
		p.Source = n.Value
		return nil
	}
	if n.Kind != yaml.MappingNode || len(n.Content) != 2 {
		return fmt.Errorf("line %d: param must be a name or a single-key map", n.Line)
	}
	p.Source, p.Ref = n.Content[0].Value, n.Content[1].Value
	return nil
}

// Condition is a stage or requirement condition.
type Condition struct {
	Op    string
	Ref   string
	Days  int
	Span  *Window
	Holds *Holds
	List  []*Condition
	Not   *Condition
	// Seeks lists privileges a recorded training programme seeks (condition seeks).
	Seeks []string
}

// UnmarshalYAML reads a name or a single-key map.
func (c *Condition) UnmarshalYAML(n *yaml.Node) error {
	n = deref(n)
	if n.Kind == yaml.ScalarNode {
		c.Op = n.Value
		return nil
	}
	if n.Kind != yaml.MappingNode || len(n.Content) != 2 {
		return fmt.Errorf("line %d: condition must be a name or a single-key map", n.Line)
	}
	c.Op = n.Content[0].Value
	v := n.Content[1]
	switch c.Op {
	case "met", "unmet", "missing":
		c.Ref = v.Value
	case "seeks":
		if v.Kind == yaml.ScalarNode {
			c.Seeks = []string{v.Value}
			return nil
		}
		return v.Decode(&c.Seeks)
	case "expires_within", "valid_until_within":
		var d struct {
			Days int `yaml:"days"`
		}
		if err := v.Decode(&d); err != nil {
			return err
		}
		c.Days = d.Days
	case "met_within":
		c.Span = &Window{}
		return v.Decode(c.Span)
	case "holds":
		c.Holds = &Holds{}
		return v.Decode(c.Holds)
	case "all", "any":
		return v.Decode(&c.List)
	case "not":
		c.Not = &Condition{}
		return v.Decode(c.Not)
	default:
		return fmt.Errorf("line %d: unknown condition %q", n.Line, c.Op)
	}
	return nil
}

// Walk visits c and nested conditions.
func (c *Condition) Walk(fn func(*Condition)) {
	if c == nil {
		return
	}
	fn(c)
	for _, s := range c.List {
		s.Walk(fn)
	}
	c.Not.Walk(fn)
}

// Holds tests what the holder holds.
type Holds struct {
	Classes      []string `yaml:"classes"`
	ULKinds      []string `yaml:"ulKinds"`
	LicenceKinds []string `yaml:"licenceKinds"`
	Privileges   []string `yaml:"privileges"`
	Credentials  []string `yaml:"credentials"`
	SameLicence  bool     `yaml:"sameLicence"`
	Valid        bool     `yaml:"valid"`
	Every        bool     `yaml:"every"`
	// ExpiryRecorded counts only items with a recorded expiry date.
	ExpiryRecorded bool `yaml:"expiryRecorded"`
	// SameCategory counts only ratings of the subject's aircraft category.
	SameCategory bool `yaml:"sameCategory"`
	// Details narrows privileges by their detail: $subject is the subject's ultralight kind
	// (else class), none a privilege recorded without a detail.
	Details []string `yaml:"details"`
	// Authorities narrows to items on licences of these authorities.
	Authorities []string `yaml:"authorities"`
}

// EventHook is a restored_by entry.
type EventHook struct {
	Event  string  `yaml:"event"`
	Filter *Filter `yaml:"filter"`
}

// Validity derives an expiry date from an anchor date.
type Validity struct {
	From    string   `yaml:"from"`
	Periods []Period `yaml:"periods"`
	// AgeOn is the anchor the holder's age is taken on (issued or valid_from); default From.
	AgeOn string `yaml:"age_on"`
	// RecordedWins uses the derived expiry only when no expiry is recorded.
	RecordedWins bool `yaml:"recorded_wins"`
	// EndOffsetDays moves the derived end (-1: the day before, §§ 187(2), 188(2) BGB).
	EndOffsetDays int `yaml:"end_offset_days"`
}

// Period is one validity period; the first whose conditions hold applies.
type Period struct {
	When       map[string]int `yaml:"when"`
	Months     int            `yaml:"months"`
	CapAtAge   *int           `yaml:"cap_at_age"`
	EndOfMonth bool           `yaml:"end_of_month"`
}
