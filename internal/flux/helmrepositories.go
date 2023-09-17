package flux

import (
	"context"
	"fmt"
	fluxSourceV1Beta2 "github.com/fluxcd/source-controller/api/v1beta2"
	"k8s.io/client-go/rest"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

func HelmRepositories(cfg *rest.Config) Lister {
	return lister{
		client: makeClient(cfg, fluxSourceV1Beta2.SchemeBuilder),
		list:   getHelmRepositories,
	}
}

func getHelmRepositories(ctx context.Context, c client.Client) ([]Resource, error) {
	var fluxResources []Resource
	var opts client.ListOptions

	for {
		var resources fluxSourceV1Beta2.HelmRepositoryList
		if err := c.List(ctx, &resources, &opts); err != nil {
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
