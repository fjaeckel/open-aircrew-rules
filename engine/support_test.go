package engine

import "testing"

func TestNewCatalogue(t *testing.T) {
	c := NewCatalogue("r", &Vocabulary{}, &Keys{}, []*Rule{{ID: "a"}, {ID: "a", File: "x"}, {ID: "b"}})
	if len(c.Rules) != 2 || len(c.Errors) != 1 {
		t.Fatalf("rules %d errors %v", len(c.Rules), c.Errors)
	}
	if _, ok := c.Rule("b"); !ok {
		t.Error("rule b")
	}
}

func TestDateJSON(t *testing.T) {
	var d Date
	if err := d.UnmarshalJSON([]byte(`"2026-09-28"`)); err != nil || d.String() != "2026-09-28" {
		t.Errorf("%v %v", d, err)
	}
	if b, _ := d.MarshalJSON(); string(b) != `"2026-09-28"` {
		t.Errorf("marshal %s", b)
	}
	if err := d.UnmarshalJSON([]byte(`12`)); err == nil {
		t.Error("a number is no date")
	}
	if err := d.UnmarshalJSON([]byte(`"2026-13-01"`)); err == nil {
		t.Error("month 13")
	}
}
