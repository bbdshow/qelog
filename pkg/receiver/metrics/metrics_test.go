package metrics

import (
	"context"
	"testing"
	"time"

	"github.com/huzhongqing/qelog/pkg/common/model"

	"github.com/huzhongqing/qelog/pkg/storage"

	"github.com/huzhongqing/qelog/infra/mongo"
	"github.com/huzhongqing/qelog/pkg/config"
)

func TestMetricsStatistics(t *testing.T) {
	cfg := config.MockDevConfig()
	database, err := mongo.NewDatabase(context.Background(), cfg.MongoDB.URI, cfg.MongoDB.DataBase)
	if err != nil {
		t.Fatal(err)
	}
	store := storage.New(database)

	SetIncIntervalSec(3)

	m := NewMetrics(store)
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
