package server

import (
	_ "embed"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/spf13/cobra"
)

var (
	listenAddress string // bind address for the server

	//go:embed index.html.tmpl
	tmplData string             // template used for 200 OK responses
	tmpl     *template.Template // parsed tmplData
)

// Command is the command that runs the govanity server.
var Command = &cobra.Command{
	Use:     "server [flags] package_mapping ...",
	Short:   "Run the govanity server",
	Example: "  govanity server go.codello.dev/govanity=github.com/codello/govanity",
	Long: `govanity is a simple server for Go vanity URLs.

The configuration is done entirely via the command line by specifying package
mappings. A package mapping looks like this:
  prefix=[vcs:]repo-root
This maps the package prefix to the repo root using the vcs protocol. If vcs is
not present, it defaults to git.

If the prefix starts with a / the prefix will be matched independently of the
hostname of the request. The server will then prepend the hostname of the
request to the package name. Use this if you want to map the prefix to a repo,
regardless of the hostname used to resolve the package.
If the prefix does not start with a / the first path component is assumed to be
the hostname. The prefix will only match if the request hostname matches this
part of the prefix and the part matches the corresponding prefix.
`,
	RunE: runServer,
	Args: cobra.MinimumNArgs(1),
}

func init() {
	// Command setup
	Command.Flags().StringVarP(&listenAddress, "listen-address", "l", ":8080", "The address on which the server runs.")

	// Prepare template
	tmpl = template.Must(template.New("package").Funcs(map[string]any{
		"hasPrefix": strings.HasPrefix,
	}).Parse(tmplData))
}

// runServer actually runs the server command.
func runServer(_ *cobra.Command, args []string) error {
	for _, spec := range args {
		prefix, vcs, repoRoot, err := ParsePackageMapping(spec)
		if err != nil {
			_, _ = fmt.Fprintf(os.Stderr, "%s\n", err)
			os.Exit(1)
		}
		h := PackageHandler(prefix, vcs, repoRoot)
		http.Handle(prefix, h)
		http.Handle(prefix+"/", h)
	}
	log.Printf("Running on %s\n", listenAddress)
	return http.ListenAndServe(listenAddress, nil)
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
