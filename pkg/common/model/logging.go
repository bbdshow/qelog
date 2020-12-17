package model

import (
	"github.com/huzhongqing/qelog/libs/mongo"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Logging struct {
	ID         primitive.ObjectID `bson:"_id,omitempty"`
	UniqueKey  string             `bson:"uk"`
	Module     string             `bson:"m"`
	IP         string             `bson:"ip"`
	Level      int                `bson:"l"`
	Short      string             `bson:"s"`
	Full       string             `json:"f"`
	Condition1 string             `bson:"c1"`
	Condition2 string             `bson:"c2"`
	Condition3 string             `bson:"c3"`
	Time       int64              `bson:"t"`  // 日志打印时间
	TimeStamp  int64              `bson:"ts"` // 秒, 用于建立秒级别索引

}

// 因为有分片的机制，那么同一collection下面，同一uniqueKey module 占多数情况。
// 所以时间可作为较大范围过滤，结合Limit可以较快返回
// 此索引因为时间粒度关系，存储会变得比较大
func LoggingIndexMany(collectionName string) []mongo.Index {
	return []mongo.Index{
		{
			Collection: collectionName,
			Keys: bson.M{
				"uk": 1,
				"m":  1,
				"ts": -1,
				// uk, m, t 是必要查询条件，所以放在最前面
				"l": 1,
				"s": 1,
				// level 与 short 一般作为常用可选查询，建立索引, level优先筛选范围更大
				"c1": 1,
				"c2": 1,
				"c3": 1,
				// 条件索引，一般前面筛选后，还有大量日志，才会用到条件筛选，
			},
			Background: true,
		},
	}
}
