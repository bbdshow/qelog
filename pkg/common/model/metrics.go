package model

import (
	"time"

	"github.com/huzhongqing/qelog/libs/mongo"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

const (
	CollectionNameModuleMetrics = "module_metrics"
)

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
