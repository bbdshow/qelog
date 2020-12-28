package storage

import (
	"context"

	"github.com/huzhongqing/qelog/pkg/common/model"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// module 注册不会太多，超过 500 个，就不太适合此系统了
func (store *Store) FindAllModule(ctx context.Context) ([]*model.Module, error) {
	docs := make([]*model.Module, 0)
	coll := store.database.Collection(model.CollectionNameModule)
	err := store.database.Find(ctx, coll, bson.M{}, &docs)
	return docs, err
}

func (store *Store) FindModuleList(ctx context.Context, filter bson.M, result interface{}, opt *options.FindOptions) (int64, error) {
	c, err := store.database.FindAndCount(ctx, store.database.Collection(model.CollectionNameModule), filter, result, opt)
	return c, err
}

func (store *Store) InsertModule(ctx context.Context, doc *model.Module) error {
	_, err := store.database.Collection(doc.CollectionName()).InsertOne(ctx, doc)
	return err
}

func (store *Store) FindOneModule(ctx context.Context, filter bson.M, doc *model.Module) (bool, error) {
	return store.database.FindOne(ctx, store.database.Collection(doc.CollectionName()), filter, doc)
}

func (store *Store) UpdateModule(ctx context.Context, filter, update bson.M) error {
	uRet, err := store.database.Collection(model.CollectionNameModule).UpdateOne(ctx, filter, update)
	if err != nil {
		return err
	}
	if uRet.MatchedCount <= 0 {
		return ErrNotMatched
	}
	return nil
}

func (store *Store) DeleteModule(ctx context.Context, id primitive.ObjectID) error {
	filter := bson.M{
		"_id": id,
	}
	_, err := store.database.Collection(model.CollectionNameModule).DeleteOne(ctx, filter)
	return err
}
