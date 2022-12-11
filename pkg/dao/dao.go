package dao

import (
	"context"
	"fmt"
	"math"
	"time"

	"github.com/bbdshow/bkit/db/mongo"
	"github.com/bbdshow/qelog/pkg/conf"
)

// Dao all database operation
type Dao struct {
	cfg *conf.Config
	//custom mongo shard
	mongo *mongo.Groups
	// admin mongo inst
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

// ReceiverInst query receiver server db inst by name
func (d *Dao) ReceiverInst(dbName string) (*mongo.Database, error) {
	db, err := d.mongo.GetInstance(dbName)
	if err != nil {
		return nil, fmt.Errorf("%s %v", dbName, err)
	}
	return db, nil
}

// CtxAfterSecDeadline if not deadline, return defSec, if defSec <= 0, return int32 max sec duration
func (d *Dao) CtxAfterSecDeadline(ctx context.Context, defSec int32) time.Duration {
	deadline, ok := ctx.Deadline()
	if !ok {
		if defSec <= 0 {
			defSec = math.MaxInt32
		}
		return time.Duration(defSec) * time.Second
	}
	sec := int32(deadline.Sub(time.Now()).Seconds())
	if sec <= 0 {
		sec = defSec
	}
	return time.Duration(sec) * time.Second
}
