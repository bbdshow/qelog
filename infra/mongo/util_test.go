package mongo

import (
	"context"
	"fmt"
	"testing"
)

func TestUtil_CollectionsStats(t *testing.T) {
	database, err := NewDatabase(context.Background(), _uri, _database)
	if err != nil {
		t.Fatal(err)
	}
	util := NewUtil(database)

	stats, err := util.CollStats(nil, []string{"module"})
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println(stats)
}

func TestUtil_HostInfo(t *testing.T) {
	database, err := NewDatabase(context.Background(), _uri, _database)
	if err != nil {
		t.Fatal(err)
	}
	util := NewUtil(database)

	info, err := util.HostInfo(nil)
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println(info)
}

func TestUtil_DBStats(t *testing.T) {
	database, err := NewDatabase(context.Background(), _uri, _database)
	if err != nil {
		t.Fatal(err)
	}
	util := NewUtil(database)

	info, err := util.DBStats(nil)
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println(info)
}
