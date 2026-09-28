package credentials

import "testing"

func load(t *testing.T) *Catalogue {
	t.Helper()
	cat, err := Load("..")
	if err != nil {
		t.Fatal(err)
	}
	for _, e := range cat.Errors {
		t.Error(e)
	}
	return cat
}
