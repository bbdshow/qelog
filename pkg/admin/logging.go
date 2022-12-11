package admin

import (
	"context"
	"fmt"
	"math/rand"
	"strings"
	"time"

	"github.com/bbdshow/bkit/db/mongo"
	"github.com/bbdshow/bkit/errc"
	"github.com/bbdshow/bkit/logs"
	"github.com/bbdshow/qelog/pkg/model"
	"go.mongodb.org/mongo-driver/bson"
	"go.uber.org/zap"
)

// FindLoggingByTraceID query logging by traceId, returns all bind this traceId logging info
func (svc *Service) FindLoggingByTraceID(ctx context.Context, in *model.FindLoggingByTraceIDReq, out *model.ListResp) error {
	exists, m, err := svc.d.GetModule(ctx, bson.M{"name": in.ModuleName})
	if err != nil {
		return errc.ErrInternalErr.MultiErr(err)
	}
	if !exists {
		return errc.ErrNotFound.MultiMsg("module")
	}

	docs, err := svc.d.FindLoggingByTraceID(ctx, m, in)
	if err != nil {
		return errc.WithStack(err)
	}
	list := make([]*model.FindLoggingList, 0, len(docs))
	// there is a low probability of data being written repeatedly
	// filtering duplicate written data
	hitMap := map[string]struct{}{}
	for _, v := range docs {
		if _, ok := hitMap[v.MessageID]; ok {
			continue
		} else {
			hitMap[v.MessageID] = struct{}{}
		}
		d := &model.FindLoggingList{
			ID:             v.ID.Hex(),
			TsMill:         v.TimeMill,
			Level:          int32(v.Level),
			Short:          v.Short,
			Full:           v.Full,
			ConditionOne:   v.Condition1,
			ConditionTwo:   v.Condition2,
			ConditionThree: v.Condition3,
			IP:             v.IP,
			TraceID:        v.TraceID,
		}
		list = append(list, d)
	}

	out.Count = int64(len(list))
	out.List = list

	return nil
}

// FindLoggingList query logging by common condition
func (svc *Service) FindLoggingList(ctx context.Context, in *model.FindLoggingListReq, out *model.ListResp) error {
	exists, m, err := svc.d.GetModule(ctx, bson.M{"name": in.ModuleName})
	if err != nil {
		return errc.ErrInternalErr.MultiErr(err)
	}
	if !exists {
		return errc.ErrNotFound.MultiMsg("module")
	}

	// if condition not time filter, default setting 1 hour
	b, e := in.InitTimeSection(time.Hour)
	in.BeginTsSec = b.Unix()
	in.EndTsSec = e.Unix()

	// by querying the time, the data is calculated in which shard
	sc := mongo.NewShardCollection(m.Prefix, m.DaySpan)
	dbName := m.Database
	cName := ""
	if in.ForceDatabase != "" {
		if !svc.cfg.MongoGroup.IsReceiverDatabase(in.ForceDatabase) {
			return errc.ErrParamInvalid.MultiMsg("force database not receiver database")
		}
	}
	if in.ForceCollectionName != "" {
		if !strings.HasPrefix(in.ForceCollectionName, m.Prefix) {
			return errc.ErrParamInvalid.MultiMsg(fmt.Sprintf("force collection name not '%s' prefix", m.Prefix))
		}
		cName = in.ForceCollectionName
	} else {
		// calc collection name
		names := sc.CollNameByStartEnd(m.Bucket, in.BeginTsSec, in.EndTsSec)
		if len(names) >= 2 {
			format := "2006-01-02"
			sepTime, err := sc.SepTime(names[0])
			if err != nil {
				return errc.ErrParamInvalid.MultiErr(err)
			}
			sep := sepTime.Format(format)
			return errc.ErrParamInvalid.MultiMsg(fmt.Sprintf("Time has crossed shards,suggest time: %s -- %s || %s -- %s",
				time.Unix(in.BeginTsSec, 0).Format(format), sep, sep, time.Unix(in.EndTsSec, 0).Format(format)))
		}
		if len(names) > 0 {
			cName = names[0]
		}
	}

	c, docs, err := svc.d.FindLoggingList(ctx, dbName, cName, in)
	if err != nil {
		return errc.WithStack(err)
	}
	out.Count = c

	// there is a low probability of data being written repeatedly
	// filtering duplicate written data
	hitMap := map[string]struct{}{}
	list := make([]*model.FindLoggingList, 0, len(docs))
	for _, v := range docs {
		if _, ok := hitMap[v.MessageID]; ok {
			continue
		} else {
			hitMap[v.MessageID] = struct{}{}
		}

		d := &model.FindLoggingList{
			ID:             v.ID.Hex(),
			TsMill:         v.TimeMill,
			Level:          int32(v.Level),
			Short:          v.Short,
			Full:           v.Full,
			ConditionOne:   v.Condition1,
			ConditionTwo:   v.Condition2,
			ConditionThree: v.Condition3,
			IP:             v.IP,
			TraceID:        v.TraceID,
		}
		list = append(list, d)
	}
	out.List = list

	return nil
}

// DropLoggingCollection manual delete collection, release storage disk space
func (svc *Service) DropLoggingCollection(ctx context.Context, in *model.DropLoggingCollectionReq) error {
	exists, m, err := svc.d.GetModule(ctx, bson.M{"name": in.ModuleName})
	if err != nil {
		return errc.ErrInternalErr.MultiErr(err)
	}
	if !exists {
		return errc.ErrNotFound.MultiMsg("module")
	}
	if !strings.HasPrefix(in.Collection, m.LoggingPrefix()) {
		return errc.ErrParamInvalid.MultiMsg(fmt.Sprintf("collection name not '%s' prefix", m.LoggingPrefix()))
	}

	if err := svc.d.DropLoggingCollection(ctx, m, in.Collection); err != nil {
		return err
	}
	return nil
}

// auto delete expired collection, release storage disk space
func (svc *Service) bgDelExpiredCollection() {
	for {
		time.Sleep(time.Duration(rand.Intn(30)+30) * time.Minute)
		// find all module
		modules, err := svc.d.FindAllModule(context.Background())
		if err != nil {
			logs.Qezap.Error("bgDelExpiredCollection", zap.String("FindAllModule", err.Error()))
			continue
		}
		for _, m := range modules {
			if m.MaxMonth <= 0 {
				continue
			}
			sc := mongo.NewShardCollection(m.Prefix, m.DaySpan)
			cNames, err := svc.d.ListCollectionNames(context.Background(), m.Database, m.LoggingPrefix())
			if err != nil {
				logs.Qezap.Error("bgDelExpiredCollection", zap.String("ListCollectionNames", err.Error()))
				continue
			}
			for _, cName := range cNames {
				t, err := sc.CollNameDate(cName)
				if err != nil {
					logs.Qezap.Error("bgDelExpiredCollection", zap.String("CollNameDate", err.Error()))
					continue
				}
				if t.AddDate(0, m.MaxMonth, 0).Before(time.Now()) {
					if err := svc.d.DropLoggingCollection(context.Background(), m, cName); err != nil {
						logs.Qezap.Error("bgDelExpiredCollection", zap.String("DropLoggingCollection", err.Error()))
						continue
					}
				}
			}
		}
	}
}
