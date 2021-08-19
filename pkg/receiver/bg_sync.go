package receiver

import (
	"context"
	"github.com/bbdshow/bkit/db/mongo"
	"github.com/bbdshow/bkit/logs"
	"go.uber.org/zap"
	"time"
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

func (svc *Service) updateAlarmRuleSetting(module string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	docs, err := svc.d.FindAlarmRule(ctx, module, true)
	if err != nil {
		return err
	}
	hooks, err := svc.d.FindAllHookURL(ctx)
	if err != nil {
		return err
	}
	svc.alarm.InitRuleState(docs, hooks)
	return nil
}

func (svc *Service) bgSyncAlarmRuleSetting() {
	for {
		modules := make([]string, 0)
		svc.lock.RLock()
		for _, m := range svc.modules {
			modules = append(modules, m.m.Name)
		}
		svc.lock.RUnlock()
		for _, name := range modules {
			err := svc.updateAlarmRuleSetting(name)
			if err != nil {
				logs.Qezap.Error("bgSyncAlarmRuleSetting", zap.String("error", err.Error()))
			}
		}
		time.Sleep(time.Minute)
	}
}
