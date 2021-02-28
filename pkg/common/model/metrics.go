package model

import (
	"time"

	"github.com/huzhongqing/qelog/infra/mongo"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

const (
	CollectionNameModuleMetrics = "module_metrics"
	CollectionNameDBStats       = "db_stats"
	CollectionNameCollStats     = "coll_stats"
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
	Sum    int32            `bson:"sum"`
	Levels map[Level]int32  `bson:"levels"`
	IPs    map[string]int32 `bson:"ips"`
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
	return CollectionNameModuleMetrics
}

func ModuleMetricsIndexMany() []mongo.Index {
	return []mongo.Index{
		{
			Collection: CollectionNameModuleMetrics,
			Keys: bson.D{
				{
					Key: "module_name", Value: 1,
				},
			},
			Background: true,
		},
		{
			Collection: CollectionNameModuleMetrics,
			Keys: bson.D{
				{
					Key: "number", Value: 1,
				},
			},
			Background: true,
		},
		{
			Collection: CollectionNameModuleMetrics,
			Keys: bson.D{
				{
					Key: "size", Value: 1,
				},
			},
			Background: true,
		},
		// ttl 30天
		{
			Collection: CollectionNameModuleMetrics,
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
	return CollectionNameDBStats
}

func DBStatsIndexMany() []mongo.Index {
	return []mongo.Index{
		{
			Collection: CollectionNameDBStats,
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
	return CollectionNameCollStats
}

func CollStatsIndexMany() []mongo.Index {
	return []mongo.Index{
		{
			Collection: CollectionNameCollStats,
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
