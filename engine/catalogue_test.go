package engine

import (
	"os"
	"path/filepath"
	"testing"
)

func TestLoadErrors(t *testing.T) {
	dir := t.TempDir()
	write := func(name, body string) string {
		p := filepath.Join(dir, name)
		if err := os.MkdirAll(filepath.Dir(p), 0o755); err != nil {
			t.Fatal(err)
		}
		if err := os.WriteFile(p, []byte(body), 0o644); err != nil {
			t.Fatal(err)
		}
		return p
	}
	if _, err := LoadVocabulary(write("bad.yaml", "version: [")); err == nil {
		t.Error("bad vocabulary yaml must fail")
	}
	if _, err := LoadKeys(write("k1.yaml", "keys: [")); err == nil {
		t.Error("bad keys yaml must fail")
	}
	if _, err := LoadKeys(write("k2.yaml", "keys:\n  - {key: a}\n  - {key: a}\n")); err == nil {
		t.Error("duplicate keys must fail")
	}
	if k, err := LoadKeys(write("k3.yaml", "keys:\n  - {key: a}\n")); err != nil {
		t.Error(err)
	} else if _, ok := k.Get("a"); !ok {
		t.Error("key a")
	}
	if k, err := LoadKeys(write("k4.yaml", "keys:\n  - {key: a}\n"), write("k5.yaml", "keys:\n  - {key: b}\n")); err != nil {
		t.Error(err)
	} else if d, ok := k.Get("a"); !ok || d.Key != "a" || len(k.Keys) != 2 {
		t.Error("keys from two files")
	}
	if _, err := LoadKeys(write("k6.yaml", "keys:\n  - {key: a}\n"), write("k7.yaml", "keys:\n  - {key: a}\n")); err == nil {
		t.Error("a key in two files must fail")
	}
	if _, err := LoadVocabulary(filepath.Join(dir, "nope.yaml")); err == nil {
		t.Error("missing vocabulary must fail")
	}
	if _, err := LoadKeys(filepath.Join(dir, "nope.yaml")); err == nil {
		t.Error("missing keys must fail")
	}
	write("tree/a.yaml", "")
	write("tree/sub/b.yaml", "")
	write("tree/c.txt", "")
	if files, err := YAMLFiles(filepath.Join(dir, "tree")); err != nil || len(files) != 2 {
		t.Errorf("yaml files: %v %v", files, err)
	}
	if files, err := YAMLFiles(filepath.Join(dir, "missing")); err != nil || len(files) != 0 {
		t.Errorf("missing dir: %v %v", files, err)
	}
}

func TestHoldings(t *testing.T) {
	c := testCatalogue(t)
	rec := record(t, "licences: [{ id: l, authority: EASA, type: PPL(A) }, { id: f, authority: FAA, type: PRIVATE }]\nratings: [{ id: r, licenceId: l, class: SEP_LAND }]\n")
	got := Holdings(c, rec, AppliesTo{Subject: "licence", Authorities: []string{"easa"}}, MustDate("2026-01-01"))
	if len(got) != 1 || got[0].ID != "l" || got[0].Detail != "PPL_A" {
		t.Errorf("licences: %+v", got)
	}
	if got := Holdings(c, rec, AppliesTo{Subject: "rating", Classes: []string{"SEP_LAND"}}, MustDate("2026-01-01")); len(got) != 1 || got[0].ID != "r" {
		t.Errorf("ratings: %+v", got)
	}
}
