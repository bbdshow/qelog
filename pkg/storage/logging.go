package storage

//import (
//	"context"
//	"github.com/bbdshow/qelog/infra/mongo"
//	"github.com/bbdshow/qelog/pkg/common/model"
//
//	"github.com/bbdshow/qelog/infra/logs"
//	"go.mongodb.org/mongo-driver/bson"
//	"go.mongodb.org/mongo-driver/mongo/options"
//	"go.uber.org/zap"
//)
//
//type Logging struct {
//	db *mongo.Database
//}
//
//func NewLogging(db *mongo.Database) *Logging {
//	return &Logging{db: db}
//}

//func (l *Logging) InsertManyLogging(ctx context.Context, collectionName string, docs []interface{}) error {
//	_, err := l.db.Collection(collectionName).InsertMany(ctx, docs)
//	return err
//}

//func (l *Logging) FindCountLoggingList(ctx context.Context, collectionName string, filter bson.M, countLimit int64, opts ...*options.FindOptions) (int64, []*model.Logging, error) {
//	// 异步统计Count
//	calcCount := func(ctx context.Context) (int64, error) {
//		countOpt := options.Count()
//		countOpt.SetLimit(countLimit)
//		c, err := l.db.Collection(collectionName).CountDocuments(ctx, filter, countOpt)
//		return c, err
//	}
//	countResp := make(chan int64, 1)
//	go func() {
//		c, err := calcCount(ctx)
//		if err != nil {
//			logs.Qezap.Error("FindLoggingCount", zap.Error(err))
//		}
//		countResp <- c
//	}()
//	docs := make([]*model.Logging, 0)
//	err := l.db.Find(ctx, l.db.Collection(collectionName), filter, &docs, opts...)
//	if err != nil {
//		return 0, docs, err
//	}
//
//	select {
//	case c := <-countResp:
//		return c, docs, nil
//	}
//}
