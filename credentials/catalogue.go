package credentials

import (
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"gopkg.in/yaml.v3"

	"github.com/fjaeckel/open-aircrew-rules/engine"
)

// Catalogue is the loaded and compiled credential catalogue.
type Catalogue struct {
	Root        string
	Vocab       *Vocab
	Engine      *engine.Catalogue
	Credentials []*Credential
	Shared      map[string]*Shared
	// Compiled maps a compiled rule id to what it was compiled from.
	Compiled map[string]*Compiled
	// ByEvaluation maps "<credential id>#<evaluation id>" to its compiled rule id ("" for
	// evaluations that only reference other credentials).
	ByEvaluation map[string]string
	// Errors lists load and compile problems; the evaluations concerned are skipped.
	Errors []string
	// Refs lists every reference used, with where.
	Refs     []RefUse
	Resolver *Resolver
	byID     map[string]*Credential
}

// Compiled is one compiled rule and the evaluations that use it.
type Compiled struct {
	Rule     *engine.Rule
	Eval     *Evaluation
	Users    []string
	UsesWith bool
}

// RefUse is one reference occurrence.
type RefUse struct {
	Ref   string
	Where string
	Line  int
}

func (cat *Catalogue) addRef(r, where string, line int) {
	cat.Refs = append(cat.Refs, RefUse{Ref: r, Where: where, Line: line})
}

// Credential returns the credential with the id.
func (cat *Catalogue) Credential(id string) (*Credential, bool) {
	c, ok := cat.byID[id]
	return c, ok
}

// Load reads vocabulary.yaml, messages/keys.yaml and credentials/** under root and
// compiles every evaluation.
func Load(root string) (*Catalogue, error) {
	ev, err := engine.LoadVocabulary(filepath.Join(root, "vocabulary.yaml"))
	if err != nil {
		return nil, err
	}
	keys, err := engine.LoadKeys(filepath.Join(root, "messages", "keys.yaml"))
	if err != nil {
		return nil, err
	}
	v, err := LoadVocab(filepath.Join(root, "vocabulary.yaml"))
	if err != nil {
		return nil, err
	}
	cat := &Catalogue{Root: root, Vocab: v, Shared: map[string]*Shared{}, Compiled: map[string]*Compiled{}, ByEvaluation: map[string]string{}, byID: map[string]*Credential{}}
	cat.Resolver = NewResolver(root, v)
	files, err := engine.YAMLFiles(filepath.Join(root, "credentials"))
	if err != nil {
		return nil, err
	}
	for _, f := range files {
		b, err := os.ReadFile(f)
		if err != nil {
			return nil, err
		}
		if filepath.Base(filepath.Dir(f)) == "shared" {
			var s Shared
			if err := yaml.Unmarshal(b, &s); err != nil {
				cat.Errors = append(cat.Errors, fmt.Sprintf("%s: %v", f, err))
				continue
			}
			s.File = f
			if _, dup := cat.Shared[s.ID]; dup {
				cat.Errors = append(cat.Errors, fmt.Sprintf("%s: duplicate shared id %s", f, s.ID))
				continue
			}
			cat.Shared[s.ID] = &s
			continue
		}
		c, err := decodeCredential(b)
		if err != nil {
			cat.Errors = append(cat.Errors, fmt.Sprintf("%s: %v", f, err))
			continue
		}
		c.File = f
		if _, dup := cat.byID[c.ID]; dup {
			cat.Errors = append(cat.Errors, fmt.Sprintf("%s: duplicate credential id %s", f, c.ID))
			continue
		}
		cat.byID[c.ID] = c
		cat.Credentials = append(cat.Credentials, c)
	}
	var rules []*engine.Rule
	for _, c := range cat.Credentials {
		for _, e := range c.Evaluations {
			key := c.ID + "#" + e.ID
			re, err := cat.resolve(c, e, 0)
			if err != nil {
				cat.Errors = append(cat.Errors, fmt.Sprintf("%s: %v", key, err))
				continue
			}
			if len(re.Requires) > 0 && empty(&re.PassesIf) && empty(&re.Outcomes) {
				cat.ByEvaluation[key] = ""
				for _, id := range re.Requires {
					if _, ok := cat.byID[id]; !ok {
						cat.Errors = append(cat.Errors, fmt.Sprintf("%s: requires %s, which is no credential", key, id))
					}
				}
				continue
			}
			id := key
			if re.Shared != nil && re.Scope != nil {
				id = re.Shared.ID
			}
			cat.ByEvaluation[key] = id
			if prev, ok := cat.Compiled[id]; ok {
				prev.Users = append(prev.Users, key)
				continue
			}
			r, cp := cat.compileEvaluation(c, re, id)
			cat.Errors = append(cat.Errors, cp.errs...)
			if len(cp.errs) > 0 {
				continue
			}
			cat.Compiled[id] = &Compiled{Rule: r, Eval: re, Users: []string{key}, UsesWith: cp.usesWith}
			rules = append(rules, r)
		}
	}
	sort.Slice(rules, func(i, j int) bool { return rules[i].ID < rules[j].ID })
	cat.Engine = engine.NewCatalogue(root, ev, keys, rules)
	for _, e := range cat.Engine.Errors {
		cat.Errors = append(cat.Errors, e.Error())
	}
	return cat, nil
}

func decodeCredential(b []byte) (*Credential, error) {
	var doc yaml.Node
	if err := yaml.Unmarshal(b, &doc); err != nil {
		return nil, err
	}
	var c Credential
	if err := doc.Decode(&c); err != nil {
		return nil, err
	}
	for _, p := range pairs(doc.Content[0]) {
		if p[0].Value != "evaluations" {
			continue
		}
		for i, item := range p[1].Content {
			if i < len(c.Evaluations) {
				c.Evaluations[i].Line = item.Line
				c.Evaluations[i].Owner = &c
			}
		}
	}
	return &c, nil
}

// resolve returns the evaluation with `uses` expanded.
func (cat *Catalogue) resolve(c *Credential, e *Evaluation, depth int) (*Evaluation, error) {
	if e.Uses == "" {
		return e, nil
	}
	if depth > 3 {
		return nil, fmt.Errorf("uses %s: too deep", e.Uses)
	}
	var base *Evaluation
	if credID, evalID, ok := strings.Cut(e.Uses, "#"); ok {
		oc, found := cat.byID[credID]
		if !found {
			return nil, fmt.Errorf("uses %s: no credential %s", e.Uses, credID)
		}
		for _, x := range oc.Evaluations {
			if x.ID == evalID {
				b, err := cat.resolve(oc, x, depth+1)
				if err != nil {
					return nil, err
				}
				cp := *b
				base = &cp
			}
		}
		if base == nil {
			return nil, fmt.Errorf("uses %s: %s has no evaluation %s", e.Uses, credID, evalID)
		}
	} else {
		s, found := cat.Shared[e.Uses]
		if !found {
			return nil, fmt.Errorf("uses %s: no shared evaluation", e.Uses)
		}
		b, err := instantiate(s, &e.With)
		if err != nil {
			return nil, fmt.Errorf("uses %s: %v", e.Uses, err)
		}
		base = b
		base.Shared = s
	}
	base.ID = e.ID
	base.Uses = ""
	base.From = e.Uses
	base.Owner = c
	base.Line = e.Line
	if e.Asks != "" {
		base.Asks = e.Asks
	}
	if e.OnlyFor != nil {
		base.OnlyFor = e.OnlyFor
	}
	if e.DescriptionKey != "" {
		base.DescriptionKey = e.DescriptionKey
	}
	if e.EffectiveFrom != nil {
		base.EffectiveFrom = e.EffectiveFrom
	}
	if e.EffectiveTo != nil {
		base.EffectiveTo = e.EffectiveTo
	}
	return base, nil
}

// instantiate substitutes the `with` parameters into a shared evaluation.
func instantiate(s *Shared, with *yaml.Node) (*Evaluation, error) {
	params := map[string]*yaml.Node{}
	for _, p := range pairs(with) {
		if _, ok := s.Params[p[0].Value]; !ok {
			return nil, fmt.Errorf("%s has no parameter %q", s.ID, p[0].Value)
		}
		params[p[0].Value] = p[1]
	}
	for name := range s.Params {
		if _, ok := params[name]; !ok {
			return nil, fmt.Errorf("parameter %q not given (with:)", name)
		}
	}
	n := substitute(&s.Evaluation, params)
	var e Evaluation
	if err := n.Decode(&e); err != nil {
		return nil, err
	}
	e.Line = s.Evaluation.Line
	return &e, nil
}

// substitute returns a copy of n with scalars "$name" replaced by the parameter's node.
func substitute(n *yaml.Node, params map[string]*yaml.Node) *yaml.Node {
	n = deref(n)
	if n.Kind == yaml.ScalarNode && strings.HasPrefix(n.Value, "$") {
		if p, ok := params[n.Value[1:]]; ok {
			return p
		}
	}
	cp := *n
	cp.Content = make([]*yaml.Node, len(n.Content))
	for i, c := range n.Content {
		cp.Content[i] = substitute(c, params)
	}
	return &cp
}

// Interpretations returns every interpretation with the file it lives in.
func (cat *Catalogue) Interpretations() []struct {
	File string
	I    Interpretation
} {
	var out []struct {
		File string
		I    Interpretation
	}
	add := func(f string, l []Interpretation) {
		for _, i := range l {
			out = append(out, struct {
				File string
				I    Interpretation
			}{f, i})
		}
	}
	for _, c := range cat.Credentials {
		add(c.File, c.Interpretations)
	}
	ids := make([]string, 0, len(cat.Shared))
	for id := range cat.Shared {
		ids = append(ids, id)
	}
	sort.Strings(ids)
	for _, id := range ids {
		add(cat.Shared[id].File, cat.Shared[id].Interpretations)
	}
	return out
}
