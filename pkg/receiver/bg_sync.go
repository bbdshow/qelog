package receiver

import (
	"context"
	"time"

	"github.com/bbdshow/bkit/db/mongo"
	"github.com/bbdshow/bkit/logs"
	"github.com/bbdshow/qelog/pkg/model"
	"go.uber.org/zap"
)

func (svc *Service) updateModuleSetting() error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	docs, err := svc.d.FindAllModule(ctx)
	if err != nil {
		return err
	}
	svc.lock.Lock()
	defer svc.lock.Unlock()

	for _, v := range docs {
		svc.modules[v.Name] = &module{
			m:  v,
			sc: mongo.NewShardCollection(v.Prefix, v.DaySpan),
		}
	}
	return nil
}
func (svc *Service) bgSyncModuleSetting() {
	tick := time.NewTicker(30 * time.Second)
	for range tick.C {
		err := svc.updateModuleSetting()
		if err != nil {
			logs.Qezap.Error("bgSyncModuleSetting", zap.Error(err))
		}
	}
}

func (svc *Service) updateAlarmRuleSetting() error {
	modules := make([]string, 0)
	svc.lock.RLock()
	for _, m := range svc.modules {
		modules = append(modules, m.m.Name)
	}
	svc.lock.RUnlock()

	rules := make([]*model.AlarmRule, 0)
	for _, name := range modules {
		docs, err := svc.d.FindAlarmRule(context.Background(), name, true)
		if err != nil {
			return err
		}
		rules = append(rules, docs...)
	}
	hooks, err := svc.d.FindAllHookURL(context.Background())
	if err != nil {
		return err
	}
	svc.alarm.InitRuleState(rules, hooks)
	return nil
}

func (svc *Service) bgSyncAlarmRuleSetting() {
	for {
		if err := svc.updateAlarmRuleSetting(); err != nil {
			logs.Qezap.Error("bgSyncAlarmRuleSetting", zap.String("updateAlarmRuleSetting", err.Error()))
		}
		time.Sleep(time.Minute)
	}
}
