# Contributing

Thank you for helping keep pilot rules correct. Most contributions are one of three things:
adding or changing a credential, following a regulation change, or challenging how a rule
reads the text. Code changes to the engine, compiler or tools are welcome too. By
contributing you agree that your contribution is licensed under the MIT License (see
[LICENSE](LICENSE)) and follow the [Code of Conduct](CODE_OF_CONDUCT.md).

Read [DESIGN.md](DESIGN.md) once, and [docs/credential-format.md](docs/credential-format.md)
before writing a credential file. Converting a credential from an article-centric rule set
follows [docs/converting.md](docs/converting.md).

## Adding or changing a credential

1. **Find its place.** A credential lives at
   `credentials/<authority>/<kind dir>/<name>.yaml` (kind dirs: `licences`, `ratings`,
   `privileges`, `endorsements`, `instructors`, `examiners`, `medicals`) with the id
   `<authority>.<kind>.<name>`. `coverage/articles.yaml` names the files still pending;
   pick one of those names where it fits.
2. **Store the text first.** Every article you reference must exist verbatim under
   `sources/<authority>/` (see [Sources and copyright](#sources-and-copyright)).
3. **Write the evaluations.** Each evaluation asks one question in plain words (`asks`),
   names its article (`source`) and uses only the words of `vocabulary.yaml`
   `credential_vocabulary`. Reuse shared evaluations with `uses:`, and name the credentials
   it depends on with `requires_all:` (every one) and `requires_any:` (at least one); a
   rating names its licence, a licence its medical, never the other way round.
4. **A `ref:` on every condition.** Every count, combinator, `waived_by`, `only_if`,
   `ul_credit`, `relevant_class`, validity period, `restored_by` event, `on_fail`, written
   stage and expiring notice points to the most specific paragraph it encodes, e.g.
   `easa:FCL.740.A(b)(1)(ii)(C)`, `faa:61.57(c)(1)(iii)`, `de:LuftPersV.45a`. The gate
   checks that the exact paragraph exists in the stored text.
5. **`policy:` for anything that is not law.** A presentation or implementation choice (the
   90-day expiring notice, "unknown" when input is missing, the status an unmet recency is
   reported with) is `ref: policy:<id>`, with the id, a statement and a rationale declared
   once in `policies.yaml`. Reuse an existing policy where it says the same; keep new ones
   few, as a reviewer will ask why each one is not law.
6. **Interpretations for every judgement call.** Where the text leaves room (what a word
   means, what is assumed because a logbook cannot show it, what part of an article is not
   applied), add an interpretation with `id`, `reading`, `ref`, `affects`, and
   `approved_by: null`, `approved_on: null`. No free-text notes.
7. **Worked examples.** In `examples/<credential id>.yaml`, give every evaluation at least
   one passing and one failing example, each stating its `outcome` (the status it must
   report; the outcome classes are defined in the format reference), and add an example for
   each interpretation that changes a result (`shows: [interpretation id]`). For a licence
   with requirements, add composite examples (`composite: true`). Keep records small and
   say in `says:` what the example shows.
8. **Account for the articles.** List each new evaluation under the article of its `source`
   in `coverage/articles.yaml`, and remove the credential file from the `pending` lists it
   appeared in.
9. **Run the gate** (below) and add a line to `CHANGELOG.md` under "Unreleased".

When several contributors work at once, steps 8 and 9 and any new message key or policy go
into fragments instead of the shared files (`fragments/`, see its README), and the gate
runs with `-fragments`.

Before changing a shared evaluation (`credentials/<authority>/shared/`) or an evaluation
other credentials borrow with `uses:`, look it up in
[docs/dependencies.md](docs/dependencies.md): every credential listed there changes with it,
and its worked examples must still pass. `go generate ./...` rewrites that file; commit it
with your change.

A new vocabulary word or message key is a larger change: add it to `vocabulary.yaml` or
`messages/keys.yaml`, implement it in the compiler or engine, add a schema entry if needed,
and show it in a worked example.

## Interpretations and sign-off

Every new or changed interpretation follows the interpretation principles of DESIGN.md
section 14 (P1 missing data is unknown, P2 only a check for this credential counts, P3 count
within the validity period, P4 devices only where the text permits, P5 presentation-only
rows keep the current choice). It names the principle it follows with `principle: P1` ...
`P5`, or none when the wording of the text decides alone; a reading that departs from a
principle says why.


An interpretation is approved when a qualified reviewer agrees that the reading is the right
one. Qualified means one of: a current holder of the credential or a higher one that
includes it, an instructor or examiner for it, or someone who confirmed the reading with the
competent authority (say so in the review).

Sign-off happens in pull-request review, never by the author of the reading:

1. The reviewer approves the pull request with a comment naming the interpretations they
   sign off and their basis.
2. In the same pull request (or a follow-up), `approved_by` is set to the reviewer's GitHub
   handle and `approved_on` to the date of the review (`YYYY-MM-DD`). The two always go
   together; the gate fails on one without the other.
3. Changing the `reading` of an approved interpretation resets both fields to `null`.

The maintainers (see `.github/CODEOWNERS`) merge. Unapproved interpretations do not block a
merge; the gate lists them so they stay visible.

## When a regulation changes

1. Open a "A regulation changed" issue: the amending act, the articles, the effective date
   and the official source.
2. The pull request replaces the verbatim text under `sources/` with the new consolidated
   version, updating the header's retrieval date and consolidation.
3. Every affected evaluation gets `effective_to: <day before the change>`, and its
   successor, a new evaluation with a new id (e.g. `revalidation_2027`), gets
   `effective_from: <effective date>`. Both keep valid references: if the old paragraph
   disappears from the text, keep the old evaluation's refs pointing to what still exists
   and say so in an interpretation, or end it on the change date.
4. Update the worked examples on both sides of the effective date, the coverage map and
   `CHANGELOG.md`.

## Sources and copyright

`sources/` holds only texts that may be reused, and the gate enforces it:

- **Allowed**: US federal regulations from the eCFR (17 U.S.C. § 105), German statutes and
  ordinances from gesetze-im-internet.de (§ 5(1) UrhG), EU legal acts from EUR-Lex (reuse
  under Commission Decision 2011/833/EU with the attribution "© European Union,
  https://eur-lex.europa.eu").
- **Never**: EASA AMC, GM or Easy Access Rules, ICAO documents, association rules (DULV,
  DAeC) or anything else under copyright. Summarise AMC/GM in your own words in
  `docs/amc-gm-notes.md` and name the item; declare an association rule in
  `associations.yaml` and cite it as `assoc:<id>` next to the statute delegating it.
- Store Markdown only, copied from the official source without retyping or paraphrasing,
  with the header lines of the existing files (`Origin`, `Attribution`, URL, retrieval date,
  consolidation). See [sources/README.md](sources/README.md).
- Credential files never quote; they reference paragraphs.

## Running the gate locally

Go 1.27 or later:

```bash
go vet ./...
go test -race ./...
go generate ./... && git diff --exit-code   # generated code is up to date
go run ./cmd/rulescheck -strict              # the gate CI runs
go run ./cmd/rulescheck -strict -fragments   # the same, with unmerged fragments/ (parallel work)
go run ./cmd/rulescheck -report              # the same report, always exits 0
gofmt -l .                                    # prints nothing
```

`go run ./cmd/evaluate -as-of <date> record.yaml` evaluates a record, which helps when
writing examples.

## Commits and pull requests

- [Conventional Commits](https://www.conventionalcommits.org/): `feat(easa): add MEP land
  class rating`, `fix(faa): count full-stop landings for tailwheel passengers`,
  `docs: explain policy refs`, `chore(deps): ...`. Use the authority as the scope for
  catalogue changes.
- One credential or one regulation change per pull request where possible.
- Fill in the pull-request template: what changes, the legal basis, and the interpretations.
