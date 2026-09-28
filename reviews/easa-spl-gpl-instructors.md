# Review: EASA SPL, GPL and instructor certificates

Reviewed on 2026-09-28 against the verbatim texts in `sources/easa/` (EUR-Lex consolidated
versions 30.04.2026 of Regulation (EU) No 1178/2011 and 15.11.2021 of Regulation (EU)
2018/1976). This review is not a sign-off: every interpretation keeps
`approved_by: null` and `approved_on: null` until a qualified reviewer approves it (see
CONTRIBUTING.md, "Interpretations and sign-off").

## Summary

| Verdict | SPL | GPL | Instructors | Total |
| --- | --- | --- | --- | --- |
| supported | 13 | 9 | 23 | 45 |
| needs-decision | 3 | 1 | 3 | 7 |
| contradicted (fixed) | 0 | 0 | 6 | 6 |
| **Total** | **16** | **10** | **32** | **58** |

Contradicted and fixed:

- `easa.instructor.fi` / `refresher-and-assessment-events`: any instructor refresher or
  assessment of competence counted, including an FI(S) refresher or an examiner assessment
  under FCL.1020. The FI text asks for refresher training "as an FI" and an assessment under
  FCL.935. Both counts now require the event to be recorded for the FI or the IRI (the IRI
  uses this evaluation).
- `easa.instructor.cri` / `refresher-and-assessment-events`: the same for the CRI; both
  counts now require `rating: CRI`.
- `easa.instructor.fi-s` / `refresher-event`, `demonstration-event`, `resumption`: the
  refresher, the demonstration flight and the assessment of competence counted whatever
  certificate they were recorded for (an FE(S) demonstration under SFCL.460(b)(2) kept an
  FI(S) current). They now require `rating: FI_S`.
- `easa.instructor.fi-s` / `sfcl-certificate`: a recorded expiry date made the FI(S)
  expired. Part-SFCL gives the FI(S) no validity period, and Article 3b(2)(c) of Regulation
  (EU) 2018/1976 lets the holder of a converted Part-FCL certificate instruct after its
  endorsed expiry date while complying with SFCL.360. The recency alone decides now.

Needs a decision by the signing reviewer: the passenger message for sailplanes at night
(`launch-as-pic`), which SPL course path is evaluated (`training-course-for-students`),
whether the SFCL.130(b) credit may change the result (`other-category-credit-shown-only`),
the gyroplane class or type (`class-or-type`), the window for FI revalidation items
(`valid-until-expiry`), the aircraft category of FI instruction (`instruction-any-category`)
and the capacity of CRI instruction (`instruction-in-aeroplanes`).

Sources added: `sources/easa/fcl-935.md`, `fcl-915-sfi.md`, `fcl-915-mcci.md`,
`fcl-915-sti.md`, `fcl-930-mcci.md`, `sfcl-345.md` and `reg-2018-1976-article-3b.md`
(Article 3b of Regulation (EU) 2018/1976). FCL.900 and FCL.915, which this review also
relies on, were added at the same time by a parallel review, with the same text.

## How to use this review

Each table gives one row per interpretation of a credential file: its verdict, the
paragraphs read, the reasoning, and what was changed. A signing reviewer can:

1. read the cited paragraph in `sources/easa/` next to the reasoning;
2. for `supported`, confirm and sign off in the pull request as CONTRIBUTING.md describes;
3. for `needs-decision`, choose one of the options given (the file's current choice is
   named) and either sign off or ask for the reading to change;
4. for `contradicted`, check the fix in the credential file and its new worked example.

The reasoning is in our own words; where the regulation is quoted, it is one short
sentence at most.

## easa.licence.spl (`credentials/easa/licences/spl.yaml`)

| Interpretation | Verdict | Paragraph(s) | Reasoning | Action |
| --- | --- | --- | --- | --- |
| `sailplane-privileges-sfcl-licences` | supported | SFCL.160(a); Article 3b(1) of Reg. 2018/1976 | SFCL.160 applies to SPL holders. Article 3b(1) deems Part-FCL sailplane licences (the LAPL(S) and the old SPL) issued under Part-SFCL, so selecting GLIDER ratings on an SPL or LAPL(S) of an EU authority is right; a GLIDER rating on another licence is not a Part-SFCL licence. | None. |
| `sailplanes-include-tmg-time` | supported | SFCL.160(a)(1), (a)(1)(i)-(ii); FCL.010 | FCL.010 defines a TMG as a class of powered sailplane, and (a)(1) sets the 5 hours on sailplanes but the 15 launches and 2 training flights on sailplanes excluding TMGs. TMG time therefore counts toward the 5 hours only. | None. |
| `training-flight-as-dual` | supported | SFCL.160(a)(1)(ii), (a)(2) | A training flight with an instructor is logged as dual time received; who the instructor or examiner was is not in the record. | None. |
| `proficiency-check-in-class` | supported | SFCL.160(a)(2) | The check must be on a sailplane excluding TMGs; a check in the GLIDER class matches. | None. |
| `launch-count` | supported | SFCL.160(a)(1)(i); SFCL.155(c) | Counting each launch of a multi-launch row is the natural count. Taking a TMG's take-offs as launches is what SFCL.155(c) allows for self-launch; for recency the launches are counted on GLIDER flights only. | None. |
| `no-expiry` | supported | SFCL.160(a); SFCL.115(c) | Part-SFCL sets no licence validity; privileges depend on recency and the medical, which the file evaluates. | None. |
| `logbook-signatures-not-checked` | supported | SFCL.160(d) | Signatures are not in the record; a scope statement. | None. |
| `launch-methods-in-use` | supported | SFCL.155(a), (c) | (c) keeps each method separately with 5 launches (2 for bungee) in the last two years, without naming a role, so launches in any role count. Evaluating a method once it was logged or its training recorded is a record mapping the text does not contradict. | None. |
| `self-launch-with-tmg` | supported | SFCL.155(c) | (c) lets self-launch recency be met by self-launches, TMG take-offs or both. Not treating TMG flying alone as a self-launch method in use is prudent: (a)(2) requires self-launch training of its own. | None. |
| `other-launch-methods` | supported | SFCL.155(a)(4), (c) | Other methods are set by the competent authority and are not recorded; (c) would apply to them too, which the file says it does not evaluate. | None. |
| `launch-as-pic` | needs-decision | SFCL.160(e)(1), (e)(2) | Counting launches of GLIDER flights with PIC time, without supervised solo, matches launches as PIC in sailplanes excluding TMGs. The last sentence goes beyond the text: (e) sets a night condition only for TMGs and says nothing about sailplane passengers at night, so the day-only message is not taken from (e). Options: (1) keep the day-only message (the file's choice); (2) report current without a statement about night. | None; for decision. |
| `training-course-for-students` | needs-decision | SFCL.130(a)(2), (a)(2)(iv)-(v); MED.A.030(a) | The hours, the TMG share and dual plus supervised solo as instruction match (a)(2). The reading does not say that the course is always evaluated for sailplane privileges: (iv) applies only "if privileges for sailplanes, excluding TMGs, are sought", and the TMG path of (v) is not evaluated. Taking a medical certificate as the sign of a student follows MED.A.030(a) but is a record assumption. Options: (1) evaluate (iv) for every student (the file's choice) and say so in the reading; (2) evaluate (v) instead when the record says TMG privileges are sought. | None; for decision. |
| `cross-country-flight` | supported | SFCL.130(a)(2)(iv)(B)(a)-(b) | The solo flight of 50 km in a sailplane excluding TMGs, or the dual flight of 100 km that may be flown in a TMG, is encoded as written; unknown without a recorded distance. | None. |
| `other-category-credit-shown-only` | needs-decision | SFCL.130(b), (b)(1)-(b)(2) | (b) grants the credit ("shall be credited"); showing it without letting it change the result can report a course in progress that the text treats as complete. The reading names aeroplane, gyroplane and helicopter licences; the text covers every other category except balloons, which the encoded licence list (airships included) already follows. Options: (1) show only (the file's choice); (2) apply up to 7 hours toward the 15 hours, with the limits of (b)(1) and (b)(2). | None; for decision. |
| `theory-and-ato-not-checked` | supported | SFCL.130(a), (a)(1) | Neither the theory nor the training organisation is in the record. | None. |
| `commercial-medical-not-applied` | supported | MED.A.030(c)(1), (c)(4) | The LAPL medical is the minimum for the SPL; the class 2 medical is needed only for commercial operations, which the record does not show. | None. |

## easa.licence.gpl (`credentials/easa/licences/gpl.yaml`)

| Interpretation | Verdict | Paragraph(s) | Reasoning | Action |
| --- | --- | --- | --- | --- |
| `gyroplane-privileges-on-gpl` | supported | FCL.240.G(a); FCL.205.G(a) | Both points apply to GPL holders; a GYROPLANE rating on another licence is not a GPL. | None. |
| `class-or-type` | needs-decision | FCL.240.G(a), (b) | (a) counts experience "in the relevant class or type", and (b) shows that classes exist (the SPG class). The record has one GYROPLANE class, so flights in another gyroplane class also count, which can overstate recency for a holder with several classes. Options: (1) every gyroplane flight counts (the file's choice); (2) report unknown unless the class or type of each flight is recorded. | None; for decision. |
| `pic-dual-or-supervised` | supported | FCL.240.G(a)(1), (a)(1)(i) | The 12 hours are PIC, dual or supervised solo time as written; counting the smaller of take-offs and landings makes each pair count once. | None. |
| `ultralight-credit-mass` | supported | FCL.035(a)(5) | (a)(5) credits gyroplanes outside the scope of the Basic Regulation with a certificated take-off mass of at least 450 kg toward the 12 hours and 12 take-offs and landings, and not toward the refresher training. Leaving flights without a recorded mass out of the count is the cautious record choice. | None. |
| `refresher-as-dual` | supported | FCL.240.G(a)(1)(ii) | The refresher is at least 1 hour of total flight time with an instructor; dual time received in the class, possibly over several flights, fits "total flight time". | None. |
| `proficiency-check-in-class` | supported | FCL.240.G(a)(2); FCL.035(a)(5) | A GPL check with an examiner in the class; (a)(5) credits ultralight hours only, so a check in an ultralight gyroplane does not count. | None. |
| `no-expiry` | supported | FCL.240.G(a), (c) | The GPL has no validity in these points; the signatures of (c) are not in the record. | None. |
| `variants-recorded` | supported | FCL.240.G(b), (b)(1)-(b)(3) | (b) applies to a variant not flown in the preceding 2 years and accepts differences training, a proficiency check or refresher training in it. A flight in the variant (training and checks included) resets the 2 years; a flight without a variant cannot be attributed and makes the result unknown. | None. |
| `pax-prerequisite-since-issue` | supported | FCL.205.G(a)(2); FCL.035(a)(5) | 10 hours as PIC on gyroplanes after issue, as written. Not crediting ultralight gyroplanes follows from (a)(5), which credits them only toward FCL.240.G(a). The re-issue caveat and unknown without an issue date are record choices. | None. |
| `commercial-not-recorded` | supported | FCL.205.G(b) | Remuneration is not recorded; a scope statement. | None. |

## easa.instructor.fi (`credentials/easa/instructors/fi.yaml`)

| Interpretation | Verdict | Paragraph(s) | Reasoning | Action |
| --- | --- | --- | --- | --- |
| `part-fcl-certificate` | supported | FCL.900(a); FCL.940 | A Part-FCL FI is held with a Part-FCL licence; the balloon and sailplane instructors are Part-BFCL and Part-SFCL certificates. | None. |
| `three-years-from-valid-from` | supported | FCL.940 | 3 years, with a shorter recorded date winning, follows FCL.940. | None. |
| `valid-until-expiry` | needs-decision | FCL.940; FCL.940.FI(a)(1); FCL.935(d) | Reporting the certificate valid until its expiry date whatever the revalidation follows FCL.940. The window is open: (a)(1) only says "before the expiry date" for the 50 hours and the refresher, while the CRI text names the validity period. Options: (1) within the validity period (the file's choice: 36 months before expiry); (2) at any time before expiry. The same question was raised for examiner certificates (`tests-as-examiner-time`). Note: FCL.935(d) bars a holder who failed an assessment needed for revalidation; a failed assessment is not recorded. | None; for decision. |
| `instruction-any-category` | needs-decision | FCL.940.FI(a)(1)(i)(A)-(C); FCL.915(c)(2) | Counting instruction in any capacity and examiner time follows (A) and FCL.915(c)(2). The category is not checked: an FI(A) could be revalidated with helicopter instruction, and an FI(As) (20 hours) is held to 50. The licence the privilege is recorded on would show the category. Options: (1) any category (the file's choice); (2) the category of the licence the FI is recorded on, with 20 hours for airships and gyroplanes for a GPL. | None; for decision. |
| `refresher-and-assessment-events` | contradicted | FCL.940.FI(a)(1)(ii), (a)(1)(iii); FCL.935(a); FCL.1020 | The refresher must be training "as an FI" and the assessment one under FCL.935. The counts took any instructor_refresher and assessment_of_competence event, so an FI(S) refresher or an examiner assessment under FCL.1020 (which `easa.examiner.fcl-examiner` records the same way) revalidated an FI. | Both counts now require `for_rating: [FI, IRI]`; reading updated; existing examples record `rating: FI` (outcomes unchanged); new example `revalidation-events-for-other-certificates`. |
| `alternate-assessment-not-tracked` | supported | FCL.940.FI(a)(2), (b) | Earlier revalidations are not in the record; a scope statement. | None. |

## easa.instructor.fi-s (`credentials/easa/instructors/fi-s.yaml`)

| Interpretation | Verdict | Paragraph(s) | Reasoning | Action |
| --- | --- | --- | --- | --- |
| `sfcl-certificate` | contradicted | SFCL.360(a); Article 3b(2)(c) of Reg. 2018/1976 | Part-SFCL sets no validity for the FI(S). The only expiry date an FI(S) carries is the one endorsed on conversion from Part-FCL, and after it the holder may instruct while complying with SFCL.360. The file reported the certificate expired from the day after any recorded date, even with SFCL.360 met. Still open: whether a converted holder needs SFCL.360 before that date; the fix applies SFCL.360 throughout, which (a) supports for every FI(S) holder. | Explicit outcomes without the expired stage (current or lapsed by recency); reading rewritten; example `recency-expired` replaced by `recency-recorded-expiry-not-used` and `recency-recorded-expiry-not-used-lapsed`; the `date` parameter dropped from the current example. |
| `instruction-on-sailplanes` | supported | SFCL.360(a)(1)(ii), (b) | 30 hours or 60 launches or take-offs and landings of instruction in 3 years, with FE(S) hours credited, as written; examiner time adds hours, not launches. | None. |
| `refresher-event` | contradicted | SFCL.360(a)(1)(i) | The refresher must refresh knowledge for sailplane instructors; any instructor refresher counted, including a Part-FCL FI refresher. | Count requires `for_rating: [FI_S]`; reading updated; examples record `rating: FI_S`; new example `recency-events-for-other-certificates`. |
| `demonstration-event` | contradicted | SFCL.360(a)(2), (c); SFCL.345; SFCL.460(b)(2) | Accepting an SFCL.345 assessment in place of the demonstration is supported by (c). But any supervised instruction or assessment event counted, so the FE(S) demonstration of SFCL.460(b)(2), recorded as an assessment of competence, kept an FI(S) current. | Count requires `for_rating: [FI_S]`; reading updated; same examples. |
| `resumption` | contradicted | SFCL.360(d); SFCL.345 | Refresher plus an SFCL.345 assessment as the way back follows (d), but the assessment was any assessment event. Still open: (d) names no window; the file lets the assessment stand in for the instruction for 36 months, where it could instead restore from its date. | Count requires `for_rating: [FI_S]`; reading updated; same examples. Window left for the signing reviewer. |

## easa.instructor.cri (`credentials/easa/instructors/cri.yaml`)

| Interpretation | Verdict | Paragraph(s) | Reasoning | Action |
| --- | --- | --- | --- | --- |
| `part-fcl-certificate` | supported | FCL.900(a); FCL.940 | As for the FI. | None. |
| `three-years-from-valid-from` | supported | FCL.940 | As for the FI. | None. |
| `valid-until-expiry` | supported | FCL.940; FCL.940.CRI(a) | Here the text itself says "within the validity period", so the window is right, and the certificate stays valid until its expiry under FCL.940. | None. |
| `instruction-in-aeroplanes` | needs-decision | FCL.940.CRI(a)(1) | (a)(1) asks for 10 hours "as a CRI". The file counts any instruction in aeroplanes (TMGs included), including instruction as an FI, and does not check the split between single-engine and multi-engine aeroplanes. Options: (1) any aeroplane instruction (the file's choice); (2) instruction in class-rating training only, which the record cannot show, so the row would be untracked. | None; for decision. |
| `refresher-and-assessment-events` | contradicted | FCL.940.CRI(a)(2), (a)(3); FCL.935 | The refresher must be "as a CRI" and the assessment under FCL.935 for aeroplanes; any refresher or assessment event counted, including those recorded for an FI or an examiner certificate. | Both counts require `for_rating: [CRI]`; reading updated; examples record `rating: CRI`; new example `revalidation-events-for-other-certificates`. |
| `alternate-assessment-not-tracked` | supported | FCL.940.CRI(b), (c) | As for the FI. | None. |

## easa.instructor.iri (`credentials/easa/instructors/iri.yaml`)

| Interpretation | Verdict | Paragraph(s) | Reasoning | Action |
| --- | --- | --- | --- | --- |
| `revalidated-like-fi` | supported | FCL.940.IRI | FCL.940.IRI applies the FI revalidation as it stands. With the fix above, an event recorded for the IRI counts; so does one recorded for the FI, because the evaluation is shared. | Examples record `rating: IRI`. |

## easa.instructor.tri (`credentials/easa/instructors/tri.yaml`)

| Interpretation | Verdict | Paragraph(s) | Reasoning | Action |
| --- | --- | --- | --- | --- |
| `part-fcl-certificate` | supported | FCL.940 | A record mapping. (The reading starts "An TRI"; wording only.) | None. |
| `revalidation-not-evaluated` | supported | FCL.940.TRI(a)(1), (a)(2), (a)(5), (b) | Course parts, instruction per type and the category are not recorded. Not applying (a)(5) means a TRI(H) valid until the expiry of an FI(H) may be reported expired on its own 3 years; the reading says so. | None. |

## easa.instructor.mi (`credentials/easa/instructors/mi.yaml`)

| Interpretation | Verdict | Paragraph(s) | Reasoning | Action |
| --- | --- | --- | --- | --- |
| `same-licence-recorded-expiry` | supported | FCL.940.MI; FCL.940 | The MI is valid as long as the FI, TRI or CRI certificate is; FCL.940 excepts the MI from the 3 years. The underlying certificate's own evaluation comes in through the composite. | None. |

## easa.instructor.fti (`credentials/easa/instructors/fti.yaml`)

| Interpretation | Verdict | Paragraph(s) | Reasoning | Action |
| --- | --- | --- | --- | --- |
| `part-fcl-certificate` | supported | FCL.940 | A record mapping. | None. |
| `revalidation-not-evaluated` | supported | FCL.940.FTI(a)(1)(i)-(ii), (a)(2), (b) | Flight-test time and FTI refresher content are not recorded; a scope statement. | None. |

## easa.instructor.mcci (`credentials/easa/instructors/mcci.yaml`)

| Interpretation | Verdict | Paragraph(s) | Reasoning | Action |
| --- | --- | --- | --- | --- |
| `part-fcl-certificate` | supported | FCL.940 | A record mapping. | None. |
| `revalidation-not-evaluated` | supported | FCL.940.MCCI(a), (b); FCL.930.MCCI(a)(3) | Revalidation needs 3 hours of supervised practical instruction on the relevant FSTD type in the last 12 months; neither is identifiable in the record. The reading said FCL.930.MCCI was not stored; it is now. | Reading reworded to describe FCL.930.MCCI(a)(3), which it now cites; no result change. Source added. |
| `no-licence-requirement` | supported | FCL.900(a)(2); FCL.915.MCCI(a); FCL.910.MCCI | MCC instruction needs only the instructor certificate (FCL.900(a)(2)), the MCCI is restricted to FSTDs, and the prerequisite is to hold or have held a CPL, MPL or ATPL. | None. |

## easa.instructor.sfi (`credentials/easa/instructors/sfi.yaml`)

| Interpretation | Verdict | Paragraph(s) | Reasoning | Action |
| --- | --- | --- | --- | --- |
| `part-fcl-certificate` | supported | FCL.940 | A record mapping. | None. |
| `revalidation-not-evaluated` | supported | FCL.940.SFI(a)(1), (b), (e) | FSTD instructing hours and per-type FFS checks are not recorded; a scope statement. | None. |
| `no-licence-requirement` | supported | FCL.900(a)(2); FCL.915.SFI(a) | Synthetic flight instruction needs only the certificate; the prerequisite is to hold or have held a CPL, MPL or ATPL. | None. |

## easa.instructor.sti (`credentials/easa/instructors/sti.yaml`)

| Interpretation | Verdict | Paragraph(s) | Reasoning | Action |
| --- | --- | --- | --- | --- |
| `part-fcl-certificate` | supported | FCL.940 | A record mapping. | None. |
| `revalidation-not-evaluated` | supported | FCL.940.STI(a)(1)-(2), (b) | Course-bound FSTD instruction and the proficiency check on that FSTD are not identifiable; a scope statement. | None. |
| `no-licence-requirement` | supported | FCL.900(a)(2); FCL.915.STI(a)(1) | Synthetic flight instruction needs only the certificate. For the STI the licence prerequisite is narrower than the reading suggests (held within the 3 years before application), but it applies at issue only, so not requiring a current licence stands. | None. |

## Observations outside this scope (fix deferred)

These concern `credentials/easa/shared/`, which another review covers; this review changed
nothing there. The integration step applied the second item and left the first for decision.

- `easa.shared.sfcl-passenger-prerequisite` / `competence-flight-not-recorded`: SFCL.115(a)(2)(ii)(A)
  requires the passenger-competence training flight in addition to the hours or launches, but
  the row is informational, so the SPL passenger prerequisite reports met without it.
  Proposed: make the row a required, untracked row (result unknown until a
  training flight event is recorded), or declare the leniency as a policy. Not applied at
  integration: `reviews/easa-shared-privileges.md` marks the reading needs-decision, not
  contradicted, so it stays on the owner's decision list (`reviews/README.md`).
- `easa.shared.sfcl-passenger-prerequisite` / `fi-s-on-any-licence`: the FI(S) alternative
  checks that the FI_S privilege is not expired by its recorded date. After the fix above
  an FI(S) has no expiry that ends it. Proposed: `holds: { privileges: [FI_S] }` without
  `valid: true`, since SFCL.115(a)(2)(ii)(B) asks only that the pilot holds the certificate.
  Fixed at integration as proposed: the `instructor` stage checks `holds: { privileges:
  [FI_S] }`; the reading says that the recency of SFCL.360 (events recorded for FI_S) is
  evaluated by `easa.instructor.fi-s`; example `passenger-prerequisite-fi-s-recorded-expiry-passed`
  (current) added to `examples/easa.licence.spl.yaml`. No existing example changed its result.

## Sources added

| File | Article | Why |
| --- | --- | --- |
| `sources/easa/fcl-935.md` | FCL.935 | The instructor assessment of competence that FI and CRI revalidation name. |
| `sources/easa/fcl-915-sfi.md` | FCL.915.SFI | "Hold or have held" prerequisite behind `no-licence-requirement`. |
| `sources/easa/fcl-915-mcci.md` | FCL.915.MCCI | The same for the MCCI. |
| `sources/easa/fcl-915-sti.md` | FCL.915.STI | The same for the STI (within 3 years). |
| `sources/easa/fcl-930-mcci.md` | FCL.930.MCCI | The practical instruction that MCCI revalidation refers to. |
| `sources/easa/sfcl-345.md` | SFCL.345 | The FI(S) assessment of competence of SFCL.360(c) and (d). |
| `sources/easa/reg-2018-1976-article-3b.md` | Article 3b of Reg. 2018/1976 | Part-FCL sailplane licences deemed Part-SFCL; the FI(S) expiry date on conversion. |

The Part-FCL texts are taken unchanged from the EUR-Lex consolidated text
02011R1178-20260430 and the Part-SFCL texts from 02018R1976-20211115, both via the
Publications Office Cellar, with the header of the existing `sources/easa/` files. The
article of Regulation 2018/1976 is stored for reading only; no credential references it.
