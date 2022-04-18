package admin

import (
	"context"
	"fmt"
	"github.com/bbdshow/bkit/errc"
	"github.com/bbdshow/bkit/util/alert"
	"github.com/bbdshow/qelog/common/types"
	"github.com/bbdshow/qelog/pkg/model"

	"go.mongodb.org/mongo-driver/bson"
	"time"
)

func (svc *Service) FindAlarmRuleList(ctx context.Context, in *model.FindAlarmRuleListReq, out *model.ListResp) error {

	c, docs, err := svc.d.FindAlarmRuleList(ctx, in)
	if err != nil {
		return errc.ErrInternalErr.MultiErr(err)
	}

	out.Count = c
	list := make([]*model.FindAlarmRuleList, 0, len(docs))
	for _, v := range docs {
		d := &model.FindAlarmRuleList{
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

func (svc *Service) CreateAlarmRule(ctx context.Context, in *model.CreateAlarmRuleReq) error {
	doc := &model.AlarmRule{
		Enable:     true,
		ModuleName: in.ModuleName,
		Short:      in.Short,
		Level:      types.Level(in.Level),
		Tag:        in.Tag,
		RateSec:    in.RateSec,
		Method:     model.Method(in.Method),
		HookID:     in.HookID,
		UpdatedAt:  time.Now().Local(),
	}

	if err := svc.d.CreateAlarmRule(ctx, doc); err != nil {
		return errc.ErrInternalErr.MultiErr(err)
	}

	return nil
}

func (svc *Service) UpdateAlarmRule(ctx context.Context, in *model.UpdateAlarmRuleReq) error {
	if err := svc.d.UpdateAlarmRule(ctx, in); err != nil {
		return errc.ErrInternalErr.MultiErr(err)
	}
	return nil
}

func (svc *Service) DelAlarmRule(ctx context.Context, in *model.DelAlarmRuleReq) error {
	id, err := in.ObjectID()
	if err != nil {
		return err
	}
	if err := svc.d.DelAlarmRule(ctx, bson.M{"_id": id}); err != nil {
		return errc.ErrInternalErr.MultiErr(err)
	}
	return nil
}

// hook报警地址管理
func (svc *Service) FindHookURLList(ctx context.Context, in *model.FindHookURLListReq, out *model.ListResp) error {
	c, docs, err := svc.d.FindHookURLList(ctx, in)
	if err != nil {
		return errc.ErrInternalErr.MultiErr(err)
	}
	out.Count = c
	list := make([]*model.FindHookURLList, 0, len(docs))
	for _, v := range docs {
		d := &model.FindHookURLList{
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

func (svc *Service) CreateHookURL(ctx context.Context, in *model.CreateHookURLReq) error {
	doc := &model.HookURL{
		Name:      in.Name,
		URL:       in.URL,
		Method:    model.Method(in.Method),
		KeyWord:   in.KeyWord,
		HideText:  in.HideText,
		UpdatedAt: time.Now().Local(),
	}

	if err := svc.d.CreateHookURL(ctx, doc); err != nil {
		return errc.ErrInternalErr.MultiErr(err)
	}
	return nil
}

func (svc *Service) UpdateHookURL(ctx context.Context, in *model.UpdateHookURLReq) error {
	if err := svc.d.UpdateHookURL(ctx, in); err != nil {
		return errc.ErrInternalErr.MultiErr(err)
	}
	return nil
}

func (svc *Service) DelHookURL(ctx context.Context, in *model.DelHookURLReq) error {
	if err := svc.d.DelHookURL(ctx, in); err != nil {
		return errc.ErrInternalErr.MultiErr(err)
	}
	return nil
}

func (svc *Service) PingHookURL(ctx context.Context, in *model.PingHookURLReq) error {
	id, err := in.ObjectID()
	if err != nil {
		return err
	}
	exists, doc, err := svc.d.GetHookURL(ctx, bson.M{"_id": id})
	if err != nil {
		return errc.ErrInternalErr.MultiErr(err)
	}
	if !exists {
		return errc.ErrNotFound
	}
	switch doc.Method {
	case model.MethodDingDing:
		ding := alert.NewDingDing()
		ding.SetHookURL(doc.URL)
		if err := ding.Send(ctx, fmt.Sprintf("%s %s Ping hook Success", doc.KeyWord, doc.Name)); err != nil {
			return errc.ErrParamInvalid.MultiErr(err)
		}
	case model.MethodTelegram:
		tel := alert.NewTelegram()
		tel.SetHookURL(doc.URL)
		if err := tel.Send(ctx, fmt.Sprintf("%s %s Ping hook Success", doc.KeyWord, doc.Name)); err != nil {
			return errc.ErrParamInvalid.MultiErr(err)
		}
	}
	return nil
}
