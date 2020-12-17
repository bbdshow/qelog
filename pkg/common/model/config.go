package model

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type LogConfig struct {
	ID        primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	Token     string             `bson:"token" json:"token"`           // 访问权限
	Name      string             `bson:"name" json:"name"`             // 名称
	UniqueKey string             `bson:"unique_key" json:"unique_key"` // 唯一身份，关联日志归属，
	Modules   []Module           `bson:"modules" json:"modules"`       // 日志下属模块
	UpdatedAt time.Time          `bson:"updated_at" json:"updated_at"`
	CreatedAt time.Time          `bson:"created_at" json:"created_at"`
}

func (lc LogConfig) CollectionName() string {
	return "log_config"
}

type Module struct {
	Enable     bool   `bson:"enable" json:"enable"`
	Name       string `bson:"name" json:"name"`
	BucketName string `bson:"bucket_name" json:"bucket_name"`
}

type Bucket struct {
	ID          primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	Name        string             `bson:"name" json:"name"`
	ModuleNames []string           `bson:"module_names" json:"module_names"`
	CreatedAt   time.Time          `bson:"created_at" json:"created_at"`
}

func (b Bucket) CollectionName() string {
	return "bucket"
}
