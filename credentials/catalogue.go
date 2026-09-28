package credentials

import (
	"errors"
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
	// Policies are the declared policies (policies.yaml, plus fragments with Options.Fragments).
	Policies map[string]*Policy
	// Associations are the association documents of associations.yaml.
	Associations map[string]*Association
	Options      Options
	byID         map[string]*Credential
	// held is what each credential selects in a record (its own part).
	held map[string]engine.AppliesTo
	// resolved memoises resolve; following holds the uses nodes being expanded.
	resolved  map[*Evaluation]resolution
	following map[string]bool
}

// Options changes what Load reads.
type Options struct {
	// Fragments also reads fragments/keys/*.yaml and fragments/policies/*.yaml, and makes
	// CheckCoverage merge fragments/coverage/*.yaml (DESIGN.md section 12).
	Fragments bool
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

// Load reads vocabulary.yaml, messages/keys.yaml, policies.yaml and credentials/** under
// root and compiles every evaluation.
func Load(root string) (*Catalogue, error) { return LoadWith(root, Options{}) }

// LoadWith is Load with options.
func LoadWith(root string, o Options) (*Catalogue, error) {
	ev, err := engine.LoadVocabulary(filepath.Join(root, "vocabulary.yaml"))
	if err != nil {
		return nil, err
	}
	var moreKeys []string
	if o.Fragments {
		if moreKeys, err = fragmentFiles(root, "keys"); err != nil {
			return nil, err
		}
	}
	keys, err := engine.LoadKeys(filepath.Join(root, "messages", "keys.yaml"), moreKeys...)
	if err != nil {
		return nil, err
	}
	v, err := LoadVocab(filepath.Join(root, "vocabulary.yaml"))
	if err != nil {
		return nil, err
	}
	pol, err := LoadPolicies(root, o.Fragments)
	if err != nil {
		return nil, err
	}
	cat := &Catalogue{Root: root, Vocab: v, Policies: pol, Options: o, Shared: map[string]*Shared{}, Compiled: map[string]*Compiled{}, ByEvaluation: map[string]string{}, byID: map[string]*Credential{}, held: map[string]engine.AppliesTo{}, resolved: map[*Evaluation]resolution{}, following: map[string]bool{}}
	assoc, err := LoadAssociations(root)
	if err != nil {
		return nil, err
	}
	cat.Associations = assoc
	cat.Resolver = NewResolver(root, v, pol)
	cat.Resolver.associations = assoc
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
	cat.checkUses()
	var rules []*engine.Rule
	for _, c := range cat.Credentials {
		for _, e := range c.Evaluations {
			key := c.ID + "#" + e.ID
			re, err := cat.resolve(c, e)
			if errors.Is(err, errReported) {
				continue
			}
			if err != nil {
				cat.Errors = append(cat.Errors, fmt.Sprintf("%s: %v", key, err))
				continue
			}
			if e.isRequirement() {
				cat.ByEvaluation[key] = ""
				cat.checkRequirement(key, e)
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
	cat.checkRequirementCycles()
	for _, c := range cat.Credentials {
		cp := &compiler{v: v, cat: cat, where: c.ID}
		cat.held[c.ID] = cp.appliesTo(c, &Evaluation{})
		cat.Errors = append(cat.Errors, cp.errs...)
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

// resolve returns the evaluation with `uses` expanded, memoised per evaluation.
func (cat *Catalogue) resolve(c *Credential, e *Evaluation) (*Evaluation, error) {
	if e.Uses == "" {
		return e, nil
	}
	if r, ok := cat.resolved[e]; ok {
		return r.eval, r.err
	}
	r, err := cat.follow(c, e, c.ID+"#"+e.ID)
	cat.resolved[e] = resolution{r, err}
	return r, err
}

// resolution is a memoised resolve result.
type resolution struct {
	eval *Evaluation
	err  error
}

// follow expands the `uses` of e, which is the evaluation node named node; a node met again
// on the way is a cycle, which checkUses reports.
func (cat *Catalogue) follow(c *Credential, e *Evaluation, node string) (*Evaluation, error) {
	if cat.following[node] {
		return nil, errReported
	}
	cat.following[node] = true
	defer delete(cat.following, node)
	var base *Evaluation
	if credID, evalID, ok := strings.Cut(e.Uses, "#"); ok {
		oc, found := cat.byID[credID]
		if !found {
			return nil, fmt.Errorf("uses %s: no credential %s", e.Uses, credID)
		}
		for _, x := range oc.Evaluations {
			if x.ID != evalID {
				continue
			}
			if x.isRequirement() {
				return nil, fmt.Errorf("uses %s: it only names required credentials; name them in requires_all or requires_any instead", e.Uses)
			}
			b, err := cat.resolve(oc, x)
			if err != nil {
				return nil, err
			}
			cp := *b
			base = &cp
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
		if b.Uses != "" {
			if strings.HasPrefix(b.Uses, "$") || strings.Contains(b.Uses, "#") {
				return nil, errReported
			}
			if b, err = cat.follow(c, b, s.ID); err != nil {
				return nil, err
			}
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
