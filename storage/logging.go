package storage

import (
	"context"
	"time"
)

func (store *Store) InsertManyLogging(collectionName string, v []interface{}) error {
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	_, err := store.database.Collection(collectionName).InsertMany(ctx, v)
	return err
}
