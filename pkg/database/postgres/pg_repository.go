package postgres

import (
	"context"
	"errors"

	"github.com/loongkirin/go-family-finance/pkg/database"
	"github.com/loongkirin/go-family-finance/pkg/database/repository"
	"gorm.io/gorm"
)

type PostgresRepository[T any] struct {
	dbContext database.DbContext
}

func NewRepository[T any](dbContext database.DbContext) repository.Repository[T] {
	return &PostgresRepository[T]{
		dbContext: dbContext,
	}
}

func (r *PostgresRepository[T]) FindById(ctx context.Context, id string) (*T, error) {
	data := new(T)
	err := r.dbContext.GetSlaveDb().WithContext(ctx).Where("id=?", id).First(&data).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, gorm.ErrRecordNotFound
		}
		return nil, err
	}
	return data, nil
}

func (r *PostgresRepository[T]) Query(ctx context.Context, query *database.DbQuery) ([]T, error) {
	datas := []T{}
	whereClaues, values, order := query.GetWhereClause()
	offset := (query.PageNumber - 1) * query.PageSize
	err := r.dbContext.GetSlaveDb().WithContext(ctx).Where(whereClaues, values...).Order(order).Offset(offset).Limit(query.PageSize + 1).Find(&datas).Error
	if err != nil {
		return nil, err
	}
	return datas, nil
}

func (r *PostgresRepository[T]) Create(ctx context.Context, data *T) (*T, error) {
	err := r.dbContext.GetMasterDb().WithContext(ctx).Create(data).Error
	if err != nil {
		return nil, err
	}
	return data, nil
}

func (r *PostgresRepository[T]) Update(ctx context.Context, data *T) (*T, error) {
	err := r.dbContext.GetMasterDb().WithContext(ctx).Save(data).Error
	if err != nil {
		return nil, err
	}
	return data, nil
}

func (r *PostgresRepository[T]) Delete(ctx context.Context, id string) (bool, error) {
	data, err := r.FindById(ctx, id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return true, nil
		}
		return false, err
	}
	err = r.dbContext.GetMasterDb().WithContext(ctx).Delete(data).Error
	if err != nil {
		return false, err
	}
	return true, nil
}
