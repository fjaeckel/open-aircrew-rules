package engine

import (
	"slices"
	"strings"
)

// subjectRef is a selected subject plus the record data it refers to.
type subjectRef struct {
	Subject
	licence        *Licence
	licenceKind    string
	typeDesignator string
	variant        string
	issued         *Date
	validFrom      *Date
	expires        *Date
	// category is the subject's aircraft category ("" when neither recorded nor derivable).
	category string
	// missing names the applies_to criterion whose data the selected item lacks; the
	// evaluation is then unknown (AppliesTo.UnknownWhenMissing).
	missing string
}

func containsFold(list []string, v string) bool {
	return slices.ContainsFunc(list, func(s string) bool { return normAuthority(s) == normAuthority(v) })
}

// licenceMatches applies the authority and licence-kind parts of applies_to.
func (p *prepared) licenceMatches(a *AppliesTo, l *Licence) bool {
	if l == nil {
		return len(a.Authorities) == 0 && len(a.LicenceKinds) == 0
	}
	if len(a.Authorities) > 0 && !containsFold(a.Authorities, l.Authority) {
		return false
	}
	if containsFold(a.ExcludeAuthorities, l.Authority) {
		return false
	}
	kind := p.licenceKind[l.ID]
	if len(a.LicenceKinds) > 0 && !slices.Contains(a.LicenceKinds, kind) {
		return false
	}
	return !slices.Contains(a.ExcludeLicenceKinds, kind)
}

// ratingMatches applies the class and rating parts of applies_to. A criterion listed in
// UnknownWhenMissing whose data the rating lacks keeps the rating, named in missing.
func (p *prepared) ratingMatches(a *AppliesTo, r *Rating) (ok bool, missing string) {
	if len(a.Classes) > 0 && !slices.Contains(a.Classes, r.Class) {
		return false, ""
	}
	if slices.Contains(a.ExcludeClasses, r.Class) {
		return false, ""
	}
	miss := func(criterion string) bool {
		if !slices.Contains(a.UnknownWhenMissing, criterion) {
			return false
		}
		if missing == "" {
			missing = criterion
		}
		return true
	}
	if a.TypeRated != nil {
		switch {
		case *a.TypeRated && r.TypeDesignator == "":
			if !miss("typeRated") {
				return false, ""
			}
		case *a.TypeRated != (r.TypeDesignator != ""):
			return false, ""
		}
	}
	if len(a.ULKinds) > 0 {
		k := r.ULKind
		switch {
		case k == "" && !slices.Contains(a.ULKinds, "none"):
			if !miss("ulKinds") {
				return false, ""
			}
		case k == "":
		case !slices.Contains(a.ULKinds, k):
			return false, ""
		}
	}
	if len(a.Categories) > 0 {
		cat := p.ratingCategory(r)
		switch {
		case cat == "":
			if !miss("categories") {
				return false, ""
			}
		case !slices.Contains(a.Categories, cat):
			return false, ""
		}
	}
	return true, missing
}

// ratingCategory returns a rating's aircraft category: recorded, else its class's (unless
// the class names none), else its licence kind's.
func (p *prepared) ratingCategory(r *Rating) string {
	if r.Category != "" {
		return r.Category
	}
	if c := p.v.Category(r.Class); c != "" && c != "other" {
		return c
	}
	return p.v.LicenceCategory(p.licenceKind[r.LicenceID])
}

func (p *prepared) ratingSubject(kind string, h ratingHit) subjectRef {
	r := h.Rating
	l := p.licenceOf(r.LicenceID)
	return subjectRef{
		Subject:        Subject{Kind: kind, ID: r.ID, Class: r.Class, ULKind: r.ULKind},
		licence:        l,
		licenceKind:    p.licenceKind[r.LicenceID],
		typeDesignator: r.TypeDesignator,
		issued:         r.Issued, validFrom: r.ValidFrom, expires: r.Expires,
		category: p.ratingCategory(r),
		missing:  h.missing,
	}
}

// ratingHit is a rating applies_to selects, with the criterion whose data it lacks.
type ratingHit struct {
	*Rating
	missing string
}

// matchingRatings returns the ratings whose licence and class match applies_to.
func (p *prepared) matchingRatings(a *AppliesTo) []ratingHit {
	var out []ratingHit
	for i := range p.ratings {
		r := &p.ratings[i]
		l := p.licenceOf(r.LicenceID)
		if l == nil || !p.licenceMatches(a, l) {
			continue
		}
		if ok, missing := p.ratingMatches(a, r); ok {
			out = append(out, ratingHit{r, missing})
		}
	}
	return out
}

// subjects returns the subjects rule r applies to, in record order.
func (p *prepared) subjects(r *Rule, v *Vocabulary, asOf Date) []subjectRef {
	a := &r.AppliesTo
	var out []subjectRef
	switch a.Subject {
	case "rating":
		for _, rt := range p.matchingRatings(a) {
			out = append(out, p.ratingSubject("rating", rt))
		}
	case "passengers":
		seen := map[string]bool{}
		perType := a.TypeRated != nil && *a.TypeRated
		for _, rt := range p.matchingRatings(a) {
			class := rt.Class
			group := ""
			if a.ClassGroup != "" {
				if g, _ := v.ClassGroup(a.ClassGroup, rt.Class); g != "" {
					group, class = g, g
				}
			}
			key := normAuthority(p.licenceOf(rt.LicenceID).Authority) + "|" + class + "|" + rt.ULKind
			if perType {
				key += "|" + rt.TypeDesignator
			}
			if seen[key] {
				continue
			}
			seen[key] = true
			s := p.ratingSubject("passengers", rt)
			s.ID = rt.LicenceID
			s.Group = group
			if perType {
				s.Detail = rt.TypeDesignator
			}
			s.issued, s.validFrom, s.expires = nil, nil, nil
			out = append(out, s)
		}
	case "licence":
		for i := range p.rec.Licences {
			l := &p.rec.Licences[i]
			if p.licenceMatches(a, l) {
				out = append(out, subjectRef{
					Subject: Subject{Kind: "licence", ID: l.ID, Detail: p.licenceKind[l.ID]},
					licence: l, licenceKind: p.licenceKind[l.ID], issued: l.Issued, expires: l.Expires,
					category: v.LicenceCategory(p.licenceKind[l.ID]),
				})
			}
		}
	case "flight_review":
		for i := range p.rec.Licences {
			l := &p.rec.Licences[i]
			if p.licenceMatches(a, l) {
				out = append(out, subjectRef{Subject: Subject{Kind: "flight_review", ID: l.ID}, licence: l, licenceKind: p.licenceKind[l.ID]})
				break
			}
		}
	case "privilege":
		for i := range p.rec.Privileges {
			pv := &p.rec.Privileges[i]
			if len(a.PrivilegeKinds) > 0 && !slices.Contains(a.PrivilegeKinds, pv.Kind) {
				continue
			}
			l := p.licenceOf(pv.LicenceID)
			if !p.licenceMatches(a, l) {
				continue
			}
			s := subjectRef{
				Subject: Subject{Kind: "privilege", ID: pv.ID, Detail: pv.Detail},
				licence: l, licenceKind: p.licenceKind[pv.LicenceID],
				issued: pv.Issued, validFrom: pv.ValidFrom, expires: pv.Expires,
				category: v.LicenceCategory(p.licenceKind[pv.LicenceID]),
			}
			if len(a.Categories) > 0 && !slices.Contains(a.Categories, s.category) {
				if s.category != "" || !slices.Contains(a.UnknownWhenMissing, "categories") {
					continue
				}
				s.missing = "categories"
			}
			out = append(out, s)
		}
	case "credential":
		for i := range p.rec.Credentials {
			c := &p.rec.Credentials[i]
			if len(a.CredentialTypes) > 0 && !slices.Contains(a.CredentialTypes, c.Type) {
				continue
			}
			out = append(out, subjectRef{
				Subject: Subject{Kind: "credential", ID: c.ID, Detail: c.Type},
				issued:  c.Issued, validFrom: c.ValidFrom, expires: c.Expires,
			})
		}
	case "training":
		s := subjectRef{Subject: Subject{Kind: "training", Detail: a.Programme}}
		if len(a.Authorities) == 0 && len(a.LicenceKinds) == 0 {
			out = append(out, s)
			break
		}
		for i := range p.rec.Licences {
			l := &p.rec.Licences[i]
			if p.licenceMatches(a, l) {
				s.ID, s.licence, s.licenceKind = l.ID, l, p.licenceKind[l.ID]
				out = append(out, s)
				break
			}
		}
	case "type":
		for _, rt := range p.matchingRatings(a) {
			for i := range p.rec.Variants {
				vr := &p.rec.Variants[i]
				if vr.RatingID != rt.ID {
					continue
				}
				s := p.ratingSubject("type", rt)
				s.ID, s.Detail, s.variant = vr.ID, vr.Name, vr.Name
				s.issued, s.validFrom, s.expires = vr.Issued, nil, nil
				if want := a.DifferentEngineType; want != nil {
					switch {
					case vr.DifferentEngineType == nil:
						if !slices.Contains(a.UnknownWhenMissing, "differentEngineType") {
							continue
						}
						if s.missing == "" {
							s.missing = "differentEngineType"
						}
					case *vr.DifferentEngineType != *want:
						continue
					}
				}
				out = append(out, s)
			}
		}
	case "launch_method":
		for _, rt := range p.matchingRatings(a) {
			methods := map[string]bool{}
			for _, pv := range p.rec.Privileges {
				if pv.LicenceID == rt.LicenceID && pv.Kind == "LAUNCH_METHOD_TRAINED" {
					methods[strings.ToLower(pv.Detail)] = true
				}
			}
			for _, f := range p.flights {
				if f.Class == rt.Class && f.LaunchMethod != "" && !f.Date.After(asOf) {
					methods[f.LaunchMethod] = true
				}
			}
			for _, m := range v.LaunchMethods {
				if !methods[m] || (len(a.LaunchMethods) > 0 && !slices.Contains(a.LaunchMethods, m)) {
					continue
				}
				s := p.ratingSubject("launch_method", rt)
				s.Detail = m
				s.issued, s.validFrom, s.expires = nil, nil, nil
				out = append(out, s)
			}
		}
	}
	return out
}

// Holdings returns the record items applies_to selects on asOf: the subjects a rule with it
// evaluates, before its holds condition.
func Holdings(c *Catalogue, rec *Record, a AppliesTo, asOf Date) []Subject {
	p := prepare(rec, c.Vocabulary)
	var out []Subject
	for _, s := range p.subjects(&Rule{AppliesTo: a}, c.Vocabulary, asOf) {
		out = append(out, s.Subject)
	}
	return out
}
