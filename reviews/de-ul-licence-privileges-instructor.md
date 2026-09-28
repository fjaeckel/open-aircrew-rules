# Review: German ultralight licence, privileges, instructor rating and shared passenger evaluation

Reviewer: review agent "de" (automated first pass, not a sign-off). Date: 2026-09-28.

Scope: every interpretation in

- `credentials/de/licences/ul.yaml` (6)
- `credentials/de/instructors/ul-instructor.yaml` (7)
- `credentials/de/privileges/ul-passenger-authorisation.yaml` (2) and `ul-towing.yaml` (4)
- `credentials/de/shared/ul-passengers.yaml` (6)

The six ultralight kind ratings (`credentials/de/ratings/`) are a separate, later review; where
an interpretation here depends on one of them, the dependency is noted at the end.

## Summary

| Verdict | Count |
| --- | --- |
| supported | 15 |
| needs-decision | 10 |
| contradicted | 0 |
| **Total** | **25** |

No credential file or example was changed. Five verbatim sources were added (below). No
`approved_by` / `approved_on` field was set.

## How to use this review

- **supported**: the reading follows from the cited paragraph as stored under `sources/de/`
  (the paragraph is named in the table). A qualified reviewer can sign it off after checking
  that paragraph; nothing else is needed.
- **needs-decision**: the text leaves real room, or the reading is a choice forced by what the
  record can hold. The table gives the options and which one the file takes. The owner (or a
  qualified reviewer) picks one; if it is not the file's choice, the credential, its examples
  and the reading change together. Where a choice would need a word the vocabulary does not
  have, this is said.
- **contradicted**: none found.
- Readings that rest on DULV or DAeC rules can be at most needs-decision here, because no
  association text is stored or verified in this repository. None of the 25 interpretations
  in scope relies on an association document for its result; the passenger authorisation's
  limited validity (§ 84a(5)) is the only place the Beauftragte's own rules come in, and there
  the file only honours a recorded date.
- All quotations and paragraph labels refer to LuftPersV "Stand: zuletzt geändert durch
  Art. 2 V v. 7.12.2021 I 5190", as stored.

## Sources added

All from the official XML at www.gesetze-im-internet.de, retrieved 2026-09-28, with the header
of the existing `sources/de/*.md` files (origin `de-amtliches-werk`).

| File | Why |
| --- | --- |
| `sources/de/luftpersv-95a.md` | The instructor rating the § 96 extension applies to (issue requirements, the instructor test of § 95a(3)); needed to read "Befähigungsprüfung" and "für die Berechtigung nach § 95a" in § 96(4). |
| `sources/de/luftpersv-122.md` | § 84a(2) third sentence keeps "§ 122 Abs. 1" unaffected; the stored text shows § 122 is repealed ("weggefallen"), so that sentence has no effect today. |
| `sources/de/bgb-187.md`, `sources/de/bgb-188.md` | Start and end of periods counted in months or years. |
| `sources/de/vwvfg-31.md` | Applies §§ 187 to 193 BGB to periods in administrative procedure; evidence for the one-day question on the instructor rating's three years. |

## credentials/de/licences/ul.yaml (`de.licence.ul`)

| Interpretation | Verdict | Paragraph(s) | Reasoning | Action |
| --- | --- | --- | --- | --- |
| kind-decides-the-rule | supported | LuftPersV § 44(2); § 45(2), (2a), (4) | § 44(2) grants the privileges for the kind of Luftsportgerät entered in the licence; § 45(2) and (2a) set recency per entered kind and § 45(4) leaves the rest to the Beauftragte per kind. A kindless rating cannot be matched to any of these, so unknown with a request for the kind follows. | None. |
| ultralight-licences | needs-decision | § 45(1); § 44(1); § 31c LuftVG (via § 44(1)) | "Unbefristet erteilt" (§ 45(1) first sentence) supports no time-limit validity. But § 45(1) second sentence makes the licence valid only with a valid medical certificate (MED.A.030(b)) for Luftsportgeräte above 120 kg empty mass; the reading's "so no licence-level validity is evaluated" does not say so. `coverage/articles.yaml` (de:LuftPersV.45) and `ul-three-axis#medical-above-120-kg-not-evaluated` do explain why it is left out (no empty mass per ultralight in the record). Options: (a) keep it unevaluated and name the medical gap in this reading too (file's choice, minus the wording); (b) add `requires_any` on the EASA class 2 / LAPL medicals for the kinds that are practically always above 120 kg (three-axis, helicopter, gyroplane), which would change composite results of existing examples. | decided 2026-09-28: option 2, applied (LAPL, class 2 or class 1 medical required by the three-axis, helicopter, gyroplane and weight-shift ratings; powered paraglider and sailplane need none, see fragments/vocab-requests/decisions-de.md). |
| training-hours-partial | supported | § 42(4) no. 1 and 2; § 42(5) no. 1 and 2 | The reading describes the credits and the no. 2 exercises accurately (three-axis: up to 20 h PIC sailplane/helicopter or 5 h weight-shift; weight-shift: up to 10 h PIC of the listed aircraft, summarised as "other aircraft"). Not counting the credits under-counts, so a programme completed with credit shows "in progress": conservative, and stated. § 42(4) no. 3 (training at an approved organisation for PPL/TMG holders) is not mentioned; it is not an hours requirement. | None (optional: mention § 42(4) no. 3). |
| training-time-counted | supported | § 45b no. 1-5; § 42(4) no. 1 | § 45b counts time as student with instructor (no. 2) and as pilot on prescribed training flights with instructor (no. 3) for acquiring the licence; "Alleinflug" as supervised solo or PIC time follows. Counting later PIC time too is a stated record limitation that only matters for someone who already holds the licence. | None. |
| training-for-ul-licence-records | supported | § 42(1) | § 42(1) lists the requirements for acquiring the licence; evaluating the programme per recorded licence without deciding student status is a scope statement, and it is accurate. | None. |
| kindless-flights-unknown | supported | § 44(2); policy:unknown-input | A flight must be of the kind trained (§ 42(4) "mit aerodynamisch gesteuerten Ultraleichtflugzeugen", § 42(5) likewise); a kindless flight cannot be attributed, and treating it as unknown follows the catalogue-wide policy. | None. |

## credentials/de/instructors/ul-instructor.yaml (`de.instructor.ul-instructor`)

| Interpretation | Verdict | Paragraph(s) | Reasoning | Action |
| --- | --- | --- | --- | --- |
| three-years-from-valid-from | needs-decision | § 96(1) first sentence; § 31(1) VwVfG; §§ 187(2), 188(2) BGB | "Mit einer Gültigkeitsdauer von drei Jahren erteilt" supports three years from the start of validity, and an earlier recorded expiry winning is sound. Open point: the catalogue computes the end as the same calendar day three years on (valid from 10 Jan 2024, valid through 10 Jan 2027, example `expired-day-after`). If the BGB period rules apply (via § 31(1) VwVfG, as for administrative periods), a period beginning at the start of 10 Jan 2024 ends with the day before the corresponding day, 9 Jan 2027 (§ 188(2) second case with § 187(2)). Whether § 31 VwVfG governs the validity period of a licence entry (rather than a procedural period) is itself arguable, and the licence normally shows the date. Options: (a) keep the catalogue-wide convention (file's choice); (b) end one day earlier for DE-computed periods, which needs an engine/compiler option and changes the `expired-day-after` and `window-day-before` dates. | decided 2026-09-28: option 2, applied (authority convention DE, periods end the day before; reading cites §§ 186, 187(2), 188(2) BGB). |
| last-three-years-before-expiry | needs-decision | § 96(4) opening words | "Innerhalb der letzten drei Jahre" has no anchor in the text; it is the three years before the extension (application or grant). Options: (a) the 36 months before the expiry date (file's choice; the latest date the extension can take effect, which lets a pilot see the status before applying); (b) the three years before the date evaluated or the application date. With a three-year validity, (a) is practically the validity period. | decided 2026-09-28: option 1, applied (kept; P3). |
| instruction-as-instructor-or-examiner | needs-decision | § 96(4) no. 1; § 96(2); § 95a | Three open points. (1) "60 Starts und Landungen": read as 60 of each (60 circuits); § 45(2) and § 45a repeat the number when they mean each ("12 Starts und 12 Landungen"), § 96(4) does not, so 60 take-offs and landings in total (30 each) is a possible but weaker reading. (2) "Als Lehrer oder Prüfer für die Berechtigung nach § 95a" read as instruction or examining on ultralights in any kind; the stored § 95a shows the rating is acquired per kind (§ 95a(1) no. 1) and § 96(2) allows limiting it, so counting only the kinds the rating covers is the stricter option (the record does not hold that scope). (3) Examiner time counts in both alternatives (file) versus instructor time only for take-offs and landings; the text puts "als Lehrer oder Prüfer" after both alternatives, which supports the file. | decided 2026-09-28: split option, applied (60 take-offs and 60 landings as separate rows; only the kinds the rating's detail names, unknown without one). |
| refresher-course-event | supported | § 96(4) no. 2 | The text reads "Fortbildungslehrgang für Fluglehrer innerhalb der Gültigkeitsdauer der Lehrberechtigung"; the 12 months before renewal apply to renewal only ("vor der Erneuerung"), which is not evaluated. The body that ran or recognised the course is not recorded; saying so is required by the `with:` qualifier and is done. | None. |
| befaehigungspruefung-event | needs-decision | § 96(4) no. 3; § 95a(3); § 45(3); § 45b no. 5 | The ordinance uses "Befähigungsüberprüfung" for the pilot's proficiency check (§ 45(3), § 45b no. 5) and "Befähigungsprüfung" only in § 96(4) no. 3, in the extension of an instructor rating whose issue requires a test of instructor competence (§ 95a(3)). That supports the file's reading as an instructor assessment of competence, excluding pilot proficiency checks. The other reading (any passed proficiency check, the difference in wording being accidental) cannot be ruled out from the text. "Erfolgreiches Ablegen" (passed) taken from the event being recorded, and the examiner not checked, are stated. | decided 2026-09-28: option 1, applied (kept; P2). |
| renewal-not-evaluated | supported | § 96(4); § 96(1); § 88 | § 96(4) covers extension and renewal; only extension is evaluated, and after expiry the rating is reported expired, which § 96(1) supports. § 88 (flight-engineer instructor) is outside ultralight scope. | None. |
| towing-instruction-not-evaluated | supported | § 96(3) | "Ausreichende praktische Erfahrung im Schleppflug" is not quantified; not evaluating the towing-instruction privilege is correct. | None. |

## credentials/de/privileges/ul-passenger-authorisation.yaml (`de.privilege.ul-passenger-authorisation`)

| Interpretation | Verdict | Paragraph(s) | Reasoning | Action |
| --- | --- | --- | --- | --- |
| validity-follows-licence | supported | § 84a(5) second sentence; § 45(1) | Validity follows the ultralight licence unless the Beauftragte limits it; a recorded expiry date is honoured and the extension conditions (Beauftragte's rules) are not evaluated, as stated. Caveat shared with `de.licence.ul#ultralight-licences`: the licence's own validity also depends on a medical above 120 kg (§ 45(1) second sentence), which is not evaluated anywhere. | None. |
| acquisition-taken-as-met | supported | § 84a(2), (3), (4), (5) first sentence | A recorded authorisation is a granted one; the acquisition requirements are not conditions of its validity. The kind it is entered for (§ 84a(5) first sentence) is not compared here; where that matters (carrying passengers) it is a decision in the shared file (`authorisation-recorded`). Note: the record can hold the kind in the privilege's `detail` (examples use `detail: THREE_AXIS`), so "not recorded" is better put as "not required or compared". | None (wording note only). |

## credentials/de/privileges/ul-towing.yaml (`de.privilege.ul-towing`)

| Interpretation | Verdict | Paragraph(s) | Reasoning | Action |
| --- | --- | --- | --- | --- |
| tow-kind-not-recorded | needs-decision | § 84(5) first sentence; § 84(4) | "In der jeweils eingetragen Art" is the kind entered under § 84(4) (take-up method and object towed), as the reading says. Two choices follow. (1) Every ultralight tow counts, whatever was towed: the record has no field for the kind entered with the rating, though flights may carry `towKind` and the vocabulary has a `tow_kinds` filter; option: record the entered kind (for example in `detail`) and filter by it when present. (2) Tows flown in any other class (for example a glider tow in an SEP aeroplane) do not count: § 84(5) counts tow flights "in der jeweils eingetragen Art" and does not say in which aircraft, so counting only ultralight tows is stricter than the text. File's choice: (1) all UL tows, (2) UL only. | decided 2026-09-28: option 2 (first limb), applied (tow kind and take-up from the rating's detail; unknown without them; ultralight tows only). |
| tows-counted | needs-decision | § 84(5) first sentence | The text counts "Schleppflüge" (tow flights). Counting `towedGliders` per flight is right when a logbook line aggregates several tow flights (the usual case: every glider tow is a flight of its own), and over-counts when one flight tows several objects at once (a double tow is one Schleppflug). The reading's own justification ("the text counts tow flights, and a flight recording several tows is read as several") states the tension without resolving it. Options: (a) count towed objects (file); (b) count flights flagged as tow flights, or landings on them, which matches the text for single-flight lines. | decided 2026-09-28: option 2, applied (landings on tow-flagged flights; a double tow counts once). |
| statutory-for-every-authority | supported | § 84(5); § 84(1) | § 84(5) is statutory and does not depend on which Beauftragter issued the licence; honouring a recorded expiry date is the catalogue's validity convention. | None. |
| supervised-tows-not-restoring | supported | § 84(5) second sentence; § 84(2) no. 2; § 84(3) no. 2 | When the ten tows are missing, § 84(5) sends the holder back to the five supervised tow flights of § 84(2) no. 2 or § 84(3) no. 2; not modelling them as a restoring event is a stated limitation, and `on_fail` names them. | None. |

## credentials/de/shared/ul-passengers.yaml (`de.shared.ul-passengers`)

| Interpretation | Verdict | Paragraph(s) | Reasoning | Action |
| --- | --- | --- | --- | --- |
| same-kind-any-role | supported | § 45a first sentence | The text requires three take-offs and three landings "mit einem Luftsportgerät derselben Art" in the preceding 90 days and names no role; take-offs and landings are counted separately and the smaller decides. The 90-day window (asOf minus 90 days through asOf) matches "innerhalb der vorhergehenden 90 Tage". Kindless flights as unknown follows policy. | None. |
| authorisation-recorded | needs-decision | § 84a(1); § 84a(5) first sentence; § 84a(2) second sentence | Requiring a recorded, unexpired authorisation on the same licence is supported. Open: the authorisation is entered for the kind the pilot trained on (§ 84a(5) first sentence), but any authorisation on the licence counts for every kind. Options: (a) any kind (file); (b) compare the privilege's `detail` with the kind flown when recorded, unknown otherwise. The deemed authorisation for PPL/SPL holders is passed in by the three-axis rating through `authorised` (see "Notes for the ratings review"). | decided 2026-09-28: option 2, applied (detail compared with the kind flown; unknown without a kind). |
| authorisation-progress-informational | needs-decision | § 84a(2) first sentence; § 84a(3) | The rows (five cross-country flights with an instructor since licence issue, two with an intermediate landing and at least 200 km) match § 84a(2) for two-seat ultralight aeroplanes, and they are informational only. The reading's "for other kinds § 84a(3) leaves the requirements to the Beauftragte" is too broad: § 84a(3) applies § 42(2) to two-seat hang gliders, paragliders "oder anderen vergleichbaren Luftsportgeräten" and tandem parachutes; for ultralight helicopters and gyroplanes neither § 84a(2) (if "Ultraleichtflugzeuge" is read narrowly) nor § 84a(3) clearly applies. Options: (a) show the § 84a(2) rows for every kind as a guide (file); (b) show them only for the kinds § 84a(2) covers and a `not_recorded` row for the others. Never changes a status either way. | decided 2026-09-28: option 1, applied (kept; § 84a(3) sentence corrected; P5). |
| xc-200-km-per-flight | needs-decision | § 84a(2) first sentence; § 42(4) no. 2; § 42(5) no. 2; § 42(5a) no. 2 | The file requires 200 km on each of the two landing flights and calls this the literal reading. The text does not bear out "literal": it reads "zwei Überlandflüge mit Zwischenlandung über eine Gesamtstrecke von mindestens 200 Kilometer", while the same ordinance writes "über jeweils eine Gesamtstrecke" in § 42(4) no. 2 and § 42(5) no. 2 when it means each flight, and "einen Überlandflug ... über eine Gesamtstrecke" in § 42(5a) no. 2 for a single flight. The missing "jeweils" and the singular "eine Gesamtstrecke" favour 200 km for both flights together; the per-flight reading remains possible (drafting by analogy to § 42(4)). Options: (a) per flight (file, stricter); (b) combined over the two flights, which the vocabulary cannot express today (`min_distance_km` filters single flights; a sum over the qualifying flights would need a new word). Informational only: no status changes either way. | decided 2026-09-28: option 1, applied (kept per flight; reading reworded; P5). |
| day-only | supported | § 44(2) first sentence; § 45a second sentence | The licence grants the privileges "am Tage" (night only for parachutes); no night passenger recency follows. Tandem parachuting (§ 45a second sentence, ten jumps) is out of scope. | None. |
| practical-test-not-checked | supported | § 84a(4) | The practical test is an acquisition requirement; a recorded authorisation implies it was passed. | None. |

## Notes for the ratings review (not in scope here, no change made)

- **Deemed authorisation (`de.rating.ul-three-axis#deemed-authorisation`).** § 84a(2) second
  sentence deems the authorisation granted "mit der Erteilung des Luftfahrerscheins" to
  three-axis pilots holding a valid PPL or SPL licence. The text ties the deemed grant to the
  day the ultralight licence is issued; once granted, its validity follows the ultralight
  licence (§ 84a(5)). The rating file instead tests a PPL(A) or SPL valid on the date
  evaluated. That differs from the text both ways: a pilot who held the PPL when the ultralight
  licence was issued keeps the authorisation even if the PPL is later not valid (the file says
  not authorised), and a pilot who obtained the PPL after the ultralight licence gets no deemed
  grant (the file says authorised). The record may not show the dates needed (licence issue
  dates are optional), which is the likely reason for the choice. This looks like a
  candidate for "contradicted" in the ratings review, with options: compare the PPL/SPL issue
  date with the ultralight licence issue date when both are recorded, and fall back to the
  current reading (or unknown) otherwise. The third sentence ("§ 122 Abs. 1 bleibt
  unberührt") refers to a repealed paragraph (`sources/de/luftpersv-122.md`) and has no effect.
- **Medical above 120 kg** (§ 45(1) second sentence) is noted in
  `de.rating.ul-three-axis#medical-above-120-kg-not-evaluated`; the licence and passenger
  authorisation readings above rely on the same gap.
- **Towing message key.** `de.privilege.ul-towing#recency` uses the description key
  `dulv_ul_towing`, while the rule is statutory for every issuing body
  (`statutory-for-every-authority`). Only the key name is association-flavoured; worth a
  rename when message keys are next touched (shared file, not edited here).
