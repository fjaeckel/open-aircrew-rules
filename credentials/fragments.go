package credentials

import (
	"os"
	"path/filepath"
	"sort"
	"strings"
)

// FragmentKinds are the directories under fragments/, one per shared file that parallel
// work must not edit directly: coverage (coverage/articles.yaml), keys (messages/keys.yaml),
// policies (policies.yaml), changelog (CHANGELOG.md) and vocab-requests (vocabulary.yaml).
var FragmentKinds = []string{"coverage", "keys", "policies", "changelog", "vocab-requests"}

// fragmentFiles lists the fragment files of one kind, README.md excluded.
func fragmentFiles(root, kind string) ([]string, error) {
	entries, err := os.ReadDir(filepath.Join(root, "fragments", kind))
	if os.IsNotExist(err) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	var out []string
	for _, e := range entries {
		if e.IsDir() || e.Name() == "README.md" || strings.HasPrefix(e.Name(), ".") {
			continue
		}
		out = append(out, filepath.Join(root, "fragments", kind, e.Name()))
	}
	sort.Strings(out)
	return out, nil
}

// Fragments lists every unmerged fragment file by kind (paths relative to root).
func Fragments(root string) (map[string][]string, error) {
	out := map[string][]string{}
	for _, k := range FragmentKinds {
		fs, err := fragmentFiles(root, k)
		if err != nil {
			return nil, err
		}
		for _, f := range fs {
			out[k] = append(out[k], rel(root, f))
		}
	}
	return out, nil
}
