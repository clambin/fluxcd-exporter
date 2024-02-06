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
	return newLister(cfg, fluxKustomizeV1.SchemeBuilder, getKustomizations, logger.With("custom_resource", "kustomizations"))
}

func getKustomizations(ctx context.Context, c client.Client, opts *client.ListOptions) ([]Resource, int64, string, error) {
	var resources fluxKustomizeV1.KustomizationList

	if err := c.List(ctx, &resources, opts); err != nil {
		return nil, 0, "", fmt.Errorf("list: %w", err)
	}

	fluxResources := make([]Resource, len(resources.Items))
	for i := range resources.Items {
		fluxResources[i] = newResource(
			resources.Items[i].ObjectMeta.GetName(),
			resources.Items[i].ObjectMeta.GetNamespace(),
			resources.Items[i].TypeMeta.Kind,
			resources.Items[i].GetConditions(),
		)
	}

	return fluxResources, getRemaining(&resources), resources.Continue, nil
}
