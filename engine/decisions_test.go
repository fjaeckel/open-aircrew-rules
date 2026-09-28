package engine

import "testing"

// Engine capabilities of the owner decisions of 2026-09-28.

func TestMonthEndWindow(t *testing.T) {
	// AMC1 SFCL.160(a)(1)(ii)(d): 24 months counted to the end of the month of the flight.
	c := testCatalogue(t, dedent(`
id: t.month_end
applies_to: { subject: licence }
window: { calendar_months: 24 }
requirements: { id: f, metric: flights, min: 1, unit: flights }
stages:
  - { when: all_met, status: current, messageKey: rating.recency_current }
  - { when: always, status: lapsed, messageKey: rating.recency_not_met }
`))
	rec := "licences: [{ id: l, authority: EASA, type: SPL }]\nflights: [{ date: 2024-03-10, class: GLIDER }]"
	evs, _ := evalOne(t, c, rec, "2026-03-31")
	wantStatus(t, evs, 0, "current", "rating.recency_current")
	if evs[0].ValidUntil == nil || evs[0].ValidUntil.String() != "2026-03-31" {
		t.Errorf("valid until %v", evs[0].ValidUntil)
	}
	evs, _ = evalOne(t, c, rec, "2026-04-01")
	wantStatus(t, evs, 0, "lapsed", "rating.recency_not_met")
}

func TestEventAndFlightFilterAdditions(t *testing.T) {
	c := testCatalogue(t, dedent(`
id: t.events
applies_to: { subject: privilege }
requirements:
  all_of:
    - { id: by_faa, metric: events, min: 1, unit: events, filter: { eventKinds: [flight_review], eventAuthorities: [FAA] } }
    - { id: not_ir, metric: events, min: 1, unit: events, filter: { eventKinds: [proficiency_check], excludeEventRatings: [IR], simulator: include, fstdTypes: [FFS], categories: [aeroplane] } }
    - { id: tows, metric: landings.total, min: 1, unit: landings, filter: { towKinds: [$subject], towTakeUps: [$subject] } }
    - { id: none, metric: events, min: 1, unit: events, unknownIfNone: true, filter: { eventKinds: [safety_training] } }
    - { id: kind, metric: flights, min: 1, unit: flights, filter: { ulKinds: [$subject], categories: [$subject] } }
stages:
  - { when: all_met, status: current, messageKey: rating.recency_current }
  - { when: undetermined, status: unknown, messageKey: rating.recency_current }
  - { when: always, status: lapsed, messageKey: rating.recency_not_met }
`))
	evs, _ := evalOne(t, c, `
licences: [{ id: l, authority: DULV, type: UL }]
privileges: [{ id: p, licenceId: l, kind: UL_TOWING }]
events:
  - { date: 2026-01-01, kind: flight_review }
  - { date: 2026-01-02, kind: proficiency_check, rating: IR, ratings: [SEP_LAND], category: aeroplane, isSimulator: true }
  - { date: 2026-01-03, kind: proficiency_check, rating: IR, class: SEP_LAND }
flights:
  - { date: 2026-01-04, class: ULTRALIGHT, towKind: banner, flags: { towFlight: true }, landings: { day: 1 } }
  - { date: 2026-01-05, class: ULTRALIGHT, landings: { day: 1 } }
`, "2026-02-01")
	for id, want := range map[string]bool{"by_faa": false, "not_ir": false, "tows": false, "none": false, "kind": false} {
		if r := row(t, evs[0], id); r.Tracked != want {
			t.Errorf("%s tracked %v, want %v (%+v)", id, r.Tracked, want, r)
		}
	}
	wantStatus(t, evs, 0, "unknown", "rating.recency_current")
	evs, _ = evalOne(t, c, `
licences: [{ id: l, authority: DULV, type: UL }]
privileges: [{ id: p, licenceId: l, kind: UL_TOWING, detail: "banner; ground, THREE_AXIS" }]
events:
  - { date: 2026-01-01, kind: flight_review, authority: faa }
  - { date: 2026-01-02, kind: proficiency_check, rating: IR, ratings: [SEP_LAND], category: aeroplane, isSimulator: true, fstdType: FFS }
  - { date: 2026-01-06, kind: safety_training }
flights:
  - { date: 2026-01-04, class: ULTRALIGHT, towKind: banner, towTakeUp: ground, flags: { towFlight: true }, landings: { day: 1 } }
  - { date: 2026-01-05, class: ULTRALIGHT, towKind: banner, towTakeUp: pick_up, flags: { towFlight: true }, landings: { day: 1 } }
  - { date: 2026-01-05, class: ULTRALIGHT, ulKind: THREE_AXIS, landings: { day: 1 } }
`, "2026-02-01")
	wantStatus(t, evs, 0, "current", "rating.recency_current")
	if r := row(t, evs[0], "tows"); r.Current != 1 {
		t.Errorf("tows of the entered kind and take-up: %+v", r)
	}
}

func TestSeeksHoldsAndSum(t *testing.T) {
	c := testCatalogue(t, dedent(`
id: t.seeks
applies_to: { subject: training, programme: SPL }
requirements:
  id: s
  min: 2
  unit: flights
  sum_of:
    - { id: a, metric: flights, unit: flights, filter: { classes: [GLIDER] } }
    - { id: b, metric: full_stop_landings, unit: flights, max: 1 }
stages:
  - { when: { seeks: [TMG] }, status: not_applicable, messageKey: training.all_met }
  - { when: { unmet: s }, status: lapsed, messageKey: training.in_progress }
  - { when: { met: s }, status: current, messageKey: training.all_met }
  - { when: { holds: { classes: [SEP_LAND], authorities: [EASA], expiryRecorded: true } }, status: expired, messageKey: training.all_met }
  - { when: { holds: { privileges: [FI], details: [none] } }, status: expiring, messageKey: training.all_met }
  - { when: always, status: unknown, messageKey: training.distance_unknown }
`))
	rec := `
licences: [{ id: l, authority: EASA, type: SPL }]
ratings: [{ id: r, licenceId: l, class: SEP_LAND }]
privileges: [{ id: p, licenceId: l, kind: FI }]
trainings: [{ programme: SPL, seeks: [SAILPLANE] }]
flights: [{ date: 2026-01-01, class: GLIDER }, { date: 2026-01-02, class: SEP_LAND }]
`
	evs, _ := evalOne(t, c, rec, "2026-02-01")
	wantStatus(t, evs, 0, "expiring", "training.all_met")
	if r := row(t, evs[0], "s"); r.Tracked || r.Current != 1 || r.Metric != "sum_of" {
		t.Errorf("sum with an unknown item: %+v", r)
	}
	evs, _ = evalOne(t, c, rec+"\n", "2026-01-01")
	wantStatus(t, evs, 0, "expiring", "training.all_met")
	evs, _ = evalOne(t, c, `
licences: [{ id: l, authority: EASA, type: SPL }]
flights: [{ date: 2026-01-01, class: GLIDER, fullStopLandings: 3 }, { date: 2026-01-02, class: GLIDER, fullStopLandings: 0 }]
`, "2026-02-01")
	wantStatus(t, evs, 0, "current", "training.all_met")
	if r := row(t, evs[0], "s"); r.Current != 3 {
		t.Errorf("capped sum: %+v", r)
	}
	evs, _ = evalOne(t, c, "licences: [{ id: l, authority: EASA, type: SPL }]\ntrainings: [{ programme: spl, seeks: [TMG] }]", "2026-02-01")
	wantStatus(t, evs, 0, "not_applicable", "training.all_met")
	evs, _ = evalOne(t, c, "licences: [{ id: l, authority: EASA, type: SPL }]\nflights: [{ date: 2026-01-01, class: TMG, fullStopLandings: 0 }]", "2026-02-01")
	wantStatus(t, evs, 0, "lapsed", "training.in_progress")
}

func TestSelectionMissingAndValidityOptions(t *testing.T) {
	c := testCatalogue(t, dedent(`
id: t.missing
applies_to: { subject: rating, ulKinds: [THREE_AXIS], unknownWhenMissing: [ulKinds] }
validity: { from: valid_from, age_on: valid_from, periods: [{ when: { age_under: 40 }, months: 12 }] }
stages:
  - { when: always, status: current, messageKey: rating.recency_current }
`))
	evs, _ := evalOne(t, c, `
holder: { dateOfBirth: 2000-01-01 }
licences: [{ id: l, authority: DULV, type: UL }]
ratings:
  - { id: a, licenceId: l, class: ULTRALIGHT }
  - { id: b, licenceId: l, class: ULTRALIGHT, ulKind: THREE_AXIS, validFrom: 2026-01-10 }
  - { id: c, licenceId: l, class: ULTRALIGHT, ulKind: WEIGHT_SHIFT }
`, "2026-02-01")
	if len(evs) != 2 || evs[0].MessageKey != InputMissingKey || evs[0].MessageParams["input"] != "ulKinds" || evs[1].ExpiresOn.String() != "2027-01-10" {
		t.Errorf("%+v", evs)
	}
	c = testCatalogue(t, dedent(`
id: t.cat
applies_to: { subject: privilege, categories: [helicopter], unknownWhenMissing: [categories] }
stages:
  - { when: always, status: current, messageKey: rating.recency_current }
`))
	evs, _ = evalOne(t, c, `
licences: [{ id: h, authority: EASA, type: PPL(H) }, { id: x, authority: EASA, type: FI }, { id: a, authority: EASA, type: PPL(A) }]
privileges: [{ id: p1, licenceId: h, kind: FI }, { id: p2, licenceId: x, kind: FI }, { id: p3, licenceId: a, kind: FI }]
`, "2026-02-01")
	if len(evs) != 2 || evs[0].Status != "current" || evs[1].MessageKey != InputMissingKey {
		t.Errorf("%+v", evs)
	}
	c.Rules[0].AppliesTo.UnknownWhenMissing = nil
	evs, _ = evalOne(t, c, "licences: [{ id: x, authority: EASA, type: FI }]\nprivileges: [{ id: p2, licenceId: x, kind: FI }]", "2026-02-01")
	if len(evs) != 0 {
		t.Errorf("category not known and not asked for: %+v", evs)
	}
}

// TestRecordedExpiryMinimum: the recorded expiry (another level's date) is a lower bound;
// without a derived end it holds until it passes, then the expiry is unknown.
func TestRecordedExpiryMinimum(t *testing.T) {
	c := testCatalogue(t, dedent(`
id: t.min
applies_to: { subject: credential }
validity: { from: issued, age_on: issued, recorded_min: true, periods: [{ when: { age_under: 40 }, months: 60 }, { months: 24 }] }
stages: [{ when: always, status: current, messageKey: credential.valid }]
`))
	for _, tc := range []struct{ rec, asOf, want string }{
		{"holder: { dateOfBirth: 1980-05-20 }\ncredentials: [{ id: m, type: OTHER, issued: 2025-03-01, expires: 2026-03-01 }]", "2026-09-15", "2027-03-01"},
		{"holder: { dateOfBirth: 1980-05-20 }\ncredentials: [{ id: m, type: OTHER, issued: 2025-03-01, expires: 2028-03-01 }]", "2026-09-15", "2028-03-01"},
		{"credentials: [{ id: m, type: OTHER, issued: 2025-03-01, expires: 2026-03-01 }]", "2026-03-01", "2026-03-01"},
		{"credentials: [{ id: m, type: OTHER, issued: 2025-03-01, expires: 2026-03-01 }]", "2026-03-02", ""},
	} {
		evs, _ := evalOne(t, c, tc.rec, tc.asOf)
		got := ""
		if evs[0].ExpiresOn != nil {
			got = evs[0].ExpiresOn.String()
		}
		if got != tc.want {
			t.Errorf("%s on %s: expires %q, want %q", tc.rec, tc.asOf, got, tc.want)
		}
	}
}
