package credentials

import (
	"slices"
	"strings"
	"testing"

	"github.com/fjaeckel/open-aircrew-rules/engine"
)

func composites(t *testing.T, cat *Catalogue, rec string, asOf string) map[string]Composite {
	t.Helper()
	var r engine.Record
	if err := yamlDecode(rec, &r); err != nil {
		t.Fatal(err)
	}
	out := map[string]Composite{}
	for _, c := range cat.Evaluate(&r, engine.MustDate(asOf)).Credentials {
		out[c.Credential+"|"+c.Subject.ID] = c
	}
	return out
}

// TestComposites checks requirement groups, licence narrowing, limitations and not held.
func TestComposites(t *testing.T) {
	cat := load(t)
	got := composites(t, cat, `
holder: { dateOfBirth: 1976-01-01 }
licences:
  - { id: l-spl, authority: EASA, type: SPL }
  - { id: l-cpl, authority: EASA, type: CPL(A) }
  - { id: l-ppl, authority: EASA, type: PPL(A) }
ratings:
  - { id: r-gld, licenceId: l-spl, class: GLIDER }
  - { id: r-sep, licenceId: l-cpl, class: SEP_LAND, expires: 2027-06-30 }
  - { id: r-sep2, licenceId: l-ppl, class: SEP_LAND, expires: 2027-06-30 }
privileges:
  - { id: p-tow, licenceId: l-spl, kind: SAILPLANE_TOWING }
credentials:
  - { id: c-med, type: EASA_CLASS2_MEDICAL, issued: 2026-06-01 }
  - { id: c-lang, type: LANG_ICAO_LEVEL4, issued: 2023-04-15 }
flights:
  - { date: 2026-08-01, class: GLIDER, launchMethod: winch, minutes: { total: 30, pic: 30 }, launches: 1, landings: { day: 1 } }
`, "2026-09-15")

	spl := got["easa.licence.spl|l-spl"]
	if spl.Status != "lapsed" || spl.DecidedBy.Evaluation != "recency" {
		t.Errorf("spl: %+v", spl)
	}
	if !slices.ContainsFunc(spl.Limitations, func(m Member) bool { return m.Evaluation == "launch_methods_towed" }) {
		t.Errorf("spl limitations: %+v", spl.Limitations)
	}
	tow := got["easa.privilege.sfcl-sailplane-towing|p-tow"]
	if m := tow.Members[1]; tow.Status != "lapsed" || tow.DecidedBy.Evaluation != "recency" || m.Kind != "requires_all" || m.Status != "lapsed" || m.Credential != "easa.licence.spl" {
		t.Errorf("towing: %+v", tow)
	}
	// The rating on the CPL(A) is decided by the CPL(A), not by the PPL(A) on another licence.
	sep := got["easa.rating.sep-land|r-sep"]
	if d := sep.DecidedBy; sep.Status != "expired" || d.Evaluation != "licence" || d.Credential != "easa.licence.cpl-a" {
		t.Errorf("sep on cpl: %+v", sep)
	}
	if sep2 := got["easa.rating.sep-land|r-sep2"]; sep2.Status != "expiring" || sep2.Members[1].Status != "current" {
		t.Errorf("sep on ppl: %+v", sep2)
	}
	if ppl := got["easa.licence.ppl-a|l-ppl"]; ppl.Status != "current" || len(ppl.Members) != 3 {
		t.Errorf("ppl: %+v", ppl)
	}
}

// TestCompositeNotApplicable covers a credential with no deciding member, which counts as
// met when another credential requires it.
func TestCompositeNotApplicable(t *testing.T) {
	root := writeRoot(t, map[string]string{
		"credentials/easa/licences/x.yaml": `
credential: X
id: easa.licence.x
kind: licence
authority: EASA
selects: { licence: { kinds: [PPL_A] } }
evaluations:
  - id: medical
    asks: The medical.
    requires_all: [easa.medical.m]
`,
		"credentials/easa/medicals/m.yaml": `
credential: M
id: easa.medical.m
kind: medical
authority: EASA
selects: { credential: { kinds: [EASA_CLASS2_MEDICAL] } }
evaluations:
  - id: gone
    asks: An evaluation that ended.
    source: easa:FCL.740.A(b)(1)
    effective_to: 2020-01-01
    outcomes: [{ when: always, status: current, message: credential.valid, ref: easa:FCL.740.A(b)(1) }]
`,
	})
	cat, err := Load(root)
	if err != nil || len(cat.Errors) > 0 {
		t.Fatal(err, cat.Errors)
	}
	got := composites(t, cat, "licences: [{ id: l, authority: EASA, type: PPL(A) }]\ncredentials: [{ id: c, type: EASA_CLASS2_MEDICAL }]\n", "2026-01-01")
	if m := got["easa.medical.m|c"]; m.Status != "not_applicable" || m.DecidedBy != nil {
		t.Errorf("medical: %+v", m)
	}
	if x := got["easa.licence.x|l"]; x.Status != "current" || x.DecidedBy.Credential != "easa.medical.m" {
		t.Errorf("licence: %+v", x)
	}
}

func TestRelated(t *testing.T) {
	rec := &engine.Record{
		Ratings:    []engine.Rating{{ID: "r", LicenceID: "l", Class: "SEP_LAND"}, {ID: "r2", LicenceID: "l", Class: "TMG"}},
		Privileges: []engine.Privilege{{ID: "p", LicenceID: "l"}},
		Variants:   []engine.Variant{{ID: "v", RatingID: "r"}},
	}
	for _, tc := range []struct {
		s    engine.Subject
		want string
	}{
		{engine.Subject{Kind: "licence", ID: "l"}, "licence:l"},
		{engine.Subject{Kind: "flight_review", ID: "l"}, "licence:l"},
		{engine.Subject{Kind: "training"}, ""},
		{engine.Subject{Kind: "rating", ID: "r"}, "rating:r licence:l"},
		{engine.Subject{Kind: "launch_method", ID: "r", Detail: "winch"}, "rating:r licence:l"},
		{engine.Subject{Kind: "type", ID: "v"}, "rating:r licence:l"},
		{engine.Subject{Kind: "passengers", ID: "l", Class: "TMG"}, "licence:l rating:r2"},
		{engine.Subject{Kind: "privilege", ID: "p"}, "privilege:p licence:l"},
		{engine.Subject{Kind: "credential", ID: "c"}, "credential:c"},
		{engine.Subject{Kind: "other", ID: "x"}, ""},
	} {
		if got := strings.Join(related(rec, tc.s), " "); got != tc.want {
			t.Errorf("%+v: %q, want %q", tc.s, got, tc.want)
		}
	}
	if licenceOf(rec, engine.Subject{Kind: "credential", ID: "c"}) != "" || licenceOf(rec, engine.Subject{Kind: "privilege", ID: "p"}) != "l" {
		t.Error("licenceOf")
	}
	for s, want := range map[string]int{"current": 0, "expiring": 1, "unknown": 2, "odd": 2, "expired": 3, "lapsed": 3, "not_applicable": -1} {
		if rank(s) != want {
			t.Errorf("rank(%s)", s)
		}
	}
	ms := []Member{{Evaluation: "a", Status: "expired"}, {Evaluation: "b", Status: "current"}, {Evaluation: "c", Status: "not_applicable"}}
	if st, d := decide(ms, true); st != "current" || d.Evaluation != "b" {
		t.Errorf("best: %s %+v", st, d)
	}
	if st, d := decide(ms, false); st != "expired" || d.Evaluation != "a" {
		t.Errorf("worst: %s %+v", st, d)
	}
	if st, d := decide(ms[2:], false); st != "not_applicable" || d != nil {
		t.Errorf("none: %s %+v", st, d)
	}
}

func TestOutcomeClass(t *testing.T) {
	for _, tc := range []struct{ status, root, want string }{
		{"current", "unmet", Passing}, {"expiring", "met", Passing}, {"expiring", "", Passing},
		{"expiring", "unmet", Failing}, {"expiring", "unknown", Failing}, {"expired", "", Failing},
		{"lapsed", "met", Failing}, {"unknown", "", Failing}, {"not_applicable", "", Failing},
	} {
		if got := OutcomeClass(tc.status, tc.root); got != tc.want {
			t.Errorf("%s/%s: %s", tc.status, tc.root, got)
		}
	}
	if CompositeClass("expiring") != Passing || CompositeClass("unknown") != Failing {
		t.Error("composite class")
	}
}

func TestRunCompositeExample(t *testing.T) {
	cat := load(t)
	var rec engine.Record
	if err := yamlDecode("licences: [{ id: a, authority: EASA, type: PPL(A) }, { id: b, authority: EASA, type: PPL(A) }]\n", &rec); err != nil {
		t.Fatal(err)
	}
	x := Example{Composite: true, Outcome: "current", AsOf: engine.MustDate("2026-01-01"), Record: rec}
	if d := cat.RunExample("easa.licence.ppl-a", x); len(d) != 1 || !strings.Contains(d[0], "holds 2") {
		t.Errorf("two held: %v", d)
	}
	x.Expect.Subject = engine.Subject{Kind: "licence", ID: "a"}
	x.Expect.DecidedBy = &ExpectDecider{Evaluation: "language", Credential: "easa.medical.class-1"}
	d := strings.Join(cat.RunExample("easa.licence.ppl-a", x), "\n")
	if !strings.Contains(d, "status: want current, got unknown") || !strings.Contains(d, "decidedBy: want language easa.medical.class-1, got medical ") {
		t.Errorf("diffs: %s", d)
	}
	y := Example{Evaluation: "medical", Outcome: "unknown", AsOf: x.AsOf, Record: rec, Expect: ExampleExpect{Subject: x.Expect.Subject, DecidedBy: &ExpectDecider{}}}
	if d := cat.RunExample("easa.licence.ppl-a", y); len(d) != 1 || !strings.Contains(d[0], "belongs to a composite") {
		t.Errorf("decidedBy on an evaluation: %v", d)
	}
	if _, ok := cat.Credential("easa.licence.ppl-a"); !ok {
		t.Error("Credential")
	}
	if (&Findings{}).Problems() != 0 {
		t.Error("Problems")
	}
}

// TestConversionWords compiles and runs the count words and qualifiers added for the
// conversion of the remaining credentials.
func TestConversionWords(t *testing.T) {
	root := writeRoot(t, map[string]string{
		"credentials/easa/ratings/v.yaml": `
credential: V
id: easa.rating.v
kind: rating
authority: EASA
selects: { ratings: { classes: [MEP_LAND] } }
evaluations:
  - id: variant
    about: variants
    asks: Every new word at once.
    source: easa:FCL.740.A(b)(1)
    counting: { within_months: 12, in_variant: true }
    passes_if:
      ref: easa:FCL.740.A(b)(1)
      all_of:
        - route_sectors: { min: 1, flagged: { examinerOnBoard: true }, min_landings: 1, max_engines: 2, max_mtom_kg: 5700, ref: easa:FCL.740.A(b)(1) }
        - solo_time: { min_hours: 1, ref: easa:FCL.740.A(b)(1) }
        - instruction_given_time: { min_hours: 1, ref: easa:FCL.740.A(b)(1) }
        - instruction_or_examining_time: { id: instructing, min_hours: 1, ref: easa:FCL.740.A(b)(1) }
        - mountain_landings: { min: 1, ref: easa:FCL.740.A(b)(1) }
        - tows: { min: 1, tow_kinds: [banner], ref: easa:FCL.740.A(b)(1) }
    outcomes: recency
`,
	})
	cat, err := Load(root)
	if err != nil || len(cat.Errors) > 0 {
		t.Fatal(err, cat.Errors)
	}
	var rec engine.Record
	if err := yamlDecode(`
licences: [{ id: l, authority: EASA, type: CPL(A) }]
ratings: [{ id: r, licenceId: l, class: MEP_LAND }]
variants: [{ id: v, ratingId: r, name: PA34 }]
flights:
  - { date: 2026-05-01, class: MEP_LAND, variant: PA34, engines: 2, mtomKg: 2000, cruiseMinutes: 30, mountainLandings: 1, towKind: banner, minutes: { total: 120, pic: 60, dualGiven: 60 }, takeoffs: { day: 1 }, landings: { day: 1 }, flags: { examinerOnBoard: true, towFlight: true } }
  - { date: 2026-05-02, class: MEP_LAND, variant: OTHER, engines: 2, mtomKg: 2000, cruiseMinutes: 30, minutes: { total: 600, pic: 600 }, landings: { day: 1 } }
`, &rec); err != nil {
		t.Fatal(err)
	}
	evs, _ := engine.EvaluateRule(cat.Engine, cat.Compiled["easa.rating.v#variant"].Rule, &rec, engine.MustDate("2026-09-15"))
	if len(evs) != 1 || evs[0].Subject.Kind != "type" || evs[0].Status != "current" || !limitation(evs[0].Subject.Kind) {
		t.Fatalf("%+v", evs)
	}
	for _, row := range evs[0].Requirements {
		if !row.Met || (row.ID == "solo_time" && row.Current != 60) {
			t.Errorf("row %+v", row)
		}
	}
}
