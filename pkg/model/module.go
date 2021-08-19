package model

import (
	"github.com/bbdshow/bkit/db/mongo"
	"time"

	"go.mongodb.org/mongo-driver/bson"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

const (
	CNModule = "module"
)

// Module 接入应用模块初始化
type Module struct {
	ID        primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	Name      string             `bson:"name" json:"name"`
	Desc      string             `bson:"desc" json:"desc"`
	Bucket    string             `bson:"bucket" json:"bucket"`
	Database  string             `bson:"database" json:"database"`
	DaySpan   int                `bson:"day_span" json:"day_span"`
	MaxMonth  int                `bson:"max_month" json:"max_month"`
	Prefix    string             `bson:"prefix" json:"prefix"`
	UpdatedAt time.Time          `bson:"updated_at" json:"updated_at"`
}

func (m Module) CollectionName() string {
	return CNModule
}
func ModuleIndexMany() []mongo.Index {
	return []mongo.Index{
		{
			Collection: CNModule,
			Keys: bson.D{{
				Key: "name", Value: 1,
			}},
			Unique:     true,
			Background: true,
		},
		{
			Collection: CNModule,
			Keys: bson.D{{
				Key: "bucket", Value: 1,
			}},
			Unique:     true,
			Background: true,
		},
	}
}
