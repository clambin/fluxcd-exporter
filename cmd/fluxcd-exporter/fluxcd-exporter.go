package main

import (
	"errors"
	"flag"
	"github.com/clambin/fluxcd-exporter/internal/collector"
	"github.com/go-logr/logr"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"log/slog"
	"net/http"
	"os"
	"sigs.k8s.io/controller-runtime/pkg/client/config"
	"sigs.k8s.io/controller-runtime/pkg/log"
)

var (
	debug = flag.Bool("debug", false, "switch on debug logging")
	addr  = flag.String("addr", ":9090", "metrics listener port")
)

func main() {
	flag.Parse()

	var opts slog.HandlerOptions
	if *debug {
		opts.Level = slog.LevelDebug
	}
	logger := slog.New(slog.NewJSONHandler(os.Stderr, &opts))

	log.SetLogger(logr.FromSlogHandler(logger.Handler()))

	c := collector.Collector{
		Config: config.GetConfigOrDie(),
		Logger: logger,
	}
	prometheus.MustRegister(c)

	http.Handle("/metrics", promhttp.Handler())

	if err := http.ListenAndServe(*addr, nil); !errors.Is(err, http.ErrServerClosed) {
		panic(err)
	}
}
