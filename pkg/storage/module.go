package storage

import (
	"context"

	"go.mongodb.org/mongo-driver/mongo/options"

	"github.com/huzhongqing/qelog/pkg/common/entity"

	"go.mongodb.org/mongo-driver/bson/primitive"

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

func (store *Store) FindModuleList(ctx context.Context, in *entity.FindModuleListReq) (int64, []*model.Module, error) {
	filter := bson.M{}
	if in.Name != "" {
		filter["name"] = primitive.Regex{
			Pattern: in.Name,
			Options: "i",
		}
	}
	opt := options.Find()
	in.SetPage(opt)
	opt.SetSort(bson.M{"_id": -1})
	docs := make([]*model.Module, 0, in.PageSize)
	c, err := store.database.FindAndCount(ctx, store.database.Collection(model.CollectionNameModule), filter, &docs, opt)
	return c, docs, err
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
