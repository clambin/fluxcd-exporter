package flux

import (
	"context"
	fluxKustomizeV1 "github.com/fluxcd/kustomize-controller/api/v1"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"sigs.k8s.io/controller-runtime/pkg/client/fake"
	"sigs.k8s.io/controller-runtime/pkg/scheme"
	"testing"
)

func TestKustomizations(t *testing.T) {
	testCases := []testCase{
		{
			name:          "valid",
			schemeBuilder: fluxKustomizeV1.SchemeBuilder,
			objects: []runtime.Object{
				&fluxKustomizeV1.Kustomization{
					TypeMeta: v1.TypeMeta{
						Kind:       "Kustomization",
						APIVersion: "",
					},
					ObjectMeta: v1.ObjectMeta{
						Name:      "foo",
						Namespace: "bar",
					},
					Status: fluxKustomizeV1.KustomizationStatus{
						Conditions: []v1.Condition{
							{Type: "Ready", Status: "True"},
						},
					},
				},
			},
			wantErr: assert.NoError,
			want: Resources{{
				Name:      "foo",
				Namespace: "bar",
				Kind:      "Kustomization",
				Conditions: map[string]string{
					"ready": "True",
				},
			}},
		},
		{
			name:          "empty",
			schemeBuilder: fluxKustomizeV1.SchemeBuilder,
			wantErr:       assert.NoError,
			want:          Resources(nil),
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
			k := lister{client: c, list: getKustomizations}

			resources, err := k.List(context.Background())
			tt.wantErr(t, err)
			if err == nil {
				assert.Equal(t, tt.want, resources)
			}
		})
	}
}
