package manager

import (
	"context"
	"time"

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

func (srv *Service) CreateModuleRegister(ctx context.Context, in *entity.CreateModuleRegisterReq) error {
	doc := &model.ModuleRegister{
		ModuleName:     in.ModuleName,
		Desc:           in.Desc,
		DBIndex:        in.DBIndex,
		HistoryDBIndex: make([]int32, 0),
		UpdatedAt:      time.Now().Local(),
	}
	return srv.store.InsertModuleRegister(ctx, doc)
}
