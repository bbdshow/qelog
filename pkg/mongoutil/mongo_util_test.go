package mongoutil

import (
	"context"
	"fmt"
	"testing"

	"github.com/huzhongqing/qelog/libs/mongo"
)

var (
	database = "qelog"
	uri      = "mongodb://127.0.0.1:27017/admin"
)

func TestMongodbUtil_CollectionsStats(t *testing.T) {
	database, err := mongo.NewDatabase(context.Background(), uri, database)
	if err != nil {
		t.Fatal(err)
	}
	mu := NewMongodbUtil(database)

	stats, err := mu.CollStats(nil, []string{"module"})
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println(stats)
}

func TestMongodbUtil_HostInfo(t *testing.T) {
	database, err := mongo.NewDatabase(context.Background(), uri, database)
	if err != nil {
		t.Fatal(err)
	}
	mu := NewMongodbUtil(database)

	info, err := mu.HostInfo(nil)
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println(info)
}

func TestMongodbUtil_DBStats(t *testing.T) {
	database, err := mongo.NewDatabase(context.Background(), uri, database)
	if err != nil {
		t.Fatal(err)
	}
	mu := NewMongodbUtil(database)

	info, err := mu.DBStats(nil)
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println(info)
}
