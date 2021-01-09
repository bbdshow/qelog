package storage

import (
	"context"
	"errors"

	"github.com/huzhongqing/qelog/libs/logs"
	"github.com/huzhongqing/qelog/libs/mongo"
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

func (store *Store) ListAllCollectionNames(ctx context.Context) ([]string, error) {
	names, err := store.database.ListAllCollectionNames(ctx)
	return names, handlerError(err)
}

func (store *Store) UpsertCollectionIndexMany(indexs []mongo.Index) error {
	err := store.database.UpsertCollectionIndexMany(indexs)
	return handlerError(err)
}

// 记录所有的数据库操作错误
func handlerError(err error) error {
	if err != nil {
		if logs.Qezap != nil {
			logs.Qezap.Error("Store Operation Error", zap.String("error", err.Error()))
		}
	}
	return err
}
