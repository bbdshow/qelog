package dao

import (
	"context"
	"github.com/bbdshow/bkit/errc"
	"github.com/bbdshow/qelog/pkg/model"
)

func (d *Dao) CreateManyLogger(ctx context.Context, dbName, cName string, docs []interface{}) error {
	inst, err := d.mongo.GetInstance(dbName)
	if err != nil {
		return errc.WithStack(err)
	}
	_, err = inst.Collection(cName).InsertMany(ctx, docs)
	return errc.WithStack(err)
}

func (d *Dao) ListCollectionNames(ctx context.Context, dbName string, prefix ...string) ([]string, error) {
	inst, err := d.mongo.GetInstance(dbName)
	if err != nil {
		return nil, err
	}
	return inst.ListCollectionNames(ctx, prefix...)
}

func (d *Dao) CreateLoggerIndex(dbName, cName string) error {
	inst, err := d.mongo.GetInstance(dbName)
	if err != nil {
		return err
	}
	return inst.UpsertCollectionIndexMany(model.LoggingIndexMany(cName))
}
