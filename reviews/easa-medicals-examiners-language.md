# Review: EASA medicals, examiners, language proficiency and the powered-lift licence

Reviewed on 2026-09-28 against the verbatim texts in `sources/easa/` (EUR-Lex consolidated
versions 30.04.2026 of Regulation (EU) No 1178/2011 and 15.11.2021 of Regulation (EU)
2018/1976). This review is not a sign-off: every interpretation keeps
`approved_by: null` and `approved_on: null` until a qualified reviewer approves it (see
CONTRIBUTING.md, "Interpretations and sign-off").

## Summary

| Verdict | Count |
| --- | --- |
| supported | 20 |
| needs-decision | 4 |
| contradicted (fixed) | 3 |
| **Total** | **27** |

Contradicted and fixed:

- `easa.endorsement.language-proficiency` / `level-6-not-evaluated`: a level 6 endorsement
  was not selected, so a licence holder with only a level 6 endorsement was reported unknown
  (endorsement not held). A `level_6` evaluation now reports it valid.
- `easa.examiner.fe-s` / `demonstration-event`: any assessment of competence counted,
  including one as a pilot or instructor. Only an assessment recorded for FE_S counts now.
- `easa.rating.powered-lift-type` / `part-fcl-type-ratings`: the MPL was accepted as the
  licence of a powered-lift type rating; FCL.720.PL names only CPL and ATPL. The MPL was
  removed and the reading now cites FCL.720.PL.

Needs a decision by the signing reviewer: the date on which the holder's age is taken for
the medical validity bands (`counted-from`, three files), and the period in which the six
tests of an examiner revalidation must fall (`tests-as-examiner-time`).

Sources added: `sources/easa/fcl-720-pl.md`, `sources/easa/fcl-1000.md`,
`sources/easa/fcl-1020.md`.

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

## easa.medical.class-1 (`credentials/easa/medicals/class-1.yaml`)

| Interpretation | Verdict | Paragraph(s) | Reasoning | Action |
| --- | --- | --- | --- | --- |
| `single-pilot-cat-not-applied` | supported | MED.A.045(a)(2)(i) | The 6-month period at 40 or more depends on the holder being engaged in single-pilot commercial air transport carrying passengers. That is a fact about the holder's operations that the record does not hold, so the file cannot apply it; a shorter expiry printed on the certificate still wins. The consequence (a 12-month period reported for such a pilot unless the recorded date is shorter) is stated openly. | None. |
| `counted-from` | needs-decision | MED.A.045(a)(5), (a)(2)(ii) | The anchor is right: (a)(5) counts from the examination date for initial issue and renewal and from the previous expiry for revalidation, which the file maps to validFrom, else the issue date. The text does not say on which date the age is taken. Options: (1) the age on the date the period starts (the file's choice); (2) the age on the examination date; (3) for class 1 only, shorten a running certificate when the holder turns 60. The file's choice follows the structure of (a)(3) and (a)(4), where the legislator wrote an explicit age cap when it wanted one; (a)(2) has none, which argues against option (3). | Decided 2026-09-28: option 2 (age on the examination date), applied: `valid_for.age_on: issue`; the class 1 file also gained the class 2 and LAPL levels (bug fix, finding on class 1 validity). |
| `earlier-expiry-wins` | supported | MED.A.045(a) | (a) sets the maximum period; a certificate may carry a shorter validity, never a longer one, so taking the earlier of the derived and the recorded date respects the text. Without a date of birth the age band cannot be chosen and only the recorded date is used, which is the conservative reading. | None. |
| `renewal-not-modelled` | supported | MED.A.045(c) | (c) describes what examination a holder of an expired certificate needs; it does not restore the old certificate. Reporting an expired certificate as expired is exactly what the text implies. | None. |

## easa.medical.class-2 (`credentials/easa/medicals/class-2.yaml`)

| Interpretation | Verdict | Paragraph(s) | Reasoning | Action |
| --- | --- | --- | --- | --- |
| `age-bands` | supported | MED.A.045(a)(3)(i)-(iii) | The three bands (until 40, between 40 and 50, above 50) meet without a gap only if "between 40 and 50" ends before the 50th birthday; the cap sentence of (ii), about a certificate issued before 50, confirms that 49-year-olds belong to the 24-month band. The caps end the certificate on the 42nd or 51st birthday; the engine treats the certificate as valid through that day and expired the day after, matching "shall cease to be valid after the licence holder reaches the age of 42". | None. |
| `counted-from` | needs-decision | MED.A.045(a)(5), (a)(3) | As for class 1: the anchor follows (a)(5); the date on which the age is taken is not written. Options: (1) the age on the first day of the period (the file's choice); (2) the age on the examination date. They differ when a revalidation examination (up to 45 days early, (b)) falls before the 40th or 50th birthday and the previous expiry after it: option (1) gives the shorter band, option (2) the longer one with its age cap. | Decided 2026-09-28: option 2 (age on the examination date), applied: `valid_for.age_on: issue`. |
| `earlier-expiry-wins` | supported | MED.A.045(a) | As for class 1. | None. |
| `renewal-not-modelled` | supported | MED.A.045(c) | As for class 1. | None. |

## easa.medical.lapl (`credentials/easa/medicals/lapl.yaml`)

| Interpretation | Verdict | Paragraph(s) | Reasoning | Action |
| --- | --- | --- | --- | --- |
| `age-bands` | supported | MED.A.045(a)(4)(i)-(ii) | (i) runs until the holder reaches 40 and (ii) covers holders above 40, so 40 belongs to (ii) and there is no gap. The 42nd-birthday cap in (i) is encoded, and there is no later cap for LAPL certificates, which the file rightly does not add. | None. |
| `counted-from` | needs-decision | MED.A.045(a)(5), (a)(4) | Same question as for class 2, at the 40th birthday only. The file takes the age on the first day of the period. | Decided 2026-09-28: option 2 (age on the examination date), applied: `valid_for.age_on: issue`. |
| `earlier-expiry-wins` | supported | MED.A.045(a) | As for class 1. | None. |
| `renewal-not-modelled` | supported | MED.A.045(c) | (c)(3) sets the LAPL renewal assessment; reporting the expired certificate as expired is consistent with it. | None. |

## easa.examiner.fcl-examiner (`credentials/easa/examiners/fcl-examiner.yaml`)

| Interpretation | Verdict | Paragraph(s) | Reasoning | Action |
| --- | --- | --- | --- | --- |
| `one-privilege-kind` | supported | FCL.1025(a); FCL.1000 | FCL.1025 applies the same validity and revalidation to every Part-FCL examiner certificate, so evaluating all categories alike follows the text. Sailplane and balloon examiners are governed by Part-SFCL and Part-BFCL, so excluding the SPL, LAPL(S) and BPL is right. Worth knowing: FCL.1000(b)(1) allows a specific certificate for new aircraft or courses limited to 1 year; the file handles it only through a recorded expiry date. | None. |
| `three-years-from-valid-from` | supported | FCL.1025(a), (b) | (a) makes the certificate valid for 3 years regardless of the revalidation conditions, which only matter for the next period; reporting an unmet revalidation as current (with the not-met message, and expiring in the last 90 days) matches that. Capping a later recorded date at 3 years follows (a). | None. |
| `tests-as-examiner-time` | needs-decision | FCL.1025(b)(1) | The count of six tests or checks conducted is fine as flights or FSTD sessions logged with examiner time. The window is the open question: (b)(1) only says the tests must be conducted before the expiry date and names no start. Options: (1) within the certificate's validity period (the file's choice: 36 months before expiry); (2) at any time before expiry, since the first issue. Option (1) is the natural reading of a revalidation condition and the more cautious one. | Decided 2026-09-28: option 1 (within the validity period), applied: reading restated (P3), no change to results. |
| `refresher-event` | supported | FCL.1025(b)(2) | The course must fall in the 12 months before expiry; who provided or approved it cannot be recorded. The event kind carries no examiner category, so an FE(S) refresher would also count for a Part-FCL certificate; that follows from the record, not from the text. | None. |
| `assessed-test-event` | supported | FCL.1025(b)(3)(i)-(ii); FCL.1020 | (b)(3) accepts one of the tests either assessed by an inspector or senior examiner or complying with FCL.1020, and FCL.1020 is a demonstration of competence in the examiner role. An assessment of competence recorded for the rating EXAMINER is that event; excluding a pilot's own assessment is correct. | None. |
| `not-applied` | supported | FCL.1025(b), (c), (d) | The reading is a correct scope statement. Its reference pointed to (b)(3)(ii) (compliance with FCL.1020), but the combined revalidation of several categories is the closing sentence of (b). | Reference changed from FCL.1025(b)(3)(ii) to FCL.1025(b); reading unchanged. |

## easa.examiner.fe-s (`credentials/easa/examiners/fe-s.yaml`)

| Interpretation | Verdict | Paragraph(s) | Reasoning | Action |
| --- | --- | --- | --- | --- |
| `sfcl-certificate` | supported | SFCL.460(a) | SFCL.460 governs every FE(S) certificate. Recording it as the FE_S privilege on any licence except FAA and the ultralight associations is a record mapping that the text does not contradict. | None. |
| `five-years-from-valid-from` | supported | SFCL.460(a) | Five years from the start of the certificate, with a shorter recorded date winning, follows (a). | None. |
| `refresher-event` | supported | SFCL.460(b)(1) | The refresher must fall "during the validity period", which is the window the file uses; the provider cannot be recorded. | None. |
| `demonstration-event` | contradicted | SFCL.460(b)(2) | (b)(2) asks for a demonstration of the ability to conduct skill tests, proficiency checks or assessments of competence, that is, examiner ability. The evaluation counted any assessment_of_competence event, so an FI(S) assessment of competence (which the FI(S) file records the same way) or a pilot's own assessment revalidated the FE(S) certificate. The Part-FCL examiner file already restricts the same event to its examiner rating. | Count now requires `for_rating: [FE_S]`; reading updated; the existing examples record `rating: FE_S` (outcomes unchanged); new example `revalidation-instructor-assessment-does-not-count`. |
| `not-applied` | supported | SFCL.460(c), (d), (e) | Combined revalidation, renewal and the continued-compliance condition are correctly named as not evaluated. | None. |

## easa.endorsement.language-proficiency (`credentials/easa/endorsements/language-proficiency.yaml`)

| Interpretation | Verdict | Paragraph(s) | Reasoning | Action |
| --- | --- | --- | --- | --- |
| `issue-is-assessment` | supported | FCL.055(c), (a) | (c) sets re-evaluation every 4 or 6 years, and (a) requires the endorsement to state its validity date. Counting from the assessment and letting an earlier recorded date win follows both; that the record's issue date is the assessment date is a record assumption. | None. |
| `level-6-not-evaluated` | contradicted | FCL.055(c), (a) | (c) exempts level 6 from re-evaluation, and (a) is satisfied by any endorsement of level 4 or higher. Because level 6 was not selected, the licences that require this endorsement (for example the PPL(A)) found it not held and reported unknown for a level 6 holder, which the text does not support. | Level 6 is now selected with a new evaluation `level_6`: current while held, expired only once a recorded validity date has passed. Reading and `affects` updated; examples `level-6-no-expiry`, `level-6-recorded-expiry-passed`, `composite-level-6`; coverage fragment `fragments/coverage/review-easa.yaml`. Checked: a PPL(A) with a valid medical and a level 6 endorsement is now current. |
| `radio-use-and-ir-not-checked` | supported | FCL.055(a), (d), (e) | Whether the holder uses the radio, and the English requirement for IR holders, are facts the record does not hold; stating that they are not evaluated is correct. | `affects` extended to `level_6`. |

## easa.rating.powered-lift-type (`credentials/easa/ratings/powered-lift-type.yaml`)

| Interpretation | Verdict | Paragraph(s) | Reasoning | Action |
| --- | --- | --- | --- | --- |
| `part-fcl-type-ratings` | contradicted (in part) | FCL.720.PL(a)(1), (b)(1), (c)(1); FCL.740.PL(a) | Now that FCL.720.PL is stored: aeroplane pilots need a CPL/IR(A) with ATPL theory or an ATPL(A), helicopter pilots a CPL/IR(H) with ATPL theory or an ATPL/IR(H), and pilots of both categories at least a CPL(H). That backs the professional-licence requirement and the aeroplane or helicopter licence as the one the rating is held on. It does not name the MPL, which the file accepted. The other prerequisites concern first issue only and remain unevaluated; the operational suitability data may determine otherwise, which the record cannot show. | `easa.licence.mpl` removed from the `licence` requirement; reading rewritten without the "Unverified" note and citing FCL.720.PL; `licence` added to `affects`; new composite example `composite-on-mpl-unknown`. Source `sources/easa/fcl-720-pl.md` added. |

## Sources added

| File | Article | Why |
| --- | --- | --- |
| `sources/easa/fcl-720-pl.md` | FCL.720.PL | Basis of the powered-lift licence prerequisite. |
| `sources/easa/fcl-1000.md` | FCL.1000 | Examiner certificates in general; the 1-year specific certificate noted under `one-privilege-kind`. |
| `sources/easa/fcl-1020.md` | FCL.1020 | The examiner assessment of competence that FCL.1025(b)(3)(ii) and (c) refer to. |

All three are taken unchanged from the EUR-Lex consolidated text 02011R1178-20260430 via
the Publications Office Cellar, with the header of the existing `sources/easa/` files.
