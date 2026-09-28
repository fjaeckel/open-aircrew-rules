# Security Policy

## Reporting a vulnerability

Please report vulnerabilities in the evaluator or the tooling privately, through
[GitHub private vulnerability reporting](https://github.com/fjaeckel/open-aircrew-rules/security/advisories/new),
or by email to **hej@ninerlog.com**. Do not open a public issue.

Include a description, steps to reproduce, the impact, and a suggested fix if you have one. We
acknowledge reports within 48 hours.

## What is not a security issue

A rule that evaluates a regulation incorrectly is a correctness issue, not a vulnerability. Report it
publicly with the "This rule looks wrong" issue form, so that everyone relying on the data can see
it.

## Scope

- The Go evaluator and compiler
- The command-line tools (validation, gate, source fetching)
- Parsing of YAML catalogue and record input
