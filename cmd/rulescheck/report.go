package main

import (
	"fmt"
	"io"
	"sort"
	"strings"

	"github.com/fjaeckel/open-aircrew-rules/credentials"
)

// report collects every finding, grouped for printing.
type report struct {
	files        int
	shared       int
	schema       []string
	f            *credentials.Findings
	cov          *credentials.CoverageResult
	vocabUnimpl  map[string][]string
	sources      sourcesResult
	coverage     string
	coverageFail bool
	frags        map[string][]string
	fragFail     []string
}

func (r *report) problems() int {
	n := len(r.schema) + len(r.sources.problems) + count(r.vocabUnimpl) + len(r.fragFail)
	if r.f != nil {
		n += r.f.Problems()
	}
	if r.cov != nil {
		n += len(r.cov.Problems)
	}
	if r.coverageFail {
		n++
	}
	return n
}

func (r *report) print(w io.Writer, o options) {
	if r.f == nil {
		fmt.Fprintf(w, "rulescheck: %d files validated\n", r.files)
		section(w, "Schema validation", len(r.schema))
		list(w, r.schema)
		summary(w, r, o)
		return
	}
	f := r.f
	fmt.Fprintf(w, "rulescheck: %d credential files, %d shared evaluations, %d evaluations (%d compiled rules), %d worked examples, %d references, %d files validated\n",
		f.Credentials, r.shared, f.Evaluations, f.Compiled, f.Examples, f.RefsChecked, r.files)
	section(w, "Schema validation", len(r.schema))
	list(w, r.schema)
	section(w, "Load and compile", len(f.Load)+len(f.Files))
	list(w, f.Load)
	list(w, f.Files)
	section(w, "References resolve to a paragraph of sources/ or a declared policy", len(f.Refs))
	list(w, f.Refs)
	section(w, "Message keys and statuses", len(f.Keys))
	list(w, f.Keys)
	section(w, "Worked examples (every evaluation has a passing and a failing one)", len(f.ExampleErrs))
	list(w, f.ExampleErrs)
	section(w, "No two credentials select the same record item", len(f.Overlaps))
	list(w, f.Overlaps)
	section(w, "Interpretations are well formed", len(f.Interpret))
	list(w, f.Interpret)
	section(w, "Vocabulary entries the engine implements", count(r.vocabUnimpl))
	groups(w, r.vocabUnimpl)
	c := r.cov
	section(w, "Coverage: every article in scope is evaluated, pending or not evaluated with a reason", len(c.Problems))
	list(w, c.Problems)
	fmt.Fprintf(w, "  %d articles: %d evaluated, %d pending, %d not evaluated\n", c.Evaluated+c.PendingOnly+c.NotEvaluated, c.Evaluated, c.PendingOnly, c.NotEvaluated)
	r.printSources(w)
	fail := 0
	if r.coverageFail {
		fail = 1
	}
	section(w, "Engine and credentials statement coverage", fail)
	fmt.Fprintf(w, "  %s\n", r.coverage)
	section(w, "No unmerged fragments", len(r.fragFail))
	list(w, r.fragFail)

	auths := keys(c.PendingByAuthority)
	per := make([]string, 0, len(auths))
	for _, a := range auths {
		per = append(per, fmt.Sprintf("%s %d", a, c.PendingByAuthority[a]))
	}
	fmt.Fprintf(w, "\n== Articles with pending credentials (report only): %d (%s)\n", len(c.Pending), strings.Join(per, ", "))
	list(w, c.Pending)
	fmt.Fprintf(w, "\n== Interpretations awaiting approval (report only): %d\n", len(f.Unapproved))
	byFile := map[string][]string{}
	for _, u := range f.Unapproved {
		file, rest, _ := strings.Cut(u, "#")
		id, _, _ := strings.Cut(rest, ":")
		byFile[file] = append(byFile[file], id)
	}
	for _, k := range keys(byFile) {
		fmt.Fprintf(w, "  %s (%d): %s\n", k, len(byFile[k]), strings.Join(byFile[k], ", "))
	}
	if o.fragments {
		fmt.Fprintf(w, "\n== Fragments pending integration (report only, -fragments): %d\n", count(r.frags))
		groups(w, r.frags)
	}
	summary(w, r, o)
}

func (r *report) printSources(w io.Writer) {
	origins := make([]string, 0, len(r.sources.origins))
	for o, n := range r.sources.origins {
		origins = append(origins, fmt.Sprintf("%s %d", o, n))
	}
	sort.Strings(origins)
	section(w, "Sources: allowed origins only ("+strings.Join(origins, ", ")+")", len(r.sources.problems))
	list(w, r.sources.problems)
}

func summary(w io.Writer, r *report, o options) {
	if n := r.problems(); n > 0 {
		mode := ""
		if o.report {
			mode = " (-report: exit 0)"
		}
		fmt.Fprintf(w, "\nFAIL: %d problems%s\n", n, mode)
		return
	}
	fmt.Fprintln(w, "\nOK")
}

func section(w io.Writer, title string, problems int) {
	state := "ok"
	if problems > 0 {
		state = fmt.Sprintf("%d problems", problems)
	}
	fmt.Fprintf(w, "\n== %s: %s\n", title, state)
}

func list(w io.Writer, l []string) {
	for _, s := range l {
		fmt.Fprintf(w, "  - %s\n", s)
	}
}

func groups(w io.Writer, m map[string][]string) {
	for _, k := range keys(m) {
		if len(m[k]) > 0 {
			fmt.Fprintf(w, "  %s (%d): %s\n", k, len(m[k]), strings.Join(m[k], ", "))
		}
	}
}

func count(m map[string][]string) int {
	n := 0
	for _, l := range m {
		n += len(l)
	}
	return n
}
