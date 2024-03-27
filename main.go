package main

import (
	"fmt"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"log/slog"
	"net/http"
	"os"

	"github.com/lmittmann/tint"
	"github.com/spf13/cobra"
)

var (
	version       bool
	logFormat     string
	listenAddress string // bind address for the server
)

func init() {
	cmd.SetHelpCommand(&cobra.Command{
		Use:    "no-help",
		Hidden: true,
	})
	cmd.Flags().BoolVarP(&version, "version", "v", false, "Display hte version and exit.")
	cmd.Flags().StringVar(&logFormat, "log", "color", "The logging format. Specify 'json' for JSON logs.")
	cmd.Flags().StringVarP(&listenAddress, "listen-address", "l", ":8080", "The address on which the server runs.")
}

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
	RunE: func(_ *cobra.Command, args []string) error {
		if version {
			printVersion()
			return nil
		}

		var logger *slog.Logger
		switch logFormat {
		case "color":
			logger = slog.New(tint.NewHandler(os.Stdout, nil))
		case "json":
			logger = slog.New(slog.NewJSONHandler(os.Stdout, nil))
		default:
			logger = slog.New(slog.NewTextHandler(os.Stdout, nil))
		}
		for _, spec := range args {
			m, err := ParsePackageMapping(spec)
			if err != nil {
				return err
			}
			handler := CacheControl(m)
			http.Handle(fmt.Sprintf("GET %s", m.Prefix), handler)
			http.Handle(fmt.Sprintf("GET %s/", m.Prefix), handler)
		}
		http.HandleFunc("GET /health", healthcheck)
		http.Handle("GET /metrics", promhttp.Handler())
		logger.Info(fmt.Sprintf("Running on %s", listenAddress))
		return http.ListenAndServe(listenAddress, RequestLogger(logger)(http.DefaultServeMux))
	},
}

func main() {
	if cmd.Execute() != nil {
		os.Exit(1)
	}
}
