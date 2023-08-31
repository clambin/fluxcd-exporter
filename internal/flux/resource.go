package flux

import (
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"log/slog"
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

func (r Resource) LogValue() slog.Value {
	conditions := make([]string, 0, len(r.Conditions))
	for key := range r.Conditions {
		conditions = append(conditions, key)
	}
	slices.Sort(conditions)

	grp := make([]any, len(conditions))
	for index, key := range conditions {
		grp[index] = slog.String(key, r.Conditions[key])
	}

	return slog.GroupValue(
		slog.String("kind", r.Kind),
		slog.String("namespace", r.Namespace),
		slog.String("name", r.Name),
		slog.Group("conditions", grp...),
	)
}
