package model

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

const (
	LogConfigCollectionName = "log_config"
	BucketCollectionName    = "bucket"
)

type LogConfig struct {
	ID        primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	Token     string             `bson:"token" json:"token"` // 访问权限
	Name      string             `bson:"name" json:"name"`   // 名称
	Desc      string             `bson:"desc" json:"desc"`
	UniqueKey string             `bson:"unique_key" json:"unique_key"` // 唯一身份，关联日志归属，
	Modules   []Module           `bson:"modules" json:"modules"`       // 日志下属模块
	UpdatedAt time.Time          `bson:"updated_at" json:"updated_at"`
	CreatedAt time.Time          `bson:"created_at" json:"created_at"`
}

type Module struct {
	Enable     bool   `bson:"enable" json:"enable"`
	Name       string `bson:"name" json:"name"`
	BucketName string `bson:"bucket_name" json:"bucket_name"`
}

func (lc LogConfig) CollectionName() string {
	return LogConfigCollectionName
}

type Bucket struct {
	ID        primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	Name      string             `bson:"name" json:"name"`
	Desc      string             `bson:"desc" json:"desc"`
	TTLMonth  int32              `bson:"ttl_month" json:"ttl_month"` // 这个仓库保存时间 /月
	CreatedAt time.Time          `bson:"created_at" json:"created_at"`
}

func (b Bucket) CollectionName() string {
	return BucketCollectionName
}
