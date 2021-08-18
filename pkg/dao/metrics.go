package dao

import (
	"context"
	"fmt"
	"github.com/bbdshow/bkit/errc"
	"github.com/bbdshow/qelog/pkg/model"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
	"time"
)

func (d *Dao) IncrModuleMetrics(ctx context.Context, in *model.MetricsState) error {
	filter := bson.M{
		"module_name":  in.ModuleName,
		"created_date": in.Date,
	}

	opt := options.Update()
	opt.SetUpsert(true)

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
