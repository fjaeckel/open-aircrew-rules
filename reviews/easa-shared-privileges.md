# Review: EASA shared evaluations and privileges

Reviewed on 2026-09-28 against the verbatim texts in `sources/easa/` (EUR-Lex consolidated
versions 30.04.2026 of Regulation (EU) No 1178/2011 and 15.11.2021 of Regulation (EU)
2018/1976). Scope: every interpretation in `credentials/easa/shared/*.yaml` (6 files, 23
interpretations) and `credentials/easa/privileges/*.yaml` (8 files, 34 interpretations). This
review is not a sign-off: every interpretation keeps `approved_by: null` and
`approved_on: null` until a qualified reviewer approves it (see CONTRIBUTING.md,
"Interpretations and sign-off").

## Summary

| Verdict | Count |
| --- | --- |
| supported | 42 |
| needs-decision | 11 |
| contradicted (fixed) | 2 |
| contradicted (fixed at integration) | 2 |
| **Total** | **57** |

Contradicted and fixed:

- `easa.shared.variants` / `refresher-any-variant` (now `refresher-sep-variants`): refresher
  training restored a variant of any rating, although FCL.710(d)(3) allows it only for SEP
  variants. The restoring event is now limited to SEP land and SEP sea ratings. No worked
  example changes its result (the only refresher examples are SEP land ones).
- `easa.privilege.spl-tmg` / `extension-credit`: TMG privileges on a LAPL(A) credited the
  TMG extension as if they were a TMG class rating. SFCL.150(c)(1) credits a class rating;
  a LAPL(A) has none and is credited only under (c)(2), with the FCL.140.A recency the file
  does not check. The LAPL(A) was removed from the credit; two examples added.

Contradicted, fix deferred by this review and applied at integration (the fix needed changes
to credential files outside this review):

- `easa.shared.passengers-day` / `class-for-type` and `easa.shared.passengers-night` /
  `class-for-type`: for type-rated helicopters, powered-lift aircraft and airships, flights
  in any type of the class counted; FCL.060(b)(1) and (b)(2)(i) ask for the same type.
  Fixed at integration as proposed below.

Needs a decision by the signing reviewer: the start of an instructor certificate's 3 years
(`period-start`), what "holds an IR" requires (`ir-held`), the untracked passenger
competence flight (`competence-flight-not-recorded`), SEP variants with the same engine
type (`same-engine-sep-variants-evaluated`), dual cloud flights and the duration of a
restoring check (`dual-counts`, `check-restores`), which tows count (`tows-counted`, four
files) and the aircraft of the 12 hours of SFCL.160(b)(1) (`twelve-hours-on-sailplanes`).

Sources added: `sources/easa/fcl-900.md`, `sources/easa/fcl-915.md`.

### Deferred fix: passengers in the same type (fixed at integration)

Applied at integration as proposed here. Change, following the FAA pattern (`credentials/faa/shared/passengers-day.yaml`):

- `easa.shared.passengers-day` and `easa.shared.passengers-night`: add a parameter
  `type_rated` ("false for class ratings, true for type ratings"), and set
  `only_for: { type_rated: $type_rated }` and `counting: { ..., in_type: $type_rated }`.
  The engine already builds one passengers subject per type designator when `type_rated`
  is true (`engine/subjects.go`).
- Pass `with: { type_rated: true }` from `easa.rating.helicopter-type`,
  `easa.rating.powered-lift-type`, `easa.rating.airship-type` and `easa.licence.lapl-h`,
  and `with: { type_rated: false }` from `easa.rating.sep-land`, `sep-sea`, `set-land`,
  `set-sea`, `mep-land`, `mep-sea`, `tmg`, `easa.licence.lapl-a` and `easa.licence.gpl`.
- Reword both `class-for-type` readings: "the recorded class for class ratings, the
  recorded type designator for type ratings; a flight without a type designator is
  unknown input". FCL.060(b)(5) (cross-type credit for similar non-complex helicopters)
  stays unapplied under `no-cross-credit`.

Affected examples (the `passengers_day` and `passengers_night` examples of
`easa.rating.helicopter-type`, `easa.rating.powered-lift-type`, `easa.rating.airship-type`
and `easa.licence.lapl-h`, four each): all record flights in the rated type, and none
changed its outcome; their subject now carries the type designator as detail. Each of the
four files has a new example `passengers-day-other-type` (expired): 3 take-offs and landings
in another type of the class do not count. The class-rated users' examples are unchanged.

## How to use this review

Each table gives one row per interpretation of a credential file: its verdict, the
paragraphs read, the reasoning, and what was changed. A signing reviewer can:

1. read the cited paragraph in `sources/easa/` next to the reasoning;
2. for `supported`, confirm and sign off in the pull request as CONTRIBUTING.md describes;
3. for `needs-decision`, choose one of the options given (the file's current choice is
   named) and either sign off or ask for the reading to change;
4. for `contradicted`, check the fix in the credential file and its new worked example, or,
   when the fix is deferred, the proposal above.

A shared evaluation is used by many credentials (`docs/dependencies.md`); a decision on a
shared interpretation applies to all of them. The reasoning is in our own words; where the
regulation is quoted, it is one short sentence at most.

## easa.shared.instructor-certificate-validity (`credentials/easa/shared/instructor-certificate-validity.yaml`)

| Interpretation | Verdict | Paragraph(s) | Reasoning | Action |
| --- | --- | --- | --- | --- |
| `period-start` | needs-decision | FCL.940 | FCL.940 fixes the length (3 years) but not the start, and revalidation starts a new period. Options: (1) derive 3 years from validFrom, else the issue date, and cap a later recorded expiry (the file's choice); (2) trust the recorded expiry and only fall back to 3 years when none is recorded. Option (1) reports a revalidated certificate as expired when the record keeps the original issue date and no validFrom; option (2) accepts a mistyped expiry beyond 3 years. | None; for decision. |
| `shorter-validities-not-applied` | supported | FCL.940; FCL.900(b)(1); FCL.915(e)(2) | FCL.940 reserves FCL.900(b)(1), a specific certificate valid at most 1 year, and FCL.915(e)(2), the yearly refresher that conditions the UPRT instructing privileges (a privilege condition rather than a certificate validity). Neither is visible in the record; relying on a shorter recorded expiry is the only way to reflect them. The MI exception does not arise, as no MI credential uses this evaluation. | None. |
| `validity-not-revalidation` | supported | FCL.940 | FCL.940 only sets the period; the revalidation conditions are in the certificate-specific points. Saying that a current result is not a revalidation is accurate. | None. |

## easa.shared.medical-certificate (`credentials/easa/shared/medical-certificate.yaml`)

| Interpretation | Verdict | Paragraph(s) | Reasoning | Action |
| --- | --- | --- | --- | --- |
| `higher-class-counts` | supported | MED.A.030(c)(1), (c)(2), (c)(5) | (c)(1) and (c)(2) ask for "at least" an LAPL or class 2 certificate, so a higher class satisfies them; (c)(5) asks for class 1, where no higher class exists. | None. |
| `entered-expiry` | supported | MED.A.030(c) | (c) requires a valid certificate. The derived validity of MED.A.045 is evaluated by the medical credential, which every licence using this evaluation also requires (`docs/dependencies.md`), so a certificate recorded without an expiry is still checked in the licence composite. | None. |

## easa.shared.passengers-day (`credentials/easa/shared/passengers-day.yaml`)

| Interpretation | Verdict | Paragraph(s) | Reasoning | Action |
| --- | --- | --- | --- | --- |
| `class-for-type` | contradicted (fixed at integration) | FCL.060(b)(1) | The text asks for an aircraft "of the same type or class". For a type-rated helicopter, powered-lift aircraft or airship the relevant unit is the type, yet the reading pooled every type of the class, so 3 landings in one helicopter type allowed passengers in another. | Fixed at integration: parameter `type_rated` (`only_for` and `in_type`), passed true by the four type-rated credentials and false by the nine class-rated ones; reading reworded; examples `passengers-day-other-type` added (see "Deferred fix" above). |
| `approaches-with-landings` | supported | FCL.060(b)(1) | Every landing is preceded by an approach, and logbooks do not record VFR approaches. Counting the smaller of take-offs and landings never credits more than the text asks. | None. |
| `pilot-flying-recorded` | supported | FCL.060(b)(1) | The text requires the take-offs and landings "as a pilot flying"; a flight that does not say so cannot count, and unknown is the honest result. | None. |
| `operations-not-checked` | supported | FCL.060(b)(1) | The single-pilot or multi-pilot condition depends on the holder's privileges and on how the flight was operated, neither recorded; the reading names it as not checked. | None. |
| `no-cross-credit` | supported | FCL.060(b)(3), (b)(4), (b)(5) | (b)(3) concerns cruise relief co-pilots, a different operation; (b)(4) and (b)(5) only add credit through the operational suitability data. Not applying them never reports current wrongly. | None. |

## easa.shared.passengers-night (`credentials/easa/shared/passengers-night.yaml`)

| Interpretation | Verdict | Paragraph(s) | Reasoning | Action |
| --- | --- | --- | --- | --- |
| `ir-held` | needs-decision | FCL.060(b)(2)(ii) | The text says only that the pilot "holds an IR". Options: (1) an IR on the same licence, not expired, one without an expiry counting (the file's choice); (2) any IR the pilot holds, of any category and whatever its validity; (3) an IR of the same aircraft category, valid. Option (1) excludes an IR(H) for aeroplane passengers and an expired IR, which the text does not do explicitly. | None; for decision. |
| `separate-flights` | supported | FCL.060(b)(2)(i) | The text asks for 1 take-off, approach and landing at night without requiring them on one flight; the approach goes with the landing as by day. | None. |
| `class-for-type` | contradicted (fixed at integration) | FCL.060(b)(2)(i) | As by day: the same type is required for type-rated aircraft. | Fixed at integration, as by day. |
| `pilot-flying-recorded` | supported | FCL.060(b)(2)(i) | As by day. | None. |
| `night-privilege-separate` | supported | FCL.060(b)(2) | FCL.060 sets recency only; whether the holder may fly at night is a matter of the licence and night rating. | None. |

## easa.shared.sfcl-passenger-prerequisite (`credentials/easa/shared/sfcl-passenger-prerequisite.yaml`)

| Interpretation | Verdict | Paragraph(s) | Reasoning | Action |
| --- | --- | --- | --- | --- |
| `competence-flight-not-recorded` | needs-decision | SFCL.115(a)(2)(ii)(A), (d) | The training flight is a mandatory part of (A), and (d) requires it to be entered in the logbook, so a record could hold it. Options: (1) show it as an untracked row and decide on the hours or launches alone (the file's choice, which can report current without the flight); (2) report unknown until a training flight event is recorded, which needs a record field or event kind; (3) keep option (1) but say in the result that the flight is assumed. | None; for decision. |
| `since-recorded-issue` | supported | SFCL.115(a)(2)(ii)(A); SFCL.160(a), (b) | The text counts experience "after the issue of the SPL" as PIC "on sailplanes". Part-SFCL writes "sailplanes, excluding TMGs" wherever TMGs are left out (SFCL.160(a)), so here sailplanes include TMGs, and a TMG take-off stands for a launch ("launches or take-offs and landings"). A recorded re-issue date is later than the first issue, which can only delay the result. | None. |
| `fi-s-on-any-licence` | supported | SFCL.115(a)(2)(ii)(B) | (B) asks only that the holder holds an FI(S) certificate; the licence it is recorded on is irrelevant. | Changed at integration (deferred item of `reviews/easa-spl-gpl-instructors.md`): the FI(S) no longer needs to be unexpired by its recorded date, as Part-SFCL gives it no validity period; reading reworded; example `passenger-prerequisite-fi-s-recorded-expiry-passed` added to `examples/easa.licence.spl.yaml`. |

## easa.shared.variants (`credentials/easa/shared/variants.yaml`)

| Interpretation | Verdict | Paragraph(s) | Reasoning | Action |
| --- | --- | --- | --- | --- |
| `variant-is-recorded-variant` | supported | FCL.710(d) | (d) concerns pilots who extended their privileges to a variant, which is what a variant recorded on the rating means; a flight without a variant name cannot be attributed. | None. |
| `restoring-events-for-two-years` | supported | FCL.710(d)(1), (d)(2), (e) | Each of the listed completions is flown in the variant, so it ends the "not flown within the preceding 2 years" state as a flight would. Not checking the provider and signature of (b), (c) and (e) is stated. | None. |
| `refresher-any-variant` (now `refresher-sep-variants`) | contradicted (fixed) | FCL.710(d)(3) | (d)(3) makes refresher training available only for a variant within the SEP class rating with a particular engine type. The file accepted it for every variant, including MEP, SET and type-rating variants, where only (d)(1) and (d)(2) apply. The engine type cannot be told apart, but the class can. | Restoring event limited to `classes: [SEP_LAND, SEP_SEA]`; interpretation renamed and reworded. No example outcome changes (the refresher examples are SEP land). Fixed at integration: example `variants-refresher-does-not-restore` in `examples/easa.rating.mep-land.yaml` (lapsed) shows that refresher training does not restore an MEP variant. |
| `same-engine-sep-variants-evaluated` | needs-decision | FCL.710(da), (a), (d)(3) | (da) removes the 2-year rule for SEP variants with the same engine type, and the record does not hold the engine type. Options: (1) evaluate every recorded SEP variant (the file's choice; a same-engine variant not flown in 2 years is reported lapsed although (d) does not apply); (2) evaluate only variants named after an engine type (piston, electric, hybrid, per Article 2 point (8c)); (3) report SEP variants unknown. The existing examples record engine-type variants (ELECTRIC), which fit option (1) and (2) alike. | None; for decision. |
| `initial-training-not-logged` | supported | FCL.710(a), (d) | A variant on the rating means the extension of (a) took place; without a later flight or restoring event, (d) applies whether or not the initial training is recorded. | None. |

## easa.privilege.cloud-flying (`credentials/easa/privileges/cloud-flying.yaml`)

| Interpretation | Verdict | Paragraph(s) | Reasoning | Action |
| --- | --- | --- | --- | --- |
| `ifr-time-is-cloud-time` | supported | SFCL.215(e) | Cloud flying is flight by reference to instruments, and the record has no other field for it; using IFR time on GLIDER flights is a record mapping the text does not contradict. Cloud flights in a TMG with the engine stopped ((a)(1)) are not counted, which can only under-count. | None. |
| `dual-counts` | needs-decision | SFCL.215(e), (f)(2) | (e) counts flights "as PIC exercising cloud flying privileges"; (f)(2) lets a holder who fell short make up the missing time or flights with an FI(S). Options: (1) count PIC and dual cloud flights alike (the file's choice; a dual flight before a shortfall counts too, and the instructor's qualification is not checked); (2) count PIC flights only for (e), and dual flights only as the (f)(2) make-up after a shortfall. The results differ only for a pilot who flies dual cloud time while still current. | None; for decision. |
| `check-restores` | needs-decision | SFCL.215(f)(1) | (f)(1) restores the privileges by a proficiency check with an FE(S) but does not say which check or for how long. Options: (1) a check recorded for CLOUD_FLYING, counting for 24 months (the file's choice); (2) the same check, restoring until the (e) experience can next be met (a shorter period); (3) any sailplane proficiency check with an FE(S). | None; for decision. |
| `bir-ir-credit` | supported | SFCL.215(g) | (g) credits a BIR or an IR(A) issued under Part-FCL. Excluding FAA and helicopter IRs follows; an IR recorded on its own IR licence as an IR(A) is a record mapping. | None. |

## easa.privilege.fcl-banner-towing (`credentials/easa/privileges/fcl-banner-towing.yaml`)

| Interpretation | Verdict | Paragraph(s) | Reasoning | Action |
| --- | --- | --- | --- | --- |
| `part-fcl-licences-only` | supported | FCL.805(a); SFCL.205(a) | FCL.805(a) applies to licences with aeroplane or TMG privileges under Part-FCL; the SPL rating is set by SFCL.205. | None. |
| `tows-counted` | needs-decision | FCL.805(e) | (e) asks for "a minimum of 5 tows" to exercise "the sailplane or banner towing ratings" without saying whether a tow of the other kind counts. Options: (1) banner tows only for the banner rating (the file's choice); (2) any tow, as the two SFCL files read SFCL.205(f). The two families should read the parallel texts the same way. | None; for decision. |
| `trained-aircraft-not-checked` | supported | FCL.805(d) | The limitation to the aircraft and towing method trained is a scope of the privilege, not recency; the reading says it is not checked. | None. |
| `recorded-expiry` | supported | FCL.805(a) | FCL.805 sets no validity; honouring a recorded end date is the conservative handling of the holder's entry. | None. |
| `resumption-with-instructor` | supported | FCL.805(f) | (f) requires the missing tows with or under the supervision of an instructor, which is the remedy named. | None. |

## easa.privilege.fcl-sailplane-towing (`credentials/easa/privileges/fcl-sailplane-towing.yaml`)

| Interpretation | Verdict | Paragraph(s) | Reasoning | Action |
| --- | --- | --- | --- | --- |
| `part-fcl-licences-only` | supported | FCL.805(a); SFCL.205(a) | As for banner towing. | None. |
| `tows-counted` | needs-decision | FCL.805(e) | As for banner towing: options (1) sailplane tows only (the file's choice), (2) any tow. | None; for decision. |
| `trained-aircraft-not-checked` | supported | FCL.805(d) | As for banner towing (the towing-method limit of (d) concerns banners only). | None. |
| `recorded-expiry` | supported | FCL.805(a) | As for banner towing. | None. |
| `resumption-with-instructor` | supported | FCL.805(f) | As for banner towing. | None. |

## easa.privilege.mountain (`credentials/easa/privileges/mountain.yaml`)

| Interpretation | Verdict | Paragraph(s) | Reasoning | Action |
| --- | --- | --- | --- | --- |
| `part-fcl-authorities` | supported | FCL.815(a) | (a) lets LAPL and PPL holders obtain the rating (and the professional licences include PPL privileges); the licence requirement of the file names the aeroplane licences. Not checking the licence kind on the privilege itself is a record mapping. | None. |
| `mountain-landings-recorded` | supported | FCL.815(d)(1), (a) | (d)(1) counts landings on a designated surface; (a) limits the privilege to aeroplanes and TMGs. Unknown where the record cannot tell is the repository convention. | None. |
| `check-recorded` | supported | FCL.815(d)(2) | A check recorded for MOUNTAIN is the proficiency check of (d)(2); its content and signature are stated as not checked. | None. |
| `wheels-and-skis-not-distinguished` | supported | FCL.815(a)(1)-(3), (d) | (d) sets one recency condition for the rating, whichever of wheels or skis it covers; only the scope of the privilege differs, and the reading says it is not distinguished. | None. |
| `two-years-rolling` | supported | FCL.815(d) | "During the last 2 years" read up to the date evaluated is the plain meaning. | None. |

## easa.privilege.sfcl-banner-towing (`credentials/easa/privileges/sfcl-banner-towing.yaml`)

| Interpretation | Verdict | Paragraph(s) | Reasoning | Action |
| --- | --- | --- | --- | --- |
| `tows-counted` | needs-decision | SFCL.205(f) | (f) asks for "a minimum of five tows" for "the sailplane towing or banner towing rating". Options: (1) any tow in any aircraft except ultralights (the file's choice); (2) banner tows only, as the Part-FCL file reads FCL.805(e); (3) option (1) or (2) restricted to TMGs, the aircraft SFCL.205(a) concerns. | None; for decision. |
| `sfcl-licences-only` | supported | SFCL.205(a); FCL.805(a) | SFCL.205(a) concerns SPL holders with TMG privileges. | None. |

## easa.privilege.sfcl-launch-methods (`credentials/easa/privileges/sfcl-launch-methods.yaml`)

| Interpretation | Verdict | Paragraph(s) | Reasoning | Action |
| --- | --- | --- | --- | --- |
| `recorded-training` | supported | SFCL.155(a), (b) | The training of (a) is completed once and entered and signed under (b); a recorded privilege is the record of that entry. Not recounting training launches is stated. | None. |
| `no-expiry` | supported | SFCL.155(a) | (a) sets no expiry; recency is (c). Honouring a recorded end date is conservative. | None. |
| `recency-with-spl` | supported | SFCL.155(c) | `easa.licence.spl` evaluates 5 launches per method in the last 24 months and 2 for bungee, with self-launch counting TMG take-offs, as (c) says. | None. |
| `sfcl-licences-only` | supported | SFCL.155(a) | (a) concerns SPL holders. | None. |

## easa.privilege.sfcl-sailplane-towing (`credentials/easa/privileges/sfcl-sailplane-towing.yaml`)

| Interpretation | Verdict | Paragraph(s) | Reasoning | Action |
| --- | --- | --- | --- | --- |
| `tows-counted` | needs-decision | SFCL.205(f) | As for SFCL banner towing: options (1) any tow except in ultralights (the file's choice), (2) sailplane tows only, (3) either restricted to TMGs. | None; for decision. |
| `sfcl-licences-only` | supported | SFCL.205(a) | As for SFCL banner towing. | None. |
| `resumption-with-instructor` | supported | SFCL.205(g) | (g) sets the make-up with or under an instructor for the sailplane towing rating (it does not mention the banner rating, whose file rightly has no such remedy). | None. |

## easa.privilege.spl-tmg (`credentials/easa/privileges/spl-tmg.yaml`)

| Interpretation | Verdict | Paragraph(s) | Reasoning | Action |
| --- | --- | --- | --- | --- |
| `twelve-hours-on-sailplanes` | needs-decision | SFCL.160(b)(1) | (b)(1) names TMGs for the included items but no aircraft for the 12 hours. Options: (1) sailplanes including TMGs (the file's choice; Part-SFCL writes "sailplanes, excluding TMGs" where it excludes them, and the dual time is under an FI(S)); (2) flight time in any aircraft. Option (2) would credit aeroplane time. | None; for decision. |
| `training-flight` | supported | SFCL.160(b)(1)(iii), (b)(2) | (iii) asks for one training flight of at least 1 hour total time with an instructor; a TMG flight of 60 minutes with dual time is that. The instructor's and examiner's qualification is stated as not checked. | None. |
| `part-fcl-exemption` | supported | SFCL.160(c) | (c) exempts holders of a Part-FCL licence including TMG privileges; a LAPL(A) with TMG privileges is such a licence, and the text does not require it to be current. Requiring the TMG class rating not to be expired is the conservative side. | None. |
| `tmg-passengers` | supported | SFCL.160(e)(2) | (e)(2) asks for 3 take-offs and landings as PIC in TMGs in 90 days and, for night passengers, one of them at night. Applying the night rows only with a TMG night privilege and reporting day passengers separately follows. | None. |
| `extension-for-every-spl` | supported | SFCL.150(a), (b); SFCL.130(a)(2)(v) | The limitation of (a) depends on the aircraft of the skill test, which is not recorded. SFCL.130(a)(2) counts supervised solo as flight instruction (15 hours including 10 dual and 2 supervised solo), so instruction = dual + supervised solo follows; the landing away of (v)(B) is stated as not checked. | None. |
| `extension-credit` | contradicted (fixed) | SFCL.150(c)(1), (c)(2) | (c)(1) credits a class rating for TMGs, (c)(2) TMG privileges together with the FCL.140.A recency. A LAPL(A) has no class ratings, so its TMG privileges fall under (c)(2); the file credited them under (c)(1) without the recency, while its reading claimed (c)(2) was not applied. | `LAPL_A` removed from the `part_fcl_credit` condition; reading rewritten; examples `extension-credited-by-ppl-tmg-rating` (current) and `extension-not-credited-by-lapl-a-tmg` (lapsed). Only `easa.privilege.spl-tmg#extension_training` is affected; it is a training result and never decides a composite. |

## Sources added

| File | Article | Why |
| --- | --- | --- |
| `sources/easa/fcl-900.md` | FCL.900 | FCL.900(b)(1), the specific instructor certificate valid for at most 1 year, reserved by FCL.940 (`shorter-validities-not-applied`). |
| `sources/easa/fcl-915.md` | FCL.915 | FCL.915(e)(2), the yearly UPRT refresher reserved by FCL.940 (`shorter-validities-not-applied`). |

Both are taken unchanged from the EUR-Lex consolidated text 02011R1178-20260430 via the
Publications Office Cellar, with the header of the existing `sources/easa/` files.
