package types

//go:generate mockgen -destination=./mocks/mock_store.go -package=mocks github.com/go-obvious/server-example/types Store

import (
	"context"
)

type Store interface {
	All(context.Context, *string) (ProductRange, error)
	Get(context.Context, string) (*Product, error)
	Put(context.Context, Product) error
	Delete(context.Context, string) error
}
