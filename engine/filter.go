package engine

import (
	"slices"
	"strings"
)

// match is a tri-state filter result.
type match int

const (
	matchNo match = iota
	matchYes
	matchUnknown
)

// matcher accumulates filter checks; any no wins, then any unknown.
type matcher struct {
	res     match
	unknown []string
}

func newMatcher() *matcher { return &matcher{res: matchYes} }

func (m *matcher) no() { m.res = matchNo }

func (m *matcher) check(ok bool) {
	if !ok {
		m.res = matchNo
	}
}

func (m *matcher) unknownBy(filter string) {
	if m.res == matchYes {
		m.res = matchUnknown
	}
	if !slices.Contains(m.unknown, filter) {
		m.unknown = append(m.unknown, filter)
	}
}

// resolvedFilter is a merged filter with $subject replaced.
type resolvedFilter struct {
	*Filter
	classes map[string]bool
	credit  map[string]ULCredit
	anyOf   []*resolvedFilter
	// unresolved names the filters whose $subject the subject cannot supply: every item is
	// then unknown for them, never silently excluded.
	unresolved []string
}

const subjectToken = "$subject"

// resolveList replaces $subject in l by vals; ok is false when l has $subject and vals is empty.
func resolveList(l []string, vals ...string) ([]string, bool) {
	out := make([]string, 0, len(l))
	ok := true
	for _, v := range l {
		if v != subjectToken {
			out = append(out, v)
			continue
		}
		n := len(out)
		for _, x := range vals {
			if x != "" {
				out = append(out, x)
			}
		}
		ok = ok && len(out) > n
	}
	return out, ok
}

// detailTokens returns the words of a privilege detail that belong to domain, spelled as
// the domain spells them ("THREE_AXIS, weight_shift" -> THREE_AXIS, WEIGHT_SHIFT).
func detailTokens(detail string, domain []string) []string {
	var out []string
	for _, w := range strings.FieldsFunc(detail, func(r rune) bool { return strings.ContainsRune(",;/+ \t", r) }) {
		for _, d := range domain {
			if strings.EqualFold(w, d) && !slices.Contains(out, d) {
				out = append(out, d)
			}
		}
	}
	return out
}

// resolveFilter replaces $subject and expands class pools and groups.
func (e *evalCtx) resolveFilter(f *Filter) *resolvedFilter {
	if f == nil {
		f = &Filter{}
	}
	c := *f
	r := &resolvedFilter{Filter: &c}
	resolve := func(name string, l []string, vals ...string) []string {
		out, ok := resolveList(l, vals...)
		if !ok {
			r.unresolved = append(r.unresolved, name)
		}
		return out
	}
	c.Classes = resolve("classes", f.Classes, e.subj.Class)
	ul := []string{e.subj.ULKind}
	if e.subj.ULKind == "" {
		ul = detailTokens(e.subj.Detail, e.v.ULKinds)
	}
	c.ULKinds = resolve("ulKinds", f.ULKinds, ul...)
	c.LaunchMethods = resolve("launchMethods", f.LaunchMethods, e.subj.Detail)
	c.TypeDesignators = resolve("typeDesignators", f.TypeDesignators, e.subj.typeDesignator)
	c.Variants = resolve("variants", f.Variants, e.subj.variant)
	c.Categories = resolve("categories", f.Categories, e.subj.category)
	c.TowKinds = resolve("towKinds", f.TowKinds, detailTokens(e.subj.Detail, e.v.TowKinds)...)
	c.TowTakeUps = resolve("towTakeUps", f.TowTakeUps, detailTokens(e.subj.Detail, e.v.TowTakeUps)...)
	if f.Has("classes") {
		r.classes = map[string]bool{}
		for _, cl := range c.Classes {
			r.classes[cl] = true
			if f.ClassGroup != "" {
				_, members := e.v.ClassGroup(f.ClassGroup, cl)
				for _, m := range members {
					r.classes[m] = true
				}
			}
		}
		held := e.heldClassesOnLicence()
		for _, pool := range f.HeldClassPools {
			hit := false
			for _, cl := range pool {
				hit = hit || r.classes[cl]
			}
			if !hit {
				continue
			}
			for _, cl := range pool {
				if held[cl] {
					r.classes[cl] = true
				}
			}
		}
		r.credit = map[string]ULCredit{}
		for _, uc := range f.ULCredit {
			if r.classes[uc.Class] {
				r.credit[uc.Class] = uc
			}
		}
	}
	for _, a := range f.Any {
		r.anyOf = append(r.anyOf, e.resolveFilter(a))
	}
	return r
}

// heldClassesOnLicence returns the rating classes on the subject's licence.
func (e *evalCtx) heldClassesOnLicence() map[string]bool {
	out := map[string]bool{}
	if e.subj.licence == nil {
		return out
	}
	for _, r := range e.p.ratings {
		if r.LicenceID == e.subj.licence.ID {
			out[r.Class] = true
		}
	}
	return out
}

// matchFlight applies a resolved filter to a flight.
func (e *evalCtx) matchFlight(f *resolvedFilter, fl *Flight) *matcher {
	m := newMatcher()
	for _, u := range f.unresolved {
		m.unknownBy(u)
	}
	switch f.Simulator {
	case "only":
		m.check(fl.IsSimulator)
	case "include":
	default:
		m.check(!fl.IsSimulator)
	}
	if f.classes != nil {
		e.matchFlightClass(f, fl, m)
	}
	if len(f.ExcludeClasses) > 0 {
		m.check(!slices.Contains(f.ExcludeClasses, fl.Class))
	}
	if len(f.Categories) > 0 {
		m.check(slices.Contains(f.Categories, e.v.Category(fl.Class)))
	}
	if f.Has("ulKinds") {
		switch {
		case fl.Class != "ULTRALIGHT":
			m.no()
		case fl.ULKind == "":
			m.unknownBy("ulKinds")
		default:
			m.check(slices.Contains(f.ULKinds, fl.ULKind))
		}
	}
	if f.Has("launchMethods") {
		m.check(fl.LaunchMethod != "" && slices.Contains(f.LaunchMethods, fl.LaunchMethod))
	}
	stringFilter(m, f.Has("typeDesignators"), "typeDesignators", f.TypeDesignators, fl.TypeDesignator)
	stringFilter(m, f.Has("variants"), "variants", f.Variants, fl.Variant)
	if len(f.Roles) > 0 {
		m.check(slices.ContainsFunc(f.Roles, func(r string) bool { return roleMinutes(fl.Minutes, r) > 0 }))
	}
	for _, name := range f.WithMinutes {
		v, _ := fl.Minutes.Minute(name)
		m.check(v > 0)
	}
	for _, name := range f.WithoutMinutes {
		v, _ := fl.Minutes.Minute(name)
		m.check(v == 0)
	}
	if len(f.FSTDTypes) > 0 && fl.IsSimulator {
		if fl.FSTDType == "" {
			m.unknownBy("fstdTypes")
		} else {
			m.check(slices.Contains(f.FSTDTypes, fl.FSTDType))
		}
	}
	for name, want := range f.Flags {
		got, _ := fl.Flags.Get(name)
		m.check(got == want)
	}
	if f.MinLandings != nil {
		m.check(fl.Landings.Total() >= *f.MinLandings)
	}
	if f.MinDistanceKm != nil {
		if fl.DistanceKm == nil {
			m.unknownBy("minDistanceKm")
		} else {
			m.check(*fl.DistanceKm >= *f.MinDistanceKm)
		}
	}
	maxFilter(m, "maxEngines", f.MaxEngines, fl.Engines)
	maxFilter(m, "maxMtomKg", f.MaxMTOMKg, fl.MTOMKg)
	boolFilter(m, "tailwheel", f.Tailwheel, fl.Tailwheel)
	boolFilter(m, "soleManipulator", f.SoleManipulator, fl.SoleManipulator)
	boolFilter(m, "pilotFlying", f.PilotFlying, fl.PilotFlying)
	if len(f.TowKinds) > 0 {
		switch {
		case !fl.Flags.TowFlight:
			m.no()
		case fl.TowKind == "":
			m.unknownBy("towKinds")
		default:
			m.check(slices.Contains(f.TowKinds, fl.TowKind))
		}
	}
	if len(f.TowTakeUps) > 0 {
		switch {
		case !fl.Flags.TowFlight:
			m.no()
		case fl.TowTakeUp == "":
			m.unknownBy("towTakeUps")
		default:
			m.check(slices.Contains(f.TowTakeUps, fl.TowTakeUp))
		}
	}
	if len(f.anyOf) > 0 {
		e.matchAny(m, f.anyOf, func(sub *resolvedFilter) *matcher { return e.matchFlight(sub, fl) })
	}
	return m
}

func (e *evalCtx) matchFlightClass(f *resolvedFilter, fl *Flight, m *matcher) {
	if f.classes[fl.Class] {
		return
	}
	if fl.Class != "ULTRALIGHT" || len(f.credit) == 0 {
		m.no()
		return
	}
	best := matchNo
	var unknownBy string
	for _, uc := range f.credit {
		switch {
		case fl.ULKind == "":
			best, unknownBy = matchUnknown, "ulCredit"
		case !slices.Contains(uc.ULKinds, fl.ULKind):
		case uc.MinMTOMKg != nil && fl.MTOMKg == nil:
			if best != matchYes {
				best, unknownBy = matchUnknown, "ulCredit"
			}
		case uc.MinMTOMKg != nil && *fl.MTOMKg < *uc.MinMTOMKg:
		case uc.FixedEngine && fl.FixedEngine == nil:
			if best != matchYes {
				best, unknownBy = matchUnknown, "ulCredit"
			}
		case uc.FixedEngine && !*fl.FixedEngine:
		default:
			best = matchYes
		}
		if best == matchYes {
			return
		}
	}
	if best == matchUnknown {
		m.unknownBy(unknownBy)
		return
	}
	m.no()
}

func (e *evalCtx) matchAny(m *matcher, subs []*resolvedFilter, fn func(*resolvedFilter) *matcher) {
	best := matchNo
	var unknown []string
	for _, s := range subs {
		r := fn(s)
		if r.res == matchYes {
			return
		}
		if r.res == matchUnknown {
			best = matchUnknown
			unknown = append(unknown, r.unknown...)
		}
	}
	if best == matchNo {
		m.no()
		return
	}
	for _, u := range unknown {
		m.unknownBy(u)
	}
}

func stringFilter(m *matcher, on bool, name string, want []string, got string) {
	if !on {
		return
	}
	if got == "" {
		m.unknownBy(name)
		return
	}
	m.check(slices.Contains(want, got))
}

// maxFilter matches when the recorded value is at most max; an absent value is unknown.
func maxFilter(m *matcher, name string, max, got *int) {
	if max == nil {
		return
	}
	if got == nil {
		m.unknownBy(name)
		return
	}
	m.check(*got <= *max)
}

func boolFilter(m *matcher, name string, want, got *bool) {
	if want == nil {
		return
	}
	if got == nil {
		m.unknownBy(name)
		return
	}
	m.check(*got == *want)
}

// matchEvent applies a resolved filter to an event.
func (e *evalCtx) matchEvent(f *resolvedFilter, ev *Event) *matcher {
	m := newMatcher()
	for _, u := range f.unresolved {
		m.unknownBy(u)
	}
	switch f.Simulator {
	case "only":
		m.check(ev.IsSimulator)
	case "include":
	default:
		m.check(!ev.IsSimulator)
	}
	if len(f.EventKinds) > 0 {
		m.check(slices.Contains(f.EventKinds, ev.Kind))
	}
	rs := ev.ratings()
	if len(f.EventRatings) > 0 {
		m.check(slices.ContainsFunc(rs, func(r string) bool { return slices.Contains(f.EventRatings, r) }))
	}
	if len(f.ExcludeEventRatings) > 0 && len(rs) > 0 {
		m.check(slices.ContainsFunc(rs, func(r string) bool { return !slices.Contains(f.ExcludeEventRatings, r) }))
	}
	if len(f.EventAuthorities) > 0 {
		if ev.Authority == "" {
			m.unknownBy("eventAuthorities")
		} else {
			m.check(containsFold(f.EventAuthorities, ev.Authority))
		}
	}
	if len(f.FSTDTypes) > 0 && ev.IsSimulator {
		if ev.FSTDType == "" {
			m.unknownBy("fstdTypes")
		} else {
			m.check(slices.Contains(f.FSTDTypes, ev.FSTDType))
		}
	}
	if f.classes != nil {
		if ev.Class == "" {
			m.unknownBy("classes")
		} else {
			m.check(f.classes[ev.Class])
		}
	}
	if len(f.ExcludeClasses) > 0 {
		m.check(!slices.Contains(f.ExcludeClasses, ev.Class))
	}
	if len(f.Categories) > 0 {
		cat := ev.Category
		if cat == "" {
			cat = e.v.Category(ev.Class)
		}
		if cat == "" {
			m.unknownBy("categories")
		} else {
			m.check(slices.Contains(f.Categories, cat))
		}
	}
	if f.Has("ulKinds") {
		switch {
		case ev.Class != "ULTRALIGHT":
			m.no()
		case ev.ULKind == "":
			m.unknownBy("ulKinds")
		default:
			m.check(slices.Contains(f.ULKinds, ev.ULKind))
		}
	}
	stringFilter(m, f.Has("typeDesignators"), "typeDesignators", f.TypeDesignators, ev.TypeDesignator)
	stringFilter(m, f.Has("variants"), "variants", f.Variants, ev.Variant)
	if len(f.anyOf) > 0 {
		e.matchAny(m, f.anyOf, func(sub *resolvedFilter) *matcher { return e.matchEvent(sub, ev) })
	}
	return m
}
