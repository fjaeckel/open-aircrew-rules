# Review: FAA shared evaluations and pilot certificates

Reviewed on 2026-09-28 against the verbatim texts in `sources/faa/` (eCFR Title 14, up to
date as of 2026-09-24). This review is not a sign-off: every interpretation keeps
`approved_by: null` and `approved_on: null` until a qualified reviewer approves it
(CONTRIBUTING.md, "Interpretations and sign-off").

Files: `credentials/faa/shared/` (flight-review, medical-duration-commercial,
medical-duration-private, passengers-day, passengers-glider, passengers-night,
passengers-tailwheel) and `credentials/faa/licences/` (atp-airplane, atp-rotorcraft,
commercial-airplane, commercial-rotorcraft, glider, powered-lift-airship, private-airplane,
private-rotorcraft, recreational, sport, student). The instructor, privilege, rating and
medical files were reviewed in `reviews/faa-instructors-privileges-ratings-medicals.md`.

## Summary

| Verdict | Count |
| --- | --- |
| supported | 44 |
| needs-decision | 8 |
| contradicted (fixed) | 3 |
| **Total** | **55** |

By group: shared 26 (19 supported, 5 needs-decision, 2 contradicted), licences 29
(25 supported, 3 needs-decision, 1 contradicted).

Contradicted and fixed:

- `faa.shared.flight-review` / `student-exemption`: the flight review is evaluated once per
  pilot, on the first FAA certificate in the record. When that was a student certificate,
  the `student` stage returned `not_applicable` for the pilot even if the record also held a
  private or higher certificate, so a lapsed review was never reported (verified with the
  evaluator: student certificate first, private second, review of 2022, result
  `not_applicable`). 61.56(g) relieves a student pilot in training, not the holder of
  another pilot certificate acting on its privileges (61.56(c)). The stage now also requires
  that no sport, recreational, private, commercial, airline transport or glider certificate
  is held; example `faa.licence.student` / `flight-review-student-certificate-listed-first`
  added. The shared file feeds the eleven licence credentials of this review only; no
  other worked example, in or outside this review, changes its result.
- `faa.shared.medical-duration-commercial` / `commercial-glider-balloon`: the reading said
  commercial balloon privileges need no medical. 61.23(a)(2)(iii) requires at least a
  second-class certificate for commercial balloon privileges for compensation or hire; only
  balloon flight training ((b)(5)) and glider privileges ((b)(3)) need none. Reading
  corrected; the file's behaviour (the evaluation is reported for every commercial and ATP
  holder) was already the safe one, so no result changes.
- `faa.licence.glider` / `powered-classes-elsewhere`: the reading gave "a motorglider logged
  as SEP land" as an example of an aircraft that "is not a glider". Under 14 CFR 1.1 a glider
  is an aircraft whose free flight does not depend principally on an engine, so a
  self-launching glider is a glider. Reading corrected to say that such flights count here
  only when logged in the glider class; no result changes.

Needs a decision by the signing reviewer (the file's current choice is named first):

- `faa.shared.flight-review` / `who-and-endorsement-not-checked`: (a) every recorded
  proficiency check or practical test counts, including a check logged for another
  authority's licence (an EASA proficiency check sets the same flag); (b) count only checks
  for an FAA certificate, rating or operating privilege given by an FAA examiner (61.1
  "Examiner"), which needs a way to tie an event to an authority.
- `faa.shared.flight-review` / `simulator-counts`: (a) a review in a device counts, the
  part 142 conditions taken as met; (b) exclude device sessions, as the passenger
  evaluations do for the same part 142 condition.
- `faa.shared.passengers-day` / `class-or-type` and `faa.shared.passengers-night` /
  `class-or-type`: (a) count by the record's class, so SEP land and SET land flights do not
  count for each other; (b) pool SEP with SET (land and sea) for FAA ratings, since 61.5(b)(2)
  knows only single-engine and multiengine land and sea classes (needs a vocabulary
  request: `pooled_with_held` pools only classes the holder rates).
- `faa.shared.passengers-night` / `logged-night-takeoffs`: (a) logged night takeoffs (the
  1.1 night, from the end of evening civil twilight) count, which can over-count takeoffs
  between the end of civil twilight and 1 hour after sunset; (b) count the full-stop night
  landings for both, as the tailwheel evaluation does; (c) record takeoffs in the 61.57(b)
  period (vocabulary and record change).
- `faa.licence.powered-lift-airship` / `medical-grade-not-recorded`: (a) one file, at least
  a second-class certificate for every grade; (b) split by certificate grade, as the
  rotorcraft files are (private: any class or BasicMed; commercial: first or second; ATP:
  first). The record does name the grade through the licence kind, except on a certificate
  recorded as a glider pilot certificate.
- `faa.licence.sport` / `sport-classes`: (a) powered parachutes and weight-shift-control
  aircraft share the class OTHER, so takeoffs and landings in one count for the other,
  although 61.57(a)(1)(ii) requires the same category and class; (b) add classes for them
  (vocabulary request).
- `faa.licence.sport` / `medical-or-drivers-license`: (a) no medical requirement, so the
  composite never shows a missing medical or driver's license outside gliders and balloons;
  (b) require a medical certificate or BasicMed (a sport pilot flying on a driver's license
  alone reads unknown); (c) add a driver's-license credential and require it or a medical.

Sources added: `sources/faa/1.1.md` (general definitions: night, category, class, glider,
airplane, flight time), `sources/faa/61.1.md` (definitions: examiner, authorized
instructor), `sources/faa/61.5.md` (certificates and ratings issued) and
`sources/faa/61.157.md` (ratings for which the ATP practical test is given), all from the
eCFR versioner API as of 2026-09-24.

## How to use this review

Each table has one row per interpretation: its verdict, the paragraphs read, the reasoning
and what was changed. A signing reviewer can:

1. For `supported` rows, check the cited paragraph and approve the interpretation as it
   stands.
2. For `needs-decision` rows, pick an option (the file's current choice is named first); a
   different choice is a follow-up change, often needing a vocabulary request.
3. For `contradicted` rows, review the fix in the credential file and the examples named.
4. Read the observations at the end: they are not interpretations but gaps a reviewer may
   want to track.

Verdicts: `supported` means the cited paragraph says what the reading says (a documented
limitation counts as supported when it describes the text correctly and the file does what
it says); `needs-decision` means the text or the record leaves real room; `contradicted`
means the reading or the file's behaviour disagrees with the text.

## faa.shared.flight-review

Checked specifically: the window and the substitutes. 61.56(c) counts "since the beginning
of the 24th calendar month before the month" of the flight; `within_calendar_months: 24`
is the 24 calendar months before the month plus the month to date, the same. (d)(1) covers
proficiency checks and practical tests by an examiner, check airman or Armed Force for a
pilot certificate, rating or operating privilege; (d)(2) the flight instructor practical
tests; (e) FAA-sponsored proficiency program phases; all three are in the `review` node.
The subject is the first FAA licence in record order (engine `flight_review` subject).

| Interpretation | Verdict | Paragraph(s) | Reasoning | Action |
| --- | --- | --- | --- | --- |
| once-per-pilot | supported | 61.56(c) | The prohibition is addressed to the person ("no person may act as pilot in command ... unless ... that person has"), whatever certificate is used, so one result per pilot fits. The order dependency it creates is handled in `student-exemption`. | none |
| who-and-endorsement-not-checked | needs-decision | 61.56(c)(1), (c)(2), (d); 61.1(b) | Not checking the giver and the endorsement is a stated limitation. But a proficiency check logged for another authority (an EASA licence proficiency check sets the same `proficiencyCheck` flag) also counts, while 61.1 defines an examiner as a person authorized by the Administrator for a certificate or rating issued under part 61. Options above. | none |
| ipc-is-no-substitute | supported | 61.56(d)(1); 61.57(d)(3)(iv) | An IPC may be given by an authorized instructor, and a recorded IPC does not show who gave it; (d)(1) requires an examiner, check airman or Armed Force. Excluding it can only under-report (an IPC given by an examiner is not credited). | none |
| glider-flights-part-of-review | supported | 61.56(a), 61.56(b) | (b) lets the three instructional glider flights replace the 1 hour of flight training "required in paragraph (a)", i.e. within a review that still needs its ground training and endorsement. | none |
| student-exemption | contradicted | 61.56(g), 61.56(c) | Assuming a current solo endorsement is harmless for the student composite, which fails on the solo endorsement anyway. But with a student certificate first in the record the exemption also covered a private or higher certificate (see Summary). | Stage `student` now requires that no other FAA pilot certificate is held; reading reworded; example `flight-review-student-certificate-listed-first` added to `examples/faa.licence.student.yaml`. |
| simulator-counts | needs-decision | 61.56(i)(1), (i)(2), (i)(3) | (i) allows a simulator or FTD only within an approved part 142 course, in a device representing an aircraft the pilot is rated for, and without landing approval the 61.57 takeoffs and landings are also required. Taking these as met over-reports a review logged in a non-142 device, while the passenger evaluations exclude devices for the same unprovable condition. Options above. | none |
| instructor-relief-not-modelled | supported | 61.56(f) | (f) removes the hour of ground training only; the review and its date are unchanged. | none |

## faa.shared.medical-duration-commercial

Checked specifically: the table. 61.23(d)(1)(iii) and (d)(2)(i) give first- and
second-class certificates 12 months for commercial privileges at any age, to the end of
the last day of the 12th month after the month of examination; the file matches.

| Interpretation | Verdict | Paragraph(s) | Reasoning | Action |
| --- | --- | --- | --- | --- |
| commercial-glider-balloon | contradicted | 61.23(a)(2)(ii), (a)(2)(iii), (b)(3), (b)(5) | Glider commercial privileges need no medical, but commercial balloon privileges for compensation or hire need at least second class; only balloon flight training is exempt. The reading said balloon privileges in general need none. The behaviour (evaluate every commercial and ATP holder) is unchanged and safe. | Reading corrected; refs to (a)(2)(ii), (a)(2)(iii), (b)(3), (b)(5). |

## faa.shared.medical-duration-private

| Interpretation | Verdict | Paragraph(s) | Reasoning | Action |
| --- | --- | --- | --- | --- |
| no-medical-operations-not-distinguished | supported | 61.23(b), 61.23(c) | The record does not say which operation is flown; reporting the duration to every certificate holder is the conservative choice and is stated. | none |
| examination-date | supported | 61.23(d) | The table counts from "the month of the date of examination shown on the medical certificate" and takes the age "on the date of examination". | none |

## faa.shared.passengers-day

Checked specifically: sole manipulator and per class counting. 61.57(a)(1)(i) requires the
pilot to have been sole manipulator; 61.51(e)(1) shows that pilot-in-command time can be
logged without it, so a separate field is needed. (a)(1)(ii) counts in the same category,
class and type "if a class or type rating is required". The 90 days count both ends.

| Interpretation | Verdict | Paragraph(s) | Reasoning | Action |
| --- | --- | --- | --- | --- |
| class-or-type | needs-decision | 61.57(a)(1)(ii); 61.5(b)(2); 1.1 "Class" | Class versus type is right. But the record splits one FAA class in two: 61.5(b)(2) knows single-engine land and sea, while the record has SEP and SET; takeoffs in a single-engine turbine logged SET land do not count for an SEP land rating (verified with the evaluator). Options above. | none |
| sole-manipulator-recorded | supported | 61.57(a)(1)(i) | A missing field gives unknown input, never a silent pass. | none |
| no-simulator-credit | supported | 61.57(a)(3) | Device credit needs an approved-for-landings device used in a part 142 course; not counting devices only under-reports. | none |
| exceptions-not-modelled | supported | 61.57(a)(1), 61.57(e) | The part 121, 125 and 135 exceptions depend on employment; aircraft certificated for more than one pilot need the recency without passengers too, which the passenger evaluation does not show separately, as stated. | none |
| students-excluded | supported | 61.89(a)(1) | A student pilot may not act as pilot in command of an aircraft carrying a passenger. | none |

## faa.shared.passengers-glider

Checked specifically: the glider three-flight alternative belongs to 61.56(b) (the review)
and has no counterpart in 61.57; passenger recency needs three takeoffs and three landings
in a glider. A glider needs no class or type rating (61.5(b)(1)(iii) has no glider class).

| Interpretation | Verdict | Paragraph(s) | Reasoning | Action |
| --- | --- | --- | --- | --- |
| launch-is-takeoff | supported | 61.57(a)(1) | A glider's takeoff is its launch, whatever the method; landings alone do not meet "three takeoffs and three landings". | none |
| sole-manipulator-recorded | supported | 61.57(a)(1)(i) | As for powered aircraft. | none |
| no-type-no-simulator | supported | 61.57(a)(1)(ii), 61.57(a)(3) | No class or type rating is required for a glider, so the category counts; device credit is not provable. See observation 2 for self-launching gliders logged as TMG. | none |
| night-not-evaluated | supported | 61.57(b)(1) | 61.57(b) applies to gliders carrying persons at night, but launches are recorded as one count without a day or night split. The gap is stated; the result concerns passengers only and never decides a composite. | none |
| students-and-exceptions | supported | 61.89(a)(1), 61.57(e) | As for powered aircraft. | none |

## faa.shared.passengers-night

Checked specifically: the night. 61.57(b)(1) uses its own period, from 1 hour after sunset
to 1 hour before sunrise, for both the flight and the three takeoffs and full-stop
landings. The logbook's "day or night" (61.51(b)(3)(i)) follows the 1.1 definition, from the
end of evening civil twilight, which starts earlier.

| Interpretation | Verdict | Paragraph(s) | Reasoning | Action |
| --- | --- | --- | --- | --- |
| logged-night-takeoffs | needs-decision | 61.57(b)(1); 1.1 "Night"; 61.51(b)(3)(i) | The reading describes the difference correctly, but the choice over-counts: a takeoff logged at night 40 minutes after sunset counts although it is outside the 61.57(b) period. Options above. | none |
| class-or-type | needs-decision | 61.57(b)(1)(ii); 61.5(b)(2) | Same as by day, including the SEP/SET split. | none |
| no-simulator-credit | supported | 61.57(b)(2) | A full flight simulator counts only in a part 142 course with the visual system set to the period. | none |
| night-limitation-not-recorded | supported | 61.110(b)(1), 61.110(c); 61.101(e)(6); 61.315(c)(5), 61.329 | The "Night flying prohibited" limitation is not in the record; recreational pilots may not fly between sunset and sunrise; sport pilots only with the 61.329 endorsement. | none |

## faa.shared.passengers-tailwheel

| Interpretation | Verdict | Paragraph(s) | Reasoning | Action |
| --- | --- | --- | --- | --- |
| tailwheel-from-logbook | supported | 61.57(a)(1)(ii) | The condition depends on the airplane "to be flown"; the logbook is the only evidence, and missing tailwheel fields give unknown. | none |
| full-stops-stand-for-both | supported | 61.57(a)(1)(ii) | Every full-stop landing ends a segment that began with a takeoff from a standstill, so the number of such takeoffs equals the number of full-stop landings; the count stands for both. | none |

## faa.licence.atp-airplane

| Interpretation | Verdict | Paragraph(s) | Reasoning | Action |
| --- | --- | --- | --- | --- |
| atp-privileges-medical | supported | 61.23(a)(1)(i), 61.23(a)(2)(ii) | Pilot-in-command privileges of an ATP certificate need first class. As stated, lower privileges are not distinguished: a second-class holder reads unknown, and a first-class certificate past its ATP duration (6 months from age 40) reads expired, although it may still support commercial privileges. | none |

## faa.licence.atp-rotorcraft

| Interpretation | Verdict | Paragraph(s) | Reasoning | Action |
| --- | --- | --- | --- | --- |
| rating-items | supported | 61.56(c); 61.157(a)(1) | The flight review is per person. 61.157(a)(1) gives the ATP practical test for the helicopter class only, so a gyroplane rating on an ATP certificate carries lower privileges, as stated. | none |
| atp-privileges-medical | supported | 61.23(a)(1)(i) | As for airplanes. A gyroplane rating is held to the first-class requirement too, which the reading states; a second-class holder reads unknown for it. | none |

## faa.licence.commercial-airplane

| Interpretation | Verdict | Paragraph(s) | Reasoning | Action |
| --- | --- | --- | --- | --- |
| commercial-privileges-medical | supported | 61.23(a)(2)(ii), 61.23(a)(3)(i) | Commercial privileges other than balloon or glider need second class. Not stated but following from it: a first- or second-class certificate past its 12 commercial months reads expired, although it still supports private privileges. | none |

## faa.licence.commercial-rotorcraft

| Interpretation | Verdict | Paragraph(s) | Reasoning | Action |
| --- | --- | --- | --- | --- |
| rating-items | supported | 61.56(c) | As for the other rotorcraft files. | none |
| commercial-privileges-medical | supported | 61.23(a)(2)(ii), 61.23(a)(3)(i) | As for airplanes. | none |

## faa.licence.glider

| Interpretation | Verdict | Paragraph(s) | Reasoning | Action |
| --- | --- | --- | --- | --- |
| private-and-commercial | supported | 61.57(a)(1), 61.23(b)(3) | 61.57 is the same for every grade, and (b)(3) exempts any pilot certificate with a glider rating flown in a glider. | none |
| no-medical | supported | 61.23(b)(3) | As above. | none |
| powered-classes-elsewhere | contradicted | 61.57(a)(1)(ii); 1.1 "Glider" | Evaluating a powered class rating by its own category is right, but the example called a motorglider "not a glider", while 1.1 makes a self-launching glider a glider. | Reading corrected; ref to 1.1 added. |

## faa.licence.powered-lift-airship

| Interpretation | Verdict | Paragraph(s) | Reasoning | Action |
| --- | --- | --- | --- | --- |
| rating-items | supported | 61.56(c) | As for the rotorcraft files. | none |
| medical-grade-not-recorded | needs-decision | 61.23(a)(1)(i), (a)(2)(ii), (a)(3)(i) | The reading's consequences are stated correctly, but its premise is only partly true: the rating's licence kind shows a private, commercial or ATP certificate, and a private certificate carries private privileges only. Options above. | none |

## faa.licence.private-airplane

| Interpretation | Verdict | Paragraph(s) | Reasoning | Action |
| --- | --- | --- | --- | --- |
| glider-certificate-powered-ratings | supported | 61.57(a)(1), 61.23(a)(3)(i) | 61.57 does not depend on the grade, and evaluating at private level asks for the least medical a powered certificate can use. | none |

## faa.licence.private-rotorcraft

| Interpretation | Verdict | Paragraph(s) | Reasoning | Action |
| --- | --- | --- | --- | --- |
| rating-items | supported | 61.56(c) | The flight review is per person and decides the certificate's composite. | none |
| tailwheel-airplanes-only | supported | 61.57(a)(1)(ii) | The condition names "an airplane with a tailwheel". | none |
| glider-certificate-powered-ratings | supported | 61.57(a)(1), 61.23(a)(3)(i) | As for airplanes. | none |

## faa.licence.recreational

Checked specifically: 61.101(g). Fewer than 400 flight hours and no pilot-in-command time
"in an aircraft" in the 180 days before the flight bars acting as pilot in command until
flight training and an instructor's endorsement.

| Interpretation | Verdict | Paragraph(s) | Reasoning | Action |
| --- | --- | --- | --- | --- |
| four-hundred-hours-total | supported | 61.101(g); 1.1 "Flight time" | Flight hours are flight time in an aircraft, so devices are excluded. Including the day evaluated in the 180 days only adds time logged that day, the convention used for all day windows. | none |
| pic-time-in-any-aircraft | supported | 61.101(g) | The text says "in an aircraft", with no category limit; dual received is not PIC time for a recreational pilot. | none |
| endorsement-restores | supported | 61.101(g) | The text lifts the bar "until" the endorsement and sets no duration; letting the endorsement count like PIC time for 180 days, after which the same bar arises again, follows the paragraph's own period. Any instructor_endorsement event counts, as stated. | none |
| flight-limits-not-evaluated | supported | 61.101(a), 61.101(e) | These are limits on a given flight. | none |
| no-night-privilege | supported | 61.101(e)(6), 61.101(i)(3) | Sunset to sunrise covers the whole 61.57(b) period; the (i) exception is for solo training flights without passengers. | none |
| tailwheel-airplanes-only | supported | 61.57(a)(1)(ii) | As above. | none |

## faa.licence.sport

Checked specifically: the night regime (61.315(c)(5) and 61.329 as amended by Amdt.
61-159), unchanged from the previous review; and the medical. 61.23(c)(1)(ii) lets a sport
pilot use a driver's license, (b)(2) needs neither in a glider or balloon, and
(c)(1)(vi) requires one of them at night even in a glider or balloon.

| Interpretation | Verdict | Paragraph(s) | Reasoning | Action |
| --- | --- | --- | --- | --- |
| night-regime-change | supported | 61.329, 61.315(c)(5) | Dates and the SPORT_NIGHT condition match the amendment (see the previous review, `faa.privilege.sport-night`). | none |
| sport-classes | needs-decision | 61.315(a); 61.5(b)(1)(vi), (vii), (b)(5), (b)(6); 61.57(a)(1)(ii) | Recording powered parachutes and weight-shift-control aircraft as OTHER pools two categories, so recency in one counts for the other. Options above. | none |
| one-passenger | supported | 61.315(c)(4) | A flight limit, not recency. | none |
| no-glider-night | supported | 61.57(b)(1) | Same gap as `passengers-glider` / `night-not-evaluated`; stated. | none |
| tailwheel-airplanes-only | supported | 61.57(a)(1)(ii) | As above. | none |
| medical-or-drivers-license | needs-decision | 61.23(c)(1)(ii), 61.23(b)(2), 61.23(c)(1)(vi), 61.303 | The reading is accurate (no driver's-license credential file exists, although the record type US_DRIVERS_LICENSE does), but the result is a composite that never checks the medical qualification of a powered sport pilot. Options above. | none |

## faa.licence.student

| Interpretation | Verdict | Paragraph(s) | Reasoning | Action |
| --- | --- | --- | --- | --- |
| endorsement-as-event | supported | 61.87(n), 61.87(o), 61.87(p)(4) | The endorsement must be given within the 90 days preceding the flight; make and model and who gave the training are not checked, and the night solo endorsement is not evaluated, as stated. | none |
| medical-for-every-student | supported | 61.23(a)(3)(i), 61.23(b)(1), 61.23(c)(1)(i) | The certificate sought is not recorded; asking for a medical or BasicMed can only read unknown where none is needed. | none |

## Observations (not interpretations)

1. **SEP and SET are one FAA class.** 61.5(b)(2) issues single-engine land and sea class
   ratings; the record classes SEP_LAND/SET_LAND and SEP_SEA/SET_SEA split them. Besides
   passenger recency (see `class-or-type`), the `about: passengers` subjects are per
   recorded class, so a pilot with both recorded gets two results for one FAA class.
2. **Self-launching gliders logged as TMG.** TMG is an aeroplane class in the vocabulary;
   an FAA glider pilot whose self-launching glider flights are logged as TMG gets no glider
   passenger credit for them. A pooling rule for FAA glider ratings (vocabulary request)
   would close this together with observation 1.
3. **Licence kind from the type alias.** The student exemption now checks that no FAA_GLIDER,
   FAA_PRIVATE and so on is held. Licence kinds come from the type alias whatever the
   authority, so a non-FAA licence recorded with the bare type `GLIDER` or `PRIVATE` would
   also end the exemption. `holds` has no authority filter.
4. **Stepped-down medicals.** The commercial and ATP composites require the medical
   credential as a whole; a first- or second-class certificate past its commercial or ATP
   duration reads expired for the certificate, although it still supports lower
   privileges. The readings name the lower-class case (unknown) but not this one.
