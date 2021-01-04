package manager

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"

	"github.com/huzhongqing/qelog/pkg/common/model"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"

	"github.com/huzhongqing/qelog/pkg/httputil"

	"github.com/huzhongqing/qelog/pkg/common/entity"
)

func (srv *Service) MetricsCount(ctx context.Context, out *entity.MetricsCountResp) error {
	dbStats, err := srv.mongoUtil.DBStats(ctx)
	if err != nil {
		return httputil.ErrSystemException.MergeError(err)
	}
	//now := time.Now()
	//today := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, time.Local)
	//mc, err := srv.store.MetricsModuleCountByDate(ctx, today)
	//if err != nil {
	//	return httputil.ErrSystemException.MergeError(err)
	//}
	//
	//mCount := entity.ModuleCount{
	//	Modules:     mc.Modules,
	//	Numbers:     mc.Numbers,
	//	LoggingSize: mc.LoggingSize,
	//}
	dbCount := entity.DBCount{
		DBName:      dbStats.DB,
		Collections: dbStats.Collections,
		DataSize:    dbStats.DataSize,
		StorageSize: dbStats.StorageSize,
		IndexSize:   dbStats.IndexSize,
		Objects:     dbStats.Objects,
		Indexs:      dbStats.Indexes,
	}

	//out.ModuleCount = mCount
	out.DBCount = dbCount

	return nil
}

func (srv *Service) MetricsModuleList(ctx context.Context, in *entity.MetricsModuleListReq, out *entity.ListResp) error {
	y, m, d := time.Unix(in.DateTsSec, 0).Date()
	date := time.Date(y, m, d, 0, 0, 0, 0, time.Local)

	filter := bson.M{
		"created_date": date.UTC(),
	}

	if in.ModuleName != "" {
		filter["module_name"] = primitive.Regex{
			Pattern: in.ModuleName,
			Options: "i",
		}
	}

	docs := make([]*model.ModuleMetrics, 0, in.Limit)
	opt := options.Find()
	in.SetPage(opt)
	opt.SetSort(bson.M{"number": -1})
	opt.SetProjection(bson.M{"sections": 0})
	c, err := srv.store.FindMetricsModuleList(ctx, filter, &docs, opt)
	if err != nil {
		return httputil.ErrSystemException.MergeError(err)
	}
	out.Count = c

	list := make([]*entity.MetricsModuleList, 0, len(docs))
	for _, v := range docs {
		d := &entity.MetricsModuleList{
			ModuleName:   v.ModuleName,
			Number:       v.Number,
			Size:         v.Size,
			CreatedTsSec: v.CreatedDate.Unix(),
		}
		list = append(list, d)
	}

	out.List = list

	return nil
}
