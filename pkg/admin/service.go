package admin

import (
	"context"
	"go.mongodb.org/mongo-driver/bson"
	"sync"
	"time"

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

	if err := svc.initData(); err != nil {
		panic(err)
	}

	bgOp := func() {
		go svc.bgDelExpiredCollection()
		go svc.bgMetricsCollectionStats()
		go svc.bgMetricsDBStats()
	}
	svc.once.Do(bgOp)

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

func (svc *Service) initData() error {
	if svc.cfg.Logging != nil {
		if svc.cfg.Logging.Module != "" {
			ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
			defer cancel()
			exists, _, err := svc.d.GetModule(ctx, bson.M{"name": svc.cfg.Logging.Module})
			if err != nil {
				return err
			}
			if !exists {
				if err := svc.CreateModule(ctx, &model.CreateModuleReq{
					Name:     svc.cfg.Logging.Module,
					Desc:     "self-access",
					DaySpan:  0,
					MaxMonth: 6,
				}); err != nil {
					return err
				}
			}

		}
	}

	return nil
}

func (svc *Service) Close() {
	if svc.d != nil {
		svc.d.Close()
	}
}
