package flux

import (
	"cmp"
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
	return slog.GroupValue(
		slog.String("kind", r.Kind),
		slog.String("namespace", r.Namespace),
		slog.String("name", r.Name),
		slog.Group("conditions", logConditions(r.Conditions)...),
	)
}

func logConditions(conditions map[string]string) []any {
	attribs := make([]any, 0, len(conditions))

	for key, val := range conditions {
		attribs = append(attribs, slog.String(key, val))
	}

	slices.SortFunc(attribs, func(a, b any) int {
		return cmp.Compare(a.(slog.Attr).Key, b.(slog.Attr).Key)
	})

	return attribs
}
