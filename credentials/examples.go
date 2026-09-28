package credentials

import (
	"fmt"
	"os"
	"path/filepath"
	"slices"

	"gopkg.in/yaml.v3"

	"github.com/fjaeckel/open-aircrew-rules/engine"
)

// ExampleFile is examples/<credential id>.yaml.
type ExampleFile struct {
	Credential string    `yaml:"credential"`
	Examples   []Example `yaml:"examples"`
	File       string    `yaml:"-"`
}

// Example is one worked example: a record, a date and the expected outcome of one evaluation,
// or of the credential's composite (Composite).
type Example struct {
	Name       string        `yaml:"name"`
	Evaluation string        `yaml:"evaluation"`
	Composite  bool          `yaml:"composite"`
	Says       string        `yaml:"says"`
	Shows      []string      `yaml:"shows"`
	AsOf       engine.Date   `yaml:"asOf"`
	Record     engine.Record `yaml:"record"`
	// Outcome is the status the evaluation or composite must report.
	Outcome string        `yaml:"outcome"`
	Expect  ExampleExpect `yaml:"expect"`
}

// ExampleExpect is the rest of the expected result; absent fields are not compared.
type ExampleExpect struct {
	Subject      engine.Subject                `yaml:"subject"`
	Message      string                        `yaml:"message"`
	Params       map[string]any                `yaml:"params"`
	ExpiresOn    *engine.Date                  `yaml:"expiresOn"`
	ValidUntil   *engine.Date                  `yaml:"validUntil"`
	Requirements map[string]engine.ExpectedRow `yaml:"requirements"`
	DecidedBy    *ExpectDecider                `yaml:"decidedBy"`
}

// ExpectDecider names the composite member expected to decide (evaluation id, and for a
// requirement group the required credential).
type ExpectDecider struct {
	Evaluation string `yaml:"evaluation"`
	Credential string `yaml:"credential"`
}

// Outcome classes of a worked example (docs/credential-format.md, "Worked examples").
const (
	Passing = "passing"
	Failing = "failing"
)

// OutcomeClass classifies a result: passing when current, or expiring while the requirement
// tree is met or absent (root "met" or ""; the expiring notice of a valid credential);
// failing otherwise (expired, lapsed, unknown, not_applicable, or expiring with the
// requirements unmet or undetermined).
func OutcomeClass(status, root string) string {
	switch {
	case status == "current":
		return Passing
	case status == "expiring" && (root == "met" || root == ""):
		return Passing
	}
	return Failing
}

// CompositeClass classifies a composite: passing when the credential may be exercised
// (current or expiring), failing otherwise.
func CompositeClass(status string) string {
	if status == "current" || status == "expiring" {
		return Passing
	}
	return Failing
}

// LoadExamples reads examples/*.yaml under root.
func LoadExamples(root string) ([]*ExampleFile, error) {
	files, err := engine.YAMLFiles(filepath.Join(root, "examples"))
	if err != nil {
		return nil, err
	}
	var out []*ExampleFile
	for _, f := range files {
		b, err := os.ReadFile(f)
		if err != nil {
			return nil, err
		}
		var ef ExampleFile
		if err := yaml.Unmarshal(b, &ef); err != nil {
			return nil, fmt.Errorf("%s: %w", f, err)
		}
		ef.File = f
		out = append(out, &ef)
	}
	return out, nil
}

// RunExample evaluates one example and returns its differences.
func (cat *Catalogue) RunExample(credID string, x Example) []string {
	_, d := cat.RunExampleClass(credID, x)
	return d
}

// RunExampleClass evaluates one example and returns the outcome class of the result and the
// differences from what the example expects.
func (cat *Catalogue) RunExampleClass(credID string, x Example) (string, []string) {
	if x.Outcome == "" {
		return "", []string{"outcome is required"}
	}
	if x.Composite {
		return cat.runComposite(credID, x)
	}
	key := credID + "#" + x.Evaluation
	id, ok := cat.ByEvaluation[key]
	if !ok {
		return "", []string{fmt.Sprintf("no evaluation %s", key)}
	}
	if id == "" {
		return "", []string{fmt.Sprintf("%s only names required credentials; show it with a composite example", key)}
	}
	r, ok := cat.Engine.Rule(id)
	if !ok {
		return "", []string{fmt.Sprintf("%s did not compile", key)}
	}
	evs, traces := engine.EvaluateRule(cat.Engine, r, &x.Record, x.AsOf)
	var hit []int
	for i, ev := range evs {
		if engine.SubjectMatches(engine.Subject{Kind: ev.Subject.Kind, ID: x.Expect.Subject.ID, Class: x.Expect.Subject.Class, ULKind: x.Expect.Subject.ULKind, Detail: x.Expect.Subject.Detail}, ev.Subject) {
			hit = append(hit, i)
		}
	}
	switch len(hit) {
	case 0:
		return "", []string{fmt.Sprintf("no result for the subject (evaluated %d subjects)", len(evs))}
	case 1:
	default:
		return "", []string{fmt.Sprintf("%d results; name the subject in expect.subject", len(hit))}
	}
	ev := evs[hit[0]]
	want := engine.Expect{
		Subject: x.Expect.Subject, Status: x.Outcome, MessageKey: x.Expect.Message,
		MessageParams: x.Expect.Params, ExpiresOn: x.Expect.ExpiresOn, ValidUntil: x.Expect.ValidUntil,
		Requirements: x.Expect.Requirements,
	}
	d := engine.DiffExpect(want, ev)
	if x.Expect.DecidedBy != nil {
		d = append(d, "expect.decidedBy belongs to a composite example")
	}
	return OutcomeClass(ev.Status, traces[hit[0]].Root), d
}

func (cat *Catalogue) runComposite(credID string, x Example) (string, []string) {
	if x.Evaluation != "" {
		return "", []string{"a composite example names no evaluation"}
	}
	if x.Expect.Message != "" || x.Expect.Params != nil || x.Expect.ExpiresOn != nil || x.Expect.ValidUntil != nil || x.Expect.Requirements != nil {
		return "", []string{"a composite example expects only subject and decidedBy"}
	}
	res := cat.Evaluate(&x.Record, x.AsOf)
	var hit []Composite
	for _, c := range res.Credentials {
		if c.Credential == credID && engine.SubjectMatches(engine.Subject{Kind: c.Subject.Kind, ID: x.Expect.Subject.ID}, c.Subject) {
			hit = append(hit, c)
		}
	}
	switch len(hit) {
	case 0:
		return "", []string{fmt.Sprintf("the record holds no %s for the subject", credID)}
	case 1:
	default:
		return "", []string{fmt.Sprintf("the record holds %d; name the subject in expect.subject", len(hit))}
	}
	c := hit[0]
	var d []string
	if c.Status != x.Outcome {
		d = append(d, fmt.Sprintf("status: want %s, got %s", x.Outcome, c.Status))
	}
	if w := x.Expect.DecidedBy; w != nil {
		got := ExpectDecider{}
		if c.DecidedBy != nil {
			got = ExpectDecider{Evaluation: c.DecidedBy.Evaluation, Credential: c.DecidedBy.Credential}
		}
		if w.Evaluation != got.Evaluation || (w.Credential != "" && w.Credential != got.Credential) {
			d = append(d, fmt.Sprintf("decidedBy: want %s %s, got %s %s", w.Evaluation, w.Credential, got.Evaluation, got.Credential))
		}
	}
	return CompositeClass(c.Status), d
}

// Statuses returns the statuses a compiled rule's stages can report.
func Statuses(r *engine.Rule) []string {
	var out []string
	for _, s := range r.Stages {
		if !slices.Contains(out, s.Status) {
			out = append(out, s.Status)
		}
	}
	return out
}
