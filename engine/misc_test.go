package engine

import (
	"encoding/json"
	"slices"
	"strings"
	"testing"

	"gopkg.in/yaml.v3"
)

func TestDates(t *testing.T) {
	if _, err := ParseDate("2026-02-30"); err == nil {
		t.Error("invalid date must fail")
	}
	defer func() {
		if recover() == nil {
			t.Error("MustDate must panic on bad input")
		}
	}()
	d := MustDate("2024-01-31")
	for _, tc := range []struct{ got, want string }{
		{d.AddMonths(1).String(), "2024-02-29"},
		{d.AddMonths(-2).String(), "2023-11-30"},
		{d.LastOfMonth().String(), "2024-01-31"},
		{MustDate("2024-02-10").LastOfMonth().String(), "2024-02-29"},
		{Date{}.String(), ""},
	} {
		if tc.got != tc.want {
			t.Errorf("got %s, want %s", tc.got, tc.want)
		}
	}
	if MustDate("1990-06-15").YearsBetween(MustDate("2026-06-14")) != 35 || MustDate("1990-06-15").YearsBetween(MustDate("2026-06-15")) != 36 {
		t.Error("YearsBetween")
	}
	if !d.Equal(MustDate("2024-01-31")) || d.Compare(d.AddDays(1)) != -1 {
		t.Error("Equal/Compare")
	}
	y, _ := yaml.Marshal(struct{ D Date }{d})
	j, _ := json.Marshal(struct{ D Date }{d})
	if !strings.Contains(string(y), "2024-01-31") || !strings.Contains(string(j), `"2024-01-31"`) {
		t.Errorf("marshal: %s %s", y, j)
	}
	var bad struct{ D Date }
	if err := yaml.Unmarshal([]byte("d: 31.01.2024"), &bad); err == nil {
		t.Error("bad YAML date must fail")
	}
	MustDate("nope")
}

// TestWindowExit checks the exit date of every window kind against a day-by-day search.
func TestWindowExit(t *testing.T) {
	e := &evalCtx{}
	for _, w := range []*Window{{Kind: "rolling_days", N: 90}, {Kind: "rolling_months", N: 1}, {Kind: "rolling_months", N: 24}, {Kind: "calendar_months", N: 6}} {
		for x := MustDate("2023-12-25"); x.Before(MustDate("2024-04-05")); x = x.AddDays(1) {
			got := windowExit(w, x)
			want := x
			for {
				s, _ := e.windowRange(w, want)
				if !s.contains(x) {
					break
				}
				want = want.AddDays(1)
			}
			if !got.Equal(want) {
				t.Fatalf("%s: exit(%s) = %s, want %s", w, x, got, want)
			}
		}
	}
	if windowMoving(nil) || windowMoving(&Window{Kind: "lifetime"}) {
		t.Error("lifetime is not moving")
	}
	if _, ok := e.windowRange(&Window{Kind: "bogus"}, MustDate("2024-01-01")); ok {
		t.Error("unknown window kind is unknown")
	}
	if s, ok := e.windowRange(nil, MustDate("2024-01-01")); !ok || !s.open {
		t.Error("no window is lifetime")
	}
}

func TestYAMLForms(t *testing.T) {
	var r Rule
	err := yaml.Unmarshal([]byte(`
window: { lifetime: true }
stages:
  - { when: { not: { any: [{ met: x }, { expires_within: { days: 3 } }] } }, status: current, messageKey: k }
`), &r)
	if err != nil {
		t.Fatal(err)
	}
	if r.Window.String() != "lifetime" || r.Stages[0].When.Not.List[1].Days != 3 {
		t.Errorf("parsed %+v", r)
	}
	var ops []string
	r.Stages[0].When.Walk(func(c *Condition) { ops = append(ops, c.Op) })
	if !slices.Equal(ops, []string{"not", "any", "met", "expires_within"}) {
		t.Errorf("walk: %v", ops)
	}
	for _, bad := range []string{
		"window: 3",
		"window: { rolling_days: x }",
		"stages: [{ when: [a], status: s, messageKey: k }]",
		"stages: [{ when: { bogus: 1 }, status: s, messageKey: k }]",
		"stages: [{ when: { expires_within: 3 }, status: s, messageKey: k }]",
		"stages: [{ when: { all: 3 }, status: s, messageKey: k }]",
		"stages: [{ when: { not: [1] }, status: s, messageKey: k }]",
		"stages: [{ when: { met_within: 3 }, status: s, messageKey: k }]",
		"stages: [{ when: { holds: 3 }, status: s, messageKey: k }]",
		"stages: [{ when: always, status: s, messageKey: k, params: { a: [1] } }]",
		"filter: { classes: 3 }",
	} {
		var x Rule
		if err := yaml.Unmarshal([]byte(bad), &x); err == nil {
			t.Errorf("want an error for %q", bad)
		}
	}
	var f Filter
	if err := yaml.Unmarshal([]byte("{ roles: [pic], simulator: only }"), &f); err != nil || !slices.Equal(f.Keys(), []string{"roles", "simulator"}) {
		t.Errorf("filter keys %v %v", f.Keys(), err)
	}
	var nilFilter *Filter
	if nilFilter.Keys() != nil || nilFilter.Has("roles") {
		t.Error("nil filter has no keys")
	}
	merged := MergeFilter(&f, &Filter{Simulator: "include", keys: map[string]bool{"simulator": true}})
	if merged.Simulator != "include" || !slices.Equal(merged.Roles, []string{"pic"}) || MergeFilter(nil, &f) != &f || MergeFilter(&f, nil) != &f {
		t.Errorf("merge %+v", merged)
	}
}

func TestRecordAccessors(t *testing.T) {
	c := testCatalogue(t)
	var fl Flags
	for _, name := range c.Vocabulary.FlightFlags {
		if _, ok := fl.Get(name); !ok {
			t.Errorf("flag %s not readable", name)
		}
	}
	if _, ok := fl.Get("nope"); ok {
		t.Error("unknown flag")
	}
	var m Minutes
	for _, name := range c.Vocabulary.MinuteFields {
		if _, ok := m.Minute(name); !ok {
			t.Errorf("minute field %s not readable", name)
		}
	}
	if _, ok := m.Minute("nope"); ok {
		t.Error("unknown minute field")
	}
	m = Minutes{PIC: 1, Dual: 2, SPIC: 3, PICUS: 4, SIC: 5, DualGiven: 6, Examiner: 7}
	for i, role := range c.Vocabulary.Roles {
		if roleMinutes(m, role) != i+1 {
			t.Errorf("role %s", role)
		}
	}
	if roleMinutes(m, "nope") != 0 {
		t.Error("unknown role")
	}
	for alias, kind := range map[string]string{"ppl ( a )": "PPL_A", "FI(S)": "INSTRUCTOR", "whatever": ""} {
		if got := c.Vocabulary.ClassifyLicence(alias, "EASA", ""); got != kind {
			t.Errorf("classify %q = %q, want %q", alias, got, kind)
		}
	}
	if c.Vocabulary.ClassifyLicence("PPL", "daec", "") != "UL" || c.Vocabulary.ClassifyLicence("x", "", " spl ") != "SPL" || c.Vocabulary.Category("NOPE") != "" {
		t.Error("classify authority/kind")
	}
}

func TestMetricsAllImplemented(t *testing.T) {
	c := testCatalogue(t)
	impl := Implemented()
	for section, names := range map[string][]string{
		"metrics": keysOf(c.Vocabulary.Metrics), "filters": keysOf(c.Vocabulary.Filters), "windows": keysOf(c.Vocabulary.Windows),
		"combinators": keysOf(c.Vocabulary.Combinators), "stage_conditions": keysOf(c.Vocabulary.StageConditions),
		"rule_events": keysOf(c.Vocabulary.RuleEvents), "param_sources": keysOf(c.Vocabulary.ParamSources),
		"hatches": keysOf(c.Vocabulary.Hatches), "subjects": keysOf(c.Vocabulary.Subjects),
	} {
		for _, n := range names {
			if !slices.Contains(impl[section], n) {
				t.Errorf("%s: %s is declared but not implemented", section, n)
			}
		}
		if len(impl[section]) != len(names) {
			t.Errorf("%s: engine implements %v, vocabulary declares %v", section, impl[section], names)
		}
	}
	two, one, three := 2, 1, 3
	yes := true
	f := &Flight{
		Minutes:  Minutes{Total: 1, PIC: 2, Dual: 3, SPIC: 4, PICUS: 5, SIC: 6, DualGiven: 7, Examiner: 8, MultiPilot: 9, Night: 10, IFR: 11, ActualInstrument: 12, SimulatedInstrument: 13, CrossCountry: 14},
		Takeoffs: DayNight{Day: 1, Night: 2}, Landings: DayNight{Day: 3, Night: 4},
		FullStopLandings: &two, FullStopNightLandings: &one, Launches: 5, Approaches: 6, Holds: 7,
		InterceptAndTrack: &yes, CruiseMinutes: &two, MountainLandings: &three, Flags: Flags{TowFlight: true},
		NightPeriodTakeoffs: &one,
	}
	want := map[string]float64{
		"flights": 1, "route_sectors": 0, "minutes.total": 1, "minutes.pic": 2, "minutes.dual": 3, "minutes.spic": 4,
		"minutes.dualGiven": 7, "minutes.ifr": 11,
		"minutes.picOrDualOrSpic": 9, "minutes.dualOrSpic": 7, "minutes.picOrSpic": 6, "minutes.instructorOrExaminer": 15,
		"takeoffs.total": 3, "takeoffs.night": 2, "landings.total": 7, "landings.night": 4, "mountain_landings": 3,
		"takeoffs_and_landings": 3, "full_stop_landings": 2, "full_stop_night_landings": 1, "launches": 5, "approaches": 6,
		"holds": 7, "intercept_and_track": 1, "not_recorded": 0, "training_flights": 1, "longest_training_flight_minutes": 1, "tows": 1,
		"takeoffs.night_period": 1, "longest_flight_minutes": 1,
	}
	for name, fn := range flightMetricFuncs {
		got := fn(f)
		if got.known == (name == "not_recorded") || got.v != want[name] {
			t.Errorf("%s = %+v, want %v", name, got, want[name])
		}
	}
	empty := &Flight{}
	for _, name := range []string{"route_sectors", "full_stop_landings", "full_stop_night_landings", "intercept_and_track", "mountain_landings"} {
		if flightMetricFuncs[name](empty).known {
			t.Errorf("%s of an empty flight must be unknown", name)
		}
	}
	if !flightMetricFuncs["takeoffs.night_period"](empty).known || flightMetricFuncs["takeoffs.night_period"](&Flight{Takeoffs: DayNight{Night: 1}}).known {
		t.Error("takeoffs.night_period: 0 without night take-offs, unknown with them and no period count")
	}
	if flightMetricFuncs["tows"](empty).v != 0 || flightMetricFuncs["tows"](&Flight{TowedGliders: &two, Flags: Flags{TowFlight: true}}).v != 2 {
		t.Error("tows")
	}
}

func keysOf[V any](m map[string]V) []string {
	var out []string
	for k := range m {
		out = append(out, k)
	}
	return out
}

func TestEventFilters(t *testing.T) {
	c := testCatalogue(t)
	e := &evalCtx{cat: c, v: c.Vocabulary, p: prepare(&Record{}, c.Vocabulary), subj: subjectRef{Subject: Subject{Class: "SEP_LAND", ULKind: "", Detail: "THREE_AXIS"}, typeDesignator: "C172", variant: "diesel"}}
	ev := func(s string) *Event {
		var x Event
		if err := yaml.Unmarshal([]byte(s), &x); err != nil {
			t.Fatal(err)
		}
		return &x
	}
	filter := func(s string) *resolvedFilter {
		var f Filter
		if err := yaml.Unmarshal([]byte(s), &f); err != nil {
			t.Fatal(err)
		}
		return e.resolveFilter(&f)
	}
	for _, tc := range []struct {
		filter, event string
		want          match
	}{
		{"{ simulator: only }", "{ kind: ipc, isSimulator: true }", matchYes},
		{"{ simulator: only }", "{ kind: ipc }", matchNo},
		{"{ classes: [$subject] }", "{ kind: ipc }", matchUnknown},
		{"{ classes: [$subject] }", "{ kind: ipc, class: SEP_LAND }", matchYes},
		{"{ excludeClasses: [SEP_LAND] }", "{ kind: ipc, class: SEP_LAND }", matchNo},
		{"{ categories: [aeroplane] }", "{ kind: ipc }", matchUnknown},
		{"{ categories: [aeroplane] }", "{ kind: ipc, class: GLIDER }", matchNo},
		{"{ ulKinds: [$subject] }", "{ kind: ipc, class: SEP_LAND }", matchNo},
		{"{ ulKinds: [$subject] }", "{ kind: ipc, class: ULTRALIGHT }", matchUnknown},
		{"{ ulKinds: [$subject] }", "{ kind: ipc, class: ULTRALIGHT, ulKind: THREE_AXIS }", matchYes},
		{"{ typeDesignators: [$subject] }", "{ kind: ipc, typeDesignator: C172 }", matchYes},
		{"{ variants: [$subject] }", "{ kind: ipc }", matchUnknown},
		{"{ any: [{ classes: [GLIDER] }, { categories: [aeroplane] }] }", "{ kind: ipc, class: SEP_LAND }", matchYes},
		{"{ any: [{ classes: [GLIDER] }, { categories: [sailplane] }] }", "{ kind: ipc, class: SEP_LAND }", matchNo},
		{"{ any: [{ classes: [GLIDER] }, { categories: [sailplane] }] }", "{ kind: ipc }", matchUnknown},
	} {
		if got := e.matchEvent(filter(tc.filter), ev(tc.event)).res; got != tc.want {
			t.Errorf("%s on %s = %v, want %v", tc.filter, tc.event, got, tc.want)
		}
	}
	fl := &Flight{Class: "SEP_LAND", TowKind: "glider"}
	for _, tc := range []struct {
		filter string
		want   match
	}{
		{"{ towKinds: [glider] }", matchNo},
		{"{ fstdTypes: [FFS] }", matchYes},
		{"{ classes: [GLIDER], ulCredit: [{ class: GLIDER, ulKinds: [SAILPLANE] }] }", matchNo},
	} {
		if got := e.matchFlight(filter(tc.filter), fl).res; got != tc.want {
			t.Errorf("%s = %v, want %v", tc.filter, got, tc.want)
		}
	}
	ul := &Flight{Class: "ULTRALIGHT", ULKind: "THREE_AXIS"}
	if e.matchFlight(filter("{ classes: [SEP_LAND], ulCredit: [{ class: SEP_LAND, ulKinds: [GYROPLANE] }] }"), ul).res != matchNo {
		t.Error("UL kind not credited")
	}
}

// TestUnresolvedSubjectFilter: a $subject the subject cannot supply leaves a recorded item
// unknown, never a silent non-match; items the filter excludes on other grounds stay out.
func TestUnresolvedSubjectFilter(t *testing.T) {
	c := testCatalogue(t)
	e := &evalCtx{cat: c, v: c.Vocabulary, p: prepare(&Record{}, c.Vocabulary), subj: subjectRef{Subject: Subject{Kind: "privilege"}}}
	for _, tc := range []struct {
		filter string
		fl     Flight
		want   match
	}{
		{"{ ulKinds: [$subject] }", Flight{Class: "ULTRALIGHT", ULKind: "THREE_AXIS"}, matchUnknown},
		{"{ ulKinds: [$subject] }", Flight{Class: "SEP_LAND"}, matchNo},
		{"{ ulKinds: [$subject] }", Flight{Class: "ULTRALIGHT", ULKind: "THREE_AXIS", IsSimulator: true}, matchNo},
		{"{ classes: [$subject] }", Flight{Class: "SEP_LAND"}, matchUnknown},
		{"{ typeDesignators: [$subject] }", Flight{Class: "SEP_LAND", TypeDesignator: "C172"}, matchUnknown},
		{"{ variants: [$subject] }", Flight{Class: "SEP_LAND", Variant: "diesel"}, matchUnknown},
		{"{ launchMethods: [$subject] }", Flight{Class: "GLIDER", LaunchMethod: "winch"}, matchUnknown},
	} {
		var f Filter
		if err := yaml.Unmarshal([]byte(tc.filter), &f); err != nil {
			t.Fatal(err)
		}
		if got := e.matchFlight(e.resolveFilter(&f), &tc.fl).res; got != tc.want {
			t.Errorf("%s on %+v = %v, want %v", tc.filter, tc.fl, got, tc.want)
		}
	}
	var f Filter
	if err := yaml.Unmarshal([]byte("{ ulKinds: [$subject], classes: [$subject] }"), &f); err != nil {
		t.Fatal(err)
	}
	if got := e.matchEvent(e.resolveFilter(&f), &Event{Kind: "proficiency_check", Class: "ULTRALIGHT", ULKind: "THREE_AXIS"}).res; got != matchUnknown {
		t.Errorf("event = %v, want unknown", got)
	}
}

func TestWalkAndDiffExpect(t *testing.T) {
	var r Rule
	if err := yaml.Unmarshal([]byte(`
requirements:
  any_of:
    - { id: check, metric: events, min: 1, unit: check, nameKey: requirement.proficiency_check }
    - { id: experience, all_of: [{ id: a, metric: flights, min: 1, unit: flights, nameKey: requirement.flights }, { id: b, metric: flights, min: 2, unit: flights, nameKey: requirement.flights }] }
`), &r); err != nil {
		t.Fatal(err)
	}
	n := 0
	r.Requirements.Walk(func(*Node) { n++ })
	if n != 5 {
		t.Errorf("walk visited %d nodes", n)
	}
	var none *Node
	none.Walk(func(*Node) { t.Error("nil walk") })

	f, tr, fa, one := 1.0, true, false, "x"
	d := MustDate("2020-01-01")
	x := Expect{Subject: Subject{Kind: "rating", ID: "r"}, Status: "current", MessageKey: "x", MessageParams: map[string]any{"days": 1, "date": d.t}, ExpiresOn: &d, WindowOpensAt: &d, ValidUntil: &d,
		Requirements: map[string]ExpectedRow{
			"total_time":        {Current: &f, Required: &f, Met: &tr, Tracked: &fa, ValidUntil: &d, LastDate: &d, MessageKey: &one, RemedyKey: &one, RemedyParams: map[string]any{"missing": 1, "nope": 2}},
			"nope":              {},
			"proficiency_check": {Absent: true},
		}}
	ev := Evaluation{
		Subject: Subject{Kind: "rating", ID: "r", Class: "SEP_LAND"}, Status: "expiring", MessageKey: "rating.expiring",
		MessageParams: map[string]any{"days": 2, "date": "2021-01-01"},
		Requirements: []RequirementResult{
			{ID: "total_time", Current: 2, Required: 3, Met: false, Tracked: true, RemedyParams: map[string]any{"missing": 3}},
			{ID: "proficiency_check"},
		},
	}
	joined := strings.Join(DiffExpect(x, ev), "\n")
	for _, want := range []string{"status: want current", "messageParams.days", "messageParams.date: want 2020-01-01", "expiresOn", "windowOpensAt", "validUntil: want 2020-01-01, got (none)", "total_time: current", "required", "met", "tracked", "lastDate", "messageKey", "remedyKey", "remedyParams.missing", "remedyParams.nope: want 2, got none", "requirement nope: want a row", "proficiency_check: want absent"} {
		if !strings.Contains(joined, want) {
			t.Errorf("diffs miss %q:\n%s", want, joined)
		}
	}
	if !SubjectMatches(Subject{Kind: "rating"}, ev.Subject) || SubjectMatches(Subject{Kind: "rating", ID: "gone"}, ev.Subject) {
		t.Error("subject matching")
	}
}

// TestSupersedes checks that Evaluate drops a superseded rule's evaluation of a subject the
// superseding rule evaluates, and keeps it for subjects only the superseded rule selects.
func TestSupersedes(t *testing.T) {
	c := testCatalogue(t, `
id: test.x.fallback
applies_to: { subject: rating }
stages: [{ when: always, status: unknown, messageKey: rating.no_expiry_date }]
`, `
id: test.x.gyroplane
supersedes: [test.x.fallback]
applies_to: { subject: rating, classes: [GYROPLANE] }
stages: [{ when: always, status: current, messageKey: rating.valid_until }]
`)
	rec := record(t, `
licences:
  - { id: l-caa, authority: CAA, type: GPL }
ratings:
  - { id: r-gyro, licenceId: l-caa, class: GYROPLANE, expires: 2027-01-31 }
  - { id: r-sep, licenceId: l-caa, class: SEP_LAND, expires: 2027-01-31 }
`)
	byRule := map[string][]string{}
	for _, ev := range Evaluate(c, rec, MustDate("2026-09-15")) {
		byRule[ev.RuleID] = append(byRule[ev.RuleID], ev.Subject.ID)
	}
	if got := byRule["test.x.fallback"]; !slices.Equal(got, []string{"r-sep"}) {
		t.Errorf("the fallback keeps only the rating no superseding rule evaluates: %v", got)
	}
	if got := byRule["test.x.gyroplane"]; !slices.Equal(got, []string{"r-gyro"}) {
		t.Errorf("the gyroplane rule evaluates the gyroplane rating: %v", got)
	}
}
