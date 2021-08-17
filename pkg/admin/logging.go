package admin

import (
	"context"
	"errors"
	"fmt"
	"github.com/bbdshow/qelog/pkg/config"
	"github.com/bbdshow/qelog/pkg/storage"
	"math/rand"
	"strings"
	"time"

	apiTypes "github.com/bbdshow/qelog/api/types"

	"github.com/bbdshow/qelog/infra/httputil"
	"github.com/bbdshow/qelog/infra/logs"
	"github.com/bbdshow/qelog/infra/mongo"
	"github.com/bbdshow/qelog/pkg/common/entity"
	"github.com/bbdshow/qelog/pkg/common/model"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.uber.org/zap"
)

func (srv *Service) FindLoggingByTraceID(ctx context.Context, in *entity.FindLoggingByTraceIDReq, out *entity.ListResp) error {
	tid, err := apiTypes.TraceIDFromHex(in.TraceID)
	if err != nil {
		return httputil.ErrArgsInvalid.MergeError(err)
	}
	// 如果查询条件存在TraceID, 则时间范围从 traceID 里面去解析
	// 在TraceTime前后2小时
	tidTime := tid.Time()
	b := tidTime.Add(-2 * time.Hour)
	e := tidTime.Add(2 * time.Hour)
	collectionNames := make([]string, 0, 2)
	if in.ForceCollectionName != "" {
		if strings.HasPrefix(in.ForceCollectionName, "logging") {
			collectionNames = append(collectionNames, in.ForceCollectionName)
		}
	} else {
		collectionNames = append(collectionNames, srv.sc.ScopeCollectionNames(in.ShardingIndex, b.Unix(), e.Unix())...)
	}
	count := int64(0)
	list := make([]*entity.FindLoggingList, 0)
	for _, coll := range collectionNames {
		filter := bson.M{
			"m":  in.ModuleName,
			"ti": in.TraceID,
		}
		findOpt := options.Find().SetSort(bson.M{"ts": 1})
		// 正序，调用流

		shardSlot, err := model.ShardingDB.ShardSlotDB(in.ShardingIndex)
		if err != nil {
			return httputil.ErrArgsInvalid.MergeError(err)
		}
		c, docs, err := storage.NewLogging(shardSlot).FindCountLoggingList(ctx, coll, filter, 50000, findOpt)
		if err != nil {
			return httputil.ErrSystemException.MergeError(err)
		}
		count += c
		// 过滤掉重复写入的数据
		hitMap := map[string]struct{}{}

		for _, v := range docs {
			if _, ok := hitMap[v.MessageID]; ok {
				continue
			} else {
				hitMap[v.MessageID] = struct{}{}
			}

			d := &entity.FindLoggingList{
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
	}

	out.Count = count
	out.List = list

	return nil
}

func (srv *Service) FindLoggingList(ctx context.Context, in *entity.FindLoggingListReq, out *entity.ListResp) error {

	s := time.Now()
	// 如果没有传入时间，则默认设置一个间隔时间
	b, e := in.InitTimeSection(time.Hour)
	// 计算查询时间应该在哪个分片
	collectionName := ""
	if in.ForceCollectionName != "" {
		if strings.HasPrefix(in.ForceCollectionName, "logging") {
			collectionName = in.ForceCollectionName
		}
	} else {
		// 计算集合名
		names := srv.sc.ScopeCollectionNames(in.ShardingIndex, b.Unix(), e.Unix())
		if len(names) >= 2 {
			format := "2006-01-02 15:04:05"
			suggestTime, _ := srv.sc.SuggestSpanTime(names[0])
			suggest := suggestTime.Format(format)
			return httputil.NewError(httputil.ErrCodeOpException,
				fmt.Sprintf("查询时间已跨分片集合,未不影响结果,建议查询时间: %s -- %s 或者 %s -- %s",
					b.Format(format), suggest, suggest, e.Format(format)))
		}
		if len(names) > 0 {
			collectionName = names[0]
		}
	}

	filter := bson.M{
		"m": strings.TrimSpace(in.ModuleName),
	}

	// 查询条件必须存在时间
	filter["ts"] = bson.M{"$gte": b.Unix(), "$lt": e.Unix()}

	if in.Level > -2 {
		filter["l"] = in.Level
	}

	if in.Short != "" {
		if _, ok := filter["l"]; !ok {
			return httputil.ErrArgsInvalid.MergeString("必需传入[等级]，才能使用[短消息]筛选条件")
		}
		filter["s"] = primitive.Regex{
			Pattern: in.Short,
			Options: "i",
		}
	}

	if in.IP != "" {
		if _, ok := filter["s"]; !ok {
			return httputil.ErrArgsInvalid.MergeString("必需传入[短消息]，才能使用[IP]筛选条件")
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

	findOpt := in.SetPage(options.Find()).SetSort(bson.M{"ts": -1})

	shardSlot, err := model.ShardingDB.ShardSlotDB(in.ShardingIndex)
	if err != nil {
		return httputil.ErrArgsInvalid.MergeError(err)
	}
	c, docs, err := storage.NewLogging(shardSlot).FindCountLoggingList(ctx, collectionName, filter, 50000, findOpt)
	if err != nil {
		return httputil.ErrSystemException.MergeError(err)
	}
	out.Count = c

	// 去除极低可能重复写入的日志信息
	hitMap := map[string]struct{}{}
	list := make([]*entity.FindLoggingList, 0, len(docs))
	for _, v := range docs {
		if _, ok := hitMap[v.MessageID]; ok {
			continue
		} else {
			hitMap[v.MessageID] = struct{}{}
		}

		d := &entity.FindLoggingList{
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

	logs.Qezap.Info("日志查询", zap.String("耗时", time.Now().Sub(s).String()),
		zap.String("分片", shardSlot.Name()),
		zap.Any("集合", collectionName),
		zap.Any("条件", filter))

	return nil
}

func (srv *Service) DropLoggingCollection(ctx context.Context, in *entity.DropLoggingCollectionReq) error {
	//  先检查 collectionName
	dbColl := strings.Split(in.Name, ".")
	if len(dbColl) < 2 {
		return httputil.ErrArgsInvalid.MergeError(errors.New("name"))
	}
	database := dbColl[0]
	collectionName := dbColl[1]
	if !strings.HasPrefix(collectionName, "logging") {
		return httputil.ErrArgsInvalid.MergeError(errors.New("drop only logging prefix collection name"))
	}

	// 根据host找到db
	uri := ""
	mainCfg := model.ShardingDB.MainConfig()
	shardingCfg := model.ShardingDB.ShardSlotsConfig()

	mainHost := strings.Join(mongo.URIToHosts(mainCfg.URI), ",")
	if mainHost == in.Host && database == mainCfg.DataBase {
		uri = mainCfg.URI
	}
	if uri == "" {
		for _, s := range shardingCfg {
			host := strings.Join(mongo.URIToHosts(s.URI), ",")
			if host == in.Host && database == s.DataBase {
				uri = s.URI
				break
			}
		}
	}

	if uri == "" {
		return httputil.ErrNotFound
	}

	db, err := mongo.NewDatabase(ctx, uri, database)
	if err != nil {
		return httputil.ErrSystemException.MergeError(err)
	}
	defer db.Client().Disconnect(ctx)

	if err := db.Collection(collectionName).Drop(ctx); err != nil {
		return httputil.ErrSystemException.MergeError(err)
	}

	filter := bson.M{
		"host": in.Host,
		"db":   database,
		"name": in.Name,
	}
	// 同时删除主库集合统计数据
	_, err = model.MainDB.Collection(model.CollectionNameCollStats).DeleteMany(ctx, filter)
	if err != nil {
		return httputil.ErrSystemException.MergeError(err)
	}
	return nil
}

// backgroundDelExpiredCollection 删除已经过期了的集合
// 月为单位
func (srv *Service) backgroundDelExpiredCollection(maxAgeMonth int) {
	if maxAgeMonth <= 0 {
		// 永久保存
		return
	}
	for {
		time.Sleep(time.Duration(rand.Intn(5)+5) * time.Second)
		for i := 1; i <= config.Global.ShardingIndexSize; i++ {
			db, err := model.ShardingDB.ShardSlotDB(i)
			if err != nil {
				continue
			}
			// 找到所有 logging 开头的集合
			names, err := db.ListCollectionNames(context.Background(), "logging")
			if err != nil {
				logs.Qezap.Error("ListCollectionNames", zap.Error(err))
				continue
			}
			expiredNames := make([]string, 0)
			// 判断是否过期
			for _, v := range names {
				date, err := srv.sc.CollectionNameToTime(v)
				if err != nil {
					logs.Qezap.Error("NameDecodeDate", zap.Error(err))
					continue
				}
				y, m, _ := time.Now().Date()
				expiredTime := time.Date(y, m, 0, 0, 0, 0, 0, time.Local).
					AddDate(0, -maxAgeMonth, 0)

				if expiredTime.Equal(date) || expiredTime.After(date) {
					expiredNames = append(expiredNames, v)
				}
			}

			for _, v := range expiredNames {
				if err := db.Collection(v).Drop(context.Background()); err != nil {
					logs.Qezap.Error("DropCollection", zap.Error(err))
					continue
				}
				time.Sleep(3 * time.Second)
			}
		}

		time.Sleep(time.Duration(rand.Intn(30)+30) * time.Minute)
	}

}
