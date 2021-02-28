package metrics

import (
	"testing"
	"time"

	"github.com/huzhongqing/qelog/pkg/storage"

	"github.com/huzhongqing/qelog/tests"

	"github.com/huzhongqing/qelog/pkg/common/model"
)

func init() {
	tests.InitTestDepends()
}

func TestMetricsStatistics(t *testing.T) {

	store := storage.ShardingDB
	mainStore, _ := store.MainStore()
	SetIncIntervalSec(3)

	m := NewMetrics(mainStore)
	docs := []*model.Logging{
		{
			Level: 0,
		},
		{
			Level: 1,
		},
		{
			Level: 2,
		},
	}

	m.Statistics("example", "127.0.0.1", docs)
	time.Sleep(5 * time.Second)
	m.Statistics("example", "127.0.0.2", docs)
	m.Statistics("example", "0:0:0:0:0:ffff:0:0", docs)

	m.Sync()
	time.Sleep(time.Second)
}
