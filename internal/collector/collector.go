package collector

import (
	"context"
	"github.com/clambin/fluxcd-exporter/internal/flux"
	"github.com/prometheus/client_golang/prometheus"
	"k8s.io/client-go/rest"
	"log/slog"
	"sync"
)

var resourceInfoMetric = prometheus.NewDesc(
	prometheus.BuildFQName("gotk", "resource", "info"),
	"TODO",
	[]string{"name", "exported_namespace", "customresource_kind", "ready"},
	nil,
)

var defaultListers = []func(config *rest.Config, logger *slog.Logger) flux.Lister{
	flux.Kustomizations,
	flux.HelmReleases,
	flux.GitRepositories,
	flux.HelmRepositories,
	flux.Buckets,
}

type Collector struct {
	Config  *rest.Config
	Logger  *slog.Logger
	listers []func(config *rest.Config, logger *slog.Logger) flux.Lister
}

func (c Collector) Describe(ch chan<- *prometheus.Desc) {
	ch <- resourceInfoMetric
}

func (c Collector) Collect(ch chan<- prometheus.Metric) {
	listers := c.listers
	if listers == nil {
		listers = defaultListers
	}

	var wg sync.WaitGroup
	for _, lister := range listers {
		wg.Add(1)
		go func(lister func(config *rest.Config, logger *slog.Logger) flux.Lister) {
			defer wg.Done()
			c.getResources(lister(c.Config, c.Logger), ch)
		}(lister)
	}
	wg.Wait()
}

func (c Collector) getResources(l flux.Lister, ch chan<- prometheus.Metric) {
	fluxResources, err := l.List(context.Background())
	if err != nil {
		c.Logger.Error("failed to get flux metrics", "err", err)
		ch <- prometheus.NewInvalidMetric(prometheus.NewDesc("fluxcd_exporter_error", "Error getting custom resource status", nil, nil), err)
		return
	}
	for _, fluxResource := range fluxResources {
		c.Logger.Debug("flux custom resource found", "custom_resource", fluxResource)
		ch <- prometheus.MustNewConstMetric(resourceInfoMetric, prometheus.GaugeValue, 1.0,
			fluxResource.Name,
			fluxResource.Namespace,
			fluxResource.Kind,
			fluxResource.Conditions["ready"],
		)
	}
}
