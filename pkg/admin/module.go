package admin

import (
	"context"
	"fmt"
	"github.com/bbdshow/bkit/errc"
	"github.com/bbdshow/bkit/gen/str"
	"github.com/bbdshow/qelog/pkg/model"
	"go.mongodb.org/mongo-driver/bson"
	"time"
)

func (svc *Service) FindModuleList(ctx context.Context, in *model.FindModuleListReq, out *model.ListResp) error {
	c, docs, err := svc.d.FindModuleList(ctx, in)
	if err != nil {
		return errc.ErrInternalErr.MultiErr(err)
	}

	out.Count = c
	list := make([]*model.FindModuleList, 0, len(docs))
	for _, v := range docs {
		d := &model.FindModuleList{
			ID:           v.ID.Hex(),
			Name:         v.Name,
			Desc:         v.Desc,
			Bucket:       v.Bucket,
			DaySpan:      v.DaySpan,
			MaxMonth:     v.MaxMonth,
			UpdatedTsSec: v.UpdatedAt.Unix(),
		}
		list = append(list, d)
	}
	out.List = list
	return nil
}

func (svc *Service) CreateModule(ctx context.Context, in *model.CreateModuleReq) error {
	doc := &model.Module{
		Name:      in.Name,
		Desc:      in.Desc,
		Bucket:    str.RandAlphaNumString(6, true),
		DaySpan:   in.DaySpan,
		MaxMonth:  in.MaxMonth,
		Database:  svc.cfg.MongoGroup.RandReceiverDatabase(),
		UpdatedAt: time.Now(),
	}
	if err := svc.d.CreateModule(ctx, doc); err != nil {
		return errc.ErrInternalErr.MultiErr(err)
	}
	return nil
}

func (svc *Service) UpdateModule(ctx context.Context, in *model.UpdateModuleReq) error {
	if in.Database != "" {
		if !svc.cfg.MongoGroup.IsExists(in.Database) {
			return errc.ErrNotFound.MultiMsg(fmt.Sprintf("%s database", in.Database))
		}
	}
	if err := svc.d.UpdateModule(ctx, in); err != nil {
		return errc.ErrInternalErr.MultiErr(err)
	}
	return nil
}

func (svc *Service) DelModule(ctx context.Context, in *model.DelModuleReq) error {
	id, err := in.ObjectID()
	if err != nil {
		return err
	}
	exists, doc, err := svc.d.GetModule(ctx, bson.M{"_id": id})
	if err != nil {
		return errc.ErrInternalErr.MultiErr(err)
	}
	if !exists {
		return nil
	}
	if doc.Name != in.Name {
		return errc.ErrNotFound
	}
	if err := svc.d.DelModule(ctx, bson.M{"_id": id}); err != nil {
		return errc.ErrInternalErr.MultiErr(err)
	}
	return nil
}
