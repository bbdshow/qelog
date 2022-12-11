package dao

import (
	"context"
	"time"

	"github.com/bbdshow/bkit/db/mongo"
	"github.com/bbdshow/bkit/errc"
	"github.com/bbdshow/qelog/pkg/model"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// FindModuleList common db CRUD operation
func (d *Dao) FindModuleList(ctx context.Context, in *model.FindModuleListReq) (int64, []*model.Module, error) {
	filter := bson.M{}
	if in.Name != "" {
		filter["name"] = primitive.Regex{
			Pattern: in.Name,
			Options: "i",
		}
	}
	opt := in.SetPage(options.Find()).SetSort(bson.M{"_id": -1})
	docs := make([]*model.Module, 0, in.Limit)
	c, err := d.adminInst.FindCount(ctx, model.CNModule, filter, &docs, opt, nil)
	return c, docs, errc.WithStack(err)
}

// FindAllModule common db CRUD operation
func (d *Dao) FindAllModule(ctx context.Context) ([]*model.Module, error) {
	docs := make([]*model.Module, 0)
	err := d.adminInst.Find(ctx, model.CNModule, bson.M{}, &docs)
	return docs, errc.WithStack(err)
}

// CreateModule common db CRUD operation
func (d *Dao) CreateModule(ctx context.Context, doc *model.Module) error {
	_, err := d.adminInst.Collection(model.CNModule).InsertOne(ctx, doc)
	return errc.WithStack(err)
}

// GetModule common db CRUD operation
func (d *Dao) GetModule(ctx context.Context, filter bson.M) (bool, *model.Module, error) {
	doc := &model.Module{}
	exists, err := d.adminInst.FindOne(ctx, model.CNModule, filter, doc)
	return exists, doc, errc.WithStack(err)
}

// UpdateModule common db CRUD operation
func (d *Dao) UpdateModule(ctx context.Context, in *model.UpdateModuleReq) error {
	id, err := in.ObjectID()
	if err != nil {
		return err
	}
	exists, doc, err := d.GetModule(ctx, bson.M{"_id": id})
	if err != nil {
		return errc.WithStack(err)
	}
	if !exists {
		return mongo.ErrNoDocuments
	}

	update := bson.M{}
	fields := bson.M{}

	if doc.DaySpan != in.DaySpan {
		fields["day_span"] = in.DaySpan
	}
	if doc.MaxMonth != in.MaxMonth {
		fields["max_month"] = in.MaxMonth
	}
	if doc.Bucket != in.Bucket {
		fields["bucket"] = in.Bucket
	}

	if doc.Desc != in.Desc {
		fields["desc"] = in.Desc
	}
	if doc.Database != in.Database {
		fields["database"] = in.Database
	}
	if doc.Prefix != in.Prefix {
		fields["prefix"] = in.Prefix
	}

	if len(fields) > 0 {
		fields["updated_at"] = time.Now().Local()
		update["$set"] = fields
	}
	if len(update) == 0 {
		return nil
	}
	filter := bson.M{
		"_id":        doc.ID,
		"updated_at": doc.UpdatedAt,
	}

	uRet, err := d.adminInst.Collection(model.CNModule).UpdateOne(ctx, filter, update)
	if err != nil {
		return errc.WithStack(err)
	}
	if uRet.MatchedCount <= 0 {
		return mongo.ErrNotMatched
	}
	return nil
}

// DelModule common db CRUD operation
func (d *Dao) DelModule(ctx context.Context, filter bson.M) error {
	_, err := d.adminInst.Collection(model.CNModule).DeleteOne(ctx, filter)
	return errc.WithStack(err)
}
