package admin

import (
	"context"
	"fmt"
	"github.com/bbdshow/qelog/infra/mongo"
	"sort"
	"strings"
	"sync"
	"time"

	"github.com/bbdshow/qelog/infra/alert"
	"github.com/bbdshow/qelog/infra/httputil"
	"github.com/bbdshow/qelog/pkg/config"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"

	"go.mongodb.org/mongo-driver/bson"

	"github.com/bbdshow/qelog/pkg/common/entity"
	"github.com/bbdshow/qelog/pkg/common/model"
	"github.com/bbdshow/qelog/pkg/storage"
)

type Service struct {
	alarmRuleStore     *storage.AlarmRule
	moduleStore        *storage.Module
	moduleMetricsStore *storage.ModuleMetrics
	sc                 mongo.ShardingCollection
	once               sync.Once
}

func NewService() *Service {
	srv := &Service{
		alarmRuleStore:     storage.NewAlarmRule(model.MainDB),
		moduleStore:        storage.NewModule(model.MainDB),
		moduleMetricsStore: storage.NewModuleMetrics(model.MainDB),
		sc:                 mongo.NewShardingCollection("logging", config.Global.DaySpan),
	}
	srv.once.Do(func() {
		go srv.backgroundDelExpiredCollection(config.Global.MaxAgeMonth)
	})
	return srv
}

func (srv *Service) FindModuleList(ctx context.Context, in *entity.FindModuleListReq, out *entity.ListResp) error {
	filter := bson.M{}
	if in.Name != "" {
		filter["name"] = primitive.Regex{
			Pattern: in.Name,
			Options: "i",
		}
	}
	opt := in.SetPage(options.Find()).SetSort(bson.M{"_id": -1})
	c, docs, err := srv.moduleStore.FindCountModule(ctx, filter, opt)
	if err != nil {
		return httputil.ErrSystemException.MergeError(err)
	}

	out.Count = c
	list := make([]*entity.FindModuleList, 0, len(docs))
	for _, v := range docs {
		d := &entity.FindModuleList{
			ID:                   v.ID.Hex(),
			Name:                 v.Name,
			Desc:                 v.Desc,
			ShardingIndex:        v.ShardingIndex,
			HistoryShardingIndex: v.HistoryShardingIndex,
			UpdatedTsSec:         v.UpdatedAt.Unix(),
		}
		list = append(list, d)
	}
	out.List = list
	return nil
}

func (srv *Service) CreateModule(ctx context.Context, in *entity.CreateModuleReq) error {
	doc := &model.Module{
		Name:                 in.Name,
		Desc:                 in.Desc,
		ShardingIndex:        in.ShardingIndex,
		HistoryShardingIndex: make([]int, 0),
		UpdatedAt:            time.Now().Local(),
	}
	if err := srv.moduleStore.InsertModule(ctx, doc); err != nil {
		return httputil.ErrSystemException.MergeError(err)
	}
	return nil
}

func (srv *Service) UpdateModule(ctx context.Context, in *entity.UpdateModuleReq) error {
	id, err := in.ObjectID()
	if err != nil {
		return err
	}

	doc := &model.Module{}
	if ok, err := srv.moduleStore.FindOneModule(ctx, bson.M{"_id": id}, doc); err != nil {
		return httputil.ErrSystemException.MergeError(err)
	} else if !ok {
		return httputil.ErrNotFound
	}
	update := bson.M{}
	fields := bson.M{}
	if doc.ShardingIndex != in.ShardingIndex {
		fields["sharding_index"] = in.ShardingIndex
		update["$addToSet"] = bson.M{"history_sharding_index": doc.ShardingIndex}
	}
	if doc.Desc != in.Desc {
		fields["desc"] = in.Desc
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

	if err := srv.moduleStore.UpdateModule(ctx, filter, update); err != nil {
		return httputil.ErrSystemException.MergeError(err)
	}
	return nil
}

func (srv *Service) DeleteModule(ctx context.Context, in *entity.DeleteModuleReq) error {
	id, err := in.ObjectID()
	if err != nil {
		return err
	}
	doc := &model.Module{}
	if ok, err := srv.moduleStore.FindOneModule(ctx, bson.M{"_id": id}, doc); err != nil {
		return httputil.ErrSystemException.MergeError(err)
	} else if !ok {
		return httputil.ErrNotFound
	}
	if doc.Name != in.Name {
		return httputil.ErrNotFound
	}
	if err := srv.moduleStore.DeleteModule(ctx, bson.M{"_id": id}); err != nil {
		return httputil.ErrSystemException.MergeError(err)
	}
	return nil
}

type AscShardingIndexState []entity.ShardingIndexState

func (asc AscShardingIndexState) Len() int           { return len(asc) }
func (asc AscShardingIndexState) Swap(i, j int)      { asc[i], asc[j] = asc[j], asc[i] }
func (asc AscShardingIndexState) Less(i, j int) bool { return asc[i].Use < asc[j].Use }

func (srv *Service) GetShardingIndex(ctx context.Context, out *entity.GetShardingIndexResp) error {
	_, docs, err := srv.moduleStore.FindCountModule(ctx, bson.M{})
	if err != nil {
		return httputil.ErrSystemException.MergeError(err)
	}
	state := make(map[int]int)
	for i := 1; i <= model.ShardingIndexSize; i++ {
		state[i] = 0
	}
	for _, v := range docs {
		num, ok := state[v.ShardingIndex]
		if ok {
			state[v.ShardingIndex] = num + 1
		}
	}

	states := make([]entity.ShardingIndexState, 0, len(state))
	for k, v := range state {
		states = append(states, entity.ShardingIndexState{
			Index: k,
			Use:   v,
		})
	}
	sort.Sort(AscShardingIndexState(states))

	// 找到最小的，作为推荐
	suggestDBIndex := model.ShardingIndexSize
	if len(states) > 0 {
		suggestDBIndex = states[0].Index
	}

	out.SuggestIndex = suggestDBIndex
	out.ShardingIndexSize = model.ShardingIndexSize
	out.UseState = states

	return nil
}

func (srv *Service) FindAlarmRuleList(ctx context.Context, in *entity.FindAlarmRuleListReq, out *entity.ListResp) error {
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
		filter["short"] = in.Short
	}

	opt := in.SetPage(options.Find()).SetSort(bson.M{"_id": -1})

	c, docs, err := srv.alarmRuleStore.FindCountAlarmRule(ctx, filter, opt)
	if err != nil {
		return httputil.ErrSystemException.MergeError(err)
	}

	out.Count = c
	list := make([]*entity.FindAlarmRuleList, 0, len(docs))
	for _, v := range docs {
		d := &entity.FindAlarmRuleList{
			ID:           v.ID.Hex(),
			Enable:       v.Enable,
			ModuleName:   v.ModuleName,
			Short:        v.Short,
			Level:        v.Level.Int32(),
			Tag:          v.Tag,
			RateSec:      v.RateSec,
			Method:       v.Method.Int32(),
			HookID:       v.HookID,
			UpdatedTsSec: v.UpdatedAt.Unix(),
		}
		list = append(list, d)
	}
	out.List = list

	return nil
}

func (srv *Service) CreateAlarmRule(ctx context.Context, in *entity.CreateAlarmRuleReq) error {
	doc := &model.AlarmRule{
		Enable:     true,
		ModuleName: in.ModuleName,
		Short:      in.Short,
		Level:      model.Level(in.Level),
		Tag:        in.Tag,
		RateSec:    in.RateSec,
		Method:     model.Method(in.Method),
		HookID:     in.HookID,
		UpdatedAt:  time.Now().Local(),
	}

	if err := srv.alarmRuleStore.InsertAlarmRule(ctx, doc); err != nil {
		return httputil.ErrSystemException.MergeError(err)
	}

	return nil
}

func (srv *Service) UpdateAlarmRule(ctx context.Context, in *entity.UpdateAlarmRuleReq) error {
	id, err := in.ObjectID()
	if err != nil {
		return err
	}

	doc := &model.AlarmRule{}
	if ok, err := srv.alarmRuleStore.FindOneAlarmRule(ctx, bson.M{"_id": id}, doc); err != nil {
		return httputil.ErrSystemException.MergeError(err)
	} else if !ok {
		return httputil.ErrNotFound
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

	if err := srv.alarmRuleStore.UpdateAlarmRule(ctx, filter, update); err != nil {
		return httputil.ErrSystemException.MergeError(err)
	}

	return nil
}

func (srv *Service) DeleteAlarmRule(ctx context.Context, in *entity.DeleteAlarmRuleReq) error {
	id, err := in.ObjectID()
	if err != nil {
		return err
	}

	if err := srv.alarmRuleStore.DeleteAlarmRule(ctx, bson.M{"_id": id}); err != nil {
		return httputil.ErrSystemException.MergeError(err)
	}
	return nil
}

// hook报警地址管理
func (srv *Service) FindHookURLList(ctx context.Context, in *entity.FindHookURLListReq, out *entity.ListResp) error {
	filter := bson.M{}
	if in.ID != "" {
		id, err := in.ObjectID()
		if err != nil {
			return err
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
	c, docs, err := srv.alarmRuleStore.FindCountHookURL(ctx, filter, opt)
	if err != nil {
		return httputil.ErrSystemException.MergeError(err)
	}

	out.Count = c
	list := make([]*entity.FindHookURLList, 0, len(docs))
	for _, v := range docs {
		d := &entity.FindHookURLList{
			ID:           v.ID.Hex(),
			Name:         v.Name,
			URL:          v.URL,
			Method:       v.Method.Int32(),
			KeyWord:      v.KeyWord,
			HideText:     v.HideText,
			UpdatedTsSec: v.UpdatedAt.Unix(),
		}
		if d.HideText == nil {
			d.HideText = make([]string, 0)
		}
		list = append(list, d)
	}
	out.List = list

	return nil
}

func (srv *Service) CreateHookURL(ctx context.Context, in *entity.CreateHookURLReq) error {
	doc := &model.HookURL{
		Name:      in.Name,
		URL:       in.URL,
		Method:    model.Method(in.Method),
		KeyWord:   in.KeyWord,
		HideText:  in.HideText,
		UpdatedAt: time.Now().Local(),
	}

	if err := srv.alarmRuleStore.InsertHookURL(ctx, doc); err != nil {
		return httputil.ErrSystemException.MergeError(err)
	}
	return nil
}

func (srv *Service) UpdateHookURL(ctx context.Context, in *entity.UpdateHookURLReq) error {
	id, err := in.ObjectID()
	if err != nil {
		return err
	}

	doc := &model.HookURL{}
	if ok, err := srv.alarmRuleStore.FindOneHookURL(ctx, bson.M{"_id": id}, doc); err != nil {
		return httputil.ErrSystemException.MergeError(err)
	} else if !ok {
		return httputil.ErrNotFound
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
	if err := srv.alarmRuleStore.UpdateManyAlarmRule(ctx, bson.M{"hook_id": id.Hex()},
		bson.M{"$set": bson.M{"updated_at": time.Now().Local()}}); err != nil {
		return httputil.ErrSystemException.MergeError(err)
	}

	if err := srv.alarmRuleStore.UpdateHookURL(ctx, filter, update); err != nil {
		return httputil.ErrSystemException.MergeError(err)
	}

	return nil
}

func (srv *Service) DelHookURL(ctx context.Context, in *entity.DelHookURLReq) error {
	id, err := in.ObjectID()
	if err != nil {
		return err
	}
	// 是否存在绑定
	opt := options.Find().SetLimit(1)
	c, _, err := srv.alarmRuleStore.FindCountAlarmRule(ctx, bson.M{"hook_id": in.ID}, opt)
	if err != nil {
		return httputil.ErrSystemException.MergeError(err)
	}
	if c != 0 {
		return httputil.ErrOpException.MergeString("hook URL cited")
	}

	if err := srv.alarmRuleStore.DelHookURL(ctx, bson.M{"_id": id}); err != nil {
		return httputil.ErrSystemException.MergeError(err)
	}
	return nil
}

func (srv *Service) PingHookURL(ctx context.Context, in *entity.PingHookURLReq) error {
	hook := &model.HookURL{}
	id, err := in.ObjectID()
	if err != nil {
		return err
	}
	ok, err := srv.alarmRuleStore.FindOneHookURL(ctx, bson.M{"_id": id}, hook)
	if err != nil {
		return httputil.ErrSystemException.MergeError(err)
	}
	if !ok {
		return httputil.ErrNotFound
	}

	switch hook.Method {
	case model.MethodDingDing:
		ding := alert.NewDingDing()
		ding.SetHookURL(hook.URL)
		if err := ding.Send(ctx, fmt.Sprintf("%s %s Ping hook Success", hook.KeyWord, hook.Name)); err != nil {
			return httputil.ErrOpException.MergeError(err)
		}
	}
	return nil
}
