package storage

import (
	"context"

	"github.com/huzhongqing/qelog/model"
)

func (store *Store) InsertManyLogging(ctx context.Context, bucket string, docs []*model.Logging) error {

	//_, err := store.database.Collection().InsertMany(ctx, docs)
	return nil
}

func (store *Store) shardingLogging(bucket string, docs []*model.Logging) map[string][]interface{} {
	out := make(map[string][]interface{})
	for _, v := range docs {
		name := store.sharding.GenerateName(bucket, v.MillTimeStamp/1e3)
		val, ok := out[name]
		if ok {
			out[name] = append(val, v)
		} else {
			out[name] = []interface{}{v}
		}
	}
	return out
}
