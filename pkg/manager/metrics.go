package manager

import (
	"context"
	"encoding/binary"
	"fmt"
	"math/rand"
	"sort"
	"strings"
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
		"created_date": date,
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

func (srv *Service) MetricsModuleTrend(ctx context.Context, in *entity.MetricsModuleTrendReq, out *entity.MetricsModuleTrendResp) error {
	now := time.Now()
	lastDay := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, time.Local).AddDate(0, 0, -in.LastDay)
	filter := bson.M{
		"module_name":  in.ModuleName,
		"created_date": bson.M{"$gte": lastDay},
	}
	fmt.Println(filter)
	docs := make([]*model.ModuleMetrics, 0, in.LastDay)
	if err := srv.store.FindModuleMetrics(ctx, filter, &docs, nil); err != nil {
		return httputil.ErrSystemException.MergeError(err)
	}

	number := int64(0)
	size := int64(0)
	ascTsNumbers := make([]model.TsNumbers, 0, in.LastDay*24)
	for _, v := range docs {
		number += v.Number
		size += v.Size
		for ts, numbers := range v.Sections {
			ascTsNumbers = append(ascTsNumbers, model.TsNumbers{
				Ts:      ts,
				Numbers: numbers,
			})
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
		for lvl, num := range v.Levels {
			data, ok := levelMapData[lvl]
			if ok {
				levelMapData[lvl] = append(data, num)
			} else {
				data = make([]int32, 0, in.LastDay*24)
				levelMapData[lvl] = append(data, num)
			}
		}
		for ip, num := range v.IPs {
			data, ok := ipMapData[ip]
			if ok {
				ipMapData[ip] = append(data, num)
			} else {
				data = make([]int32, 0, in.LastDay*24)
				ipMapData[ip] = append(data, num)
			}
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
	return fmt.Sprintf("rgba(%d,%d,%d,1)", rand.Int31n(100), rand.Int31n(255), rand.Int31n(255))
}
