// Command rulescheck is the gate over the credential catalogue (DESIGN.md section 10): it
// validates every YAML file against schema/, loads and compiles the credentials, resolves
// every reference and policy, runs the worked examples, checks the coverage map of the
// articles in scope and the provenance of sources/, and measures engine and credentials
// statement coverage. It prints a grouped report and exits non-zero on any problem, unless
// -report is given. Files under fragments/ fail the default run; -fragments merges them in
// memory and lists them instead (DESIGN.md section 12).
//
//	go run ./cmd/rulescheck -strict              # CI gate
//	go run ./cmd/rulescheck -strict -fragments   # during parallel work, before integration
//	go run ./cmd/rulescheck -report              # the same report, always exits 0
package main

import (
	"flag"
	"fmt"
	"io"
	"os"
)

func main() {
	os.Exit(run(os.Args[1:], os.Stdout))
}

type options struct {
	root        string
	report      bool
	coverage    bool
	fragments   bool
	minCoverage float64
}

func run(args []string, out io.Writer) int {
	fs := flag.NewFlagSet("rulescheck", flag.ContinueOnError)
	fs.SetOutput(out)
	var o options
	strict := fs.Bool("strict", false, "fail on every problem (the default; spelled out for CI)")
	fs.StringVar(&o.root, "root", ".", "module root")
	fs.BoolVar(&o.report, "report", false, "print the report but exit 0")
	fs.BoolVar(&o.coverage, "coverage", true, "measure engine and credentials statement coverage with go test")
	fs.BoolVar(&o.fragments, "fragments", false, "merge fragments/ in memory and report them instead of failing on them")
	fs.Float64Var(&o.minCoverage, "min-coverage", 95, "minimum engine and credentials statement coverage in percent")
	if err := fs.Parse(args); err != nil {
		return 2
	}
	if *strict && o.report {
		fmt.Fprintln(out, "rulescheck: -strict and -report are exclusive")
		return 2
	}
	r, err := check(o)
	if err != nil {
		fmt.Fprintf(out, "rulescheck: %v\n", err)
		return 2
	}
	r.print(out, o)
	if r.problems() > 0 && !o.report {
		return 1
	}
	return 0
}
