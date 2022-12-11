package dao

import (
	"context"
	"fmt"
	"time"

	"github.com/bbdshow/bkit/db/mongo"
	"github.com/bbdshow/bkit/errc"
	"github.com/bbdshow/bkit/util/itime"
	"github.com/bbdshow/qelog/pkg/model"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// IncrModuleMetrics module metrics info insert
func (d *Dao) IncrModuleMetrics(ctx context.Context, in *model.MetricsState) error {
	filter := bson.M{
		"module_name":  in.ModuleName,
		"created_date": in.Date,
	}

	opt := options.Update().SetUpsert(true)

	fields := bson.M{
		"number": in.Number,
		"size":   in.Size,
		fmt.Sprintf("sections.%d.sum", in.Section): in.Number,
	}
	for k, v := range in.Levels {
		fields[fmt.Sprintf("sections.%d.levels.%d", in.Section, k.Int32())] = v
	}
	for k, v := range in.IPs {
		fields[fmt.Sprintf("sections.%d.ips.%s", in.Section, k)] = v
	}

	update := bson.M{
		"$inc": fields,
	}
	if _, err := d.adminInst.Collection(model.CNModuleMetrics).UpdateOne(ctx, filter, update, opt); err != nil {
		return errc.WithStack(err)
	}

	return nil
}

// GetModuleMetricsCountByDate aggregate module this time period logging count and data size
func (d *Dao) GetModuleMetricsCountByDate(ctx context.Context, date time.Time) (*model.ModuleCount, error) {
	coll := d.adminInst.Collection(model.CNModuleMetrics)

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
	out := &model.ModuleCount{}
	if len(val) > 0 {
		out.Numbers = val[0].Numbers
		out.LoggingSize = val[0].Sizes
		out.Modules = val[0].Count
	}
	return out, nil
}

// GetDBStats common db CRUD operation
func (d *Dao) GetDBStats(ctx context.Context, filter bson.M) (bool, *model.DBStats, error) {
	doc := &model.DBStats{}
	exists, err := d.adminInst.FindOne(ctx, model.CNDBStats, filter, doc)
	return exists, doc, errc.WithStack(err)
}

// FindDBStats common db CRUD operation
func (d *Dao) FindDBStats(ctx context.Context, filter bson.M) ([]*model.DBStats, error) {
	docs := make([]*model.DBStats, 0)
	opt := options.Find().SetSort(bson.M{"_id": -1})
	err := d.adminInst.Find(ctx, model.CNDBStats, filter, &docs, opt)
	return docs, errc.WithStack(err)
}

// UpsertDBStats common db CRUD operation
func (d *Dao) UpsertDBStats(ctx context.Context, in *model.DBStats) error {
	filter := bson.M{
		"host": in.Host,
		"db":   in.DB,
	}

	opt := options.Update().SetUpsert(true)

	update := bson.M{
		"$set": bson.M{
			"collections":  in.Collections,
			"objects":      in.Objects,
			"data_size":    in.DataSize,
			"storage_size": in.StorageSize,
			"index_size":   in.IndexSize,
			"updated_at":   time.Now(),
		},
		"$setOnInsert": bson.M{
			"created_at": time.Now(),
		},
	}
	_, err := d.adminInst.Collection(model.CNDBStats).UpdateOne(ctx, filter, update, opt)
	return err
}

// GetCollStats common db CRUD operation
func (d *Dao) GetCollStats(ctx context.Context, filter bson.M) (bool, *model.CollStats, error) {
	doc := &model.CollStats{}
	exists, err := d.adminInst.FindOne(ctx, model.CNCollStats, filter, doc)
	return exists, doc, errc.WithStack(err)
}

// FindCollStats common db CRUD operation
func (d *Dao) FindCollStats(ctx context.Context, filter bson.M) ([]*model.CollStats, error) {
	docs := make([]*model.CollStats, 0)
	opt := options.Find().SetSort(bson.M{"_id": -1})
	err := d.adminInst.Find(ctx, model.CNCollStats, filter, &docs, opt)
	return docs, errc.WithStack(err)
}

// UpsertCollStats common db CRUD operation
func (d *Dao) UpsertCollStats(ctx context.Context, in *model.CollStats) error {
	filter := bson.M{
		"module_name": in.ModuleName,
		"host":        in.Host,
		"db":          in.DB,
		"name":        in.Name,
	}

	opt := options.Update().SetUpsert(true)

	update := bson.M{
		"$set": bson.M{
			"size":             in.Size,
			"count":            in.Count,
			"avg_obj_size":     in.AvgObjSize,
			"storage_size":     in.StorageSize,
			"capped":           in.Capped,
			"total_index_size": in.TotalIndexSize,
			"index_sizes":      in.IndexSizes,
			"updated_at":       time.Now(),
		},
		"$setOnInsert": bson.M{
			"created_at": time.Now(),
		},
	}
	_, err := d.adminInst.Collection(model.CNCollStats).UpdateOne(ctx, filter, update, opt)
	return err
}

// ReadDBStats exec db command
func (d *Dao) ReadDBStats(ctx context.Context, db *mongo.Database) (mongo.DBStatsResp, error) {
	c := mongo.NewCommand(db)
	return c.DBStats(ctx)
}

// ReadCollStats exec db command
func (d *Dao) ReadCollStats(ctx context.Context, db *mongo.Database, collection string) (mongo.CollStatsResp, error) {
	c := mongo.NewCommand(db)
	return c.CollStats(ctx, collection)
}

// FindMetricsModule common db CRUD operation
func (d *Dao) FindMetricsModule(ctx context.Context, filter bson.M) ([]*model.ModuleMetrics, error) {
	docs := make([]*model.ModuleMetrics, 0)
	err := d.adminInst.Find(ctx, model.CNModuleMetrics, filter, &docs)
	return docs, errc.WithStack(err)
}

// FindMetricsModuleList common db CRUD operation
func (d *Dao) FindMetricsModuleList(ctx context.Context, in *model.MetricsModuleListReq) (int64, []*model.ModuleMetrics, error) {
	date := itime.UnixSecToDate(in.DateTsSec)
	filter := bson.M{
		"created_date": date,
	}
	if in.ModuleName != "" {
		filter["module_name"] = primitive.Regex{
			Pattern: in.ModuleName,
			Options: "i",
		}
	}

	opt := in.SetPage(options.Find()).SetSort(bson.M{"number": -1}).SetProjection(bson.M{"sections": 0})

	docs := make([]*model.ModuleMetrics, 0)
	c, err := d.adminInst.FindCount(ctx, model.CNModuleMetrics, filter, &docs, opt, nil)
	return c, docs, errc.WithStack(err)
}
