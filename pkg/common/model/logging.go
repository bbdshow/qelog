package model

import (
	"bytes"
	"strconv"
	"time"

	"github.com/huzhongqing/qelog/libs/mongo"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

const (
	MaxDBIndex          int32 = 16
	LoggingShardingTime       = "200601"
)

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
	TimeSec    int64              `bson:"ts"` // 秒, 用于建立秒级别索引
	MessageID  string             `bson:"mi"` // 如果重复写入，可以通过此ID区分
	Size       int                `bson:"-"`
}

func (l Logging) Key() string {
	buf := bytes.Buffer{}
	buf.WriteString(l.Module)
	buf.WriteString("_")
	buf.WriteString(l.Short)
	buf.WriteString("_")
	buf.WriteString(l.Level.String())
	return buf.String()
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

func LoggingCollectionName(dbIndex int32, unix int64) string {
	buff := bytes.Buffer{}
	buff.WriteString("logging_")
	buff.WriteString(strconv.Itoa(int(dbIndex)))
	buff.WriteString("_")
	y, m, _ := time.Unix(unix, 0).Date()
	buff.WriteString(strconv.Itoa(y))
	buff.WriteString(strconv.Itoa(int(m)))
	//name := fmt.Sprintf("logging_%d_%s",
	//	dbIndex, time.Unix(unix, 0).Format(LoggingShardingTime))
	return buff.String()
}

// 因为有分片的机制，那么同一collection下面，同一uniqueKey module 占多数情况。
// 所以时间可作为较大范围过滤，结合Limit可以较快返回
// 此索引因为时间粒度关系，存储会变得比较大
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
				// c2,c3 不建立索引，是优化索引大小
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
