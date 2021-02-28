package model

import (
	"time"

	"go.mongodb.org/mongo-driver/bson"

	"github.com/huzhongqing/qelog/infra/mongo"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

const (
	CollectionNameModule = "module"
)

// Module 接入应用模块初始化
type Module struct {
	ID                   primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	Name                 string             `bson:"name" json:"name"`
	Desc                 string             `bson:"desc" json:"desc"`
	ShardingIndex        int                `bson:"sharding_index"`
	HistoryShardingIndex []int              `bson:"history_sharding_index"`
	UpdatedAt            time.Time          `bson:"updated_at" json:"updated_at"`
}

func (m Module) CollectionName() string {
	return CollectionNameModule
}

func ModuleIndexMany() []mongo.Index {
	return []mongo.Index{
		{
			Collection: CollectionNameModule,
			Keys: bson.D{{
				Key: "name", Value: 1,
			}},
			Unique:     true,
			Background: true,
		},
	}
}
