package credentials

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestCatalogueGreen(t *testing.T) {
	cat := load(t)
	ex, err := LoadExamples("..")
	if err != nil {
		t.Fatal(err)
	}
	f := Check(cat, ex)
	for _, l := range [][]string{f.Files, f.Refs, f.Keys, f.ExampleErrs, f.Overlaps, f.Interpret} {
		for _, s := range l {
			t.Error(s)
		}
	}
	cov := CheckCoverage(cat)
	for _, s := range cov.Problems {
		t.Error(s)
	}
	if cov.Evaluated+cov.NotEvaluated+cov.PendingOnly != 113 || len(cov.Pending) == 0 {
		t.Errorf("coverage counts: %+v", cov)
	}
}

func TestParseRef(t *testing.T) {
	r, err := ParseRef("easa:FCL.740.A(b)(1)(ii)(C)@sha256:ab12")
	if err != nil || r.Prefix != "easa" || r.Article != "FCL.740.A" || strings.Join(r.Labels, ",") != "b,1,ii,C" || r.Anchor != "sha256:ab12" {
		t.Fatalf("%+v %v", r, err)
	}
	for _, bad := range []string{"FCL.740", "easa:", "easa:FCL.740(b", "Easa:FCL.740"} {
		if _, err := ParseRef(bad); err == nil {
			t.Errorf("%q parsed", bad)
		}
	}
}

func TestResolve(t *testing.T) {
	v, err := LoadVocab("../vocabulary.yaml")
	if err != nil {
		t.Fatal(err)
	}
	rs := NewResolver("..", v)
	for _, ok := range []string{"easa:FCL.740.A(b)(1)(ii)(C)", "easa:FCL.740.A(b)(2)", "easa:SFCL.130(a)(2)(iv)(B)(a)", "easa:SFCL.130(a)(2)(v)(B)", "faa:61.57(c)(1)(iii)", "faa:61.56(i)", "policy:expiring-notice", "de:LuftPersV.45a"} {
		if err := rs.Resolve(ok); err != nil {
			t.Errorf("%s: %v", ok, err)
		}
	}
	for bad, want := range map[string]string{
		"easa:FCL.740.A(b)(9)":      "paragraph (b)(9) not found",
		"easa:FCL.740.A(c)(1)(ii)":  "paragraph (c)(1)(ii) not found",
		"easa:FCL.999":              "no source file",
		"xx:FCL.740":                "unknown ref prefix",
		"policy:nope":               "policy ref not declared",
		"policy:expiring-notice(a)": "has no paragraph labels",
		"not a ref":                 "want <prefix>",
	} {
		err := rs.Resolve(bad)
		if err == nil || !strings.Contains(err.Error(), want) {
			t.Errorf("%s: got %v, want %q", bad, err, want)
		}
	}
	if got := v.Cite(Ref{Prefix: "faa", Article: "61.57", Labels: []string{"c", "1"}}); got != "14 CFR 61.57(c)(1)" {
		t.Errorf("cite %q", got)
	}
	if got := v.Cite(Ref{Raw: "zz:x", Prefix: "zz"}); got != "zz:x" {
		t.Errorf("cite %q", got)
	}
}

func TestLabelOutline(t *testing.T) {
	text := "# X\n\n## Text\n\n(a) *Head.* (1) one\n\n(i) roman\n\n(ii) two\n\n(A) upper\n\n(b) b\n\n(c) c\n\n(d) d\n\n(e) e\n\n(f) f\n\n(g) g\n\n(h) h\n\n(i) letter i\n\n(7) odd\n"
	l := parseLabels(text)
	var got []string
	for _, x := range l {
		got = append(got, strings.Join(x.path, ""))
	}
	if strings.Join(got, " ") != "a a1 a1i a1ii a1iiA b c d e f g h i i7" {
		t.Errorf("outline %v", got)
	}
	if ordinal("x1", "num") != 0 || ordinal("ab", "lower") != 0 || ordinal("a", "upper") != 0 || ordinal("a", "other") != 0 {
		t.Error("ordinal of an impossible label")
	}
	if !findPath(l, nil) {
		t.Error("empty path")
	}
}

// writeRoot builds a module root with the real vocabulary, keys and one source file.
func writeRoot(t *testing.T, files map[string]string) string {
	t.Helper()
	root := t.TempDir()
	cp := func(from, to string) {
		b, err := os.ReadFile(filepath.Join("..", from))
		if err != nil {
			t.Fatal(err)
		}
		write(t, root, to, string(b))
	}
	cp("vocabulary.yaml", "vocabulary.yaml")
	cp("messages/keys.yaml", "messages/keys.yaml")
	cp("sources/easa/fcl-740-a.md", "sources/easa/fcl-740-a.md")
	for p, s := range files {
		write(t, root, p, s)
	}
	return root
}

func write(t *testing.T, root, p, s string) {
	t.Helper()
	full := filepath.Join(root, p)
	if err := os.MkdirAll(filepath.Dir(full), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(full, []byte(s), 0o644); err != nil {
		t.Fatal(err)
	}
}

func TestBrokenCredentials(t *testing.T) {
	root := writeRoot(t, map[string]string{
		"credentials/easa/ratings/a.yaml": `
credential: A
id: easa.rating.a
kind: rating
authority: EASA
selects: { ratings: { classes: [SEP_LAND], authorities: [EASA] } }
evaluations:
  - id: one
    asks: Broken on purpose.
    source: easa:FCL.740.A(b)(9)
    relevant_class: { pooled_with_held: [SEP_LAND, TMG] }
    counting: { within_days: x, flavour: 3, ever: yes }
    passes_if:
      all_of:
        - flight_time: { in_class: maybe, with: examiner, ref: easa:FCL.740.A(b)(1) }
        - bogus_count: { min: 1 }
        - takeoffs: 3
        - takeoffs: { min: 1, ref: easa:FCL.740.A(b)(1) }
          landings: { min: 1, ref: easa:FCL.740.A(b)(1) }
        - landings: { min: 1, ul_credit: { SEP_LAND: THREE_AXIS }, as: pilot_flying, waived_by: { categories: [aeroplane] }, only_if: always, ref: easa:FCL.740.A(b)(1) }
        - id: nested
          any_of:
            - n_of: { n: 1, of: [{ launches: { min: 1, by_this_launch_method: maybe, ref: easa:FCL.740.A(b)(1) } }] }
              ref: easa:FCL.740.A(b)(1)
              all_of: []
    restored_by:
      - { ipc: { simulator: include }, proficiency_check: {} }
    valid_for: { counted_from: yesterday, periods: [{ months: 1 }] }
    on_fail: { what: renewal }
    outcomes: { preset: nope, flavour: x }
    description_key: nope_key
  - id: two
    about: nowhere
    asks: Unknown about.
    source: easa:FCL.740.A(b)(1)
    outcomes: [{ when: always, status: current, message: no.such.key, bogus: 1 }]
  - id: three
    about: licence
    asks: No licence part.
    source: easa:FCL.740.A(b)(1)
    outcomes: recency
  - id: four
    asks: Presets without their options.
    source: easa:FCL.740.A(b)(1)
    passes_if: { any_of: [{ takeoffs: { min: 1, ref: easa:FCL.740.A(b)(1) } }], ref: easa:FCL.740.A(b)(1) }
    outcomes: { preset: passengers, expiring_notice: { days: 1, colour: red } }
  - id: five
    asks: Revalidation without notice.
    source: easa:FCL.740.A(b)(1)
    outcomes: { preset: revalidation }
  - id: six
    asks: Validity with a bad option.
    source: easa:FCL.740.A(b)(1)
    outcomes: { preset: validity, date_of_birth: sometimes, expiring_notice: { days: 5, ref: policy:expiring-notice } }
  - id: seven
    uses: easa.rating.bee#none
  - id: eight
    uses: easa.shared.nope
  - id: nine
    uses: nobody#x
  - id: ten
    requires: [easa.rating.nope]
  - id: eleven
    asks: No outcomes.
    source: easa:FCL.740.A(b)(1)
  - id: twelve
    asks: Unknown keys.
    source: easa:FCL.740.A(b)(1)
    outcomes: [{ when: always, status: current, message: no.such.key, ref: policy:not-applicable }]
    description_key: nope_key
  - id: one
    asks: Repeated id.
    source: easa:FCL.740.A(b)(1)
    outcomes: []
interpretations:
  - { id: x, reading: r, affects: [nope], approved_by: someone, approved_on: null }
  - { id: x, reading: "", affects: [], ref: "policy:nope", approved_by: null, approved_on: null }
`,
		"credentials/easa/ratings/b.yaml": `
credential: B
id: easa.rating.bee
kind: flavour
authority: FAA
selects: { ratings: { classes: [SEP_LAND, MEP_LAND], authorities: [EASA] } }
evaluations:
  - id: p
    uses: easa.shared.p
    with: { extra: 1 }
  - id: q
    uses: easa.shared.p
`,
		"credentials/easa/shared/p.yaml": `
shared: P
id: easa.shared.p
authority: EASA
params: { n: count }
evaluation: { asks: x, source: "easa:FCL.740.A(b)(1)", outcomes: recency }
`,
		"credentials/easa/shared/unused.yaml": `
shared: U
id: easa.shared.unused-wrong
authority: EASA
evaluation: { asks: x }
`,
		"credentials/easa/shared/zdup.yaml": "shared: D\nid: easa.shared.p\n",
		"credentials/easa/shared/bad.yaml":  "shared: [\n",
		"credentials/easa/ratings/bad.yaml": "credential: [\n",
		"credentials/easa/ratings/dup.yaml": "credential: A\nid: easa.rating.a\nkind: rating\n",
		"examples/easa.rating.a.yaml": `
credential: easa.rating.a
examples:
  - { name: x, evaluation: nope, asOf: 2026-01-01, record: {}, expect: { status: current } }
  - { name: x, evaluation: ten, asOf: 2026-01-01, record: {}, expect: { status: current } }
`,
		"examples/wrong-name.yaml": "credential: easa.rating.bee\nexamples: []\n",
		"examples/gone.yaml":       "credential: easa.rating.gone\nexamples: []\n",
	})
	cat, err := Load(root)
	if err != nil {
		t.Fatal(err)
	}
	ex, err := LoadExamples(root)
	if err != nil {
		t.Fatal(err)
	}
	f := Check(cat, ex)
	var all []string
	for _, l := range [][]string{f.Load, f.Files, f.Refs, f.Keys, f.ExampleErrs, f.Overlaps, f.Interpret, f.Unapproved} {
		all = append(all, l...)
	}
	report := strings.Join(all, "\n")
	for _, want := range []string{
		"paragraph (b)(9) not found", "relevant_class has no ref", "within_days needs a number", `"flavour" is not a qualifier`,
		"ever: write ever: true", "in_class: write in_class: true", `"bogus_count" is not a count word`, "write { min: 3, ref: ... }",
		"an item is one count word", "ul_credit.SEP_LAND: want a list", "ul_credit has no ref", "waived_by has no ref",
		"waived_by: list the events", "only_if has no ref", "more than one combinator", "by_this_launch_method: write",
		"all_of has no ref", "restored_by: one event per item", "restored_by proficiency_check has no ref", "counted_from is issue or valid_from",
		"valid_for.periods[0] has no ref", "on_fail has no ref", `preset "nope" is not in the vocabulary`, `outcomes: unknown key "flavour"`,
		`about "nowhere" is not in the vocabulary`, `outcomes[0]: unknown key "bogus"`, "outcomes[0] has no ref",
		"selects no licence part", "passengers preset needs requirements", `expiring_notice: unknown key "colour"`,
		"preset revalidation needs expiring_notice", "date_of_birth is when_no_expiry or required", "has no evaluation none",
		"no shared evaluation", "no credential nobody", "requires easa.rating.nope, which is no credential",
		`evaluation id "one" missing or repeated`, "no outcomes", `affects "nope", which is no evaluation here`, "id missing or repeated",
		"needs reading and affects", "no ref to the paragraph", "approved_by and approved_on go together", "policy ref not declared",
		"description key nope_key", "message no.such.key", `kind "flavour" is not one of the credential kinds`,
		"must be <authority>.<kind>.<file name>", "does not start with the authority", "has no parameter \"extra\"", `parameter "n" not given`,
		"must be <authority>.shared.<file name>", "no credential uses easa.shared.unused-wrong", "duplicate shared id", "duplicate credential id",
		"no evaluation easa.rating.a#nope", "only references other credentials", `example name "x" missing or repeated`,
		"file name must be easa.rating.bee.yaml", "no credential easa.rating.gone", "easa.rating.a and easa.rating.bee both select ratings items",
	} {
		if !strings.Contains(report, want) {
			t.Errorf("report misses %q", want)
		}
	}
	if t.Failed() {
		t.Log(report)
	}
}

func TestExamplesAndCoverage(t *testing.T) {
	root := writeRoot(t, map[string]string{
		"credentials/easa/ratings/a.yaml": `
credential: A
id: easa.rating.a
kind: rating
authority: EASA
selects: { ratings: { classes: [SEP_LAND], authorities: [EASA, LBA] } }
evaluations:
  - id: two
    about: passengers
    asks: Passengers again.
    source: easa:FCL.740.A(b)(1)
    passes_if: { takeoffs: { min: 1, ref: easa:FCL.740.A(b)(1) } }
    outcomes: passengers
  - id: one
    about: passengers
    asks: Passengers.
    source: easa:FCL.740.A(b)(1)
    counting: { within_days: 90, in_class: true }
    passes_if: { takeoffs: { min: 1, ref: easa:FCL.740.A(b)(1) } }
    outcomes: passengers
`,
		"examples/easa.rating.a.yaml": `
credential: easa.rating.a
examples:
  - name: two-subjects
    evaluation: one
    shows: [nothing]
    asOf: 2026-01-01
    record:
      licences: [{ id: l1, authority: EASA, type: PPL(A) }, { id: l2, authority: LBA, type: CPL(A) }]
      ratings: [{ id: r1, licenceId: l1, class: SEP_LAND }, { id: r2, licenceId: l2, class: SEP_LAND }]
    expect: { status: current }
  - name: none
    evaluation: one
    asOf: 2026-01-01
    record: {}
    expect: { status: current }
  - name: no-status
    evaluation: one
    asOf: 2026-01-01
    record: {}
    expect: {}
  - name: lapsed
    evaluation: two
    asOf: 2026-01-01
    record: {}
    expect: { status: expired }
`,
		"coverage/articles.yaml": `
articles:
  - { article: "easa:FCL.740.A", pending: [easa/ratings/a], note: planned }
  - { article: "easa:FCL.740.A", evaluations: [easa.rating.a#nope] }
  - { article: "easa:FCL.999" }
  - { article: "easa:FCL.740", not_evaluated: why, pending: [easa/ratings/b] }
  - { article: "easa:FCL.740.B", pending: [easa/ratings/b] }
`,
		"scope/articles-easa.yaml": "articles:\n  - { cite: FCL.740.A }\n  - { cite: FCL.740 }\n  - { cite: FCL.740.B }\n  - { cite: FCL.060 }\n",
	})
	cat, err := Load(root)
	if err != nil {
		t.Fatal(err)
	}
	ex, _ := LoadExamples(root)
	f := Check(cat, ex)
	cov := CheckCoverage(cat)
	report := strings.Join(append(f.ExampleErrs, cov.Problems...), "\n")
	for _, want := range []string{
		"2 results; name the subject", "no result for the subject", "expect.status is required", `shows "nothing"`,
		"no passing example", "no failing example", "listed twice", "pending credentials/easa/ratings/a.yaml exists",
		"easa.rating.a#nope is no compiled evaluation", "no source file",
		"needs evaluations, pending or not_evaluated", "article easa:FCL.060 of the scope is missing",
		"pending needs a note", "not an article of scope",
		"does not list easa.rating.a#one",
	} {
		if !strings.Contains(report, want) {
			t.Errorf("report misses %q", want)
		}
	}
	if t.Failed() {
		t.Log(report)
	}
	if len(cov.Pending) != 3 || cov.PendingByAuthority["easa"] != 3 {
		t.Errorf("pending: %v %v", cov.Pending, cov.PendingByAuthority)
	}
	if got := CheckCoverage(&Catalogue{Root: t.TempDir()}); len(got.Problems) != 1 {
		t.Errorf("missing coverage file: %v", got.Problems)
	}
	bad := t.TempDir()
	write(t, bad, "coverage/articles.yaml", "articles: [\n")
	if got := CheckCoverage(&Catalogue{Root: bad}); len(got.Problems) != 1 {
		t.Errorf("bad coverage file: %v", got.Problems)
	}
	write(t, bad, "coverage/articles.yaml", "articles: []\n")
	write(t, bad, "scope/articles-easa.yaml", "articles: [\n")
	if got := CheckCoverage(&Catalogue{Root: bad}); len(got.Problems) != 1 {
		t.Errorf("bad scope file: %v", got.Problems)
	}
	if ScopeRef("de", "LuftPersV § 45a") != "de:LuftPersV.45a" || ScopeRef("faa", "14 CFR 61.57") != "faa:61.57" || ArticleOf("bad") != "bad" {
		t.Error("scope refs")
	}
}

func TestLoadErrors(t *testing.T) {
	if _, err := Load(t.TempDir()); err == nil {
		t.Error("no vocabulary")
	}
	if _, err := LoadVocab(filepath.Join(t.TempDir(), "x.yaml")); err == nil {
		t.Error("no file")
	}
	root := t.TempDir()
	write(t, root, "v.yaml", "version: 1\n")
	if _, err := LoadVocab(filepath.Join(root, "v.yaml")); err == nil {
		t.Error("no section")
	}
	write(t, root, "bad.yaml", ": [\n")
	if _, err := LoadVocab(filepath.Join(root, "bad.yaml")); err == nil {
		t.Error("bad yaml")
	}
	write(t, root, "examples/bad.yaml", "examples: [\n")
	if _, err := LoadExamples(root); err == nil {
		t.Error("bad examples")
	}
}

// TestEffectiveDates checks that effective_from and effective_to reach the compiled rule, that
// a using evaluation may set them, and that an inverted period is a compile error.
func TestEffectiveDates(t *testing.T) {
	root := writeRoot(t, map[string]string{
		"credentials/easa/ratings/a.yaml": `
credential: A
id: easa.rating.a
kind: rating
authority: EASA
selects: { ratings: { classes: [SEP_LAND], authorities: [EASA] } }
evaluations:
  - id: old
    about: passengers
    asks: Passengers under the old text.
    source: easa:FCL.740.A(b)(1)
    effective_to: 2026-12-31
    counting: { within_days: 90, in_class: true }
    passes_if: { takeoffs: { min: 3, ref: easa:FCL.740.A(b)(1) } }
    outcomes: passengers
  - id: new
    about: passengers
    asks: Passengers under the new text.
    source: easa:FCL.740.A(b)(1)
    effective_from: 2027-01-01
    counting: { within_days: 90, in_class: true }
    passes_if: { takeoffs: { min: 5, ref: easa:FCL.740.A(b)(1) } }
    outcomes: passengers
  - id: bad
    about: passengers
    asks: Inverted.
    source: easa:FCL.740.A(b)(1)
    effective_from: 2027-01-01
    effective_to: 2026-01-01
    passes_if: { takeoffs: { min: 1, ref: easa:FCL.740.A(b)(1) } }
    outcomes: passengers
`,
		"credentials/easa/ratings/b.yaml": `
credential: B
id: easa.rating.b
kind: rating
authority: EASA
selects: { ratings: { classes: [MEP_LAND], authorities: [EASA] } }
evaluations:
  - id: reused
    uses: easa.rating.a#old
    effective_from: 2026-06-01
    effective_to: 2026-06-30
`,
	})
	cat, err := Load(root)
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(strings.Join(cat.Errors, "\n"), "easa.rating.a#bad: effective_to is before effective_from") {
		t.Errorf("errors: %v", cat.Errors)
	}
	old := cat.Compiled["easa.rating.a#old"].Rule
	if old.EffectiveTo == nil || old.EffectiveTo.String() != "2026-12-31" || old.EffectiveFrom != nil {
		t.Errorf("old: %v %v", old.EffectiveFrom, old.EffectiveTo)
	}
	if r := cat.Compiled["easa.rating.a#new"].Rule; r.EffectiveFrom == nil || r.EffectiveFrom.String() != "2027-01-01" {
		t.Errorf("new: %v", r.EffectiveFrom)
	}
	if r := cat.Compiled["easa.rating.b#reused"].Rule; r.EffectiveFrom.String() != "2026-06-01" || r.EffectiveTo.String() != "2026-06-30" {
		t.Errorf("reused: %v %v", r.EffectiveFrom, r.EffectiveTo)
	}
}
