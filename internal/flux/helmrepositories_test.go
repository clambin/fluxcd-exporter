package flux

import (
	"context"
	fluxSourceV1Beta2 "github.com/fluxcd/source-controller/api/v1beta2"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"log/slog"
	"sigs.k8s.io/controller-runtime/pkg/client/fake"
	"sigs.k8s.io/controller-runtime/pkg/scheme"
	"testing"
)

func TestHelmRepositories(t *testing.T) {
	testCases := []testCase{
		{
			name:          "valid",
			schemeBuilder: fluxSourceV1Beta2.SchemeBuilder,
			objects: []runtime.Object{
				&fluxSourceV1Beta2.HelmRepository{
					TypeMeta: v1.TypeMeta{
						Kind:       "HelmRepository",
						APIVersion: "",
					},
					ObjectMeta: v1.ObjectMeta{
						Name:      "foo",
						Namespace: "bar",
					},
					Status: fluxSourceV1Beta2.HelmRepositoryStatus{
						Conditions: []v1.Condition{
							{Type: "Ready", Status: "True"},
						},
					},
				},
			},
			wantErr: assert.NoError,
			want: []Resource{{
				Name:      "foo",
				Namespace: "bar",
				Kind:      "HelmRepository",
				Conditions: map[string]string{
					"ready": "True",
				},
			}},
		},
		{
			name:          "empty",
			schemeBuilder: fluxSourceV1Beta2.SchemeBuilder,
			wantErr:       assert.NoError,
			want:          []Resource(nil),
		},
		{
			name:          "error",
			schemeBuilder: &scheme.Builder{},
			wantErr:       assert.Error,
		},
	}

	for _, tt := range testCases {
		t.Run(tt.name, func(t *testing.T) {
			schema, err := tt.schemeBuilder.Build()
			require.NoError(t, err)

			c := fake.NewClientBuilder().WithScheme(schema).WithRuntimeObjects(tt.objects...).Build()
			k := lister{client: c, list: getHelmRepositories, logger: slog.Default()}

			resources, err := k.List(context.Background())
			tt.wantErr(t, err)
			if err == nil {
				assert.Equal(t, tt.want, resources)
			}
		})
	}
}
