// Command evaluate loads the credential catalogue, evaluates one record on a date and prints
// every evaluation and, per held credential, the composite answer as JSON (DESIGN.md
// section 5): { "evaluations": [...], "credentials": [...] }.
//
//	go run ./cmd/evaluate -as-of 2026-09-28 record.yaml
//	go run ./cmd/evaluate -root path/to/open-aircrew-rules - < record.json
package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"os"
	"strings"
	"time"

	"gopkg.in/yaml.v3"

	"github.com/fjaeckel/open-aircrew-rules/credentials"
	"github.com/fjaeckel/open-aircrew-rules/engine"
)

func main() {
	os.Exit(run(os.Args[1:], os.Stdin, os.Stdout, os.Stderr))
}

func run(args []string, stdin io.Reader, stdout, stderr io.Writer) int {
	fs := flag.NewFlagSet("evaluate", flag.ContinueOnError)
	fs.SetOutput(stderr)
	root := fs.String("root", ".", "repository root (the directory holding vocabulary.yaml and credentials/)")
	asOf := fs.String("as-of", time.Now().UTC().Format(time.DateOnly), "date to evaluate on (YYYY-MM-DD)")
	only := fs.String("credential", "", "report only the evaluations of this credential id")
	fs.Usage = func() {
		fmt.Fprintln(stderr, "usage: evaluate [-root dir] [-as-of YYYY-MM-DD] [-credential id] <record.yaml|record.json|->")
		fs.PrintDefaults()
	}
	if err := fs.Parse(args); err != nil {
		return 2
	}
	if fs.NArg() != 1 {
		fs.Usage()
		return 2
	}
	date, err := engine.ParseDate(*asOf)
	if err != nil {
		fmt.Fprintf(stderr, "evaluate: -as-of: %v\n", err)
		return 2
	}
	var in []byte
	if fs.Arg(0) == "-" {
		in, err = io.ReadAll(stdin)
	} else {
		in, err = os.ReadFile(fs.Arg(0))
	}
	if err != nil {
		fmt.Fprintf(stderr, "evaluate: %v\n", err)
		return 2
	}
	var rec engine.Record
	if err := yaml.Unmarshal(in, &rec); err != nil {
		fmt.Fprintf(stderr, "evaluate: record: %v\n", err)
		return 2
	}
	cat, err := credentials.Load(*root)
	if err != nil {
		fmt.Fprintf(stderr, "evaluate: %v\n", err)
		return 2
	}
	if len(cat.Errors) > 0 {
		fmt.Fprintf(stderr, "evaluate: the catalogue does not compile:\n  %s\n", strings.Join(cat.Errors, "\n  "))
		return 2
	}
	res := cat.Evaluate(&rec, date)
	out := credentials.Result{Evaluations: []engine.Evaluation{}, Credentials: []credentials.Composite{}}
	for _, ev := range res.Evaluations {
		if *only == "" || belongsTo(cat, ev.RuleID, *only) {
			out.Evaluations = append(out.Evaluations, ev)
		}
	}
	for _, c := range res.Credentials {
		if *only == "" || c.Credential == *only {
			out.Credentials = append(out.Credentials, c)
		}
	}
	enc := json.NewEncoder(stdout)
	enc.SetIndent("", "  ")
	if err := enc.Encode(out); err != nil {
		fmt.Fprintf(stderr, "evaluate: %v\n", err)
		return 2
	}
	return 0
}

// belongsTo reports whether a compiled rule serves an evaluation of the credential.
func belongsTo(cat *credentials.Catalogue, ruleID, credential string) bool {
	c, ok := cat.Compiled[ruleID]
	if !ok {
		return false
	}
	for _, u := range c.Users {
		if strings.HasPrefix(u, credential+"#") {
			return true
		}
	}
	return false
}
