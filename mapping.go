package main

import (
	_ "embed"
	"fmt"
	"html/template"
	"net/http"
	"strings"
)

var (
	//go:embed index.html.tmpl
	tmplData string             // template used for 200 OK responses
	tmpl     *template.Template // parsed tmplData
)

func init() {
	// Prepare template
	tmpl = template.Must(template.New("package").Funcs(map[string]any{
		"hasPrefix": strings.HasPrefix,
	}).Parse(tmplData))
}

type PackageMapping struct {
	Prefix   string
	VCS      string
	repoRoot string
}

// ParsePackageMapping converts a "prefix=[vcs:]repo-root" spec into prefix, vcs and repoRoot.
// If a parsing error occurs, err will be non-nil.
func ParsePackageMapping(spec string) (*PackageMapping, error) {
	prefix, repoRoot, ok := strings.Cut(spec, "=")
	if !ok {
		return nil, fmt.Errorf("not a valid package mapping: %s", spec)
	}
	prefix = strings.TrimSpace(prefix)
	if prefix == "" {
		return nil, fmt.Errorf("missing prefix in package mapping: %s", spec)
	}

	if strings.HasSuffix(prefix, "/") {
		return nil, fmt.Errorf("package cannot end in a slash: %s", prefix)
	}

	vcs, repoRoot, ok := strings.Cut(repoRoot, ":")
	if ok {
		cleanVCS := strings.ToLower(strings.TrimSpace(vcs))
		switch cleanVCS {
		case "bzr", "fossil", "git", "hg", "svn":
			vcs = cleanVCS
		default:
			repoRoot = vcs + ":" + repoRoot
			vcs = ""
		}
	} else {
		vcs, repoRoot = repoRoot, vcs
	}
	if vcs == "" {
		vcs = "git"
	}
	repoRoot = strings.TrimSpace(repoRoot)
	if repoRoot == "" {
		return nil, fmt.Errorf("repository root cannot be empty: %s", spec)
	}
	return &PackageMapping{
		Prefix:   prefix,
		VCS:      vcs,
		repoRoot: repoRoot,
	}, nil
}

// ServeHTTP implements the http.Handler interface, renders the package meta page for the package prefix.
// The rendered page will indicate that the source for the page can be found at repoRoot using vcs.
func (m *PackageMapping) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	prefix := m.Prefix
	if strings.HasPrefix(prefix, "/") {
		prefix = r.Host + prefix
	}
	_ = tmpl.Execute(w, map[string]string{
		"Prefix":   prefix,
		"VCS":      m.VCS,
		"RepoRoot": m.repoRoot,
		"URL":      r.Host + r.URL.Path,
	})
}
