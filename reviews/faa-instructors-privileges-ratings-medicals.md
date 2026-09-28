# Review: FAA instructors, privileges, ratings and medicals

Reviewed on 2026-09-28 against the verbatim texts in `sources/faa/` (eCFR Title 14, up to
date as of 2026-09-24) and, for effective dates, the Federal Register final rules stored as
excerpts. This review is not a sign-off: every interpretation keeps `approved_by: null` and
`approved_on: null` until a qualified reviewer approves it (CONTRIBUTING.md,
"Interpretations and sign-off").

Files: `credentials/faa/instructors/` (flight-instructor, ground-instructor),
`credentials/faa/privileges/` (category-ii-iii, glider-towing, sport-night),
`credentials/faa/ratings/instrument-airplane.yaml`, `credentials/faa/medicals/` (basicmed,
first-class; second-class and third-class have no interpretations of their own and use the
shared medical-duration evaluations).

## Summary

| Verdict | Count |
| --- | --- |
| supported | 32 |
| needs-decision | 3 |
| contradicted (fixed) | 1 |
| **Total** | **36** |

By group: instructors 14 (12 supported, 1 needs-decision, 1 contradicted), privileges 10
(10 supported), ratings 7 (5 supported, 2 needs-decision), medicals 5 (5 supported).

Contradicted and fixed:

- `faa.instructor.flight-instructor` / `other-means-not-recorded`: the reading treated
  61.197(a)(1) as "the practical test that led to the certificate". The paragraph starts the
  24 calendar months from the month the FAA issued the certificate, whatever led to the
  issue, so every certificate issued on or after 1 December 2024 is current at least
  through 31 December 2026. Three worked examples reported such a certificate (issued
  2024-12-05) as lapsed or current only to August 2026. The reading now states (a)(1) as
  written and the limitation (the recorded issue date is not counted); the three examples
  moved to 2027 dates, with their outcomes unchanged.

Needs a decision by the signing reviewer:

- `faa.instructor.flight-instructor` / `regime-by-recorded-expiry`: which regime applies,
  decided by the recorded expiry (file) or by the recorded issue date (needs a stage
  condition on the issue date).
- `faa.rating.instrument-airplane` / `any-powered-category`: tasks in any powered category
  count (file), or only in the category whose instrument privileges are kept (needs the
  category on the record's IR rating).
- `faa.rating.instrument-airplane` / `grace-from-last-met`: the timing is right, but the
  grace stage reports `expiring`, which the composite treats as "may be exercised", while
  61.57(c) bars IFR as pilot in command during those six months.

Sources added: `sources/faa/61.67.md`, `sources/faa/61.68.md` (eCFR versioner, as of
2026-09-24), `sources/faa/fr-2024-22018.md` (89 FR 80020, excerpts: DATES, the regulatory
evaluation on removing the expiration date, amendatory instructions 6 and 10) and
`sources/faa/fr-2025-13972.md` (90 FR 35034, excerpts: DATES, amendatory instructions 37,
46 and 49).

Also changed (no verdict involved): example `says:` texts in
`examples/faa.instructor.flight-instructor.yaml` and `examples/faa.medical.second-class.yaml`
named retired rule ids (`faa.14cfr61...`); they now name the current evaluations. Two
composite examples were added to `examples/faa.privilege.glider-towing.yaml` for
`powered-certificate`.

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

## faa.instructor.flight-instructor

Checked specifically: the regime switch. 61.19(d)(1) and (d)(2) split certificates by issue
date (on or after / before 1 December 2024); the final rule is effective 1 December 2024
(89 FR 80020, DATES; the only later instruction, number 10, amends 61.51 from 1 March 2027
and is not used here). `recent_experience` has `effective_from: 2024-12-01`, matching the
date in 61.19(d)(1). 61.425(a) says "issued after December 1, 2024" where 61.19(d)(1) says
"on or after"; the file follows 61.19(d)(1), which only matters for a certificate issued on
1 December 2024 itself.

| Interpretation | Verdict | Paragraph(s) | Reasoning | Action |
| --- | --- | --- | --- | --- |
| regime-by-recorded-expiry | needs-decision | 61.19(d)(1), 61.19(d)(2) | The text keys the regime to the issue date, and the record holds an issue date. The file uses the recorded expiry as the proxy because no stage condition reads the issue date; an old certificate recorded without its expiry is evaluated under the new regime (and read lapsed or current by recent experience). Options: (a) recorded expiry decides (file); (b) the issue date decides, an old certificate with no expiry recorded reports unknown (needs a vocabulary word such as an `issued_before` stage condition). | decided 2026-09-28: option 1, applied (file's choice kept; reading rewritten) |
| old-regime-not-bounded | supported | 61.19(d)(2); 89 FR 80020 DATES | A certificate issued, renewed or reinstated in November 2024 expires at the end of November 2026, so the expiry evaluation must keep running past the rule change. | none |
| renewal-is-a-new-record | supported | 61.197(e), 61.199(a), 61.19(d)(1); 89 FR 80020 regulatory evaluation | 61.197(e) lets an unexpired old certificate be renewed by recent experience before its expiration month; a certificate issued from 1 December 2024 has no expiry, and the rule's regulatory evaluation speaks of the holder receiving a permanent certificate without an expiration date. Recording it as a new privilege fits. | none |
| practical-test-events | supported | 61.197(b)(1), 61.197(d) | (b)(1) counts a practical test for a rating on the flight instructor certificate or an additional flight instructor rating; a pilot instrument rating test is neither. (d) allows a simulator under a part 142 course, taken as met. | none |
| refresher-counts-24-months | supported | 61.197(a)(2), 61.197(b)(2)(iii) | Completing the course is one of the (b) means, and (a)(2) starts the 24 calendar months from the month it is accomplished. The 3-calendar-month limit in (b)(2)(iii) concerns when the course is completed relative to submitting the documentation, which is not recorded. | none |
| other-means-not-recorded | contradicted | 61.197(a)(1), 61.197(b)(2)(i), (ii), (iv), (v) | The list of unrecordable means is right. The last clause misread (a)(1): the period starts "from" the month the FAA issued the certificate, for any issue, not only for the practical test behind it. As a result the examples dated a December 2024 certificate's recent experience from older events and reported it lapsed in August 2026, although (a)(1) keeps it current through 31 December 2026. | Reading rewritten (the issue month starts a period; the recorded issue date is not counted, so a recently issued certificate reads lapsed until its test or refresher is recorded). Examples `recent-experience-practical-test-reinstatement`, `recent-experience-refresher-reinstatement` and `recent-experience-window-first-day` moved to 2027 dates so that the issue month no longer governs; outcomes unchanged. |
| carry-forward-not-modelled | supported | 61.197(a)(3) | (a)(3) lets experience in the 3 calendar months before the last month keep the period running from that last month; not modelling it can only shorten the period shown, as the reading says. | none |
| reinstatement-window | supported | 61.199(a)(1), (a)(2), (a)(3) | Within 3 calendar months after the last month a refresher course or a practical test reinstates; later only a practical test. The engine's 3-calendar-month look-back matches this, and the late-refresher blind spot is stated. | none |
| sport-rating-same-rule | supported | 61.425(a), 61.425(b), 61.427(a), 61.427(b) | 61.425 sends both regimes to 61.197, and 61.427 has the same two reinstatement steps (refresher within 3 calendar months, else a practical test). | none |
| no-pilot-certificate-requirement | supported | 61.23(b)(7) | No medical is needed to exercise flight instructor privileges when not acting as pilot in command or required crewmember; 61.23(a)(3)(ii) applies only when acting as such. | none |

## faa.instructor.ground-instructor

Checked specifically: recency. 61.217 requires one of four activities during the preceding
12 calendar months; the evaluation counts 12 calendar months plus the current month, the
engine's calendar-month window.

| Interpretation | Verdict | Paragraph(s) | Reasoning | Action |
| --- | --- | --- | --- | --- |
| ground-instruction-not-recorded | supported | 61.217(a), 61.217(b) | Ground instruction given is not in the record; keeping it as an untracked alternative makes a missing record unknown rather than lapsed, consistent with policy `unknown-input`. | none |
| flight-instruction-given | supported | 61.217(b) | (b) counts activity as a flight instructor giving ground or flight training; logged instruction given, in an aircraft or a device, is such activity. Any amount counts, as the text sets no minimum. | none |
| events-as-waivers | supported | 61.217(c), 61.217(d) | The refresher course with its graduation certificate and the endorsement of demonstrated knowledge are events; the text needs only one of the four. Who endorsed and the subject areas are not checked, as stated. | none |
| no-expiry | supported | 61.19(e) | A ground instructor certificate is issued without an expiration date. | none |

## faa.privilege.category-ii-iii

Checked specifically: the section. The duration of a Category II or III pilot
authorization is 61.21 ("for other than part 121 and part 135 use"); 61.67 and 61.68 are the
issue requirements (certificate, instrument rating, experience, practical test). The refs
to 61.21 are correct and were kept; 61.67 and 61.68 were stored for the observation below.

| Interpretation | Verdict | Paragraph(s) | Reasoning | Action |
| --- | --- | --- | --- | --- |
| renewal-date | supported | 61.21(a) | The authorization expires at the end of the sixth calendar month after the month of issue or renewal; the derived expiry from validFrom (else the issue date) with `end_of_month` matches. | none |
| renewal-rules-not-modelled | supported | 61.21(b), 61.21(c), 61.21(d) | Per-type renewal and the 12-calendar-month limit are not modelled; a test passed in the month before expiry counts as passed in the expiry month under (d), and the reading says how to record it. | none |

## faa.privilege.glider-towing

Checked specifically: "accompanied". 61.69(a)(6)(i) requires three actual or simulated tows
"while accompanied by a qualified pilot"; the phrase follows both kinds of tow. Its
predecessor (62 FR 16220, 1997, not stored) read "three actual glider tows while
accompanied", which confirms that actual tows need the accompanying pilot too. Unlike
(a)(4), where only the simulated procedures carry the accompaniment, (a)(6)(i) gives no
room for counting solo tows.

| Interpretation | Verdict | Paragraph(s) | Reasoning | Action |
| --- | --- | --- | --- | --- |
| accompanied-tows | supported | 61.69(a)(6)(i) | See above. Checking that the accompanying pilot meets 61.69 is not possible from the record, as stated. | none |
| simulated-tows | supported | 61.69(a)(6)(i) | Simulated tows count under the text but cannot be told from other flights; not counting them only under-reports. | none |
| towed-glider-flights | supported | 61.69(a)(6)(ii) | Three flights as pilot in command of a glider towed by an aircraft; aerotowed glider flights with PIC time are exactly that. Towed unpowered ultralights have no class in the record. | none |
| calendar-months-before-flight | supported | 61.69(a)(6) | "Within 24 calendar months before the flight" read as the 24 calendar months before the month plus the month to date, the same calendar-month construction used for 61.57 and 61.197. | none |
| prerequisites-not-evaluated | supported | 61.69(a)(2), (a)(3), (a)(4), (a)(5) | These are one-time qualifications with no recency and no expiry. | none |
| powered-certificate | supported | 61.69(a)(1) | (a)(1) requires a private, commercial or airline transport certificate with a powered category rating. The requires_any list names the airplane and the private and commercial rotorcraft credentials; an ATP certificate with only a rotorcraft rating is still covered, because `faa.licence.atp-airplane` selects the ATP licence item itself (shown by the new example `composite-atp-helicopter`). A certificate with only a glider rating has no powered credential and reads unknown (`composite-glider-certificate-only`). | Two composite examples added. See observation 2. |

## faa.privilege.sport-night

Checked specifically: the effective date. 61.329 and 61.23(c)(1)(vi) were added, and
61.315(c)(5) changed to "except as provided in § 61.329", by amendatory instructions 49, 37
and 46 of 90 FR 35034, effective 22 October 2025 (none of the three is among the
instructions delayed to 24 July 2026). This file has only requirement entries, which carry
no dates; the date is applied in `faa.licence.sport` (`no_night_privileges` to 2025-10-21,
night passenger recency from 2025-10-22), which agrees with the rule.

| Interpretation | Verdict | Paragraph(s) | Reasoning | Action |
| --- | --- | --- | --- | --- |
| endorsement-as-privilege | supported | 61.329(a), 61.329(c) | The endorsement certifies the 3 hours of night training (with the night cross-country and 10 full-stop night takeoffs and landings) and night proficiency in the category and class; recording it as the privilege and not recounting the training fits. | none |
| medical-for-night | supported | 61.329(b), 61.23(c)(1)(vi) | 61.329(b) requires a part 67 medical or the 61.113(i) conditions, which include the course and examination of 61.23(c)(3); 61.23(c)(1)(vi) overrides the glider and balloon exceptions of (b)(1), (b)(2) and (b)(6) at night. A driver's license alone does not meet it. | none |

## faa.rating.instrument-airplane

Checked specifically: the windows. 61.57(c)(1) counts the 6 calendar months preceding the
month of the flight; 61.57(d)(1) requires an IPC after more than six calendar months of not
meeting (c). The engine's calendar-month windows give currency through the end of the sixth
month after the tasks and the IPC-free recovery through the end of the twelfth.

| Interpretation | Verdict | Paragraph(s) | Reasoning | Action |
| --- | --- | --- | --- | --- |
| any-powered-category | needs-decision | 61.57(c)(1), 61.57(c)(2) | The text requires the tasks in the category "as appropriate, for the instrument rating privileges to be maintained", and a device representing that category. The record's IR rating does not name a category, so the file counts any powered category. Options: (a) any powered category (file; a helicopter's approaches keep an airplane IR current); (b) airplane only, as the credential's name says (wrong for helicopter IR holders); (c) record the category on the IR rating and count by it (vocabulary and record change). | decided 2026-09-28: option 3, applied (the IR rating records `category`; `only_for.categories` with `if_missing: unknown`, `in_category: true`; examples `instrument-experience-helicopter-category`, `instrument-experience-other-category-not-counted`, `instrument-experience-category-not-recorded`) |
| grace-from-last-met | needs-decision | 61.57(c), 61.57(d)(1) | The timing is right: after the last day the tasks were met, the pilot can still regain currency without an IPC until the end of the sixth calendar month after lapse. But the `grace` stage reports `expiring`, and the composite reads `expiring` as "may be exercised", while under 61.57(c) the pilot may not act as pilot in command under IFR during that time. Options: (a) keep `expiring` (file; docs/credential-format.md names this stage as a failing `expiring`); (b) report `lapsed` with `policy:recency-lapsed` and keep the message, so the composite does not show IFR privileges as usable. | decided 2026-09-28: option 2, applied (grace stage `lapsed` with `policy:recency-lapsed`, message kept; example `instrument-experience-grace-last-day` expiring to lapsed, `instrument-experience-grace-first-day` added) |
| ipc-restores | supported | 61.57(d)(1), 61.57(d)(2) | An IPC re-establishes currency; (d)(2) allows an aircraft or a representative device. Who conducted it ((d)(3)) is not checked, as stated. | fix F1 applied 2026-09-28: the IPC counts only in the rating's category (61.57(d)(2)); example `instrument-experience-ipc-other-category` |
| intercept-track-recorded | supported | 61.57(c)(1)(iii) | The task is required; missing input gives unknown under policy `unknown-input`. | none |
| no-recorded-expiry | supported | 61.19(c)(1) | Pilot certificates, and the ratings on them, are issued without an expiration date. | none |
| exceptions-not-modelled | supported | 61.57(e) | The part 121, 125 and 135 exceptions depend on employment the record does not show. | none |
| no-privileges-basis | supported | 61.315(c)(12), 61.315(c)(13), 61.101(e)(9), 61.101(e)(10) | Sport and recreational pilots may not fly below 3 statute miles visibility or without visual reference to the surface, so there are no instrument privileges to keep current. | none |

## faa.medical.basicmed

Checked specifically: Part 68 and 61.23(c)(3). The course is required during the 24 and the
examination during the 48 calendar months before acting as pilot in command; 61.113(i)
requires a valid U.S. driver's license.

| Interpretation | Verdict | Paragraph(s) | Reasoning | Action |
| --- | --- | --- | --- | --- |
| valid-until-both-count | supported | 61.23(c)(3)(i)(C), 61.23(c)(3)(i)(D) | Same calendar-month construction as elsewhere in part 61; the result lasts while both still count. | none |
| drivers-license | supported | 61.113(i) | A valid U.S. driver's license is a condition; missing input gives unknown, an expired license lapsed. | none |
| not-modelled | supported | 61.23(c)(3)(i)(B), 61.23(c)(3)(i)(E), 61.23(c)(3)(ii), 68.9, 61.113(i) | These are conditions on the medical history, care and the flight, none of which the record holds. | none |

## faa.medical.first-class

Checked specifically: the 61.23(d) table. First class: ATP privileges 12 months under 40 and
6 months from 40 ((d)(1)(i), (ii)); commercial 12 months at any age ((d)(1)(iii)); private,
recreational, student, sport and instructor-as-PIC 60 months under 40 and 24 from 40
((d)(1)(iv), (v)). Second class: commercial 12 months ((d)(2)(i)); private 60/24 ((d)(2)(ii),
(iii)). Third class: 60/24 ((d)(3)(i), (ii)). Each runs to the end of the last day of the
month, counted from the month of examination, with the age on the examination date. The
first-class file and the two shared duration evaluations match the table; the examples
check both age bands and the month ends.

| Interpretation | Verdict | Paragraph(s) | Reasoning | Action |
| --- | --- | --- | --- | --- |
| steps-down | supported | 61.23(d) | The table gives one duration per operation for the same certificate; one evaluation per privilege level expresses that. | none |
| part-121-age-60-not-modelled | supported | 61.23(d)(1)(ii) row; 61.23(a)(1)(iii) | The 6-month duration for a part 121 flightcrew member aged 60 or more depends on operations the record does not show. The age is taken on the examination date, as the table's second column says. | none |

## Observations (not interpretations)

1. **Category II/III requirements.** 61.67(a)(1) and 61.68(a)(1) require a private or
   commercial certificate with an instrument rating, or an ATP certificate. The `certificate`
   requirement of `faa.privilege.category-ii-iii` names the pilot certificates only; it does
   not require the instrument rating or its currency. A reviewer may want the composite to
   require `faa.rating.instrument-airplane` for private and commercial holders.
2. **Requirement lists.** The `requires_any` lists of `faa.privilege.category-ii-iii` and
   `faa.privilege.glider-towing` do not name `faa.licence.atp-rotorcraft` or
   `faa.licence.powered-lift-airship`. The outcome does not change for an ATP holder
   (`faa.licence.atp-airplane` selects the licence item), but a powered-lift or airship
   rating on a glider-only certificate is not recognised as a powered category. Adding the
   two ids would not change any existing example.
3. **CFI issue month.** Following the fix, a flight instructor certificate issued in the
   last 24 calendar months with no recorded practical test or refresher reads lapsed,
   although 61.197(a)(1) keeps it current. A vocabulary request (count the privilege's
   issue month as an event, or a stage condition on the issue date) would close this and
   the `regime-by-recorded-expiry` decision together.

Status after the owner decisions of 2026-09-28 (docs/decisions-2026-09-28.md): observations
1 and 2 are fixed (F7): `faa.privilege.category-ii-iii` names every powered-category
certificate credential and requires an instrument rating or an airline transport pilot
certificate (requirement `instrument_rating`), and `faa.privilege.glider-towing` names the
ATP rotorcraft and the three powered-lift and airship credentials. Observation 3 is closed
by the decision to keep `regime-by-recorded-expiry`. The Legal Interpretation to Ken
Williams cited for `any-powered-category` carries no legible date; the year 2012 is the one
in the FAA file title (F4).
