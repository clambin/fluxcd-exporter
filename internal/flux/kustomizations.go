package flux

import (
	"context"
	"fmt"
	fluxKustomizeV1 "github.com/fluxcd/kustomize-controller/api/v1"
	"k8s.io/client-go/rest"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

func Kustomizations(cfg *rest.Config) Lister {
	return lister{
		client: makeClient(cfg, fluxKustomizeV1.SchemeBuilder),
		list:   getKustomizations,
	}
}

func getKustomizations(ctx context.Context, c client.Client) ([]Resource, error) {
	var fluxResources []Resource
	var opts client.ListOptions

	for {
		var resources fluxKustomizeV1.KustomizationList
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

		remaining := resources.GetRemainingItemCount()
		if remaining == nil || *remaining == 0 {
			break
		}
		opts.Continue = resources.Continue
	}

	return fluxResources, nil
}
