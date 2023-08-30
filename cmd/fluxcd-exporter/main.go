package main

import (
	"errors"
	"flag"
	"fluxcd-exporter/internal/collector"
	"fmt"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
	"log/slog"
	"net/http"
	"os"
	"path/filepath"
)

var debug = flag.Bool("debug", false, "switch on debug logging")

func main() {
	flag.Parse()

	var opts slog.HandlerOptions
	if *debug {
		opts.Level = slog.LevelDebug
	}
	logger := slog.New(slog.NewTextHandler(os.Stderr, &opts))

	cfg, err := getK8sConfig()
	if err != nil {
		panic(err)
	}

	c := collector.Collector{
		Config: cfg,
		Logger: logger,
	}
	prometheus.MustRegister(c)

	http.Handle("/metrics", promhttp.Handler())

	if err = http.ListenAndServe(":8080", nil); !errors.Is(err, http.ErrServerClosed) {
		panic(err)
	}
}
func getK8sConfig() (*rest.Config, error) {
	if cfg, err := rest.InClusterConfig(); err == nil {
		return cfg, nil
	}

	// not running inside cluster. try to connect as external client
	userHomeDir, err := os.UserHomeDir()
	if err != nil {
		return nil, fmt.Errorf("user home dir: %w", err)
	}
	kubeConfigPath := filepath.Join(userHomeDir, ".kube", "config")
	slog.Debug("not running inside cluster. using kube config", "filename", kubeConfigPath)

	cfg, err := clientcmd.BuildConfigFromFlags("", kubeConfigPath)
	return cfg, err
}
