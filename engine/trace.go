package engine

import (
	"slices"
	"sort"
)

// Trace records what one evaluation exercised: the stage reached, the state of each
// requirement, the any_of branches met, window edges, restoring events, unknown inputs and
// where asOf fell relative to the expiry date. Tests use it to show that a feature ran.
type Trace struct {
	Stage        string
	Requirements map[string]string
	AnyOf        []string
	NOf          map[string]bool
	Edges        map[string]bool
	Events       []string
	Unknown      []string
	Expiry       string
}

func newTrace() *Trace {
	return &Trace{Requirements: map[string]string{}, NOf: map[string]bool{}, Edges: map[string]bool{}}
}

func (t *Trace) leaf(lf *leaf, st leafState, v *Vocabulary) {
	switch {
	case st.met:
		t.Requirements[lf.node.ID] = "met"
	case st.tracked:
		t.Requirements[lf.node.ID] = "unmet"
	default:
		t.Requirements[lf.node.ID] = "untracked"
	}
	for _, u := range st.unknown {
		if u == lf.node.Metric && st.tracked {
			continue
		}
		if !slices.Contains(t.Unknown, u) {
			t.Unknown = append(t.Unknown, u)
		}
	}
	if !st.inWindow || st.span.open || lf.window == nil {
		return
	}
	for _, it := range lf.items {
		switch {
		case it.date.Equal(st.span.from):
			t.Edges[lf.window.String()+":edge-in"] = true
		case it.date.Equal(st.span.from.AddDays(-1)):
			t.Edges[lf.window.String()+":edge-out"] = true
		}
	}
}

// edgeWindows returns the distinct bounded windows used by a rule's leaves.
func edgeWindows(r *Rule) []string {
	var out []string
	var walk func(n *Node, w *Window)
	walk = func(n *Node, w *Window) {
		if n == nil {
			return
		}
		if n.Window != nil {
			w = n.Window
		}
		if n.IsLeaf() {
			if w != nil && w.Kind != "lifetime" && !slices.Contains(out, w.String()) {
				out = append(out, w.String())
			}
			return
		}
		for _, c := range n.Children() {
			walk(c, w)
		}
	}
	walk(r.Requirements, r.Window)
	sort.Strings(out)
	return out
}

func windowTag(windows []string, w, edge string) string {
	if len(windows) == 1 {
		return "window:" + edge
	}
	return "window:" + w + ":" + edge
}

// ObservedTags returns what one evaluation trace exercised as sorted tags
// ("stage:<id>", "requirement:<id>:met|unmet|untracked", "any_of:<node>:<branch>",
// "window:edge-in", "event:<kind>", "unknown:<input>", "expiry:before|on|after").
func ObservedTags(r *Rule, t *Trace) []string {
	set := map[string]bool{}
	if t.Stage != "" {
		set["stage:"+t.Stage] = true
	}
	for id, st := range t.Requirements {
		set["requirement:"+id+":"+st] = true
	}
	for _, a := range t.AnyOf {
		set["any_of:"+a] = true
	}
	for id, met := range t.NOf {
		if met {
			set["n_of:"+id+":met"] = true
		} else {
			set["n_of:"+id+":unmet"] = true
		}
	}
	ws := edgeWindows(r)
	for e := range t.Edges {
		for _, w := range ws {
			if e == w+":edge-in" {
				set[windowTag(ws, w, "edge-in")] = true
			}
			if e == w+":edge-out" {
				set[windowTag(ws, w, "edge-out")] = true
			}
		}
	}
	for _, e := range t.Events {
		set["event:"+e] = true
	}
	for _, u := range t.Unknown {
		set["unknown:"+u] = true
	}
	if t.Expiry != "" {
		set["expiry:"+t.Expiry] = true
	}
	return sortedKeys(set)
}

func sortedKeys(m map[string]bool) []string {
	out := make([]string, 0, len(m))
	for k := range m {
		out = append(out, k)
	}
	sort.Strings(out)
	return out
}
