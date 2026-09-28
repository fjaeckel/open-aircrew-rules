# Owner decisions of 2026-09-28: checklist

Frederic Jung decided on 2026-09-28:

1. The five interpretation principles are confirmed (DESIGN.md section 14).
2. Every recommendation in `reviews/recommendations-easa.md`, `recommendations-faa.md` and
   `recommendations-de.md` is accepted: all 76 rows of the decision table in
   `reviews/README.md`, including the low-confidence ones.
3. The bugs found beyond the listed options are to be fixed.

This file lists, for each row and each bug, the format capability it needs and whether the
format has it now. The capabilities are described in DESIGN.md section 15 and
docs/credential-format.md. What is left is credential content: the credential files, their
worked examples and their interpretations (name the principle each follows with
`principle:`). New message keys, policies and coverage entries go through `fragments/`
while several contributors work at once (DESIGN.md section 12); an association document
not yet in `associations.yaml` is declared there.

Legend for "Available": **existing** means the format could already say it; **new** means a
capability added for these decisions; **none** means the row keeps the file's choice and needs
no change beyond the wording the recommendation names.

## EASA (38 rows)

| # | Credential | Interpretation | Accepted option | Capability | Available |
| --- | --- | --- | --- | --- | --- |
| 1 | `easa.medical.class-1` | `counted-from` | (2) age on the examination date | `valid_for.age_on: issue` (the period still runs from `counted_from`) | new |
| 2 | `easa.medical.class-2` | `counted-from` | (2) age on the examination date | `valid_for.age_on: issue` | new |
| 3 | `easa.medical.lapl` | `counted-from` | (2) age on the examination date | `valid_for.age_on: issue` | new |
| 4 | `easa.examiner.fcl-examiner` | `tests-as-examiner-time` | (1) keep | - | none |
| 5 | `easa.shared.instructor-certificate-validity` | `period-start` | (2) trust the recorded expiry; 3 years only without one | `valid_for.recorded_expiry_wins: true` | new |
| 6 | `easa.shared.passengers-night` | `ir-held` | (3) a valid IR of the same aircraft category | `holds: { classes: [IR], sameCategory: true, valid: true, expiryRecorded: true }`; an `unknown` stage for an IR of the category without an expiry (`sameCategory` without `expiryRecorded`). The IR's category is its recorded `category`, else the licence kind's (`licence_kinds.*.category`) | new |
| 7 | `easa.shared.sfcl-passenger-prerequisite` | `competence-flight-not-recorded` | (2) unknown until a training flight event is recorded | count `passenger_training_flight` (event kind of that name) with `unknown_if_none: true`, required (not informational) | new |
| 8 | `easa.shared.variants` | `same-engine-sep-variants-evaluated` | (2) only variants differing by engine type; unidentified ones unknown | variant record field `differentEngineType`; `only_for: { different_engine_type: true, if_missing: unknown }` | new |
| 9 | `easa.privilege.cloud-flying` | `dual-counts` | (2) PIC only for (e), dual only as the (f)(2) make-up | `as: [pic]` on the (e) rows; dual rows only in the (f)(2) branch | existing |
| 10 | `easa.privilege.cloud-flying` | `check-restores` | (1) keep | - | none |
| 11 | `easa.privilege.fcl-banner-towing` | `tows-counted` | (1) keep | - | none |
| 12 | `easa.privilege.fcl-sailplane-towing` | `tows-counted` | (1) keep | - | none |
| 13 | `easa.privilege.sfcl-banner-towing` | `tows-counted` | (2) banner tows only, any aircraft except ultralights | `tow_kinds: [banner]`, `excluding_classes: [ULTRALIGHT]` | existing |
| 14 | `easa.privilege.sfcl-sailplane-towing` | `tows-counted` | (2) sailplane tows only, any aircraft except ultralights | `tow_kinds: [glider]` (decide in the reading whether `ul_glider` counts), `excluding_classes: [ULTRALIGHT]` | existing |
| 15 | `easa.privilege.spl-tmg` | `twelve-hours-on-sailplanes` | (1) keep | - | none |
| 16 | `easa.rating.sep-land` | `examiner-not-checked` | (a)(2) an FFS check of the class counts; (b)(2) exclude checks recorded only for the IR or BIR | on the check: `simulator: include, fstd: [FFS]` (events now carry `fstdType`, flag-derived ones the flight's) and `excluding_ratings: [IR, BIR]` (events may list further `ratings`, so a combined check recorded for the class too still counts) | new |
| 17 | `easa.rating.sep-land` | `exemption-period` | (2) the validity period (24 months) | `waived_by: { ..., within_months_before_expiry: 24 }` | existing |
| 18 | `easa.rating.sep-land` | `land-sea-pooling-not-applied` | (2) a third route over both classes, per-class minima, both held on the same licence | a third `any_of` branch with `only_if: { holds: { classes: [SEP_LAND, SEP_SEA], sameLicence: true, valid: true, every: true } }`; totals with `classes: [SEP_LAND, SEP_SEA]`, minima with `classes: [SEP_LAND]` and `classes: [SEP_SEA]` | existing |
| 19 | `easa.rating.sep-land` | `ul-credit-time-only` | (2) no ultralight motorglider time unless the record shows a fixed engine and propeller | flight field `fixedEngine`; `ul_credit: { TMG: { ul_kinds: [THREE_AXIS_MOTORGLIDER], fixed_engine: true } }` | new |
| 20 | `easa.rating.sep-sea` | `no-land-sea-pooling` | (2) as for SEP land | as row 18 | existing |
| 21 | `easa.rating.bir` | `alternate-check-not-tracked` | (3) experience route only after a BIR check in an aeroplane in the 24 months before expiry; unknown without one | in the experience `all_of`: `proficiency_check: { for_rating: [BIR], categories: [aeroplane], within_months_before_expiry: 24, unknown_if_none: true }` | new |
| 22 | `easa.licence.cpl-a` | `class-1-only` | (2) usable within PPL privileges, commercial privileges limited | medical requirement `requires_any: [easa.medical.class-2, easa.medical.class-1@class_2]` deciding; `requires_any: [easa.medical.class-1@class_1]` (or the class 1 check) with `limits: commercial_privileges`. Needs the class 1 levels of bug E1. The shared `easa.shared.medical-certificate` checks the recorded expiry with `valid: true`; the licence files must not let it decide against a class 1 certificate still valid at class 2 | new |
| 23 | `easa.licence.cpl-h` | `class-1-only` | (2) | as row 22 | new |
| 24 | `easa.licence.cpl-as` | `class-1-only` | (2) | as row 22 | new |
| 25 | `easa.licence.atpl-a` | `class-1-only` | (2) | as row 22 (scope `commercial_privileges` covers CPL and ATPL privileges; `airline_transport_privileges` is also declared) | new |
| 26 | `easa.licence.atpl-h` | `class-1-only` | (2) | as row 25 | new |
| 27 | `easa.rating.helicopter-type` | `validity-period-is-twelve-months` | (2) an FSTD check counts toward the 2 hours | on the 2-hour row: `simulator: include` and `any_flight_of: [{ simulator: exclude }, { simulator: only, flagged: { proficiencyCheck: true } }]` | existing |
| 28 | `easa.licence.lapl-a` | `lapl-proficiency-check` | (2) only a check recorded for the LAPL(A) | `for_rating: [LAPL_A]` (say the rating value in the reading) | existing |
| 29 | `easa.licence.lapl-a` | `ul-credit-time-only` | (2) also toward the per-class hour | `ul_credit` on the per-class row | existing |
| 30 | `easa.licence.lapl-a` | `variants-sep-only` | (2) unknown unless the record marks a different engine type | as row 8 | new |
| 31 | `easa.licence.lapl-h` | `check-on-type` | (2) only a check in the helicopter | drop the simulator inclusion (default `simulator: exclude`) | existing |
| 32 | `easa.licence.spl` | `launch-as-pic` | (1) keep | - | none |
| 33 | `easa.licence.spl` | `training-course-for-students` | (2) evaluate (v) when the record says TMG privileges are sought | record `trainings: [{ programme: SPL, seeks: [TMG] }]`; (iv) rows `only_if: { not: { seeks: TMG } }`, (v) rows `only_if: { seeks: TMG }` | new |
| 34 | `easa.licence.spl` | `other-category-credit-shown-only` | (2) up to 7 hours toward the 15 hours, limits of (b)(1) and (b)(2) | `sum_of` (15 hours = instruction + the `sfcl-130b-credit` row, `max_hours: 7`); the credit is not added to the (a)(2)(ii), (iv)(B) and (v)(B) rows; launches: a `sum_of` with the credited launches capped `max: 10` | new |
| 35 | `easa.licence.gpl` | `class-or-type` | (1) refined: unknown for a gyroplane type privilege whose flights record no type | `in_type: true` on the type privilege (a flight without a designator is unknown input; a subject without one now gives unknown items, not a silent non-match) | existing |
| 36 | `easa.instructor.fi` | `valid-until-expiry` | (1) keep | - | none |
| 37 | `easa.instructor.fi` | `instruction-any-category` | (2) the category of the licence the FI is on; 20 hours for airships; gyroplane hours for a GPL | `in_category: true` (a privilege's category is its licence kind's); the 20-hour row with `only_if: { holds: { licenceKinds: [PPL_AS, CPL_AS], sameLicence: true } }` beside the 50-hour row with the negation | new |
| 38 | `easa.instructor.cri` | `instruction-in-aeroplanes` | (2) class-rating instruction only; untracked | `not_recorded` row in place of the instruction count | existing |

## FAA (11 rows)

| # | Credential | Interpretation | Accepted option | Capability | Available |
| --- | --- | --- | --- | --- | --- |
| 39 | `faa.instructor.flight-instructor` | `regime-by-recorded-expiry` | (1) keep | - | none |
| 40 | `faa.rating.instrument-airplane` | `any-powered-category` | (3) record the category on the IR rating and count by it | rating field `category`; `only_for: { categories: [aeroplane, helicopter, powered_lift, airship], if_missing: unknown }` and `in_category: true` in `counting` (the record's word is `aeroplane`) | new |
| 41 | `faa.rating.instrument-airplane` | `grace-from-last-met` | (2) `lapsed` with `policy:recency-lapsed` | the grace stage's status and ref | existing |
| 42 | `faa.shared.flight-review` | `who-and-endorsement-not-checked` | (2) only FAA checks (examiner, check airman, Armed Forces) | `by_authority: [FAA]` on the check and practical-test rows (event `authority`, flight `checkAuthority`; an event without one is unknown) | new |
| 43 | `faa.shared.flight-review` | `simulator-counts` | (2) exclude device sessions for the flight review | drop `simulator: include` from the `flight_review` rows (checks of (d)(1) keep theirs) | existing |
| 44 | `faa.shared.passengers-day` | `class-or-type` | (2) pool SEP with SET, land with land, sea with sea | `relevant_class: { class_group: faa, ref }` (vocabulary `class_groups.faa`); one result per group, subject `group` | new |
| 45 | `faa.shared.passengers-night` | `class-or-type` | (2) as by day | as row 44 | new |
| 46 | `faa.shared.passengers-night` | `logged-night-takeoffs` | (3) record take-offs in the 61.57(b) period; unknown meanwhile | count `night_period_takeoffs` (record `nightPeriodTakeoffs`; 0 on a flight without night take-offs, unknown on a night flight without the field, so the interim unknown is built in) | new |
| 47 | `faa.licence.powered-lift-airship` | `medical-grade-not-recorded` | (2) split by certificate grade; glider-only certificate unknown | as the rotorcraft files: evaluations `only_for: { licence_kinds: [...] }` per grade, and an `unknown` evaluation for `FAA_GLIDER` | existing |
| 48 | `faa.licence.sport` | `sport-classes` | (2) add powered-parachute and weight-shift-control classes | classes `POWERED_PARACHUTE_LAND`, `POWERED_PARACHUTE_SEA`, `WEIGHT_SHIFT_CONTROL_LAND`, `WEIGHT_SHIFT_CONTROL_SEA` (categories `powered_parachute`, `weight_shift_control`) | new |
| 49 | `faa.licence.sport` | `medical-or-drivers-license` | (3) a driver's-license credential, required with a medical as alternative for powered categories | credential kind `document` (`credentials/faa/documents/<name>.yaml`, selecting `US_DRIVERS_LICENSE`); the requirement as an evaluation over the classes held (`only_for: { when_holding: ... }` / `holds`) or `requires_any` naming the document and the medicals; 61.23(c)(2) stated as a limitation or interpretation | new |

## Germany (27 rows)

| # | Credential | Interpretation | Accepted option | Capability | Available |
| --- | --- | --- | --- | --- | --- |
| 50 | `de.licence.ul` | `ultralight-licences` | (2) class 2 or LAPL medical for three-axis, helicopter, gyroplane, weight-shift | `requires_any` in the kind rating files (`easa.medical.class-2`, `easa.medical.lapl`, `easa.medical.class-1`, or `@class_2` once row 22's levels exist) | existing |
| 51 | `de.instructor.ul-instructor` | `three-years-from-valid-from` | (2) German periods end one day earlier | `authority_conventions.DE.validity_ends: day_before` (vocabulary, not per file; §§ 186, 187(2), 188(2) BGB). Already applied: the examples moved by one day and keep their outcomes; the reading and its `ref` still need updating | new, applied |
| 52 | `de.instructor.ul-instructor` | `last-three-years-before-expiry` | (1) keep | - | none |
| 53 | `de.instructor.ul-instructor` | `instruction-as-instructor-or-examiner` | split: 60 take-offs and 60 landings; the rating's kinds when recorded, otherwise unknown | separate `takeoffs` and `landings` rows (min 60); `in_ul_kind: true` (the kinds the privilege `detail` names; none named gives unknown items) | new |
| 54 | `de.instructor.ul-instructor` | `befaehigungspruefung-event` | (1) keep | - | none |
| 55 | `de.privilege.ul-towing` | `tow-kind-not-recorded` | (2) first limb: record the entered kind and filter by it; ultralight tows only | `by_this_tow_kind: true` and `by_this_tow_take_up: true` (privilege `detail` names the kind, e.g. `banner pick_up`); flight `towTakeUp`; tow kind `hang_glider`. A detail naming no kind makes every tow unknown input (stricter than "unknown only with tows of several kinds"; say so in the reading) | new |
| 56 | `de.privilege.ul-towing` | `tows-counted` | (2) count tow flights | `landings: { flagged: { towFlight: true } }` or `flights` with the flag | existing |
| 57 | `de.shared.ul-passengers` | `authorisation-recorded` | (2) compare the privilege's detail with the kind flown; unknown without one | `holds: { privileges: [UL_PASSENGER_AUTH], details: [$subject], sameLicence: true }` for current, `details: [none]` for unknown | new |
| 58 | `de.shared.ul-passengers` | `authorisation-progress-informational` | (1) keep; correct the § 84a(3) sentence | - | none |
| 59 | `de.shared.ul-passengers` | `xc-200-km-per-flight` | (1) keep per flight; reword the reading | - | none |
| 60 | `de.rating.ul-three-axis` | `scope-of-closing-words` | (1) keep | - | none |
| 61 | `de.rating.ul-three-axis` | `training-flight-as-dual` | (2) also a flight flagged with an instructor aboard | count `longest_flight` with `any_flight_of: [{ with_time: [dual] }, { flagged: { instructorOnBoard: true } }]` (the flag already exists) | new |
| 62 | `de.rating.ul-three-axis` | `proficiency-check-period` | (1) 24 months | - | none |
| 63 | `de.rating.ul-three-axis` | `medical-above-120-kg-not-evaluated` | (2) require a class 2 or LAPL medical | `requires_any` (as row 50) | existing |
| 64 | `de.rating.ul-weight-shift` | `association-rules-by-authority` | (1) keep; say LBA-to-DULV is unverified | - | none |
| 65 | `de.rating.ul-weight-shift` | `dulv-pic-only` | (2) 12 hours PIC or a check on a weight-shift ultralight | `any_of` with `proficiency_check: { ul_kinds: [WEIGHT_SHIFT] }` | existing |
| 66 | `de.rating.ul-weight-shift` | `daec-safety-training-not-recorded` | (2) 12 h, 12 take-offs, 12 landings, the training (unknown until recorded), or a check | count `safety_training` (event kind) with `unknown_if_none: true`, a `takeoffs` row, a check branch | new |
| 67 | `de.rating.ul-weight-shift` | `above-120-kg-not-named` | (1) keep; reading names the unclear basis | - | none |
| 68 | `de.rating.ul-gyroplane` | `association-rule-unverified` | (2) per issuer | two evaluations with `only_for: { authorities: [...] }` (as weight-shift); declare the DAeC gyroplane document in `associations.yaml` | existing |
| 69 | `de.rating.ul-gyroplane` | `gyroplane-class-credited` | (2) ultralight gyroplanes only, count take-offs, check on an ultralight gyroplane | drop the `GYROPLANE` class credit; `takeoffs` row; check with `ul_kinds: [GYROPLANE]` | existing |
| 70 | `de.rating.ul-gyroplane` | `above-120-kg-not-named` | (1) as weight-shift | - | none |
| 71 | `de.rating.ul-helicopter` | `training-flight-as-dual` | (2) as three-axis | as row 61 | new |
| 72 | `de.rating.ul-helicopter` | `proficiency-check-period` | (1) 12 months | - | none |
| 73 | `de.rating.ul-powered-paraglider` | `association-rule-unverified` | (2) 30 take-offs and 30 landings, or a check | `takeoffs` and `landings` rows; check branch | existing |
| 74 | `de.rating.ul-powered-paraglider` | `above-120-kg-not-named` | (1) as weight-shift | - | none |
| 75 | `de.rating.ul-sailplane` | `association-rule-unverified` | (2) 5 take-offs and 5 landings, any launch method, DAeC rule for every issuer | `takeoffs` and `landings` rows | existing |
| 76 | `de.rating.ul-sailplane` | `above-120-kg-not-named` | (1) as weight-shift | - | none |

Totals: 76 rows; 20 keep the file's choice (none), 23 need only existing words, 33 needed a
new capability, all of which is now available (row 51 already applied).

## Bugs beyond the listed options

| # | Source | Bug | Capability | Available |
| --- | --- | --- | --- | --- |
| E1 | recommendations-easa, "Class 1 validity for lower privileges" | A class 1 certificate past its class 1 period is reported invalid for class 2 and LAPL privileges (AMC1 MED.A.030) | `level:` on the medical's evaluations (`class_1`, `class_2`, `lapl`, each a `valid_for` with that class's periods), requirements `easa.medical.class-1@class_2` / `@lapl` in the PPL, LAPL, SPL, GPL, BPL and CPL/ATPL private-privilege requirements; composites list `levels` | new |
| E2 | recommendations-easa, "FI(As) and FI(G) hours" | FI(As) held to 50 hours instead of 20 (FCL.940.FI(a)(1)(i)(B)) | as row 37: `in_category: true` and the 20-hour row under `only_if: { holds: { licenceKinds: [...], sameLicence: true } }` | new |
| E3 | recommendations-easa, "CRI single-engine and multi-engine split" | FCL.940.CRI(a)(1) equal split not checked | under row 38 the 10 hours are untracked, so the split is part of the untracked row; if it is ever tracked: rows with `max_engines: 1` / `classes: [MEP_LAND, MEP_SEA]` under `only_if: { holds: { privileges: [CRI], details: [ME] } }` (`details`) | existing / new |
| E4 | recommendations-easa, "SFCL.160 month-end counting" | 24 months counted to the end of the month of the training flight (AMC1 SFCL.160(a)(1)(ii)(d)) | `within_calendar_months: 24` on the training-flight row: it is month-end anchoring (engine test `TestMonthEndWindow`) | existing |
| E5 | easa-class-ir-ratings, "Other findings" | IR-only and BIR-only checks count as class rating checks (SEP land and users, MEP land, SET land) | `excluding_ratings: [IR, BIR]`; event `ratings` for combined checks | new |
| E6 | CHANGELOG / easa-shared-privileges | A type rating recorded without a type designator gets no passenger result | `only_for: { type_rated: $type_rated, if_missing: unknown }` in `easa.shared.passengers-day` and `-night` (the EASA type ratings are all type-rated; FAA class ratings legitimately have no designator, so the FAA shared files must not add it) | new |
| F1 | recommendations-faa, `ipc-restores` | The IPC's aircraft category is not checked (61.57(d)(2)) | `restored_by: [{ ipc: { in_category: true, simulator: include }, ref }]` (event `category`, else its class's) | new |
| F2 | recommendations-faa, flight review | A review given by a non-part 61 instructor counts (61.56(c)(1)) | `by_authority: [FAA]` on the `flight_review` rows | new |
| F3 | recommendations-faa, passengers-day | FAA single-engine land and sea are separate classes; pooling must not merge them | `class_groups.faa` keeps land and sea apart | new |
| F4 | recommendations-faa, Williams letter | Letter date from the FAA file title | none (reviews wording) | none |
| F5 | faa-shared-licences observation 4 | A first- or second-class medical past its commercial or ATP duration reads expired for the certificate although it supports lower privileges | `level:` on `faa.medical.first-class` / `second-class` (`atp`, `commercial`, `private`) and `@private` / `@commercial` in the certificates' requirements | new |
| F6 | faa-shared-licences observation 3 | `holds` has no authority filter (a non-FAA `PRIVATE` ends the student exemption) | `holds.authorities: [FAA]` | new |
| F7 | faa-instructors observations 1-2 | CAT II/III requirement lists incomplete | `requires_any` edits | existing |
| F8 | faa-shared-licences observation 2 | Self-launching gliders logged as TMG get no glider passenger credit | not addressed by an accepted recommendation; a `class_groups` set can pool classes once the owner decides | existing (no decision) |
| D1 | recommendations-de, finding 1 | Check alternatives missing (weight-shift, powered paraglider) | as rows 65, 66, 73 | existing |
| D2 | recommendations-de, finding 2 | Take-offs not counted | as rows 66, 69, 73, 75 | existing |
| D3 | recommendations-de, finding 3 | Part-FCL gyroplane credit | as row 69 | existing |
| D4 | recommendations-de, finding 4 | Three-axis LL up to 120 kg has its own association rule | the finding defers it to a separate decision once the record can tell an LL entry from a UL entry; no record field was added | not implemented (needs a decision) |
| D5 | recommendations-de, finding 5 | Medical in the association texts | as rows 50, 63 | existing |
| D6 | recommendations-de, finding 6 | § 186 BGB not stored | `sources/de/bgb-186.md` added; cited by `authority_conventions.DE` | done |

## Notes for the content work

- Name the principle each new or changed interpretation follows (`principle: P1` ...).
- Row 51 is already in effect; update `three-years-from-valid-from` to describe it and cite
  `de:BGB.186`.
- `sum_of` rows, `unknown_if_none` rows and `if_missing: unknown` evaluations each need a
  failing worked example showing the unknown or unmet case (gate rule, DESIGN.md section 8).
- New keys: `selection.input_missing`, `requirement.passenger_training_flight` and
  `requirement.safety_training` exist. Limitation scopes: `commercial_privileges`,
  `airline_transport_privileges`, `night_privileges`.
