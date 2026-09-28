package engine

import (
	"fmt"
	"math"
	"slices"
	"sort"
	"strings"
	"time"
)

// Expect is one expected evaluation; absent fields are not compared.
type Expect struct {
	Subject       Subject                `yaml:"subject"`
	Status        string                 `yaml:"status"`
	MessageKey    string                 `yaml:"messageKey"`
	MessageParams map[string]any         `yaml:"messageParams"`
	ExpiresOn     *Date                  `yaml:"expiresOn"`
	WindowOpensAt *Date                  `yaml:"windowOpensAt"`
	ValidUntil    *Date                  `yaml:"validUntil"`
	Requirements  map[string]ExpectedRow `yaml:"requirements"`
}

// ExpectedRow is an expected requirement row; absent fields are not compared.
type ExpectedRow struct {
	Current      *float64       `yaml:"current"`
	Required     *float64       `yaml:"required"`
	Met          *bool          `yaml:"met"`
	Tracked      *bool          `yaml:"tracked"`
	ValidUntil   *Date          `yaml:"validUntil"`
	LastDate     *Date          `yaml:"lastDate"`
	MessageKey   *string        `yaml:"messageKey"`
	RemedyKey    *string        `yaml:"remedyKey"`
	RemedyParams map[string]any `yaml:"remedyParams"`
	Absent       bool           `yaml:"absent"`
}

func subjectMatches(want, got Subject) bool {
	return want.Kind == got.Kind &&
		(want.ID == "" || want.ID == got.ID) &&
		(want.Class == "" || want.Class == got.Class) &&
		(want.ULKind == "" || want.ULKind == got.ULKind) &&
		(want.Detail == "" || want.Detail == got.Detail)
}

func fmtSubject(s Subject) string {
	parts := []string{s.Kind}
	for _, p := range []string{s.ID, s.Class, s.ULKind, s.Detail} {
		if p != "" {
			parts = append(parts, p)
		}
	}
	return strings.Join(parts, "/")
}

func diffEvaluation(x Expect, ev Evaluation) []string {
	var d []string
	p := "subject " + fmtSubject(ev.Subject) + ": "
	cmp := func(field, want, got string) {
		if want != "" && want != got {
			d = append(d, fmt.Sprintf("%s%s: want %s, got %s", p, field, want, orNone(got)))
		}
	}
	cmp("status", x.Status, ev.Status)
	cmp("messageKey", x.MessageKey, ev.MessageKey)
	cmp("expiresOn", dateStr(x.ExpiresOn), dateStr(ev.ExpiresOn))
	cmp("windowOpensAt", dateStr(x.WindowOpensAt), dateStr(ev.WindowOpensAt))
	cmp("validUntil", dateStr(x.ValidUntil), dateStr(ev.ValidUntil))
	d = append(d, diffParams(p+"messageParams", x.MessageParams, ev.MessageParams)...)
	ids := make([]string, 0, len(x.Requirements))
	for id := range x.Requirements {
		ids = append(ids, id)
	}
	sort.Strings(ids)
	for _, id := range ids {
		want := x.Requirements[id]
		i := slices.IndexFunc(ev.Requirements, func(r RequirementResult) bool { return r.ID == id })
		if want.Absent {
			if i >= 0 {
				d = append(d, fmt.Sprintf("%srequirement %s: want absent, got a row", p, id))
			}
			continue
		}
		if i < 0 {
			d = append(d, fmt.Sprintf("%srequirement %s: want a row, got none", p, id))
			continue
		}
		d = append(d, diffRow(p+"requirement "+id+": ", want, ev.Requirements[i])...)
	}
	return d
}

func diffRow(p string, want ExpectedRow, got RequirementResult) []string {
	var d []string
	num := func(field string, w *float64, g float64) {
		if w != nil && math.Abs(*w-g) > 1e-9 {
			d = append(d, fmt.Sprintf("%s%s: want %v, got %v", p, field, *w, g))
		}
	}
	boolean := func(field string, w *bool, g bool) {
		if w != nil && *w != g {
			d = append(d, fmt.Sprintf("%s%s: want %v, got %v", p, field, *w, g))
		}
	}
	str := func(field string, w *string, g string) {
		if w != nil && *w != g {
			d = append(d, fmt.Sprintf("%s%s: want %q, got %q", p, field, *w, g))
		}
	}
	num("current", want.Current, got.Current)
	num("required", want.Required, got.Required)
	boolean("met", want.Met, got.Met)
	boolean("tracked", want.Tracked, got.Tracked)
	if want.ValidUntil != nil && dateStr(want.ValidUntil) != dateStr(got.ValidUntil) {
		d = append(d, fmt.Sprintf("%svalidUntil: want %s, got %s", p, want.ValidUntil, orNone(dateStr(got.ValidUntil))))
	}
	if want.LastDate != nil && dateStr(want.LastDate) != dateStr(got.LastDate) {
		d = append(d, fmt.Sprintf("%slastDate: want %s, got %s", p, want.LastDate, orNone(dateStr(got.LastDate))))
	}
	str("messageKey", want.MessageKey, got.MessageKey)
	str("remedyKey", want.RemedyKey, got.RemedyKey)
	d = append(d, diffParams(p+"remedyParams", want.RemedyParams, got.RemedyParams)...)
	return d
}

func diffParams(p string, want, got map[string]any) []string {
	var d []string
	keys := make([]string, 0, len(want))
	for k := range want {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	for _, k := range keys {
		g, ok := got[k]
		if !ok {
			d = append(d, fmt.Sprintf("%s.%s: want %s, got none", p, k, paramString(want[k])))
			continue
		}
		if paramString(want[k]) != paramString(g) {
			d = append(d, fmt.Sprintf("%s.%s: want %s, got %s", p, k, paramString(want[k]), paramString(g)))
		}
	}
	return d
}

func paramString(v any) string {
	if t, ok := v.(time.Time); ok {
		return t.Format(time.DateOnly)
	}
	return fmt.Sprint(v)
}

func dateStr(d *Date) string {
	if d == nil {
		return ""
	}
	return d.String()
}

func orNone(s string) string {
	if s == "" {
		return "(none)"
	}
	return s
}

// DiffExpect compares one evaluation with an expectation; absent fields are not compared.
func DiffExpect(x Expect, ev Evaluation) []string { return diffEvaluation(x, ev) }

// SubjectMatches reports whether got matches the non-empty fields of want.
func SubjectMatches(want, got Subject) bool { return subjectMatches(want, got) }
