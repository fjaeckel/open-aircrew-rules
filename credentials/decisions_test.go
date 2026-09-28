package credentials

import (
	"os"
	"path/filepath"
	"slices"
	"strings"
	"testing"

	"github.com/fjaeckel/open-aircrew-rules/engine"
)

// Capabilities added for the owner decisions of 2026-09-28 (docs/decisions-2026-09-28.md).

const decisionFiles = `
credential: Helicopter type
id: easa.rating.ht
kind: rating
authority: EASA
selects: { ratings: { classes: [HELICOPTER] } }
evaluations:
  - id: pax
    about: passengers
    asks: Passengers per type.
    source: easa:FCL.740.A(b)(1)
    only_for: { type_rated: true, if_missing: unknown }
    counting: { within_days: 90, in_type: true }
    passes_if: { takeoffs_and_landings: { min: 3, ref: easa:FCL.740.A(b)(1) } }
    outcomes: passengers
---
credential: Instrument rating by category
id: faa.rating.ir
kind: rating
authority: FAA
selects: { ratings: { classes: [IR], authorities: [FAA] } }
evaluations:
  - id: experience
    asks: Approaches in the rating's category.
    source: easa:FCL.740.A(b)(1)
    only_for: { categories: [aeroplane, helicopter], if_missing: unknown }
    counting: { within_calendar_months: 6, in_category: true, simulator: include }
    passes_if: { approaches: { min: 6, ref: easa:FCL.740.A(b)(1) } }
    restored_by:
      - { ipc: { in_category: true, simulator: include }, ref: easa:FCL.740.A(b)(1) }
    outcomes: recency
---
credential: FAA airplane passengers
id: faa.licence.pa
kind: licence
authority: FAA
selects: { licence: { kinds: [FAA_PRIVATE] }, ratings: { classes: [SEP_LAND, SET_LAND, MEP_LAND], authorities: [FAA] } }
evaluations:
  - id: pax
    about: passengers
    asks: Passengers by FAA class.
    source: easa:FCL.740.A(b)(1)
    relevant_class: { class_group: faa, ref: easa:FCL.740.A(b)(1) }
    counting: { within_days: 90 }
    passes_if:
      ref: easa:FCL.740.A(b)(1)
      all_of:
        - takeoffs_and_landings: { min: 3, in_class: true, ref: easa:FCL.740.A(b)(1) }
        - night_period_takeoffs: { min: 1, in_class: true, ref: easa:FCL.740.A(b)(1) }
    outcomes: passengers
  - id: review
    about: pilot
    asks: Checks under the FAA only, no devices.
    source: easa:FCL.740.A(b)(1)
    counting: { within_calendar_months: 24 }
    passes_if: { flight_review: { by_authority: [FAA], unknown_if_none: true, ref: easa:FCL.740.A(b)(1) } }
    outcomes: recency
  - id: class_check
    asks: Class checks, not IR-only ones; FFS counts.
    source: easa:FCL.740.A(b)(1)
    about: rating
    only_for: { classes: [SEP_LAND] }
    counting: { ever: true }
    passes_if: { proficiency_check: { in_class: true, excluding_ratings: [IR, BIR], simulator: include, fstd: [FFS], ref: easa:FCL.740.A(b)(1) } }
    outcomes: recency
---
credential: Sums and caps
id: easa.licence.sums
kind: licence
authority: EASA
selects: { licence: { kinds: [SPL] } }
evaluations:
  - id: course
    asks: 15 hours with up to 7 hours of credit, and a flight with an instructor aboard.
    source: easa:FCL.740.A(b)(1)
    counting: { ever: true }
    passes_if:
      ref: easa:FCL.740.A(b)(1)
      all_of:
        - id: hours
          ref: easa:FCL.740.A(b)(1)
          min_hours: 15
          sum_of:
            - instruction_time: { classes: [GLIDER], ref: easa:FCL.740.A(b)(1) }
            - pic_time: { classes: [SEP_LAND], max_hours: 7, ref: easa:FCL.740.A(b)(1) }
        - longest_flight: { min_hours: 1, any_flight_of: [{ with_time: [dual] }, { flagged: { instructorOnBoard: true } }], ref: easa:FCL.740.A(b)(1) }
        - id: tmg
          ref: easa:FCL.740.A(b)(1)
          only_if: { seeks: TMG, ref: easa:FCL.740.A(b)(1) }
          all_of:
            - dual_time: { min_hours: 4, classes: [TMG], ref: easa:FCL.740.A(b)(1) }
    outcomes: training
---
credential: Towing and instructing ultralights
id: de.privilege.tw
kind: privilege
authority: DE
selects: { privilege: { kinds: [UL_TOWING] } }
evaluations:
  - id: tows
    asks: Tows of the entered kind.
    source: easa:FCL.740.A(b)(1)
    counting: { within_months: 24 }
    passes_if: { landings: { min: 2, by_this_tow_kind: true, by_this_tow_take_up: true, ref: easa:FCL.740.A(b)(1) } }
    outcomes: privilege_recency
  - id: validity
    asks: Three years, BGB counting.
    source: easa:FCL.740.A(b)(1)
    valid_for: { counted_from: valid_from, periods: [{ months: 36, ref: easa:FCL.740.A(b)(1) }] }
    outcomes: { preset: validity, expiring_notice: { days: 1, ref: easa:FCL.740.A(b)(1) } }
---
credential: UL instructor
id: de.instructor.ui
kind: instructor_certificate
authority: DE
selects: { privilege: { kinds: [UL_INSTRUCTOR] } }
evaluations:
  - id: instruction
    asks: Instruction in the rating's kinds.
    source: easa:FCL.740.A(b)(1)
    counting: { ever: true }
    passes_if: { takeoffs: { min: 2, in_ul_kind: true, ref: easa:FCL.740.A(b)(1) } }
    outcomes: recency
---
credential: UL passengers by authorisation kind
id: de.rating.up
kind: rating
authority: DE
selects: { ratings: { classes: [ULTRALIGHT] } }
evaluations:
  - id: pax
    about: passengers
    asks: Authorisation for the kind flown.
    source: easa:FCL.740.A(b)(1)
    outcomes:
      - { id: auth, when: { holds: { privileges: [UL_PASSENGER_AUTH], details: [$subject], sameLicence: true } }, status: current, message: pax.day_current, ref: easa:FCL.740.A(b)(1) }
      - { id: no_kind, when: { holds: { privileges: [UL_PASSENGER_AUTH], details: [none], authorities: [DULV] } }, status: unknown, message: pax.ul_authorisation_missing, ref: easa:FCL.740.A(b)(1) }
      - { id: none, when: always, status: expired, message: pax.not_current, ref: easa:FCL.740.A(b)(1) }
---
credential: Night with a valid IR of the category
id: easa.rating.nt
kind: rating
authority: EASA
selects: { ratings: { classes: [SEP_SEA] } }
evaluations:
  - id: night
    about: passengers
    asks: IR of the category with an expiry.
    source: easa:FCL.740.A(b)(1)
    outcomes:
      - { id: ir, when: { holds: { classes: [IR], sameCategory: true, valid: true, expiryRecorded: true } }, status: current, message: pax.night_waived_ir, ref: easa:FCL.740.A(b)(1) }
      - { id: none, when: always, status: expired, message: pax.not_current, ref: easa:FCL.740.A(b)(1) }
  - id: variants
    about: variants
    asks: Only engine-type variants.
    source: easa:FCL.740.A(b)(1)
    only_for: { different_engine_type: true, if_missing: unknown }
    counting: { within_months: 24 }
    passes_if: { flights: { min: 1, in_variant: true, ref: easa:FCL.740.A(b)(1) } }
    outcomes: recency
  - id: credit
    asks: Ultralight motorglider time only with a fixed engine.
    source: easa:FCL.740.A(b)(1)
    counting: { ever: true }
    passes_if: { flight_time: { min_hours: 2, classes: [TMG], ul_credit: { TMG: { ul_kinds: [THREE_AXIS_MOTORGLIDER], fixed_engine: true }, ref: easa:FCL.740.A(b)(1) }, ref: easa:FCL.740.A(b)(1) } }
    outcomes: recency
---
credential: Class 1 with lower levels
id: easa.medical.c1
kind: medical
authority: EASA
selects: { credential: { kinds: [EASA_CLASS1_MEDICAL] } }
evaluations:
  - id: class_1
    level: class_1
    asks: Class 1 validity, age on the examination date.
    source: easa:FCL.740.A(b)(1)
    valid_for: { counted_from: valid_from, age_on: issue, periods: [{ age_from: 60, months: 6, ref: easa:FCL.740.A(b)(1) }, { months: 12, ref: easa:FCL.740.A(b)(1) }] }
    outcomes: { preset: validity, date_of_birth: when_no_expiry, expiring_notice: { days: 1, ref: easa:FCL.740.A(b)(1) } }
  - id: class_2
    level: class_2
    asks: Class 2 validity of the same certificate.
    source: easa:FCL.740.A(b)(1)
    valid_for: { counted_from: valid_from, age_on: issue, periods: [{ months: 60, ref: easa:FCL.740.A(b)(1) }] }
    outcomes: { preset: validity, expiring_notice: { days: 1, ref: easa:FCL.740.A(b)(1) } }
---
credential: Instructor certificate trusting the recorded date
id: easa.instructor.rc
kind: instructor_certificate
authority: EASA
selects: { privilege: { kinds: [FI] } }
evaluations:
  - id: validity
    asks: Recorded expiry wins.
    source: easa:FCL.740.A(b)(1)
    valid_for: { counted_from: valid_from, recorded_expiry_wins: true, periods: [{ months: 36, ref: easa:FCL.740.A(b)(1) }] }
    outcomes: { preset: validity, expiring_notice: { days: 1, ref: easa:FCL.740.A(b)(1) } }
---
credential: CPL with a PPL level
id: easa.licence.cp
kind: licence
authority: EASA
selects: { licence: { kinds: [CPL_A] } }
evaluations:
  - id: medical
    asks: Private privileges with class 1 or its class 2 validity.
    requires_any: [easa.medical.c1@class_2]
  - id: medical_commercial
    asks: Commercial privileges need class 1.
    requires_any: [easa.medical.c1@class_1]
    limits: commercial_privileges
---
credential: Driver's license
id: faa.document.dl
kind: document
authority: FAA
selects: { credential: { kinds: [US_DRIVERS_LICENSE] } }
evaluations:
  - id: validity
    asks: Valid until its expiry.
    source: easa:FCL.740.A(b)(1)
    outcomes: { preset: validity, expiring_notice: { days: 1, ref: easa:FCL.740.A(b)(1) } }
interpretations:
  - { id: p1, reading: A reading., ref: easa:FCL.740.A(b)(1), affects: [validity], principle: P1, approved_by: null, approved_on: null }
  - { id: p9, reading: A reading., ref: easa:FCL.740.A(b)(1), affects: [validity], principle: P9, approved_by: null, approved_on: null }
`

func decisionsRoot(t *testing.T) *Catalogue {
	t.Helper()
	files := map[string]string{}
	for _, doc := range strings.Split(decisionFiles, "\n---\n") {
		id := strings.TrimSpace(strings.SplitN(strings.SplitN(doc, "id: ", 2)[1], "\n", 2)[0])
		parts := strings.Split(id, ".")
		files["credentials/"+parts[0]+"/"+parts[1]+"s/"+parts[2]+".yaml"] = doc
	}
	cat, err := Load(writeRoot(t, files))
	if err != nil {
		t.Fatal(err)
	}
	for _, e := range cat.Errors {
		t.Error(e)
	}
	return cat
}

func decisionRecord(t *testing.T, s string) *engine.Record {
	t.Helper()
	var r engine.Record
	if err := yamlDecode(s, &r); err != nil {
		t.Fatal(err)
	}
	return &r
}

// result returns the evaluation of a rule for the subject with id (and detail when given).
func result(t *testing.T, res Result, rule, id string) engine.Evaluation {
	t.Helper()
	for _, e := range res.Evaluations {
		if e.RuleID == rule && e.Subject.ID == id {
			return e
		}
	}
	t.Fatalf("no evaluation of %s for %s in %+v", rule, id, res.Evaluations)
	return engine.Evaluation{}
}

func wantEval(t *testing.T, e engine.Evaluation, status, key string) {
	t.Helper()
	if e.Status != status || (key != "" && e.MessageKey != key) {
		t.Errorf("%s %s: got %s/%s, want %s/%s (rows %+v)", e.RuleID, e.Subject.ID, e.Status, e.MessageKey, status, key, e.Requirements)
	}
}

func TestDecisionCapabilities(t *testing.T) {
	cat := decisionsRoot(t)
	if !slices.ContainsFunc(Check(cat, nil).Interpret, func(s string) bool { return strings.Contains(s, `principle "P9"`) }) {
		t.Error("principle P9 not reported")
	}
	asOf := engine.MustDate("2026-09-15")
	res := cat.Evaluate(decisionRecord(t, `
holder: { dateOfBirth: 1960-01-01 }
licences:
  - { id: l-h, authority: EASA, type: PPL(H) }
  - { id: l-faa, authority: FAA, type: PRIVATE }
  - { id: l-spl, authority: EASA, type: SPL }
  - { id: l-ul, authority: DULV, type: UL }
  - { id: l-cpl, authority: EASA, type: CPL(A) }
  - { id: l-ppl, authority: EASA, type: PPL(A) }
ratings:
  - { id: r-h1, licenceId: l-h, class: HELICOPTER }
  - { id: r-h2, licenceId: l-h, class: HELICOPTER, typeDesignator: R44 }
  - { id: r-ir, licenceId: l-faa, class: IR, category: aeroplane }
  - { id: r-ir2, licenceId: l-faa, class: IR }
  - { id: r-sep, licenceId: l-faa, class: SEP_LAND }
  - { id: r-set, licenceId: l-faa, class: SET_LAND }
  - { id: r-ul, licenceId: l-ul, class: ULTRALIGHT, ulKind: THREE_AXIS }
  - { id: r-ulw, licenceId: l-ul, class: ULTRALIGHT, ulKind: WEIGHT_SHIFT }
  - { id: r-sea, licenceId: l-ppl, class: SEP_SEA }
  - { id: r-irp, licenceId: l-ppl, class: IR, expires: 2027-01-01 }
variants:
  - { id: v-el, ratingId: r-sea, name: ELECTRIC, differentEngineType: true }
  - { id: v-rg, ratingId: r-sea, name: RETRACTABLE }
  - { id: v-same, ratingId: r-sea, name: PISTON, differentEngineType: false }
privileges:
  - { id: p-tow, licenceId: l-ul, kind: UL_TOWING, detail: "banner pick_up", validFrom: 2024-01-10 }
  - { id: p-ui, licenceId: l-ul, kind: UL_INSTRUCTOR, detail: WEIGHT_SHIFT }
  - { id: p-auth, licenceId: l-ul, kind: UL_PASSENGER_AUTH, detail: THREE_AXIS }
  - { id: p-fi, licenceId: l-ppl, kind: FI, validFrom: 2020-01-01, expires: 2027-06-30 }
credentials:
  - { id: m1, type: EASA_CLASS1_MEDICAL, issued: 2025-01-10, validFrom: 2025-01-10 }
  - { id: dl, type: US_DRIVERS_LICENSE, expires: 2030-01-01 }
trainings:
  - { programme: SPL, seeks: [TMG] }
events:
  - { date: 2026-08-01, kind: ipc, class: SEP_LAND }
  - { date: 2026-01-01, kind: proficiency_check, class: SEP_LAND, rating: IR }
  - { date: 2026-02-01, kind: flight_review, authority: EASA }
flights:
  - { date: 2026-09-01, class: HELICOPTER, typeDesignator: R44, takeoffs: { day: 3 }, landings: { day: 3 } }
  - { date: 2026-09-01, class: SET_LAND, takeoffs: { day: 3, night: 1 }, landings: { day: 3 }, nightPeriodTakeoffs: 1 }
  - { date: 2026-08-01, class: SEP_LAND, isSimulator: true, fstdType: FFS, checkRating: IR, flags: { proficiencyCheck: true } }
  - { date: 2026-07-01, class: SEP_LAND, isSimulator: true, fstdType: FFS, flags: { proficiencyCheck: true } }
  - { date: 2020-01-01, class: GLIDER, minutes: { total: 600, dual: 600 } }
  - { date: 2020-01-02, class: SEP_LAND, minutes: { total: 900, pic: 900 } }
  - { date: 2020-01-03, class: GLIDER, minutes: { total: 70, pic: 70 }, flags: { instructorOnBoard: true } }
  - { date: 2026-03-01, class: ULTRALIGHT, ulKind: WEIGHT_SHIFT, landings: { day: 2 }, towKind: banner, towTakeUp: pick_up, takeoffs: { day: 2 }, flags: { towFlight: true } }
  - { date: 2026-03-02, class: ULTRALIGHT, ulKind: THREE_AXIS, landings: { day: 5 }, towKind: glider, towTakeUp: ground, flags: { towFlight: true } }
  - { date: 2026-03-03, class: ULTRALIGHT, ulKind: THREE_AXIS_MOTORGLIDER, minutes: { total: 120 }, fixedEngine: true }
  - { date: 2026-03-04, class: ULTRALIGHT, ulKind: THREE_AXIS_MOTORGLIDER, minutes: { total: 120 } }
  - { date: 2026-03-05, class: SEP_SEA, variant: ELECTRIC }
`), asOf)

	// A held type rating without a type designator is unknown, not dropped.
	missing := result(t, res, "easa.rating.ht#pax", "l-h")
	var sawUnknown bool
	for _, e := range res.Evaluations {
		if e.RuleID == "easa.rating.ht#pax" && e.Status == "unknown" && e.MessageKey == engine.InputMissingKey && e.MessageParams["input"] == "typeRated" {
			sawUnknown = true
		}
	}
	if !sawUnknown || missing.Subject.Kind != "passengers" {
		t.Errorf("type rating without designator: %+v", res.Evaluations)
	}

	// Category on the IR rating: recorded counts by it; unrecorded is unknown.
	wantEval(t, result(t, res, "faa.rating.ir#experience", "r-ir"), "current", "")
	wantEval(t, result(t, res, "faa.rating.ir#experience", "r-ir2"), "unknown", engine.InputMissingKey)

	// FAA class group: one passengers result for SEP and SET land; the SET flight counts.
	var groups []engine.Evaluation
	for _, e := range res.Evaluations {
		if e.RuleID == "faa.licence.pa#pax" {
			groups = append(groups, e)
		}
	}
	if len(groups) != 1 || groups[0].Subject.Group != "SINGLE_ENGINE_LAND" || groups[0].Status != "current" {
		t.Errorf("class group: %+v", groups)
	}
	// A review under another authority does not count; with nothing counted the row is untracked.
	if r := findRow(result(t, res, "faa.licence.pa#review", "l-faa"), "flight_review"); r.Tracked {
		t.Errorf("review under another authority: %+v", r)
	}
	// The IR-only FFS check is left out; the class FFS check counts.
	wantEval(t, result(t, res, "faa.licence.pa#class_check", "r-sep"), "current", "")

	// sum_of with a capped credit, longest flight with an instructor aboard, seeks.
	course := result(t, res, "easa.licence.sums#course", "l-spl")
	if r := findRow(course, "hours"); r.Current != 17*60 || !r.Met {
		t.Errorf("sum row %+v", r)
	}
	if r := findRow(course, "training_flight"); r.Current != 600 {
		t.Errorf("longest flight %+v", r)
	}
	wantEval(t, course, "lapsed", "training.in_progress")

	// Tow kind and take-up from the privilege detail; DE periods end the day before.
	wantEval(t, result(t, res, "de.privilege.tw#tows", "p-tow"), "current", "")
	if v := result(t, res, "de.privilege.tw#validity", "p-tow"); v.ExpiresOn == nil || v.ExpiresOn.String() != "2027-01-09" {
		t.Errorf("DE validity end %v", v.ExpiresOn)
	}
	wantEval(t, result(t, res, "de.instructor.ui#instruction", "p-ui"), "current", "")

	// Passenger authorisation detail per kind.
	for _, e := range res.Evaluations {
		if e.RuleID == "de.rating.up#pax" {
			want := map[string]string{"THREE_AXIS": "current", "WEIGHT_SHIFT": "expired"}[e.Subject.ULKind]
			if e.Status != want {
				t.Errorf("authorisation for %s: %s", e.Subject.ULKind, e.Status)
			}
		}
	}

	// holds sameCategory/expiryRecorded; variants by engine type; fixed-engine credit.
	wantEval(t, result(t, res, "easa.rating.nt#night", "l-ppl"), "current", "pax.night_waived_ir")
	wantEval(t, result(t, res, "easa.rating.nt#variants", "v-el"), "current", "")
	wantEval(t, result(t, res, "easa.rating.nt#variants", "v-rg"), "unknown", engine.InputMissingKey)
	for _, e := range res.Evaluations {
		if e.RuleID == "easa.rating.nt#variants" && e.Subject.ID == "v-same" {
			t.Error("same-engine variant evaluated")
		}
	}
	if r := findRow(result(t, res, "easa.rating.nt#credit", "r-sea"), "total_time"); r.Current != 120 {
		t.Errorf("fixed-engine credit %+v", r)
	}

	// Age on the examination date; levels and limits.
	c1 := result(t, res, "easa.medical.c1#class_1", "m1")
	if c1.ExpiresOn == nil || c1.ExpiresOn.String() != "2025-07-10" {
		t.Errorf("class 1 aged 65 at examination: %v", c1.ExpiresOn)
	}
	var cpl, med *Composite
	for i, c := range res.Credentials {
		switch c.Credential {
		case "easa.licence.cp":
			cpl = &res.Credentials[i]
		case "easa.medical.c1":
			med = &res.Credentials[i]
		}
	}
	if med == nil || med.Status != "expired" || len(med.Levels) != 2 || med.Levels[1].Status != "current" {
		t.Fatalf("medical composite %+v", med)
	}
	if cpl == nil || cpl.Status != "current" || len(cpl.Limitations) != 1 || cpl.Limitations[0].Scope != "commercial_privileges" || cpl.Limitations[0].Status != "expired" {
		t.Errorf("CPL composite %+v", cpl)
	}
	if st := med.status("nope"); st != "unknown" {
		t.Errorf("unknown level %s", st)
	}

	// Recorded expiry wins over the derived 3 years; documents validate.
	if v := result(t, res, "easa.instructor.rc#validity", "p-fi"); v.ExpiresOn.String() != "2027-06-30" {
		t.Errorf("recorded expiry %v", v.ExpiresOn)
	}
	wantEval(t, result(t, res, "faa.document.dl#validity", "dl"), "current", "")
}

func findRow(e engine.Evaluation, id string) engine.RequirementResult {
	for _, r := range e.Requirements {
		if r.ID == id {
			return r
		}
	}
	return engine.RequirementResult{}
}

func TestDecisionCapabilityErrors(t *testing.T) {
	root := writeRoot(t, map[string]string{"credentials/easa/ratings/x.yaml": `
credential: X
id: easa.rating.x
kind: rating
authority: EASA
selects: { ratings: { classes: [SEP_LAND] } }
evaluations:
  - id: a
    asks: if_missing alone.
    source: easa:FCL.740.A(b)(1)
    only_for: { if_missing: unknown }
    relevant_class: { pooled_with_held: [SEP_LAND], class_group: faa, ref: easa:FCL.740.A(b)(1) }
    outcomes: recency
  - id: b
    asks: bad words.
    source: easa:FCL.740.A(b)(1)
    only_for: { if_missing: maybe }
    relevant_class: { class_group: nope, ref: easa:FCL.740.A(b)(1) }
    valid_for: { age_on: birth, periods: [{ months: 1, ref: easa:FCL.740.A(b)(1) }] }
    passes_if:
      ref: easa:FCL.740.A(b)(1)
      min: 3
      all_of:
        - sum_of: []
          ref: easa:FCL.740.A(b)(1)
        - id: s
          ref: easa:FCL.740.A(b)(1)
          sum_of:
            - takeoffs: { min: 2, ref: easa:FCL.740.A(b)(1) }
            - flight_time: { ref: easa:FCL.740.A(b)(1) }
            - { takeoffs: {}, landings: {} }
          remedy: none
        - { takeoffs: { max: many, ref: easa:FCL.740.A(b)(1) } }
        - { takeoffs: { ul_credit: { TMG: { ul_kinds: [THREE_AXIS], fixed_engine: false }, ref: easa:FCL.740.A(b)(1) }, ref: easa:FCL.740.A(b)(1) } }
    outcomes: recency
  - id: c
    asks: A level that does not exist, and bad levels and scopes.
    requires_any: [easa.rating.x@nope]
    level: Bad Level
  - id: d
    asks: Unknown scope.
    requires_all: [easa.rating.y]
    limits: everything
    level: both
`, "credentials/easa/ratings/y.yaml": `
credential: Y
id: easa.rating.y
kind: rating
authority: EASA
selects: { ratings: { classes: [MEP_LAND] } }
evaluations:
  - { id: a, asks: A., requires_all: [easa.rating.x@both] }
`})
	write(t, root, "vocabulary.yaml", strings.Replace(readFile(t, root, "vocabulary.yaml"), "validity_ends: day_before, ref:", "validity_ends: never, ref:", 1))
	cat, err := Load(root)
	if err != nil {
		t.Fatal(err)
	}
	report := strings.Join(cat.Errors, "\n")
	for _, want := range []string{
		"if_missing needs type_rated: true",
		"relevant_class names pooled_with_held or class_group",
		`if_missing is unknown, not "maybe"`,
		`class_group "nope" is not in the vocabulary`,
		"valid_for.age_on is issue or valid_from",
		"min belongs to a count or a sum_of",
		"a sum_of needs an id",
		"a sum_of adds at least one count",
		"a sum_of item has no min",
		"a sum_of item is one count word",
		"sum_of items have one unit",
		"a sum_of needs min, min_hours or min_minutes",
		"takeoffs.max needs a number",
		"ul_credit.TMG: want { ul_kinds",
		`but easa.rating.x has no level "nope"`,
		`level "Bad Level" is not a lower-case name`,
		`limits "everything" is not a limitation scope`,
		"an entry has a level or limits, not both",
		"authority_conventions.DE.validity_ends is same_day or day_before",
	} {
		if !strings.Contains(report, want) {
			t.Errorf("report misses %q", want)
		}
	}
	if t.Failed() {
		t.Log(report)
	}
}

func readFile(t *testing.T, root, p string) string {
	t.Helper()
	b, err := os.ReadFile(filepath.Join(root, p))
	if err != nil {
		t.Fatal(err)
	}
	return string(b)
}
