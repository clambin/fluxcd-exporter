package flux

import (
	fluxHelmV2Beta2 "github.com/fluxcd/helm-controller/api/v2beta2"
	fluxKustomizeV1 "github.com/fluxcd/kustomize-controller/api/v1"
	fluxSourceV1 "github.com/fluxcd/source-controller/api/v1"
	fluxSourceV1Beta2 "github.com/fluxcd/source-controller/api/v1beta2"
	"github.com/stretchr/testify/assert"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/client-go/rest"
	"sigs.k8s.io/controller-runtime/pkg/scheme"
	"testing"
)

type testCase struct {
	name          string
	schemeBuilder *scheme.Builder
	objects       []runtime.Object
	wantErr       assert.ErrorAssertionFunc
	want          []Resource
}

func Test_makeClient(t *testing.T) {
	for _, builder := range []SchemeBuilder{
		fluxSourceV1Beta2.SchemeBuilder,
		fluxSourceV1.SchemeBuilder,
		fluxHelmV2Beta2.SchemeBuilder,
		fluxKustomizeV1.SchemeBuilder,
	} {
		assert.NotNil(t, makeClient(&rest.Config{}, builder))
		assert.Panics(t, func() {
			_ = makeClient(nil, builder)
		})
	}
}
