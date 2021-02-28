package storage

import (
	"context"

	"github.com/huzhongqing/qelog/infra/logs"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.uber.org/zap"
)

func (store *Store) InsertManyLogging(ctx context.Context, name string, docs []interface{}) error {
	_, err := store.database.Collection(name).InsertMany(ctx, docs)
	return handlerError(err)
}

func (store *Store) FindLoggingList(ctx context.Context, collectionName string, filter bson.M, countLimit int64, result interface{}, opt *options.FindOptions) (int64, error) {
	// 异步统计Count
	calcCount := func(ctx context.Context) (int64, error) {
		countOpt := options.Count()
		countOpt.SetLimit(countLimit)
		c, err := store.database.Collection(collectionName).CountDocuments(ctx, filter, countOpt)
		return c, handlerError(err)
	}
	countResp := make(chan int64, 1)
	go func() {
		c, err := calcCount(ctx)
		if err != nil {
			logs.Qezap.Error("FindLoggingCount", zap.Error(err))
		}
		countResp <- c
	}()

	err := store.database.Find(ctx, store.database.Collection(collectionName), filter, result, opt)
	if err != nil {
		return 0, handlerError(err)
	}

	select {
	case c := <-countResp:
		return c, nil
	}
}
