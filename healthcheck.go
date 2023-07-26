package main

import (
	"io"
	"net/http"
)

var (
	healthcheckPath string
)

func init() {
	cmd.Flags().StringVarP(&healthcheckPath, "healthcheck-path", "", "/health", "Request path for the healthcheck endpoint.")
}

func setupHealthcheck() {
	if healthcheckPath == "" {
		return
	}
	http.HandleFunc(healthcheckPath, func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/plain")
		w.WriteHeader(http.StatusOK)
		_, _ = io.WriteString(w, "OK")
	})
}
