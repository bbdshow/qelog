package admin

import (
	"context"
	"encoding/binary"
	"fmt"
	"github.com/bbdshow/bkit/db/mongo"
	"github.com/bbdshow/bkit/errc"
	"github.com/bbdshow/bkit/logs"
	"github.com/bbdshow/bkit/util/itime"
	"github.com/bbdshow/qelog/common/types"
	"github.com/bbdshow/qelog/pkg/model"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.uber.org/zap"
	"math/rand"
	"sort"
	"strings"
	"time"
)

// bgMetricsDBStats 统计数据库状态
func (svc *Service) bgMetricsDBStats() {
	for {
		time.Sleep(time.Duration(rand.Intn(30)+60) * time.Minute)
		//time.Sleep(time.Duration(30) * time.Second)
		databases := svc.cfg.MongoGroup.Databases()
		for _, dbName := range databases {
			// find conn
			for _, conn := range svc.cfg.Mongo.Conns {
				if conn.Database == dbName {
					if err := svc.metricsDBStats(conn); err != nil {
						logs.Qezap.Error("bgMetricsDBStats", zap.String("metricsDBStats", err.Error()))
						continue
					}
				}
			}
		}
	}
}

func (svc *Service) metricsDBStats(conn mongo.Conn) error {
	host := strings.Join(mongo.URIToHosts(conn.URI), ",")
	beforeDay := itime.BeforeDayDate(1)
	filter := bson.M{
		"host":       host,
		"db":         conn.Database,
		"updated_at": bson.M{"$gt": beforeDay},
	}
	ctx := context.Background()
	exists, _, err := svc.d.GetDBStats(ctx, filter)
	if err != nil {
		return err
	}
	if exists {
		return nil
	}

	db, err := mongo.NewDatabase(ctx, conn.URI, conn.Database)
	if err != nil {
		return err
	}
	defer db.Client().Disconnect(ctx)

	stats, err := svc.d.ReadDBStats(ctx, db)
	if err != nil {
		return err
	}
	doc := &model.DBStats{
		Host:        host,
		DB:          conn.Database,
		Collections: stats.Collections,
		Objects:     stats.Objects,
		DataSize:    stats.DataSize,
		StorageSize: stats.StorageSize,
		Indexes:     stats.Indexes,
		IndexSize:   stats.Indexes,
	}
	return svc.d.UpsertDBStats(ctx, doc)
}

// bgMetricsCollectionStats 统计集合状态
func (svc *Service) bgMetricsCollectionStats() {
	for {
		time.Sleep(time.Duration(rand.Intn(30)+30) * time.Minute)
		//time.Sleep(time.Duration(30) * time.Second)
		ctx := context.Background()
		modules, err := svc.d.FindAllModule(ctx)
		if err != nil {
			logs.Qezap.Error("bgMetricsCollectionStats", zap.String("FindAllModule", err.Error()))
			continue
		}
		for _, m := range modules {
			// find conn
			for _, conn := range svc.cfg.Mongo.Conns {
				if conn.Database == m.Database {
					colls, err := svc.d.ListCollectionNames(ctx, m.Database, m.LoggingPrefix())
					if err != nil {
						logs.Qezap.Error("bgMetricsCollectionStats", zap.String("ListCollectionNames", err.Error()))
						continue
					}
					if err := svc.metricsCollStats(conn, m, colls); err != nil {
						logs.Qezap.Error("bgMetricsDBStats", zap.String("metricsCollStats", err.Error()))
						continue
					}
					time.Sleep(3 * time.Second)
				}
			}
		}
	}
}

func (svc *Service) metricsCollStats(conn mongo.Conn, m *model.Module, colls []string) error {
	validColls := make([]string, 0)
	host := strings.Join(mongo.URIToHosts(conn.URI), ",")
	beforeDay := itime.BeforeDayDate(1)
	for _, coll := range colls {
		filter := bson.M{
			"module_name": m.Name,
			"host":        host,
			"db":          conn.Database,
			"name":        coll,
			"updated_at":  bson.M{"$gt": beforeDay},
		}
		ctx := context.Background()
		exists, _, err := svc.d.GetCollStats(ctx, filter)
		if err != nil {
			return err
		}
		if !exists {
			validColls = append(validColls, coll)
		}
	}
	if len(validColls) <= 0 {
		return nil
	}

	ctx := context.Background()
	db, err := mongo.NewDatabase(ctx, conn.URI, conn.Database)
	if err != nil {
		return err
	}
	defer db.Client().Disconnect(ctx)

	for _, coll := range validColls {
		stats, err := svc.d.ReadCollStats(ctx, db, coll)
		if err != nil {
			logs.Qezap.Error("metricsCollStats", zap.String("ReadCollStats", err.Error()))
			continue
		}
		doc := &model.CollStats{
			ModuleName:     m.Name,
			Host:           host,
			DB:             conn.Database,
			Size:           stats.Size,
			Count:          stats.Count,
			AvgObjSize:     stats.AvgObjSize,
			StorageSize:    stats.StorageSize,
			Capped:         stats.Capped,
			TotalIndexSize: stats.TotalIndexSize,
			IndexSizes:     stats.IndexSizes,
		}
		// ns 去除 db
		ns := strings.Replace(stats.Ns, fmt.Sprintf("%s.", conn.Database), "", 1)
		doc.Name = ns
		if err := svc.d.UpsertCollStats(ctx, doc); err != nil {
			logs.Qezap.Error("metricsCollStats", zap.String("UpsertCollStats", err.Error()))
			continue
		}
	}
	return nil
}

func (svc *Service) MetricsDBStats(ctx context.Context, out *model.ListResp) error {
	docs, err := svc.d.FindDBStats(ctx, bson.M{})
	if err != nil {
		return errc.ErrInternalErr.MultiErr(err)
	}
	list := make([]*model.DBStat, 0, len(docs))

	for _, v := range docs {
		list = append(list, &model.DBStat{
			Host:         v.Host,
			DBName:       v.DB,
			Collections:  v.Collections,
			DataSize:     v.DataSize,
			StorageSize:  v.StorageSize,
			IndexSize:    v.IndexSize,
			Objects:      v.Objects,
			Indexs:       v.Indexes,
			UpdatedTsSec: v.UpdatedAt.Unix(),
		})
	}

	out.Count = int64(len(list))
	out.List = list
	return nil
}

// MetricsCollStats 集合统计
func (svc *Service) MetricsCollStats(ctx context.Context, in *model.MetricsCollStatsReq, out *model.ListResp) error {
	exists, m, err := svc.d.GetModule(ctx, bson.M{"name": in.ModuleName})
	if err != nil {
		return errc.ErrInternalErr.MultiErr(err)
	}
	if !exists {
		return errc.ErrNotFound.MultiMsg(in.ModuleName)
	}
	var validConn mongo.Conn
	for _, conn := range svc.cfg.Mongo.Conns {
		if m.Database == m.Database {
			validConn = conn
			break
		}
	}
	if validConn.Database == "" {
		return errc.ErrNotFound.MultiMsg(fmt.Sprintf("%s not found %s database", m.Name, m.Database))
	}
	host := strings.Join(mongo.URIToHosts(validConn.URI), ",")
	filter := bson.M{
		"module_name": m.Name,
		"host":        host,
		"db":          m.Database,
		"name": primitive.Regex{
			Pattern: fmt.Sprintf("%s_", m.LoggingPrefix()),
			Options: "i",
		},
	}
	docs, err := svc.d.FindCollStats(ctx, filter)
	if err != nil {
		return errc.ErrInternalErr.MultiErr(err)
	}

	list := make([]*model.CollStat, 0, len(docs))
	for _, v := range docs {
		d := &model.CollStat{
			ModuleName:     v.ModuleName,
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
			UpdatedTsSec:   v.UpdatedAt.Unix(),
			CreatedTsSec:   v.CreatedAt.Unix(),
		}
		list = append(list, d)
	}

	out.Count = int64(len(list))
	out.List = list

	return nil
}

// MetricsModuleList 模块列表
func (svc *Service) MetricsModuleList(ctx context.Context, in *model.MetricsModuleListReq, out *model.ListResp) error {
	c, docs, err := svc.d.FindMetricsModuleList(ctx, in)
	if err != nil {
		return errc.ErrInternalErr.MultiErr(err)
	}

	list := make([]*model.MetricsModuleList, 0, len(docs))
	for _, v := range docs {
		d := &model.MetricsModuleList{
			ModuleName:   v.ModuleName,
			Number:       v.Number,
			Size:         v.Size,
			CreatedTsSec: v.CreatedDate.Unix(),
		}
		list = append(list, d)
	}

	out.Count = c
	out.List = list
	return nil
}

// MetricsModuleTrend 模块日志趋势
func (svc *Service) MetricsModuleTrend(ctx context.Context, in *model.MetricsModuleTrendReq, out *model.MetricsModuleTrendResp) error {
	beforeDay := itime.BeforeDayDate(in.LastDay)
	filter := bson.M{
		"module_name":  in.ModuleName,
		"created_date": bson.M{"$gte": beforeDay},
	}
	docs, err := svc.d.FindMetricsModule(ctx, filter)
	if err != nil {
		return errc.ErrInternalErr.MultiErr(err)
	}

	number := int64(0)
	size := int64(0)
	ascTsNumbers := make([]model.TsNumbers, 0, in.LastDay*24)
	allLevels := make(map[types.Level]bool)
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
	levelMapData := map[types.Level][]int32{}
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

	levelSeries := make([]model.Serie, 0, len(levelMapData))
	for lvl, data := range levelMapData {
		levelSeries = append(levelSeries, model.Serie{
			Index: lvl.Int32(),
			Name:  lvl.String(),
			Type:  "bar",
			Color: levelColor(lvl),
			Data:  data,
		})
	}
	model.SortSeries(levelSeries, "ASC")
	for _, v := range levelSeries {
		legend = append(legend, v.Name)
	}
	ipSeries := make([]model.Serie, 0, len(ipMapData))
	for ip, data := range ipMapData {
		ip := strings.ReplaceAll(ip, "_", ".")
		legend = append(legend, ip)
		ipSeries = append(ipSeries, model.Serie{
			Index: int32(binary.BigEndian.Uint32([]byte(ip))),
			Name:  ip,
			Type:  "line",
			Color: ipColor(),
			Data:  data,
		})
	}
	model.SortSeries(ipSeries, "ASC")
	for _, v := range ipSeries {
		legend = append(legend, v.Name)
	}
	out.LegendData = legend
	out.LevelSeries = levelSeries
	out.IPSeries = ipSeries

	return nil
}

func levelColor(lvl types.Level) string {
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
