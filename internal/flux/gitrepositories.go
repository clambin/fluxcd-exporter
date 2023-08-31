package flux

import (
	"context"
	"fmt"
	fluxSourceV1 "github.com/fluxcd/source-controller/api/v1"
	"k8s.io/client-go/rest"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

func GitRepositories(cfg *rest.Config) Lister {
	return lister{
		client: makeClient(cfg, fluxSourceV1.SchemeBuilder),
		list:   getGitRepositories,
	}
}

func getGitRepositories(ctx context.Context, c client.Client) (Resources, error) {
	var fluxResources Resources
	var opts client.ListOptions

	for {
		var resources fluxSourceV1.GitRepositoryList
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
