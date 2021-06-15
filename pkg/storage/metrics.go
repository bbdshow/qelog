package storage

import (
	"context"
	"github.com/bbdshow/qelog/infra/mongo"
	"time"

	"github.com/bbdshow/qelog/pkg/common/entity"
	"github.com/bbdshow/qelog/pkg/common/model"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type ModuleMetrics struct {
	db *mongo.Database
}

func NewModuleMetrics(db *mongo.Database) *ModuleMetrics {
	return &ModuleMetrics{db: db}
}

func (mm *ModuleMetrics) UpdateModuleMetrics(ctx context.Context, filter, update bson.M, opts ...*options.UpdateOptions) error {
	_, err := mm.db.Collection(model.CollectionNameModuleMetrics).UpdateOne(ctx, filter, update, opts...)
	return err
}

func (mm *ModuleMetrics) GetModuleMetricsCountByDate(ctx context.Context, date time.Time) (*entity.ModuleCount, error) {
	coll := mm.db.Collection(model.CollectionNameModuleMetrics)

	pipeline := []bson.D{
		{
			{Key: "$match", Value: bson.M{"created_date": date}},
		},
		{
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
		{
			{Key: "$project", Value: bson.M{
				"numbers": 1,
				"sizes":   1,
				"count":   1,
			}},
		},
	}
	cursor, err := coll.Aggregate(ctx, pipeline)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	type counts = struct {
		Numbers int64 `bson:"numbers"`
		Sizes   int64 `bson:"sizes"`
		Count   int64 `bson:"count"`
	}
	val := make([]counts, 0)
	if err := cursor.All(ctx, &val); err != nil {
		return nil, err
	}
	out := &entity.ModuleCount{}
	if len(val) > 0 {
		out.Numbers = val[0].Numbers
		out.LoggingSize = val[0].Sizes
		out.Modules = val[0].Count
	}
	return out, nil
}

func (mm *ModuleMetrics) FindCountModuleMetrics(ctx context.Context, filter bson.M, opts ...*options.FindOptions) (int64, []*model.ModuleMetrics, error) {
	docs := make([]*model.ModuleMetrics, 0)
	c, err := mm.db.FindAndCount(ctx, mm.db.Collection(model.CollectionNameModuleMetrics), filter, &docs, opts...)
	return c, docs, err
}

func (mm *ModuleMetrics) InsertOneDBStats(ctx context.Context, doc *model.DBStats) error {
	_, err := mm.db.Collection(doc.CollectionName()).InsertOne(ctx, doc)
	return err
}

func (mm *ModuleMetrics) FindOneDBStats(ctx context.Context, filter bson.M, doc *model.DBStats, opts ...*options.FindOneOptions) (bool, error) {
	ok, err := mm.db.FindOne(ctx, mm.db.Collection(doc.CollectionName()), filter, doc, opts...)
	return ok, err
}

func (mm *ModuleMetrics) InsertManyCollStats(ctx context.Context, docs []interface{}) error {
	_, err := mm.db.Collection(model.CollectionNameCollStats).InsertMany(ctx, docs)
	return err
}

func (mm *ModuleMetrics) FindCollStats(ctx context.Context, filter bson.M, opts ...*options.FindOptions) ([]*model.CollStats, error) {
	docs := make([]*model.CollStats, 0)
	err := mm.db.Find(ctx, mm.db.Collection(model.CollectionNameCollStats), filter, &docs, opts...)
	return docs, err
}
