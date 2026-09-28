package main

import (
	"fmt"
	"path/filepath"
	"slices"
	"sort"

	"github.com/fjaeckel/open-aircrew-rules/credentials"
	"github.com/fjaeckel/open-aircrew-rules/engine"
)

// check runs every validation and gate step and returns the report.
func check(o options) (*report, error) {
	r := &report{vocabUnimpl: map[string][]string{}}
	sch, err := loadSchemas(o.root)
	if err != nil {
		return nil, err
	}
	validate := func(name, path string) {
		r.files++
		r.schema = append(r.schema, sch.validate(name, path)...)
	}
	validate("vocabulary", filepath.Join(o.root, "vocabulary.yaml"))
	validate("messages", filepath.Join(o.root, "messages", "keys.yaml"))
	if len(r.schema) > 0 {
		return r, nil
	}
	files, err := engine.YAMLFiles(filepath.Join(o.root, "credentials"))
	if err != nil {
		return nil, err
	}
	for _, f := range files {
		name := "credential"
		if filepath.Base(filepath.Dir(f)) == "shared" {
			name = "shared"
			r.shared++
		}
		validate(name, f)
	}
	exFiles, err := engine.YAMLFiles(filepath.Join(o.root, "examples"))
	if err != nil {
		return nil, err
	}
	for _, f := range exFiles {
		validate("examples", f)
	}
	scopeFiles, _ := filepath.Glob(filepath.Join(o.root, "scope", "articles-*.yaml"))
	for _, f := range scopeFiles {
		validate("scope", f)
	}
	validate("coverage", filepath.Join(o.root, "coverage", "articles.yaml"))

	cat, err := credentials.Load(o.root)
	if err != nil {
		return nil, err
	}
	ex, err := credentials.LoadExamples(o.root)
	if err != nil {
		return nil, err
	}
	r.f = credentials.Check(cat, ex)
	r.cov = credentials.CheckCoverage(cat)
	r.vocabUnimpl = unimplemented(cat.Engine.Vocabulary)
	r.sources, err = checkSources(o.root, cat.Engine.Vocabulary)
	if err != nil {
		return nil, err
	}
	switch {
	case !o.coverage:
		r.coverage = "not measured (-coverage=false)"
	default:
		pct, err := engineCoverage(o.root)
		switch {
		case err != nil:
			r.coverage, r.coverageFail = "could not measure: "+err.Error(), true
		case pct < o.minCoverage:
			r.coverage, r.coverageFail = fmt.Sprintf("%.1f%% (minimum %.0f%%)", pct, o.minCoverage), true
		default:
			r.coverage = fmt.Sprintf("%.1f%% (minimum %.0f%%)", pct, o.minCoverage)
		}
	}
	return r, nil
}

// unimplemented lists, per vocabulary section, the declared names the engine lacks.
func unimplemented(v *engine.Vocabulary) map[string][]string {
	declared := map[string][]string{
		"subjects":         keys(v.Subjects),
		"metrics":          keys(v.Metrics),
		"filters":          keys(v.Filters),
		"windows":          keys(v.Windows),
		"combinators":      keys(v.Combinators),
		"stage_conditions": keys(v.StageConditions),
		"rule_events":      keys(v.RuleEvents),
		"param_sources":    keys(v.ParamSources),
		"hatches":          keys(v.Hatches),
	}
	out := map[string][]string{}
	impl := engine.Implemented()
	for section, names := range declared {
		for _, n := range names {
			if !slices.Contains(impl[section], n) {
				out[section] = append(out[section], n)
			}
		}
	}
	return out
}

func keys[V any](m map[string]V) []string {
	out := make([]string, 0, len(m))
	for k := range m {
		out = append(out, k)
	}
	sort.Strings(out)
	return out
}
