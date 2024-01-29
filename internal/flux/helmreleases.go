package flux

import (
	"context"
	"fmt"
	fluxHelmV2Beta1 "github.com/fluxcd/helm-controller/api/v2beta1"
	"k8s.io/client-go/rest"
	"log/slog"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

func HelmReleases(cfg *rest.Config, logger *slog.Logger) Lister {
	return lister{
		client: makeClient(cfg, fluxHelmV2Beta1.SchemeBuilder),
		list:   getHelmReleases,
		logger: logger.With("custom_resource", "helmReleases"),
	}
}

func getHelmReleases(ctx context.Context, c client.Client, opts *client.ListOptions) ([]Resource, int64, string, error) {
	var resources fluxHelmV2Beta1.HelmReleaseList

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
