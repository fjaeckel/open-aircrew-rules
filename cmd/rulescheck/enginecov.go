package main

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"strings"
)

// engineCoverage runs the engine and credentials tests and returns their combined statement
// coverage.
func engineCoverage(root string) (float64, error) {
	tmp, err := os.CreateTemp("", "rulescheck-cover-*.out")
	if err != nil {
		return 0, err
	}
	tmp.Close()
	defer os.Remove(tmp.Name())
	cmd := exec.Command("go", "test", "-count=1", "-coverpkg=./engine/...,./credentials/...", "-coverprofile="+tmp.Name(), "./engine/...", "./credentials/...")
	cmd.Dir = root
	if out, err := cmd.CombinedOutput(); err != nil {
		return 0, fmt.Errorf("go test ./engine/... ./credentials/...: %v\n%s", err, out)
	}
	return profileCoverage(tmp.Name())
}

// profileCoverage computes covered statements over all blocks of a cover profile.
func profileCoverage(path string) (float64, error) {
	f, err := os.Open(path)
	if err != nil {
		return 0, err
	}
	defer f.Close()
	type block struct{ stmts, count int }
	blocks := map[string]block{}
	sc := bufio.NewScanner(f)
	for sc.Scan() {
		line := sc.Text()
		if strings.HasPrefix(line, "mode:") {
			continue
		}
		fields := strings.Fields(line)
		if len(fields) != 3 {
			continue
		}
		stmts, _ := strconv.Atoi(fields[1])
		count, _ := strconv.Atoi(fields[2])
		b := blocks[fields[0]]
		b.stmts = stmts
		b.count += count
		blocks[fields[0]] = b
	}
	total, covered := 0, 0
	for _, b := range blocks {
		total += b.stmts
		if b.count > 0 {
			covered += b.stmts
		}
	}
	if total == 0 {
		return 0, fmt.Errorf("empty cover profile")
	}
	return 100 * float64(covered) / float64(total), sc.Err()
}
