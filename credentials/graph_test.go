package credentials

import (
	"fmt"
	"strings"
	"testing"
)

const written = `
    asks: Takeoffs.
    source: easa:FCL.740.A(b)(1)
    passes_if: { takeoffs: { min: 1, ref: easa:FCL.740.A(b)(1) } }
    outcomes: recency`

// rating is a rating credential on its own class with the evaluations given.
func rating(name, class, evaluations string) string {
	return fmt.Sprintf(`
credential: %[1]s
id: easa.rating.%[1]s
kind: rating
authority: EASA
selects: { ratings: { classes: [%[2]s], authorities: [EASA] } }
evaluations:%[3]s
`, name, class, evaluations)
}

func shared(name, params, evaluation string) string {
	return fmt.Sprintf("shared: %[1]s\nid: easa.shared.%[1]s\nauthority: EASA\nparams: {%[2]s}\nevaluation: %[3]s\n", name, params, evaluation)
}

func loadErrors(t *testing.T, files map[string]string) (*Catalogue, string) {
	t.Helper()
	cat, err := Load(writeRoot(t, files))
	if err != nil {
		t.Fatal(err)
	}
	return cat, strings.Join(cat.Errors, "\n")
}

func wantErrors(t *testing.T, got string, want ...string) {
	t.Helper()
	for _, w := range want {
		if !strings.Contains(got, w) {
			t.Errorf("errors miss %q:\n%s", w, got)
		}
	}
}

func TestUsesCycles(t *testing.T) {
	for name, tc := range map[string]struct {
		files map[string]string
		want  []string
	}{
		"two nodes": {map[string]string{
			"credentials/easa/ratings/a.yaml": rating("a", "SEP_LAND", "\n  - id: x\n    uses: easa.rating.b#y"),
			"credentials/easa/ratings/b.yaml": rating("b", "MEP_LAND", "\n  - id: y\n    uses: easa.rating.a#x"),
		}, []string{"uses cycle: easa.rating.a#x -> easa.rating.b#y -> easa.rating.a#x"}},
		"self": {map[string]string{
			"credentials/easa/ratings/a.yaml": rating("a", "SEP_LAND", "\n  - id: x\n    uses: easa.rating.a#x"),
		}, []string{"uses cycle: easa.rating.a#x -> easa.rating.a#x"}},
		"three shared": {map[string]string{
			"credentials/easa/ratings/a.yaml": rating("a", "SEP_LAND", "\n  - id: x\n    uses: easa.shared.s1"),
			"credentials/easa/shared/s1.yaml": shared("s1", "", "{ uses: easa.shared.s2 }"),
			"credentials/easa/shared/s2.yaml": shared("s2", "", "{ uses: easa.shared.s3 }"),
			"credentials/easa/shared/s3.yaml": shared("s3", "", "{ uses: easa.shared.s1 }"),
		}, []string{"uses cycle: easa.shared.s1 -> easa.shared.s2 -> easa.shared.s3 -> easa.shared.s1"}},
		"through a shared evaluation back to a credential": {map[string]string{
			"credentials/easa/ratings/a.yaml": rating("a", "SEP_LAND", "\n  - id: x\n    uses: easa.shared.y"),
			"credentials/easa/shared/y.yaml":  shared("y", "", "{ uses: easa.rating.a#x }"),
		}, []string{
			"uses cycle: easa.rating.a#x -> easa.shared.y -> easa.rating.a#x",
			"easa.shared.y: uses easa.rating.a#x: a shared evaluation may use only shared evaluations",
		}},
	} {
		t.Run(name, func(t *testing.T) {
			cat, got := loadErrors(t, tc.files)
			wantErrors(t, got, tc.want...)
			if strings.Count(got, "uses cycle") != 1 || len(cat.Errors) != len(tc.want) {
				t.Errorf("want only %v:\n%s", tc.want, got)
			}
		})
	}
}

func TestUsesDeepChain(t *testing.T) {
	files := map[string]string{
		"credentials/easa/shared/leaf.yaml": shared("leaf", "n: count", `{ asks: Takeoffs., source: "easa:FCL.740.A(b)(1)", passes_if: { takeoffs: { min: $n, ref: "easa:FCL.740.A(b)(1)" } }, outcomes: recency }`),
		"credentials/easa/shared/mid.yaml":  shared("mid", "m: count", `{ uses: easa.shared.leaf, with: { n: $m } }`),
		"credentials/easa/ratings/c5.yaml":  rating("c5", "TMG", "\n  - id: e\n    uses: easa.shared.mid\n    with: { m: 4 }\n  - id: w"+written),
	}
	classes := []string{"", "SEP_LAND", "MEP_LAND", "SEP_SEA", "MEP_SEA"}
	for i := 1; i <= 4; i++ {
		files[fmt.Sprintf("credentials/easa/ratings/c%d.yaml", i)] = rating(fmt.Sprintf("c%d", i), classes[i], fmt.Sprintf("\n  - id: e\n    uses: easa.rating.c%d#e", i+1))
	}
	cat, got := loadErrors(t, files)
	if got != "" {
		t.Fatal(got)
	}
	for i := 1; i <= 5; i++ {
		c := cat.Compiled[fmt.Sprintf("easa.rating.c%d#e", i)]
		if c == nil || c.Eval.Shared.ID != "easa.shared.mid" || c.Eval.Asks != "Takeoffs." {
			t.Fatalf("c%d: %+v", i, c)
		}
	}
	if r := cat.Compiled["easa.rating.c1#e"].Rule.Requirements; r == nil || r.Min == nil || *r.Min != 4 {
		t.Errorf("parameter not passed through: %+v", r)
	}
	c1, _ := cat.Credential("easa.rating.c1")
	a, _ := cat.resolve(c1, c1.Evaluations[0])
	b, _ := cat.resolve(c1, c1.Evaluations[0])
	if a != b {
		t.Error("resolve is not memoised")
	}
	f := Check(cat, nil)
	if strings.Contains(strings.Join(f.Files, "\n"), "no credential uses") {
		t.Errorf("a shared evaluation used through another counts as used: %v", f.Files)
	}
	d := cat.Dependencies()
	if len(d.Uses) != 6 || d.Uses[0] != (Use{"easa.rating.c1#e", "easa.rating.c2#e"}) || d.Uses[5] != (Use{"easa.shared.mid", "easa.shared.leaf"}) {
		t.Errorf("uses: %v", d.Uses)
	}
}

func TestDependencyRules(t *testing.T) {
	cat, got := loadErrors(t, map[string]string{
		"credentials/easa/ratings/a.yaml": rating("a", "SEP_LAND", `
  - id: own`+written+`
  - id: needs
    requires_all: [easa.rating.a]
    requires_any: [easa.rating.b]
  - id: borrowed
    uses: easa.rating.b#needs`),
		"credentials/easa/ratings/b.yaml": rating("b", "MEP_LAND", `
  - id: needs
    requires_all: [easa.rating.a]
  - id: param
    uses: easa.shared.p
    with: { u: x }`),
		"credentials/easa/shared/p.yaml": shared("p", "u: count", "{ uses: $u }"),
	})
	wantErrors(t, got,
		"easa.rating.a#needs: requires the credential itself",
		"requirement cycle: easa.rating.a -> easa.rating.b -> easa.rating.a",
		"easa.rating.a#borrowed: uses easa.rating.b#needs: it only names required credentials",
		"easa.shared.p: uses $u: name a shared evaluation, not a parameter",
	)
	if strings.Contains(got, "requirement cycle: easa.rating.a -> easa.rating.a") {
		t.Errorf("self-requirement reported as a cycle too:\n%s", got)
	}
	d := cat.Dependencies()
	want := []Requirement{
		{"easa.rating.a", "needs", "easa.rating.a", false},
		{"easa.rating.a", "needs", "easa.rating.b", true},
		{"easa.rating.b", "needs", "easa.rating.a", false},
	}
	if fmt.Sprint(d.Requires) != fmt.Sprint(want) {
		t.Errorf("requires: %v", d.Requires)
	}
}

func TestCatalogueDependenciesClean(t *testing.T) {
	cat := load(t)
	d := cat.Dependencies()
	if len(d.Requires) == 0 || len(d.Uses) == 0 {
		t.Fatalf("empty graph: %d requires, %d uses", len(d.Requires), len(d.Uses))
	}
}
