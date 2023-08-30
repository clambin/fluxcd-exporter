package flux

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestResources_String(t *testing.T) {
	resources := Resources{
		{
			Name:      "foo",
			Namespace: "bar",
			Kind:      "Kustomization",
			Conditions: map[string]string{
				"ready": "True",
			},
		},
		{
			Name:      "snafu",
			Namespace: "foobar",
			Kind:      "GitRepository",
			Conditions: map[string]string{
				"ready":             "False",
				"artifactinstorage": "True",
			},
		},
	}

	assert.Equal(t, `gotk_resource_info{name="foo", namespace="bar", kind="Kustomization", ready="True"}
gotk_resource_info{name="snafu", namespace="foobar", kind="GitRepository", artifactinstorage="True", ready="False"}`, resources.String())
}
