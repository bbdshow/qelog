package storage

import (
	"github.com/huzhongqing/qelog/libs/mongoclient"
	"github.com/huzhongqing/qelog/libs/sharding"
)

type Store struct {
	database *mongoclient.Database
	sharding *sharding.Sharding
}

func New(database *mongoclient.Database) *Store {
	store := &Store{
		database: database,
		sharding: sharding.NewSharding(sharding.FormatMonth, "logging"),
	}
	return store
}
