package main

import (
	_ "embed"
	"fmt"
	"html/template"
	"net/http"
	"strings"
)

var (
	noCache     bool // whether to send the Cache-Control header
	cacheMaxAge int  // max-age cache control
	maxStaleAge int  // max-stale cache control
	errorAge    int  // stale-if-error cache control

	//go:embed index.html.tmpl
	tmplData string             // template used for 200 OK responses
	tmpl     *template.Template // parsed tmplData
)

func init() {
	// Command setup
	cmd.Flags().BoolVarP(&noCache, "no-cache", "", false, "Disables the Cache-Control header.")
	cmd.Flags().IntVarP(&cacheMaxAge, "cache-max-age", "", 604800, "Cache-Control max-age value.")
	cmd.Flags().IntVarP(&errorAge, "cache-stale-if-error", "", 86400, "Cache-Control stale-if-error value.")
	cmd.Flags().IntVarP(&maxStaleAge, "cache-max-stale", "", 3600, "Cache-Control max-stale value.")

	// Prepare template
	tmpl = template.Must(template.New("package").Funcs(map[string]any{
		"hasPrefix": strings.HasPrefix,
	}).Parse(tmplData))
}

// runServer actually runs the server command.
func setupServer(specs []string) error {
	for _, spec := range specs {
		prefix, vcs, repoRoot, err := ParsePackageMapping(spec)
		if err != nil {
			return err
		}
		h := PackageHandler(prefix, vcs, repoRoot)
		http.Handle(prefix, h)
		http.Handle(prefix+"/", h)
	}
	return nil
}

// ParsePackageMapping converts a "prefix=[vcs:]repo-root" spec into prefix, vcs and repoRoot.
// If a parsing error occurs, err will be non-nil.
func ParsePackageMapping(spec string) (prefix, vcs, repoRoot string, err error) {
	var ok bool
	prefix, repoRoot, ok = strings.Cut(spec, "=")
	if !ok {
		err = fmt.Errorf("not a valid package mapping: %s", spec)
		return
	}
	prefix = strings.TrimSpace(prefix)

	if strings.HasSuffix(prefix, "/") {
		err = fmt.Errorf("package cannot end in a slash: %s", prefix)
	}

	vcs, repoRoot, ok = strings.Cut(repoRoot, ":")
	if ok {
		cleanVCS := strings.ToLower(strings.TrimSpace(vcs))
		switch cleanVCS {
		case "bzr", "fossil", "git", "hg", "svc":
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
		err = fmt.Errorf("repository root cannot be empty: %s", spec)
		return
	}
	return
}

// PackageHandler returns a http.Handler that renders the package meta page for the package prefix.
// The rendered page will indicate that the source for the page can be found at repoRoot using vcs.
func PackageHandler(prefix string, vcs string, repoRoot string) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		if !noCache {
			w.Header().Set("Cache-Control", fmt.Sprintf("public, max-age=%d, max-stale=%d, stale-if-error=%d", cacheMaxAge, maxStaleAge, maxStaleAge))
		}
		w.WriteHeader(http.StatusOK)
		if strings.HasPrefix(prefix, "/") {
			prefix = r.Host + prefix
		}
		_ = tmpl.Execute(w, map[string]string{
			"Prefix":   prefix,
			"VCS":      vcs,
			"RepoRoot": repoRoot,
			"URL":      r.Host + r.URL.Path,
		})
	}
	return http.HandlerFunc(fn)
}
