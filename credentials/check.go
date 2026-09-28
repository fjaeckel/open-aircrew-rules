package credentials

import (
	"fmt"
	"path/filepath"
	"regexp"
	"slices"
	"sort"
	"strings"

	"github.com/fjaeckel/open-aircrew-rules/engine"
)

// Findings is the structural gate over the credential catalogue.
type Findings struct {
	Credentials int
	Evaluations int
	Compiled    int
	Examples    int
	Load        []string
	Files       []string
	Refs        []string
	Keys        []string
	ExampleErrs []string
	Coverage    []string
	Overlaps    []string
	Interpret   []string
	Unapproved  []string
	RefsChecked int
}

// Problems counts the failures (unapproved interpretations are reported, not counted).
func (f *Findings) Problems() int {
	return len(f.Load) + len(f.Files) + len(f.Refs) + len(f.Keys) + len(f.ExampleErrs) + len(f.Coverage) + len(f.Overlaps) + len(f.Interpret)
}

// kindDirs maps a credential kind to its directory under credentials/<authority>/.
var kindDirs = map[string]string{
	"licence": "licences", "rating": "ratings", "privilege": "privileges", "endorsement": "endorsements",
	"instructor_certificate": "instructors", "examiner_certificate": "examiners", "medical": "medicals",
}

var idShape = regexp.MustCompile(`^[a-z]+\.[a-z_]+\.[a-z0-9-]+$`)

// Check runs every structural check.
func Check(cat *Catalogue, examples []*ExampleFile) *Findings {
	f := &Findings{Credentials: len(cat.Credentials), Compiled: len(cat.Compiled), Load: slices.Clone(cat.Errors)}
	f.checkFiles(cat)
	for _, r := range cat.Refs {
		f.RefsChecked++
		if err := cat.Resolver.Resolve(r.Ref); err != nil {
			f.Refs = append(f.Refs, fmt.Sprintf("%s (line %d): %v", r.Where, r.Line, err))
		}
	}
	f.checkKeys(cat)
	f.checkExamples(cat, examples)
	f.checkInterpretations(cat)
	f.checkPolicies(cat)
	f.Overlaps = Overlaps(cat)
	sort.Strings(f.Refs)
	sort.Strings(f.Keys)
	return f
}

func (f *Findings) checkFiles(cat *Catalogue) {
	for _, c := range cat.Credentials {
		where := rel(cat.Root, c.File)
		dir := kindDirs[c.Kind]
		switch {
		case dir == "":
			f.Files = append(f.Files, fmt.Sprintf("%s: kind %q is not one of the credential kinds", where, c.Kind))
		case filepath.Base(filepath.Dir(c.File)) != dir:
			f.Files = append(f.Files, fmt.Sprintf("%s: a %s lives under %s/", where, c.Kind, dir))
		}
		if !idShape.MatchString(c.ID) || !strings.HasSuffix(c.ID, "."+strings.TrimSuffix(filepath.Base(c.File), ".yaml")) {
			f.Files = append(f.Files, fmt.Sprintf("%s: id %q must be <authority>.<kind>.<file name>", where, c.ID))
		}
		if !strings.HasPrefix(c.ID, strings.ToLower(c.Authority)+".") {
			f.Files = append(f.Files, fmt.Sprintf("%s: id %q does not start with the authority", where, c.ID))
		}
		if c.Name == "" || len(c.Evaluations) == 0 {
			f.Files = append(f.Files, fmt.Sprintf("%s: a credential has a name and at least one evaluation", where))
		}
		if c.Validity != nil {
			cat.addRef(c.Validity.Ref, c.ID+" validity", 0)
		}
		seen := map[string]bool{}
		for _, e := range c.Evaluations {
			f.Evaluations++
			if seen[e.ID] || e.ID == "" {
				f.Files = append(f.Files, fmt.Sprintf("%s: evaluation id %q missing or repeated", where, e.ID))
			}
			seen[e.ID] = true
			if e.Asks == "" && e.Uses == "" {
				f.Files = append(f.Files, fmt.Sprintf("%s#%s: evaluation has no asks", c.ID, e.ID))
			}
			key := c.ID + "#" + e.ID
			if id, ok := cat.ByEvaluation[key]; ok && id != "" && cat.Compiled[id] != nil && cat.Compiled[id].UsesWith && len(cat.interpretationsFor(c, e)) == 0 {
				f.Interpret = append(f.Interpret, fmt.Sprintf("%s uses `with:` (not evaluated) but no interpretation affects it", key))
			}
		}
	}
	for id, s := range cat.Shared {
		where := rel(cat.Root, s.File)
		if !strings.HasSuffix(id, "."+strings.TrimSuffix(filepath.Base(s.File), ".yaml")) || !strings.Contains(id, ".shared.") {
			f.Files = append(f.Files, fmt.Sprintf("%s: id %q must be <authority>.shared.<file name>", where, id))
		}
		used := false
		for _, c := range cat.Compiled {
			used = used || (c.Eval.Shared == s)
		}
		if !used {
			f.Files = append(f.Files, fmt.Sprintf("%s: no credential uses %s", where, id))
		}
	}
	sort.Strings(f.Files)
}

// interpretationsFor returns the interpretations affecting an evaluation: in its own file,
// in the shared file it uses, or in the credential whose evaluation it reuses.
func (cat *Catalogue) interpretationsFor(c *Credential, e *Evaluation) []string {
	var out []string
	add := func(l []Interpretation, id string) {
		for _, i := range l {
			if slices.Contains(i.Affects, id) || slices.Contains(i.Affects, "*") {
				out = append(out, i.ID)
			}
		}
	}
	add(c.Interpretations, e.ID)
	if e.Uses == "" {
		return out
	}
	if credID, evalID, ok := strings.Cut(e.Uses, "#"); ok {
		if oc, found := cat.byID[credID]; found {
			for _, x := range oc.Evaluations {
				if x.ID == evalID {
					out = append(out, cat.interpretationsFor(oc, x)...)
				}
			}
		}
	} else if s, found := cat.Shared[e.Uses]; found {
		add(s.Interpretations, "evaluation")
	}
	return out
}

func (f *Findings) checkKeys(cat *Catalogue) {
	check := func(id, what, key string) {
		if key == "" {
			return
		}
		if _, ok := cat.Engine.Keys.Get(key); !ok {
			f.Keys = append(f.Keys, fmt.Sprintf("%s: %s %s is not in messages/keys.yaml", id, what, key))
		}
	}
	for id, c := range cat.Compiled {
		r := c.Rule
		check(id, "description key", r.RuleDescriptionKey)
		for _, s := range r.Stages {
			check(id, "message", s.MessageKey)
			if _, ok := cat.Engine.Vocabulary.Statuses[s.Status]; !ok {
				f.Keys = append(f.Keys, fmt.Sprintf("%s: status %q is not in the vocabulary", id, s.Status))
			}
		}
		r.Requirements.Walk(func(n *engine.Node) {
			check(id, "name", n.NameKey)
			check(id, "remedy", n.RemedyKey)
			if n.Messages != nil {
				check(id, "row message", n.Messages.Met)
				check(id, "row message", n.Messages.Unmet)
				check(id, "row message", n.Messages.Untracked)
			}
		})
	}
}

func (f *Findings) checkExamples(cat *Catalogue, files []*ExampleFile) {
	seen := map[string][]string{}
	for _, ef := range files {
		where := rel(cat.Root, ef.File)
		c, ok := cat.byID[ef.Credential]
		if !ok {
			f.ExampleErrs = append(f.ExampleErrs, fmt.Sprintf("%s: no credential %s", where, ef.Credential))
			continue
		}
		if base := strings.TrimSuffix(filepath.Base(ef.File), ".yaml"); base != c.ID {
			f.ExampleErrs = append(f.ExampleErrs, fmt.Sprintf("%s: file name must be %s.yaml", where, c.ID))
		}
		names := map[string]bool{}
		for _, x := range ef.Examples {
			f.Examples++
			if names[x.Name] || x.Name == "" {
				f.ExampleErrs = append(f.ExampleErrs, fmt.Sprintf("%s: example name %q missing or repeated", where, x.Name))
			}
			names[x.Name] = true
			class, diffs := cat.RunExampleClass(c.ID, x)
			for _, d := range diffs {
				f.ExampleErrs = append(f.ExampleErrs, fmt.Sprintf("%s: %s: %s", where, x.Name, d))
			}
			key := c.ID + "#" + x.Evaluation
			if x.Composite {
				key = compositeKey
			}
			seen[key] = append(seen[key], class)
			var interps []string
			for _, e := range c.Evaluations {
				if x.Composite || e.ID == x.Evaluation {
					interps = append(interps, cat.interpretationsFor(c, e)...)
				}
			}
			for _, s := range x.Shows {
				if !slices.Contains(interps, s) {
					f.ExampleErrs = append(f.ExampleErrs, fmt.Sprintf("%s: %s: shows %q, which is no interpretation of %s", where, x.Name, s, key))
				}
			}
		}
	}
	needsComposite := false
	for _, c := range cat.Credentials {
		for _, e := range c.Evaluations {
			needsComposite = needsComposite || e.isRequirement()
			key := c.ID + "#" + e.ID
			id := cat.ByEvaluation[key]
			comp := cat.Compiled[id]
			if id == "" || comp == nil {
				continue
			}
			st := Statuses(comp.Rule)
			if slices.ContainsFunc(st, func(s string) bool { return s == "current" || s == "expiring" }) && !slices.Contains(seen[key], Passing) {
				f.ExampleErrs = append(f.ExampleErrs, fmt.Sprintf("%s: no passing example", key))
			}
			if slices.ContainsFunc(st, func(s string) bool { return s != "current" && s != "expiring" }) && !slices.Contains(seen[key], Failing) {
				f.ExampleErrs = append(f.ExampleErrs, fmt.Sprintf("%s: no failing example", key))
			}
		}
	}
	if needsComposite {
		for _, class := range []string{Passing, Failing} {
			if !slices.Contains(seen[compositeKey], class) {
				f.ExampleErrs = append(f.ExampleErrs, fmt.Sprintf("composites: no %s composite example in the catalogue", class))
			}
		}
	}
	sort.Strings(f.ExampleErrs)
}

// compositeKey collects the classes of composite examples across the catalogue.
const compositeKey = "(composite)"

func (f *Findings) checkInterpretations(cat *Catalogue) {
	check := func(file string, evals []string, l []Interpretation) {
		ids := map[string]bool{}
		for _, i := range l {
			where := rel(cat.Root, file) + ": interpretation " + i.ID
			if i.ID == "" || ids[i.ID] {
				f.Interpret = append(f.Interpret, where+": id missing or repeated")
			}
			ids[i.ID] = true
			if i.Reading == "" || len(i.Affects) == 0 {
				f.Interpret = append(f.Interpret, where+": needs reading and affects")
			}
			for _, a := range i.Affects {
				if a != "*" && !slices.Contains(evals, a) {
					f.Interpret = append(f.Interpret, fmt.Sprintf("%s: affects %q, which is no evaluation here", where, a))
				}
			}
			if len(i.Ref) == 0 {
				f.Interpret = append(f.Interpret, where+": no ref to the paragraph it interprets")
			}
			for _, r := range i.Ref {
				if err := cat.Resolver.Resolve(r); err != nil {
					f.Refs = append(f.Refs, where+": "+err.Error())
				}
			}
			if (i.ApprovedBy == nil) != (i.ApprovedOn == nil) {
				f.Interpret = append(f.Interpret, where+": approved_by and approved_on go together")
			}
			if i.ApprovedBy == nil {
				name := strings.TrimSuffix(strings.TrimPrefix(rel(cat.Root, file), "credentials/"), ".yaml")
				f.Unapproved = append(f.Unapproved, fmt.Sprintf("%s#%s: %s", name, i.ID, i.Reading))
			}
		}
	}
	for _, c := range cat.Credentials {
		var evals []string
		for _, e := range c.Evaluations {
			evals = append(evals, e.ID)
		}
		check(c.File, evals, c.Interpretations)
	}
	for _, s := range cat.Shared {
		check(s.File, []string{"evaluation"}, s.Interpretations)
	}
	sort.Strings(f.Interpret)
	sort.Strings(f.Unapproved)
}

// checkPolicies reports declared policies that no reference cites.
func (f *Findings) checkPolicies(cat *Catalogue) {
	used := map[string]bool{}
	cite := func(r string) {
		if id, ok := strings.CutPrefix(r, "policy:"); ok {
			used[id] = true
		}
	}
	for _, r := range cat.Refs {
		cite(r.Ref)
	}
	for _, x := range cat.Interpretations() {
		for _, r := range x.I.Ref {
			cite(r)
		}
	}
	for id, p := range cat.Policies {
		if !used[id] {
			f.Refs = append(f.Refs, fmt.Sprintf("%s: policy %s is cited nowhere; cite it or remove it", rel(cat.Root, p.File), id))
		}
	}
}
