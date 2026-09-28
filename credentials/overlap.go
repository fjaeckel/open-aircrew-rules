package credentials

import (
	"fmt"
	"slices"
	"sort"
	"strings"

	"github.com/fjaeckel/open-aircrew-rules/engine"
)

// Overlaps reports pairs of credential files that can claim the same record item: the same
// part (licence, ratings, privilege, credential) with intersecting kinds, classes,
// authorities and licence kinds.
func Overlaps(cat *Catalogue) []string {
	v := cat.Engine.Vocabulary
	auths := append(slices.Clone(v.RecordAuthorities), "")
	lk := []string{""}
	for k := range v.LicenceKinds {
		lk = append(lk, k)
	}
	var out []string
	cs := cat.Credentials
	for i := range cs {
		for j := i + 1; j < len(cs); j++ {
			for _, part := range []string{"licence", "ratings", "privilege", "credential"} {
				a, b := cs[i].Selects.part(part), cs[j].Selects.part(part)
				if a == nil || b == nil {
					continue
				}
				if intersects(part, a, b, auths, lk, v) {
					out = append(out, fmt.Sprintf("%s and %s both select %s items (%s)", cs[i].ID, cs[j].ID, part, describe(part, a)))
				}
			}
		}
	}
	sort.Strings(out)
	return out
}

func describe(part string, p *Part) string {
	if part == "ratings" {
		return "classes " + strings.Join(p.Classes, ", ")
	}
	return strings.Join(p.Kinds, ", ")
}

func intersects(part string, a, b *Part, auths, lk []string, v *engine.Vocabulary) bool {
	allowed := func(p *Part, universe []string, include, exclude []string) []string {
		var out []string
		for _, x := range universe {
			if len(include) > 0 && !slices.ContainsFunc(include, func(s string) bool { return strings.EqualFold(s, x) }) {
				continue
			}
			if slices.ContainsFunc(exclude, func(s string) bool { return strings.EqualFold(s, x) }) {
				continue
			}
			out = append(out, x)
		}
		return out
	}
	overlap := func(x, y []string) bool {
		return slices.ContainsFunc(x, func(s string) bool { return slices.Contains(y, s) })
	}
	if !overlap(allowed(a, auths, a.Authorities, a.NotAuthorities), allowed(b, auths, b.Authorities, b.NotAuthorities)) {
		return false
	}
	kindsA, kindsB := a.LicenceKinds, b.LicenceKinds
	if part == "licence" {
		kindsA, kindsB = append(slices.Clone(a.Kinds), a.LicenceKinds...), append(slices.Clone(b.Kinds), b.LicenceKinds...)
	}
	if !overlap(allowed(a, lk, kindsA, a.NotLicenceKinds), allowed(b, lk, kindsB, b.NotLicenceKinds)) {
		return false
	}
	switch part {
	case "ratings":
		classes := make([]string, 0, len(v.Classes))
		for c := range v.Classes {
			classes = append(classes, c)
		}
		return overlap(allowed(a, classes, a.Classes, nil), allowed(b, classes, b.Classes, nil))
	case "privilege", "credential":
		return len(a.Kinds) == 0 || len(b.Kinds) == 0 || overlap(a.Kinds, b.Kinds)
	}
	return true
}
