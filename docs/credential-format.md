# Credential file format

One YAML file per credential a pilot holds (a licence, a rating, a privilege, an
endorsement, an instructor or examiner certificate, a medical certificate) lists every
evaluation that credential needs. A pilot can read it; the compiler in `credentials/`
(package `credentials`) turns each evaluation into a rule of the engine's closed vocabulary
(DESIGN.md sections 3 and 4), and the gate (`cmd/rulescheck`) checks it structurally.

Everything a credential file may say is declared in `vocabulary.yaml`, section
`credential_vocabulary`. A word that is not declared there is an error.

## Layout

```text
credentials/<authority>/licences/<name>.yaml        kind: licence
credentials/<authority>/ratings/<name>.yaml         kind: rating
credentials/<authority>/privileges/<name>.yaml      kind: privilege
credentials/<authority>/endorsements/<name>.yaml    kind: endorsement
credentials/<authority>/instructors/<name>.yaml     kind: instructor_certificate
credentials/<authority>/examiners/<name>.yaml       kind: examiner_certificate
credentials/<authority>/medicals/<name>.yaml        kind: medical
credentials/<authority>/shared/<name>.yaml          parameterised evaluations used by several credentials
examples/<credential id>.yaml                       worked examples of one credential
coverage/articles.yaml                              every article in scope -> evaluations, pending files, or why not
scope/articles-<authority>.yaml                     the articles in scope, with a one-line summary each
sources/<authority>/<article>.md                    verbatim texts (unchanged; see sources/README.md)
```

A credential id is `<authority>.<kind>.<file name>` (`easa.rating.sep-land`), a shared
evaluation id `<authority>.shared.<file name>`. An evaluation is named
`<credential id>#<evaluation id>` (`easa.rating.sep-land#revalidation`).

## Worked example: the SEP (land) class rating

```yaml
credential: SEP (land) class rating
id: easa.rating.sep-land
kind: rating
authority: EASA
held_on: [PPL(A), CPL(A), ATPL(A), MPL]
selects:
  ratings: { classes: [SEP_LAND], authorities: [EASA, LBA], not_licence_kinds: [LAPL_A, SPL, LAPL_S, UL] }
validity: { period_months: 24, ref: "easa:FCL.740(a)(1)" }

evaluations:
  - id: revalidation
    asks: Is the rating valid, and has the holder done what keeps it valid?
    source: easa:FCL.740.A(b)(1)
    also_cites: [easa:FCL.740.A(b)(2), easa:FCL.035(a)(4)]
    relevant_class: { pooled_with_held: [SEP_LAND, TMG], ref: easa:FCL.740.A(b)(2) }
    counting: { within_months_before_expiry: 12 }
    passes_if:
      id: revalidation
      ref: easa:FCL.740.A(b)(1)
      any_of:
        - proficiency_check: { in_class: true, with: examiner, within_months_before_expiry: 3, ref: easa:FCL.740.A(b)(1)(i) }
        - id: experience
          ref: easa:FCL.740.A(b)(1)(ii)
          all_of:
            - flight_time: { min_hours: 12, in_class: true, ul_credit: &ul { SEP_LAND: [THREE_AXIS], TMG: [THREE_AXIS_MOTORGLIDER], ref: easa:FCL.035(a)(4) }, ref: easa:FCL.740.A(b)(1)(ii) }
            - pic_time: { min_hours: 6, in_class: true, ul_credit: *ul, ref: easa:FCL.740.A(b)(1)(ii)(A) }
            - takeoffs: { min: 12, in_class: true, ref: easa:FCL.740.A(b)(1)(ii)(B) }
            - landings: { min: 12, in_class: true, ref: easa:FCL.740.A(b)(1)(ii)(B) }
            - refresher:
                min_hours: 1
                in_class: true
                with: [FI, CRI]
                ref: easa:FCL.740.A(b)(1)(ii)(C)
                waived_by:
                  events: [proficiency_check, skill_test, ebt_assessment, assessment_of_competence]
                  categories: [aeroplane]
                  ref: easa:FCL.740.A(b)(1)(ii)(C)
    outcomes:
      preset: revalidation
      expiring_notice: { days: 90, ref: policy:expiring-notice }
    on_fail: { what: renewal, ref: easa:FCL.740(b) }
    description_key: easa_sep_tmg

  - id: passengers_day
    asks: May the holder carry passengers in an SEP land aeroplane?
    uses: easa.shared.passengers-day

  - id: passengers_night
    asks: May the holder carry passengers at night in an SEP land aeroplane?
    uses: easa.shared.passengers-night

  - id: licence
    asks: Is the rating held on a licence whose own requirements (medical, language) are met?
    requires: [easa.licence.ppl-a]

interpretations:
  - id: refresher-as-dual
    reading: The refresher training is 1 hour of dual time received in the class; that the instructor holds an FI or CRI certificate, and the exercises flown, are not checked.
    ref: easa:FCL.740.A(b)(1)(ii)(C)
    affects: [revalidation]
    approved_by: null
    approved_on: null
  # ... eight more (credentials/easa/ratings/sep-land.yaml)
```

Read aloud: the rating is revalidated by a proficiency check in the class in the 3 months
before expiry, or by 12 hours in the class in the 12 months before expiry including 6 hours
as PIC, 12 take-offs, 12 landings and 1 hour of refresher training, the refresher waived by
a check, test or assessment in any aeroplane. SEP land and TMG flights pool when both
ratings are held; ultralight time counts toward the hours. TMG (`easa.rating.tmg`) and SEP
sea (`easa.rating.sep-sea`) reuse this evaluation with `uses: easa.rating.sep-land#revalidation`.

## Top level

| Key | Meaning |
| --- | --- |
| `credential` | Name a pilot recognises. |
| `id`, `kind`, `authority` | As above; `kind` is one of `credential_vocabulary.kinds`. |
| `held_on` | Licences the credential is held on, for readers (not evaluated). |
| `selects` | How the credential appears in a record: `licence` (`kinds`: licence kinds), `ratings` (`classes`), `privilege` (`kinds`: privilege kinds), `credential` (`kinds`: certificate types), each with optional `authorities`, `not_authorities`, `licence_kinds`, `not_licence_kinds`. No two credential files may select the same record item (gate). |
| `validity` | The credential's validity period for readers (`period_months`, `ref`); an evaluation's `valid_for` or `outcomes` is what is evaluated. |
| `evaluations` | Every evaluation the credential needs, in reading order. |
| `interpretations` | Every reading of the text the evaluations depend on (below). |

## Evaluations

An evaluation is one of three things.

1. **Written here**: `id`, `asks` (the question in plain words), `source` (the article
   reference), optional `also_cites`, then what it is about, what counts, what passes and
   what it reports.
2. **Used from elsewhere**: `uses: <shared id>` (with `with:` for the shared file's
   parameters) or `uses: <credential id>#<evaluation id>`. The evaluation is compiled for
   this credential's `selects`; `asks` and `only_for` given here override. A shared
   evaluation with its own `scope` (the flight review: one per pilot) is compiled once,
   whichever credentials list it.
3. **A reference**: `requires: [credential ids]` says that this credential needs another
   one (the licence's medical, the rating's licence). It compiles to nothing; the other
   credential evaluates itself. It keeps "what does this credential need?" answerable
   from one file.

| Key | Meaning |
| --- | --- |
| `about` | The subject: `self` (default; the credential's own record part, the licence for kind licence), `rating`, `passengers` (one result per class and authority), `launch_methods` (one per launch method used), `licence`, `training` (a programme), `pilot` (one result per pilot). |
| `only_for` | Narrows the subjects: `authorities`, `not_authorities`, `licence_kinds`, `not_licence_kinds`, `classes`, `not_classes`, `credential_types`, `privilege_kinds`, `launch_methods`, `type_rated`, `programme`, `when_holding` (engine `holds`), `ref`. |
| `scope` | Replaces `selects` for this evaluation (same keys as `only_for`). |
| `relevant_class` | `pooled_with_held: [classes]` with `ref`: `in_class` then also counts the pool's classes the holder rates on the same licence. |
| `counting` | Qualifiers for every count of the evaluation (window and filters); the window here is also the period the outcomes speak of. |
| `passes_if` | The requirement tree (below). May be absent when the outcomes decide on holdings or dates alone. |
| `restored_by` | Events that make the tree count as met for the evaluation's period: `- { ipc: { simulator: include }, ref: faa:61.57(d)(1) }`. |
| `effective_from`, `effective_to` | The dates the evaluation applies to, both inclusive (`YYYY-MM-DD`). When a regulation changes, the old evaluation gets `effective_to` and its successor (a new id, e.g. `revalidation_2027`) `effective_from`; outside its period an evaluation reports nothing. A `uses:` evaluation may set its own. |
| `valid_for` | A derived expiry: `counted_from` (`issue` or `valid_from`), `periods` of `{ age_under?, age_from?, months, ends_at_age?, end_of_month?, ref }`; the first period whose ages hold applies, the earlier of the derived and a recorded expiry wins. |
| `outcomes` | Status and message stages (below). |
| `description_key` | Rule description key from `messages/keys.yaml`. |
| `on_fail` | What the holder does when it fails (`what`, `ref`), for readers. |

## passes_if

A node is either a **combinator** or **one count**.

```yaml
id: experience            # optional; names the node in results
ref: easa:FCL.740.A(b)(1)(ii)
within_months: 24         # qualifiers here apply to every count below
all_of: [ ... ]           # or any_of: [ ... ], or n_of: { n: 2, of: [ ... ] }
```

```yaml
- takeoffs: { min: 12, in_class: true, ref: easa:FCL.740.A(b)(1)(ii)(B) }
```

**Counts** (`credential_vocabulary.counts`): `flight_time`, `pic_time`, `dual_time`,
`supervised_solo_time`, `instruction_time` (dual + supervised solo), `pic_dual_or_solo_time`,
`refresher` (dual in the class), `cloud_flying_time`, `cloud_flights`,
`longest_training_flight`, `flights`, `training_flights`, `takeoffs`, `landings`,
`takeoffs_and_landings` (the smaller of the two), `night_takeoffs`, `night_landings`,
`full_stop_landings`, `full_stop_night_landings`, `launches`, `approaches`, `holds`,
`intercept_and_track`, `tows`, `not_recorded` (something no record holds: always
untracked), and the events `proficiency_check`, `skill_test`, `flight_review`,
`proficiency_program_phase`, `basicmed_course`, `basicmed_exam`. Each has a default row id,
name key, unit and remedy key (the requirement rows consumers already display); override
with `id`, `name`, `unit`, `remedy` (`remedy: none` drops it).

Amounts: `min` for counts (default 1), `min_hours` or `min_minutes` for times.

**Qualifiers** (`credential_vocabulary.qualifiers`), on a count, a combinator, in
`counting` or on a `restored_by` event:

| Qualifier | Counts |
| --- | --- |
| `within_days: n`, `within_months: n`, `within_calendar_months: n` | the period before the date evaluated, both ends included |
| `within_months_before_expiry: n`, `within_validity_period: true` | periods anchored on the credential's expiry |
| `since_licence_issue: true`, `since_issue: true`, `ever: true` | since an issue date, or everything |
| `in_class: true`, `classes`, `excluding_classes`, `categories`, `ul_kinds` | where it was flown |
| `ul_credit: { CLASS: [ultralight kinds], ref }` | ultralight time credited to a class |
| `in_type: true` | in the type of the rating evaluated |
| `launch_methods`, `by_this_launch_method: true` | how a sailplane was launched |
| `as: pilot_flying`, `as: sole_manipulator`, `as: [pic, dual, spic, ...]` | the pilot's role |
| `with_time: [ifr, crossCountry, ...]`, `without_time` | flights with or without time of a kind |
| `simulator: include \| only`, `fstd: [FFS]` | simulator sessions (excluded by default) |
| `tailwheel`, `min_distance_km` | aircraft and flight properties |
| `for_rating`, `also: [event kinds]` | events for a rating; further event kinds that count |
| `any_flight_of: [ {qualifiers}, ... ]` | a flight counts when it matches any set |
| `with: examiner` | who conducts it. **Not evaluated**; the gate requires an interpretation that says so |

**Modifiers of a count**: `waived_by: { events: [...], <qualifiers>, ref }` (the count or
one of these events; compiles to `any_of` named after the count), `only_if: { <engine
condition>, ref }` (the row exists only while the condition holds), `informational: true`
(shown, never decides), `credit: <escape hatch>`, `messages: { met, unmet, untracked }`.

## Outcomes

`outcomes: <preset>` or `{ preset, ... }` picks a standard stage list; `outcomes: [stages]`
writes them out. Each stage is `{ id?, when, status, message, params?, ref }`, `when` being
an engine stage condition (`all_met`, `undetermined`, `{ unmet: <row> }`, `{ holds: ... }`,
`{ met_within: ... }`, `all`, `any`, `not`, ...). Keys come from `messages/keys.yaml`.

| Preset | Stages (id: condition -> status, message) |
| --- | --- |
| `revalidation` | no_expiry -> unknown `rating.no_expiry_date`; expired -> expired `rating.expired`; *special*; window_not_open: before the period -> current `rating.window_not_open`; met_expiring: met and within `expiring_notice` -> expiring `rating.revalidation_expiring_met`; met -> current `rating.revalidation_current`; not_met -> expiring `rating.revalidation_not_met` |
| `recency` | *special*; current: met -> current `rating.recency_current`; lapsed -> lapsed `rating.recency_not_met` |
| `privilege_recency` | expired -> expired `privilege.expired`; *special*; current -> `privilege.recency_current`; lapsed -> `privilege.recency_not_met` (all with the expiry date) |
| `passengers` | *special*; current -> current `pax.day_current`; unknown: undetermined -> unknown `pax.experience_not_logged`; not_met -> expired `pax.not_current`, naming the missing count of the last listed unmet requirement |
| `validity` | date_of_birth (with `date_of_birth: when_no_expiry` or `required`) -> unknown; no_expiry -> unknown `credential.no_expiry_date`; expired -> expired `credential.expired`; *special*; expiring: within `expiring_notice` -> expiring `credential.expiring`; valid -> current `credential.valid` |
| `training` | *special*; complete -> current `training.all_met`; unknown -> unknown `training.distance_unknown`; in_progress -> lapsed `training.in_progress` |

Options: `special: [stages]` (exemptions and waivers, checked after the expiry stages and
before the requirement stages), `messages: { <stage id>: key }`, `statuses: { <stage id>:
status }`, `omit: [stage ids]`, `expiring_notice: { days, message?, ref }` (required by
`revalidation` and `validity`).

## References

Every evaluation has a `source`; every count and combinator of `passes_if`, every
`waived_by`, `only_if`, `ul_credit`, `relevant_class`, `valid_for` period, `restored_by`
event, `on_fail`, written stage and `expiring_notice` has a `ref`; every interpretation
names the paragraph it interprets. A ref is `<prefix>:<article><paragraph labels>`:

- `easa:FCL.740.A(b)(1)(ii)(C)` resolves to `sources/easa/fcl-740-a.md`;
  `faa:61.57(c)(1)(iii)` to `sources/faa/61.57.md`; `de:LuftPersV.45a` to
  `sources/de/luftpersv-45a.md` (`credential_vocabulary.ref_authorities`).
- The gate parses the paragraph outline of the source text ((a), (1), (i), (A), nested in
  whatever order the text uses) and requires the exact label path to exist.
- `policy:<id>` marks a condition with no legal text behind it (a presentation or
  implementation choice such as the 90-day expiring notice); the id must be declared in
  `credential_vocabulary.policy_refs`. Keep these few.
- `ref: [a, b]` gives several. No quote is pasted into a credential file; the text lives
  in `sources/`. An optional `@<hash>` suffix is reserved for anchoring the quoted passage
  (change detection); it is parsed and not yet verified.

## Interpretations

```yaml
- id: exemption-period
  reading: The text does not say when an exempting check must have been passed; it counts when passed within the 12 months before expiry.
  ref: easa:FCL.740.A(b)(1)(ii)(C)
  affects: [revalidation]            # evaluation ids of this file ("evaluation" in a shared file; "*" for all)
  approved_by: null                  # who signed the reading off
  approved_on: null                  # YYYY-MM-DD
```

Every judgement call is one: what a word is read to mean, what is taken as met because the
record cannot show it, what part of an article is not applied. There are no free-text notes
in credential files. The gate lists unapproved interpretations (report only) and fails on
one without `ref`, `reading` or `affects`.

## Worked examples

`examples/<credential id>.yaml` holds a few examples: at least one passing (status
`current`) and one failing example per compiled evaluation (only one kind when the
evaluation cannot report the other), plus examples for interpretations that change a result
(`shows: [interpretation ids]`).

```yaml
credential: easa.rating.tmg
examples:
  - name: revalidation-pooled-with-sep-land
    evaluation: revalidation
    shows: [pooling-same-licence]
    says: With SEP land and TMG ratings on the same licence, SEP land flights revalidate the TMG rating too.
    asOf: 2026-09-15
    record:
      licences: [{ id: l-ppl, authority: EASA, type: PPL(A) }]
      ratings:
        - { id: r-sep, licenceId: l-ppl, class: SEP_LAND, expires: 2027-03-31 }
        - { id: r-tmg, licenceId: l-ppl, class: TMG, expires: 2027-03-31 }
      flights:
        - { date: 2026-09-01, class: SEP_LAND, minutes: { total: 720, pic: 660, dual: 60 }, takeoffs: { day: 12 }, landings: { day: 12 } }
    expect: { subject: { kind: rating, id: r-tmg }, status: current, message: rating.revalidation_current }
```

`record` is the neutral input (`schema/record.schema.json`). `expect` compares what it
states: `subject` (needed when the evaluation reports several), `status`, `message`,
`params`, `expiresOn`, `validUntil`, `requirements`. `go generate` writes one test per
credential under `gen/` that runs its examples.

## Coverage

`coverage/articles.yaml` accounts for every article of `scope/articles-*.yaml`:

```yaml
- article: "easa:FCL.060"
  evaluations: [easa.rating.sep-land#passengers_day, easa.rating.sep-land#passengers_night]
  pending: [easa/ratings/mep-land, easa/licences/lapl-a]     # credential files still to write
  note: "Still to convert beyond the evaluations above: passenger recency for the remaining classes."
- article: "easa:FCL.050"
  not_evaluated: "Not a currency or validity rule: the obligation to keep the logbook itself."
```

An article lists the evaluations that encode it, the credential files planned to encode it
(`pending`, with a one-line `note`), and `not_evaluated: <reason>` for what is not
evaluated; at least one of the three. Every compiled evaluation must be listed under the
article of its `source`. A pending file that exists fails the gate: list its evaluations
instead. The gate reports pending articles without failing.

## How it compiles

Each evaluation becomes one engine rule (`engine.Rule`), evaluated by the engine:

| Credential file | Engine rule |
| --- | --- |
| `selects` part chosen by `about`, then `scope` / `only_for` | `applies_to` (subject kind, authorities, licence kinds, classes, privilege kinds, credential types, holds) |
| `counting` | rule `window` and `filter` |
| combinator node, its qualifiers | `all_of` / `any_of` / `n_of` node with `window` and `filter` (inherited by children) |
| count word + amount + qualifiers | leaf: `metric`, `min` (hours x 60), `unit`, `nameKey`, `remedyKey`, `filter` (with the word's implied filters, e.g. `eventKinds`) |
| `in_class` (+ `relevant_class`) | `classes: [$subject]` (+ `heldClassPools`) |
| `waived_by` | `any_of: [count, events leaf <word>_exemption]` |
| `only_if`, `informational`, `credit` | `when`, `informational`, `escape_hatch` |
| `restored_by` | `restored_by` hooks |
| `valid_for` | `validity` (age periods, caps, end of month) |
| `outcomes` | `stages` |
| `effective_from`, `effective_to` | `effective_from`, `effective_to` |

The compiled rule id is `<credential id>#<evaluation id>`, or the shared id for a
shared evaluation with its own `scope`.
