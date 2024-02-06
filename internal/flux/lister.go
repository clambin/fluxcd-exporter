package flux

import (
	"context"
	"fmt"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/client-go/rest"
	"log/slog"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

type Lister interface {
	List(ctx context.Context) ([]Resource, error)
}

type ListerFunc func(ctx context.Context) ([]Resource, error)

func (l ListerFunc) List(ctx context.Context) ([]Resource, error) {
	return l(ctx)
}

type listFunc func(context.Context, client.Client, *client.ListOptions) ([]Resource, int64, string, error)

type lister struct {
	client client.Client
	list   listFunc
	logger *slog.Logger
}

func (l lister) List(ctx context.Context) ([]Resource, error) {
	var fluxResources []Resource
	var opts client.ListOptions

	//opts.Limit = 5
	for {
		resources, remaining, cont, err := l.list(ctx, l.client, &opts)
		if err != nil {
			return nil, fmt.Errorf("list: %w", err)
		}

		l.logger.Debug("flux custom resources found", "len", len(resources))

		fluxResources = append(fluxResources, resources...)

		if remaining == 0 {
			break
		}

		l.logger.Debug("more flux custom resources to be retrieved", "count", remaining)
		opts.Continue = cont
	}
	return fluxResources, nil
}

type SchemeBuilder interface {
	Build() (*runtime.Scheme, error)
}

func newLister(cfg *rest.Config, builder SchemeBuilder, listFunc listFunc, logger *slog.Logger) lister {
	scheme, err := builder.Build()
	if err != nil {
		panic(fmt.Errorf("build scheme: %w", err))
	}
	c, err := client.New(cfg, client.Options{Scheme: scheme})
	if err != nil {
		panic(fmt.Errorf("new runtime client: %w", err))
	}

	return lister{
		client: c,
		list:   listFunc,
		logger: logger,
	}
}

type windowedLister interface {
	GetRemainingItemCount() *int64
}

func getRemaining(resources windowedLister) int64 {
	remaining := resources.GetRemainingItemCount()
	if remaining == nil {
		return 0
	}
	return *remaining
}
