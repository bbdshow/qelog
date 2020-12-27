package manager

import (
	"context"
	"fmt"
	"sort"
	"time"

	"go.mongodb.org/mongo-driver/bson"

	"github.com/huzhongqing/qelog/pkg/httputil"

	"github.com/huzhongqing/qelog/pkg/common/entity"
	"github.com/huzhongqing/qelog/pkg/common/model"
	"github.com/huzhongqing/qelog/pkg/storage"
)

type Service struct {
	store *storage.Store
}

func NewService(store *storage.Store) *Service {
	srv := &Service{store: store}
	return srv
}

func (srv *Service) FindModuleList(ctx context.Context, in *entity.FindModuleListReq, out *entity.ListResp) error {
	c, docs, err := srv.store.FindModuleList(ctx, in)
	if err != nil {
		return httputil.ErrSystemException.MergeError(err)
	}
	out.Count = c
	list := make([]*entity.FindModuleList, 0, len(docs))
	for _, v := range docs {
		d := &entity.FindModuleList{
			ID:             v.ID.Hex(),
			Name:           v.Name,
			Desc:           v.Desc,
			DBIndex:        v.DBIndex,
			HistoryDBIndex: v.HistoryDBIndex,
			UpdatedAt:      v.UpdatedAt.Unix(),
		}
		list = append(list, d)
	}
	out.List = list
	return nil
}

func (srv *Service) CreateModule(ctx context.Context, in *entity.CreateModuleReq) error {

	doc := &model.Module{
		Name:           in.Name,
		Desc:           in.Desc,
		DBIndex:        in.DBIndex,
		HistoryDBIndex: make([]int32, 0),
		UpdatedAt:      time.Now().Local(),
	}
	if err := srv.store.InsertModule(ctx, doc); err != nil {
		return httputil.ErrSystemException.MergeError(err)
	}
	return nil
}

func (srv *Service) UpdateModule(ctx context.Context, in *entity.UpdateModuleReq) error {
	id, err := in.ObjectID()
	if err != nil {
		return httputil.ErrArgsInvalid.MergeError(err)
	}

	doc := &model.Module{}
	if ok, err := srv.store.FindOneModule(ctx, bson.M{"_id": id}, doc); err != nil {
		return httputil.ErrSystemException.MergeError(err)
	} else if !ok {
		return httputil.ErrNotFound
	}
	update := bson.M{}
	fields := bson.M{}
	if doc.DBIndex != in.DBIndex {
		fields["db_index"] = in.DBIndex
		update["$addToSet"] = bson.M{"history_db_index": in.DBIndex}
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
	return srv.store.UpdateModule(ctx, filter, update)
}

func (srv *Service) DeleteModule(ctx context.Context, in *entity.DeleteModuleReq) error {
	id, err := in.ObjectID()
	if err != nil {
		return httputil.ErrArgsInvalid.MergeError(err)
	}
	doc := &model.Module{}
	if ok, err := srv.store.FindOneModule(ctx, bson.M{"_id": id}, doc); err != nil {
		return httputil.ErrSystemException.MergeError(err)
	} else if !ok {
		return httputil.ErrNotFound
	}
	if doc.Name != in.Name {
		return httputil.ErrNotFound
	}
	return srv.store.DeleteModule(ctx, id)
}

func (srv *Service) FindLoggingList(ctx context.Context, in *entity.FindLoggingListReq, out *entity.ListResp) error {
	// 查看改module当前的DB
	b, e := in.DefaultSection(time.Hour)
	collectionName := model.LoggingCollectionName(in.DBIndex, b.Unix())
	if collectionName != model.LoggingCollectionName(in.DBIndex, e.Unix()) {
		return httputil.ErrArgsInvalid.MergeError(fmt.Errorf("查询时间跨度不能超过时间分片设置 (分片粒度 %s)", model.LoggingShardingTime))
	}
	// 如果没有传入时间，则默认查询最近半小时
	in.BeginUnix = b.Unix()
	in.EndUnix = e.Unix()

	c, docs, err := srv.store.FindLoggingList(ctx, collectionName, in)
	if err != nil {
		return httputil.ErrSystemException.MergeError(err)
	}

	out.Count = c

	// 去除极低可能重复写入的日志信息
	hitMap := map[string]struct{}{}
	list := make([]*entity.FindLoggingList, 0, len(docs))
	for _, v := range docs {
		if _, ok := hitMap[v.MessageID]; ok {
			continue
		} else {
			hitMap[v.MessageID] = struct{}{}
		}

		d := &entity.FindLoggingList{
			ID:             v.ID.Hex(),
			TimeUnixMill:   v.Time,
			Level:          int(v.Level),
			ShortMsg:       v.Short,
			Full:           v.Full,
			ConditionOne:   v.Condition1,
			ConditionTwo:   v.Condition2,
			ConditionThree: v.Condition3,
			IP:             v.IP,
		}
		list = append(list, d)
	}
	out.List = list

	return nil
}

type AscDBIndexState []entity.DBIndexState

func (v AscDBIndexState) Len() int           { return len(v) }
func (v AscDBIndexState) Swap(i, j int)      { v[i], v[j] = v[j], v[i] }
func (v AscDBIndexState) Less(i, j int) bool { return v[i].Use < v[j].Use }

func (srv *Service) GetDBIndex(ctx context.Context, out *entity.GetDBIndexResp) error {
	docs, err := srv.store.FindAllModule(ctx)
	if err != nil {
		return httputil.ErrSystemException.MergeError(err)
	}
	state := make(map[int32]int)
	for i := int32(1); i <= model.MaxDBIndex; i++ {
		state[i] = 0
	}
	for _, v := range docs {
		num, ok := state[v.DBIndex]
		if ok {
			state[v.DBIndex] = num + 1
		}
	}

	states := make([]entity.DBIndexState, 0, len(state))
	for k, v := range state {
		states = append(states, entity.DBIndexState{
			Index: k,
			Use:   v,
		})
	}
	sort.Sort(AscDBIndexState(states))

	// 找到最小的，作为推荐
	suggestDBIndex := model.MaxDBIndex
	if len(states) > 0 {
		suggestDBIndex = states[0].Index
	}

	out.SuggestDBIndex = suggestDBIndex
	out.MaxDBIndex = model.MaxDBIndex
	out.UseState = states

	return nil
}
