package manager

import (
	"context"
	"fmt"
	"time"

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

func (srv *Service) CreateModule(ctx context.Context, in *entity.CreateModuleReq) error {
	doc := &model.Module{
		Name:           in.Name,
		Desc:           in.Desc,
		DBIndex:        in.DBIndex,
		HistoryDBIndex: make([]int32, 0),
		UpdatedAt:      time.Now().Local(),
	}
	return srv.store.InsertModule(ctx, doc)
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

	c, records, err := srv.store.FindLoggingList(ctx, collectionName, in)
	if err != nil {
		return httputil.ErrSystemException.MergeError(err)
	}

	out.Count = c

	// 去除极低可能重复写入的日志信息
	hitMap := map[string]struct{}{}
	list := make([]*entity.FindLoggingList, 0, len(records))
	for _, v := range records {
		if _, ok := hitMap[v.MessageID]; ok {
			continue
		} else {
			hitMap[v.MessageID] = struct{}{}
		}

		d := &entity.FindLoggingList{
			ID:             v.ID.Hex(),
			TimeUnixMill:   v.Time,
			Level:          v.Level,
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
