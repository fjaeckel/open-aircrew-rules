# open-aircrew-rules

Machine-readable pilot licence, rating and medical rules. Every credential a pilot can hold
(a licence, a class or type rating, a privilege, an endorsement, an instructor or examiner
certificate, a medical) is one YAML file that lists every evaluation it needs: recency,
revalidation, passenger carriage, validity. Each condition carries a reference to the exact
EASA, FAA or German paragraph behind it (or, for rules an association sets under a German
statute, to that association document by title). A Go evaluator runs the files against a pilot's
logbook record.

**Not legal advice.** The catalogue can be wrong, incomplete or out of date. It does not
replace the official text, your authority, your examiner or your instructor; see
[NOTICE](NOTICE).

## Who it is for

- **Logbook and flight-school software** that wants to show "current", "expiring" or
  "expired" with the reason and the paragraph behind it, without writing the rules itself.
- **Pilots, instructors and examiners** who want to read, check or challenge how a rule is
  applied: one file per credential, plain words, a reference on every condition.
- **Anyone tracking regulation changes**: every reading of the text is a named
  interpretation awaiting sign-off, and every article in scope is accounted for.

## A credential in 25 lines

```yaml
credential: SEP (land) class rating
id: easa.rating.sep-land
kind: rating
authority: EASA
selects:
  ratings: { classes: [SEP_LAND], authorities: [EASA, LBA], not_licence_kinds: [LAPL_A, SPL, LAPL_S, UL] }
evaluations:
  - id: revalidation
    asks: Is the rating valid, and has the holder done what keeps it valid?
    source: easa:FCL.740.A(b)(1)
    counting: { within_months_before_expiry: 12 }
    passes_if:
      ref: easa:FCL.740.A(b)(1)
      any_of:
        - proficiency_check: { in_class: true, within_months_before_expiry: 3, ref: easa:FCL.740.A(b)(1)(i) }
        - id: experience
          ref: easa:FCL.740.A(b)(1)(ii)
          all_of:
            - flight_time: { min_hours: 12, in_class: true, ref: easa:FCL.740.A(b)(1)(ii) }
            - takeoffs: { min: 12, in_class: true, ref: easa:FCL.740.A(b)(1)(ii)(B) }
    outcomes: { preset: revalidation, expiring_notice: { days: 90, ref: policy:expiring-notice } }
  - id: passengers_day
    uses: easa.shared.passengers-day
    with: { type_rated: false }
```

Shortened from [credentials/easa/ratings/sep-land.yaml](credentials/easa/ratings/sep-land.yaml),
which also has the PIC hours, landings, refresher training, the ultralight credit, night
passengers and nine interpretations. The format is described in
[docs/credential-format.md](docs/credential-format.md).

## Evaluating a record

The input is a neutral record of licences, ratings, privileges, certificates, events and
flights ([schema/record.schema.json](schema/record.schema.json)); the output is one
evaluation per credential evaluation and subject, with a status, a message key, the
requirement rows and the citations, and, per credential the record holds, a composite
answer to "may it be exercised today?" that combines its own evaluations with the
credentials it requires (`requires_all`, `requires_any`) and names what decides it.

```go
import (
	"github.com/fjaeckel/open-aircrew-rules/credentials"
	"github.com/fjaeckel/open-aircrew-rules/engine"
)

cat, err := credentials.Load("path/to/open-aircrew-rules") // compiles every credential
if err != nil {
	return err
}
res := cat.Evaluate(&record, engine.MustDate("2026-09-28"))
for _, ev := range res.Evaluations {
	fmt.Println(ev.RuleID, ev.Status, ev.MessageKey, ev.Citations)
	// easa.rating.sep-land#revalidation current rating.revalidation_current [EASA FCL.740.A(b)(1) ...]
}
for _, c := range res.Credentials {
	fmt.Println(c.Credential, c.Subject.ID, c.Status) // c.DecidedBy names the deciding member
	// easa.licence.ppl-a l-ppl expired
}
```

From the command line, with a record in YAML or JSON (the output is
`{ "evaluations": [...], "credentials": [...] }`):

```bash
go run ./cmd/evaluate -as-of 2026-09-28 record.yaml
go run ./cmd/evaluate -as-of 2026-09-28 -credential easa.rating.sep-land - < record.json
```

Missing input is unknown, never zero: a requirement that cannot be counted from the record
is reported untracked and never met, and the status says what to record.

## Repository layout

| Path | What |
| --- | --- |
| `credentials/<authority>/<kind>/<name>.yaml` | One credential and every evaluation it needs |
| `credentials/<authority>/shared/<name>.yaml` | Parameterised evaluations several credentials use |
| `examples/<credential id>.yaml` | Worked examples with their expected outcome: at least one passing and one failing per evaluation |
| `scope/articles-<authority>.yaml` | Every article in scope, with a one-line summary |
| `coverage/articles.yaml` | How each article is accounted for: evaluated, pending or not evaluated |
| `sources/` | Verbatim regulation texts of allowed origin ([sources/README.md](sources/README.md)) |
| `vocabulary.yaml`, `messages/keys.yaml` | The closed vocabulary and the message keys |
| `policies.yaml` | Every convention with no legal text behind it, with its rationale |
| `associations.yaml` | Association documents cited by title (`assoc:` refs); their text is never stored |
| `fragments/` | Unmerged changes to the shared files during parallel work ([fragments/README.md](fragments/README.md)) |
| `schema/` | JSON Schema for every YAML file |
| `credentials/*.go`, `engine/` | Loader and compiler; the evaluator |
| `cmd/evaluate`, `cmd/rulescheck`, `cmd/rulesgen` | Command-line evaluator, gate, generator |
| `gen/` | Generated constants and example tests (never edited by hand) |
| `docs/` | Format reference, [conversion recipe](docs/converting.md), own-words AMC/GM notes, [ADRs](docs/adr/) |

[DESIGN.md](DESIGN.md) is the contract the files follow.

## Status

On 2026-09-28:

| Authority | Licences | Ratings | Privileges | Endorsements | Instructor / examiner | Medicals | Shared evaluations | Articles in scope | Evaluated | Pending | Not evaluated |
| --- | --- | --- | --- | --- | --- | --- | --- | --- | --- | --- | --- |
| EASA (Part-FCL, Part-SFCL, Part-BFCL, Part-MED) | 14 | 15 | 8 | 1 | 10 / 2 | 3 | 6 | 66 | 42 | 0 | 24 |
| FAA (14 CFR Parts 61, 68) | 11 | 1 | 3 | 0 | 2 / 0 | 4 | 7 | 32 | 19 | 0 | 13 |
| Germany (LuftPersV, LuftVZO) | 1 | 6 | 2 | 0 | 1 / 0 | 0 | 1 | 15 | 7 | 0 | 8 |
| **Total** | **26** | **22** | **13** | **1** | **13 / 2** | **7** | **14** | **113** | **68** | **0** | **45** |

Every article in scope is evaluated or says why it is not (for example, not a currency rule,
or the record cannot show what it needs); no credential file is pending. The 84 credentials
hold 280 evaluations (189 compiled rules and 81 requirement entries; evaluations reused
through `uses:` share a rule), 792 worked examples (66 of them composite), 1591 references
(policies and association documents included), 147 verbatim source texts (87 EU legal acts,
40 US federal regulations, 20 German statutes and ordinances), 10 policies, 6 association
documents (none verified) and 407 interpretations, none signed off yet. See
[CHANGELOG.md](CHANGELOG.md) for the list.

## Review status

Every interpretation has had a first-pass review against the verbatim text under `sources/`.
A review is not a sign-off: **no interpretation is approved yet** (`approved_by` and
`approved_on` are empty everywhere). Each review gives, per interpretation, a verdict:
**supported** (the cited paragraph says what the reading says), **needs-decision** (the
text or the record leaves real room; the owner picks an option) or **fixed** (the reading
or the result contradicted the text and was corrected, by the review or at integration).

| Authority | Interpretations reviewed | Supported | Needs decision | Fixed | Reviews |
| --- | --- | --- | --- | --- | --- |
| EASA | 251 | 196 | 38 | 18 | [medicals, examiners, language](reviews/easa-medicals-examiners-language.md), [shared and privileges](reviews/easa-shared-privileges.md), [class ratings, IR, BIR](reviews/easa-class-ir-ratings.md), [type ratings and licences](reviews/easa-types-licences.md), [SPL, GPL, instructors](reviews/easa-spl-gpl-instructors.md) |
| FAA | 91 | 76 | 11 | 4 | [instructors, privileges, ratings, medicals](reviews/faa-instructors-privileges-ratings-medicals.md), [shared and pilot certificates](reviews/faa-shared-licences.md) |
| Germany | 65 | 37 | 27 | 1 | [licence, privileges, instructor](reviews/de-ul-licence-privileges-instructor.md), [ultralight kind ratings](reviews/de-ul-ratings.md) |
| **Total** | **407** | **309** | **76** | **23** | [reviews/README.md](reviews/README.md) |

The EASA verdicts add up to 252: `easa.rating.powered-lift-type` / `part-fcl-type-ratings`
was fixed in one review and found supported in another, and is counted in both.
[reviews/README.md](reviews/README.md) lists the 76 needs-decision items in one table, the
owner's decision list.

## Keeping up with regulation

- **A regulation changed?** Open a "A regulation changed" issue with the amending act, the
  effective date and the official source. The fix is a pull request that updates the
  verbatim text under `sources/` and ends the affected evaluations with `effective_to`,
  adding their successors with `effective_from`.
- **A rule looks wrong?** Open a "This rule looks wrong" issue naming the credential,
  evaluation and paragraph.
- **A credential is missing?** Open an "Add a credential" issue, or write it following
  [CONTRIBUTING.md](CONTRIBUTING.md); `coverage/articles.yaml` says why each article not
  evaluated is left out.
- **Interpretations** are signed off by a reviewer qualified for the credential, through
  pull-request review, which sets `approved_by` and `approved_on`.

Details: [CONTRIBUTING.md](CONTRIBUTING.md). Every change runs the gate:

```bash
go vet ./... && go test -race ./...
go generate ./... && git diff --exit-code   # generated code is up to date
go run ./cmd/rulescheck -strict              # the gate
```

## Licence

Code, catalogue, examples and tools are MIT licensed ([LICENSE](LICENSE)). The texts under
`sources/` are not relicensed: US federal works are in the public domain, German statutes
are official works, and EU legal acts are reused under Commission Decision 2011/833/EU with
attribution, as [NOTICE](NOTICE) and [sources/README.md](sources/README.md) describe. EASA
AMC/GM, ICAO and association documents are never stored.

## Provenance

Started inside the NinerLog logbook project and extracted as a standalone package.
