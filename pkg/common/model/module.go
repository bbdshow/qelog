package model

import (
	"time"

	"go.mongodb.org/mongo-driver/bson"

	"github.com/huzhongqing/qelog/libs/mongo"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

const (
	CollectionNameModuleRegister = "module_register"
)

type ModuleRegister struct {
	ID             primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	ModuleName     string             `bson:"module_name" json:"module_name"`
	Desc           string             `bson:"desc" json:"desc"`
	DBIndex        int32              `bson:"db_index" json:"db_index"`
	HistoryDBIndex []int32            `bson:"history_db_index" json:"history_db_index"`
	UpdatedAt      time.Time          `bson:"updated_at" json:"updated_at"`
}

func (m ModuleRegister) CollectionName() string {
	return CollectionNameModuleRegister
}

func ModuleRegisterIndexMany() []mongo.Index {
	return []mongo.Index{
		{
			Collection: CollectionNameModuleRegister,
			Keys:       bson.M{"module_name": 1},
			Unique:     true,
			Background: true,
		},
	}
}
