package dao

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/bbdshow/bkit/db/mongo"
	"github.com/bbdshow/bkit/errc"
	"github.com/bbdshow/bkit/logs"
	apiTypes "github.com/bbdshow/qelog/api/types"
	"github.com/bbdshow/qelog/pkg/model"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.uber.org/zap"
)

// CreateManyLogging multi insert logging to db
func (d *Dao) CreateManyLogging(ctx context.Context, dbName, cName string, docs []interface{}) error {
	inst, err := d.mongo.GetInstance(dbName)
	if err != nil {
		return errc.WithStack(err)
	}
	_, err = inst.Collection(cName).InsertMany(ctx, docs)
	return errc.WithStack(err)
}

// ListCollectionNames query this db collection by prefix
func (d *Dao) ListCollectionNames(ctx context.Context, dbName string, prefix ...string) ([]string, error) {
	inst, err := d.mongo.GetInstance(dbName)
	if err != nil {
		return nil, err
	}
	return inst.ListCollectionNames(ctx, prefix...)
}

// CreateLoggingIndex db runtime create index, when new collection created
func (d *Dao) CreateLoggingIndex(dbName, cName string) error {
	inst, err := d.mongo.GetInstance(dbName)
	if err != nil {
		return err
	}
	return inst.UpsertCollectionIndexMany(model.LoggingIndexMany(cName))
}

// FindLoggingList query logging
func (d *Dao) FindLoggingList(ctx context.Context, dbName, cName string, in *model.FindLoggingListReq) (int64, []*model.Logging, error) {
	s := time.Now()
	filter := bson.M{
		"m": strings.TrimSpace(in.ModuleName),
	}
	// condition required time arg
	filter["ts"] = bson.M{"$gte": in.BeginTsSec, "$lt": in.EndTsSec}

	if in.Level > -2 {
		filter["l"] = in.Level
	}

	if in.Short != "" {
		if _, ok := filter["l"]; !ok {
			return 0, nil, errc.ErrParamInvalid.MultiMsg("required level condition before it can used short message condition")
		}
		filter["s"] = primitive.Regex{
			Pattern: in.Short,
			Options: "i",
		}
	}

	if in.IP != "" {
		if _, ok := filter["s"]; !ok {
			return 0, nil, errc.ErrParamInvalid.MultiMsg("required short message condition before it can used IP condition")
		}
		filter["ip"] = in.IP
	}

	// condition of dependence,in order for index performance
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
	// async count, max number 50000 or max time 15s
	// too many bar counts, meaningless, and poor performance
	calcCount := func(ctx context.Context) (int64, error) {
		opt := options.Count().SetLimit(50000).SetMaxTime(d.CtxAfterSecDeadline(ctx, 15))
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
		if c <= 0 {
			c = int64(len(docs))
		}
		logs.Qezap.Info("LoggingQuery", zap.String("latency", time.Since(s).String()),
			zap.String("database", dbName),
			zap.Any("collection", cName),
			zap.Any("condition", filter))
		return c, docs, nil
	}
}

// FindLoggingByTraceID query logging by traceId
func (d *Dao) FindLoggingByTraceID(ctx context.Context, m *model.Module, in *model.FindLoggingByTraceIDReq) ([]*model.Logging, error) {
	tid, err := apiTypes.TraceIDFromHex(in.TraceID)
	if err != nil {
		return nil, errc.ErrParamInvalid.MultiErr(err)
	}
	// traceId have time info, set it time condition +-2 hour
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
		// asc logging by time
		opt := options.Find().SetSort(bson.M{"ts": 1})
		docs := make([]*model.Logging, 0)
		if err := inst.Find(ctx, cName, filter, &docs, opt); err != nil {
			return nil, errc.ErrInternalErr.MultiErr(err)
		}
		allDocs = append(allDocs, docs...)
	}
	return allDocs, nil
}

// DropLoggingCollection delete collection
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
