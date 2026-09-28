# Reviews

First-pass reviews of every interpretation in the catalogue, each against the verbatim texts
under `sources/`. A review is not a sign-off: no interpretation is approved yet, and
`approved_by` / `approved_on` stay empty until a qualified reviewer approves an
interpretation in a pull request (CONTRIBUTING.md, "Interpretations and sign-off").

Verdicts: **supported** (the cited paragraph says what the reading says), **needs-decision**
(the text or the record leaves real room; the owner picks one of the options) and **fixed**
(the reading or the result contradicted the text and was corrected, by the review itself or,
where the review deferred the fix, at integration).

## Index

| Review | Scope | Interpretations | Supported | Needs decision | Fixed |
| --- | --- | --- | --- | --- | --- |
| [easa-medicals-examiners-language.md](easa-medicals-examiners-language.md) | EASA medicals, examiners, language proficiency, powered-lift type rating | 27 | 20 | 4 | 3 |
| [easa-shared-privileges.md](easa-shared-privileges.md) | EASA shared evaluations and privileges | 57 | 42 | 11 | 4 (2 at integration) |
| [easa-class-ir-ratings.md](easa-class-ir-ratings.md) | EASA class ratings, IR(A), IR(H), IR(As), BIR | 51 | 42 | 6 | 3 (2 at integration) |
| [easa-types-licences.md](easa-types-licences.md) | EASA type ratings and licences | 59 | 47 | 10 | 2 |
| [easa-spl-gpl-instructors.md](easa-spl-gpl-instructors.md) | EASA SPL, GPL and instructor certificates | 58 | 45 | 7 | 6 |
| [faa-instructors-privileges-ratings-medicals.md](faa-instructors-privileges-ratings-medicals.md) | FAA instructors, privileges, instrument rating, medicals | 36 | 32 | 3 | 1 |
| [faa-shared-licences.md](faa-shared-licences.md) | FAA shared evaluations and pilot certificates | 55 | 44 | 8 | 3 |
| [de-ul-licence-privileges-instructor.md](de-ul-licence-privileges-instructor.md) | German ultralight licence, privileges, instructor rating, shared passenger evaluation | 25 | 15 | 10 | 0 |
| [de-ul-ratings.md](de-ul-ratings.md) | German ultralight kind ratings | 40 | 22 | 17 | 1 |
| **Total** | | **408** | **309** | **76** | **23** |

The catalogue has 407 interpretations; `easa.rating.powered-lift-type` /
`part-fcl-type-ratings` is counted twice (fixed in the medicals review, supported in the
types review). `easa.shared.sfcl-passenger-prerequisite` / `fi-s-on-any-licence`, supported,
was also changed at integration to follow the FI(S) fix of the SPL, GPL and instructors
review.

## Decisions for the owner

Every needs-decision item, one row per interpretation and credential file. The file's
current choice is the first option; choosing another changes the credential, its worked
examples and the reading together, and some options need a vocabulary or record change (as
noted). Details and reasoning are in the review file named.

| Credential | Interpretation | Question | Options (file's current choice first) | Review |
| --- | --- | --- | --- | --- |
| `easa.medical.class-1` | `counted-from` | On which date is the holder's age taken for the validity band? | (1) the age on the date the period starts; (2) the age on the examination date; (3) class 1 only: shorten a running certificate when the holder turns 60 | [easa-medicals-examiners-language](easa-medicals-examiners-language.md) |
| `easa.medical.class-2` | `counted-from` | As for class 1 (40th and 50th birthdays). | (1) the age on the first day of the period; (2) the age on the examination date | [easa-medicals-examiners-language](easa-medicals-examiners-language.md) |
| `easa.medical.lapl` | `counted-from` | As for class 2 (40th birthday). | (1) the age on the first day of the period; (2) the age on the examination date | [easa-medicals-examiners-language](easa-medicals-examiners-language.md) |
| `easa.examiner.fcl-examiner` | `tests-as-examiner-time` | In which period must the six tests of FCL.1025(b)(1) fall? | (1) within the validity period (36 months before expiry); (2) at any time before expiry | [easa-medicals-examiners-language](easa-medicals-examiners-language.md) |
| `easa.shared.instructor-certificate-validity` | `period-start` | When do an instructor certificate's 3 years start? | (1) 3 years from validFrom, else the issue date, capping a later recorded expiry; (2) trust the recorded expiry, 3 years only when none is recorded | [easa-shared-privileges](easa-shared-privileges.md) |
| `easa.shared.passengers-night` | `ir-held` | What does "holds an IR" in FCL.060(b)(2)(ii) require? | (1) an IR on the same licence, not expired (none recorded counts); (2) any IR, any category and validity; (3) a valid IR of the same aircraft category | [easa-shared-privileges](easa-shared-privileges.md) |
| `easa.shared.sfcl-passenger-prerequisite` | `competence-flight-not-recorded` | Must the passenger-competence training flight of SFCL.115(a)(2)(ii)(A) be recorded? | (1) informational untracked row, result on hours or launches alone; (2) unknown until a training flight event is recorded (record field or event kind needed); (3) as (1), saying in the result that the flight is assumed | [easa-shared-privileges](easa-shared-privileges.md), [easa-spl-gpl-instructors](easa-spl-gpl-instructors.md) |
| `easa.shared.variants` | `same-engine-sep-variants-evaluated` | How are SEP variants with the same engine type (FCL.710(da)) handled? | (1) evaluate every recorded SEP variant; (2) only variants named after an engine type; (3) report SEP variants unknown | [easa-shared-privileges](easa-shared-privileges.md) |
| `easa.privilege.cloud-flying` | `dual-counts` | Do dual cloud flights count toward SFCL.215(e)? | (1) PIC and dual alike; (2) PIC only for (e), dual only as the (f)(2) make-up | [easa-shared-privileges](easa-shared-privileges.md) |
| `easa.privilege.cloud-flying` | `check-restores` | Which check restores the privileges under SFCL.215(f)(1), and for how long? | (1) a check recorded for CLOUD_FLYING, for 24 months; (2) the same check, until (e) can next be met; (3) any sailplane proficiency check with an FE(S) | [easa-shared-privileges](easa-shared-privileges.md) |
| `easa.privilege.fcl-banner-towing` | `tows-counted` | Which tows count toward the 5 of FCL.805(e)? | (1) banner tows only; (2) any tow | [easa-shared-privileges](easa-shared-privileges.md) |
| `easa.privilege.fcl-sailplane-towing` | `tows-counted` | As for banner towing. | (1) sailplane tows only; (2) any tow | [easa-shared-privileges](easa-shared-privileges.md) |
| `easa.privilege.sfcl-banner-towing` | `tows-counted` | Which tows count toward the 5 of SFCL.205(f)? | (1) any tow in any aircraft except ultralights; (2) banner tows only; (3) (1) or (2) restricted to TMGs | [easa-shared-privileges](easa-shared-privileges.md) |
| `easa.privilege.sfcl-sailplane-towing` | `tows-counted` | As for SFCL banner towing. | (1) any tow except in ultralights; (2) sailplane tows only; (3) either restricted to TMGs | [easa-shared-privileges](easa-shared-privileges.md) |
| `easa.privilege.spl-tmg` | `twelve-hours-on-sailplanes` | On which aircraft are the 12 hours of SFCL.160(b)(1) flown? | (1) sailplanes including TMGs; (2) any aircraft | [easa-shared-privileges](easa-shared-privileges.md) |
| `easa.rating.sep-land` | `examiner-not-checked` | Which checks count as an SEP/TMG class rating check? | (a) (1) aeroplane only, (2) also an FFS flight of the class; (b) (1) any check in the class, (2) exclude checks recorded for the IR or BIR (new filter needed) | [easa-class-ir-ratings](easa-class-ir-ratings.md) |
| `easa.rating.sep-land` | `exemption-period` | In which period does a check exempt from the refresher of FCL.740.A(b)(1)(ii)(C)? | (1) the 12 months before expiry; (2) the validity period (24 months); (3) any time since the last revalidation or issue | [easa-class-ir-ratings](easa-class-ir-ratings.md) |
| `easa.rating.sep-land` | `land-sea-pooling-not-applied` | Is the SEP land/sea pooling of FCL.740.A(b)(4) applied? | (1) not applied; (2) a third route over both classes with the per-class minima, while both are held on the same licence | [easa-class-ir-ratings](easa-class-ir-ratings.md) |
| `easa.rating.sep-land` | `ul-credit-time-only` | Is ultralight motorglider time credited to the TMG class? | (1) every ultralight motorglider counts as TMG; (2) none is credited | [easa-class-ir-ratings](easa-class-ir-ratings.md) |
| `easa.rating.sep-sea` | `no-land-sea-pooling` | As `land-sea-pooling-not-applied` (decide together). | (1) not applied; (2) applied as for SEP land | [easa-class-ir-ratings](easa-class-ir-ratings.md) |
| `easa.rating.bir` | `alternate-check-not-tracked` | How is the every-second-revalidation check in an aeroplane of FCL.835(g)(3) handled? | (1) not tracked; (2) experience route only after a BIR check in an aeroplane in the 24 months before expiry; (3) as (2), unknown instead of not met without such a check | [easa-class-ir-ratings](easa-class-ir-ratings.md) |
| `easa.licence.cpl-a` | `class-1-only` | Is a CPL or ATPL with only a class 2 medical unusable altogether? | (1) licence expired with a class 2 medical; (2) usable within its PPL privileges, commercial privileges as a limitation (the composite cannot express this yet); (3) keep (1), ratings also accept the PPL privileges of a CPL or ATPL | [easa-types-licences](easa-types-licences.md) |
| `easa.licence.cpl-h` | `class-1-only` | As for the CPL(A). | As for the CPL(A) | [easa-types-licences](easa-types-licences.md) |
| `easa.licence.cpl-as` | `class-1-only` | As for the CPL(A). | As for the CPL(A) | [easa-types-licences](easa-types-licences.md) |
| `easa.licence.atpl-a` | `class-1-only` | As for the CPL(A). | As for the CPL(A) | [easa-types-licences](easa-types-licences.md) |
| `easa.licence.atpl-h` | `class-1-only` | As for the CPL(A). | As for the CPL(A) | [easa-types-licences](easa-types-licences.md) |
| `easa.rating.helicopter-type` | `validity-period-is-twelve-months` | Does the duration of an FSTD proficiency check count toward the 2 hours of FCL.740.H(a)(1)(i)? | (1) only a check flown in the helicopter; (2) an FSTD check too | [easa-types-licences](easa-types-licences.md) |
| `easa.licence.lapl-a` | `lapl-proficiency-check` | Which check counts as the LAPL(A) proficiency check of FCL.140.A(a)(2)? | (1) any aeroplane or TMG proficiency check; (2) only a check recorded for the LAPL(A) | [easa-types-licences](easa-types-licences.md) |
| `easa.licence.lapl-a` | `ul-credit-time-only` | How far are ultralight hours credited under FCL.035(a)(4)? | (1) toward the 12 hours only; (2) also toward the per-class hour of FCL.140.A(b); (3) also ultralight dual toward the refresher | [easa-types-licences](easa-types-licences.md) |
| `easa.licence.lapl-a` | `variants-sep-only` | How are SEP variants with the same engine type handled? | (1) evaluate every recorded SEP variant; (2) unknown unless the record marks a different engine type; (3) ask integrators to record only different-engine variants | [easa-types-licences](easa-types-licences.md) |
| `easa.licence.lapl-h` | `check-on-type` | Does an FSTD check count as the proficiency check of FCL.140.H(a)(2)? | (1) an FSTD check on the type counts; (2) only a check in the helicopter | [easa-types-licences](easa-types-licences.md) |
| `easa.licence.spl` | `launch-as-pic` | What does the passenger message say about sailplanes at night? | (1) the day-only message; (2) current with no statement about night | [easa-spl-gpl-instructors](easa-spl-gpl-instructors.md) |
| `easa.licence.spl` | `training-course-for-students` | Which SFCL.130(a)(2) path is evaluated? | (1) (iv) for every student, said in the reading; (2) (v) when the record says TMG privileges are sought | [easa-spl-gpl-instructors](easa-spl-gpl-instructors.md) |
| `easa.licence.spl` | `other-category-credit-shown-only` | May the SFCL.130(b) credit change the result? | (1) shown only; (2) up to 7 hours toward the 15 hours, with the limits of (b)(1) and (b)(2) | [easa-spl-gpl-instructors](easa-spl-gpl-instructors.md) |
| `easa.licence.gpl` | `class-or-type` | Do flights in another gyroplane class or type count? | (1) every gyroplane flight counts; (2) unknown unless each flight's class or type is recorded | [easa-spl-gpl-instructors](easa-spl-gpl-instructors.md) |
| `easa.instructor.fi` | `valid-until-expiry` | In which window are the FI revalidation items counted? | (1) the validity period (36 months before expiry); (2) any time before expiry | [easa-spl-gpl-instructors](easa-spl-gpl-instructors.md) |
| `easa.instructor.fi` | `instruction-any-category` | Must FI instruction be in the certificate's aircraft category? | (1) any category; (2) the category of the licence the FI is recorded on, with 20 hours for airships and gyroplanes | [easa-spl-gpl-instructors](easa-spl-gpl-instructors.md) |
| `easa.instructor.cri` | `instruction-in-aeroplanes` | Does any aeroplane instruction count as the 10 hours "as a CRI"? | (1) any aeroplane instruction; (2) class-rating instruction only (untracked, as the record cannot show it) | [easa-spl-gpl-instructors](easa-spl-gpl-instructors.md) |
| `faa.instructor.flight-instructor` | `regime-by-recorded-expiry` | Which 61.19(d) regime applies? | (1) the recorded expiry decides; (2) the issue date decides, an old certificate without expiry is unknown (needs a stage condition on the issue date) | [faa-instructors-privileges-ratings-medicals](faa-instructors-privileges-ratings-medicals.md) |
| `faa.rating.instrument-airplane` | `any-powered-category` | In which category must the 61.57(c) tasks be flown? | (1) any powered category; (2) airplane only; (3) record the category on the IR rating and count by it (vocabulary and record change) | [faa-instructors-privileges-ratings-medicals](faa-instructors-privileges-ratings-medicals.md) |
| `faa.rating.instrument-airplane` | `grace-from-last-met` | Which status does the six-month grace period report? | (1) `expiring`; (2) `lapsed` with `policy:recency-lapsed`, keeping the message | [faa-instructors-privileges-ratings-medicals](faa-instructors-privileges-ratings-medicals.md) |
| `faa.shared.flight-review` | `who-and-endorsement-not-checked` | Which proficiency checks and practical tests replace the flight review? | (1) every recorded one, including another authority's; (2) only checks for an FAA certificate, rating or privilege by an FAA examiner (needs an event authority) | [faa-shared-licences](faa-shared-licences.md) |
| `faa.shared.flight-review` | `simulator-counts` | Does a flight review in a device count? | (1) yes, the part 142 conditions taken as met; (2) exclude device sessions | [faa-shared-licences](faa-shared-licences.md) |
| `faa.shared.passengers-day` | `class-or-type` | Are SEP and SET one FAA class for passenger recency? | (1) count by the record's class (SEP and SET separate); (2) pool SEP with SET, land and sea (vocabulary request) | [faa-shared-licences](faa-shared-licences.md) |
| `faa.shared.passengers-night` | `class-or-type` | As by day. | As by day | [faa-shared-licences](faa-shared-licences.md) |
| `faa.shared.passengers-night` | `logged-night-takeoffs` | Which takeoffs count for the 61.57(b) period? | (1) logged night takeoffs (1.1 night); (2) full-stop night landings stand for both; (3) record takeoffs in the 61.57(b) period (vocabulary and record change) | [faa-shared-licences](faa-shared-licences.md) |
| `faa.licence.powered-lift-airship` | `medical-grade-not-recorded` | Which medical does a powered-lift or airship certificate need? | (1) at least second class for every grade; (2) split by certificate grade, as the rotorcraft files | [faa-shared-licences](faa-shared-licences.md) |
| `faa.licence.sport` | `sport-classes` | Do powered parachutes and weight-shift-control aircraft count for each other? | (1) both recorded as OTHER, pooled; (2) add classes for them (vocabulary request) | [faa-shared-licences](faa-shared-licences.md) |
| `faa.licence.sport` | `medical-or-drivers-license` | Does the sport pilot composite check the medical qualification? | (1) no medical requirement; (2) require a medical certificate or BasicMed; (3) add a driver's-license credential and require it or a medical | [faa-shared-licences](faa-shared-licences.md) |
| `de.licence.ul` | `ultralight-licences` | Is the § 45(1) medical for ultralights above 120 kg evaluated? | (1) left unevaluated, the gap named in the reading; (2) `requires_any` class 2 / LAPL medical for the kinds practically always above 120 kg | [de-ul-licence-privileges-instructor](de-ul-licence-privileges-instructor.md) |
| `de.instructor.ul-instructor` | `three-years-from-valid-from` | Do the three years end on the same calendar day or the day before (§§ 187, 188 BGB)? | (1) the catalogue-wide convention (same day); (2) one day earlier for German periods (engine or compiler option) | [de-ul-licence-privileges-instructor](de-ul-licence-privileges-instructor.md) |
| `de.instructor.ul-instructor` | `last-three-years-before-expiry` | Where do "the last three years" of § 96(4) end? | (1) the 36 months before expiry; (2) the three years before the date evaluated or the application | [de-ul-licence-privileges-instructor](de-ul-licence-privileges-instructor.md) |
| `de.instructor.ul-instructor` | `instruction-as-instructor-or-examiner` | How are the § 96(4) no. 1 take-offs, landings and kinds counted? | (1) 60 take-offs and 60 landings, in any ultralight kind; (2) 60 in total, or only the kinds the rating covers | [de-ul-licence-privileges-instructor](de-ul-licence-privileges-instructor.md) |
| `de.instructor.ul-instructor` | `befaehigungspruefung-event` | What is the "Befähigungsprüfung" of § 96(4) no. 3? | (1) an instructor assessment of competence only; (2) any passed proficiency check | [de-ul-licence-privileges-instructor](de-ul-licence-privileges-instructor.md) |
| `de.privilege.ul-towing` | `tow-kind-not-recorded` | Which tows count "in der jeweils eingetragen Art"? | (1) every tow in an ultralight, whatever was towed; (2) record the entered kind and filter by it, or also count tows in other aircraft | [de-ul-licence-privileges-instructor](de-ul-licence-privileges-instructor.md) |
| `de.privilege.ul-towing` | `tows-counted` | Is a tow flight a towed object or a flight? | (1) count towed objects; (2) count flights flagged as tow flights, or landings on them | [de-ul-licence-privileges-instructor](de-ul-licence-privileges-instructor.md) |
| `de.shared.ul-passengers` | `authorisation-recorded` | Does a passenger authorisation count for every kind? | (1) any kind; (2) compare the privilege's detail with the kind flown, unknown otherwise | [de-ul-licence-privileges-instructor](de-ul-licence-privileges-instructor.md) |
| `de.shared.ul-passengers` | `authorisation-progress-informational` | For which kinds are the § 84a(2) progress rows shown? | (1) every kind, as a guide; (2) only the kinds § 84a(2) covers, a `not_recorded` row for the others | [de-ul-licence-privileges-instructor](de-ul-licence-privileges-instructor.md) |
| `de.shared.ul-passengers` | `xc-200-km-per-flight` | Are the 200 km of § 84a(2) per flight or for both flights? | (1) per flight; (2) combined over the two flights (needs a new word) | [de-ul-licence-privileges-instructor](de-ul-licence-privileges-instructor.md) |
| `de.rating.ul-three-axis` | `scope-of-closing-words` | Which parts of the 12 hours of § 45(2) must be on three-axis ultralights? | (1) the training flight only; (2) also the PIC hours, take-offs and landings | [de-ul-ratings](de-ul-ratings.md) |
| `de.rating.ul-three-axis` | `training-flight-as-dual` | Does a flight logged as PIC with an instructor aboard count as the training flight? | (1) only a flight with dual time; (2) also a flight flagged as flown with an instructor (record field or flag needed) | [de-ul-ratings](de-ul-ratings.md) |
| `de.rating.ul-three-axis` | `proficiency-check-period` | How long does a § 45(3) check count? | (1) 24 months; (2) only until the experience can be counted again | [de-ul-ratings](de-ul-ratings.md) |
| `de.rating.ul-three-axis` | `medical-above-120-kg-not-evaluated` | Is the § 45(1) medical evaluated for this kind? | (1) left unevaluated; (2) require a class 2 or LAPL medical | [de-ul-ratings](de-ul-ratings.md) |
| `de.rating.ul-weight-shift` | `association-rules-by-authority` | Which association rule applies to which licence? | (1) by issuing body, LBA licences with the DULV rule; (2) one rule for all | [de-ul-ratings](de-ul-ratings.md) |
| `de.rating.ul-weight-shift` | `dulv-pic-only` | Are the DULV numbers (12 hours PIC in 24 months, no check alternative) right? | (1) as encoded; (2) as the association text says once verified (the review names no alternative) | [de-ul-ratings](de-ul-ratings.md) |
| `de.rating.ul-weight-shift` | `daec-safety-training-not-recorded` | Are the DAeC numbers (12 hours, 12 landings, untracked safety training) right? | (1) as encoded; (2) as the association text says once verified (take-offs may also be named) | [de-ul-ratings](de-ul-ratings.md) |
| `de.rating.ul-weight-shift` | `above-120-kg-not-named` | Is recency evaluated for weight-shift ultralights above 120 kg, which § 45 does not name? | (1) the association rule regardless of mass; (2) no statutory recency above 120 kg | [de-ul-ratings](de-ul-ratings.md) |
| `de.rating.ul-gyroplane` | `association-rule-unverified` | Do the DULV numbers apply to every issuer? | (1) the DULV numbers for every issuer; (2) per issuer once the DAeC rule is known | [de-ul-ratings](de-ul-ratings.md) |
| `de.rating.ul-gyroplane` | `gyroplane-class-credited` | Are Part-FCL gyroplane time and checks credited, and take-offs not counted? | (1) as encoded; (2) as the association text says once verified (the review names no alternative) | [de-ul-ratings](de-ul-ratings.md) |
| `de.rating.ul-gyroplane` | `above-120-kg-not-named` | As for weight-shift. | As for weight-shift | [de-ul-ratings](de-ul-ratings.md) |
| `de.rating.ul-helicopter` | `training-flight-as-dual` | As for three-axis. | As for three-axis | [de-ul-ratings](de-ul-ratings.md) |
| `de.rating.ul-helicopter` | `proficiency-check-period` | How long does a § 45(3) check count? | (1) 12 months; (2) only until the experience can be counted again | [de-ul-ratings](de-ul-ratings.md) |
| `de.rating.ul-powered-paraglider` | `association-rule-unverified` | Are 30 landings in 24 months right? | (1) as encoded; (2) as the association texts say once verified | [de-ul-ratings](de-ul-ratings.md) |
| `de.rating.ul-powered-paraglider` | `above-120-kg-not-named` | As for weight-shift. | As for weight-shift | [de-ul-ratings](de-ul-ratings.md) |
| `de.rating.ul-sailplane` | `association-rule-unverified` | Are 5 landings in 12 months, any launch method, right for every issuer? | (1) the DAeC rule for every issuer; (2) per issuer or as the association text says once verified | [de-ul-ratings](de-ul-ratings.md) |
| `de.rating.ul-sailplane` | `above-120-kg-not-named` | As for weight-shift. | As for weight-shift | [de-ul-ratings](de-ul-ratings.md) |

76 items: EASA 38, FAA 11, Germany 27.

### Open points inside fixed items

These interpretations were corrected, but the review leaves part of the reading to the
signing reviewer:

- `easa.licence.lapl-a` / `land-sea-split`: for the SEP land and SEP sea privileges, the
  split applies whenever both are held (the file's reading), or experience in one class keeps
  that class current and the split is needed only to cover both.
- `easa.instructor.fi-s` / `sfcl-certificate`: whether a holder of a certificate converted
  from Part-FCL needs SFCL.360 before the endorsed expiry date (the fix applies SFCL.360
  throughout).
- `easa.instructor.fi-s` / `resumption`: the window of SFCL.360(d); the assessment stands in
  for the instruction for 36 months, or restores from its date.

Gaps that no interpretation names are listed at the end of `easa-class-ir-ratings.md`
("Other findings": IR-only and BIR-only checks counted as class rating checks) and
`faa-shared-licences.md` ("Observations").
