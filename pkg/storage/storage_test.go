package storage

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/huzhongqing/qelog/libs/mongo"
	"github.com/huzhongqing/qelog/pkg/config"
)

var _database *mongo.Database

func init() {
	cfg := config.MockDevConfig()
	db, err := mongo.NewDatabase(context.Background(), cfg.MongoDB.URI, cfg.MongoDB.DataBase)
	if err != nil {
		panic(err)
	}
	_database = db
}

func TestStore_MetricsModuleCountByDate(t *testing.T) {
	store := New(_database)
	date, err := time.ParseInLocation("2006-01-02", "2021-01-04", time.Local)
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println(date.String())
	mc, err := store.MetricsModuleCountByDate(context.Background(), date)
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println(mc)
}
