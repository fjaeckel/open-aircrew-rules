# 0002. Explicit outcomes, evaluated requirements, declared policies

- Status: Accepted
- Date: 2026-09-28
- Deciders: Frederic Jung
- Amends: [0001](0001-credential-centric-catalogue.md) (points 3, 5 and 7)
- Contract: [DESIGN.md](../../DESIGN.md) sections 3, 5, 6, 8, 10 and 12

## Context

Converting the rest of the catalogue from an article-centric predecessor showed four gaps in
the format of 0001:

- A worked example only stated a status, and "passing" meant `current`. A revalidation met
  in the last 90 days reports `expiring` and was counted as failing, and nobody could tell
  from an example which outcome it was meant to show.
- `requires:` was documentation only: a licence named its medical, but nothing answered
  "may I exercise this licence today?", and the references ran both ways (a licence listed
  its ratings, the ratings their licence).
- Sailplane privileges were evaluated on a GLIDER rating of any European licence, although
  Part-SFCL privileges exist only on the SPL and the LAPL(S).
- Statuses set by convention (EASA passenger recency as `expired`, FAA as `lapsed`, the
  expiring thresholds) were hidden in the compiler's presets, and the policy ids lived in
  the vocabulary without a rationale.

## Decision

1. **Explicit outcomes.** Every worked example states `outcome:`. Outcome classes:
   passing is `current`, or `expiring` with the requirement tree met or absent; failing is
   everything else. The gate demands a passing example for every evaluation that can
   report `current` or `expiring` and a failing one for every evaluation that can report
   another status.
2. **`requires_all` and `requires_any`.** They replace `requires:` and are evaluated: a
   held credential's composite is the conjunction of its own evaluations (passenger,
   launch-method, variant and training results as limitations only) and its requirement groups,
   with `unknown` for a required credential the record does not hold. Requirements point
   from the dependent credential to the one it depends on; cycles are errors.
3. **Part-SFCL privileges on Part-SFCL licences.** The SPL credential and the cloud-flying
   privilege select only items on SPL and LAPL(S) licences of authority EASA or LBA.
4. **`policies.yaml`.** Every policy has an id, a statement and a rationale; outcome presets
   cite a policy for each status they set by convention; a status override cites its own;
   the gate fails on an unresolved or an uncited policy.

Parallel conversion work writes fragments (`fragments/`) instead of editing the shared
files; the gate accepts them only with `-fragments`.

## Consequences

- Examples are self-describing, and a converted legacy case proves equivalence by its
  expected status alone: it must equal the new example's outcome.
- Consumers get a composite per held credential (`credentials.Catalogue.Evaluate`,
  `cmd/evaluate`), whose output changed from a list to `{ evaluations, credentials }`.
- A licence no longer lists its ratings or privileges; the ratings and privileges name the
  licence they need.
- A GLIDER rating on a PPL(A) or another non-SFCL licence reports no sailplane recency.

## Alternatives considered

- **Classes by status only** (expiring always passing): would count a revalidation whose
  requirements are not met as passing.
- **Composite as the conjunction of every evaluation**: would say a pilot without passenger
  recency may not exercise the licence at all; passenger and launch-method results are
  limitations instead.
- **Policy refs written in every credential file**: repeats the same convention in dozens
  of files; declaring them once per preset stage keeps them reviewable in one place.
