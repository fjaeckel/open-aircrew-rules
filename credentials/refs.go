package credentials

import (
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"slices"
	"strings"
	"sync"
)

// Ref is a parsed `<prefix>:<article><labels>[@anchor]` reference.
type Ref struct {
	Raw     string
	Prefix  string
	Article string
	Labels  []string
	// Anchor is an optional hash of the quoted passage (change detection hook; not verified yet).
	Anchor string
}

var (
	refShape   = regexp.MustCompile(`^([a-z]+):([^@(]+)((?:\([0-9A-Za-z]{1,6}\))*)(?:@([A-Za-z0-9:_-]+))?$`)
	labelToken = regexp.MustCompile(`\(([0-9A-Za-z]{1,6})\)`)
)

// ParseRef parses a reference string.
func ParseRef(s string) (Ref, error) {
	m := refShape.FindStringSubmatch(strings.TrimSpace(s))
	if m == nil {
		return Ref{}, fmt.Errorf("ref %q: want <prefix>:<article>(label)...", s)
	}
	r := Ref{Raw: s, Prefix: m[1], Article: strings.TrimSpace(m[2]), Anchor: m[4]}
	for _, l := range labelToken.FindAllStringSubmatch(m[3], -1) {
		r.Labels = append(r.Labels, l[1])
	}
	return r, nil
}

// Cite returns the human citation, e.g. "EASA FCL.740.A(b)(1)".
func (v *Vocab) Cite(r Ref) string {
	a, ok := v.RefAuthorities[r.Prefix]
	if !ok {
		return r.Raw
	}
	lab := ""
	for _, l := range r.Labels {
		lab += "(" + l + ")"
	}
	return strings.ReplaceAll(a.Cite, "{article}", r.Article+lab)
}

// Resolver checks references against sources/ and the policies of policies.yaml.
type Resolver struct {
	root     string
	v        *Vocab
	policies map[string]*Policy
	mu       sync.Mutex
	files    map[string][]label
}

// NewResolver returns a resolver for the module rooted at root.
func NewResolver(root string, v *Vocab, policies map[string]*Policy) *Resolver {
	return &Resolver{root: root, v: v, policies: policies, files: map[string][]label{}}
}

// SourceFile returns the sources/ path of an article reference.
func (rs *Resolver) SourceFile(r Ref) (string, error) {
	a, ok := rs.v.RefAuthorities[r.Prefix]
	if !ok {
		return "", fmt.Errorf("unknown ref prefix %q", r.Prefix)
	}
	name := r.Article
	if a.File == "lower_dashed" {
		name = strings.ReplaceAll(strings.ToLower(name), ".", "-")
	}
	return filepath.Join(rs.root, "sources", a.Directory, name+".md"), nil
}

// Resolve checks one reference: the prefix is known, a policy is declared, the article
// file exists and the paragraph labels occur in order in its text.
func (rs *Resolver) Resolve(s string) error {
	r, err := ParseRef(s)
	if err != nil {
		return err
	}
	if r.Prefix == "policy" {
		if len(r.Labels) > 0 {
			return fmt.Errorf("ref %q: a policy ref has no paragraph labels", s)
		}
		if _, ok := rs.policies[r.Article]; !ok {
			return fmt.Errorf("ref %q: policy not declared in policies.yaml", s)
		}
		return nil
	}
	path, err := rs.SourceFile(r)
	if err != nil {
		return fmt.Errorf("ref %q: %v", s, err)
	}
	labels, err := rs.labels(path)
	if err != nil {
		return fmt.Errorf("ref %q: no source file %s", s, rel(rs.root, path))
	}
	if !findPath(labels, r.Labels) {
		return fmt.Errorf("ref %q: paragraph %s not found in %s", s, strings.Join(wrap(r.Labels), ""), rel(rs.root, path))
	}
	return nil
}

func wrap(l []string) []string {
	out := make([]string, len(l))
	for i, s := range l {
		out[i] = "(" + s + ")"
	}
	return out
}

func rel(root, p string) string {
	if r, err := filepath.Rel(root, p); err == nil {
		return r
	}
	return p
}

// label is one paragraph with its full label path, e.g. [b 1 ii C].
type label struct {
	path []string
}

var (
	leadLabels = regexp.MustCompile(`^((?:\([0-9A-Za-z]{1,6}\)\s*)+)`)
	heading    = regexp.MustCompile(`^\*[^*]+\*\s*`)
	romans     = []string{"i", "ii", "iii", "iv", "v", "vi", "vii", "viii", "ix", "x", "xi", "xii", "xiii", "xiv", "xv", "xvi", "xvii", "xviii", "xix", "xx"}
)

func (rs *Resolver) labels(path string) ([]label, error) {
	rs.mu.Lock()
	defer rs.mu.Unlock()
	if l, ok := rs.files[path]; ok {
		return l, nil
	}
	b, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	l := parseLabels(string(b))
	rs.files[path] = l
	return l, nil
}

// outline is one open level of the paragraph outline.
type outline struct {
	text string
	kind string
	ord  int
}

// ordinal returns the position of a label in a list of the given kind (0 if impossible).
func ordinal(s, kind string) int {
	switch kind {
	case "num":
		n := 0
		for _, r := range s {
			if r < '0' || r > '9' {
				return 0
			}
			n = n*10 + int(r-'0')
		}
		return n
	case "roman":
		return slices.Index(romans, s) + 1
	case "lower", "upper":
		if len(s) != 1 {
			return 0
		}
		c := s[0]
		if kind == "lower" && c >= 'a' && c <= 'z' {
			return int(c-'a') + 1
		}
		if kind == "upper" && c >= 'A' && c <= 'Z' {
			return int(c-'A') + 1
		}
	}
	return 0
}

func kinds(s string) []string {
	switch {
	case s[0] >= '0' && s[0] <= '9':
		return []string{"num"}
	case s[0] >= 'A' && s[0] <= 'Z':
		return []string{"upper"}
	case slices.Contains(romans, s):
		return []string{"roman", "lower"}
	}
	return []string{"lower"}
}

// place puts a label into the outline: as the next sibling of an open level, else as the
// first child of the current level.
func place(stack []outline, s string) []outline {
	ks := kinds(s)
	for d := len(stack) - 1; d >= 0; d-- {
		e := stack[d]
		if slices.Contains(ks, e.kind) && ordinal(s, e.kind) == e.ord+1 {
			return append(stack[:d], outline{s, e.kind, e.ord + 1})
		}
	}
	for _, k := range ks {
		if ordinal(s, k) == 1 {
			return append(stack, outline{s, k, 1})
		}
	}
	return append(stack, outline{s, ks[0], ordinal(s, ks[0])})
}

// parseLabels returns the label path of every paragraph that opens a line of the text
// after the first second-level heading, including labels that follow an italic heading.
func parseLabels(text string) []label {
	var out []label
	var stack []outline
	body := false
	for _, line := range strings.Split(text, "\n") {
		if strings.HasPrefix(line, "## ") {
			body = true
			continue
		}
		if !body {
			continue
		}
		rest := strings.TrimSpace(line)
		for {
			m := leadLabels.FindString(rest)
			if m == "" {
				break
			}
			for _, t := range labelToken.FindAllStringSubmatch(m, -1) {
				stack = place(stack, t[1])
				p := make([]string, len(stack))
				for i, e := range stack {
					p[i] = e.text
				}
				out = append(out, label{path: p})
			}
			rest = strings.TrimSpace(rest[len(m):])
			h := heading.FindString(rest)
			if h == "" {
				break
			}
			rest = rest[len(h):]
		}
	}
	return out
}

// findPath reports whether a paragraph with exactly the wanted label path exists.
func findPath(seq []label, want []string) bool {
	if len(want) == 0 {
		return true
	}
	for _, l := range seq {
		if slices.Equal(l.path, want) {
			return true
		}
	}
	return false
}
