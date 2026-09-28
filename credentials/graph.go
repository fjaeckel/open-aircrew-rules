package credentials

import (
	"errors"
	"fmt"
	"slices"
	"sort"
	"strings"
)

// errReported marks a resolve failure already reported by checkUses (a cycle or a layering
// violation); Load does not report it again for every evaluation it affects.
var errReported = errors.New("reported by the uses check")

// Requirement is one requires_all or requires_any edge between credentials.
type Requirement struct {
	From       string
	Evaluation string
	To         string
	Any        bool
}

// Use is one uses edge. From and To are evaluation nodes: "<credential id>#<evaluation id>"
// or a shared evaluation id.
type Use struct {
	From string
	To   string
}

// Dependencies is the requires and uses graph of the catalogue, sorted.
type Dependencies struct {
	Requires []Requirement
	Uses     []Use
}

// Dependencies returns the requires and uses edges of every credential and shared evaluation.
func (cat *Catalogue) Dependencies() Dependencies {
	var d Dependencies
	for _, c := range cat.Credentials {
		for _, e := range c.Evaluations {
			for _, to := range e.RequiresAll {
				d.Requires = append(d.Requires, Requirement{From: c.ID, Evaluation: e.ID, To: credentialOf(to)})
			}
			for _, to := range e.RequiresAny {
				d.Requires = append(d.Requires, Requirement{From: c.ID, Evaluation: e.ID, To: credentialOf(to), Any: true})
			}
		}
	}
	for from, to := range cat.usesEdges() {
		d.Uses = append(d.Uses, Use{From: from, To: to})
	}
	sort.Slice(d.Requires, func(i, j int) bool {
		a, b := d.Requires[i], d.Requires[j]
		if a.From != b.From {
			return a.From < b.From
		}
		if a.Evaluation != b.Evaluation {
			return a.Evaluation < b.Evaluation
		}
		return a.To < b.To
	})
	sort.Slice(d.Uses, func(i, j int) bool { return d.Uses[i].From < d.Uses[j].From })
	return d
}

// sharedUses returns what a shared evaluation's own `uses:` names ("" when none).
func sharedUses(s *Shared) string {
	for _, p := range pairs(&s.Evaluation) {
		if p[0].Value == "uses" {
			return p[1].Value
		}
	}
	return ""
}

// usesEdges maps every evaluation node with a `uses:` to what it names.
func (cat *Catalogue) usesEdges() map[string]string {
	next := map[string]string{}
	for _, c := range cat.Credentials {
		for _, e := range c.Evaluations {
			if e.Uses != "" {
				next[c.ID+"#"+e.ID] = e.Uses
			}
		}
	}
	for id, s := range cat.Shared {
		if u := sharedUses(s); u != "" {
			next[id] = u
		}
	}
	return next
}

// checkUses reports uses cycles with their path and shared evaluations that use a
// credential's evaluation (borrowing flows from credentials to shared evaluations, never
// back).
func (cat *Catalogue) checkUses() {
	next := cat.usesEdges()
	nodes := make([]string, 0, len(next))
	for n := range next {
		nodes = append(nodes, n)
	}
	sort.Strings(nodes)
	for _, n := range nodes {
		if _, shared := cat.Shared[n]; !shared {
			continue
		}
		switch u := next[n]; {
		case strings.HasPrefix(u, "$"):
			cat.Errors = append(cat.Errors, fmt.Sprintf("%s: uses %s: name a shared evaluation, not a parameter", n, u))
		case strings.Contains(u, "#"):
			cat.Errors = append(cat.Errors, fmt.Sprintf("%s: uses %s: a shared evaluation may use only shared evaluations, never a credential's evaluation", n, u))
		}
	}
	// Every node has at most one uses edge, so each cycle is found once by following them.
	const (
		open = 1
		done = 2
	)
	state := map[string]int{}
	for _, start := range nodes {
		var path []string
		n := start
		for state[n] == 0 {
			if _, ok := next[n]; !ok {
				break
			}
			state[n] = open
			path = append(path, n)
			n = next[n]
		}
		if state[n] == open {
			i := slices.Index(path, n)
			cat.Errors = append(cat.Errors, fmt.Sprintf("uses cycle: %s -> %s", strings.Join(path[i:], " -> "), n))
		}
		for _, p := range path {
			state[p] = done
		}
	}
}
