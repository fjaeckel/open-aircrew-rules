# Design contract

This repository holds pilot licence, rating, privilege and medical rules as reviewed plain
text (YAML), organised by credential, together with an evaluator, a compiler, a code
generator and a gate. Every contributor follows this document; where the format reference
([docs/credential-format.md](docs/credential-format.md)) goes into more detail, both agree.
The reasons are recorded in [ADR 0001](docs/adr/0001-credential-centric-catalogue.md).
Changing this document is a design change: say so in the pull request.

## 1. Scope and independence

- The catalogue answers, for one pilot's record on one date, what every credential they
  hold needs and whether it is met: recent experience, revalidation and renewal, validity
  periods, privilege recency, flight reviews and proficiency checks, passenger carriage and
  medical validity. Training syllabi and test content are out of scope, apart from the
  experience they require.
- Authorities: EASA (Part-FCL, Part-SFCL, Part-BFCL, the validity articles of Part-MED), FAA
  (14 CFR Part 61 and Part 68) and Germany (LuftPersV for ultralights). The articles in scope
  are listed in `scope/articles-<authority>.yaml` (section 9).
- The module has no product dependency. It reads only files under its root and assumes no
  database, API or user model: its input is the neutral record of section 5.
- Allowed third-party Go dependencies: `gopkg.in/yaml.v3` and
  `github.com/santhosh-tekuri/jsonschema/v6`. Anything else needs an ADR.

## 2. Layout

```text
credentials/<authority>/<kind dir>/<name>.yaml   one credential and every evaluation it needs
credentials/<authority>/shared/<name>.yaml       parameterised evaluations several credentials use
credentials/*.go                                 package credentials: loader, compiler, references, checks
examples/<credential id>.yaml                    worked examples, run by the generated tests in gen/
coverage/articles.yaml                           every article in scope and how it is accounted for
scope/articles-<authority>.yaml                  the articles in scope, one own-words summary each
sources/<authority>/<article>.md                 verbatim texts of allowed origin only (section 6)
vocabulary.yaml                                  the closed vocabulary, incl. credential_vocabulary
messages/keys.yaml                               every message, requirement, remedy and description key
schema/                                          JSON Schema 2020-12 for every YAML file
engine/, engine/hatches/                         the evaluator (package engine)
gen/                                             generated constants and tests (never edited by hand)
cmd/rulesgen, cmd/rulescheck, cmd/evaluate       generator, gate, command-line evaluator
docs/                                            format reference, AMC/GM notes, ADRs
```

Kind directories: `licences`, `ratings`, `privileges`, `endorsements`, `instructors`,
`examiners`, `medicals`. Authority directories: `easa`, `faa`, `de`.

## 3. Credentials

One file per credential a pilot holds (licence, rating, privilege, endorsement, instructor
or examiner certificate, medical) lists every evaluation it needs:

- evaluations written in the file;
- evaluations used from a shared file or from another credential (`uses:`, with `with:`
  parameters for a shared file);
- other credentials it depends on (`requires:`), which compile to nothing and keep "what
  does this credential need?" answerable from one file.

`selects` says which record items the credential is (a licence, ratings on a licence, a
licence privilege, a certificate); no two credential files may select the same item. Ids:
`<authority>.<kind>.<file name>`, `<authority>.shared.<file name>`; an evaluation is
`<credential id>#<evaluation id>`. An evaluation may carry `effective_from` and
`effective_to` (both inclusive); a regulation change ends one evaluation and starts its
successor rather than editing it in place.

## 4. Engine rules and the closed vocabulary

Credential files use only the words of `vocabulary.yaml` `credential_vocabulary`: counts,
qualifiers, outcome presets, subjects (`about`), reference authorities and policy
references. The compiler maps each evaluation onto one engine rule (`engine.Rule`) of the
closed vocabulary:

- **Subjects**: `rating`, `licence`, `privilege`, `credential`, `passengers` (one per class
  and authority, or per class and type designator for type-rated aircraft), `flight_review`
  (one per pilot), `training` (one programme), `type` (one variant), `launch_method`.
- **Metrics**: named counters over the flights or events in a window after filtering
  (`minutes.total`, `landings.total`, `takeoffs_and_landings`, `approaches`, `events`, ...).
  `not_recorded` stands for something no record holds: always untracked, never met.
- **Filters**: classes, ultralight kinds, launch methods, roles, simulator use, flags,
  event kinds, aircraft properties and more; absent input makes a flight unknown, not
  excluded.
- **Windows**: `rolling_days`, `rolling_months`, `calendar_months`,
  `before_expiry_months`, `validity_period`, `since_issue`, `lifetime` (section 11).
- **Combinators**: `all_of`, `any_of`, `n_of`.
- **Stages**: an ordered list; the first whose `when` holds sets status and message key.
  Conditions include `all_met`, `met`/`unmet`, `undetermined`, `expired`, `no_expiry`,
  `expires_within`, `valid_until_within`, `met_within`, `holds`, `missing`, `all`, `any`,
  `not`.
- **Events**: `restored_by` (an IPC, a proficiency check) makes the tree count as met from
  the event date.
- **Escape hatches**: logic implemented in Go under `engine/hatches/`, each declared in the
  vocabulary; used only where the vocabulary cannot express the text.

The engine implements every declared entry (a missing metric is a compile error through the
generated `engine/zz_metrics_gen.go`; the gate lists any other gap). Adding a word means:
the vocabulary, the engine or the compiler, a schema entry where needed, and a worked
example that uses it.

## 5. Record and evaluation

The engine takes one neutral **record** and a date `asOf`. `schema/record.schema.json` is the
authoritative definition:

```yaml
holder:      { dateOfBirth? }
licences:    [{ id, authority, type, kind?, issued?, expires? }]
ratings:     [{ id, licenceId, class, ulKind?, typeDesignator?, issued?, validFrom?, expires? }]
privileges:  [{ id, licenceId, kind, detail?, issued?, validFrom?, expires? }]
credentials: [{ id, type, issued?, validFrom?, expires? }]      # medicals, language, radio
variants:    [{ id, ratingId, name }]
events:      [{ date, kind, class?, ... }]                        # checks, tests, courses
flights:
  - date, class, ulKind?, typeDesignator?, variant?, launchMethod?, isSimulator, fstdType?
    minutes: { total, pic, dual, spic, picus, sic, dualGiven, examiner, multiPilot, night,
               ifr, actualInstrument, simulatedInstrument, crossCountry }
    takeoffs: { day, night }   landings: { day, night }
    optional inputs: fullStopNightLandings, interceptAndTrack, pilotFlying, soleManipulator,
    tailwheel, mtomKg, engines, towKind, cruiseMinutes, mountainLandings, ...
    flags: { proficiencyCheck, flightReview, ipc, trainingFlight, towFlight, ... }
```

It returns one **evaluation** per (compiled rule, subject):

```yaml
ruleId, subject: { kind, id?, class?, ulKind?, detail? }
status: current | expiring | expired | lapsed | unknown | not_applicable
messageKey, messageParams?, ruleDescriptionKey
expiresOn?, windowOpensAt?, validUntil?
requirements: [{ id, nameKey, metric, current, required, unit, met, tracked,
                 validUntil?, lastDate?, messageKey?, remedyKey?, remedyParams? }]
citations: ["EASA FCL.740.A(b)(1)", ...]
```

Keys are stable, language-neutral ids from `messages/keys.yaml`; applications translate
them. A new key is added there before a credential uses it. `engine.Evaluate` drops a
rule's evaluation of a subject when a rule that `supersedes` it evaluated the same subject.

## 6. References and sources

Every evaluation has a `source`; every condition of `passes_if` and every `waived_by`,
`only_if`, `ul_credit`, `relevant_class`, validity period, `restored_by` event, `on_fail`,
written stage and expiring notice has a `ref` to the most specific paragraph it encodes
(`easa:FCL.740.A(b)(1)(ii)(C)`, `faa:61.57(c)(1)(iii)`, `de:LuftPersV.45a`). The gate
resolves each: the prefix is declared, the article file exists under `sources/`, and the
exact paragraph label path occurs in its outline. Credential files contain no quotes; an
`@<hash>` suffix is reserved for anchoring quoted passages later.

**Source kinds.** A condition rests on one of:

- `regulation`: EU regulations and US federal regulations, stored verbatim;
- `national_law`: German statutes and ordinances, stored verbatim;
- `association`: rules an association sets under a statutory delegation (DULV, DAeC). They
  are cited by title and never stored or quoted; the delegating statute is stored;
- `policy:<id>`: a presentation or implementation choice with no legal text behind it (the
  90-day expiring notice, unknown when input is missing). Every id is declared in
  `credential_vocabulary.policy_refs`, with its reason. Keep these few.

**Copyright allow-list (enforced).** `vocabulary.yaml` `source_origins` lists the only
origins a file under `sources/` may have: `us-federal` (17 U.S.C. § 105), `de-amtliches-werk`
(§ 5(1) UrhG) and `eu-legal-act` (Commission Decision 2011/833/EU, attribution "© European
Union, https://eur-lex.europa.eu"). Every file declares its origin and attribution in its
header, lives in its origin's directory, cites only the origin's hosts, and names no
AMC, GM, Easy Access Rules, ICAO or association material in its header or headings. EASA
AMC/GM, ICAO and association documents are never stored; `docs/amc-gm-notes.md` holds
own-words notes that name them. The texts keep their own status and are not relicensed
(NOTICE, `sources/README.md`).

## 7. Interpretations

Every judgement call is an interpretation with `id`, `reading`, `ref` (the paragraph it
interprets), `affects` (the evaluations it changes) and `approved_by` / `approved_on`:
what a word is read to mean, what is taken as met because the record cannot show it, what
part of an article is not applied. There are no free-text notes in credential files. The
gate lists unapproved interpretations and does not fail on them; an evaluation that uses a
documentation-only qualifier (`with: examiner`) must have one. Sign-off follows
CONTRIBUTING.md.

## 8. Worked examples

`examples/<credential id>.yaml` holds, per compiled evaluation, at least one passing
(`current`) and one failing example, unless the evaluation cannot report one of them, plus
examples that show an interpretation changing a result (`shows:`). `go generate` writes one
test per credential under `gen/` that runs them.

## 9. Articles in scope and coverage

`scope/articles-<authority>.yaml` lists every article in scope with its cite, official URL,
subject, a one-line own-words summary and the verbatim source file. `coverage/articles.yaml`
accounts for each one with any of:

- `evaluations`: the compiled evaluations that encode it;
- `pending`: the credential files (`<authority>/<kind dir>/<name>`) planned to encode it,
  with a one-line `note`;
- `not_evaluated`: why (the rest of) it is not evaluated.

Every compiled evaluation is listed under the article of its `source`. A pending file that
already exists fails the gate.

## 10. The gate

`go run ./cmd/rulescheck -strict` fails on: schema errors in any YAML file; credentials that
do not load or compile; references that do not resolve; unknown message keys or statuses;
missing or failing worked examples; two credentials that select the same record item;
malformed interpretations; vocabulary entries the engine does not implement; articles of the
scope missing from the coverage map, unknown evaluations, evaluations not listed under their
source article, pending files that exist; source files that break the allow-list; and engine
statement coverage below 95 %. It reports without failing: articles with pending
credentials, and interpretations awaiting approval. CI also runs `go vet`, the tests with
the race detector, the generator up-to-date check and a vulnerability scan.

## 11. Decisions that apply throughout

1. **Dates, not instants.** Windows work on calendar dates and include both ends.
   `rolling_days: n` covers `asOf - n days` through `asOf`. `calendar_months: n` covers the n
   calendar months before the month of `asOf` plus the current month to date.
   `rolling_months` clamps to the end of the month (31 March minus one month is the last day
   of February). `before_expiry_months` counts back from the expiry date.
2. **Valid through the expiry date.** A credential with expiry date D is valid on D and
   expired from D + 1, for ratings, privileges, certificates and medicals alike.
3. **Authorities match case-insensitively** after trimming (`easa`, `EASA`, ` Easa `).
4. **Unknown is not zero.** Absent optional input is unknown. A requirement whose metric is
   unknown for every flight in its window is `tracked: false` and never met; the result is
   `unknown` with a message that says what to record, never a silent pass or fail.
5. **The text as written.** An evaluation encodes the article, not the behaviour of any
   application. Where the record cannot show what the text requires, an interpretation says
   what is assumed.
6. **Source kinds and the copyright allow-list** of section 6.
7. **No double evaluation.** Two credentials never select the same record item; a shared
   evaluation used by several credentials is compiled for each, or once when it has its own
   scope (the flight review, one per pilot).
