package model

import (
	"fmt"

	"github.com/huzhongqing/qelog/infra/mongo"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

const (
	LoggingShardingTime = "200601"
)

// 分片索引容量
// 配置化把存储对象，映射到不同的索引下
// 举例：如果存在4个存储对象，则分配规则为 [1,2] = db1 [3,4] = db2 ...类推
// 通过此类设计，实现一个简单的存储横向扩展。
// 横向扩展时，应在原有基础上增加此值，预留给扩展的DB实例，这样以前的数据可不用迁移
var (
	ShardingIndexSize int = 8
)

func SetShardingIndexSize(size int) {
	ShardingIndexSize = size
}

type Logging struct {
	ID         primitive.ObjectID `bson:"_id,omitempty"`
	Module     string             `bson:"m"`
	IP         string             `bson:"ip"`
	Level      Level              `bson:"l"`
	Short      string             `bson:"s"`
	Full       string             `bson:"f"`
	Condition1 string             `bson:"c1"`
	Condition2 string             `bson:"c2"`
	Condition3 string             `bson:"c3"`
	TraceID    string             `bson:"ti"`
	TimeMill   int64              `bson:"tm"` // 日志打印时间
	TimeSec    int64              `bson:"ts"` // 秒, 用于建立秒级别索引, ts 返回结果排序, 所以会存在毫秒级别一定的误差
	MessageID  string             `bson:"mi"` // 如果重复写入，可以通过此ID忽略返回结果
	Size       int                `bson:"-"`
}

func (l Logging) Key() string {
	return fmt.Sprintf("%s_%s_%s", l.Module, l.Short, l.Level)
}

type Level int32

func (lvl Level) Int32() int32 {
	return int32(lvl)
}
func (lvl Level) String() string {
	v := "UNKNOWN"
	switch lvl {
	case -1:
		v = "DEBUG"
	case 0:
		v = "INFO"
	case 1:
		v = "WARN"
	case 2:
		v = "ERROR"
	case 3:
		v = "DPANIC"
	case 4:
		v = "PANIC"
	case 5:
		v = "FATAL"
	}
	return v
}

// 多条件查询只建立了一个联合索引，减少索引大小，提升写入速度
// 结合查询条件限制，保证此联合索引命中
func LoggingIndexMany(collectionName string) []mongo.Index {
	return []mongo.Index{
		{
			Collection: collectionName,
			Keys: bson.D{
				// m, ts 是必要查询条件，所以放在最前面
				{
					Key: "m", Value: 1,
				},
				{
					Key: "ts", Value: 1,
				},
				// level 与 short 一般作为常用可选查询，建立索引,
				// level筛选频率较高，同时索引的大小和速度比较平均
				{
					Key: "l", Value: 1,
				},
				{
					Key: "s", Value: 1,
				},
				// 条件索引，一般前面筛选后，还有大量日志，才会用到条件筛选，
				// 且查询语句不能跳跃条件查询
				// 正确示例 c1 c2 c3 或 c1 或 c1 c2
				{
					Key: "c1", Value: 1,
				},
				// c2,c3 不建立索引，是优化索引大小，及写入速度
				//{
				//	Key: "c2", Value: 1,
				//},
				//{
				//	Key: "c3", Value: 1,
				//},
			},
			Background: true,
		},
		{
			Collection: collectionName,
			Keys: bson.D{
				// trace_id 作为单独索引，当排查问题作为查询条件更快
				{Key: "ti", Value: -1},
			},
			Background: true,
		},
	}
}
