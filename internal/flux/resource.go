package flux

import (
	"fmt"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"slices"
	"strings"
)

type Resource struct {
	Name       string
	Namespace  string
	Kind       string
	Conditions map[string]string
}

func newResource(name, namespace, kind string, conditions []v1.Condition) Resource {
	cond := make(map[string]string)
	for _, condition := range conditions {
		cond[strings.ToLower(condition.Type)] = string(condition.Status)
	}
	return Resource{
		Name:       name,
		Namespace:  namespace,
		Kind:       kind,
		Conditions: cond,
	}
}

func (r Resource) String() string {
	var parts []string
	parts = append(parts, fmt.Sprintf(`name="%s"`, r.Name))
	parts = append(parts, fmt.Sprintf(`namespace="%s"`, r.Namespace))
	parts = append(parts, fmt.Sprintf(`kind="%s"`, r.Kind))
	keys := make([]string, 0, len(r.Conditions))
	for key := range r.Conditions {
		keys = append(keys, key)
	}
	slices.Sort(keys)
	for _, key := range keys {
		parts = append(parts, fmt.Sprintf(`%s="%s"`, strings.ToLower(key), r.Conditions[key]))
	}

	return "gotk_resource_info{" + strings.Join(parts, ", ") + "}"

}

type Resources []Resource

func (f Resources) String() string {
	var resources []string
	for _, r := range f {
		resources = append(resources, r.String())
	}
	return strings.Join(resources, "\n")
}
