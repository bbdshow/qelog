package dao

import (
	"context"
	"fmt"
	"github.com/bbdshow/bkit/db/mongo"
	"github.com/bbdshow/bkit/errc"
	"github.com/bbdshow/bkit/logs"
	apiTypes "github.com/bbdshow/qelog/api/types"
	"github.com/bbdshow/qelog/pkg/model"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.uber.org/zap"
	"strings"
	"time"
)

func (d *Dao) CreateManyLogging(ctx context.Context, dbName, cName string, docs []interface{}) error {
	inst, err := d.mongo.GetInstance(dbName)
	if err != nil {
		return errc.WithStack(err)
	}
	_, err = inst.Collection(cName).InsertMany(ctx, docs)
	return errc.WithStack(err)
}

func (d *Dao) ListCollectionNames(ctx context.Context, dbName string, prefix ...string) ([]string, error) {
	inst, err := d.mongo.GetInstance(dbName)
	if err != nil {
		return nil, err
	}
	return inst.ListCollectionNames(ctx, prefix...)
}

func (d *Dao) CreateLoggingIndex(dbName, cName string) error {
	inst, err := d.mongo.GetInstance(dbName)
	if err != nil {
		return err
	}
	return inst.UpsertCollectionIndexMany(model.LoggingIndexMany(cName))
}

func (d *Dao) FindLoggingList(ctx context.Context, dbName, cName string, in *model.FindLoggingListReq) (int64, []*model.Logging, error) {
	s := time.Now()
	filter := bson.M{
		"m": strings.TrimSpace(in.ModuleName),
	}
	// 查询条件必须存在时间
	filter["ts"] = bson.M{"$gte": in.BeginTsSec, "$lt": in.EndTsSec}

	if in.Level > -2 {
		filter["l"] = in.Level
	}

	if in.Short != "" {
		if _, ok := filter["l"]; !ok {
			return 0, nil, errc.ErrParamInvalid.MultiMsg("必需传入[等级]，才能使用[短消息]筛选条件")
		}
		filter["s"] = primitive.Regex{
			Pattern: in.Short,
			Options: "i",
		}
	}

	if in.IP != "" {
		if _, ok := filter["s"]; !ok {
			return 0, nil, errc.ErrParamInvalid.MultiMsg("必需传入[短消息]，才能使用[IP]筛选条件")
		}
		filter["ip"] = in.IP
	}
	// 必需要有前置条件，才能查询后面的，以便命中索引
	if in.ConditionOne != "" {
		filter["c1"] = in.ConditionOne
		if in.ConditionTwo != "" {
			filter["c2"] = in.ConditionTwo
			if in.ConditionThree != "" {
				filter["c3"] = in.ConditionThree
			}
		}
	}

	inst, err := d.mongo.GetInstance(dbName)
	if err != nil {
		return 0, nil, errc.ErrParamInvalid.MultiErr(err)
	}
	// 异步统计Count
	calcCount := func(ctx context.Context) (int64, error) {
		opt := options.Count().SetLimit(50000).SetMaxTime(3 * time.Second)
		c, err := inst.Collection(cName).CountDocuments(ctx, filter, opt)
		return c, err
	}
	countResp := make(chan int64, 1)
	go func() {
		c, err := calcCount(ctx)
		if err != nil {
			logs.Qezap.Error("FindLoggingCount", zap.Error(err))
		}
		countResp <- c
	}()

	docs := make([]*model.Logging, 0, in.Limit)
	opt := in.SetPage(options.Find()).SetSort(bson.M{"ts": -1})
	err = inst.Find(ctx, cName, filter, &docs, opt)
	if err != nil {
		return 0, docs, errc.ErrInternalErr.MultiErr(err)
	}

	select {
	case c := <-countResp:
		logs.Qezap.Info("日志查询", zap.String("耗时", time.Since(s).String()),
			zap.String("数据库", dbName),
			zap.Any("集合", cName),
			zap.Any("条件", filter))
		return c, docs, nil
	}
}

func (d *Dao) FindLoggingByTraceID(ctx context.Context, m *model.Module, in *model.FindLoggingByTraceIDReq) ([]*model.Logging, error) {
	tid, err := apiTypes.TraceIDFromHex(in.TraceID)
	if err != nil {
		return nil, errc.ErrParamInvalid.MultiErr(err)
	}
	// 如果查询条件存在TraceID, 则时间范围从 traceID 里面去解析
	// 在TraceTime前后2小时
	tidTime := tid.Time()
	b := tidTime.Add(-2 * time.Hour)
	e := tidTime.Add(2 * time.Hour)
	dbName := m.Database
	cNames := make([]string, 0, 2)
	sc := mongo.NewShardCollection(m.Prefix, m.DaySpan)
	if in.ForceDatabase != "" {
		if !d.cfg.MongoGroup.IsReceiverDatabase(in.ForceDatabase) {
			return nil, errc.ErrParamInvalid.MultiMsg("force database not receiver database")
		}
	}
	if in.ForceCollectionName != "" {
		if !strings.HasPrefix(in.ForceCollectionName, m.Prefix) {
			return nil, errc.ErrParamInvalid.MultiMsg(fmt.Sprintf("force collection name not '%s' prefix", m.Prefix))
		}
		cNames = append(cNames, in.ForceCollectionName)
	} else {
		cNames = append(cNames, sc.CollNameByStartEnd(m.Bucket, b.Unix(), e.Unix())...)
	}

	inst, err := d.mongo.GetInstance(dbName)
	if err != nil {
		return nil, errc.ErrParamInvalid.MultiErr(err)
	}

	allDocs := make([]*model.Logging, 0)
	for _, cName := range cNames {
		filter := bson.M{
			"m":  in.ModuleName,
			"ti": in.TraceID,
		}
		opt := options.Find().SetSort(bson.M{"ts": 1})
		// 正序，调用流
		docs := make([]*model.Logging, 0)
		if err := inst.Find(ctx, cName, filter, &docs, opt); err != nil {
			return nil, errc.ErrInternalErr.MultiErr(err)
		}
		allDocs = append(allDocs, docs...)
	}
	return allDocs, nil
}

func (d *Dao) DropLoggingCollection(ctx context.Context, m *model.Module, cName string) error {
	inst, err := d.mongo.GetInstance(m.Database)
	if err != nil {
		return errc.ErrParamInvalid.MultiErr(err)
	}
	if err := inst.Collection(cName).Drop(ctx); err != nil {
		return errc.ErrInternalErr.MultiErr(err)
	}

	filter := bson.M{
		"module_name": m.Name,
		"name":        cName,
	}
	_, _ = d.adminInst.Collection(model.CNCollStats).DeleteOne(ctx, filter)

	return nil
}
