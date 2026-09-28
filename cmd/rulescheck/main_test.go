package main

import (
	"bytes"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// TestGateGreen runs the strict gate over the module (without measuring engine coverage,
// which the engine's own tests and CI do) and prints the report with -v.
func TestGateGreen(t *testing.T) {
	var out bytes.Buffer
	if code := run([]string{"-strict", "-coverage=false", "-root", "../.."}, &out); code != 0 {
		t.Fatalf("exit %d\n%s", code, out.String())
	}
	for _, want := range []string{
		"== Schema validation: ok", "== Load and compile: ok", "== References resolve", "== Worked examples", "== Coverage: every article in scope",
		"== Sources: allowed origins only", "== Articles with pending credentials (report only)", "== Interpretations awaiting approval (report only)",
		"not measured (-coverage=false)", "\nOK\n",
	} {
		if !strings.Contains(out.String(), want) {
			t.Errorf("report misses %q", want)
		}
	}
	t.Log(out.String())
}

func TestFlags(t *testing.T) {
	var out bytes.Buffer
	if code := run([]string{"-strict", "-report"}, &out); code != 2 {
		t.Errorf("-strict -report: exit %d", code)
	}
	if code := run([]string{"-bogus"}, &out); code != 2 {
		t.Errorf("unknown flag: exit %d", code)
	}
	if code := run([]string{"-root", t.TempDir()}, &out); code != 2 {
		t.Errorf("no schemas: exit %d", code)
	}
}

// newRoot copies the vocabulary, the message keys, the schemas and one source text into a
// temporary module root.
func newRoot(t *testing.T) string {
	t.Helper()
	root := t.TempDir()
	copyFile(t, "../../vocabulary.yaml", filepath.Join(root, "vocabulary.yaml"))
	copyFile(t, "../../messages/keys.yaml", filepath.Join(root, "messages/keys.yaml"))
	copyFile(t, "../../policies.yaml", filepath.Join(root, "policies.yaml"))
	writeFile(t, filepath.Join(root, "associations.yaml"), "associations: []\n")
	copyFile(t, "../../sources/faa/61.57.md", filepath.Join(root, "sources/faa/61.57.md"))
	schemas, err := filepath.Glob("../../schema/*.schema.json")
	if err != nil || len(schemas) == 0 {
		t.Fatalf("schemas: %v", err)
	}
	for _, s := range schemas {
		copyFile(t, s, filepath.Join(root, "schema", filepath.Base(s)))
	}
	return root
}

// TestBroken checks that the gate reports the problems a contributor is likely to make.
func TestBroken(t *testing.T) {
	root := newRoot(t)
	vocab, err := os.ReadFile(filepath.Join(root, "vocabulary.yaml"))
	if err != nil {
		t.Fatal(err)
	}
	// A metric the engine does not implement.
	patched := strings.Replace(string(vocab), "\nmetrics:\n", "\nmetrics:\n  bogus_metric: { source: flights, aggregate: sum, reads: [], units: [flights], description: x }\n", 1)
	writeFile(t, filepath.Join(root, "vocabulary.yaml"), patched)
	writeFile(t, filepath.Join(root, "credentials/faa/ratings/x.yaml"), `
credential: X
id: faa.rating.x
kind: rating
authority: FAA
selects: { ratings: { classes: [SEP_LAND], authorities: [FAA] } }
evaluations:
  - id: passengers
    about: passengers
    asks: May the holder carry passengers?
    source: faa:61.57(a)
    counting: { within_days: 90, in_class: true }
    passes_if: { takeoffs: { min: 3, ref: faa:61.57(z)(9) } }
    outcomes: passengers
    wings: 2
`)
	writeFile(t, filepath.Join(root, "scope/articles-faa.yaml"), "articles:\n  - { cite: 14 CFR 61.57, url: \"https://www.ecfr.gov/x\", subject: passengers, summary: x, source_file: sources/faa/61.57.md }\n  - { cite: 14 CFR 61.58 }\n")
	writeFile(t, filepath.Join(root, "coverage/articles.yaml"), "articles:\n  - { article: \"faa:61.57\", pending: [faa/ratings/x], note: planned }\n")
	writeFile(t, filepath.Join(root, "sources/README.md"), "# Sources\n")
	writeFile(t, filepath.Join(root, "sources/faa/raw.xml"), "<xml/>\n")
	writeFile(t, filepath.Join(root, "sources/faa/no-origin.md"), "# X\n\n- Source URL: https://www.ecfr.gov/x\n\n---\n\ntext\n")
	writeFile(t, filepath.Join(root, "sources/de/amc.md"), "# AMC1 FCL.060\n\n- Origin: eu-legal-act\n- Attribution: none\n- URL: https://www.easa.europa.eu/document-library/easy-access-rules\n\n## Text\n\ncopied\n")
	writeFile(t, filepath.Join(root, "associations.yaml"), "associations:\n  - { id: dulv:ul-rule, publisher: DULV, title: X, delegated_by: \"faa:61.57\", verified: false }\n")
	writeFile(t, filepath.Join(root, "sources/de/ul-rule.md"), "# X\n\n- Origin: de-amtliches-werk\n- Attribution: § 5(1) UrhG\n\n## Text\n\ntext\n")
	writeFile(t, filepath.Join(root, "sources/de/dulv.md"), "# Rule\n\n- Origin: de-amtliches-werk\n- Attribution: § 5(1) UrhG\n\n## DULV Ausbildungsrichtlinie\n\ntext\n")
	var out bytes.Buffer
	if code := run([]string{"-root", root, "-coverage=false"}, &out); code != 1 {
		t.Fatalf("exit %d\n%s", code, out.String())
	}
	report := out.String()
	for _, want := range []string{
		"x.yaml: at /evaluations/0: validation failed", "missing properties 'url'", "article faa:61.58 of the scope is missing",
		"paragraph (z)(9) not found", "faa.rating.x#passengers: no passing example",
		"pending credentials/faa/ratings/x.yaml exists", "does not list faa.rating.x#passengers",
		"metrics (1): bogus_metric",
		"sources/faa/raw.xml: only Markdown texts", "sources/faa/no-origin.md: header has no allowed origin",
		"sources/de/amc.md: origin eu-legal-act belongs in sources/easa/", "sources/de/amc.md: header needs an \"- Attribution:\" line containing \"© European Union",
		"which is not a host of origin eu-legal-act", "names \"AMC\" in its header or a heading", "names \"DULV\" in its header or a heading",
		"sources/de/ul-rule.md: association document dulv:ul-rule is never stored", "association dulv:ul-rule is cited nowhere",
		"FAIL:",
	} {
		if !strings.Contains(report, want) {
			t.Errorf("report misses %q", want)
		}
	}
	out.Reset()
	if code := run([]string{"-root", root, "-coverage=false", "-report"}, &out); code != 0 || !strings.Contains(out.String(), "(-report: exit 0)") {
		t.Errorf("-report: exit %d", code)
	}
	if t.Failed() {
		t.Log(report)
	}
}

// TestBrokenVocabulary stops after schema validation when the vocabulary is invalid.
func TestBrokenVocabulary(t *testing.T) {
	root := newRoot(t)
	writeFile(t, filepath.Join(root, "vocabulary.yaml"), "version: 1\n")
	var out bytes.Buffer
	if code := run([]string{"-root", root, "-coverage=false"}, &out); code != 1 || !strings.Contains(out.String(), "missing properties") {
		t.Errorf("exit %d\n%s", code, out.String())
	}
}

// TestEngineCoverage measures coverage in a root without an engine (the gate reports that
// it could not measure) and with an unreachable minimum.
func TestEngineCoverage(t *testing.T) {
	root := newRoot(t)
	var out bytes.Buffer
	if code := run([]string{"-root", root}, &out); code != 1 || !strings.Contains(out.String(), "could not measure") {
		t.Errorf("exit %d\n%s", code, out.String())
	}
	if testing.Short() {
		t.Skip("runs the engine tests")
	}
	out.Reset()
	if code := run([]string{"-root", "../..", "-min-coverage", "101"}, &out); code != 1 || !strings.Contains(out.String(), "(minimum 101%)") {
		t.Errorf("exit %d\n%s", code, out.String())
	}
}

func TestProfileCoverage(t *testing.T) {
	p := filepath.Join(t.TempDir(), "c.out")
	writeFile(t, p, "mode: set\na.go:1.1,2.2 3 1\na.go:1.1,2.2 3 0\nb.go:1.1,2.2 1 0\nbad line\n")
	pct, err := profileCoverage(p)
	if err != nil || pct != 75 {
		t.Errorf("coverage %v %v", pct, err)
	}
	writeFile(t, p, "mode: set\n")
	if _, err := profileCoverage(p); err == nil {
		t.Error("empty profile must fail")
	}
	if _, err := profileCoverage(filepath.Join(t.TempDir(), "none")); err == nil {
		t.Error("missing profile must fail")
	}
}

func copyFile(t *testing.T, from, to string) {
	t.Helper()
	b, err := os.ReadFile(from)
	if err != nil {
		t.Fatal(err)
	}
	writeFile(t, to, string(b))
}

func writeFile(t *testing.T, path, body string) {
	t.Helper()
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(path, []byte(body), 0o644); err != nil {
		t.Fatal(err)
	}
}

// TestFragmentsFlag checks that fragments fail the default run and are validated and listed
// with -fragments.
func TestFragmentsFlag(t *testing.T) {
	root := newRoot(t)
	writeFile(t, filepath.Join(root, "coverage/articles.yaml"), "articles: []\n")
	writeFile(t, filepath.Join(root, "fragments/coverage/README.md"), "readme\n")
	writeFile(t, filepath.Join(root, "fragments/coverage/x.yaml"), "articles:\n  - { article: \"faa:61.57\", colour: red }\n")
	writeFile(t, filepath.Join(root, "fragments/keys/x.yaml"), "keys:\n  - { key: x.y, kind: message, params: [] }\n")
	writeFile(t, filepath.Join(root, "fragments/policies/x.yaml"), "policies: []\n")
	writeFile(t, filepath.Join(root, "fragments/changelog/x.md"), "- x\n")
	var out bytes.Buffer
	if code := run([]string{"-root", root, "-coverage=false"}, &out); code != 1 || !strings.Contains(out.String(), "fragments/keys/x.yaml: unmerged fragment") {
		t.Errorf("default: exit %d\n%s", code, out.String())
	}
	out.Reset()
	run([]string{"-root", root, "-coverage=false", "-fragments"}, &out)
	for _, want := range []string{"== No unmerged fragments: ok", "== Fragments pending integration (report only, -fragments): 4", "fragments/coverage/x.yaml: at /articles/0", "fragments/policies/x.yaml: at /policies"} {
		if !strings.Contains(out.String(), want) {
			t.Errorf("-fragments report misses %q\n%s", want, out.String())
		}
	}
	if err := os.RemoveAll(filepath.Join(root, "fragments/keys")); err != nil {
		t.Fatal(err)
	}
	writeFile(t, filepath.Join(root, "fragments/keys"), "not a directory\n")
	if code := run([]string{"-root", root, "-coverage=false"}, &out); code != 2 {
		t.Errorf("unreadable fragments: exit %d", code)
	}
}
