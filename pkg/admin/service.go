package admin

import (
	"sync"

	"github.com/bbdshow/qelog/pkg/conf"
	"github.com/bbdshow/qelog/pkg/dao"
	"github.com/bbdshow/qelog/pkg/model"
)

// Service admin
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

	// admin db inst, create collection and index
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
	if svc.d != nil {
		svc.d.Close()
	}
}
