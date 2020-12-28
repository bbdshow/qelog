package model

import (
	"fmt"
	"time"

	"go.mongodb.org/mongo-driver/bson"

	"github.com/huzhongqing/qelog/libs/mongo"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

const (
	CollectionNameAlarmRule = "alarm_rule"
)

const (
	MethodDingDing = iota + 1
)

type AlarmRule struct {
	ID         primitive.ObjectID `bson:"_id,omitempty"`
	Enable     bool               `bson:"enable"`       // 是否开启
	ModuleName string             `bson:"module_name" ` // 哪个模块
	Short      string             `bson:"short"`        // 命中的短消息
	Level      Level              `bson:"level"`        // 命中日志等级
	Tag        string             `bson:"tag"`          // 报警Tag
	RateSec    int64              `bson:"rate_sec"`     // 多少s之内，只发送一次
	Method     int32              `bson:"method"`       // 支持方式  1-钉钉
	HookURL    string             `bson:"hook_url"`     // 发送链接
	UpdatedAt  time.Time          `bson:"updated_at"`
}

func (AlarmRule) CollectionName() string {
	return CollectionNameAlarmRule
}

func (am AlarmRule) Key() string {
	return fmt.Sprintf("%s_%s_%d", am.ModuleName, am.Short, am.Level)
}

func AlarmRuleIndexMany() []mongo.Index {
	return []mongo.Index{{
		Collection: CollectionNameAlarmRule,
		Keys: bson.M{
			"module_name": 1,
			"short":       1,
			"level":       1,
		},
		Unique:     true,
		Background: true,
	}}
}
