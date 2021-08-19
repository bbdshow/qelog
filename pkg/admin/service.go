package admin

import (
	"github.com/bbdshow/qelog/pkg/conf"
	"github.com/bbdshow/qelog/pkg/dao"

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
	})
	return svc
}

func (svc *Service) Close() {
	svc.d.Close()
}
