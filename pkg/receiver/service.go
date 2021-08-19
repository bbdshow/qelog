package receiver

import (
	"github.com/bbdshow/bkit/db/mongo"
	"github.com/bbdshow/qelog/pkg/conf"
	"github.com/bbdshow/qelog/pkg/dao"
	"github.com/bbdshow/qelog/pkg/model"
	"github.com/bbdshow/qelog/pkg/receiver/alarm"
	"github.com/bbdshow/qelog/pkg/receiver/metrics"

	"sync"
)

type Service struct {
	cfg *conf.Config
	d   *dao.Dao

	lock        sync.RWMutex
	modules     map[string]*module
	collections map[string]struct{}

	alarm   *alarm.Alarm
	metrics *metrics.Metrics
}

type module struct {
	m  *model.Module
	sc mongo.ShardCollection
}

func NewService(cfg *conf.Config) *Service {
	svc := &Service{
		d:           dao.New(cfg),
		lock:        sync.RWMutex{},
		modules:     map[string]*module{},
		collections: map[string]struct{}{},
	}
	svc.metrics = metrics.NewMetrics(svc.d)

	if err := svc.updateModuleSetting(); err != nil {
		panic(err)
	}

	go svc.bgSyncModuleSetting()

	if cfg.Receiver.AlarmEnable {
		svc.alarm = alarm.NewAlarm()
		go svc.bgSyncAlarmRuleSetting()
	}

	if cfg.Receiver.MetricsEnable {
		svc.metrics = metrics.NewMetrics(svc.d)
		metrics.SetIncIntervalSec(30)
	}

	return svc
}

func (svc *Service) Close() {
	if svc.d != nil {
		svc.d.Close()
	}

	if svc.metrics != nil {
		svc.metrics.Sync()
	}
}
