package storage

import (
	"context"
	"github.com/huzhongqing/qelog/infra/mongo"

	"github.com/huzhongqing/qelog/pkg/common/model"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Module struct {
	db *mongo.Database
}

func NewModule(db *mongo.Database) *Module {
	return &Module{db: db}
}

func (m *Module) FindCountModule(ctx context.Context, filter bson.M, opts ...*options.FindOptions) (int64, []*model.Module, error) {
	docs := make([]*model.Module, 0)
	c, err := m.db.FindAndCount(ctx, m.db.Collection(model.CollectionNameModule), filter, &docs, opts...)
	return c, docs, err
}

func (m *Module) InsertModule(ctx context.Context, doc *model.Module) error {
	_, err := m.db.Collection(doc.CollectionName()).InsertOne(ctx, doc)
	return err
}

func (m *Module) FindOneModule(ctx context.Context, filter bson.M, doc *model.Module, opts ...*options.FindOneOptions) (bool, error) {
	ok, err := m.db.FindOne(ctx, m.db.Collection(doc.CollectionName()), filter, doc, opts...)
	return ok, err
}

func (m *Module) UpdateModule(ctx context.Context, filter, update bson.M, opts ...*options.UpdateOptions) error {
	uRet, err := m.db.Collection(model.CollectionNameModule).UpdateOne(ctx, filter, update, opts...)
	if err != nil {
		return err
	}
	if uRet.MatchedCount <= 0 {
		return mongo.ErrNotMatched
	}
	return nil
}

func (m *Module) DeleteModule(ctx context.Context, filter bson.M) error {
	_, err := m.db.Collection(model.CollectionNameModule).DeleteOne(ctx, filter)
	return err
}
