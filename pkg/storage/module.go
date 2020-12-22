package storage

import (
	"context"

	"go.mongodb.org/mongo-driver/bson"

	"github.com/huzhongqing/qelog/pkg/common/model"
)

// module 注册不会太多，超过 500 个，就不太适合此系统了
func (store *Store) FindAllModuleRegister(ctx context.Context) ([]*model.ModuleRegister, error) {
	docs := make([]*model.ModuleRegister, 0)
	coll := store.database.Collection(model.CollectionNameModuleRegister)
	err := store.database.Find(ctx, coll, bson.M{}, &docs)
	return docs, err
}
