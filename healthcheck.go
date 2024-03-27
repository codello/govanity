package main

import (
	"io"
	"net/http"
)

// healthcheck implements an HTTP endpoint that performs a health check.
func healthcheck(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/plain")
	w.WriteHeader(http.StatusOK)
	_, _ = io.WriteString(w, "OK")
}
