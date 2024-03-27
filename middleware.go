package main

import (
	"fmt"
	"log/slog"
	"net/http"
)

var (
	noCache     bool // whether to send the Cache-Control header
	cacheMaxAge int  // max-age cache control
	maxStaleAge int  // max-stale cache control
	errorAge    int  // stale-if-error cache control
)

func init() {
	// Command setup
	cmd.Flags().BoolVarP(&noCache, "no-cache", "", false, "Disables the Cache-Control header.")
	cmd.Flags().IntVarP(&cacheMaxAge, "cache-max-age", "", 604800, "Cache-Control max-age value.")
	cmd.Flags().IntVarP(&errorAge, "cache-stale-if-error", "", 86400, "Cache-Control stale-if-error value.")
	cmd.Flags().IntVarP(&maxStaleAge, "cache-max-stale", "", 3600, "Cache-Control max-stale value.")
}

// RequestLogger is an HTTP middleware that logs requests.
func RequestLogger(logger *slog.Logger) func(handler http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			logger.InfoContext(r.Context(), fmt.Sprintf("%s %s%s", r.Method, r.Host, r.URL), "remoteAddr", r.RemoteAddr)
			next.ServeHTTP(w, r)
		})
	}
}

// CacheControl is an HTTP middleware that adds a Cache-Control header to the response.
func CacheControl(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if !noCache {
			w.Header().Set("Cache-Control", fmt.Sprintf("public, max-age=%d, max-stale=%d, stale-if-error=%d", cacheMaxAge, maxStaleAge, maxStaleAge))
		}
		next.ServeHTTP(w, r)
	})
}
