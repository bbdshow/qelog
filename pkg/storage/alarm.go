package storage

//import (
//	"context"
//	"github.com/bbdshow/qelog/infra/mongo"
//
//	"github.com/bbdshow/qelog/pkg/common/model"
//	"go.mongodb.org/mongo-driver/bson"
//	"go.mongodb.org/mongo-driver/mongo/options"
//)
//
//type AlarmRule struct {
//	db *mongo.Database
//}
//
//func NewAlarmRule(db *mongo.Database) *AlarmRule {
//	return &AlarmRule{db: db}
//}
//
//// 当报警规则超过上千的规则，可优化语句
//func (ar *AlarmRule) FindAlarmRule(ctx context.Context, filter bson.M, opts ...*options.FindOptions) ([]*model.AlarmRule, error) {
//	docs := make([]*model.AlarmRule, 0)
//	err := ar.db.Find(ctx, ar.db.Collection(model.CollectionNameAlarmRule), filter, &docs, opts...)
//	return docs, err
//}
//
//func (ar *AlarmRule) FindCountAlarmRule(ctx context.Context, filter bson.M, opts ...*options.FindOptions) (int64, []*model.AlarmRule, error) {
//	docs := make([]*model.AlarmRule, 0)
//	c, err := ar.db.FindAndCount(ctx, ar.db.Collection(model.CollectionNameAlarmRule), filter, &docs, opts...)
//	return c, docs, err
//}
//
//func (ar *AlarmRule) InsertAlarmRule(ctx context.Context, doc *model.AlarmRule) error {
//	_, err := ar.db.Collection(doc.CollectionName()).InsertOne(ctx, doc)
//	return err
//}
//
//func (ar *AlarmRule) FindOneAlarmRule(ctx context.Context, filter bson.M, doc *model.AlarmRule, opts ...*options.FindOneOptions) (bool, error) {
//	return ar.db.FindOne(ctx, ar.db.Collection(doc.CollectionName()), filter, doc, opts...)
//}
//
//func (ar *AlarmRule) UpdateAlarmRule(ctx context.Context, filter, update bson.M, opts ...*options.UpdateOptions) error {
//	uRet, err := ar.db.Collection(model.CollectionNameAlarmRule).UpdateOne(ctx, filter, update, opts...)
//	if err != nil {
//		return err
//	}
//	if uRet.MatchedCount <= 0 {
//		return mongo.ErrNotMatched
//	}
//	return nil
//}
//
//func (ar *AlarmRule) UpdateManyAlarmRule(ctx context.Context, filter, update bson.M, opts ...*options.UpdateOptions) error {
//	_, err := ar.db.Collection(model.CollectionNameAlarmRule).UpdateMany(ctx, filter, update, opts...)
//	if err != nil {
//		return err
//	}
//	return nil
//}
//
//func (ar *AlarmRule) DeleteAlarmRule(ctx context.Context, filter bson.M) error {
//	_, err := ar.db.Collection(model.CollectionNameAlarmRule).DeleteOne(ctx, filter)
//	return err
//}
//
//func (ar *AlarmRule) FindCountHookURL(ctx context.Context, filter bson.M, opts ...*options.FindOptions) (int64, []*model.HookURL, error) {
//	docs := make([]*model.HookURL, 0)
//	c, err := ar.db.FindAndCount(ctx, ar.db.Collection(model.CollectionNameHookURL), filter, &docs, opts...)
//	return c, docs, err
//}
//
//func (ar *AlarmRule) FindOneHookURL(ctx context.Context, filter bson.M, doc *model.HookURL, opts ...*options.FindOneOptions) (bool, error) {
//	return ar.db.FindOne(ctx, ar.db.Collection(doc.CollectionName()), filter, doc, opts...)
//}
//
//func (ar *AlarmRule) InsertHookURL(ctx context.Context, doc *model.HookURL) error {
//	_, err := ar.db.Collection(doc.CollectionName()).InsertOne(ctx, doc)
//	return err
//}
//
//func (ar *AlarmRule) UpdateHookURL(ctx context.Context, filter, update bson.M, opts ...*options.UpdateOptions) error {
//	uRet, err := ar.db.Collection(model.CollectionNameHookURL).UpdateOne(ctx, filter, update, opts...)
//	if err != nil {
//		return err
//	}
//	if uRet.MatchedCount <= 0 {
//		return mongo.ErrNotMatched
//	}
//	return nil
//}
//
//func (ar *AlarmRule) DelHookURL(ctx context.Context, filter bson.M) error {
//	_, err := ar.db.Collection(model.CollectionNameHookURL).DeleteOne(ctx, filter)
//	return err
//}
