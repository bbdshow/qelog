package storage

import (
	"context"
)

func (store *Store) InsertManyLogging(ctx context.Context, name string, docs []interface{}) error {
	_, err := store.database.Collection(name).InsertMany(ctx, docs)
	return WrapError(err)
}
