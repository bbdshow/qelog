package storage

import (
	"context"

	"github.com/huzhongqing/qelog/libs/logs"
	"go.uber.org/zap"

	"go.mongodb.org/mongo-driver/mongo/options"

	"go.mongodb.org/mongo-driver/bson/primitive"

	"github.com/huzhongqing/qelog/pkg/common/entity"
	"github.com/huzhongqing/qelog/pkg/common/model"
	"go.mongodb.org/mongo-driver/bson"
)

func (store *Store) InsertManyLogging(ctx context.Context, name string, docs []interface{}) error {
	_, err := store.database.Collection(name).InsertMany(ctx, docs)
	return WrapError(err)
}

func (store *Store) FindLoggingList(ctx context.Context, collectionName string, in *entity.FindLoggingListReq) (int64, []*model.Logging, error) {
	filter := bson.M{
		"m": in.ModuleName,
	}

	// 必须存在时间
	filter["ts"] = bson.M{"$gte": in.BeginUnix, "$lt": in.EndUnix}

	if in.ShortMsg != "" {
		// 区分大小写
		filter["s"] = primitive.Regex{
			Pattern: in.ShortMsg,
		}
	}

	if in.Level >= 0 {
		filter["l"] = in.Level
	}

	if in.ConditionOne != "" {
		filter["c1"] = in.ConditionOne
		if in.ConditionTwo != "" {
			filter["c2"] = in.ConditionTwo
			if in.ConditionThree != "" {
				filter["c3"] = in.ConditionThree
			}
		}
	}
	// 异步统计Count
	calcCount := func(ctx context.Context) (int64, error) {
		countOpt := options.Count()
		countOpt.SetLimit(50000)
		return store.database.Collection(collectionName).CountDocuments(ctx, filter, countOpt)
	}
	countResp := make(chan int64, 1)
	go func() {
		c, err := calcCount(ctx)
		if err != nil {
			logs.Qezap.Error("FindLogging", zap.String("count", err.Error()))
		}
		countResp <- c
	}()

	findOpt := options.Find()
	in.SetPage(findOpt)
	findOpt.SetSort(bson.M{"ts": -1})

	docs := make([]*model.Logging, 0, in.PageSize)
	err := store.database.Find(ctx, store.database.Collection(collectionName), filter, &docs, findOpt)
	if err != nil {
		return 0, nil, err
	}

	select {
	case c := <-countResp:
		return c, docs, nil
	}
}
