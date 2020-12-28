package storage

import (
	"context"

	"github.com/huzhongqing/qelog/pkg/common/model"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func (store *Store) UpsertModuleMetrics(ctx context.Context, filter, update bson.M, opt *options.UpdateOptions) error {
	_, err := store.database.Collection(model.CollectionNameModuleMetrics).UpdateOne(ctx, filter, update, opt)
	return err
}
