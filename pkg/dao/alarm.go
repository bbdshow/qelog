package dao

import (
	"context"
	"fmt"
	"github.com/bbdshow/bkit/db/mongo"
	"github.com/bbdshow/bkit/errc"
	"github.com/bbdshow/qelog/pkg/model"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
	"strings"
	"time"
)

func (d *Dao) FindAlarmRule(ctx context.Context, module string, enable bool) ([]*model.AlarmRule, error) {
	filter := bson.M{
		"module_name": module,
		"enable":      enable,
	}
	docs := make([]*model.AlarmRule, 0)
	err := d.adminInst.Find(ctx, model.CNAlarmRule, filter, &docs)
	return docs, errc.WithStack(err)
}

func (d *Dao) FindAlarmRuleList(ctx context.Context, in *model.FindAlarmRuleListReq) (int64, []*model.AlarmRule, error) {
	filter := bson.M{}
	if in.ModuleName != "" {
		filter["module_name"] = primitive.Regex{
			Pattern: in.ModuleName,
			Options: "i",
		}
	}
	if in.Enable > 0 {
		filter["enable"] = in.Enable == 1
	}

	if in.Short != "" {
		filter["short"] = primitive.Regex{
			Pattern: in.Short,
			Options: "i",
		}
	}

	opt := in.SetPage(options.Find()).SetSort(bson.M{"_id": -1})
	docs := make([]*model.AlarmRule, 0, in.Limit)
	c, err := d.adminInst.FindCount(ctx, model.CNAlarmRule, filter, &docs, opt, nil)
	return c, docs, errc.WithStack(err)
}

func (d *Dao) CreateAlarmRule(ctx context.Context, in *model.AlarmRule) error {
	_, err := d.adminInst.Collection(model.CNAlarmRule).InsertOne(ctx, in)
	return errc.WithStack(err)
}

func (d *Dao) GetAlarmRule(ctx context.Context, filter bson.M) (bool, *model.AlarmRule, error) {
	doc := &model.AlarmRule{}
	exists, err := d.adminInst.FindOne(ctx, model.CNAlarmRule, filter, doc)
	return exists, doc, errc.WithStack(err)
}

func (d *Dao) UpdateAlarmRule(ctx context.Context, in *model.UpdateAlarmRuleReq) error {
	id, err := in.ObjectID()
	if err != nil {
		return err
	}
	exists, doc, err := d.GetAlarmRule(ctx, bson.M{"_id": id})
	if err != nil {
		return errc.WithStack(err)
	}
	if !exists {
		return mongo.ErrNoDocuments
	}

	update := bson.M{}
	fields := bson.M{}
	if in.Enable != doc.Enable {
		fields["enable"] = in.Enable
	}
	if in.Short != doc.Short {
		fields["short"] = in.Short
	}
	if in.RateSec != doc.RateSec {
		fields["rate_sec"] = in.RateSec
	}
	if in.Level != doc.Level.Int32() {
		fields["level"] = in.Level
	}
	if in.Tag != doc.Tag {
		fields["tag"] = in.Tag
	}
	if in.Method != doc.Method.Int32() {
		fields["method"] = in.Method
	}
	if in.HookID != doc.HookID {
		fields["hook_id"] = in.HookID
	}

	if len(fields) > 0 {
		fields["updated_at"] = time.Now().Local()
		update["$set"] = fields
	}
	if len(update) == 0 {
		return nil
	}

	filter := bson.M{
		"_id":        id,
		"updated_at": doc.UpdatedAt,
	}
	uRet, err := d.adminInst.Collection(model.CNAlarmRule).UpdateOne(ctx, filter, update)
	if err != nil {
		return errc.WithStack(err)
	}
	if uRet.MatchedCount <= 0 {
		return mongo.ErrNotMatched
	}
	return nil
}

func (d *Dao) DelAlarmRule(ctx context.Context, filter bson.M) error {
	_, err := d.adminInst.Collection(model.CNAlarmRule).DeleteOne(ctx, filter)
	return errc.WithStack(err)
}

func (d *Dao) FindHookURLList(ctx context.Context, in *model.FindHookURLListReq) (int64, []*model.HookURL, error) {
	filter := bson.M{}
	if in.ID != "" {
		id, err := in.ObjectID()
		if err != nil {
			return 0, nil, err
		}
		filter["_id"] = id
	}
	if in.Name != "" {
		filter["name"] = in.Name
	}
	if in.KeyWord != "" {
		filter["key_word"] = in.KeyWord
	}
	if in.Method > 0 {
		filter["method"] = in.Method
	}

	opt := in.SetPage(options.Find()).SetSort(bson.M{"_id": -1})
	docs := make([]*model.HookURL, 0, in.Limit)
	c, err := d.adminInst.FindCount(ctx, model.CNHookURL, filter, &docs, opt, nil)
	return c, docs, errc.WithStack(err)
}

func (d *Dao) FindAllHookURL(ctx context.Context) ([]*model.HookURL, error) {
	docs := make([]*model.HookURL, 0)
	err := d.adminInst.Find(ctx, model.CNHookURL, bson.M{}, &docs)
	return docs, errc.WithStack(err)
}

func (d *Dao) GetHookURL(ctx context.Context, filter bson.M) (bool, *model.HookURL, error) {
	doc := &model.HookURL{}
	exists, err := d.adminInst.FindOne(ctx, model.CNHookURL, filter, doc)
	return exists, doc, errc.WithStack(err)
}

func (d *Dao) CreateHookURL(ctx context.Context, in *model.HookURL) error {
	_, err := d.adminInst.Collection(model.CNHookURL).InsertOne(ctx, in)
	return errc.WithStack(err)
}

func (d *Dao) UpdateHookURL(ctx context.Context, in *model.UpdateHookURLReq) error {
	id, err := in.ObjectID()
	if err != nil {
		return err
	}

	exists, doc, err := d.GetHookURL(ctx, bson.M{"_id": id})
	if err != nil {
		return errc.WithStack(err)
	}
	if !exists {
		return mongo.ErrNoDocuments
	}
	update := bson.M{}
	fields := bson.M{}
	if in.Name != doc.Name {
		fields["name"] = in.Name
	}

	if in.Method != doc.Method.Int32() {
		fields["method"] = in.Method
	}
	if in.URL != doc.URL {
		fields["url"] = in.URL
	}

	if in.KeyWord != doc.KeyWord {
		fields["key_word"] = in.KeyWord
	}

	if strings.Join(in.HideText, ",") != strings.Join(doc.HideText, ",") {
		fields["hide_text"] = in.HideText
	}

	if len(fields) > 0 {
		fields["updated_at"] = time.Now().Local()
		update["$set"] = fields
	}
	if len(update) == 0 {
		return nil
	}

	filter := bson.M{
		"_id":        id,
		"updated_at": doc.UpdatedAt,
	}
	// 更新已经引用的
	_, err = d.adminInst.Collection(model.CNAlarmRule).UpdateMany(ctx, bson.M{"hook_id": id.Hex()},
		bson.M{"$set": bson.M{"updated_at": time.Now()}})
	if err != nil {
		return errc.WithStack(err)
	}

	uRet, err := d.adminInst.Collection(model.CNHookURL).UpdateOne(ctx, filter, update)
	if err != nil {
		return errc.WithStack(err)
	}
	if uRet.MatchedCount <= 0 {
		return mongo.ErrNotMatched
	}
	return nil
}

func (d *Dao) DelHookURL(ctx context.Context, in *model.DelHookURLReq) error {
	id, err := in.ObjectID()
	if err != nil {
		return err
	}
	// 是否存在绑定
	c, err := d.adminInst.Collection(model.CNAlarmRule).CountDocuments(ctx, bson.M{"hook_id": in.ID}, options.Count().SetLimit(1))
	if err != nil {
		return errc.WithStack(err)
	}
	if c != 0 {
		return fmt.Errorf("hook url be referenced")
	}
	_, err = d.adminInst.Collection(model.CNHookURL).DeleteOne(ctx, bson.M{"_id": id})
	return errc.WithStack(err)
}
