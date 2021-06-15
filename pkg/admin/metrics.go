package admin

import (
	"context"
	"encoding/binary"
	"fmt"
	"math/rand"
	"sort"
	"strings"
	"time"

	"github.com/bbdshow/qelog/infra/httputil"
	"github.com/bbdshow/qelog/infra/mongo"
	"go.mongodb.org/mongo-driver/bson/primitive"

	"github.com/bbdshow/qelog/pkg/common/model"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"

	"github.com/bbdshow/qelog/pkg/common/entity"
)

func (srv *Service) MetricsDBStats(ctx context.Context, out *entity.ListResp) error {
	// 先查看最后一条， 如果超时就去库里查询
	mainCfg := model.ShardingDB.MainConfig()
	shardingCfg := model.ShardingDB.ShardSlotsConfig()
	mainHost := strings.Join(mongo.URIToHosts(mainCfg.URI), ",")

	dbStats := make([]entity.DBStats, 0)
	mainStats, err := srv.readDBStatsAndInsert(ctx, mainCfg.URI, mainHost, mainCfg.DataBase)
	if err != nil {
		return httputil.ErrSystemException.MergeError(err)
	}
	dbStats = append(dbStats, entity.DBStats{
		ShardingIndex: []int{0},
		Host:          mainStats.Host,
		DBName:        mainStats.DB,
		Collections:   mainStats.Collections,
		DataSize:      mainStats.DataSize,
		StorageSize:   mainStats.StorageSize,
		IndexSize:     mainStats.IndexSize,
		Objects:       mainStats.Objects,
		Indexs:        mainStats.Indexes,
		CreatedTsSec:  mainStats.CreatedAt.Unix(),
	})
	// 去获取 sharding 的DB状态
	for _, v := range shardingCfg {
		host := strings.Join(mongo.URIToHosts(v.URI), ",")
		sStats, err := srv.readDBStatsAndInsert(ctx, v.URI, host, v.DataBase)
		if err != nil {
			return httputil.ErrSystemException.MergeError(err)
		}

		dbStats = append(dbStats, entity.DBStats{
			ShardingIndex: v.Index,
			Host:          sStats.Host,
			DBName:        sStats.DB,
			Collections:   sStats.Collections,
			DataSize:      sStats.DataSize,
			StorageSize:   sStats.StorageSize,
			IndexSize:     sStats.IndexSize,
			Objects:       sStats.Objects,
			Indexs:        sStats.Indexes,
			CreatedTsSec:  sStats.CreatedAt.Unix(),
		})
	}

	out.Count = int64(len(dbStats))
	out.List = dbStats

	return nil
}

func (srv *Service) MetricsCollStats(ctx context.Context, in *entity.MetricsCollStatsReq, out *entity.ListResp) error {
	uri := ""
	mainCfg := model.ShardingDB.MainConfig()
	shardingCfg := model.ShardingDB.ShardSlotsConfig()

	mainHost := strings.Join(mongo.URIToHosts(mainCfg.URI), ",")
	if mainHost == in.Host {
		uri = mainCfg.URI
	}
	if uri == "" {
		for _, s := range shardingCfg {
			host := strings.Join(mongo.URIToHosts(s.URI), ",")
			if host == in.Host {
				uri = s.URI
				break
			}
		}
	}

	if uri == "" {
		return httputil.ErrNotFound
	}
	collStats, err := srv.readCollStatsAndInsert(ctx, uri, in.Host, in.DBName)
	if err != nil {
		return err
	}

	list := make([]*entity.CollStats, 0, len(collStats))
	for _, v := range collStats {
		d := &entity.CollStats{
			Host:           v.Host,
			DBName:         v.DB,
			Name:           v.Name,
			Size:           v.Size,
			Count:          v.Count,
			AvgObjSize:     v.AvgObjSize,
			StorageSize:    v.StorageSize,
			Capped:         v.Capped,
			TotalIndexSize: v.TotalIndexSize,
			IndexSizes:     v.IndexSizes,
			CreatedTsSec:   v.CreatedAt.Unix(),
		}
		list = append(list, d)
	}

	out.Count = int64(len(list))
	out.List = list

	return nil
}

func (srv *Service) MetricsModuleList(ctx context.Context, in *entity.MetricsModuleListReq, out *entity.ListResp) error {
	y, m, d := time.Unix(in.DateTsSec, 0).Date()
	date := time.Date(y, m, d, 0, 0, 0, 0, time.Local)

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
	c, docs, err := srv.moduleMetricsStore.FindCountModuleMetrics(ctx, filter, opt)
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

func (srv *Service) MetricsModuleTrend(ctx context.Context, in *entity.MetricsModuleTrendReq, out *entity.MetricsModuleTrendResp) error {
	now := time.Now()
	lastDay := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, time.Local).AddDate(0, 0, -in.LastDay)
	filter := bson.M{
		"module_name":  in.ModuleName,
		"created_date": bson.M{"$gte": lastDay},
	}

	_, docs, err := srv.moduleMetricsStore.FindCountModuleMetrics(ctx, filter, nil)
	if err != nil {
		return httputil.ErrSystemException.MergeError(err)
	}

	number := int64(0)
	size := int64(0)
	ascTsNumbers := make([]model.TsNumbers, 0, in.LastDay*24)
	allLevels := make(map[model.Level]bool)
	allIps := make(map[string]bool)
	for _, v := range docs {
		number += v.Number
		size += v.Size
		for ts, numbers := range v.Sections {
			ascTsNumbers = append(ascTsNumbers, model.TsNumbers{
				Ts:      ts,
				Numbers: numbers,
			})
			// 需要找出这个时间段的所有 等级 和 ip， 在后面排序时，如果某个时间段不存在值， 则应该默认 0
			for lvl := range numbers.Levels {
				allLevels[lvl] = true
			}
			for ip := range numbers.IPs {
				allIps[ip] = true
			}
		}
	}
	sort.Sort(model.AscTsNumbers(ascTsNumbers))

	xData := make([]string, 0, in.LastDay*24)
	levelMapData := map[model.Level][]int32{}
	ipMapData := map[string][]int32{}
	for _, v := range ascTsNumbers {
		t := time.Unix(v.Ts, 0)
		xData = append(xData, fmt.Sprintf("%d-%d %d:00", t.Month(), t.Day(), t.Hour()))
		//xData = append(xData, t.Format("2006-01-02 15:04:05"))
		for lvl := range allLevels {
			data, ok := levelMapData[lvl]
			if !ok {
				data = make([]int32, 0, in.LastDay*24)
				levelMapData[lvl] = data
			}
			// 这个时间段没有这个错误等级，默认设置为 0
			num := v.Levels[lvl]
			levelMapData[lvl] = append(data, num)
		}

		for ip := range allIps {
			data, ok := ipMapData[ip]
			if !ok {
				data = make([]int32, 0, in.LastDay*24)
				ipMapData[ip] = data
			}
			// 这个时间段没有这个错误等级，默认设置为 0
			num := v.IPs[ip]
			ipMapData[ip] = append(data, num)
		}
	}
	out.XData = xData
	out.Number = number
	out.Size = size

	legend := make([]string, 0)

	levelSeries := make([]entity.Serie, 0, len(levelMapData))
	for lvl, data := range levelMapData {
		levelSeries = append(levelSeries, entity.Serie{
			Index: lvl.Int32(),
			Name:  lvl.String(),
			Type:  "bar",
			Color: levelColor(lvl),
			Data:  data,
		})
	}
	entity.SortSeries(levelSeries, "ASC")
	for _, v := range levelSeries {
		legend = append(legend, v.Name)
	}
	ipSeries := make([]entity.Serie, 0, len(ipMapData))
	for ip, data := range ipMapData {
		ip := strings.ReplaceAll(ip, "_", ".")
		legend = append(legend, ip)
		ipSeries = append(ipSeries, entity.Serie{
			Index: int32(binary.BigEndian.Uint32([]byte(ip))),
			Name:  ip,
			Type:  "line",
			Color: ipColor(),
			Data:  data,
		})
	}
	entity.SortSeries(ipSeries, "ASC")
	for _, v := range ipSeries {
		legend = append(legend, v.Name)
	}
	out.LegendData = legend
	out.LevelSeries = levelSeries
	out.IPSeries = ipSeries

	return nil
}

func levelColor(lvl model.Level) string {
	switch lvl.String() {
	case "DEBUG":
		return "rgba(144,202,249,1)"
	case "INFO":
		return "rgba(30,150,243,1)"
	case "WARN":
		return "rgba(251,192,45,1)"
	case "ERROR":
		return "rgba(244,67,54,1)"
	case "DPANIC":
		return "rgba(211,47,47,1)"
	case "PANIC":
		return "rgba(198,40,40,1)"
	case "FATAL":
		return "rgba(0,0,0,1)"
	}
	return "rgba(255,255,255,1)"
}
func ipColor() string {
	return fmt.Sprintf("rgba(%d,%d,%d,1)", rand.Int31n(100)+150, rand.Int31n(80)+100, rand.Int31n(135)+100)
}

func (srv *Service) readDBStatsAndInsert(ctx context.Context, uri, host, database string) (*model.DBStats, error) {
	filter := bson.M{
		"host":       host,
		"db":         database,
		"created_at": bson.M{"$gte": time.Now().Add(-10 * time.Minute)},
	}
	opt := options.FindOne()
	opt.SetSort(bson.M{"_id": -1})
	latestDBStats := &model.DBStats{}
	ok, err := srv.moduleMetricsStore.FindOneDBStats(ctx, filter, latestDBStats, opt)
	if err != nil {
		return nil, err
	}
	if ok {
		// 有效
		return latestDBStats, nil
	}

	db, err := mongo.NewDatabase(ctx, uri, database)
	if err != nil {
		return nil, err
	}
	defer db.Client().Disconnect(ctx)
	util := mongo.NewUtil(db)

	stats, err := util.DBStats(ctx)
	if err != nil {
		return nil, err
	}
	doc := &model.DBStats{
		Host:        host,
		DB:          database,
		Collections: stats.Collections,
		Objects:     stats.Objects,
		DataSize:    stats.DataSize,
		StorageSize: stats.StorageSize,
		Indexes:     stats.Indexes,
		IndexSize:   stats.Indexes,
		CreatedAt:   time.Now(),
	}

	_ = srv.moduleMetricsStore.InsertOneDBStats(ctx, doc)

	return doc, nil
}

func (srv *Service) readCollStatsAndInsert(ctx context.Context, uri, host, database string) ([]*model.CollStats, error) {
	filter := bson.M{
		"host":       host,
		"db":         database,
		"created_at": bson.M{"$gte": time.Now().Add(-30 * time.Minute)},
	}
	opt := options.Find().SetSort(bson.M{"name": -1})
	latestCollStats, err := srv.moduleMetricsStore.FindCollStats(ctx, filter, opt)
	if err != nil {
		return nil, err
	}

	db, err := mongo.NewDatabase(ctx, uri, database)
	if err != nil {
		return nil, err
	}
	defer db.Client().Disconnect(ctx)

	names, err := db.ListCollectionNames(ctx)
	if err != nil {
		return nil, err
	}
	if len(latestCollStats) >= len(names) {
		return latestCollStats, nil
	}

	// 去读取最新的结果
	util := mongo.NewUtil(db)

	collStats, err := util.CollStats(ctx, names)
	if err != nil {
		return nil, err
	}

	latestCollStats = make([]*model.CollStats, 0, len(collStats))
	insertDocs := make([]interface{}, 0, len(collStats))

	for _, v := range collStats {
		if v.Ok == 0 {
			continue
		}
		doc := &model.CollStats{
			Host:           host,
			DB:             database,
			Name:           v.Ns,
			Size:           v.Size,
			Count:          v.Count,
			AvgObjSize:     v.AvgObjSize,
			StorageSize:    v.StorageSize,
			Capped:         v.Capped,
			TotalIndexSize: v.TotalIndexSize,
			IndexSizes:     v.IndexSizes,
			CreatedAt:      time.Now(),
		}
		insertDocs = append(insertDocs, doc)
		latestCollStats = append(latestCollStats, doc)
	}

	_ = srv.moduleMetricsStore.InsertManyCollStats(ctx, insertDocs)

	return latestCollStats, nil
}
