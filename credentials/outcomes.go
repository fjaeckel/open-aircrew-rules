package credentials

import (
	"fmt"

	"gopkg.in/yaml.v3"

	"github.com/fjaeckel/open-aircrew-rules/engine"
)

// outcomeSpec is the mapping form of `outcomes`.
type outcomeSpec struct {
	Preset         string
	Special        []engine.Stage
	Messages       map[string]string
	Statuses       map[string]string
	Omit           []string
	NoticeDays     int
	NoticeMessage  string
	DateOfBirth    string
	hasNotice      bool
	explicitStages []engine.Stage
}

func cond(op string) engine.Condition { return engine.Condition{Op: op} }

func param(src string) engine.ParamSource { return engine.ParamSource{Source: src} }

func needed(id string) map[string]engine.ParamSource {
	return map[string]engine.ParamSource{"needed": {Source: "needed", Ref: id}}
}

// outcomes compiles the outcomes node into stages.
func (cp *compiler) outcomes(n *yaml.Node, r *engine.Rule) []engine.Stage {
	n = deref(n)
	if empty(n) {
		cp.errs = append(cp.errs, cp.where+": no outcomes")
		return nil
	}
	spec := outcomeSpec{}
	switch n.Kind {
	case yaml.ScalarNode:
		spec.Preset = n.Value
	case yaml.SequenceNode:
		return cp.stages(n, "outcomes")
	default:
		for _, p := range pairs(n) {
			k, v := p[0].Value, p[1]
			switch k {
			case "preset":
				spec.Preset = v.Value
			case "special":
				spec.Special = cp.stages(v, "outcomes.special")
			case "messages":
				_ = v.Decode(&spec.Messages)
			case "statuses":
				_ = v.Decode(&spec.Statuses)
			case "omit":
				_ = v.Decode(&spec.Omit)
			case "date_of_birth":
				spec.DateOfBirth = v.Value
			case "expiring_notice":
				spec.hasNotice = true
				for _, q := range pairs(v) {
					switch q[0].Value {
					case "days":
						_ = q[1].Decode(&spec.NoticeDays)
					case "message":
						spec.NoticeMessage = q[1].Value
					case "ref":
						cp.ref(q[1], "outcomes.expiring_notice")
					default:
						cp.fail(q[0], "expiring_notice: unknown key %q", q[0].Value)
					}
				}
			default:
				cp.fail(p[0], "outcomes: unknown key %q", k)
			}
		}
	}
	if _, ok := cp.v.OutcomePresets[spec.Preset]; !ok {
		cp.fail(n, "outcomes: preset %q is not in the vocabulary", spec.Preset)
		return nil
	}
	st := cp.preset(spec, r)
	var out []engine.Stage
	for _, s := range st {
		if containsStr(spec.Omit, s.ID) {
			continue
		}
		base := baseID(s.ID)
		if m, ok := spec.Messages[base]; ok {
			s.MessageKey = m
		}
		if m, ok := spec.Statuses[base]; ok {
			s.Status = m
		}
		out = append(out, s)
	}
	return out
}

func containsStr(l []string, s string) bool {
	for _, x := range l {
		if x == s {
			return true
		}
	}
	return false
}

// baseID maps generated chain stage ids (not_met_<req>) to their base.
func baseID(id string) string {
	if len(id) > len("not_met_") && id[:len("not_met_")] == "not_met_" {
		return "not_met"
	}
	return id
}

func (cp *compiler) preset(s outcomeSpec, r *engine.Rule) []engine.Stage {
	notice := func(def string) (int, string) {
		if !s.hasNotice || s.NoticeDays <= 0 {
			cp.errs = append(cp.errs, fmt.Sprintf("%s: preset %s needs expiring_notice { days, ref }", cp.where, s.Preset))
		}
		m := s.NoticeMessage
		if m == "" {
			m = def
		}
		return s.NoticeDays, m
	}
	date := map[string]engine.ParamSource{"date": param("expiry_date")}
	switch s.Preset {
	case "revalidation":
		days, _ := notice("")
		out := []engine.Stage{
			{ID: "no_expiry", When: cond("no_expiry"), Status: "unknown", MessageKey: "rating.no_expiry_date"},
			{ID: "expired", When: cond("expired"), Status: "expired", MessageKey: "rating.expired"},
		}
		out = append(out, s.Special...)
		return append(out,
			engine.Stage{ID: "window_not_open", When: cond("before_window"), Status: "current", MessageKey: "rating.window_not_open", Params: map[string]engine.ParamSource{"date": param("window_opens_at")}},
			engine.Stage{ID: "met_expiring", When: engine.Condition{Op: "all", List: []*engine.Condition{{Op: "all_met"}, {Op: "expires_within", Days: days}}}, Status: "expiring", MessageKey: "rating.revalidation_expiring_met", Params: map[string]engine.ParamSource{"days": param("days_to_expiry")}},
			engine.Stage{ID: "met", When: cond("all_met"), Status: "current", MessageKey: "rating.revalidation_current"},
			engine.Stage{ID: "not_met", When: cond("always"), Status: "expiring", MessageKey: "rating.revalidation_not_met"},
		)
	case "recency":
		return append(s.Special,
			engine.Stage{ID: "current", When: cond("all_met"), Status: "current", MessageKey: "rating.recency_current"},
			engine.Stage{ID: "lapsed", When: cond("always"), Status: "lapsed", MessageKey: "rating.recency_not_met"},
		)
	case "privilege_recency":
		out := []engine.Stage{{ID: "expired", When: cond("expired"), Status: "expired", MessageKey: "privilege.expired", Params: date}}
		out = append(out, s.Special...)
		return append(out,
			engine.Stage{ID: "current", When: cond("all_met"), Status: "current", MessageKey: "privilege.recency_current", Params: date},
			engine.Stage{ID: "lapsed", When: cond("always"), Status: "lapsed", MessageKey: "privilege.recency_not_met", Params: date},
		)
	case "passengers":
		out := append(s.Special,
			engine.Stage{ID: "current", When: cond("all_met"), Status: "current", MessageKey: "pax.day_current"},
			engine.Stage{ID: "unknown", When: cond("undetermined"), Status: "unknown", MessageKey: "pax.experience_not_logged"},
		)
		return append(out, cp.notMetChain(r, "expired", "pax.not_current")...)
	case "training":
		return append(s.Special,
			engine.Stage{ID: "complete", When: cond("all_met"), Status: "current", MessageKey: "training.all_met"},
			engine.Stage{ID: "unknown", When: cond("undetermined"), Status: "unknown", MessageKey: "training.distance_unknown"},
			engine.Stage{ID: "in_progress", When: cond("always"), Status: "lapsed", MessageKey: "training.in_progress"},
		)
	case "validity":
		days, msg := notice("credential.expiring")
		var out []engine.Stage
		switch s.DateOfBirth {
		case "":
		case "when_no_expiry":
			out = append(out, engine.Stage{ID: "date_of_birth", When: engine.Condition{Op: "all", List: []*engine.Condition{{Op: "missing", Ref: "date_of_birth"}, {Op: "no_expiry"}}}, Status: "unknown", MessageKey: "credential.date_of_birth_required"})
		case "required":
			out = append(out, engine.Stage{ID: "date_of_birth", When: engine.Condition{Op: "missing", Ref: "date_of_birth"}, Status: "unknown", MessageKey: "credential.date_of_birth_required"})
		default:
			cp.errs = append(cp.errs, cp.where+": date_of_birth is when_no_expiry or required")
		}
		out = append(out,
			engine.Stage{ID: "no_expiry", When: cond("no_expiry"), Status: "unknown", MessageKey: "credential.no_expiry_date"},
			engine.Stage{ID: "expired", When: cond("expired"), Status: "expired", MessageKey: "credential.expired", Params: date},
		)
		out = append(out, s.Special...)
		return append(out,
			engine.Stage{ID: "expiring", When: engine.Condition{Op: "expires_within", Days: days}, Status: "expiring", MessageKey: msg, Params: map[string]engine.ParamSource{"days": param("days_to_expiry"), "date": param("expiry_date")}},
			engine.Stage{ID: "valid", When: cond("always"), Status: "current", MessageKey: "credential.valid", Params: date},
		)
	}
	return nil
}

// notMetChain reports the missing count of the last listed unmet requirement: for required
// leaves L1..Ln of the root, stages "unmet Ln" ... "unmet L2", then "always" with L1.
func (cp *compiler) notMetChain(r *engine.Rule, status, msg string) []engine.Stage {
	root := r.Requirements
	var leaves []string
	switch {
	case root == nil:
	case root.IsLeaf():
		leaves = []string{root.ID}
	case root.AllOf != nil:
		for _, c := range root.AllOf {
			if c.Informational {
				continue
			}
			if !c.IsLeaf() {
				cp.errs = append(cp.errs, cp.where+": the passengers preset needs requirements that are counts under one all_of; write the outcomes explicitly")
				return nil
			}
			leaves = append(leaves, c.ID)
		}
	default:
		cp.errs = append(cp.errs, cp.where+": the passengers preset needs requirements that are counts under one all_of; write the outcomes explicitly")
		return nil
	}
	if len(leaves) == 0 {
		return []engine.Stage{{ID: "not_met", When: cond("always"), Status: status, MessageKey: msg}}
	}
	var out []engine.Stage
	for i := len(leaves) - 1; i >= 1; i-- {
		out = append(out, engine.Stage{ID: "not_met_" + leaves[i], When: engine.Condition{Op: "unmet", Ref: leaves[i]}, Status: status, MessageKey: msg, Params: needed(leaves[i])})
	}
	return append(out, engine.Stage{ID: "not_met", When: cond("always"), Status: status, MessageKey: msg, Params: needed(leaves[0])})
}

// stages compiles an explicit list of { id?, when, status, message, params?, ref }.
func (cp *compiler) stages(n *yaml.Node, what string) []engine.Stage {
	var out []engine.Stage
	for i, item := range deref(n).Content {
		item = deref(item)
		s := engine.Stage{}
		where := fmt.Sprintf("%s[%d]", what, i)
		hasRef := false
		for _, p := range pairs(item) {
			k, v := p[0].Value, p[1]
			switch k {
			case "id":
				s.ID = v.Value
			case "when":
				if err := v.Decode(&s.When); err != nil {
					cp.fail(v, "%s.when: %v", where, err)
				}
			case "status":
				s.Status = v.Value
			case "message":
				s.MessageKey = v.Value
			case "params":
				if err := v.Decode(&s.Params); err != nil {
					cp.fail(v, "%s.params: %v", where, err)
				}
			case "ref":
				hasRef = true
				cp.ref(v, where)
			default:
				cp.fail(p[0], "%s: unknown key %q", where, k)
			}
		}
		if !hasRef {
			cp.fail(item, "%s has no ref", where)
		}
		if s.ID == "" {
			s.ID = fmt.Sprintf("%s_%d", s.Status, i)
		}
		out = append(out, s)
	}
	return out
}
