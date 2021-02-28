package storage

import (
	"context"
	"time"

	"github.com/huzhongqing/qelog/pkg/common/entity"
	"github.com/huzhongqing/qelog/pkg/common/model"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func (store *Store) UpsertModuleMetrics(ctx context.Context, filter, update bson.M, opt *options.UpdateOptions) error {
	_, err := store.database.Collection(model.CollectionNameModuleMetrics).UpdateOne(ctx, filter, update, opt)
	return handlerError(err)
}

func (store *Store) MetricsModuleCountByDate(ctx context.Context, date time.Time) (*entity.ModuleCount, error) {
	coll := store.database.Collection(model.CollectionNameModuleMetrics)

	pipeline := mongo.Pipeline{
		bson.D{
			{Key: "$match", Value: bson.M{"created_date": date}},
		},
		bson.D{
			{Key: "$group", Value: bson.M{
				"_id": bson.D{
					{"$dateToString", bson.D{
						{"format", "%Y-%m-%d%H"},
						{"date", "$created_date"},
					}},
				},
				"numbers": bson.M{"$sum": "$number"},
				"sizes":   bson.M{"$sum": "$size"},
				"count":   bson.M{"$sum": 1},
			}},
		},
		bson.D{
			{Key: "$project", Value: bson.M{
				"numbers": 1,
				"sizes":   1,
				"count":   1,
			}},
		},
	}
	cursor, err := coll.Aggregate(ctx, pipeline)
	if err != nil {
		return nil, handlerError(err)
	}
	defer cursor.Close(ctx)

	type counts = struct {
		Numbers int64 `bson:"numbers"`
		Sizes   int64 `bson:"sizes"`
		Count   int64 `bson:"count"`
	}
	val := make([]counts, 0)
	if err := cursor.All(ctx, &val); err != nil {
		return nil, handlerError(err)
	}
	out := &entity.ModuleCount{}
	if len(val) > 0 {
		out.Numbers = val[0].Numbers
		out.LoggingSize = val[0].Sizes
		out.Modules = val[0].Count
	}
	return out, nil
}

func (store *Store) FindMetricsModuleList(ctx context.Context, filter bson.M, result interface{}, opt *options.FindOptions) (int64, error) {
	c, err := store.database.FindAndCount(ctx, store.database.Collection(model.CollectionNameModuleMetrics), filter, result, opt)
	return c, handlerError(err)
}

func (store *Store) FindModuleMetrics(ctx context.Context, filter bson.M, result interface{}, opt *options.FindOptions) error {
	err := store.database.Find(ctx, store.database.Collection(model.CollectionNameModuleMetrics), filter, result, opt)
	return handlerError(err)
}

func (store *Store) InsertOneDBStats(ctx context.Context, doc *model.DBStats) error {
	_, err := store.database.Collection(doc.CollectionName()).InsertOne(ctx, doc)
	return handlerError(err)
}

func (store *Store) FindOneDBStats(ctx context.Context, filter bson.M, doc *model.DBStats, opt *options.FindOneOptions) (bool, error) {
	ok, err := store.database.FindOne(ctx, store.database.Collection(doc.CollectionName()), filter, doc, opt)
	return ok, handlerError(err)
}

func (store *Store) InsertManyCollStats(ctx context.Context, docs []interface{}) error {
	_, err := store.database.Collection(model.CollectionNameCollStats).InsertMany(ctx, docs)
	return handlerError(err)
}

func (store *Store) FindCollStats(ctx context.Context, filter bson.M, doc interface{}, opt *options.FindOptions) error {
	err := store.database.Find(ctx, store.database.Collection(model.CollectionNameCollStats), filter, doc, opt)
	return handlerError(err)
}
