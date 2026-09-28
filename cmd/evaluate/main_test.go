package main

import (
	"bytes"
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/fjaeckel/open-aircrew-rules/credentials"
)

const record = `
licences:
  - { id: l-ppl, authority: EASA, type: PPL(A) }
ratings:
  - { id: r-sep, licenceId: l-ppl, class: SEP_LAND, expires: 2027-03-31 }
flights:
  - { date: 2026-09-01, class: SEP_LAND, minutes: { total: 60, pic: 60 }, takeoffs: { day: 3 }, landings: { day: 3 }, pilotFlying: true }
`

func TestEvaluate(t *testing.T) {
	var out, errs bytes.Buffer
	if code := run([]string{"-root", "../..", "-as-of", "2026-09-28", "-credential", "easa.rating.sep-land", "-"}, strings.NewReader(record), &out, &errs); code != 0 {
		t.Fatalf("exit %d: %s", code, errs.String())
	}
	var res credentials.Result
	if err := json.Unmarshal(out.Bytes(), &res); err != nil {
		t.Fatal(err)
	}
	if len(res.Credentials) != 1 || res.Credentials[0].Credential != "easa.rating.sep-land" || res.Credentials[0].Status != "unknown" || res.Credentials[0].DecidedBy.Evaluation != "licence" {
		t.Errorf("composites: %+v", res.Credentials)
	}
	got := map[string]string{}
	for _, ev := range res.Evaluations {
		if !strings.HasPrefix(ev.RuleID, "easa.rating.sep-land#") {
			t.Errorf("unexpected %s", ev.RuleID)
		}
		got[ev.RuleID] = ev.Status
	}
	if got["easa.rating.sep-land#passengers_day"] != "current" || got["easa.rating.sep-land#revalidation"] == "" {
		t.Errorf("evaluations: %v", got)
	}

	file := filepath.Join(t.TempDir(), "r.yaml")
	if err := os.WriteFile(file, []byte(record), 0o644); err != nil {
		t.Fatal(err)
	}
	out.Reset()
	if code := run([]string{"-root", "../..", "-as-of", "2026-09-28", file}, nil, &out, &errs); code != 0 || !strings.Contains(out.String(), "easa.licence.ppl-a#") {
		t.Errorf("all credentials: exit %d", code)
	}
}

func TestEvaluateErrors(t *testing.T) {
	bad := t.TempDir()
	if err := os.WriteFile(filepath.Join(bad, "vocabulary.yaml"), []byte("version: 1\ncredential_vocabulary: {}\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	for _, tc := range []struct {
		name  string
		args  []string
		stdin string
		want  string
	}{
		{"no args", nil, "", "usage"},
		{"bad flag", []string{"-nope", "x"}, "", "not defined"},
		{"bad date", []string{"-as-of", "2026-13-01", "-"}, "", "-as-of"},
		{"missing file", []string{filepath.Join(bad, "none.yaml")}, "", "no such file"},
		{"bad record", []string{"-"}, "licences: [", "record"},
		{"no catalogue", []string{"-root", t.TempDir(), "-"}, "{}", "vocabulary.yaml"},
		{"broken catalogue", []string{"-root", bad, "-"}, "{}", "evaluate:"},
	} {
		var out, errs bytes.Buffer
		if code := run(tc.args, strings.NewReader(tc.stdin), &out, &errs); code != 2 || !strings.Contains(errs.String(), tc.want) {
			t.Errorf("%s: exit %d, stderr %q", tc.name, code, errs.String())
		}
	}
	if belongsTo(&credentials.Catalogue{}, "x#y", "x") {
		t.Error("a rule the catalogue does not compile belongs to no credential")
	}
}
