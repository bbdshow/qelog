package model

import (
	"github.com/bbdshow/bkit/db/mongo"
	"github.com/bbdshow/qelog/common/types"
	"strings"
	"sync/atomic"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

const (
	CNModuleMetrics = "module_metrics"
	CNDBStats       = "db_stats"
	CNCollStats     = "coll_stats"
)

// ModuleMetrics 模块写入日志数据分布统计
type ModuleMetrics struct {
	ID          primitive.ObjectID `bson:"_id,omitempty"`
	ModuleName  string             `bson:"module_name"`
	Number      int64              `bson:"number"`
	Size        int64              `bson:"size"`
	Sections    map[int64]Numbers  `bson:"sections"`     // key 小时精度的时间戳
	CreatedDate time.Time          `bson:"created_date"` // 创建日期
}

type Numbers struct {
	Sum    int32                 `bson:"sum"`
	Levels map[types.Level]int32 `bson:"levels"`
	IPs    map[string]int32      `bson:"ips"`
}

type TsNumbers struct {
	Ts int64
	Numbers
}

type AscTsNumbers []TsNumbers

func (v AscTsNumbers) Len() int           { return len(v) }
func (v AscTsNumbers) Swap(i, j int)      { v[i], v[j] = v[j], v[i] }
func (v AscTsNumbers) Less(i, j int) bool { return v[i].Ts < v[j].Ts }

func (mm ModuleMetrics) CollectionName() string {
	return CNModuleMetrics
}

func ModuleMetricsIndexMany() []mongo.Index {
	return []mongo.Index{
		{
			Collection: CNModuleMetrics,
			Keys: bson.D{
				{
					Key: "module_name", Value: 1,
				},
			},
			Background: true,
		},
		{
			Collection: CNModuleMetrics,
			Keys: bson.D{
				{
					Key: "number", Value: 1,
				},
			},
			Background: true,
		},
		{
			Collection: CNModuleMetrics,
			Keys: bson.D{
				{
					Key: "size", Value: 1,
				},
			},
			Background: true,
		},
		// ttl 30天
		{
			Collection: CNModuleMetrics,
			Keys: bson.D{
				{
					Key: "created_date", Value: 1,
				},
			},
			Background:         true,
			ExpireAfterSeconds: 86400 * 30,
		},
	}
}

// DBStats 数据库容量统计
type DBStats struct {
	ID          primitive.ObjectID `bson:"_id,omitempty"`
	Host        string             `bson:"host"`
	DB          string             `bson:"db"`
	Collections int32              `bson:"collections"`
	Objects     int64              `bson:"objects"`
	DataSize    int64              `bson:"data_size"`
	StorageSize int64              `bson:"storage_size"`
	Indexes     int64              `bson:"indexes"`
	IndexSize   int64              `bson:"index_size"`
	CreatedAt   time.Time          `bson:"created_at"`
}

func (ds DBStats) CollectionName() string {
	return CNDBStats
}

func DBStatsIndexMany() []mongo.Index {
	return []mongo.Index{
		{
			Collection: CNDBStats,
			Keys: bson.D{
				{
					Key: "host", Value: 1,
				},
				{
					Key: "db", Value: 1,
				},
			},
			Background: true,
		},
	}
}

// CollStats 集合容量统计
type CollStats struct {
	ID             primitive.ObjectID `bson:"_id,omitempty"`
	Host           string             `bson:"host"`
	DB             string             `bson:"db"`
	Name           string             `bson:"name"`
	Size           int64              `bson:"size"`
	Count          int64              `bson:"count"`
	AvgObjSize     int64              `bson:"avg_obj_size"`
	StorageSize    int64              `bson:"storage_size"`
	Capped         bool               `bson:"capped"`
	TotalIndexSize int64              `bson:"total_index_size"`
	IndexSizes     map[string]int64   `bson:"index_sizes"`
	CreatedAt      time.Time          `bson:"created_at"`
}

func (cs CollStats) CollectionName() string {
	return CNCollStats
}

func CollStatsIndexMany() []mongo.Index {
	return []mongo.Index{
		{
			Collection: CNCollStats,
			Keys: bson.D{
				{
					Key: "host", Value: 1,
				},
				{
					Key: "db", Value: 1,
				},
				{
					Key: "name", Value: 1,
				},
			},
			Background: true,
		},
	}
}

type MetricsState struct {
	Date           time.Time
	Section        int64
	ModuleName     string
	Number         int32
	Size           int32
	Levels         map[types.Level]int32
	IPs            map[string]int32
	IncIntervalSec int64
}

func (s *MetricsState) IncrNumber(n int32) {
	atomic.AddInt32(&s.Number, n)
}
func (s *MetricsState) IncrSize(n int32) {
	atomic.AddInt32(&s.Size, n)
}

func (s *MetricsState) IncrLevel(lvl types.Level, n int32) {
	v, ok := s.Levels[lvl]
	if ok {
		s.Levels[lvl] = v + n
		return
	}
	s.Levels[lvl] = n
}

func (s *MetricsState) IncrIP(ip string, n int32) {
	if ip == "" {
		return
	}
	strs := strings.Split(ip, ".")
	if len(strs) <= 1 {
		// ipv6
		strs = strings.Split(ip, ":")
	}
	// 使用 _ 链接，便于mongodb更新
	ip = strings.Join(strs, "_")

	v, ok := s.IPs[ip]
	if ok {
		s.IPs[ip] = v + n
		return
	}
	s.IPs[ip] = n
}

func (s *MetricsState) IsIncr() bool {
	// 超过一定时间，就可以写入了
	return time.Now().Unix()-s.Section >= s.IncIntervalSec
}
