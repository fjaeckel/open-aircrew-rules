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
- `easa.rating.sep-land` and `easa.rating.sep-sea`: variant recency (FCL.710(d)) through
  `easa.shared.variants`, with the FCL.710 legacy cases as worked examples.
- `easa.rating.mep-land` and `easa.rating.mep-sea` MEP class ratings: revalidation (proficiency
  check or EBT assessment in the 3 months before expiry, plus 10 route sectors or 1 with an
  examiner), passenger recency by day and by night, variants, and their licence.
- `easa.rating.set-land` and `easa.rating.set-sea` SET class ratings: revalidation by
  proficiency check, passenger recency, variants, and their licence.
- `easa.rating.ir-a`, `easa.rating.ir-h` and `easa.rating.ir-as` instrument ratings:
  revalidation by IR proficiency check (aeroplane, helicopter or airship category), and the
  licence each is held on.
- `easa.rating.bir` basic instrument rating: 1-year validity and revalidation (check, or 6 hours
  IFR as PIC with 3 approaches and a training flight), and its licence.
- `easa.rating.aeroplane-type` aeroplane type ratings (class OTHER): validity by the recorded
  expiry date, and variants.
- `easa.rating.helicopter-type`, `easa.rating.powered-lift-type`, `easa.rating.airship-type`:
  type rating revalidation (FCL.740.H, FCL.740.PL, FCL.740.As), passenger recency, variants,
  and the licence each is held on (a helicopter or airship licence; a professional aeroplane or
  helicopter licence for powered-lift types, since Part-FCL has no powered-lift licence).
- `easa.licence.cpl-a`, `easa.licence.atpl-a`, `easa.licence.mpl`: class 1 medical and language
  endorsement.
- `easa.licence.ppl-h`, `easa.licence.cpl-h`, `easa.licence.atpl-h`, `easa.licence.ppl-as` and
  `easa.licence.cpl-as` helicopter and airship licences: the class 2 (private) or class 1
  (commercial, airline transport) medical and the language endorsement.
- `easa.licence.lapl-a`: LAPL(A) recency (FCL.140.A, with the SEP land/sea split and ultralight
  credit), the 10-hour passenger prerequisite (FCL.105.A(b)), passenger recency, SEP variants,
  LAPL medical and language endorsement.
- `easa.licence.lapl-h`: LAPL(H) recency per type (FCL.140.H), passenger recency, variants, LAPL
  medical and language endorsement.
- `easa.licence.gpl` GPL: recency per gyroplane rating (FCL.240.G(a)) with ultralight credit for
  gyroplanes of at least 450 kg (FCL.035(a)(5)), SPG variant recency (FCL.240.G(b)), the
  passenger prerequisite after issue (FCL.205.G(a)(2)), passenger recency by day and night
  (FCL.060), the class 2 medical (MED.A.030(c)(2)) and the language endorsement.
- `easa.licence.bpl` BPL: the LAPL medical (MED.A.030(c)(1)); its recency (BFCL.160) is reported
  unknown because the record holds no balloon class or group, so the licence is never reported
  usable on a guess.
- `easa.privilege.fcl-sailplane-towing`, `easa.privilege.fcl-banner-towing`: Part-FCL towing
  recency (FCL.805(e)); `easa.privilege.mountain`: mountain rating recency (FCL.815(d)).
- `easa.privilege.sfcl-launch-methods` SPL launch method training records (SFCL.155(a), (b)),
  requiring the SPL; the per-method recency stays with `easa.licence.spl`.
- `easa.instructor.fi` and `easa.instructor.iri` FI and IRI certificates: 3-year validity
  (FCL.940) with revalidation, two of 50 hours of instruction, instructor refresher training
  and an assessment of competence in the last 12 months (FCL.940.FI, FCL.940.IRI).
- `easa.instructor.cri` CRI certificate: 3-year validity with revalidation, two of 10 hours of
  instruction in aeroplanes, refresher training and an assessment of competence (FCL.940.CRI).
- `easa.instructor.tri`, `easa.instructor.sfi`, `easa.instructor.mcci`, `easa.instructor.sti`
  and `easa.instructor.fti`: the 3-year validity of FCL.940 through the shared evaluation
  `easa.shared.instructor-certificate-validity`; their revalidation needs data the record does
  not hold and is stated as an interpretation.
- `easa.instructor.mi` MI certificate: valid while an FI, TRI or CRI certificate on the same
  licence is (FCL.940.MI), requiring one of those credentials.
- `easa.instructor.fi-s` FI(S) certificate: SFCL.360 recency (refresher training and 30 hours
  or 60 launches of instruction in 3 years, a demonstration of instructional ability in 9
  years, resumption by an assessment of competence).
- `easa.examiner.fcl-examiner` Part-FCL examiner certificates: 3-year validity and revalidation
  (six tests conducted, refresher course and an assessed test in the last 12 months,
  FCL.1025).
- `easa.examiner.fe-s` FE(S) certificate: 5-year validity and revalidation (SFCL.460).
- The instructor and examiner certificates require the licence they are held on (the SFI,
  STI and MCCI excepted: their holders need only have held a licence).
- Shared evaluations: `easa.shared.passengers-day`, `easa.shared.passengers-night`,
  `easa.shared.medical-certificate`, `easa.shared.sfcl-passenger-prerequisite`,
  `easa.shared.variants` (variant not flown within 2 years, FCL.710(d), FCL.140.A(c)),
  `easa.shared.instructor-certificate-validity`.
- Verbatim texts of Part-FCL FCL.720.PL, FCL.1000 and FCL.1020 (EUR-Lex consolidated
  version 30.04.2026) under `sources/easa/`, the evidence for the powered-lift licence
  prerequisite and the examiner readings.
- `easa.endorsement.language-proficiency#level_6`: an expert level (6) endorsement is now
  selected and reported valid (it is never re-evaluated) unless a recorded validity date has
  passed.

- Reviews of every interpretation of the EASA shared evaluations and privileges
  (`reviews/easa-shared-privileges.md`), the class ratings, instrument ratings and BIR
  (`reviews/easa-class-ir-ratings.md`), the type ratings and licences
  (`reviews/easa-types-licences.md`) and the SPL, GPL and instructor certificates
  (`reviews/easa-spl-gpl-instructors.md`); none is a sign-off.
- Verbatim texts under `sources/easa/`: FCL.600 and Appendix 9, Section A, points 1 to 1f
  (an extract); FCL.900 and FCL.915; FCL.205.As, FCL.305, FCL.405.A, FCL.505, FCL.700,
  FCL.720.H and FCL.720.As; FCL.935, FCL.915.SFI, FCL.915.MCCI, FCL.915.STI, FCL.930.MCCI,
  SFCL.345 and Article 3b of Regulation (EU) 2018/1976.
- `easa.licence.lapl-a#recency_tmg`: recency of the LAPL(A) TMG privilege (FCL.140.A(a)),
  without the SEP land/sea split of FCL.140.A(b).

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
- `faa.licence.commercial-airplane` commercial pilot certificate (airplane) and
  `faa.licence.atp-airplane` airline transport pilot certificate (airplane): passenger
  recency by day and night, by class and by type, in tailwheel airplanes, the flight review,
  and the medical each requires (first or second class; first class).
- `faa.licence.private-rotorcraft`, `faa.licence.commercial-rotorcraft` and
  `faa.licence.atp-rotorcraft` (helicopter and gyroplane ratings): passenger recency by day
  and night, by class and by type, the flight review, and the medical each requires.
- `faa.licence.powered-lift-airship` powered-lift and airship ratings on a private,
  commercial, airline transport or glider pilot certificate: passenger recency by day and
  night, by class and by type, the flight review and the medical.
- `faa.licence.glider` glider category rating on a glider, private or commercial pilot
  certificate: passenger recency (3 launches and 3 landings as sole manipulator) and the
  flight review; no medical. Airplane and rotorcraft ratings recorded on a glider pilot
  certificate are evaluated as private pilot ratings.
- `faa.licence.recreational` recreational pilot certificate: 61.101(g) pilot-in-command
  recency below 400 hours (restored by an instructor's endorsement), day and tailwheel
  passenger recency, no night privileges, the flight review and the medical.
- `faa.licence.sport` sport pilot certificate: passenger recency by class, type, glider and
  tailwheel; night passenger recency from 22 October 2025 (61.329, with the night
  endorsement) and not applicable before; the flight review.
- `faa.licence.student` student pilot certificate: the 90-day solo endorsement (61.87(n)),
  the flight review exemption and the medical.
- `faa.instructor.flight-instructor` flight instructor certificate (sport pilot rating
  included): expiry of certificates issued before 1 December 2024 (61.19(d)(2)) and, from
  1 December 2024, recent experience within 24 calendar months with the 61.199 / 61.427
  reinstatement steps.
- `faa.instructor.ground-instructor` ground instructor recent experience (61.217).
- `faa.privilege.category-ii-iii` Category II and III pilot authorizations: six-calendar-month
  validity (61.21(a)).
- `faa.privilege.glider-towing` glider and unpowered ultralight vehicle towing: recency
  within 24 calendar months (accompanied tows or towed glider flights, 61.69(a)(6)).
- `faa.privilege.sport-night` sport pilot night endorsement: the sport pilot certificate and
  the medical certificate or BasicMed it requires (61.329(b)).
- Shared evaluations: `faa.shared.flight-review`, `faa.shared.passengers-day`,
  `faa.shared.passengers-night`, `faa.shared.passengers-tailwheel`,
  `faa.shared.passengers-glider`, `faa.shared.medical-duration-private`,
  `faa.shared.medical-duration-commercial`.
- `faa.privilege.glider-towing`: composite examples for a towing privilege on an airline
  transport pilot certificate with a helicopter rating (current) and on a certificate with
  only a glider rating (unknown, no powered category).
- Sources: 14 CFR 61.67 and 61.68 (Category II and III pilot authorization requirements),
  and excerpts (dates and amendatory instructions) of the final rules 89 FR 80020
  (flight instructor certificates, effective 1 December 2024) and 90 FR 35034 (sport pilot
  night privileges, effective 22 October 2025).
- Review of every FAA instructor, privilege, rating and medical interpretation:
  `reviews/faa-instructors-privileges-ratings-medicals.md`.

- Review of every interpretation of the FAA shared evaluations and pilot certificates:
  `reviews/faa-shared-licences.md`.
- Sources: 14 CFR 1.1 (definitions, including night, category, class and glider), 61.1
  (definitions, including examiner), 61.5 (certificates and ratings issued) and 61.157
  (airline transport pilot practical test ratings) (`sources/faa/`).

#### Germany

- `de.licence.ul` ultralight pilot licence (Luftfahrerschein für Luftsportgeräteführer): the
  training hours for three-axis (§ 42(4) no. 1) and weight-shift (§ 42(5) no. 1) ultralights,
  and an ultralight rating recorded without a kind (§ 44(2)), reported unknown.
- `de.rating.ul-three-axis` (§ 45(2), (3)) and `de.rating.ul-helicopter` (§ 45(2a), (3)):
  recency of the statute; `de.rating.ul-gyroplane`, `de.rating.ul-weight-shift`,
  `de.rating.ul-powered-paraglider` and `de.rating.ul-sailplane`: the recency the
  associations set under § 45(4) (DULV, DAeC; cited by title, not verified). Each kind has
  passenger recency (§ 45a, with the § 84a authorisation) through `de.shared.ul-passengers`
  and requires the licence.
- `de.privilege.ul-towing` ultralight towing rating: 10 tows in 24 months (LuftPersV § 84(5)),
  and its licence.
- `de.privilege.ul-passenger-authorisation` ultralight passenger authorisation: its validity,
  which follows the licence unless an expiry date is entered (§ 84a(5)), and its licence.
- `de.instructor.ul-instructor` ultralight instructor rating: three years from valid-from,
  extended by two of instruction given, a refresher course and an assessment of competence
  in the last three years (§ 96(1), (4)), and its licence.
- Verbatim sources for the review of the German ultralight licence, privileges and instructor
  rating: LuftPersV § 95a (instructor rating) and § 122 (repealed, referred to by § 84a(2)),
  BGB §§ 187 and 188 and VwVfG § 31 (how periods are counted).
- `reviews/de-ul-licence-privileges-instructor.md`: first-pass review of the 25
  interpretations of `de.licence.ul`, `de.instructor.ul-instructor`,
  `de.privilege.ul-passenger-authorisation`, `de.privilege.ul-towing` and
  `de.shared.ul-passengers` (15 supported, 10 needing a decision, none contradicted; no
  evaluation result changed).

- `reviews/de-ul-ratings.md`: first-pass review of the 40 interpretations of the six
  ultralight kind ratings (22 supported, 17 needing a decision, 1 contradicted and fixed).

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
- `ul_kinds` in `selects.ratings`, `only_for` and `scope` (with `none` for a rating recorded
  without a kind); the overlap check treats ratings of disjoint ultralight kinds as distinct.
- Association references `assoc:<publisher>:<document>`, declared in `associations.yaml`
  (publisher, title, delegating statute, `verified`; no text stored) and cited next to the
  delegating statute.
- Count words `ifr_time`, `practical_test`, `assessment_of_competence`,
  `instructor_refresher`, `examiner_refresher`, `supervised_instruction`,
  `differences_training`, `solo_endorsement`, `instructor_endorsement`.
- `ul_credit` takes a minimum mass: `{ CLASS: { ul_kinds: [...], min_mtom_kg: n } }`.
- Licence kinds `PPL_H`, `CPL_H`, `ATPL_H` (replacing `HELICOPTER`), `PPL_AS`, `CPL_AS`.
- `cmd/evaluate`: evaluate a record from the command line and print JSON.
- 147 verbatim source texts: 87 EU legal acts, 40 US federal regulations, 20 German
  statutes and ordinances.
- `reviews/README.md`: an index of the review files and the owner's list of every
  interpretation that needs a decision; the README gives the review status per authority.

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
- `easa.endorsement.language-proficiency`: a licence whose holder has only a level 6
  endorsement recorded was reported unknown (endorsement not held); it now meets the
  language requirement (FCL.055(c)).
- `easa.examiner.fe-s#revalidation`: the demonstration of examiner ability counts only an
  assessment of competence recorded for the rating FE_S; an assessment of competence as a
  pilot or instructor no longer revalidates the FE(S) certificate (SFCL.460(b)(2)).
- `easa.rating.powered-lift-type#licence`: the rating requires a CPL or ATPL (aeroplanes or
  helicopters), the licences FCL.720.PL names; an MPL no longer satisfies it, and the
  interpretation `part-fcl-type-ratings` now cites FCL.720.PL instead of being marked
  unverified.
- `easa.examiner.fcl-examiner`: interpretation `not-applied` cites FCL.1025(b) (where the
  combined revalidation of several examiner categories is written) instead of
  FCL.1025(b)(3)(ii).

- `easa.licence.lapl-a#recency` now covers the SEP land and SEP sea privileges only, and the
  new `easa.licence.lapl-a#recency_tmg` evaluates the TMG privilege. The land/sea split of
  FCL.140.A(b) (1 hour and 6 take-offs and landings in each SEP class) no longer applies to
  the TMG privilege: a holder of SEP land, SEP sea and TMG privileges who flew the 12 hours
  in a TMG is now current for the TMG privilege (review `reviews/easa-types-licences.md`).
- `easa.rating.aeroplane-type` interpretation `revalidation-not-evaluated`: the reading no
  longer suggests a 2-year validity for type ratings; FCL.740(a)(1) sets 1 year for type
  ratings unless the OSD determines otherwise (wording only, no result changes).
- `easa.rating.set-land#revalidation` (and `easa.rating.set-sea#revalidation`, which uses
  it): a proficiency check recorded as a flight in an FFS of the class, flagged as a
  proficiency check, now revalidates the rating (Appendix 9, Section A, points 1c and 1e,
  allow the class rating check in an FFS). A check event or a check in another FSTD recorded
  in a simulator still does not count. Interpretation `check-in-aeroplane` rewritten; new
  requirement row `ffs_check` (review `reviews/easa-class-ir-ratings.md`).
- `easa.rating.sep-land`, `easa.rating.sep-sea` and `easa.rating.tmg` `#licence`: the rating
  is held on a PPL(A), CPL(A), ATPL(A) or MPL, as for the MEP and SET ratings. A rating on a
  CPL(A), ATPL(A) or MPL was reported unknown (licence not held) however valid that licence
  was; it is now decided by the licence it is recorded on (fix deferred by
  `reviews/easa-class-ir-ratings.md`, applied at integration; new composite example
  `composite-rating-on-cpl`).
- `easa.shared.passengers-day` and `easa.shared.passengers-night` take a parameter
  `type_rated`. For the type ratings (`easa.rating.helicopter-type`,
  `easa.rating.powered-lift-type`, `easa.rating.airship-type`) and the LAPL(H) types,
  passenger recency is evaluated per type designator and counts only flights in that type
  (FCL.060(b)(1), (b)(2)(i)); previously flights in any type of the class counted. A type
  rating recorded without a type designator gets no passenger result, and a flight in the
  class without one is unknown input. The class ratings, the LAPL(A) and the GPL are
  unchanged. Interpretation `class-for-type` reworded; examples
  `passengers-day-other-type` added (fix deferred by `reviews/easa-shared-privileges.md`,
  applied at integration).
- `easa.shared.sfcl-passenger-prerequisite` (the SPL and the SPL TMG privileges): an FI(S)
  certificate meets the passenger prerequisite whatever its recorded expiry date, as
  Part-SFCL gives the FI(S) no validity period; whether its privileges may be exercised
  (SFCL.360, events recorded for FI_S) is evaluated by `easa.instructor.fi-s`.
  Interpretation `fi-s-on-any-licence` reworded; example
  `passenger-prerequisite-fi-s-recorded-expiry-passed` added (deferred item of
  `reviews/easa-spl-gpl-instructors.md`, applied at integration).
- `easa.privilege.spl-tmg#extension_training`: TMG privileges on a LAPL(A) no longer credit
  the TMG extension of the SPL. SFCL.150(c)(1) credits a TMG class rating, which only a
  PPL(A), CPL(A), ATPL(A) or MPL carries; the LAPL(A) route of SFCL.150(c)(2) also needs the
  FCL.140.A recency, which is not checked, so the training is evaluated instead. Interpretation
  `extension-credit` reworded; two worked examples added (review
  `reviews/easa-shared-privileges.md`).
- `easa.shared.variants` (used by the SEP, SET and MEP class ratings, the aeroplane,
  helicopter, airship and powered-lift type ratings, the LAPL(A) and the LAPL(H)): refresher
  training now restores only a variant of an SEP land or SEP sea class rating, as
  FCL.710(d)(3) allows it only for SEP variants with a particular type of engine; variants of
  other ratings are restored by differences training or a proficiency check. Interpretation
  `refresher-any-variant` replaced by `refresher-sep-variants`. No worked example changes its
  result; example `variants-refresher-does-not-restore` (MEP land) added at integration.
- `easa.instructor.fi#revalidation` (also used by `easa.instructor.iri`): refresher training
  and the assessment of competence count only when recorded for the FI or the IRI
  certificate (event `rating: FI` or `IRI`). Refresher training recorded for an FI(S) or an
  assessment of competence recorded for an examiner certificate no longer revalidates an FI
  (FCL.940.FI(a)(1)(ii), (iii); FCL.935) (review `reviews/easa-spl-gpl-instructors.md`).
- `easa.instructor.cri#revalidation`: refresher training and the assessment of competence
  count only when recorded for the CRI certificate (event `rating: CRI`)
  (FCL.940.CRI(a)(2), (3)).
- `easa.instructor.fi-s#recency`: refresher training, the demonstration of instructional
  ability and the assessment of competence count only when recorded for the FI(S)
  certificate (event `rating: FI_S`), so an FE(S) or Part-FCL instructor event no longer
  keeps an FI(S) current (SFCL.360(a), (c), (d); SFCL.345).
- `easa.instructor.fi-s#recency`: a recorded expiry date no longer reports the FI(S)
  expired. Part-SFCL sets no validity period; for a certificate converted from Part-FCL,
  Article 3b(2)(c) of Regulation (EU) 2018/1976 lets the holder instruct after the endorsed
  expiry date while complying with SFCL.360, so the recency alone decides (current or
  lapsed).
- `easa.instructor.mcci` interpretation `revalidation-not-evaluated`: the reading now cites
  the stored FCL.930.MCCI(a)(3) (wording only, no result changes).

#### FAA

- The private pilot certificate no longer lists the instrument rating; the instrument
  rating requires the certificate instead (a private, commercial or airline transport
  airplane certificate).
- The student solo endorsement, the flight instructor practical test and refresher course and
  the ground instructor refresher course and endorsement are counted with their own count
  words instead of stand-ins; no result changes.
- `faa.instructor.flight-instructor#recent_experience`: interpretation `other-means-not-recorded`
  now reads 61.197(a)(1) as it is written (the 24 calendar months may start from the month
  the FAA issued the certificate, whatever led to the issue) and says that the recorded
  issue date is not counted, so a certificate issued in the last 24 calendar months reads
  lapsed until the test or refresher course behind it is recorded. Three worked examples
  that reported a certificate issued in December 2024 as lapsed or expiring in 2026 (when
  61.197(a)(1) keeps it current through 31 December 2026) now use dates in 2027; their
  outcomes are unchanged.
- Worked examples of `faa.instructor.flight-instructor` and `faa.medical.second-class` name
  evaluations by their current ids instead of retired rule ids.

- `faa.shared.flight-review` (used by every FAA pilot certificate credential): the student
  pilot exemption of 61.56(g) now applies only while the pilot holds no other FAA pilot
  certificate. The flight review is evaluated once, on the first FAA certificate in the
  record; when that was a student certificate, a pilot who also held a private (or higher)
  certificate got `not_applicable` and needed no flight review. Interpretation
  `student-exemption` reworded; worked example
  `faa.licence.student` / `flight-review-student-certificate-listed-first` added. No other
  worked example changes its result (review `reviews/faa-shared-licences.md`).
- `faa.shared.medical-duration-commercial`: interpretation `commercial-glider-balloon`
  reworded. Commercial balloon privileges for compensation or hire need at least a
  second-class certificate (61.23(a)(2)(iii)); only balloon flight training and glider
  privileges need none. No result changes.
- `faa.licence.glider`: interpretation `powered-classes-elsewhere` reworded. A self-launching
  glider is a glider under 14 CFR 1.1; the reading no longer calls a motorglider "not a
  glider", and states that its flights count here only when logged in the glider class. No
  result changes.

#### EASA and Germany

- Instructor and examiner refresher training, assessments of competence and demonstrations
  are counted with their own count words; the untracked informational rows that stood for
  them are gone and no example changes status. The BIR counts IFR time with `ifr_time`.
- `de.instructor.ul-instructor#extension`: the Befähigungsprüfung of § 96(4) item 3 is an
  assessment of competence; a proficiency check taken for a pilot rating no longer counts.

#### Germany

- `de.rating.ul-three-axis#passengers`: a PPL(A) or SPL valid on the date evaluated no longer
  counts as the passenger authorisation. LuftPersV § 84a(2) second sentence deems the
  authorisation granted on the day the ultralight licence is issued to a pilot who then holds
  a valid PPL or SPL, which the record cannot show; a deemed authorisation counts when it is
  recorded as a passenger authorisation on the ultralight licence, and otherwise the result is
  unknown and asks for it (previously current for any pilot with a PPL(A) or SPL valid today,
  including one obtained after the ultralight licence) (review `reviews/de-ul-ratings.md`).

#### Tools and format

- Paragraph labels inserted by amendment (`(2a)`, `(5a)`, `(da)`) resolve as siblings of the
  label they follow: `de:LuftPersV.45(2a)` and `easa:FCL.710(da)` are now cited exactly.
- Worked examples state their expected `outcome` instead of `expect.status`. Passing
  examples are `current`, or `expiring` with the requirements met; every other outcome is
  failing (docs/credential-format.md).
- `requires:` is replaced by `requires_all:` and `requires_any:`, which are evaluated.
- Outcome presets cite a policy for each status they set by convention; a `statuses:`
  override is written `{ status, ref: policy:<id> }`. `credential_vocabulary.policy_refs`
  moved to `policies.yaml`.
- The gate measures statement coverage of the engine and the credentials package together.
- `uses:` cycles are detected by following the uses graph and reported with their full path
  (`uses cycle: a#x -> easa.shared.y -> a#x`), replacing the fixed depth limit; chains of any
  depth resolve. A shared evaluation may use another shared evaluation, never a credential's
  evaluation; `uses:` may not name a requirement-only entry, and a credential may not
  require itself. The gate fails on each.
- `docs/dependencies.md` (written by `go generate`): the requires and uses graphs per
  authority as Mermaid diagrams, a table per credential and the most depended-on shared
  evaluations and credentials.
