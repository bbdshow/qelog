package manager

import (
	"math/rand"
	"time"

	"github.com/huzhongqing/qelog/pkg/config"
)

// backgroundDelExpiredCollection 删除已经过期了的集合
// 月为单位
func (srv *Service) backgroundDelExpiredCollection(intervalMonth int) {
	if intervalMonth <= 0 {
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
			store.ListAllCollectionNames()

		}

		time.Sleep(time.Duration(rand.Intn(30)+30) * time.Minute)
	}

}
