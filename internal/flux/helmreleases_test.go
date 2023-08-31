package flux

import (
	"context"
	"github.com/fluxcd/helm-controller/api/v2beta1"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"sigs.k8s.io/controller-runtime/pkg/client/fake"
	"sigs.k8s.io/controller-runtime/pkg/scheme"
	"testing"
)

func TestHelmReleases(t *testing.T) {
	testCases := []testCase{
		{
			name:          "valid",
			schemeBuilder: v2beta1.SchemeBuilder,
			objects: []runtime.Object{
				&v2beta1.HelmRelease{
					TypeMeta: v1.TypeMeta{
						Kind:       "HelmRelease",
						APIVersion: "",
					},
					ObjectMeta: v1.ObjectMeta{
						Name:      "foo",
						Namespace: "bar",
					},
					Status: v2beta1.HelmReleaseStatus{
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
				Kind:      "HelmRelease",
				Conditions: map[string]string{
					"ready": "True",
				},
			}},
		},
		{
			name:          "empty",
			schemeBuilder: v2beta1.SchemeBuilder,
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
			k := lister{client: c, list: getHelmReleases}

			resources, err := k.List(context.Background())
			tt.wantErr(t, err)
			if err == nil {
				assert.Equal(t, tt.want, resources)
			}
		})
	}
}
