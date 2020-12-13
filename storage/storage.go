package storage

import "github.com/huzhongqing/qelog/model/mongoclient"

type Store struct {
	database *mongoclient.Database
}

func New(database *mongoclient.Database) *Store {
	store := &Store{database: database}
	return store
}
