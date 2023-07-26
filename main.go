package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/spf13/cobra"
)

// cmd is the starting point of the govanity program.
var cmd = &cobra.Command{
	Use:     "govanity [flags] package_mapping ...",
	Short:   "govanity is a simple vanity URL server for Go packages",
	Example: "  govanity codello.dev/govanity=https://github.com/codello/govanity",
	Long: `
govanity is a simple server for Go vanity URLs.

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
	CompletionOptions: cobra.CompletionOptions{HiddenDefaultCmd: true},
	RunE:              run,
}

var (
	version       bool
	listenAddress string // bind address for the server
)

func init() {
	cmd.SetHelpCommand(&cobra.Command{
		Use:    "no-help",
		Hidden: true,
	})
	cmd.Flags().BoolVarP(&version, "version", "v", false, "Display hte version and exit.")
	cmd.Flags().StringVarP(&listenAddress, "listen-address", "l", ":8080", "The address on which the server runs.")
}

func run(_ *cobra.Command, args []string) error {
	if version {
		printVersion()
		return nil
	}
	setupHealthcheck()
	if err := setupServer(args); err != nil {
		return err
	}
	startMetricsServer()
	log.Printf("Running on %s\n", listenAddress)
	return http.ListenAndServe(listenAddress, logRequest(http.DefaultServeMux))
}

func logRequest(handler http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Printf("%s %s %s%s\n", r.RemoteAddr, r.Method, r.Host, r.URL)
		handler.ServeHTTP(w, r)
	})
}

func main() {
	if err := cmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
