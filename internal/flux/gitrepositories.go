package flux

import (
	"context"
	"fmt"
	fluxSourceV1 "github.com/fluxcd/source-controller/api/v1"
	"k8s.io/client-go/rest"
	"log/slog"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

func GitRepositories(cfg *rest.Config, logger *slog.Logger) Lister {
	return lister{
		client: makeClient(cfg, fluxSourceV1.SchemeBuilder),
		list:   getGitRepositories,
		logger: logger.With("custom_resource", "gitRepositories"),
	}
}

func getGitRepositories(ctx context.Context, c client.Client, opts *client.ListOptions) ([]Resource, int64, string, error) {
	var resources fluxSourceV1.GitRepositoryList

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
