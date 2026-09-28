# Review: German ultralight kind ratings

Reviewer: review agent "de-ul-ratings" (automated first pass, not a sign-off). Date: 2026-09-28.

Scope: every interpretation in the six ultralight kind ratings under `credentials/de/ratings/`:

- `ul-three-axis.yaml` (11)
- `ul-weight-shift.yaml` (7)
- `ul-gyroplane.yaml` (6)
- `ul-helicopter.yaml` (6)
- `ul-powered-paraglider.yaml` (5)
- `ul-sailplane.yaml` (5)

The licence, the passenger authorisation and towing privileges, the instructor rating and the
shared passenger evaluation were reviewed in `reviews/de-ul-licence-privileges-instructor.md`
and are not repeated here.

## Summary

| Verdict | Count |
| --- | --- |
| supported | 22 |
| needs-decision | 17 |
| contradicted | 1 (fixed) |
| **Total** | **40** |

One credential and one example changed (`de.rating.ul-three-axis#deemed-authorisation`, below).
No source was added: every paragraph cited is already stored under `sources/de/`
(LuftPersV §§ 44, 45, 45b, 84a, 122; LuftVZO § 1). No `approved_by` / `approved_on` field was
set.

## How to use this review

- **supported**: the reading follows from the cited paragraph as stored under `sources/de/`.
  A qualified reviewer can sign it off after checking that paragraph.
- **needs-decision**: the text leaves real room, the reading is a choice forced by what the
  record can hold, or the numbers come from an association rule. The table gives the options
  and the file's choice. If the owner picks another option, the credential, its examples and
  the reading change together.
- **contradicted**: the file's result differed from the text; fixed here, with a changelog
  fragment.
- Readings that rest on DULV or DAeC rules (weight-shift, gyroplane, powered paraglider, UL
  sailplane recency numbers) are at most needs-decision: no association text is stored or
  verified in this repository, only the delegation in § 45(4) second sentence (details set by
  the Beauftragte under § 31c LuftVG "entsprechend § 42 Absatz 2").
- Quotations and labels refer to LuftPersV "zuletzt geändert durch Art. 2 V v. 7.12.2021
  I 5190" and LuftVZO "zuletzt geändert durch Art. 28 V v. 11.12.2024 I Nr. 411", as stored.

## Changes made

| File | Change |
| --- | --- |
| `credentials/de/ratings/ul-three-axis.yaml` | `passengers.authorised` no longer accepts a PPL(A) or SPL valid on the date evaluated; only a recorded, valid `UL_PASSENGER_AUTH` on the same licence. Reading of `deemed-authorisation` rewritten (grant on issue, validity follows the UL licence, § 122 repealed); refs now § 84a(2), § 84a(5), § 122. |
| `examples/de.rating.ul-three-axis.yaml` | `pax-deemed-by-ppl`: outcome current becomes unknown (`pax.ul_authorisation_missing`), with the informational progress rows. No other example changed. |
| `fragments/changelog/review-de-ul-ratings.md` | Changelog line for the change. |
| `fragments/coverage/review-de-ul-ratings.yaml` | Replaces the `de:LuftPersV.84a` `not_evaluated` text, which said the deemed grant is evaluated by the three-axis passengers evaluation. |

## credentials/de/ratings/ul-three-axis.yaml (`de.rating.ul-three-axis`)

| Interpretation | Verdict | Paragraph(s) | Reasoning | Action |
| --- | --- | --- | --- | --- |
| three-axis-includes-motorgliders | supported | § 45(2) first sentence; § 44(2) | The kind entered is "aerodynamisch gesteuerte Ultraleichtflugzeuge", which names the control method; an ultralight motorglider is an aerodynamically controlled ultralight aeroplane, and neither § 45 nor § 84a names a separate motorglider kind. The kind list itself (what a Beauftragter enters) is not in the statute. | None. |
| scope-of-closing-words | needs-decision | § 45(2) second sentence | "In den 12 Stunden" makes PIC time, take-offs, landings and the training flight parts of the 12 hours, which already include TMG and SEP time (first sentence); the closing words "auf aerodynamisch gesteuerten Ultraleichtflugzeugen" can qualify only the training flight (file) or every element. § 45(2a) puts "auf einem Ultraleichthubschrauber" after the training flight too, but there the whole 6 hours are helicopter time, so it decides nothing. Options: (a) training flight only (file, more generous); (b) PIC hours, take-offs and landings also on three-axis ULs only, which lapses pilots whose PIC time or circuits are mostly SEP/TMG. | Owner decision. |
| counted-aircraft | supported | § 45(2) first sentence; § 45(3) first sentence | "Reisemotorsegler" is the TMG class; "einmotorige Landflugzeuge mit Kolbentriebwerk" excludes turbine, multi-engine and sea aeroplanes; other UL kinds are not named. No launch method is named, so every flight counts. | None. |
| takeoffs-and-landings | supported | § 45(2) second sentence | "12 Starts und 12 Landungen" are two separate counts; counting each as logged follows. | None. |
| training-flight-as-dual | needs-decision | § 45(2) second sentence; § 45b no. 3 | "Ein Übungsflug von mindestens einer Stunde Flugzeit in Begleitung eines Fluglehrers": one flight of at least 60 minutes is right. Requiring dual time is narrower than "in Begleitung": § 45b no. 3 counts "die Flugzeit als Luftfahrzeugführer bei vorgeschriebenen Übungsflügen mit Fluglehrer", so the ordinance expects such a flight may be logged as pilot (PIC) with the instructor aboard. Options: (a) a flight with dual time (file; a flight logged as PIC with an instructor aboard does not count, conservative); (b) also a flight flagged as flown with an instructor, which needs a record field or flag the vocabulary does not have. | Owner decision; mention § 45b no. 3 in the reading. |
| proficiency-check-period | needs-decision | § 45(3) first and third sentences | § 45(3) names no period: the check replaces "die Voraussetzungen nach Absatz 2", which are 24-month requirements, so counting a check for 24 months (file) is the natural reading; a shorter effect (only until the experience can be counted again) is possible. The aircraft list (three-axis UL, TMG, SEP(land)) matches the first sentence; the examiner's recognition and the logbook signature (third sentence) are not checked, as stated. | Owner decision; file's choice is the better-supported one. |
| licence-without-expiry | supported | § 45(1) first sentence | "Wird unbefristet erteilt"; ignoring a recorded rating expiry and letting recency decide follows. | None. |
| flight-time-credit | supported | § 45b no. 1-5 | § 45b counts instructor, student, training-flight, examiner and applicant time as flight time "für ... den Nachweis für die Ausübung der Rechte"; ordinary PIC time is flight time without needing the list. Instructor time counts under no. 1 only during training and prescribed training flights, which is what instruction given practically always is. | None. |
| kindless-flights-unknown | supported | § 45(2); policy:unknown-input | A UL flight with no kind cannot be shown to be on "aerodynamisch gesteuerten Ultraleichtflugzeugen". | None. |
| deemed-authorisation | contradicted (fixed) | § 84a(2) second and third sentences; § 84a(5); § 122 | The authorisation "gilt mit der Erteilung des Luftfahrerscheins ... als erteilt" for three-axis pilots who hold a valid PPL or SPL: the test is at the day of issue of the UL licence, and once granted the authorisation's validity follows the UL licence (§ 84a(5) second sentence), not the PPL/SPL. The file tested a PPL(A)/SPL valid on the date evaluated and said so ("not on the day the ultralight licence was issued"). That wrongly authorises a pilot whose PPL was obtained after the UL licence, and denies one whose PPL has since lapsed. The engine's `holds` cannot compare the PPL/SPL validity with the UL licence's issue date. § 122, which the third sentence keeps unaffected, is repealed (`sources/de/luftpersv-122.md`). | Fixed: the PPL/SPL alternative removed; a deemed authorisation counts when recorded as `UL_PASSENGER_AUTH` on the UL licence (§ 84a(5) first sentence enters it there), otherwise unknown asking for it. Example `pax-deemed-by-ppl` now unknown. Future option (needs a vocabulary word): honour a PPL(A)/SPL whose validity covers the UL licence's issue date when both dates are recorded. |
| medical-above-120-kg-not-evaluated | needs-decision | § 45(1) second sentence | The reading states the text correctly and the gap is real (no empty mass per ultralight). Three-axis ULs are practically always above 120 kg, so the rating (and its `licence` composite) can report current for a pilot whose licence is not valid for want of a medical. Options, as in the earlier review for `de.licence.ul#ultralight-licences`: (a) leave unevaluated (file); (b) require a class 2 or LAPL medical for this kind, which changes composite results. | Owner decision. |

## credentials/de/ratings/ul-weight-shift.yaml (`de.rating.ul-weight-shift`)

| Interpretation | Verdict | Paragraph(s) | Reasoning | Action |
| --- | --- | --- | --- | --- |
| association-rules-by-authority | needs-decision | § 45(4); assoc:dulv:ul-recency-weight-shift; assoc:daec:ul-recency-weight-shift | The statute requires "ausreichende fliegerische Übung" and delegates the details (§ 45(4) second sentence), which supports a rule per Beauftragter. The numbers and the mapping (LBA-issued licences follow the DULV rule) are association matters not verified here. Options: (a) rule by issuing body, LBA with DULV (file); (b) one rule for all. | Owner decision (association text). |
| dulv-pic-only | needs-decision | § 45(4); assoc:dulv:ul-recency-weight-shift | 12 hours PIC in 24 months with no check alternative comes from the association rule, not the statute. | Owner decision (association text). |
| daec-safety-training-not-recorded | needs-decision | § 45(4); assoc:daec:ul-recency-weight-shift | 12 hours and 12 landings, plus an untracked safety training row, come from the association rule; the reading already flags that take-offs may also be named. | Owner decision (association text). |
| luftvzo-numbering | supported | § 45(4) first sentence; LuftVZO § 1(4) | The stored LuftVZO § 1(4) first sentence has no numbered items, so "Satz 1 Nummer 1" cannot be resolved; read without the number it names one- and two-seat Luftsportgeräte of at most 120 kg empty mass. The textual finding is exact; whether this kind is evaluated under § 45(4) when heavier is the separate `above-120-kg-not-named` decision. | None. |
| above-120-kg-not-named | needs-decision | § 45(1), (2), (2a), (4) | Correct that no paragraph names weight-shift ULs above 120 kg (most trikes): § 45(2) is three-axis only, (2a) helicopters, (4) the 120 kg class. Options: (a) apply the association rule regardless of mass (file); (b) report no statutory recency above 120 kg (not determinable without the mass). | Owner decision. |
| licence-without-expiry | supported | § 45(1) first sentence | As for three-axis. | None. |
| kindless-flights-unknown | supported | § 44(2) first sentence; policy:unknown-input | The licence covers "Luftsportgerät der im Luftfahrerschein eingetragenen Art"; a kindless flight cannot be attributed. | None. |

## credentials/de/ratings/ul-gyroplane.yaml (`de.rating.ul-gyroplane`)

| Interpretation | Verdict | Paragraph(s) | Reasoning | Action |
| --- | --- | --- | --- | --- |
| association-rule-unverified | needs-decision | § 45(4); assoc:dulv:ul-recency-gyroplane | Numbers from the DULV rule, cited only; application to DAeC-issued licences unverified, as stated. Options: (a) DULV numbers for every issuer (file); (b) per issuer once the DAeC rule is known. | Owner decision (association text). |
| gyroplane-class-credited | needs-decision | § 45(4); assoc:dulv:ul-recency-gyroplane | Crediting Part-FCL gyroplane time and checks, and not counting take-offs, depends on the association rule; the statute says nothing. | Owner decision (association text). |
| luftvzo-numbering | supported | § 45(4); LuftVZO § 1(4) | As for weight-shift. | None. |
| above-120-kg-not-named | needs-decision | § 45(1), (2), (2a), (4) | UL gyroplanes are practically always above 120 kg, so § 45(4) read literally does not cover them. Options as for weight-shift. | Owner decision. |
| licence-without-expiry | supported | § 45(1) first sentence | As for three-axis. | None. |
| kindless-flights-unknown | supported | § 44(2); policy:unknown-input | As for weight-shift. | None. |

## credentials/de/ratings/ul-helicopter.yaml (`de.rating.ul-helicopter`)

| Interpretation | Verdict | Paragraph(s) | Reasoning | Action |
| --- | --- | --- | --- | --- |
| only-ultralight-helicopters | supported | § 45(2a); § 45(3) second sentence | Hours "auf Ultraleichthubschraubern", the check "auf einem Ultraleichthubschrauber"; unlike § 45(2), no other class and no PIC time is named. | None. |
| takeoffs-and-landings | supported | § 45(2a) second sentence | "Sechs Starts und sechs Landungen", counted separately. | None. |
| training-flight-as-dual | needs-decision | § 45(2a) second sentence; § 45b no. 3 | Same point as three-axis: "in Begleitung eines Fluglehrers" may be logged as pilot with the instructor aboard (§ 45b no. 3). Options as there. | Owner decision. |
| proficiency-check-period | needs-decision | § 45(3) second sentence | No period named; counting a check for the 12 months that (2a) uses (file) is the natural reading. | Owner decision; file's choice better supported. |
| licence-without-expiry | supported | § 45(1) first sentence | As for three-axis. The § 45(1) second-sentence medical (UL helicopters are above 120 kg) is not evaluated, as for the licence. | None. |
| kindless-flights-unknown | supported | § 45(2a); policy:unknown-input | As for three-axis. | None. |

## credentials/de/ratings/ul-powered-paraglider.yaml (`de.rating.ul-powered-paraglider`)

| Interpretation | Verdict | Paragraph(s) | Reasoning | Action |
| --- | --- | --- | --- | --- |
| association-rule-unverified | needs-decision | § 45(4); assoc:dulv:ul-recency-powered-paraglider; assoc:daec:ul-recency-powered-paraglider | 30 landings in 24 months comes from the association rules; the statute sets no number. | Owner decision (association text). |
| luftvzo-numbering | supported | § 45(4); LuftVZO § 1(4) | As for weight-shift; most powered paragliders are in the 120 kg class, so the literal reading fits this kind. | None. |
| above-120-kg-not-named | needs-decision | § 45(1), (2), (2a), (4) | Options as for weight-shift; of little practical weight for this kind. | Owner decision. |
| licence-without-expiry | supported | § 45(1) first sentence | As for three-axis. | None. |
| kindless-flights-unknown | supported | § 44(2); policy:unknown-input | As for weight-shift. | None. |

## credentials/de/ratings/ul-sailplane.yaml (`de.rating.ul-sailplane`)

| Interpretation | Verdict | Paragraph(s) | Reasoning | Action |
| --- | --- | --- | --- | --- |
| association-rule-unverified | needs-decision | § 45(4); assoc:daec:ul-recency-sailplane | 5 landings in 12 months, any launch method, comes from the DAeC rule and is applied to every issuer; unverified. | Owner decision (association text). |
| luftvzo-numbering | supported | § 45(4); LuftVZO § 1(4) | As for weight-shift. | None. |
| above-120-kg-not-named | needs-decision | § 45(1), (2), (2a), (4) | Options as for weight-shift. | Owner decision. |
| licence-without-expiry | supported | § 45(1) first sentence | As for three-axis. | None. |
| kindless-flights-unknown | supported | § 44(2); policy:unknown-input | As for weight-shift. | None. |

## Checks

`go test ./...` and `go run ./cmd/rulescheck -strict -fragments` pass with the change
(fragments listed: one changelog, one coverage).
