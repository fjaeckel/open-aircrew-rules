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

// Example is one worked example: a record, a date and the expected result of one evaluation.
type Example struct {
	Name       string        `yaml:"name"`
	Evaluation string        `yaml:"evaluation"`
	Says       string        `yaml:"says"`
	Shows      []string      `yaml:"shows"`
	AsOf       engine.Date   `yaml:"asOf"`
	Record     engine.Record `yaml:"record"`
	Expect     ExampleExpect `yaml:"expect"`
}

// ExampleExpect is the expected result; absent fields are not compared.
type ExampleExpect struct {
	Subject      engine.Subject                `yaml:"subject"`
	Status       string                        `yaml:"status"`
	Message      string                        `yaml:"message"`
	Params       map[string]any                `yaml:"params"`
	ExpiresOn    *engine.Date                  `yaml:"expiresOn"`
	ValidUntil   *engine.Date                  `yaml:"validUntil"`
	Requirements map[string]engine.ExpectedRow `yaml:"requirements"`
}

// Outcome classifies an example: pass when the expected status is current, fail for any
// other status, "" without one.
func (x Example) Outcome() string {
	switch x.Expect.Status {
	case "":
		return ""
	case "current":
		return "pass"
	}
	return "fail"
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
	key := credID + "#" + x.Evaluation
	id, ok := cat.ByEvaluation[key]
	if !ok {
		return []string{fmt.Sprintf("no evaluation %s", key)}
	}
	if id == "" {
		return []string{fmt.Sprintf("%s only references other credentials; it has no examples", key)}
	}
	r, ok := cat.Engine.Rule(id)
	if !ok {
		return []string{fmt.Sprintf("%s did not compile", key)}
	}
	if x.Expect.Status == "" {
		return []string{"expect.status is required"}
	}
	evs, _ := engine.EvaluateRule(cat.Engine, r, &x.Record, x.AsOf)
	var hit []engine.Evaluation
	for _, ev := range evs {
		if engine.SubjectMatches(engine.Subject{Kind: ev.Subject.Kind, ID: x.Expect.Subject.ID, Class: x.Expect.Subject.Class, ULKind: x.Expect.Subject.ULKind, Detail: x.Expect.Subject.Detail}, ev.Subject) {
			hit = append(hit, ev)
		}
	}
	switch len(hit) {
	case 0:
		return []string{fmt.Sprintf("no result for the subject (evaluated %d subjects)", len(evs))}
	case 1:
	default:
		return []string{fmt.Sprintf("%d results; name the subject in expect.subject", len(hit))}
	}
	want := engine.Expect{
		Subject: x.Expect.Subject, Status: x.Expect.Status, MessageKey: x.Expect.Message,
		MessageParams: x.Expect.Params, ExpiresOn: x.Expect.ExpiresOn, ValidUntil: x.Expect.ValidUntil,
		Requirements: x.Expect.Requirements,
	}
	return engine.DiffExpect(want, hit[0])
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
