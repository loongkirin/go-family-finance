package repository

import (
	"context"

	"github.com/loongkirin/go-family-finance/pkg/database"
)

type Repository[T any] interface {
	FindById(ctx context.Context, id string) (*T, error)
	Query(ctx context.Context, query *database.DbQuery) ([]T, error)
	Create(ctx context.Context, data *T) (*T, error)
	Update(ctx context.Context, data *T) (*T, error)
	Delete(ctx context.Context, id string) (bool, error)
}
