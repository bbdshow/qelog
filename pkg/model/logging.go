package model

import (
	"fmt"

	"github.com/bbdshow/bkit/db/mongo"
	"github.com/bbdshow/qelog/pkg/types"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Logging struct {
	ID         primitive.ObjectID `bson:"_id,omitempty"`
	Module     string             `bson:"m"`
	IP         string             `bson:"ip"`
	Level      types.Level        `bson:"l"`
	Short      string             `bson:"s"`
	Full       string             `bson:"f"`
	Condition1 string             `bson:"c1"`
	Condition2 string             `bson:"c2"`
	Condition3 string             `bson:"c3"`
	TraceID    string             `bson:"ti"`
	TimeMill   int64              `bson:"tm"` // logging print time, /mill
	TimeSec    int64              `bson:"ts"` // /sec ,used to create index, order by ts. if used tm created,index too large
	MessageID  string             `bson:"mi"` // logging idempotent
	Size       int                `bson:"-"`
}

func (l Logging) Key() string {
	return fmt.Sprintf("%s_%s_%s", l.Module, l.Short, l.Level)
}

// LoggingIndexMany
// only one joint index is created. reduce the index size, written performance
// optimize query conditions to ensure index matching.
func LoggingIndexMany(collectionName string) []mongo.Index {
	return []mongo.Index{
		{
			Collection: collectionName,
			Keys: bson.D{
				// m, ts required condition, put it front
				{
					Key: "m", Value: 1,
				},
				{
					Key: "ts", Value: 1,
				},
				// level as a priority filter
				{
					Key: "l", Value: 1,
				},
				{
					Key: "s", Value: 1,
				},
				// condition rule combined eg:[c1]  [c1 & c2] [c1 & c2 & c3]
				{
					Key: "c1", Value: 1,
				},
				// c2,c3 no indexing, because data filter little
			},
			Background: true,
		},
		{
			Collection: collectionName,
			Keys: bson.D{
				// trace_id as single index,improve query performance.
				{Key: "ti", Value: -1},
			},
			Background: true,
		},
	}
}
