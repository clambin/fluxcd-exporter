package collector

import (
	"bytes"
	"context"
	"fluxcd-exporter/internal/flux"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/testutil"
	"github.com/stretchr/testify/assert"
	"k8s.io/client-go/rest"
	"log/slog"
	"testing"
)

func TestCollector_Collect(t *testing.T) {
	c := Collector{
		Logger: slog.Default(),
		listers: []func(config *rest.Config) flux.Lister{
			func(config *rest.Config) flux.Lister {
				return flux.ListerFunc(func(ctx context.Context) (flux.Resources, error) {
					return flux.Resources{{
						Name:       "foo",
						Namespace:  "bar",
						Kind:       "snafu",
						Conditions: map[string]string{"ready": "False"},
					}}, nil
				})
			},
		}}

	r := prometheus.NewPedanticRegistry()
	r.MustRegister(c)
	assert.NoError(t, testutil.GatherAndCompare(r, bytes.NewBufferString(`
# HELP gotk_resource_info TODO
# TYPE gotk_resource_info gauge
gotk_resource_info{kind="snafu",name="foo",namespace="bar",ready="False"} 1
`)))
}
