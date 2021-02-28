package manager

import (
	"context"
	"math/rand"
	"time"

	"github.com/huzhongqing/qelog/infra/logs"
	"go.uber.org/zap"

	"github.com/huzhongqing/qelog/pkg/config"
)

// backgroundDelExpiredCollection 删除已经过期了的集合
// 月为单位
func (srv *Service) backgroundDelExpiredCollection(maxAgeMonth int) {
	if maxAgeMonth <= 0 {
		// 永久保存
		return
	}
	for {
		time.Sleep(time.Duration(rand.Intn(5)+5) * time.Second)
		for i := 1; i <= config.Global.ShardingIndexSize; i++ {
			store, err := srv.sharding.GetStore(i)
			if err != nil {
				continue
			}
			// 找到所有 logging 开头的集合
			names, err := store.ListCollectionNames(context.Background(), "logging")
			if err != nil {
				logs.Qezap.Error("ListCollectionNames", zap.Error(err))
				continue
			}
			expiredNames := make([]string, 0)
			// 判断是否过期
			for _, v := range names {
				date, err := srv.lcn.NameDecodeDate(v)
				if err != nil {
					logs.Qezap.Error("NameDecodeDate", zap.Error(err))
					continue
				}
				y, m, _ := time.Now().Date()
				expiredTime := time.Date(y, m, 0, 0, 0, 0, 0, time.Local).
					AddDate(0, -maxAgeMonth, 0)

				if expiredTime.Equal(date) || expiredTime.After(date) {
					expiredNames = append(expiredNames, v)
				}
			}

			for _, v := range expiredNames {
				if err := store.Database().Collection(v).Drop(context.Background()); err != nil {
					logs.Qezap.Error("DropCollection", zap.Error(err))
					continue
				}
			}
		}

		time.Sleep(time.Duration(rand.Intn(30)+30) * time.Minute)
	}

}
