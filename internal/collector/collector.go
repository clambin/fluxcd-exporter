package collector

import (
	"context"
	"fluxcd-exporter/internal/flux"
	"fmt"
	"github.com/prometheus/client_golang/prometheus"
	"k8s.io/client-go/rest"
	"log/slog"
)

var resourceInfoMetric = prometheus.NewDesc(
	prometheus.BuildFQName("gotk", "resource", "info"),
	"TODO",
	[]string{"name", "exported_namespace", "customresource_kind", "ready"},
	nil,
)

var defaultListers = []func(config *rest.Config) flux.Lister{
	flux.Kustomizations,
	flux.HelmReleases,
	flux.GitRepositories,
	flux.HelmRepositories,
	flux.Buckets,
}

type Collector struct {
	Config  *rest.Config
	Logger  *slog.Logger
	listers []func(config *rest.Config) flux.Lister
}

func (c Collector) Describe(ch chan<- *prometheus.Desc) {
	ch <- resourceInfoMetric
}

func (c Collector) Collect(ch chan<- prometheus.Metric) {
	fluxResources, err := c.getResources()
	if err != nil {
		ch <- prometheus.NewInvalidMetric(prometheus.NewDesc("fluxmon_error", "Error getting custom resource status", nil, nil), err)
		c.Logger.Error("failed to get resource status", "err", err)
		return
	}
	for _, fluxResource := range fluxResources {
		ch <- prometheus.MustNewConstMetric(resourceInfoMetric, prometheus.GaugeValue, 1.0,
			fluxResource.Name,
			fluxResource.Namespace,
			fluxResource.Kind,
			fluxResource.Conditions["ready"],
		)
	}
}

func (c Collector) getResources() (flux.Resources, error) {
	var fluxResources flux.Resources

	listers := c.listers
	if listers == nil {
		listers = defaultListers
	}
	for _, lister := range listers {
		resources, err := lister(c.Config).List(context.Background())
		if err != nil {
			return nil, fmt.Errorf("flux: %w", err)
		}
		fluxResources = append(fluxResources, resources...)
	}
	return fluxResources, nil
}
