package credentials

import (
	"fmt"
	"slices"
	"sort"
	"strings"

	"gopkg.in/yaml.v3"

	"github.com/fjaeckel/open-aircrew-rules/engine"
)

// compiler turns one evaluation into an engine rule.
type compiler struct {
	v        *Vocab
	cat      *Catalogue
	where    string
	errs     []string
	usesWith bool
	pool     []string
	group    string
	inSum    bool
}

func (cp *compiler) fail(n *yaml.Node, format string, a ...any) {
	line := 0
	if n != nil {
		line = n.Line
	}
	cp.errs = append(cp.errs, fmt.Sprintf("%s (line %d): %s", cp.where, line, fmt.Sprintf(format, a...)))
}

func (cp *compiler) ref(n *yaml.Node, what string) {
	if n == nil || n.Kind == 0 {
		cp.fail(n, "%s has no ref", what)
		return
	}
	var l RefList
	if err := n.Decode(&l); err != nil || len(l) == 0 {
		cp.fail(n, "%s: ref must be a string or a list of strings", what)
		return
	}
	for _, r := range l {
		cp.cat.addRef(r, cp.where+" "+what, n.Line)
	}
	cp.delegated(l, what, n.Line)
}

// delegated requires the statute delegating an association document next to its ref.
func (cp *compiler) delegated(l RefList, what string, line int) {
	for _, r := range l {
		id, ok := strings.CutPrefix(r, AssociationPrefix+":")
		if !ok {
			continue
		}
		if a, found := cp.cat.Associations[id]; found && !slices.Contains(l, a.DelegatedBy) {
			cp.errs = append(cp.errs, fmt.Sprintf("%s (line %d): %s cites %s without the statute delegating it (%s)", cp.where, line, what, r, a.DelegatedBy))
		}
	}
}

func (cp *compiler) refList(l RefList, what string, line int, required bool) {
	if len(l) == 0 {
		if required {
			cp.errs = append(cp.errs, fmt.Sprintf("%s (line %d): %s has no ref", cp.where, line, what))
		}
		return
	}
	for _, r := range l {
		cp.cat.addRef(r, cp.where+" "+what, line)
	}
	cp.delegated(l, what, line)
}

// pairs returns the key/value pairs of a mapping node.
func pairs(n *yaml.Node) [][2]*yaml.Node {
	n = deref(n)
	if n == nil || n.Kind != yaml.MappingNode {
		return nil
	}
	out := make([][2]*yaml.Node, 0, len(n.Content)/2)
	for i := 0; i+1 < len(n.Content); i += 2 {
		out = append(out, [2]*yaml.Node{n.Content[i], deref(n.Content[i+1])})
	}
	return out
}

func deref(n *yaml.Node) *yaml.Node {
	for n != nil && n.Kind == yaml.AliasNode {
		n = n.Alias
	}
	return n
}

func empty(n *yaml.Node) bool { return n == nil || n.Kind == 0 }

// compileEvaluation builds the engine rule of one resolved evaluation.
func (cat *Catalogue) compileEvaluation(c *Credential, e *Evaluation, id string) (*engine.Rule, *compiler) {
	cp := &compiler{v: cat.Vocab, cat: cat, where: fmt.Sprintf("%s#%s", c.ID, e.ID)}
	if e.Shared != nil {
		cp.where = fmt.Sprintf("%s#%s (uses %s)", c.ID, e.ID, e.Shared.ID)
	}
	r := &engine.Rule{
		ID:                 id,
		Title:              e.Asks,
		Authority:          strings.ToLower(c.Authority),
		RuleDescriptionKey: e.DescriptionKey,
		EffectiveFrom:      e.EffectiveFrom,
		EffectiveTo:        e.EffectiveTo,
		File:               c.File,
	}
	if e.EffectiveFrom != nil && e.EffectiveTo != nil && e.EffectiveTo.Before(*e.EffectiveFrom) {
		cp.errs = append(cp.errs, cp.where+": effective_to is before effective_from")
	}
	if e.Source == "" {
		cp.errs = append(cp.errs, cp.where+": evaluation has no source")
	} else {
		cp.cat.addRef(e.Source, cp.where+" source", e.Line)
		if pr, err := ParseRef(e.Source); err == nil {
			r.Citations = append(r.Citations, cat.Vocab.Cite(pr))
		}
	}
	for _, s := range e.AlsoCites {
		cp.cat.addRef(s, cp.where+" also_cites", e.Line)
		if pr, err := ParseRef(s); err == nil {
			r.Citations = append(r.Citations, cat.Vocab.Cite(pr))
		}
	}
	if e.RelevantClass != nil {
		cp.pool = e.RelevantClass.PooledWithHeld
		cp.group = e.RelevantClass.ClassGroup
		if (len(cp.pool) > 0) == (cp.group != "") {
			cp.errs = append(cp.errs, cp.where+": relevant_class names pooled_with_held or class_group")
		}
		if _, ok := cat.ev.ClassGroups[cp.group]; cp.group != "" && !ok {
			cp.errs = append(cp.errs, fmt.Sprintf("%s: relevant_class.class_group %q is not in the vocabulary", cp.where, cp.group))
		}
		cp.refList(e.RelevantClass.Ref, "relevant_class", e.Line, true)
	}
	r.AppliesTo = cp.appliesTo(c, e)
	if cp.group != "" && r.AppliesTo.Subject == "passengers" {
		r.AppliesTo.ClassGroup = cp.group
	}
	if !empty(&e.Counting) {
		w, f := cp.qualifiers(&e.Counting, pairs(&e.Counting), true)
		r.Window, r.Filter = w, f
	}
	if !empty(&e.PassesIf) {
		r.Requirements = cp.node(&e.PassesIf, "passes_if")
	}
	if !empty(&e.RestoredBy) {
		r.RestoredBy = cp.restoredBy(&e.RestoredBy)
	}
	if e.ValidFor != nil {
		r.Validity = cp.validity(e.ValidFor, e.Line)
		if conv, ok := cat.Vocab.AuthorityConventions[c.Authority]; ok && conv.ValidityEnds == "day_before" {
			r.Validity.EndOffsetDays = -1
		}
	}
	if e.OnFail != nil {
		cp.refList(e.OnFail.Ref, "on_fail", e.Line, true)
	}
	r.Stages = cp.outcomes(&e.Outcomes, r)
	return r, cp
}

// selfPart returns the record part an evaluation about the credential itself evaluates.
func selfPart(c *Credential) (string, *Part) {
	s := c.Selects
	if c.Kind == "licence" && s.Licence != nil {
		return "licence", s.Licence
	}
	switch {
	case s.Ratings != nil:
		return "ratings", s.Ratings
	case s.Privilege != nil:
		return "privilege", s.Privilege
	case s.Credential != nil:
		return "credential", s.Credential
	case s.Licence != nil:
		return "licence", s.Licence
	}
	return "", nil
}

func (s Selects) part(name string) *Part {
	switch name {
	case "licence":
		return s.Licence
	case "ratings":
		return s.Ratings
	case "privilege":
		return s.Privilege
	case "credential":
		return s.Credential
	}
	return nil
}

func (cp *compiler) appliesTo(c *Credential, e *Evaluation) engine.AppliesTo {
	about := e.About
	if about == "" {
		about = "self"
	}
	var a engine.AppliesTo
	var partName string
	var part *Part
	if about == "self" {
		partName, part = selfPart(c)
		a.Subject = cp.v.Selects[partName]
	} else {
		def, ok := cp.v.About[about]
		if !ok {
			cp.errs = append(cp.errs, fmt.Sprintf("%s: about %q is not in the vocabulary", cp.where, about))
			return a
		}
		a.Subject = def.Subject
		partName, part = def.Selects, c.Selects.part(def.Selects)
	}
	switch {
	case e.Scope != nil:
		cp.applyScope(&a, e.Scope)
		cp.refList(e.Scope.Ref, "scope", e.Line, false)
	case part == nil:
		cp.errs = append(cp.errs, fmt.Sprintf("%s: the credential selects no %s part for an evaluation about %s; give the evaluation a scope", cp.where, partName, about))
	default:
		a.Authorities, a.ExcludeAuthorities = part.Authorities, part.NotAuthorities
		a.LicenceKinds, a.ExcludeLicenceKinds = part.LicenceKinds, part.NotLicenceKinds
		switch partName {
		case "licence":
			a.LicenceKinds = append(slices.Clone(part.Kinds), part.LicenceKinds...)
		case "ratings":
			a.Classes = part.Classes
			a.ULKinds = part.ULKinds
		case "privilege":
			a.PrivilegeKinds = part.Kinds
		case "credential":
			a.CredentialTypes = part.Kinds
		}
	}
	if e.OnlyFor != nil {
		cp.applyScope(&a, e.OnlyFor)
		cp.refList(e.OnlyFor.Ref, "only_for", e.Line, false)
	}
	return a
}

func (cp *compiler) applyScope(a *engine.AppliesTo, s *Scope) {
	switch s.IfMissing {
	case "":
	case "unknown":
		for criterion, given := range map[string]bool{
			"typeRated": s.TypeRated != nil && *s.TypeRated, "ulKinds": len(s.ULKinds) > 0,
			"categories": len(s.Categories) > 0, "differentEngineType": s.DifferentEngineType != nil,
		} {
			if given {
				a.UnknownWhenMissing = append(a.UnknownWhenMissing, criterion)
			}
		}
		sort.Strings(a.UnknownWhenMissing)
		if len(a.UnknownWhenMissing) == 0 {
			cp.errs = append(cp.errs, cp.where+": if_missing needs type_rated: true, ul_kinds, categories or different_engine_type beside it")
		}
	default:
		cp.errs = append(cp.errs, fmt.Sprintf("%s: if_missing is unknown, not %q", cp.where, s.IfMissing))
	}
	applyScope(a, s)
}

func applyScope(a *engine.AppliesTo, s *Scope) {
	set := func(dst *[]string, src []string) {
		if len(src) > 0 {
			*dst = src
		}
	}
	set(&a.Authorities, s.Authorities)
	set(&a.ExcludeAuthorities, s.NotAuthorities)
	set(&a.LicenceKinds, s.LicenceKinds)
	set(&a.ExcludeLicenceKinds, s.NotLicenceKinds)
	set(&a.Classes, s.Classes)
	set(&a.ExcludeClasses, s.NotClasses)
	set(&a.CredentialTypes, s.CredentialTypes)
	set(&a.PrivilegeKinds, s.PrivilegeKinds)
	set(&a.LaunchMethods, s.LaunchMethods)
	set(&a.ULKinds, s.ULKinds)
	if s.TypeRated != nil {
		a.TypeRated = s.TypeRated
	}
	set(&a.Categories, s.Categories)
	if s.DifferentEngineType != nil {
		a.DifferentEngineType = s.DifferentEngineType
	}
	if s.Programme != "" {
		a.Programme = s.Programme
	}
	if s.WhenHolding != nil {
		a.Holds = s.WhenHolding
	}
}

var nodeKeys = []string{"id", "ref", "only_if", "informational", "all_of", "any_of", "n_of", "sum_of", "min", "min_hours", "min_minutes", "name", "unit", "remedy"}

// node compiles a combinator or a single-count item.
func (cp *compiler) node(n *yaml.Node, path string) *engine.Node {
	n = deref(n)
	ps := pairs(n)
	if len(ps) == 0 {
		cp.fail(n, "%s: want a mapping", path)
		return nil
	}
	var comb string
	for _, p := range ps {
		if k := p[0].Value; k == "all_of" || k == "any_of" || k == "n_of" || k == "sum_of" {
			if comb != "" {
				cp.fail(p[0], "%s: more than one combinator", path)
			}
			comb = k
		}
	}
	if comb == "" {
		if len(ps) != 1 {
			cp.fail(n, "%s: an item is one count word ({ %s: {...} }); found %d keys", path, ps[0][0].Value, len(ps))
			return nil
		}
		return cp.leaf(ps[0][0], ps[0][1], path)
	}
	out := &engine.Node{}
	var quals [][2]*yaml.Node
	var hasRef bool
	min := -1.0
	for _, p := range ps {
		k, v := p[0].Value, p[1]
		if comb != "sum_of" && slices.Contains([]string{"min", "min_hours", "min_minutes", "name", "unit", "remedy"}, k) {
			cp.fail(p[0], "%s: %s belongs to a count or a sum_of", path, k)
			continue
		}
		switch k {
		case "min", "min_minutes":
			_ = v.Decode(&min)
		case "min_hours":
			var h float64
			_ = v.Decode(&h)
			min = h * 60
		case "name":
			out.NameKey = v.Value
		case "unit":
			out.Unit = v.Value
		case "remedy":
			out.RemedyKey = v.Value
		case "sum_of":
			for i, c := range v.Content {
				if kid := cp.sumItem(c, fmt.Sprintf("%s.sum_of[%d]", path, i)); kid != nil {
					out.SumOf = append(out.SumOf, kid)
				}
			}
			if len(out.SumOf) == 0 {
				out.SumOf = []*engine.Node{}
			}
		case "id":
			out.ID = v.Value
		case "ref":
			hasRef = true
			cp.ref(v, path+" "+comb)
		case "only_if":
			out.When = cp.condition(v, path+" only_if")
		case "informational":
			out.Informational = v.Value == "true"
		case "all_of", "any_of":
			var kids []*engine.Node
			for i, c := range v.Content {
				if kid := cp.node(c, fmt.Sprintf("%s.%s[%d]", path, k, i)); kid != nil {
					kids = append(kids, kid)
				}
			}
			if k == "all_of" {
				out.AllOf = kids
			} else {
				out.AnyOf = kids
			}
		case "n_of":
			nof := &engine.NOf{}
			for _, q := range pairs(v) {
				switch q[0].Value {
				case "n":
					_ = q[1].Decode(&nof.N)
				case "of":
					for i, c := range q[1].Content {
						if kid := cp.node(c, fmt.Sprintf("%s.n_of[%d]", path, i)); kid != nil {
							nof.Of = append(nof.Of, kid)
						}
					}
				}
			}
			out.NOf = nof
		default:
			quals = append(quals, p)
		}
	}
	if !hasRef {
		cp.fail(n, "%s: %s has no ref", path, comb)
	}
	if comb == "sum_of" {
		cp.sumDefaults(n, out, min, path)
	}
	out.Window, out.Filter = cp.qualifiers(n, quals, true)
	return out
}

// sumItem compiles one count of a sum_of: a count word without a minimum.
func (cp *compiler) sumItem(n *yaml.Node, path string) *engine.Node {
	ps := pairs(n)
	if len(ps) != 1 || cp.v.Counts[ps[0][0].Value].Metric == "" {
		cp.fail(n, "%s: a sum_of item is one count word", path)
		return nil
	}
	for _, p := range pairs(ps[0][1]) {
		if k := p[0].Value; k == "min" || k == "min_hours" || k == "min_minutes" || k == "waived_by" {
			cp.fail(p[0], "%s: a sum_of item has no %s; the sum_of has the minimum", path, k)
		}
	}
	cp.inSum = true
	defer func() { cp.inSum = false }()
	return cp.leaf(ps[0][0], ps[0][1], path)
}

// sumDefaults checks a sum_of's id and minimum and takes unit, name and remedy from its
// first item when not given.
func (cp *compiler) sumDefaults(n *yaml.Node, out *engine.Node, min float64, path string) {
	if out.ID == "" {
		cp.fail(n, "%s: a sum_of needs an id (it is a requirement row)", path)
	}
	if min < 0 {
		cp.fail(n, "%s: a sum_of needs min, min_hours or min_minutes", path)
		min = 0
	}
	out.Min = &min
	if len(out.SumOf) == 0 {
		cp.fail(n, "%s: a sum_of adds at least one count", path)
		return
	}
	first := out.SumOf[0]
	for _, c := range out.SumOf[1:] {
		if c.Unit != first.Unit {
			cp.fail(n, "%s: sum_of items have one unit (%s, %s)", path, first.Unit, c.Unit)
		}
	}
	if out.Unit == "" {
		out.Unit = first.Unit
	}
	if out.NameKey == "" {
		out.NameKey = first.NameKey
	}
	switch out.RemedyKey {
	case "":
		out.RemedyKey = first.RemedyKey
	case "none":
		out.RemedyKey = ""
	}
}

var leafKeys = []string{"min", "min_hours", "min_minutes", "max", "max_hours", "max_minutes", "unknown_if_none", "id", "name", "unit", "remedy", "messages", "informational", "only_if", "credit", "ref", "waived_by"}

// leaf compiles one count word with its parameters.
func (cp *compiler) leaf(key, val *yaml.Node, path string) *engine.Node {
	word := key.Value
	def, ok := cp.v.Counts[word]
	if !ok {
		cp.fail(key, "%s: %q is not a count word or combinator in the vocabulary", path, word)
		return nil
	}
	out := &engine.Node{ID: def.ID, Metric: def.Metric, Unit: def.Unit, NameKey: def.Name, RemedyKey: def.Remedy}
	path += "." + word
	var quals [][2]*yaml.Node
	min := -1.0
	var waiver *yaml.Node
	hasRef := false
	ps := pairs(val)
	if val.Kind == yaml.ScalarNode {
		cp.fail(val, "%s: write { min: %s, ref: ... }", path, val.Value)
	}
	for _, p := range ps {
		k, v := p[0].Value, p[1]
		switch k {
		case "min":
			_ = v.Decode(&min)
		case "min_hours":
			var h float64
			_ = v.Decode(&h)
			min = h * 60
		case "min_minutes":
			_ = v.Decode(&min)
		case "max", "max_minutes", "max_hours":
			var x float64
			if err := v.Decode(&x); err != nil {
				cp.fail(v, "%s.%s needs a number", path, k)
			}
			if k == "max_hours" {
				x *= 60
			}
			out.Max = &x
		case "unknown_if_none":
			out.UnknownIfNone = v.Value == "true"
		case "id":
			out.ID = v.Value
		case "name":
			out.NameKey = v.Value
		case "unit":
			out.Unit = v.Value
		case "remedy":
			out.RemedyKey = v.Value
			if v.Value == "none" {
				out.RemedyKey = ""
			}
		case "messages":
			out.Messages = &engine.ReqMessages{}
			_ = v.Decode(out.Messages)
		case "informational":
			out.Informational = v.Value == "true"
		case "only_if":
			out.When = cp.condition(v, path+" only_if")
		case "credit":
			out.EscapeHatch = v.Value
		case "ref":
			hasRef = true
			cp.ref(v, path)
		case "waived_by":
			waiver = v
		default:
			quals = append(quals, p)
		}
	}
	if !hasRef {
		cp.fail(key, "%s has no ref", path)
	}
	switch {
	case cp.inSum:
		min = 0
	case min < 0:
		if def.Amount == "time" {
			cp.fail(key, "%s: give min_hours or min_minutes", path)
		}
		min = 1
	}
	out.Min = &min
	w, f := cp.qualifiers(val, quals, true)
	f = mergeImplied(cp, def, f)
	out.Window, out.Filter = w, f
	if waiver == nil {
		return out
	}
	return cp.waived(word, out, waiver, path)
}

// mergeImplied adds the filters a count word implies (event kinds, withMinutes ...).
func mergeImplied(cp *compiler, def CountDef, f *engine.Filter) *engine.Filter {
	if len(def.Events) == 0 && len(def.Filter) == 0 {
		return f
	}
	m := map[string]any{}
	for k, v := range def.Filter {
		m[k] = v
	}
	if len(def.Events) > 0 {
		kinds := slices.Clone(def.Events)
		if f != nil && f.Has("eventKinds") {
			kinds = append(kinds, f.EventKinds...)
		}
		m["eventKinds"] = kinds
	}
	return engine.MergeFilter(f, cp.filter(m))
}

// waived wraps a leaf in any_of with the event-based exemption.
func (cp *compiler) waived(word string, leaf *engine.Node, w *yaml.Node, path string) *engine.Node {
	one := 1.0
	ex := &engine.Node{ID: word + "_exemption", Metric: "events", Unit: "check", NameKey: "requirement." + word + "_exemption", Min: &one}
	var quals [][2]*yaml.Node
	var events []string
	hasRef := false
	for _, p := range pairs(w) {
		switch p[0].Value {
		case "events":
			_ = p[1].Decode(&events)
		case "id":
			ex.ID = p[1].Value
		case "name":
			ex.NameKey = p[1].Value
		case "ref":
			hasRef = true
			cp.ref(p[1], path+" waived_by")
		default:
			quals = append(quals, p)
		}
	}
	if !hasRef {
		cp.fail(w, "%s waived_by has no ref", path)
	}
	if len(events) == 0 {
		cp.fail(w, "%s waived_by: list the events that waive it", path)
	}
	win, f := cp.qualifiers(w, quals, true)
	ex.Window = win
	ex.Filter = engine.MergeFilter(f, cp.filter(map[string]any{"eventKinds": events}))
	return &engine.Node{ID: word, AnyOf: []*engine.Node{leaf, ex}}
}

// qualifiers compiles qualifier words into a window and a filter.
func (cp *compiler) qualifiers(at *yaml.Node, ps [][2]*yaml.Node, allowWindow bool) (*engine.Window, *engine.Filter) {
	var win *engine.Window
	m := map[string]any{}
	for _, p := range ps {
		k, v := p[0].Value, p[1]
		def, ok := cp.v.Qualifiers[k]
		if !ok {
			cp.fail(p[0], "%q is not a qualifier in the vocabulary", k)
			continue
		}
		switch {
		case def.Document:
			cp.usesWith = true
		case def.Window != "":
			if !allowWindow {
				cp.fail(p[0], "%s: no window here", k)
				continue
			}
			if win != nil {
				cp.fail(p[0], "%s: a second window", k)
			}
			win = &engine.Window{Kind: def.Window, Anchor: def.Anchor}
			switch def.Window {
			case "lifetime", "validity_period", "since_issue":
				if v.Value != "true" {
					cp.fail(v, "%s: write %s: true", k, k)
				}
			default:
				if err := v.Decode(&win.N); err != nil {
					cp.fail(v, "%s needs a number", k)
				}
			}
		default:
			cp.filterQualifier(k, def, v, m)
		}
	}
	if len(m) == 0 {
		return win, nil
	}
	return win, cp.filter(m)
}

func (cp *compiler) filterQualifier(k string, def QualifierDef, v *yaml.Node, m map[string]any) {
	switch k {
	case "in_class":
		if v.Value != "true" {
			cp.fail(v, "in_class: write in_class: true")
			return
		}
		m["classes"] = []string{"$subject"}
		if len(cp.pool) > 0 {
			m["heldClassPools"] = [][]string{cp.pool}
		}
		if cp.group != "" {
			m["classGroup"] = cp.group
		}
	case "ul_credit":
		var credits []map[string]any
		for _, p := range pairs(v) {
			if p[0].Value == "ref" {
				cp.ref(p[1], "ul_credit")
				continue
			}
			credit := map[string]any{"class": p[0].Value}
			if p[1].Kind == yaml.MappingNode {
				var c struct {
					ULKinds     []string `yaml:"ul_kinds"`
					MinMTOMKg   *int     `yaml:"min_mtom_kg"`
					FixedEngine *bool    `yaml:"fixed_engine"`
				}
				err := p[1].Decode(&c)
				extra := 0
				if c.MinMTOMKg != nil {
					extra++
				}
				if c.FixedEngine != nil {
					extra++
				}
				if err != nil || len(c.ULKinds) == 0 || extra == 0 || len(pairs(p[1])) != 1+extra || (c.FixedEngine != nil && !*c.FixedEngine) {
					cp.fail(p[1], "ul_credit.%s: want { ul_kinds: [...], min_mtom_kg: n, fixed_engine: true } (min_mtom_kg or fixed_engine or both)", p[0].Value)
				}
				credit["ulKinds"] = c.ULKinds
				if c.MinMTOMKg != nil {
					credit["minMtomKg"] = *c.MinMTOMKg
				}
				if c.FixedEngine != nil && *c.FixedEngine {
					credit["fixedEngine"] = true
				}
			} else {
				var kinds []string
				if err := p[1].Decode(&kinds); err != nil {
					cp.fail(p[1], "ul_credit.%s: want a list of ultralight kinds", p[0].Value)
				}
				credit["ulKinds"] = kinds
			}
			credits = append(credits, credit)
		}
		if !slices.ContainsFunc(pairs(v), func(p [2]*yaml.Node) bool { return p[0].Value == "ref" }) {
			cp.fail(v, "ul_credit has no ref")
		}
		sort.Slice(credits, func(i, j int) bool { return credits[i]["class"].(string) < credits[j]["class"].(string) })
		m["ulCredit"] = credits
	case "as":
		var roles []string
		if v.Kind == yaml.ScalarNode {
			roles = []string{v.Value}
		} else {
			_ = v.Decode(&roles)
		}
		switch {
		case len(roles) == 1 && roles[0] == "pilot_flying":
			m["pilotFlying"] = true
		case len(roles) == 1 && roles[0] == "sole_manipulator":
			m["soleManipulator"] = true
		default:
			m["roles"] = roles
		}
	case "any_flight_of":
		var alts []any
		for _, c := range v.Content {
			_, f := cp.qualifiers(c, pairs(c), false)
			alts = append(alts, filterMap(f))
		}
		m["any"] = alts
	default:
		if def.Value == "subject" {
			if v.Value == "false" {
				return
			}
			if v.Value != "true" {
				cp.fail(v, "%s: write %s: true", k, k)
				return
			}
			m[def.Filter] = []string{"$subject"}
			return
		}
		var x any
		_ = v.Decode(&x)
		m[def.Filter] = x
	}
}

// filter turns a map with engine filter keys into an engine filter.
func (cp *compiler) filter(m map[string]any) *engine.Filter {
	b, err := yaml.Marshal(m)
	if err != nil {
		cp.errs = append(cp.errs, cp.where+": "+err.Error())
		return nil
	}
	var f engine.Filter
	if err := yaml.Unmarshal(b, &f); err != nil {
		cp.errs = append(cp.errs, cp.where+": "+err.Error())
		return nil
	}
	return &f
}

// filterMap returns the given keys of a filter as a map (for nesting under any).
func filterMap(f *engine.Filter) map[string]any {
	out := map[string]any{}
	if f == nil {
		return out
	}
	b, _ := yaml.Marshal(f)
	var all map[string]any
	_ = yaml.Unmarshal(b, &all)
	for _, k := range f.Keys() {
		out[k] = all[k]
	}
	return out
}

// condition decodes a stage or only_if condition; a `ref` key beside it is recorded.
func (cp *compiler) condition(n *yaml.Node, what string) *engine.Condition {
	n = deref(n)
	var rest []*yaml.Node
	hasRef := false
	if n.Kind == yaml.MappingNode {
		for _, p := range pairs(n) {
			if p[0].Value == "ref" {
				hasRef = true
				cp.ref(p[1], what)
				continue
			}
			rest = append(rest, p[0], p[1])
		}
		n = &yaml.Node{Kind: yaml.MappingNode, Tag: "!!map", Content: rest, Line: n.Line}
	}
	if !hasRef {
		cp.fail(n, "%s has no ref", what)
	}
	c := &engine.Condition{}
	if err := n.Decode(c); err != nil {
		cp.fail(n, "%s: %v", what, err)
	}
	return c
}

func (cp *compiler) restoredBy(n *yaml.Node) []engine.EventHook {
	var out []engine.EventHook
	for _, item := range deref(n).Content {
		item = deref(item)
		h := engine.EventHook{}
		hasRef := false
		for _, p := range pairs(item) {
			if p[0].Value == "ref" {
				hasRef = true
				cp.ref(p[1], "restored_by")
				continue
			}
			if h.Event != "" {
				cp.fail(p[0], "restored_by: one event per item")
			}
			h.Event = p[0].Value
			_, h.Filter = cp.qualifiers(p[1], pairs(p[1]), false)
		}
		if !hasRef {
			cp.fail(item, "restored_by %s has no ref", h.Event)
		}
		out = append(out, h)
	}
	return out
}

func (cp *compiler) validity(v *ValidFor, line int) *engine.Validity {
	out := &engine.Validity{From: "issued", RecordedWins: v.RecordedExpiryWins, RecordedMin: v.RecordedExpiryMinimum}
	if v.RecordedExpiryWins && v.RecordedExpiryMinimum {
		cp.errs = append(cp.errs, fmt.Sprintf("%s: valid_for has both recorded_expiry_wins and recorded_expiry_minimum", cp.where))
	}
	switch v.CountedFrom {
	case "", "issue":
	case "valid_from":
		out.From = "valid_from"
	default:
		cp.errs = append(cp.errs, fmt.Sprintf("%s: valid_for.counted_from is issue or valid_from", cp.where))
	}
	switch v.AgeOn {
	case "":
	case "issue":
		out.AgeOn = "issued"
	case "valid_from":
		out.AgeOn = "valid_from"
	default:
		cp.errs = append(cp.errs, fmt.Sprintf("%s: valid_for.age_on is issue or valid_from", cp.where))
	}
	cp.refList(v.Ref, "valid_for", line, false)
	for i, p := range v.Periods {
		cp.refList(p.Ref, fmt.Sprintf("valid_for.periods[%d]", i), line, true)
		per := engine.Period{Months: p.Months, CapAtAge: p.EndsAtAge, EndOfMonth: p.EndOfMonth}
		if p.AgeUnder != nil || p.AgeFrom != nil {
			per.When = map[string]int{}
			if p.AgeUnder != nil {
				per.When["age_under"] = *p.AgeUnder
			}
			if p.AgeFrom != nil {
				per.When["age_at_least"] = *p.AgeFrom
			}
		}
		out.Periods = append(out.Periods, per)
	}
	return out
}
