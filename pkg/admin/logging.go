package admin

import (
	"context"
	"fmt"
	"github.com/bbdshow/bkit/db/mongo"
	"github.com/bbdshow/bkit/errc"
	"github.com/bbdshow/bkit/logs"
	"github.com/bbdshow/qelog/pkg/model"
	"go.mongodb.org/mongo-driver/bson"
	"go.uber.org/zap"
	"math/rand"
	"strings"
	"time"
)

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
	// 过滤掉重复写入的数据
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

func (svc *Service) FindLoggingList(ctx context.Context, in *model.FindLoggingListReq, out *model.ListResp) error {
	exists, m, err := svc.d.GetModule(ctx, bson.M{"name": in.ModuleName})
	if err != nil {
		return errc.ErrInternalErr.MultiErr(err)
	}
	if !exists {
		return errc.ErrNotFound.MultiMsg("module")
	}

	// 如果没有传入时间，则默认设置一个间隔时间
	b, e := in.InitTimeSection(time.Hour)
	in.BeginTsSec = b.Unix()
	in.EndTsSec = e.Unix()

	// 计算查询时间应该在哪个分片
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
		// 计算集合名
		names := sc.CollNameByStartEnd(m.Bucket, in.BeginTsSec, in.EndTsSec)
		if len(names) >= 2 {
			format := "2006-01-02"
			sepTime, err := sc.SepTime(names[0])
			if err != nil {
				return errc.ErrParamInvalid.MultiErr(err)
			}
			sep := sepTime.Format(format)
			return errc.ErrParamInvalid.MultiMsg(fmt.Sprintf("查询时间已跨分片集合,未不影响结果,建议查询时间: %s -- %s 或者 %s -- %s",
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

	// 去除极低可能重复写入的日志信息
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

// bgDelExpiredCollection 删除已经过期了的集合
func (svc *Service) bgDelExpiredCollection() {
	for {
		time.Sleep(time.Duration(rand.Intn(30)+30) * time.Minute)
		// 查找所有的 module
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
