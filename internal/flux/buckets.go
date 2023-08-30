package flux

import (
	"context"
	"fmt"
	fluxSourceV1Beta2 "github.com/fluxcd/source-controller/api/v1beta2"
	"k8s.io/client-go/rest"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

func Buckets(cfg *rest.Config) Lister {
	return lister{
		client: makeClient(cfg, fluxSourceV1Beta2.SchemeBuilder),
		list:   getBuckets,
	}
}

func getBuckets(ctx context.Context, c client.Client) (Resources, error) {
	var fluxResources Resources
	opts := client.ListOptions{Namespace: ""}

	for {
		var resources fluxSourceV1Beta2.BucketList
		if err := c.List(ctx, &resources, &opts); err != nil {
			return nil, fmt.Errorf("list: %w", err)
		}

		for _, resource := range resources.Items {
			fluxResources = append(fluxResources, newResource(
				resource.Name,
				resource.Namespace,
				resource.Kind,
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
