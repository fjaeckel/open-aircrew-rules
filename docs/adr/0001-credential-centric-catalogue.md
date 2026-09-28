# 0001. One file per credential, with references and interpretations made explicit

- Status: Accepted
- Date: 2026-09-28
- Deciders: Frederic Jung
- Contract: [DESIGN.md](../../DESIGN.md); format: [credential-format.md](../credential-format.md)

## Context

Pilots, instructors, examiners and logbook applications all ask the same question: "what
does my SEP rating (my PPL, my class 2 medical) need, and do I meet it today?". The answer is
spread over several regulations: a revalidation article, two passenger articles, the medical
articles, the language article. A machine-readable catalogue of these rules has to be
reviewable by people who know the regulations but do not read code, and precise enough for
software to evaluate a logbook against it.

Three things make that hard:

- **The unit of review.** Organised by article, a reviewer has to read many files to see
  what one credential needs, and cannot see whether anything is missing.
- **Hidden judgement.** Much of the text needs reading before it can be counted ("refresher
  training" as dual time, "hold an IR" as a valid IR on the same licence). When such
  readings sit in free-text notes next to plain facts, nobody can review them as such or
  sign them off.
- **Tests nobody reads.** Exhaustive fixtures derived from the rule structure test the tool,
  not the rules, and do not show a reader how a rule behaves.

## Decision

1. **One file per credential** (licence, rating, privilege, endorsement, instructor or
   examiner certificate, medical) under `credentials/<authority>/<kind>/`, listing every
   evaluation the credential needs. Evaluations shared by many credentials (passenger
   recency, the medical a licence needs, the flight review) live in
   `credentials/<authority>/shared/` and are listed by each credential with `uses:`;
   dependencies on other credentials are listed with `requires:`.
2. **A near-English vocabulary** (`flight_time`, `takeoffs`, `within_days`, `in_class`,
   `as: pilot_flying`, `waived_by`, `only_if`, ...) declared in `vocabulary.yaml` and
   **compiled onto a small evaluation engine** with a closed vocabulary of subjects, metrics,
   filters, windows, combinators and status stages. Outcome presets generate the status
   stages so authors rarely write them.
3. **References at fine granularity, no quotes.** Every evaluation has a `source`, every
   condition a `ref` to the most specific paragraph it encodes (`easa:FCL.740.A(b)(1)(ii)(C)`).
   The gate resolves each against verbatim texts in `sources/` by parsing their paragraph
   outline. Conditions with no legal text behind them are marked `policy:<id>`.
4. **Interpretations are first-class**: each reading of the text has an id, the paragraph
   it interprets, the evaluations it affects, and `approved_by` / `approved_on`. The gate
   lists the unapproved ones.
5. **Worked examples instead of fixtures**: a few plain examples per credential, at least
   one passing and one failing per evaluation and one per interpretation that changes a
   result.
6. **A structural gate**: schemas; every reference, `uses:` and `requires:` resolves; message
   keys exist; examples pass and cover each evaluation; no two credentials select the same
   record item; every article in scope is accounted for (`coverage/articles.yaml`:
   evaluated, pending or not evaluated with a reason); only allow-listed source texts are
   stored; engine statement coverage of at least 95 %.
7. **No product-specific material.** The catalogue encodes the texts, not the behaviour of
   any application. Differences between an application and the catalogue belong to that
   application.

## Consequences

- Reading one file answers what a credential needs; reviewing a rule means reviewing its
  conditions, their references and its interpretations.
- Parameterised shared evaluations and `uses:` of another credential's evaluation avoid
  copying; the price is one indirection when reading.
- Articles whose credential is not written yet are visible as `pending` in the coverage map,
  so the catalogue is a complete map of the regulation even while it grows.
- Licensing: code, catalogue, examples and tools are MIT; texts under `sources/` keep their
  own status (US federal works, 17 U.S.C. § 105; German official works, § 5(1) UrhG; EU
  legal acts reused under Commission Decision 2011/833/EU with attribution) and are not
  relicensed.

## Alternatives considered

- **Article files with generated per-credential views.** Readers would get a view, but the
  reviewed artefact would still be the article file with its hidden readings.
- **An engine built around credentials.** Unnecessary: an evaluation maps one to one onto an
  engine rule, and the engine's semantics can be verified on their own.
- **Quotes in credential files.** Duplicates `sources/`, bloats the files and invites
  paraphrase; paragraph references resolved by the gate, with a hash anchor later, give the
  same traceability.

## History

The catalogue was first organised by regulation article: one file per article, mixing
verbatim quotes, a subject selector, a requirement tree, status stages and free-text notes,
gated by exhaustive coverage-tag fixtures. The credential-centric structure replaced it; the
converted credentials were checked against the article rules' fixtures and gave identical
results before those were retired. The engine, its closed vocabulary, the neutral record,
the message keys and the source allow-list carried over unchanged.
