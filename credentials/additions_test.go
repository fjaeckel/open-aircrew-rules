package credentials

import (
	"slices"
	"strings"
	"testing"
)

// TestFormatAdditions covers ul_kinds in selects and only_for, association refs and the
// minimum-mass form of ul_credit.
func TestFormatAdditions(t *testing.T) {
	root := writeRoot(t, map[string]string{
		"associations.yaml": `
associations:
  - { id: dulv:used, publisher: DULV, title: Used, delegated_by: "easa:FCL.740.A(b)(1)", verified: false }
  - { id: dulv:unused, publisher: DULV, title: Unused, delegated_by: "easa:FCL.740.A(b)(1)", verified: false }
  - { id: daec:wrong-publisher, publisher: NOBODY, title: X, delegated_by: "policy:expiring-notice", verified: false }
  - { id: dulv:bad-statute, publisher: DAEC, title: X, delegated_by: "easa:FCL.740.A(b)(9)", verified: true }
`,
		"credentials/easa/ratings/u3.yaml": `
credential: U3
id: easa.rating.u3
kind: rating
authority: EASA
selects: { ratings: { classes: [ULTRALIGHT], ul_kinds: [THREE_AXIS] } }
evaluations:
  - id: one
    asks: Three-axis only.
    source: easa:FCL.740.A(b)(1)
    only_for: { ul_kinds: [THREE_AXIS] }
    passes_if:
      ref: [easa:FCL.740.A(b)(1), assoc:dulv:used]
      all_of:
        - flight_time: { min_hours: 1, ul_credit: { SEP_LAND: { ul_kinds: [THREE_AXIS], min_mtom_kg: 450 }, ref: easa:FCL.740.A(b)(1) }, ref: [easa:FCL.740.A(b)(1), assoc:dulv:used] }
    outcomes: recency
`,
		"credentials/easa/ratings/ub.yaml": `
credential: UB
id: easa.rating.ub
kind: rating
authority: EASA
selects: { ratings: { classes: [TMG] } }
evaluations:
  - id: one
    asks: Broken on purpose.
    source: easa:FCL.740.A(b)(1)
    passes_if:
      ref: easa:FCL.740.A(b)(1)
      all_of:
        - flight_time: { min_hours: 1, ref: [assoc:dulv:used] }
        - takeoffs: { min: 1, ul_credit: { SEP_LAND: { ul_kinds: [THREE_AXIS] }, ref: easa:FCL.740.A(b)(1) }, ref: [assoc:dulv:missing, assoc:dulv:used(a)] }
    outcomes: recency
`,
		"credentials/easa/ratings/uw.yaml": `
credential: UW
id: easa.rating.uw
kind: rating
authority: EASA
selects: { ratings: { classes: [ULTRALIGHT], ul_kinds: [WEIGHT_SHIFT, BOGUS] } }
evaluations:
  - id: one
    asks: Weight-shift only.
    source: easa:FCL.740.A(b)(1)
    scope: { classes: [ULTRALIGHT], ul_kinds: [NOPE] }
    outcomes: recency
`,
		"credentials/easa/ratings/un.yaml": `
credential: UN
id: easa.rating.un
kind: rating
authority: EASA
selects: { ratings: { classes: [ULTRALIGHT], ul_kinds: [none, THREE_AXIS] } }
evaluations:
  - id: one
    asks: No kind.
    source: easa:FCL.740.A(b)(1)
    outcomes: recency
`,
	})
	cat, err := Load(root)
	if err != nil {
		t.Fatal(err)
	}
	f := Check(cat, nil)
	var all []string
	for _, l := range [][]string{f.Load, f.Files, f.Refs, f.Overlaps} {
		all = append(all, l...)
	}
	report := strings.Join(all, "\n")
	for _, want := range []string{
		"ul_credit.SEP_LAND: want { ul_kinds: [...], min_mtom_kg: n }",
		`selects.ratings: ul_kinds "BOGUS" is not an ultralight kind`,
		`easa.rating.uw#one: ul_kinds "NOPE" is not an ultralight kind`,
		"association document not declared in associations.yaml",
		"an association ref has no paragraph labels",
		"cites assoc:dulv:used without the statute delegating it (easa:FCL.740.A(b)(1))",
		"association dulv:unused is cited nowhere",
		`publisher "NOBODY" is not a record authority`,
		"the id starts with the publisher (nobody:<document>)",
		"delegated_by names the statute paragraph that delegates it",
		"paragraph (b)(9) not found",
		"easa.rating.u3 and easa.rating.un both select ratings items",
	} {
		if !strings.Contains(report, want) {
			t.Errorf("report misses %q", want)
		}
	}
	if strings.Contains(report, "easa.rating.u3 and easa.rating.uw") || strings.Contains(report, "easa.rating.un and easa.rating.uw") {
		t.Error("disjoint ultralight kinds reported as overlapping")
	}
	if strings.Contains(report, "dulv:used is cited nowhere") {
		t.Error("a cited association reported as unused")
	}
	if t.Failed() {
		t.Log(report)
	}
	c := cat.Compiled["easa.rating.u3#one"]
	if c == nil {
		t.Fatal("u3 not compiled")
	}
	if !slices.Equal(c.Rule.AppliesTo.ULKinds, []string{"THREE_AXIS"}) {
		t.Errorf("applies_to ulKinds %v", c.Rule.AppliesTo.ULKinds)
	}
	if uc := c.Rule.Requirements.AllOf[0].Filter.ULCredit; len(uc) != 1 || uc[0].MinMTOMKg == nil || *uc[0].MinMTOMKg != 450 {
		t.Errorf("ul_credit %+v", uc)
	}
	if got := cat.held["easa.rating.un"].ULKinds; !slices.Equal(got, []string{"none", "THREE_AXIS"}) {
		t.Errorf("held ulKinds %v", got)
	}
}

func TestLoadAssociations(t *testing.T) {
	if m, err := LoadAssociations(t.TempDir()); err != nil || len(m) != 0 {
		t.Errorf("missing file: %v %v", m, err)
	}
	dup := writeRoot(t, map[string]string{"associations.yaml": "associations:\n  - { id: dulv:x }\n  - { id: dulv:x }\n"})
	if _, err := LoadAssociations(dup); err == nil || !strings.Contains(err.Error(), "declared twice") {
		t.Errorf("duplicate: %v", err)
	}
	bad := writeRoot(t, map[string]string{"associations.yaml": "associations: [\n"})
	if _, err := LoadAssociations(bad); err == nil {
		t.Error("malformed file loaded")
	}
	if _, err := Load(bad); err == nil {
		t.Error("catalogue loaded with a malformed associations.yaml")
	}
}
