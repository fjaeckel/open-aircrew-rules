package credentials

import (
	"fmt"
	"regexp"
	"slices"
	"strings"

	"github.com/fjaeckel/open-aircrew-rules/engine"
)

// Composite answers, for one credential a record holds, "may it be exercised on the date?":
// the conjunction of its own evaluations and its requirements (DESIGN.md section 5).
type Composite struct {
	Credential string         `yaml:"credential" json:"credential"`
	Subject    engine.Subject `yaml:"subject" json:"subject"`
	Status     string         `yaml:"status" json:"status"`
	// DecidedBy is the member that sets Status (the first in file order with that status).
	DecidedBy *Member  `yaml:"decidedBy,omitempty" json:"decidedBy,omitempty"`
	Members   []Member `yaml:"members" json:"members"`
	// Limitations are evaluations of part of the privileges (passengers, a launch method, a
	// variant, a training programme) and entries with limits: reported, never deciding.
	Limitations []Member `yaml:"limitations,omitempty" json:"limitations,omitempty"`
	// Levels answers the question per named level of the credential (level:), each from
	// that level's members and those without a level.
	Levels []Level `yaml:"levels,omitempty" json:"levels,omitempty"`
}

// Level is a composite restricted to one named level of a credential's privileges.
type Level struct {
	Level     string  `yaml:"level" json:"level"`
	Status    string  `yaml:"status" json:"status"`
	DecidedBy *Member `yaml:"decidedBy,omitempty" json:"decidedBy,omitempty"`
}

// Member is one part of a composite: an own evaluation or a requirement group.
type Member struct {
	// Evaluation is the evaluation id in the credential file.
	Evaluation string `yaml:"evaluation" json:"evaluation"`
	// Kind is evaluation, requires_all or requires_any.
	Kind   string `yaml:"kind" json:"kind"`
	Status string `yaml:"status" json:"status"`
	// RuleID, Subject and MessageKey are those of the evaluation result (Kind evaluation).
	RuleID     string          `yaml:"ruleId,omitempty" json:"ruleId,omitempty"`
	Subject    *engine.Subject `yaml:"subject,omitempty" json:"subject,omitempty"`
	MessageKey string          `yaml:"messageKey,omitempty" json:"messageKey,omitempty"`
	// Credential is the required credential that decides a requirement group; Reason is
	// not_held when the record holds none that fits.
	Credential string `yaml:"credential,omitempty" json:"credential,omitempty"`
	Reason     string `yaml:"reason,omitempty" json:"reason,omitempty"`
	// Level is the entry's level; Scope the privileges a limitation (limits:) concerns.
	Level string `yaml:"level,omitempty" json:"level,omitempty"`
	Scope string `yaml:"scope,omitempty" json:"scope,omitempty"`
}

// credentialOf returns the credential id of a requirement id (<credential id>[@<level>]).
func credentialOf(id string) string {
	c, _, _ := strings.Cut(id, "@")
	return c
}

// levels returns the distinct levels of a credential's entries, in file order.
func levels(c *Credential) []string {
	var out []string
	for _, e := range c.Evaluations {
		if e.Level != "" && !slices.Contains(out, e.Level) {
			out = append(out, e.Level)
		}
	}
	return out
}

// Result is everything a record evaluates to: one evaluation per compiled rule and subject,
// and one composite per held credential.
type Result struct {
	Evaluations []engine.Evaluation `yaml:"evaluations" json:"evaluations"`
	Credentials []Composite         `yaml:"credentials" json:"credentials"`
}

// Evaluate evaluates a record on asOf: every compiled rule, then every held credential.
func (cat *Catalogue) Evaluate(rec *engine.Record, asOf engine.Date) Result {
	evs := engine.Evaluate(cat.Engine, rec, asOf)
	return Result{Evaluations: evs, Credentials: cat.Composites(rec, asOf, evs)}
}

// isRequirement reports whether an evaluation only names credentials this one needs.
func (e *Evaluation) isRequirement() bool { return len(e.RequiresAll)+len(e.RequiresAny) > 0 }

func (cat *Catalogue) checkRequirement(key string, e *Evaluation) {
	if e.Uses != "" || e.Source != "" || !empty(&e.PassesIf) || !empty(&e.Outcomes) {
		cat.Errors = append(cat.Errors, fmt.Sprintf("%s: an evaluation with requires_all or requires_any has only id, asks and those", key))
	}
	for _, ref := range append(slices.Clone(e.RequiresAll), e.RequiresAny...) {
		id, level, _ := strings.Cut(ref, "@")
		c, ok := cat.byID[id]
		if !ok {
			cat.Errors = append(cat.Errors, fmt.Sprintf("%s: requires %s, which is no credential", key, id))
		} else if level != "" && !slices.Contains(levels(c), level) {
			cat.Errors = append(cat.Errors, fmt.Sprintf("%s: requires %s, but %s has no level %q", key, ref, id, level))
		}
		if e.Owner != nil && id == e.Owner.ID {
			cat.Errors = append(cat.Errors, fmt.Sprintf("%s: requires the credential itself", key))
		}
	}
}

// requires lists the credentials a credential names in requires_all or requires_any.
func requires(c *Credential) []string {
	var out []string
	for _, e := range c.Evaluations {
		for _, id := range append(slices.Clone(e.RequiresAll), e.RequiresAny...) {
			out = append(out, credentialOf(id))
		}
	}
	return out
}

// checkLevelsAndLimits reports malformed level names and limitation scopes the vocabulary
// lacks.
func (cat *Catalogue) checkLevelsAndLimits(c *Credential) {
	for _, e := range c.Evaluations {
		where := c.ID + "#" + e.ID
		if e.Level != "" && !levelName.MatchString(e.Level) {
			cat.Errors = append(cat.Errors, fmt.Sprintf("%s: level %q is not a lower-case name", where, e.Level))
		}
		if _, ok := cat.Vocab.LimitationScopes[e.Limits]; e.Limits != "" && !ok {
			cat.Errors = append(cat.Errors, fmt.Sprintf("%s: limits %q is not a limitation scope of the vocabulary", where, e.Limits))
		}
		if e.Limits != "" && e.Level != "" {
			cat.Errors = append(cat.Errors, where+": an entry has a level or limits, not both")
		}
	}
}

var levelName = regexp.MustCompile(`^[a-z][a-z0-9_]*$`)

// checkRequirementCycles reports credentials that require themselves through others.
func (cat *Catalogue) checkRequirementCycles() {
	const (
		open = 1
		done = 2
	)
	state := map[string]int{}
	var path []string
	var visit func(id string)
	visit = func(id string) {
		switch state[id] {
		case open:
			i := slices.Index(path, id)
			cat.Errors = append(cat.Errors, fmt.Sprintf("requirement cycle: %s -> %s", strings.Join(path[i:], " -> "), id))
			return
		case done:
			return
		}
		c, ok := cat.byID[id]
		if !ok {
			return
		}
		state[id] = open
		path = append(path, id)
		for _, r := range requires(c) {
			if r == id {
				continue // reported by checkRequirement
			}
			visit(r)
		}
		path = path[:len(path)-1]
		state[id] = done
	}
	for _, c := range cat.Credentials {
		visit(c.ID)
	}
}

// rank orders statuses for a composite: 0 current, 1 expiring, 2 unknown, 3 expired or
// lapsed; not_applicable is -1 and never decides.
func rank(status string) int {
	switch status {
	case "current":
		return 0
	case "expiring":
		return 1
	case "expired", "lapsed":
		return 3
	case "not_applicable":
		return -1
	}
	return 2
}

// composer evaluates composites for one record.
type composer struct {
	cat  *Catalogue
	rec  *engine.Record
	asOf engine.Date
	evs  map[string][]engine.Evaluation
	memo map[string]*Composite
	held map[string][]engine.Subject
}

// Composites returns one composite per credential item the record holds, by credential id
// and then in record order. evs are the record's evaluations on asOf.
func (cat *Catalogue) Composites(rec *engine.Record, asOf engine.Date, evs []engine.Evaluation) []Composite {
	cm := &composer{cat: cat, rec: rec, asOf: asOf, evs: map[string][]engine.Evaluation{}, memo: map[string]*Composite{}, held: map[string][]engine.Subject{}}
	for _, ev := range evs {
		cm.evs[ev.RuleID] = append(cm.evs[ev.RuleID], ev)
	}
	ids := make([]string, 0, len(cat.Credentials))
	for _, c := range cat.Credentials {
		ids = append(ids, c.ID)
	}
	slices.Sort(ids)
	out := []Composite{}
	for _, id := range ids {
		for _, s := range cm.holdings(id) {
			out = append(out, *cm.composite(id, s))
		}
	}
	return out
}

func (cm *composer) holdings(id string) []engine.Subject {
	if h, ok := cm.held[id]; ok {
		return h
	}
	h := engine.Holdings(cm.cat.Engine, cm.rec, cm.cat.held[id], cm.asOf)
	cm.held[id] = h
	return h
}

func itemKey(kind, id string) string { return kind + ":" + id }

// related returns the record items (kind:id) an evaluation subject is about, with the
// licences they are held on.
func related(rec *engine.Record, s engine.Subject) []string {
	var out []string
	rating := func(id string) {
		for _, r := range rec.Ratings {
			if r.ID == id {
				out = append(out, itemKey("rating", r.ID), itemKey("licence", r.LicenceID))
			}
		}
	}
	switch s.Kind {
	case "licence", "flight_review", "training":
		if s.ID != "" {
			out = append(out, itemKey("licence", s.ID))
		}
	case "rating", "launch_method":
		rating(s.ID)
	case "type":
		for _, v := range rec.Variants {
			if v.ID == s.ID {
				rating(v.RatingID)
			}
		}
	case "passengers":
		out = append(out, itemKey("licence", s.ID))
		for _, r := range rec.Ratings {
			if r.LicenceID == s.ID && r.Class == s.Class && r.ULKind == s.ULKind {
				out = append(out, itemKey("rating", r.ID))
			}
		}
	case "privilege":
		out = append(out, itemKey("privilege", s.ID))
		for _, p := range rec.Privileges {
			if p.ID == s.ID {
				out = append(out, itemKey("licence", p.LicenceID))
			}
		}
	case "credential":
		out = append(out, itemKey("credential", s.ID))
	}
	return out
}

// licenceOf returns the licence a held item is on ("" for a certificate such as a medical).
func licenceOf(rec *engine.Record, s engine.Subject) string {
	for _, k := range related(rec, s) {
		if id, ok := strings.CutPrefix(k, "licence:"); ok {
			return id
		}
	}
	return ""
}

// limitation reports whether an evaluation subject is part of the privileges only.
func limitation(kind string) bool {
	return kind == "passengers" || kind == "launch_method" || kind == "type" || kind == "training"
}

func (cm *composer) composite(id string, item engine.Subject) *Composite {
	key := id + "|" + item.Kind + ":" + item.ID
	if c, ok := cm.memo[key]; ok {
		return c
	}
	out := &Composite{Credential: id, Subject: item, Status: "unknown", Members: []Member{}}
	cm.memo[key] = out // a cycle (refused by Load) would read unknown
	c := cm.cat.byID[id]
	self := itemKey(item.Kind, item.ID)
	add := func(e *Evaluation, m Member) {
		m.Level = e.Level
		if e.Limits != "" {
			m.Scope = e.Limits
			out.Limitations = append(out.Limitations, m)
			return
		}
		out.Members = append(out.Members, m)
	}
	for _, e := range c.Evaluations {
		if e.isRequirement() {
			if m, ok := cm.group(e.ID, "requires_all", e.RequiresAll, item); ok {
				add(e, m)
			}
			if m, ok := cm.group(e.ID, "requires_any", e.RequiresAny, item); ok {
				add(e, m)
			}
			continue
		}
		rule := cm.cat.ByEvaluation[id+"#"+e.ID]
		var worst *engine.Evaluation
		for i, ev := range cm.evs[rule] {
			if slices.Contains(related(cm.rec, ev.Subject), self) && (worst == nil || rank(ev.Status) > rank(worst.Status)) {
				worst = &cm.evs[rule][i]
			}
		}
		if worst == nil {
			continue
		}
		subj := worst.Subject
		m := Member{Evaluation: e.ID, Kind: "evaluation", Status: worst.Status, RuleID: rule, Subject: &subj, MessageKey: worst.MessageKey}
		if limitation(subj.Kind) {
			m.Level = e.Level
			out.Limitations = append(out.Limitations, m)
			continue
		}
		add(e, m)
	}
	out.Status, out.DecidedBy = decide(out.Members, false)
	for _, l := range levels(c) {
		st, d := decide(levelMembers(out.Members, l), false)
		out.Levels = append(out.Levels, Level{Level: l, Status: st, DecidedBy: d})
	}
	return out
}

// levelMembers returns the members of one level and those without a level.
func levelMembers(ms []Member, level string) []Member {
	var out []Member
	for _, m := range ms {
		if m.Level == "" || m.Level == level {
			out = append(out, m)
		}
	}
	return out
}

// status returns a composite's status, or that of one of its levels.
func (c *Composite) status(level string) string {
	if level == "" {
		return c.Status
	}
	for _, l := range c.Levels {
		if l.Level == level {
			return l.Status
		}
	}
	return "unknown"
}

// decide returns the worst status of the members (the best with best) and the first member
// that has it; not_applicable when no member decides.
func decide(ms []Member, best bool) (string, *Member) {
	var pick *Member
	for i := range ms {
		r := rank(ms[i].Status)
		if r < 0 {
			continue
		}
		if pick == nil || (!best && r > rank(pick.Status)) || (best && r < rank(pick.Status)) {
			pick = &ms[i]
		}
	}
	if pick == nil {
		return "not_applicable", nil
	}
	d := *pick
	return d.Status, &d
}

// group evaluates requires_all (every credential) or requires_any (at least one) for a held
// item. A required credential counts by the best composite among the items the record holds
// of it, narrowed to the same licence when both are on one. requires_all: the worst of the
// listed credentials, one not held being unknown. requires_any: the best of the listed
// credentials the record holds; unknown when it holds none of them.
func (cm *composer) group(eval, kind string, ids []string, item engine.Subject) (Member, bool) {
	if len(ids) == 0 {
		return Member{}, false
	}
	lic := licenceOf(cm.rec, item)
	var parts []Member
	for _, ref := range ids {
		id, level, _ := strings.Cut(ref, "@")
		var cands []Member
		for _, h := range cm.holdings(id) {
			if l := licenceOf(cm.rec, h); lic != "" && l != "" && l != lic {
				continue
			}
			st := cm.composite(id, h).status(level)
			if st == "not_applicable" {
				st = "current"
			}
			cands = append(cands, Member{Status: st, Credential: ref})
		}
		if len(cands) == 0 {
			parts = append(parts, Member{Status: "unknown", Credential: ref, Reason: "not_held"})
			continue
		}
		_, b := decide(cands, true)
		parts = append(parts, *b)
	}
	if kind == "requires_any" {
		if held := slices.DeleteFunc(slices.Clone(parts), func(m Member) bool { return m.Reason == "not_held" }); len(held) > 0 {
			parts = held
		}
	}
	_, d := decide(parts, kind == "requires_any")
	return Member{Evaluation: eval, Kind: kind, Status: d.Status, Credential: d.Credential, Reason: d.Reason}, true
}
