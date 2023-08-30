package flux

import (
	"context"
	"fmt"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/client-go/rest"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

type Lister interface {
	List(ctx context.Context) (Resources, error)
}

type ListerFunc func(ctx context.Context) (Resources, error)

func (l ListerFunc) List(ctx context.Context) (Resources, error) {
	return l(ctx)
}

type lister struct {
	client client.Client
	list   func(context.Context, client.Client) (Resources, error)
}

func (l lister) List(ctx context.Context) (Resources, error) {
	return l.list(ctx, l.client)
}

type SchemeBuilder interface {
	Build() (*runtime.Scheme, error)
}

func makeClient(cfg *rest.Config, builder SchemeBuilder) client.Client {
	scheme, err := builder.Build()
	if err != nil {
		panic(fmt.Errorf("build scheme: %w", err))
	}
	c, err := client.New(cfg, client.Options{Scheme: scheme})
	if err != nil {
		panic(fmt.Errorf("new runtime client: %w", err))
	}
	return c
}
