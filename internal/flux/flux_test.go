package flux

import (
	"github.com/stretchr/testify/assert"
	"k8s.io/apimachinery/pkg/runtime"
	"sigs.k8s.io/controller-runtime/pkg/scheme"
)

type testCase struct {
	name          string
	schemeBuilder *scheme.Builder
	objects       []runtime.Object
	wantErr       assert.ErrorAssertionFunc
	want          []Resource
}
