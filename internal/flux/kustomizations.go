package flux

import (
	"context"
	"fmt"
	fluxKustomizeV1 "github.com/fluxcd/kustomize-controller/api/v1"
	"k8s.io/client-go/rest"
	"log/slog"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

func Kustomizations(cfg *rest.Config, logger *slog.Logger) Lister {
	return lister{
		client: makeClient(cfg, fluxKustomizeV1.SchemeBuilder),
		list:   getKustomizations,
		logger: logger.With("custom_resource", "kustomizations"),
	}
}

func getKustomizations(ctx context.Context, c client.Client, opts *client.ListOptions) ([]Resource, int64, string, error) {
	var fluxResources []Resource
	var resources fluxKustomizeV1.KustomizationList

	if err := c.List(ctx, &resources, opts); err != nil {
		return nil, 0, "", fmt.Errorf("list: %w", err)
	}

	for _, resource := range resources.Items {
		fluxResources = append(fluxResources, newResource(
			resource.ObjectMeta.GetName(),
			resource.ObjectMeta.GetNamespace(),
			resource.TypeMeta.Kind,
			resource.GetConditions(),
		))
	}

	return fluxResources, getRemaining(&resources), resources.Continue, nil
}
