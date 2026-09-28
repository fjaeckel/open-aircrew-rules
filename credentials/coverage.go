package credentials

import (
	"fmt"
	"os"
	"path/filepath"
	"slices"
	"sort"
	"strings"

	"gopkg.in/yaml.v3"
)

// Coverage is coverage/articles.yaml.
type Coverage struct {
	Articles []CoverageEntry `yaml:"articles"`
}

// CoverageEntry accounts for one article in scope: the evaluations that encode it, the
// credential files planned to encode it (pending, with a note), and why (the rest of) it is
// not evaluated. At least one of the three is given.
type CoverageEntry struct {
	Article      string   `yaml:"article"`
	Evaluations  []string `yaml:"evaluations"`
	Pending      []string `yaml:"pending"`
	Note         string   `yaml:"note"`
	NotEvaluated string   `yaml:"not_evaluated"`
}

// CoverageFragment is one file under fragments/coverage/: changes to coverage/articles.yaml
// that the integration step merges.
type CoverageFragment struct {
	Articles []struct {
		Article      string   `yaml:"article"`
		Evaluations  []string `yaml:"evaluations"`
		Pending      []string `yaml:"pending"`
		DonePending  []string `yaml:"done_pending"`
		Note         string   `yaml:"note"`
		NotEvaluated string   `yaml:"not_evaluated"`
	} `yaml:"articles"`
}

// mergeCoverage applies one coverage fragment: evaluations and pending are added,
// done_pending removed from pending, note and not_evaluated replaced when given; an article
// not listed yet is added. A note without pending files left is dropped.
func mergeCoverage(cov *Coverage, file string) []string {
	var probs []string
	b, err := os.ReadFile(file)
	var fr CoverageFragment
	if err == nil {
		err = yaml.Unmarshal(b, &fr)
	}
	if err != nil {
		return []string{fmt.Sprintf("%s: %v", file, err)}
	}
	for _, a := range fr.Articles {
		i := slices.IndexFunc(cov.Articles, func(e CoverageEntry) bool { return e.Article == a.Article })
		if i < 0 {
			cov.Articles = append(cov.Articles, CoverageEntry{Article: a.Article})
			i = len(cov.Articles) - 1
		}
		e := &cov.Articles[i]
		e.Evaluations = append(e.Evaluations, a.Evaluations...)
		for _, d := range a.DonePending {
			j := slices.Index(e.Pending, d)
			if j < 0 {
				probs = append(probs, fmt.Sprintf("%s: %s: done_pending %s is not pending", file, a.Article, d))
				continue
			}
			e.Pending = slices.Delete(e.Pending, j, j+1)
		}
		for _, p := range a.Pending {
			if !slices.Contains(e.Pending, p) {
				e.Pending = append(e.Pending, p)
			}
		}
		if a.Note != "" {
			e.Note = a.Note
		}
		if len(e.Pending) == 0 {
			e.Note = ""
		}
		if a.NotEvaluated != "" {
			e.NotEvaluated = a.NotEvaluated
		}
	}
	return probs
}

// CoverageResult is the outcome of CheckCoverage.
type CoverageResult struct {
	// Problems fail the gate.
	Problems []string
	// Pending lists the articles still waiting for a credential ("<article>: <files>").
	Pending []string
	// Evaluated, PendingOnly and NotEvaluated count the articles by how they are accounted
	// for, in that precedence: an article with evaluations counts as evaluated, one with
	// pending files and no evaluation as pending.
	Evaluated, NotEvaluated, PendingOnly int
	// PendingByAuthority counts articles with pending files per reference prefix.
	PendingByAuthority map[string]int
}

// ScopeRef turns a scope cite ("FCL.740.A", "14 CFR 61.57", "LuftPersV § 45a") into an
// article reference of the given authority.
func ScopeRef(authority, cite string) string {
	switch authority {
	case "faa":
		return "faa:" + strings.TrimPrefix(cite, "14 CFR ")
	case "de":
		return "de:" + strings.ReplaceAll(strings.ReplaceAll(cite, " § ", "."), " ", "")
	}
	return authority + ":" + cite
}

// ArticleOf returns the article part of a reference ("easa:FCL.740.A(b)(1)" -> "easa:FCL.740.A").
func ArticleOf(ref string) string {
	r, err := ParseRef(ref)
	if err != nil {
		return ref
	}
	return r.Prefix + ":" + r.Article
}

// LoadScope reads scope/articles-*.yaml and returns the article references in scope.
func LoadScope(root string) (map[string]bool, error) {
	scope := map[string]bool{}
	files, _ := filepath.Glob(filepath.Join(root, "scope", "articles-*.yaml"))
	for _, f := range files {
		var doc struct {
			Articles []struct {
				Cite string `yaml:"cite"`
			} `yaml:"articles"`
		}
		b, err := os.ReadFile(f)
		if err == nil {
			err = yaml.Unmarshal(b, &doc)
		}
		if err != nil {
			return nil, fmt.Errorf("%s: %w", f, err)
		}
		auth := strings.TrimSuffix(strings.TrimPrefix(filepath.Base(f), "articles-"), ".yaml")
		for _, a := range doc.Articles {
			scope[ScopeRef(auth, a.Cite)] = true
		}
	}
	return scope, nil
}

// CheckCoverage checks coverage/articles.yaml against the articles in scope and the
// compiled evaluations: every article in scope is listed once and accounted for, every
// evaluation exists, every pending credential file is still to be written, and every
// compiled evaluation is listed under the article of its source.
func CheckCoverage(cat *Catalogue) *CoverageResult {
	res := &CoverageResult{PendingByAuthority: map[string]int{}}
	fail := func(format string, args ...any) { res.Problems = append(res.Problems, fmt.Sprintf(format, args...)) }
	b, err := os.ReadFile(filepath.Join(cat.Root, "coverage", "articles.yaml"))
	if err != nil {
		fail("coverage/articles.yaml: %v", err)
		return res
	}
	var cov Coverage
	if err := yaml.Unmarshal(b, &cov); err != nil {
		fail("coverage/articles.yaml: %v", err)
		return res
	}
	if cat.Options.Fragments {
		frags, err := fragmentFiles(cat.Root, "coverage")
		if err != nil {
			fail("fragments/coverage: %v", err)
		}
		for _, f := range frags {
			for _, p := range mergeCoverage(&cov, f) {
				fail("%s", p)
			}
		}
	}
	scope, err := LoadScope(cat.Root)
	if err != nil {
		fail("%v", err)
	}
	listed := map[string][]string{}
	for _, a := range cov.Articles {
		where := "coverage " + a.Article
		if listed[a.Article] != nil {
			fail("%s: listed twice", where)
		}
		listed[a.Article] = append([]string{}, a.Evaluations...)
		if !scope[a.Article] {
			fail("%s: not an article of scope/articles-*.yaml", where)
		}
		if err := cat.Resolver.Resolve(a.Article); err != nil {
			fail("%s: %v", where, err)
		}
		switch {
		case len(a.Evaluations)+len(a.Pending) == 0 && a.NotEvaluated == "":
			fail("%s: needs evaluations, pending or not_evaluated", where)
		case len(a.Pending) > 0 && a.Note == "":
			fail("%s: pending needs a note", where)
		}
		for _, e := range a.Evaluations {
			if id, ok := cat.ByEvaluation[e]; !ok || id == "" {
				fail("%s: %s is no compiled evaluation", where, e)
			}
		}
		for _, p := range a.Pending {
			if _, err := os.Stat(filepath.Join(cat.Root, "credentials", filepath.FromSlash(p)+".yaml")); err == nil {
				fail("%s: pending credentials/%s.yaml exists; list its evaluations instead", where, p)
			}
		}
		switch {
		case len(a.Evaluations) > 0:
			res.Evaluated++
		case len(a.Pending) > 0:
			res.PendingOnly++
		default:
			res.NotEvaluated++
		}
		if len(a.Pending) > 0 {
			prefix, _, _ := strings.Cut(a.Article, ":")
			res.PendingByAuthority[prefix]++
			res.Pending = append(res.Pending, a.Article+": "+strings.Join(a.Pending, ", "))
		}
	}
	for a := range scope {
		if listed[a] == nil {
			fail("coverage: article %s of the scope is missing", a)
		}
	}
	for _, c := range cat.Credentials {
		for _, e := range c.Evaluations {
			key := c.ID + "#" + e.ID
			comp := cat.Compiled[cat.ByEvaluation[key]]
			if comp == nil {
				continue
			}
			art := ArticleOf(comp.Eval.Source)
			if !slices.Contains(listed[art], key) {
				fail("coverage %s: does not list %s, whose source it is", art, key)
			}
		}
	}
	sort.Strings(res.Problems)
	sort.Strings(res.Pending)
	return res
}
