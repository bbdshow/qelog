package mongoutil

import (
	"context"

	"go.mongodb.org/mongo-driver/bson"

	"github.com/huzhongqing/qelog/libs/mongo"
)

/*
collStats
dbStats
hostInfo
*/
type MongodbUtil struct {
	database *mongo.Database
}

func NewMongodbUtil(database *mongo.Database) *MongodbUtil {
	return &MongodbUtil{database: database}
}

type CollStatsResp struct {
	Ok             int              `json:"ok"`
	Ns             string           `json:"ns"`
	Size           int64            `json:"size"`
	Count          int64            `json:"count"`
	AvgObjSize     int64            `json:"avgObjSize"`
	StorageSize    int64            `json:"storageSize"`
	Capped         bool             `json:"capped"`
	TotalIndexSize int64            `json:"totalIndexSize"`
	IndexSizes     map[string]int64 `json:"indexSizes"`
}

func (mu *MongodbUtil) CollStats(ctx context.Context, colls []string) ([]CollStatsResp, error) {
	out := make([]CollStatsResp, 0, len(colls))
	for _, coll := range colls {
		stats := CollStatsResp{}
		err := mu.database.RunCommand(ctx, bson.D{{Key: "collStats", Value: coll}}).Decode(&stats)
		if err != nil {
			return nil, err
		}
		out = append(out, stats)
	}
	return out, nil
}

type DBStatsResp struct {
	Ok          int32  `json:"ok"`
	DB          string `json:"db"`
	Collections int32  `json:"collections"`
	Objects     int64  `json:"objects"`
	DataSize    int64  `json:"dataSize"`
	StorageSize int64  `json:"storageSize"`
	Indexes     int32  `json:"indexes"`
	IndexSize   int64  `json:"indexSize"`
}

func (mu *MongodbUtil) DBStats(ctx context.Context) (DBStatsResp, error) {
	out := DBStatsResp{}
	err := mu.database.RunCommand(ctx, bson.D{{Key: "dbStats", Value: 1}}).Decode(&out)
	return out, err
}

type HostInfoResp struct {
	System map[string]interface{} `json:"system"`
	Os     map[string]interface{} `json:"os"`
}

func (mu *MongodbUtil) HostInfo(ctx context.Context) (HostInfoResp, error) {
	out := HostInfoResp{}
	err := mu.database.RunCommand(ctx, bson.D{{Key: "hostInfo", Value: 1}}).Decode(&out)
	return out, err
}
