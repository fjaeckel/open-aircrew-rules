# Changelog

All notable changes to the catalogue and the tools are listed here, written for pilots and
integrators. The format follows [Keep a Changelog](https://keepachangelog.com/en/1.1.0/);
catalogue changes are grouped by authority and name the credential or evaluation
(`<credential id>#<evaluation id>`).

## Unreleased — initial catalogue

### Added

#### EASA

- `easa.licence.ppl-a` PPL(A): the class 2 medical it needs, and the medical and language
  proficiency endorsement it requires.
- `easa.licence.spl` SPL: sailplane recency, launch-method recency (winch and car, aerotow,
  self-launch, bungee), passenger recency and the passenger prerequisite after licence
  issue, the LAPL medical and the medical it requires, and the SFCL.130 training course.
- `easa.rating.sep-land` SEP (land) class rating: revalidation (proficiency check, or 12 hours
  with the experience and refresher training; SEP land and TMG pooled on the same licence;
  ultralight time credited), passenger recency by day and by night, and its licence.
- `easa.rating.sep-sea` SEP (sea) and `easa.rating.tmg` TMG class ratings: the same
  revalidation and passenger evaluations.
- `easa.privilege.spl-tmg` TMG privileges of the SPL: recency, passengers, the passenger
  prerequisite and the extension training.
- `easa.privilege.sfcl-sailplane-towing` and `easa.privilege.sfcl-banner-towing` (SPL towing
  ratings) and `easa.privilege.cloud-flying` (sailplane cloud flying): recency and the
  licence they are held on.
- `easa.medical.class-1`, `easa.medical.class-2`, `easa.medical.lapl`: validity by age at
  examination.
- `easa.endorsement.language-proficiency`: validity of level 4 and level 5.
- Shared evaluations: `easa.shared.passengers-day`, `easa.shared.passengers-night`,
  `easa.shared.medical-certificate`, `easa.shared.sfcl-passenger-prerequisite`.

#### FAA

- `faa.licence.private-airplane` private pilot certificate (airplane): passenger recency by
  day and night, by class and by type, in tailwheel airplanes, the flight review, and the
  medical it requires.
- `faa.rating.instrument-airplane` instrument rating (airplane): instrument experience with
  the six-month grace and the IPC, the case of no instrument privileges, and the
  certificate it requires.
- `faa.medical.first-class`, `faa.medical.second-class`, `faa.medical.third-class`: duration
  for ATP, commercial and private privileges by age.
- `faa.medical.basicmed` BasicMed: the course, the comprehensive examination and the
  driver's license.
- Shared evaluations: `faa.shared.flight-review`, `faa.shared.passengers-day`,
  `faa.shared.passengers-night`, `faa.shared.passengers-tailwheel`,
  `faa.shared.medical-duration-private`, `faa.shared.medical-duration-commercial`.

#### Germany

- No credential yet. The 15 articles in scope (LuftPersV, LuftVZO) are accounted for in
  `coverage/articles.yaml`: 8 pending, 7 not evaluated with a reason.

#### Tools and format

- Composite results: per credential the record holds, whether it may be exercised on the
  date, combining its own evaluations with the credentials it requires and naming the
  member that decides (`credentials.Catalogue.Evaluate`; `cmd/evaluate` now prints
  `{ "evaluations": [...], "credentials": [...] }`).
- `policies.yaml`: every convention with no legal text behind it, with a statement and a
  rationale; the gate fails on a policy that is not declared or not cited.
- `fragments/` and `rulescheck -fragments` for contributors working in parallel.
- The credential format, its compiler onto the evaluation engine, and the gate
  (`cmd/rulescheck`) with reference resolution against `sources/`, worked examples,
  coverage of every article in scope (evaluated, pending or not evaluated) and the source
  copyright allow-list.
- `effective_from` / `effective_to` on evaluations, for regulation changes.
- `cmd/evaluate`: evaluate a record from the command line and print JSON.
- 113 verbatim source texts: 66 EU legal acts, 32 US federal regulations, 15 German
  statutes and ordinances.

### Changed

#### EASA

- `easa.licence.spl` sailplane recency, launch-method and passenger evaluations apply only
  to GLIDER ratings on an SPL or LAPL(S) licence of authority EASA or LBA (Part-SFCL); a
  GLIDER rating on any other licence, such as a PPL(A), is no longer evaluated as sailplane
  privileges (interpretation `sailplane-privileges-sfcl-licences` replaces
  `sailplane-privileges-any-licence`). `easa.privilege.cloud-flying` likewise selects only
  cloud-flying privileges on those licences.
- Ratings and privileges name the licence they need (`requires_any` / `requires_all`); the
  PPL(A) no longer lists its class ratings and the SPL no longer lists its privileges.

#### FAA

- The private pilot certificate no longer lists the instrument rating; the instrument
  rating requires the certificate instead.

#### Tools and format

- Worked examples state their expected `outcome` instead of `expect.status`. Passing
  examples are `current`, or `expiring` with the requirements met; every other outcome is
  failing (docs/credential-format.md).
- `requires:` is replaced by `requires_all:` and `requires_any:`, which are evaluated.
- Outcome presets cite a policy for each status they set by convention; a `statuses:`
  override is written `{ status, ref: policy:<id> }`. `credential_vocabulary.policy_refs`
  moved to `policies.yaml`.
- The gate measures statement coverage of the engine and the credentials package together.
