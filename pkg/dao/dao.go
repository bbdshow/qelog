package dao

import (
	"fmt"
	"github.com/bbdshow/bkit/db/mongo"
	"github.com/bbdshow/qelog/pkg/conf"
)

type Dao struct {
	cfg   *conf.Config
	mongo *mongo.Groups

	adminInst *mongo.Database
}

func New(cfg *conf.Config) *Dao {
	d := &Dao{
		cfg: cfg,
	}
	mc, err := mongo.NewGroups(cfg.Mongo)
	if err != nil {
		panic(err)
	}
	d.mongo = mc

	d.adminInst = d.AdminInst()

	return d
}

func (d *Dao) Close() {
	if d.mongo != nil {
		_ = d.mongo.Disconnect()
	}
}

func (d *Dao) AdminInst() *mongo.Database {
	dbName := d.cfg.MongoGroup.AdminDatabase
	db, err := d.mongo.GetInstance(dbName)
	if err != nil {
		panic(fmt.Sprintf("%s %v", dbName, err))
	}
	return db
}

func (d *Dao) ReceiverInst(dbName string) (*mongo.Database, error) {
	db, err := d.mongo.GetInstance(dbName)
	if err != nil {
		return nil, fmt.Errorf("%s %v", dbName, err)
	}
	return db, nil
}
