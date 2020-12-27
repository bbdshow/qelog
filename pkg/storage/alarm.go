package storage

import (
	"context"

	"github.com/huzhongqing/qelog/pkg/common/model"
	"go.mongodb.org/mongo-driver/bson"
)

// 当报警规则超过上千的规则，可优化语句
func (store *Store) FindAllAlarmRule(ctx context.Context) ([]*model.AlarmRule, error) {
	docs := make([]*model.AlarmRule, 0)
	coll := store.database.Collection(model.CollectionNameAlarmRule)
	err := store.database.Find(ctx, coll, bson.M{}, &docs)
	return docs, err
}
