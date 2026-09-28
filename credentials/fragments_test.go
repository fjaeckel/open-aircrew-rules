package credentials

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// TestFragments checks that fragments are ignored by default and merged with Fragments.
func TestFragments(t *testing.T) {
	root := writeRoot(t, map[string]string{
		"credentials/easa/ratings/a.yaml": `
credential: A
id: easa.rating.a
kind: rating
authority: EASA
selects: { ratings: { classes: [SEP_LAND], authorities: [EASA] } }
evaluations:
  - id: one
    about: passengers
    asks: Passengers.
    source: easa:FCL.740.A(b)(1)
    passes_if: { takeoffs: { min: 1, ref: easa:FCL.740.A(b)(1) } }
    outcomes: [{ when: always, status: current, message: frag.key, ref: policy:frag-policy }]
`,
		"scope/articles-easa.yaml": "articles:\n  - { cite: FCL.740.A }\n  - { cite: FCL.740 }\n",
		"coverage/articles.yaml": `
articles:
  - { article: "easa:FCL.740.A", pending: [easa/ratings/a, easa/ratings/b], note: planned }
`,
		"fragments/coverage/README.md": "readme\n",
		"fragments/coverage/a.yaml": `
articles:
  - { article: "easa:FCL.740.A", evaluations: [easa.rating.a#one], done_pending: [easa/ratings/a, easa/ratings/z], pending: [easa/ratings/c, easa/ratings/b], note: still b and c }
  - { article: "easa:FCL.740", not_evaluated: why }
`,
		"fragments/coverage/b.yaml":   "articles: [\n",
		"fragments/keys/a.yaml":       "keys:\n  - { key: frag.key, kind: message, params: [] }\n",
		"fragments/policies/a.yaml":   "policies:\n  - { id: frag-policy, statement: s, rationale: r }\n  - { id: frag-unused, statement: s, rationale: r }\n",
		"fragments/changelog/a.md":    "- x\n",
		"fragments/vocab-requests/.x": "hidden\n",
	})
	fr, err := Fragments(root)
	if err != nil || len(fr["coverage"]) != 2 || len(fr["keys"]) != 1 || len(fr["changelog"]) != 1 || len(fr["vocab-requests"]) != 0 {
		t.Fatalf("fragments %v %v", fr, err)
	}

	strict, err := Load(root)
	if err != nil {
		t.Fatal(err)
	}
	f := Check(strict, nil)
	if !strings.Contains(strings.Join(append(f.Refs, f.Keys...), "\n"), "policy not declared") {
		t.Errorf("strict load must not see fragment policies: %v", f.Refs)
	}

	cat, err := LoadWith(root, Options{Fragments: true})
	if err != nil {
		t.Fatal(err)
	}
	f = Check(cat, nil)
	report := strings.Join(append(f.Refs, f.Keys...), "\n")
	if strings.Contains(report, "frag.key") || strings.Contains(report, "frag-policy ") || !strings.Contains(report, "policy frag-unused is cited nowhere") {
		t.Errorf("merged keys and policies: %s", report)
	}
	cov := CheckCoverage(cat)
	probs := strings.Join(cov.Problems, "\n")
	for _, want := range []string{"done_pending easa/ratings/z is not pending", "fragments/coverage/b.yaml"} {
		if !strings.Contains(probs, want) {
			t.Errorf("coverage misses %q:\n%s", want, probs)
		}
	}
	if strings.Contains(probs, "does not list easa.rating.a#one") || strings.Contains(probs, "easa/ratings/a.yaml exists") || len(cov.Pending) != 1 || cov.Pending[0] != "easa:FCL.740.A: easa/ratings/b, easa/ratings/c" || cov.NotEvaluated != 1 {
		t.Errorf("merged coverage: %+v\n%s", cov, probs)
	}
}

func TestMergeCoverageDropsNote(t *testing.T) {
	root := t.TempDir()
	write(t, root, "f.yaml", "articles:\n  - { article: a, done_pending: [x/ratings/y] }\n")
	cov := &Coverage{Articles: []CoverageEntry{{Article: "a", Pending: []string{"x/ratings/y"}, Note: "n"}}}
	if p := mergeCoverage(cov, filepath.Join(root, "f.yaml")); len(p) != 0 || cov.Articles[0].Note != "" || len(cov.Articles[0].Pending) != 0 {
		t.Errorf("%v %+v", p, cov)
	}
	if p := mergeCoverage(cov, filepath.Join(root, "none.yaml")); len(p) != 1 {
		t.Errorf("missing file: %v", p)
	}
}

func TestLoadErrorsPoliciesAndFragments(t *testing.T) {
	root := t.TempDir()
	if _, err := LoadPolicies(root, false); err == nil {
		t.Error("no policies.yaml")
	}
	write(t, root, "policies.yaml", "policies: [\n")
	if _, err := LoadPolicies(root, false); err == nil {
		t.Error("bad policies.yaml")
	}
	write(t, root, "policies.yaml", "policies:\n  - { id: a }\n")
	write(t, root, "fragments/policies/x.yaml", "policies:\n  - { id: a }\n")
	if _, err := LoadPolicies(root, true); err == nil || !strings.Contains(err.Error(), "declared twice") {
		t.Errorf("duplicate: %v", err)
	}
	write(t, root, "fragments/keys", "a file, not a directory\n")
	if _, err := Fragments(root); err == nil {
		t.Error("fragments/keys is a file")
	}
	if _, err := LoadPolicies(root, false); err != nil {
		t.Error(err)
	}

	full := writeRoot(t, nil)
	write(t, full, "fragments/keys", "a file\n")
	if _, err := LoadWith(full, Options{Fragments: true}); err == nil {
		t.Error("keys fragments unreadable")
	}
	if err := os.Remove(filepath.Join(full, "fragments/keys")); err != nil {
		t.Fatal(err)
	}
	write(t, full, "fragments/policies", "a file\n")
	if _, err := LoadWith(full, Options{Fragments: true}); err == nil {
		t.Error("policy fragments unreadable")
	}
	if err := os.Remove(filepath.Join(full, "policies.yaml")); err != nil {
		t.Fatal(err)
	}
	if _, err := Load(full); err == nil {
		t.Error("no policies.yaml")
	}
	write(t, full, "fragments/coverage", "a file\n")
	write(t, full, "coverage/articles.yaml", "articles: []\n")
	if got := CheckCoverage(&Catalogue{Root: full, Options: Options{Fragments: true}}); len(got.Problems) != 1 || !strings.Contains(got.Problems[0], "fragments/coverage") {
		t.Error("coverage fragments unreadable")
	}
	write(t, full, "messages/keys.yaml", "keys: [\n")
	if _, err := Load(full); err == nil {
		t.Error("bad keys")
	}
}
