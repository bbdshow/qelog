package storage

import (
	"context"

	"go.mongodb.org/mongo-driver/bson"

	"github.com/huzhongqing/qelog/pkg/common/model"
)

// module 注册不会太多，超过 500 个，就不太适合此系统了
func (store *Store) FindAllModule(ctx context.Context) ([]*model.Module, error) {
	docs := make([]*model.Module, 0)
	coll := store.database.Collection(model.CollectionNameModule)
	err := store.database.Find(ctx, coll, bson.M{}, &docs)
	return docs, err
}

func (store *Store) InsertModule(ctx context.Context, doc *model.Module) error {
	_, err := store.database.Collection(doc.CollectionName()).InsertOne(ctx, doc)
	return err
}
