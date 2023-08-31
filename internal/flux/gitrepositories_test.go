package flux

import (
	"context"
	fluxSourceV1 "github.com/fluxcd/source-controller/api/v1"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"sigs.k8s.io/controller-runtime/pkg/client/fake"
	"sigs.k8s.io/controller-runtime/pkg/scheme"
	"testing"
)

func TestGitRepositories(t *testing.T) {
	testCases := []testCase{
		{
			name:          "valid",
			schemeBuilder: fluxSourceV1.SchemeBuilder,
			objects: []runtime.Object{
				&fluxSourceV1.GitRepository{
					TypeMeta: v1.TypeMeta{
						Kind:       "GitRepository",
						APIVersion: "",
					},
					ObjectMeta: v1.ObjectMeta{
						Name:      "foo",
						Namespace: "bar",
					},
					Status: fluxSourceV1.GitRepositoryStatus{
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
				Kind:      "GitRepository",
				Conditions: map[string]string{
					"ready": "True",
				},
			}},
		},
		{
			name:          "empty",
			schemeBuilder: fluxSourceV1.SchemeBuilder,
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
			k := lister{client: c, list: getGitRepositories}

			resources, err := k.List(context.Background())
			tt.wantErr(t, err)
			if err == nil {
				assert.Equal(t, tt.want, resources)
			}
		})
	}
}
