package collector

import (
	"bytes"
	"context"
	"fmt"
	"github.com/clambin/fluxcd-exporter/internal/flux"
	"github.com/prometheus/client_golang/prometheus/testutil"
	"github.com/stretchr/testify/assert"
	"k8s.io/client-go/rest"
	"log/slog"
	"testing"
)

func TestCollector_Collect(t *testing.T) {
	testCases := []struct {
		name    string
		lister  flux.ListerFunc
		wantErr assert.ErrorAssertionFunc
		want    string
	}{
		{
			name: "valid",
			lister: func(ctx context.Context) (flux.Resources, error) {
				return flux.Resources{{
					Name:       "foo",
					Namespace:  "bar",
					Kind:       "snafu",
					Conditions: map[string]string{"ready": "False"},
				}}, nil
			},
			wantErr: assert.NoError,
			want: `
# HELP gotk_resource_info TODO
# TYPE gotk_resource_info gauge
gotk_resource_info{customresource_kind="snafu",exported_namespace="bar",name="foo",ready="False"} 1
`,
		},
		{
			name: "error",
			lister: func(ctx context.Context) (flux.Resources, error) {
				return nil, fmt.Errorf("failed")
			},
			wantErr: assert.Error,
			want:    "",
		},
	}

	for _, tt := range testCases {
		t.Run(tt.name, func(t *testing.T) {
			c := Collector{
				Logger: slog.Default(),
				listers: []func(config *rest.Config) flux.Lister{
					func(config *rest.Config) flux.Lister {
						return tt.lister
					},
				},
			}
			tt.wantErr(t, testutil.CollectAndCompare(c, bytes.NewBufferString(tt.want)))
		})
	}
}
