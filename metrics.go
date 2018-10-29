package main

import (
	"fmt"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"log"
	"net/http"
)

var (
	// namespace for all metrics of the application
	namespace = "gohole"
)

// runPrometheusServer starts an HTTP server which exposes
// the application metrics in the Prometheus format.
func runPrometheusServer() {
	port := getEnvOrDefault("PROMETHEUS_PORT", "9090")

	fmt.Printf("Starting HTTP server with metrics on TCP port %s...\n", port)
	server := &http.Server{Addr: "0.0.0.0:" + port}
	http.Handle("/metrics", promhttp.Handler())
	err := server.ListenAndServe()
	if err != nil {
		log.Fatal(err)
	}
}
