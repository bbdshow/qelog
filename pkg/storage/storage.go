package storage

import (
	"context"
	"errors"

	"github.com/huzhongqing/qelog/infra/logs"
	"github.com/huzhongqing/qelog/infra/mongo"
	"go.uber.org/zap"
)

var (
	ErrNotMatched         = errors.New("not matched")
	ErrMainDBNotFound     = errors.New("main db not found")
	ErrShardingDBNotFound = errors.New("sharding db not found")
)

type Store struct {
	database *mongo.Database
}

func New(database *mongo.Database) *Store {
	store := &Store{
		database: database,
	}
	return store
}

func (store *Store) Database() *mongo.Database {
	return store.database
}

// ListCollectionNames
func (store *Store) ListCollectionNames(ctx context.Context, prefix ...string) ([]string, error) {
	names, err := store.database.ListCollectionNames(ctx, prefix...)
	return names, handlerError(err)
}

// UpsertCollectionIndexMany
func (store *Store) UpsertCollectionIndexMany(indexs []mongo.Index) error {
	err := store.database.UpsertCollectionIndexMany(indexs)
	return handlerError(err)
}

func handlerError(err error) error {
	if err != nil {
		if logs.Qezap != nil {
			logs.Qezap.Error("StoreOperationError", zap.String("error", err.Error()))
		}
	}
	return err
}
