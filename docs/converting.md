# Converting a credential from the legacy rule set

The catalogue started from an article-centric rule set: one file per regulation article
with its logic verified against the text, and cases with expected results. This recipe turns
those rules into credential files ([credential-format.md](credential-format.md)) and proves
the conversion equivalent without running the old engine: **every legacy case becomes a
worked example whose `outcome` equals the case's expected status**.

## The legacy source (read only)

In a checkout of the legacy rule set, under `rules/`:

| Path | What it holds | Carry over |
| --- | --- | --- |
| `catalogue/<authority>/<instrument>/<rule id>.yaml` | one rule per article or point: `applies_to`, `window`, `requirements`, `stages`, `validity`, `restored_by`, `source.quote`, `notes`, `divergences`, `support` | the article logic, the refs, the judgement calls |
| `cases/<rule id>/<case>.yaml` | `record`, `asOf`, `expect: [{ subject, status, messageKey, ... }]` | every case, as a worked example |
| `migration.yaml` | legacy rules already converted, and the subjects still `remaining` for later credentials | skip what is converted |

Never copy:

- **`divergences`**: they describe how another application behaves, not the text. They are
  not requirements, not interpretations and not examples.
- **`source.quote`**: credential files do not quote; they reference paragraphs, and the text
  lives in `sources/`.
- **Product or application names** in notes, keys or readings: the catalogue is neutral.
- **Free-text `notes`** as notes: every judgement call in them becomes an interpretation
  (below); the rest is dropped.

## 1. Pick the file and its name

1. Find the legacy rules for the credential and the pending name in
   `coverage/articles.yaml` (e.g. `easa/ratings/mep-land`). The file is
   `credentials/<authority>/<kind dir>/<name>.yaml`, the id `<authority>.<kind>.<name>`.
2. Names are what a pilot calls the credential, lower-case and dashed: `mep-land`, `ir-a`,
   `lapl-a`, `aeroplane-type`, `fi`. Where Part-FCL and Part-SFCL both have a privilege of
   the same name, prefix it: `fcl-sailplane-towing` and `sfcl-sailplane-towing`. If a
   pending name does not fit this, use the better name and move the pending entry in your
   coverage fragment (`done_pending` the old, `pending` the new).
3. One file per credential a pilot holds. Several legacy rules usually become one credential
   file (the MEP revalidation, passenger and night passenger rules all go into
   `easa/ratings/mep-land`); one legacy rule may serve several credentials, which then
   reuse one evaluation through `uses:`.

## 2. `selects`

Translate the legacy `applies_to` of the credential's own part:

| Legacy | Credential file |
| --- | --- |
| `subject: rating`, `classes` | `selects.ratings.classes` |
| `subject: licence`, `licenceKinds` | `selects.licence.kinds` |
| `subject: privilege`, `privilegeKinds` | `selects.privilege.kinds` |
| `subject: credential`, `credentialTypes` | `selects.credential.kinds` |
| `authorities`, `excludeAuthorities` | `authorities`, `not_authorities` |
| `licenceKinds`, `excludeLicenceKinds` (on ratings, privileges) | `licence_kinds`, `not_licence_kinds` |

No two credential files may select the same record item: the gate reports overlaps. Part-FCL
class ratings exclude `LAPL_A, SPL, LAPL_S, UL`; Part-SFCL privileges select
`licence_kinds: [SPL, LAPL_S]` with `authorities: [EASA, LBA]`. An evaluation about part of
the credential narrows with `about:` and `only_for:`.

## 3. Evaluations

For each legacy rule of the credential, one evaluation:

- `id` (short, snake case: `revalidation`, `passengers_day`), `asks` (the question in plain
  words), `source` (the most specific paragraph the rule encodes), `also_cites`.
- `window` becomes a qualifier: `before_expiry_months` -> `within_months_before_expiry`,
  `rolling_days` -> `within_days`, `rolling_months` -> `within_months`, `calendar_months` ->
  `within_calendar_months`, `validity_period` -> `within_validity_period`, `since_issue` ->
  `since_licence_issue` or `since_issue`, `lifetime` -> `ever`. On the evaluation it goes
  into `counting`, on a node next to the count.
- `requirements` become `passes_if`. Each leaf's `metric` maps to the count word with that
  metric in `vocabulary.yaml` `credential_vocabulary.counts` (`minutes.total` ->
  `flight_time`, `minutes.pic` -> `pic_time`, `route_sectors` -> `route_sectors`,
  `events` with `eventKinds: [proficiency_check]` -> `proficiency_check` ...); `min` in
  minutes becomes `min_hours` or `min_minutes`. Filters become qualifiers: `classes:
  [$subject]` -> `in_class: true`, `typeDesignators: [$subject]` -> `in_type: true`,
  `variants: [$subject]` -> `in_variant: true`, `flags` -> `flagged`, `roles` -> `as`,
  `withMinutes` -> `with_time`, `fstdTypes` -> `fstd`, and so on (the qualifier table of
  the format reference). Keep the legacy row `id`, `nameKey` and `remedyKey` where they
  differ from the count word's defaults (`id:`, `name:`, `remedy:`), so the rows stay the
  same.
- `stages` become an outcome preset when they are the preset's stages (compare with the
  preset table), with `messages:`, `statuses:`, `omit:` and `special:` for the
  differences; otherwise written stages, each with a `ref`.
- `validity` becomes `valid_for`, `restored_by` becomes `restored_by`,
  `ruleDescriptionKey` becomes `description_key`.
- A legacy rule marked `support: not_supported` is not converted: its reason goes into
  the article's `not_evaluated` in your coverage fragment, in your own words. For
  `partial`, convert what is supported and state the rest as interpretations.

When the vocabulary has no word for what a rule needs, do not edit `vocabulary.yaml`:
write a vocabulary request (`fragments/vocab-requests/`) and leave the evaluation pending,
or express it with existing words and an interpretation saying what is approximated.

## 4. Shared evaluations: `uses:`

Before writing an evaluation, look under `credentials/<authority>/shared/` and in the other
credentials: passenger recency by day and night, the medical a licence needs, the flight
review are written once. Use them with `uses: <shared id>` (and `with:` for parameters) or
`uses: <credential id>#<evaluation id>`. Do not change an existing shared file or another
family's credential during parallel work: if a shared evaluation needs a new parameter,
write a new shared file and say in your changelog fragment that integration may merge the
two.

## 5. `ref:` on every condition

Every count, combinator, `waived_by`, `only_if`, `ul_credit`, `relevant_class`, validity
period, `restored_by` event, `on_fail`, written stage and expiring notice has a `ref` to the
most specific paragraph it encodes: `easa:FCL.740.A(a)(2)(ii)`, not `easa:FCL.740.A`. The
legacy `source.cite` and the quotes show which paragraph; the gate resolves every label
path against the outline of the text under `sources/` and fails when the paragraph does
not exist there. If an article you need is not stored yet, store its verbatim text first
(`sources/README.md`).

## 6. `policy:` for what is not law

A status, threshold or row set by convention cites a policy: `ref: policy:<id>`. Use the
policies of `policies.yaml` (`expiring-notice`, `unknown-input`, `recency-lapsed`,
`passengers-expired`, ...); the presets cite theirs already. A new policy goes into
`fragments/policies/` with a statement and a rationale, and must be cited.

## 7. Interpretations

Every judgement call the legacy rule makes, in its notes, its filters or its stages, is
an interpretation: what a word is read to mean ("refresher training" as dual time), what is
taken as met because the record cannot show it (who conducted a check), what part of the
article is not applied. Give each an `id`, `reading`, the `ref` it interprets, `affects`,
and `approved_by: null`, `approved_on: null`. A `with:` qualifier (who conducts or signs)
is not evaluated and always needs one.

## 8. Requirements

Name what the credential needs with a requirement entry: `requires_all: [ids]` (every
one) and `requires_any: [ids]` (at least one), e.g. a licence `requires_any` its medicals
and `requires_all` the language endorsement, a rating `requires_any` the licences it is
held on, a privilege `requires_all` its licence. Requirements point from the dependent
credential to the one it depends on, never back; a cycle fails the gate. A required
credential that is not written yet cannot be named: leave the entry out and add it when
that credential exists (say so in your changelog fragment). For a licence, add composite
examples (`composite: true`): one where it may be exercised and one where a requirement
decides that it may not.

## 9. Worked examples from the legacy cases

Every legacy case of a converted rule becomes a worked example:

| Legacy case | Worked example |
| --- | --- |
| file name | `name` |
| `description` | `says` |
| `asOf`, `record` | `asOf`, `record` (unchanged) |
| `expect[i].status` | `outcome` |
| `expect[i].subject` | `expect.subject` |
| `messageKey`, `messageParams` | `expect.message`, `expect.params` |
| `expiresOn`, `validUntil`, `requirements` | `expect.expiresOn`, `expect.validUntil`, `expect.requirements` |
| `covers`, `windowOpensAt` | dropped |

- **The outcome must equal the case's expected status.** That is the equivalence proof: if
  the example fails, the conversion differs from the verified legacy logic. Fix the
  conversion, not the outcome. The only exception is a deliberate change of this
  catalogue (such as the Part-SFCL selection of sailplane privileges): then change the
  outcome, say why in `says`, and name the change in your changelog fragment.
- A case expecting several subjects becomes one example per subject the credential
  selects (`<case>-<subject id>`), or one whose `expect.subject` names the subject that
  matters. Subjects this credential does not select belong to the credential that does.
- A case with `expect: []` (the rule does not apply) cannot be an example of its own. Make
  sure `selects` excludes those items, and show it with an example whose record also holds
  an excluded item while `expect` names no subject: the example then fails if the excluded
  item is evaluated too ("2 results").
- Add `shows: [interpretation ids]` where an example shows a reading changing the result.
- The gate needs, per evaluation, a passing example (`current`, or `expiring` with the
  requirements met) if the evaluation can report `current` or `expiring`, and a failing one
  if it can report any other status. Write one more when the legacy cases have none.

## 10. Shared files: write fragments

During the parallel round nobody edits `coverage/articles.yaml`, `messages/keys.yaml`,
`policies.yaml`, `CHANGELOG.md` or `vocabulary.yaml`. Write, named after your work
(`<authority>-<kind dir>-<name>`):

- `fragments/coverage/<name>.yaml`: the new evaluations under the article of each
  `source`, `done_pending` for the pending file you wrote, a rewritten `note` if other
  files remain pending there, and `not_evaluated` reasons for unsupported rules;
- `fragments/keys/<name>.yaml`: message keys the legacy rules use that
  `messages/keys.yaml` lacks;
- `fragments/policies/<name>.yaml`: new policies, if any;
- `fragments/changelog/<name>.md`: the lines for "Unreleased";
- `fragments/vocab-requests/<name>-<word>.md`: words the vocabulary lacks.

Each directory's README has the format.

## 11. Run the gate

```bash
go generate ./...                              # the example test of the new credential
go test ./gen/ ./credentials/
go run ./cmd/rulescheck -strict -fragments     # must print OK
gofmt -l .                                     # prints nothing
```

With `-fragments` the gate merges the coverage, key and policy fragments in memory, checks
everything as if they were merged and lists the fragments pending integration. The default
run (without the flag) fails while any fragment exists; integration merges them.

## Checklist

- [ ] File and id follow the pending name (or a better one, moved in the coverage fragment)
- [ ] `selects` overlaps no other credential; Part-SFCL privileges only on SPL and LAPL(S)
- [ ] Every evaluation has `asks` and `source`; every condition a `ref` to the most specific paragraph
- [ ] Conventions cite a `policy:`; judgement calls are interpretations; no divergence copied
- [ ] Shared evaluations reused with `uses:`; requirements with `requires_all` / `requires_any`
- [ ] Every legacy case is an example with `outcome` equal to its expected status
- [ ] Fragments for coverage, keys, policies, changelog and vocabulary requests
- [ ] `go run ./cmd/rulescheck -strict -fragments` prints OK
