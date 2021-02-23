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

type Module struct {
	ID             primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	Name           string             `bson:"name" json:"name"`
	Desc           string             `bson:"desc" json:"desc"`
	DBIndex        int32              `bson:"db_index" json:"db_index"`
	HistoryDBIndex []int32            `bson:"history_db_index" json:"history_db_index"`
	UpdatedAt      time.Time          `bson:"updated_at" json:"updated_at"`
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
