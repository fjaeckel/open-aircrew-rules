# Review: EASA class ratings, instrument ratings and the BIR

Reviewed on 2026-09-28 against the verbatim texts in `sources/easa/` (EUR-Lex consolidated
version 30.04.2026 of Regulation (EU) No 1178/2011). Scope: every interpretation of
`credentials/easa/ratings/{sep-land,sep-sea,tmg,mep-land,mep-sea,set-land,set-sea,ir-a,ir-as,ir-h,bir}.yaml`.
This review is not a sign-off: every interpretation keeps `approved_by: null` and
`approved_on: null` until a qualified reviewer approves it (see CONTRIBUTING.md,
"Interpretations and sign-off").

## Summary

| Verdict | Count |
| --- | --- |
| supported | 42 |
| needs-decision | 6 |
| contradicted (fixed) | 1 |
| contradicted (fixed at integration) | 2 |
| **Total** | **51** |

Contradicted and fixed:

- `easa.rating.set-land` / `check-in-aeroplane`: the reading excluded every check recorded
  in a simulator because FCL.740.A(b)(3) names no FSTD. Appendix 9, Section A, points 1c and
  1e allow a class rating proficiency check in an FFS (or in FSTDs combined with the
  aircraft). A check recorded as an FFS flight of the class now counts; a check event or an
  FNPT-only check still does not, because the record cannot show it was an FFS or that the
  aircraft part was flown.

Contradicted, fix deferred by this review and applied at integration:

- `easa.rating.sep-land` and `easa.rating.tmg` / `part-fcl-licences` (and the `licence`
  entry of `sep-sea`): the reading and `held_on` accept a rating on a CPL(A), ATPL(A) or
  MPL, but the `licence` requirement names only `easa.licence.ppl-a`. An SEP or TMG rating on
  a CPL(A) therefore has an `unknown` composite (licence not held) however valid the CPL(A)
  and its class 1 medical are. The fix (`requires_any: [easa.licence.ppl-a,
  easa.licence.cpl-a, easa.licence.atpl-a, easa.licence.mpl]`, as in the MEP and SET files)
  was tried and passes every worked example. It is deferred because `TestComposites` in
  `credentials/composite_test.go` pins the current result ("sep on cpl" expects `unknown`,
  `not_held`); that test is shared code outside this review's files. Proposed change to the
  test: expect `expired`, decided by `licence` with credential `easa.licence.cpl-a` (its
  record holds only a class 2 medical), which still shows that the rating cannot use the
  PPL(A) held on another licence. Fixed at integration as proposed: the three files now
  require any of the four aeroplane licences, `TestComposites` "sep on cpl" expects
  `expired` decided by `licence` with `easa.licence.cpl-a`, and the composite example
  `composite-rating-on-cpl` (expiring) in `examples/easa.rating.sep-land.yaml` shows an SEP
  rating on a CPL(A) with a class 1 medical. No existing worked example changed its result.

Needs a decision by the signing reviewer: the exempting-check period for the SEP/TMG
refresher (`exemption-period`), whether the SEP land and SEP sea pooling of FCL.740.A(b)(4)
is applied (`land-sea-pooling-not-applied`, `no-land-sea-pooling`), which checks count as an
SEP/TMG class rating check (`examiner-not-checked`), the ultralight motorglider credit to
the TMG class (`ul-credit-time-only`) and the alternate BIR check in an aeroplane
(`alternate-check-not-tracked`).

Sources added: `sources/easa/fcl-600.md`, `sources/easa/appendix-9.md` (Section A, points 1
to 1f only).

## How to use this review

Each table gives one row per interpretation of a credential file: its verdict, the
paragraphs read, the reasoning, and what was changed. A signing reviewer can:

1. read the cited paragraph in `sources/easa/` next to the reasoning;
2. for `supported`, confirm and sign off in the pull request as CONTRIBUTING.md describes;
3. for `needs-decision`, choose one of the options given (the file's current choice is
   named) and either sign off or ask for the reading to change;
4. for `contradicted`, check the fix in the credential file and its worked examples, or,
   where the fix is deferred, the proposed change.

The reasoning is in our own words; where the regulation is quoted, it is one short
sentence at most. "Other findings" at the end lists gaps that no interpretation names.

## easa.rating.sep-land (`credentials/easa/ratings/sep-land.yaml`)

| Interpretation | Verdict | Paragraph(s) | Reasoning | Action |
| --- | --- | --- | --- | --- |
| `part-fcl-licences` | contradicted (in part), fixed at integration | FCL.740.A(b)(1); FCL.740(a)(1) | Taking every EASA or LBA licence except the LAPL(A), SPL, LAPL(S) and ultralight licences is sound: FCL.740.A(b)(1) revalidates class ratings on whatever Part-FCL licence carries them, and the LAPL(A) has no class ratings. But the file's `licence` requirement accepts only the PPL(A), so a rating on a CPL(A), ATPL(A) or MPL (all in `held_on`) was reported unknown as a whole. | Fixed at integration: `licence` now `requires_any` the PPL(A), CPL(A), ATPL(A) or MPL; `TestComposites` updated; example `composite-rating-on-cpl` added. |
| `valid-through-expiry` | supported | FCL.740.A(b)(1)(i), (ii) | The expiry date is the last day of validity, and the 3 and 12 months "preceding the expiry date" end on it. Counting the expiry date inside both periods is the natural reading. | None. |
| `examiner-not-checked` | needs-decision | FCL.740.A(b)(1)(i); Appendix 9 Section A, points 1c, 1e | That the examiner and the Appendix 9 content cannot be checked is correctly stated. Two further choices sit in the reading but are not stated: (a) a check recorded in a simulator never counts, although points 1c and 1e allow a non-complex non-high-performance aeroplane or TMG check in an FFS; (b) any check event in the class counts, including one recorded for the IR or the BIR only, which is not a class rating check. Options for (a): (1) aeroplane only (file's choice, cautious); (2) also a check recorded as an FFS flight of the class, as now done for SET. Options for (b): (1) any check in the class (file's choice; a combined class and IR/BIR check is recorded as one event for the IR/BIR); (2) exclude checks recorded for IR or BIR, which needs a new filter because a class rating check is often recorded without a rating. | Decided 2026-09-28: (a) option 2 (an FFS check of the class counts) and (b) option 2 (checks recorded only for the IR or BIR excluded), applied: `simulator: include, fstd: [FFS]`, `excluding_ratings: [IR, BIR]`. |
| `refresher-as-dual` | supported | FCL.740.A(b)(1)(ii)(C) | One hour of dual time received in the class stands for the refresher flown with an FI or CRI; the instructor's certificate and the exercises chosen cannot be recorded. With the SEP/TMG pool the refresher may be flown in either class, which (b)(2) allows. | None. |
| `exemption-period` | needs-decision | FCL.740.A(b)(1)(ii)(C) | The exemption names the checks, tests and assessments but no period. Options: (1) within the 12 months before expiry, the period of point (ii) (file's choice); (2) within the rating's validity period (24 months); (3) at any time since the rating was last revalidated or issued. Option (1) is the cautious reading and matches the period in which the rest of (ii) is flown. | Decided 2026-09-28: option 2 (the 24-month validity period), applied: `within_months_before_expiry: 24` on the waiver. |
| `exemption-any-aeroplane` | supported | FCL.740.A(b)(1)(ii)(C) | "In any class or type of aeroplane" covers every aeroplane class; TMG class ratings are revalidated under the aeroplane point FCL.740.A(b)(1), so counting TMG checks is consistent. See "Other findings" on checks recorded for the IR or BIR only. | None. |
| `pooling-same-licence` | supported | FCL.740.A(b)(2) | (b)(2) lets the holder of both ratings complete point (1), the proficiency check included, in either class or a mix, and revalidate both. Requiring both ratings on the same licence is the ordinary case of "hold both". | None. |
| `land-sea-pooling-not-applied` | needs-decision | FCL.740.A(b)(4) | (b)(4) lets a holder of SEP land and SEP sea complete the experience route in either class or a mix, with at least 1 hour of PIC time and 6 take-offs and 6 landings in each class. Not applying it means a pilot who met (b)(4) is reported as not met and, after the expiry date, as expired, although the text treats the requirements as fulfilled for both ratings. Options: (1) not applied (file's choice; never reports a rating valid that is not); (2) apply it as a third route (12 h, 6 h PIC and 12/12 over both classes, with the per-class minima, only while both ratings are held on the same licence), which the vocabulary can express with `classes` and `only_if: { holds: ... }` but which changes `easa.rating.sep-land#revalidation` and every file that uses it. | Decided 2026-09-28: option 2 (a third route over both classes with the per-class minima), applied; the TMG rating now has its own copy of the evaluation without that route. |
| `ul-credit-time-only` | needs-decision | FCL.035(a)(4)(i), (ii); FCL.010 (TMG) | Crediting Annex I aeroplane hours to the flight time and PIC time only is supported: (a)(4) credits the flight time requirements, whereas (a)(5) for gyroplanes names take-offs and landings expressly. Leaving the refresher out is cautious, since (a)(4)(ii) requires an authorised aircraft for training flights, which the record cannot show. Open point: FCL.010 defines a TMG by an integrally mounted, non-retractable engine and propeller, and (a)(4)(i) credits only the same class. Options: (1) every ultralight motorglider counts as TMG (file's choice); (2) no ultralight motorglider time is credited to the TMG class, because the record cannot show that the engine is fixed. | Decided 2026-09-28: option 2 (no ultralight motorglider time unless the record shows a fixed engine and propeller), applied: `ul_credit` with `fixed_engine: true`. |

## easa.rating.sep-sea (`credentials/easa/ratings/sep-sea.yaml`)

| Interpretation | Verdict | Paragraph(s) | Reasoning | Action |
| --- | --- | --- | --- | --- |
| `no-land-sea-pooling` | needs-decision | FCL.740.A(b)(4); FCL.035(a)(4) | Same decision as `land-sea-pooling-not-applied` above; the two files should decide together. Crediting no ultralight time to the sea class is cautious: the record has no sea ultralight kind, so the same-class condition of FCL.035(a)(4)(i) cannot be shown. That the SEP/TMG pool does not include SEP sea is correct, as (b)(2) names SEP land only. | Decided 2026-09-28: option 2 (as for SEP land), applied through the shared SEP land evaluation. |

## easa.rating.tmg (`credentials/easa/ratings/tmg.yaml`)

| Interpretation | Verdict | Paragraph(s) | Reasoning | Action |
| --- | --- | --- | --- | --- |
| `part-fcl-licences` | contradicted (in part), fixed at integration | FCL.740.A(b)(1) | Separating the TMG class rating from the SPL's TMG privileges is right (the latter are Part-SFCL). As for SEP land, the `licence` requirement accepts only the PPL(A), so a TMG rating on a CPL(A) or ATPL(A) was reported unknown. | Fixed at integration; same change as SEP land. |

## easa.rating.mep-land (`credentials/easa/ratings/mep-land.yaml`)

| Interpretation | Verdict | Paragraph(s) | Reasoning | Action |
| --- | --- | --- | --- | --- |
| `part-fcl-licences` | supported | FCL.740.A(a) | Multi-engine class ratings are held on the PPL(A), CPL(A), ATPL(A) and MPL, which is what the `licence` requirement lists; the LAPL(A) cannot carry one. | None. |
| `validity-period-is-twelve-months` | supported | FCL.740(a)(1); FCL.740.A(a)(2) | Only single-pilot single-engine class ratings get 2 years, so an MEP rating has 1 year and "the period of validity" is the 12 months before expiry. | None. |
| `check-and-sectors-both-required` | supported | FCL.740.A(a)(1), (a)(2) | Points (1) and (2) are joined by "and"; there is no experience route for multi-engine ratings. | None. |
| `check-in-any-fstd` | supported | FCL.740.A(a)(1); FCL.740(a)(2) | (a)(1) allows the check in "an FSTD representing that class", without naming a level, and EBT practical assessment counts in full. See "Other findings" on IR-only checks. | None. |
| `route-sector-cruise` | supported | FCL.010 (route sector); FCL.740.A(a)(2)(i) | The definition requires take-off, departure, a cruise of at least 15 minutes, approach and landing; a logged flight with a recorded cruise of 15 minutes has these phases. "As pilot" names no role. Reporting unknown when cruise times are missing is the right handling of missing input. A logbook line with several legs counts once, which can only under-count. | None. |
| `examiner-sector` | supported | FCL.740.A(a)(2)(ii) | (ii) accepts one sector in the aeroplane or an FFS with an examiner and says it may be flown during the check. | None. |
| `early-check-not-credited` | supported | FCL.740(a)(1) | An early check starts a new validity period from the check date, which only a newly recorded expiry date can express. Not crediting it against the old date can only under-report and is stated openly. | None. |
| `cat-exemption-not-applied` | supported | FCL.740.A(a)(3) | Whether the holder works for a CAT operator and passed a combined operator check is not in the record; the exemption can only shorten the requirements, so not applying it is cautious. | None. |
| `multi-engine-types-separate` | supported | FCL.740.A(a) | A scope statement: (a) covers class and type ratings alike, and type ratings are another credential file. | None. |

## easa.rating.mep-sea (`credentials/easa/ratings/mep-sea.yaml`)

| Interpretation | Verdict | Paragraph(s) | Reasoning | Action |
| --- | --- | --- | --- | --- |
| `part-fcl-licences` | supported | FCL.740.A(a) | As for MEP land. There is no land/sea pooling for multi-engine classes, so each class revalidating from its own check and sectors follows the text. | None. |

## easa.rating.set-land (`credentials/easa/ratings/set-land.yaml`)

| Interpretation | Verdict | Paragraph(s) | Reasoning | Action |
| --- | --- | --- | --- | --- |
| `part-fcl-licences` | supported | FCL.740.A(b)(3) | As for MEP land; the `licence` requirement lists all four aeroplane licences. | None. |
| `check-only` | supported | FCL.740.A(b)(3) | (b)(3) names only the proficiency check in the 3 months before expiry; there is no experience route. | None. |
| `check-in-aeroplane` | contradicted (fixed) | FCL.740.A(b)(3); Appendix 9 Section A, points 1c, 1e | The reasoning that (b)(3) names no FSTD is right, but (b)(3) requires the check to follow Appendix 9, and point 1c there allows a class rating check for single-pilot aeroplanes in an FFS, in FSTDs combined with the aircraft, or in the aircraft. Excluding an FFS check reported a revalidated rating as not met. | `passes_if` now also accepts a flight in an FFS of the class flagged as a proficiency check (row `ffs_check`, refs FCL.740.A(b)(3) and Appendix 9); reading rewritten to say why check events and other FSTDs recorded in a simulator still do not count. Example `revalidation-simulator-check` now expects the FFS check to revalidate (outcome `expiring`, `rating.revalidation_expiring_met`); new example `revalidation-fnpt-check-does-not-count`. Source `sources/easa/appendix-9.md` added. |
| `validity-two-years` | supported | FCL.740(a)(1) | SET is a single-pilot single-engine class rating, so 2 years unless the OSD says otherwise; using the recorded expiry covers an OSD exception. | None. |
| `early-check-not-credited` | supported | FCL.740(a)(1) | As for MEP land. | None. |

## easa.rating.set-sea (`credentials/easa/ratings/set-sea.yaml`)

| Interpretation | Verdict | Paragraph(s) | Reasoning | Action |
| --- | --- | --- | --- | --- |
| `part-fcl-licences` | supported | FCL.740.A(b)(3) | As for SET land; the check is taken in the relevant class, so SET sea checks only. The FFS change above applies here too through `uses:`. | None. |

## easa.rating.ir-a (`credentials/easa/ratings/ir-a.yaml`)

| Interpretation | Verdict | Paragraph(s) | Reasoning | Action |
| --- | --- | --- | --- | --- |
| `ir-on-aeroplane-licence` | supported | FCL.600; FCL.625.A(a); FCL.625(b)(1) | An IR is tied to a category of aircraft (FCL.600), and the record holds no category for it, so taking the category from the licence it is recorded on is the available mapping. | None. |
| `ir-check-recorded` | supported | FCL.625.A(a)(2), (a)(3) | Whether combined with a class or type check or not, the IR is revalidated by a check that covers the IR sections; requiring the check to be recorded for the IR, and not counting a class rating check that is not, follows from that. | None. |
| `fstd-level-not-checked` | supported | FCL.625.A(a)(4) | The FNPT II/FFS level and the alternate check in an aeroplane apply only to a check not combined with a class or type rating check, which the record cannot tell apart; the limitation is stated. | None. |
| `class-rating-held` | supported | FCL.625.A(a)(1) | Stated openly; the class or type rating is evaluated by its own credential, and without it the IR cannot be used in that class anyway. | None. |
| `early-check-not-credited` | supported | FCL.625(b)(2) | As for MEP land, under the IR's own early-revalidation paragraph. | None. |
| `renewal-not-modelled` | supported | FCL.625(c), (d) | Renewal and the 7-year rule concern what an expired IR needs; reporting it expired until a new date is recorded matches the text. | None. |
| `cross-credit-not-applied` | supported | FCL.625.A(b) | Appendix 8 cross-credit can only reduce what a check must contain; not applying it is cautious. | None. |

## easa.rating.ir-as (`credentials/easa/ratings/ir-as.yaml`)

| Interpretation | Verdict | Paragraph(s) | Reasoning | Action |
| --- | --- | --- | --- | --- |
| `ir-on-airship-licence` | supported | FCL.600; FCL.625.As; FCL.625(b)(1) | As for IR(A), mapped from the airship licence. | None. |
| `ir-check-recorded` | supported | FCL.625.As(a), (b) | Both routes are a proficiency check for the type (combined, or sections 5 and 1 alone); an FTD 2/3 or FFS may be used only for the uncombined check. The unchecked conditions are named. | None. |
| `early-check-not-credited` | supported | FCL.625(b)(2) | As for IR(A). | None. |
| `renewal-not-modelled` | supported | FCL.625(c), (d) | As for IR(A). | None. |

## easa.rating.ir-h (`credentials/easa/ratings/ir-h.yaml`)

| Interpretation | Verdict | Paragraph(s) | Reasoning | Action |
| --- | --- | --- | --- | --- |
| `ir-on-helicopter-licence` | supported | FCL.600; FCL.625.H(a); FCL.625(b)(1) | As for IR(A), mapped from the helicopter licence. | None. |
| `ir-check-recorded` | supported | FCL.625.H(a)(1)-(3) | As for IR(A); the type rating and the Appendix 9 content are named as not checked. | None. |
| `fstd-level-not-checked` | supported | FCL.625.H(b) | As for IR(A), with FTD 2/3 or FFS and a helicopter for every alternate uncombined check. | None. |
| `early-check-not-credited` | supported | FCL.625(b)(2) | As for IR(A). | None. |
| `renewal-not-modelled` | supported | FCL.625(c), (d); FCL.625.H(c) | As for IR(A), including the Appendix 8 cross-credit. | None. |

## easa.rating.bir (`credentials/easa/ratings/bir.yaml`)

| Interpretation | Verdict | Paragraph(s) | Reasoning | Action |
| --- | --- | --- | --- | --- |
| `bir-as-privilege` | supported | FCL.835(g)(1) | One year from the start of the period, with an earlier recorded date winning, follows (g)(1); after an early check (g)(4) the new period appears through the recorded dates. | None. |
| `validity-period-is-twelve-months` | supported | FCL.835(g)(1), (g)(2)(ii) | With a 1-year validity, "within the validity period" is the 12 months before expiry. | None. |
| `ifr-pic-time` | supported | FCL.835(g)(2)(ii) | IFR time and approaches on flights logged as PIC is the record's best view of "6 hours as PIC under IFR". Worth knowing: a flight with some PIC time counts all of its IFR time, and only aeroplane flights count, although the text names no category; both are stated in the reading. | None. |
| `training-flight-as-dual` | supported | FCL.835(g)(2)(ii) | One flight of at least 1 hour with an instructor; the instructor's BIR training privileges cannot be recorded. | None. |
| `bir-check-recorded` | supported | FCL.835(g)(2)(i), (g)(8) | A check recorded for the BIR, possibly combined with a single-pilot class rating check, as (g)(8) allows. | None. |
| `alternate-check-not-tracked` | needs-decision | FCL.835(g)(3) | (g)(3) makes every second revalidation a check in an aeroplane, so the experience route is not always available. Offering it always can report a BIR current when the text requires a check. Options: (1) not tracked (file's choice); (2) offer the experience route only when a BIR check in an aeroplane is recorded in the 24 months before expiry (the previous revalidation was a check); (3) as (2), but report unknown instead of not met when no such check is recorded, since earlier history may be missing. | Decided 2026-09-28: option 3 (experience route only after a BIR check in an aeroplane in the 24 months before expiry; unknown without one), applied. |
| `not-met-stays-valid` | supported | FCL.835(g)(1), (g)(5) | The BIR stays valid until its expiry date; only a failed check (g)(5) stops the privileges earlier, and that is not recorded. | None. |
| `multi-engine-and-renewal-not-modelled` | supported | FCL.835(g)(7), (g)(4), (g)(6) | The record does not say whether a BIR is multi-engine, so (g)(7) cannot be applied; renewal and early revalidation work through recorded dates. The (g)(7) gap can over-credit a single-engine check for a multi-engine BIR and is stated. | None. |

## Other findings

Gaps that no interpretation names; for the signing reviewer to decide whether they need one.

- **IR-only and BIR-only checks as class rating checks.** `easa.rating.sep-land#revalidation`
  (and its users), `easa.rating.mep-land#revalidation` and `easa.rating.set-land#revalidation`
  count every proficiency check event in the class, whatever rating it was recorded for. A
  check recorded for the IR (FCL.625.A(a)(3): IR sections only) or the BIR therefore also
  revalidates the class rating, and in SEP land it also exempts from the refresher, where
  FCL.740.A(b)(1)(ii)(C)(1) names "a class or type rating proficiency check". Since a
  combined class and IR/BIR check (FCL.740.A(a)(4), (b)(5)) is recorded as one event for the
  IR or BIR, excluding such events would miss combined checks. A fix needs either a record
  convention for combined checks or a filter that excludes checks recorded only for IR/BIR.
- **SEP/TMG licence requirement.** See the summary (fixed at integration).

## Sources added

| File | Article | Why |
| --- | --- | --- |
| `sources/easa/fcl-600.md` | FCL.600 | An IR is appropriate to a category of aircraft; the basis for mapping an IR to IR(A), IR(H) or IR(As) by licence. |
| `sources/easa/appendix-9.md` | Appendix 9, Section A, points 1 to 1f (extract) | Where class rating checks for single-pilot aeroplanes may be conducted (FFS, FSTDs with the aircraft, or the aircraft); basis of the `check-in-aeroplane` fix. Only this extract is stored. |

Both are taken unchanged from the EUR-Lex consolidated text 02011R1178-20260430 via the
Publications Office Cellar, with the header of the existing `sources/easa/` files.
