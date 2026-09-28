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
