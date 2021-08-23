package admin

import (
	"github.com/bbdshow/qelog/pkg/conf"
	"github.com/bbdshow/qelog/pkg/dao"
	"github.com/bbdshow/qelog/pkg/model"

	"sync"
)

type Service struct {
	cfg  *conf.Config
	d    *dao.Dao
	once sync.Once
}

func NewService(cfg *conf.Config) *Service {
	svc := &Service{
		cfg: cfg,
		d:   dao.New(cfg),
	}
	svc.once.Do(func() {
		go svc.bgDelExpiredCollection()
		go svc.bgMetricsCollectionStats()
		go svc.bgMetricsDBStats()
	})

	if err := svc.d.AdminInst().UpsertCollectionIndexMany(
		model.ModuleIndexMany(),
		model.AlarmRuleIndexMany(),
		model.DBStatsIndexMany(),
		model.ModuleMetricsIndexMany(),
		model.CollStatsIndexMany(),
	); err != nil {
		panic(err)
	}

	return svc
}

func (svc *Service) Close() {
	svc.d.Close()
}
