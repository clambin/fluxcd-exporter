package flux

import (
	"context"
	"fmt"
	fluxHelmV2Beta1 "github.com/fluxcd/helm-controller/api/v2beta1"
	"k8s.io/client-go/rest"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

func HelmReleases(cfg *rest.Config) Lister {
	return lister{
		client: makeClient(cfg, fluxHelmV2Beta1.SchemeBuilder),
		list:   getHelmReleases,
	}
}

func getHelmReleases(ctx context.Context, c client.Client) ([]Resource, error) {
	var fluxResources []Resource
	var opts client.ListOptions

	for {
		var resources fluxHelmV2Beta1.HelmReleaseList

		err := c.List(ctx, &resources, &opts)
		if err != nil {
			return nil, fmt.Errorf("list: %w", err)
		}

		for _, resource := range resources.Items {
			fluxResources = append(fluxResources, newResource(
				resource.ObjectMeta.GetName(),
				resource.ObjectMeta.GetNamespace(),
				resource.TypeMeta.Kind,
				resource.GetConditions(),
			))
		}

		if resources.RemainingItemCount == nil || *resources.RemainingItemCount == 0 {
			break
		}
		opts.Continue = resources.Continue
	}
	return fluxResources, nil
}
