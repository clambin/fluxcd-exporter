package flux

import (
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"log/slog"
	"strings"
)

type Resource struct {
	Name       string
	Namespace  string
	Kind       string
	Conditions map[string]string
}

type Resources []Resource

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
	var grp []any
	for key, val := range r.Conditions {
		grp = append(grp, slog.String(key, val))
	}
	return slog.GroupValue(
		slog.String("kind", r.Kind),
		slog.String("namespace", r.Namespace),
		slog.String("name", r.Name),
		slog.Group("conditions", grp...),
	)
}
