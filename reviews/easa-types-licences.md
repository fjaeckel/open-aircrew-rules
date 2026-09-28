# Review: EASA type ratings and licences

Reviewed on 2026-09-28 against the verbatim texts in `sources/easa/` (EUR-Lex consolidated
version 30.04.2026 of Regulation (EU) No 1178/2011, and the consolidated Part-BFCL text of
Regulation (EU) 2018/395 already stored). Scope: the aeroplane, airship, helicopter and
powered-lift type ratings, and the PPL(A), PPL(H), PPL(As), CPL(A), CPL(H), CPL(As),
ATPL(A), ATPL(H), MPL, LAPL(A), LAPL(H) and BPL licences. This review is not a sign-off:
every interpretation keeps `approved_by: null` and `approved_on: null` until a qualified
reviewer approves it (see CONTRIBUTING.md, "Interpretations and sign-off").

## Summary

| Verdict | Count |
| --- | --- |
| supported | 47 |
| needs-decision | 10 |
| contradicted (fixed) | 2 |
| **Total** | **59** |

Contradicted and fixed:

- `easa.licence.lapl-a` / `land-sea-split`: the split of FCL.140.A(b) was applied to every
  LAPL(A) privilege, including TMG, whenever the holder had both SEP land and SEP sea. Point
  (b) concerns those two SEP privileges only. The recency is now two evaluations: `recency`
  (SEP land, SEP sea, with the split) and `recency_tmg` (TMG, without it).
- `easa.rating.aeroplane-type` / `revalidation-not-evaluated`: the reading listed "2 years"
  among possible validity periods of a type rating; FCL.740(a)(1) gives 2 years only to
  single-pilot single-engine class ratings. Wording corrected; no result changes.

Needs a decision by the signing reviewer:

- whether a CPL or ATPL holder with only a class 2 medical is reported unable to use the
  licence at all (`class-1-only`, five files);
- whether an FSTD proficiency check counts toward the 2 hours of a helicopter type rating
  (`validity-period-is-twelve-months`);
- three LAPL(A) readings (`lapl-proficiency-check`, `ul-credit-time-only`,
  `variants-sep-only`) and one LAPL(H) reading (`check-on-type`).

Sources added: `sources/easa/fcl-205-as.md`, `fcl-305.md`, `fcl-405-a.md`, `fcl-505.md`,
`fcl-700.md`, `fcl-720-h.md`, `fcl-720-as.md`.

## How to use this review

Each table has one row per interpretation of a credential file. It gives the verdict, the
paragraphs read, the reasoning and what was changed. A signing reviewer can:

1. read the cited paragraph in `sources/easa/` next to the reasoning;
2. for `supported`, confirm and sign off in the pull request as CONTRIBUTING.md describes;
3. for `needs-decision`, choose one of the options (the file's current choice is named),
   then either sign off or ask for the reading to change;
4. for `contradicted`, check the fix in the credential file and its new worked examples.

The reasoning is in our own words. Where the regulation is quoted, the quote is one short
sentence at most.

## easa.rating.aeroplane-type (`credentials/easa/ratings/aeroplane-type.yaml`)

| Interpretation | Verdict | Paragraph(s) | Reasoning | Action |
| --- | --- | --- | --- | --- |
| `other-class-is-aeroplane-type` | supported | FCL.740(a)(1); FCL.705 | This maps the record onto the rule. Class OTHER is the only place an aeroplane type rating can be recorded. Other entries recorded as OTHER are treated the same way, but only the aeroplane licences satisfy the `licence` requirement, so a non-aeroplane OTHER entry ends up unknown in the composite and never current by mistake. | None. |
| `revalidation-not-evaluated` | contradicted (wording) | FCL.740(a)(1); FCL.740.A(a), (b) | Not evaluating revalidation is sound. The record cannot tell whether FCL.740.A(a) (multi-engine and multi-pilot types) or (b) applies, and the OSD may set another period. However, the reading listed "2 years" as a possible type-rating validity. FCL.740(a)(1) gives 2 years only to single-pilot single-engine class ratings. | Reading changed to "1 year for a type rating unless the OSD determines otherwise". No results change. |
| `renewal-not-modelled` | supported | FCL.740(b) | Renewal is an assessment, possible training and a check. It does not bring the expired rating back, so leaving the rating expired until a new date is recorded matches the text. | None. |
| `passengers-not-evaluated` | supported | FCL.060(b)(1) | FCL.060(b)(1) counts take-offs and landings in the same type or class. Passenger subjects are built per class, and class OTHER does not identify a type, so a result would pool unrelated types. Leaving it out is an openly stated gap, not a misreading. | None. |

## easa.rating.airship-type (`credentials/easa/ratings/airship-type.yaml`)

| Interpretation | Verdict | Paragraph(s) | Reasoning | Action |
| --- | --- | --- | --- | --- |
| `part-fcl-type-ratings` | supported | FCL.740.As(a); FCL.205.As(a); FCL.305(a); FCL.720.As | The PPL(As) and the CPL(As) are the only airship licences. Part-FCL has no ATPL for airships, and the airship privileges come from FCL.205.As and FCL.305. FCL.720.As states experience requirements for a first issue only and names no licence, so requiring the airship licence the rating is on is right. | None. |
| `type-is-type-designator` | supported | FCL.740.As(a)(1), (a)(2) | Both the check and the 2 hours must be in "the relevant type". Matching by type designator, and returning unknown when none is recorded, is the correct cautious mapping. | None. |
| `validity-period-is-twelve-months` | supported | FCL.740(a)(1); FCL.740.As(a)(2) | Type ratings are valid for 1 year (unless the OSD says otherwise), so the validity period ends at the recorded expiry and starts 12 months earlier. The last sentence of (a)(2) lets the check count toward the 2 hours. | None. |
| `check-in-aircraft` | supported | FCL.740.As(a)(1) | (a)(1) requires the check "in the relevant type of airship" and, unlike FCL.740.H(a)(1)(ii)(A), names no FSTD. Not counting an FSTD check follows the text. | None. |
| `combined-ir-not-modelled` | supported | FCL.740.As(a)(3); FCL.740(b) | (a)(3) is permissive: it lets the IR(As) check be combined, but it is no condition of the type rating. Leaving it out does not change any type-rating result. | None. |

## easa.rating.helicopter-type (`credentials/easa/ratings/helicopter-type.yaml`)

| Interpretation | Verdict | Paragraph(s) | Reasoning | Action |
| --- | --- | --- | --- | --- |
| `part-fcl-type-ratings` | supported | FCL.740.H(a); FCL.700(a)(1); FCL.205.H; FCL.720.H | FCL.700(a)(1) exempts LAPL privileges from the type-rating requirement, so helicopter types on a LAPL(H) follow FCL.140.H. Leaving them out here is right. The PPL(H), CPL(H) and ATPL(H) hold helicopter type ratings. FCL.720.H states first-issue prerequisites only. | None. |
| `type-is-type-designator` | supported | FCL.740.H(a)(1)(i), (a)(2)(ii) | Every step is "in the relevant helicopter type". Matching by designator, with unknown when it is absent, is correct. | None. |
| `validity-period-is-twelve-months` | needs-decision | FCL.740.H(a)(1)(i), (a)(1)(ii)(A)-(B); FCL.740(a)(1) | The 12-month window is right. The question is FSTD time. (ii)(A) allows the check in an FSTD, and (ii)(B) lets the duration of that check count toward the time of point (i). Point (i), however, asks for 2 hours as pilot in the helicopter type, and (B) calls that time "flight time". Options: (1) only a check flown in the helicopter counts toward the 2 hours (the file's choice; example `revalidation-fstd-check`); (2) the duration of an FSTD check counts too, because (B) refers to the check of (A), which may be in an FSTD. Option (1) is the cautious one. Option (2) reads (B) more literally. | None; for decision. |
| `check-in-fstd` | supported | FCL.740.H(a)(1)(ii)(A) | (A) names "an FSTD representing that type" explicitly. Appendix 9 content and who conducted the check cannot be recorded. | None. |
| `light-single-engine-route` | supported | FCL.740.H(a)(2), (a)(2)(ii) | (a)(2) opens the refresher route only for single-engine types of 3 175 kg or less. Reading this from the engine count and mass on the flights, and giving unknown when they are missing, is correct. | None. |
| `refresher-as-dual` | supported | FCL.740.H(a)(2)(ii)(B) | (B) sets 1 hour with an instructor in the 3 months before expiry and allows the aircraft, an FSTD or both. Counting dual time in the type over that window matches it. Who chose the exercises cannot be recorded. | None. |
| `group-revalidation-not-applied` | supported | FCL.740.H(b), (c), (d) | (b) to (d) are further ways to revalidate, so leaving them out can only under-report. The reading says so. | None. |
| `early-check-not-credited` | supported | FCL.740(a)(1) | FCL.740(a)(1) starts the new period on the date of an early check. The holder records that new expiry date, and the file then evaluates the new period. | None. |

## easa.rating.powered-lift-type (`credentials/easa/ratings/powered-lift-type.yaml`)

The review `reviews/easa-medicals-examiners-language.md` already corrected
`part-fcl-type-ratings`. It was checked again here and is unchanged.

| Interpretation | Verdict | Paragraph(s) | Reasoning | Action |
| --- | --- | --- | --- | --- |
| `part-fcl-type-ratings` | supported | FCL.720.PL(a)(1), (b)(1), (c)(1); FCL.740.PL(a) | FCL.720.PL names only the CPL and ATPL of aeroplanes or helicopters, which matches the `licence` requirement now that the MPL is gone. The other prerequisites concern the first issue. | None. |
| `type-is-type-designator` | supported | FCL.740.PL(a)(2)(i) | The route sectors and the check are "of the relevant type". The designator mapping is right. | None. |
| `validity-period-is-twelve-months` | supported | FCL.740.PL(a)(2); FCL.740(a)(1) | The 1-year validity makes "during the period of validity" the 12 months before expiry. | None. |
| `route-sector-cruise` | supported | FCL.010 ("Route sector"); FCL.740.PL(a)(2)(ii) | FCL.010 defines a route sector as a flight with a cruise of at least 15 minutes, among other phases. (a)(2)(ii) allows one sector in the aircraft or an FFS with an examiner, during the check if wished. | None. |
| `check-in-aircraft` | supported | FCL.740.PL(a)(1) | (a)(1) requires the check "in the relevant type of powered-lift" and names no FSTD. | None. |
| `cat-exemption-not-applied` | supported | FCL.740.PL(a)(3) | The exemption depends on the holder's employer and on a combined operator check. The record cannot show either. Leaving it out can only under-report. | None. |

## easa.licence.ppl-a, ppl-h, ppl-as (`credentials/easa/licences/ppl-*.yaml`)

| File / interpretation | Verdict | Paragraph(s) | Reasoning | Action |
| --- | --- | --- | --- | --- |
| `ppl-a` / `licence-authorities` | supported | MED.A.030(c)(2) | This maps the record onto the licence. The medical accepted (class 2, or class 1 as a higher class) follows "at least a valid class 2 medical certificate". | None. |
| `ppl-h` / `licence-authorities` | supported | MED.A.030(c)(2) | As for the PPL(A). | None. |
| `ppl-as` / `licence-authorities` | supported | MED.A.030(c)(2) | As for the PPL(A). (c)(2) covers every PPL, airships included. | None. |

## easa.licence.cpl-a, cpl-h, cpl-as, atpl-a, atpl-h (`credentials/easa/licences/{cpl,atpl}-*.yaml`)

The five files carry the same two interpretations and get the same verdicts.

| File / interpretation | Verdict | Paragraph(s) | Reasoning | Action |
| --- | --- | --- | --- | --- |
| `cpl-a`, `cpl-h`, `cpl-as`, `atpl-a`, `atpl-h` / `licence-authorities` | supported | MED.A.030(c)(5) | (c)(5) requires a class 1 medical for the CPL and the ATPL. The licence mapping is a record choice. | None. |
| `cpl-a`, `cpl-h`, `cpl-as`, `atpl-a`, `atpl-h` / `class-1-only` | needs-decision | MED.A.030(c)(2), (c)(5); FCL.305(a)(1); FCL.505(a)(1) | Commercial privileges do need class 1. But FCL.305(a)(1) and FCL.505(a)(1) give the CPL and ATPL holder "all the privileges" of a PPL, and (c)(2) needs only class 2 for PPL privileges. With a class 2 medical the holder may therefore still fly as a private pilot. The file reports the whole licence as expired. Every rating whose `licence` requirement is met only by that CPL or ATPL is then reported expired as well, though it may be used privately. Options: (1) report the licence expired with a class 2 medical (the file's choice); (2) report it usable within its PPL privileges, with the commercial privileges as a limitation (the composite has no way to say this yet); (3) keep (1), but let the ratings also accept the PPL privileges of a CPL or ATPL. Option (1) is safe for commercial flying and too strict for private flying. | None; for decision. Note: the parallel change adding `easa.licence.cpl-a` to the SEP land `licence` requirement meets this case (`credentials/composite_test.go`, "sep on cpl"). |

## easa.licence.mpl (`credentials/easa/licences/mpl.yaml`)

| Interpretation | Verdict | Paragraph(s) | Reasoning | Action |
| --- | --- | --- | --- | --- |
| `licence-authorities` | supported | MED.A.030(c)(5) | (c)(5) names the MPL. The mapping is a record choice. | None. |
| `class-1-only` | supported | MED.A.030(c)(5); FCL.405.A(a), (b)(1) | Unlike the CPL and ATPL, the MPL does not include PPL privileges by itself. FCL.405.A(b)(1) grants them only on application, as additional privileges. Reporting an MPL with only a class 2 medical as not usable therefore follows the text. | None. |

## easa.licence.lapl-a (`credentials/easa/licences/lapl-a.yaml`)

| Interpretation | Verdict | Paragraph(s) | Reasoning | Action |
| --- | --- | --- | --- | --- |
| `lapl-privileges` | supported | FCL.140.A(a); FCL.105.A(a) | FCL.105.A(a) lists the privileges as SEP(land), SEP(sea) and TMG. (a) counts experience "as pilots of aeroplanes or TMGs", so pooling all aeroplane classes is right. Flights recorded as class OTHER (type-rated aeroplanes) do not count, which is stricter than the text; it is a record limit. | `affects` extended to `recency_tmg`. |
| `refresher-as-dual` | supported | FCL.140.A(a)(1)(ii), (d) | 1 hour with an instructor. The instructor's choice of exercises and the signature of (d) cannot be recorded. | `affects` extended to `recency_tmg`. |
| `lapl-proficiency-check` | needs-decision | FCL.140.A(a)(2) | (a)(2) asks for a LAPL(A) proficiency check whose programme is based on the LAPL(A) skill test. The file counts any proficiency check recorded in an aeroplane or TMG, including a PPL class-rating check under Appendix 9. Options: (1) any aeroplane proficiency check (the file's choice); (2) only a check recorded for the LAPL(A), for example with `for_rating`, reporting other checks as not counting. Option (1) assumes a class-rating check at least covers the LAPL(A) skill-test items, which the text does not say. | `affects` extended to `recency_tmg`; for decision. |
| `ul-credit-time-only` | needs-decision | FCL.035(a)(4)(i)-(ii); FCL.140.A(a)(1), (b) | Crediting ultralight hours toward the 12 hours, and not toward take-offs or landings (which are not flight time), follows (a)(4). Two parts go further than the text. First, (a)(4)(ii) contemplates training flights with an instructor in an authorised ultralight, so the refresher hour could be credited too; the file never credits it, because the authorisation cannot be recorded. Second, (a)(4) credits the hours "in full" to the same class, so they could count toward the 1 hour per class of (b); the file does not count them. Options: (1) keep both exclusions (the file's choice, cautious); (2) credit ultralight hours toward the per-class hour of (b); (3) also credit ultralight dual toward the refresher. | `affects` extended to `recency_tmg`; for decision. |
| `land-sea-split` | contradicted (fixed, in part) | FCL.140.A(a)(1), (b) | (b) applies to a holder of both a SEP(land) and a SEP(sea) privilege, and its effect is valid "for both privileges", meaning those two. The `only_if` tested only whether both were held, not which privilege was evaluated, so the TMG privilege also needed 1 hour and 6 take-offs and landings in each SEP class. A pilot who flew all 12 hours in a TMG was reported lapsed for TMG. That condition has no basis in (a). The rest of the reading stays open: for the two SEP privileges, the file requires the split whenever both are held, so 12 hours flown only on land keeps neither SEP privilege current. The other reading is that experience in one class keeps that class current, and the split is needed only to cover both. | `recency` now has `only_for: classes [SEP_LAND, SEP_SEA]`. The new evaluation `recency_tmg` has the same experience without the split. Reading rewritten; `affects: [recency, recency_tmg]`. Examples: `recency-tmg-no-land-sea-split` (current), `recency-tmg-sep-privileges-need-split` (lapsed), `recency-tmg-lapsed`; `recency-towed-tmg` moved to `recency_tmg` (outcome unchanged). Coverage fragment `fragments/coverage/review-easa-types-licences.yaml`. For the SEP part, the signing reviewer chooses between the file's reading and the per-class reading. |
| `no-expiry` | supported | FCL.140.A(a) | The LAPL(A) privileges have no validity period. "In the last 2 years" is a window looking back from the day the privileges are used. | `affects` extended to `recency_tmg`. |
| `passenger-pic-since-issue` | supported | FCL.105.A(b)(1); FCL.035(a)(4) | (b)(1) counts PIC time on aeroplanes or TMG "after the issuance of the licence". FCL.035(a)(4) credits ultralight hours only toward FCL.140.A(a)(1) and FCL.740.A(b)(1)(ii), not toward FCL.105.A. The issue-date edge is handled correctly. As in `lapl-privileges`, time recorded as class OTHER does not count. | None. |
| `previously-held` | supported | FCL.105.A(b)(2) | (b)(2) exempts holders who previously held an ATPL(A), MPL(A), CPL(A) or PPL(A), and whether that licence is still valid does not matter. Treating any such licence in the record as previously held follows it. | None. |
| `prerequisite-not-met-expired` | supported | FCL.105.A(b)(1) | Until the 10 hours are flown, carrying passengers is not allowed. Reporting that as expired is the same convention as for passenger recency. Consider citing `policy:passengers-expired` on the `not_met` stage, since the status is a convention. | None. |
| `variants-sep-only` | needs-decision | FCL.140.A(c); FCL.710(d), (da) | (c) and FCL.710(da) apply the 2-year rule only to SEP variants with a different type of engine (Article 2, point (8c)). Leaving TMG variants out is right. Every other recorded SEP variant is evaluated, because the engine type is not recorded, so a variant with the same engine type (for which (da) disables the rule) can be reported lapsed. Options: (1) evaluate every recorded SEP variant (the file's choice, cautious); (2) report such variants unknown unless the record marks a different engine type; (3) ask integrators to record only variants with a different engine type. | None; for decision. |
| `privileges-scope-not-checked` | supported | FCL.105.A(a) | The mass and occupant limits describe which flights the licence covers. They are not recency or validity conditions. | `affects` extended to `recency_tmg`. |

## easa.licence.lapl-h (`credentials/easa/licences/lapl-h.yaml`)

| Interpretation | Verdict | Paragraph(s) | Reasoning | Action |
| --- | --- | --- | --- | --- |
| `lapl-h-types` | supported | FCL.140.H(a); FCL.105.H; FCL.700(a)(1) | (a) makes the LAPL(H) recency per "specific type". FCL.700(a)(1) keeps LAPL privileges outside the type-rating system, so separating them from `easa.rating.helicopter-type` is right. | None. |
| `type-is-type-designator` | supported | FCL.140.H(a), (a)(1)(i) | Both the relevant-type wording of (a) and the "helicopters of that type" of (a)(1)(i) support matching by designator, with unknown when none is recorded. | None. |
| `takeoffs-approaches-landings` | supported | FCL.140.H(a)(1)(i) | Every take-off and landing implies an approach. The smaller of take-offs and landings is a fair measure of the six. | None. |
| `refresher-as-dual` | supported | FCL.140.H(a)(1)(ii), (b), (c) | (b) allows the aircraft, an FSTD or both. The instructor's choice and the signature of (c) cannot be recorded. | None. |
| `check-on-type` | needs-decision | FCL.140.H(a)(2), (b) | (a)(2) asks for a check "on the specific type" and does not mention an FSTD. (b) allows an FSTD for the refresher training only. That argues the check must be flown in the helicopter, which is also how the airship and powered-lift files read a check that names no FSTD. Options: (1) an FSTD check on the type counts (the file's choice); (2) only a check in the helicopter counts (drop `simulator: include` from the check). Option (2) matches the structure of FCL.140.H. | None; for decision. |
| `twelve-months-rolling` | supported | FCL.140.H(a) | No validity period applies. "In the last 12 months" looks back from the date evaluated. | None. |

## easa.licence.bpl (`credentials/easa/licences/bpl.yaml`)

| Interpretation | Verdict | Paragraph(s) | Reasoning | Action |
| --- | --- | --- | --- | --- |
| `balloons-not-modelled` | supported | BFCL.160(a), (b), (d); BFCL.115(c) | BFCL.160 is per balloon class, adds 3 hours for each further class, and ties the hot-air group to the group of the training flight or check. None of these can be read from the record. Reporting unknown every time, rather than current in classes or groups the pilot may not fly, is the only safe choice. | None. |
| `commercial-credit-not-applied` | supported | BFCL.160(f) | The credit depends on a BFCL.215 check or training flight in the relevant class, which cannot be recorded. With the recency always unknown, it would change nothing. | None. |
| `commercial-medical-not-applied` | supported | MED.A.030(c)(3)(i)-(ii) | The class 2 requirement depends on the kind of operation (commercial passenger ballooning, or more than 4 persons on board). The record does not show it. The LAPL medical of (c)(1) is the right base requirement. | None. |
| `licence-authorities` | supported | MED.A.030(c)(1) | (c)(1) names the BPL issued under Part-BFCL. Reading a LAPL(B) entry as a BPL is a record mapping for licences converted to the BPL. | None. |

## Sources added

| File | Article | Why |
| --- | --- | --- |
| `sources/easa/fcl-205-as.md` | FCL.205.As | PPL(As) privileges; basis of the airship licence requirement. |
| `sources/easa/fcl-305.md` | FCL.305 | CPL privileges include the PPL ones (`class-1-only`); airship licence requirement. |
| `sources/easa/fcl-405-a.md` | FCL.405.A | MPL privileges; PPL(A) privileges only on application (`mpl` / `class-1-only`). |
| `sources/easa/fcl-505.md` | FCL.505 | ATPL privileges include the PPL and CPL ones (`class-1-only`). |
| `sources/easa/fcl-700.md` | FCL.700 | LAPL privileges need no class or type rating (`lapl-h-types`, `part-fcl-type-ratings`). |
| `sources/easa/fcl-720-h.md` | FCL.720.H | Helicopter type-rating prerequisites: first issue only, no licence named. |
| `sources/easa/fcl-720-as.md` | FCL.720.As | Airship type-rating prerequisites: first issue only, no licence named. |

All seven are taken unchanged from the EUR-Lex consolidated text 02011R1178-20260430,
fetched again from the Publications Office Cellar on 2026-09-28 (it was byte-identical to
the copy the existing sources came from). They use the header of the existing
`sources/easa/` files. None is cited by a credential file yet; the review cites them.
