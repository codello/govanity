package main

import (
	"log"
	"net/http"

	"github.com/prometheus/client_golang/prometheus/promhttp"
)

var (
	metricsAddress string
)

func init() {
	cmd.Flags().StringVarP(&metricsAddress, "metrics-address", "", ":9090", "The address on which the metrics server runs.")
}

func startMetricsServer() {
	if metricsAddress != "" {
		mux := http.NewServeMux()
		mux.Handle("/metrics", promhttp.Handler())
		log.Printf("Metrics running on %s\n", metricsAddress)
		go func() {
			_ = http.ListenAndServe(metricsAddress, mux)
		}()
	}
}
