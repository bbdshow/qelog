package model

import (
	"fmt"
	"time"

	"go.mongodb.org/mongo-driver/bson"

	"github.com/huzhongqing/qelog/infra/mongo"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

const (
	CollectionNameAlarmRule = "alarm_rule"
	CollectionNameHookURL   = "hook_url"
)

const (
	MethodDingDing = iota + 1
)

// AlarmRule
type AlarmRule struct {
	ID         primitive.ObjectID `bson:"_id,omitempty"`
	Enable     bool               `bson:"enable"`       // 是否开启
	ModuleName string             `bson:"module_name" ` // 哪个模块
	Short      string             `bson:"short"`        // 命中的短消息
	Level      Level              `bson:"level"`        // 命中日志等级
	Tag        string             `bson:"tag"`          // 报警Tag
	RateSec    int64              `bson:"rate_sec"`     // 多少s之内，只发送一次
	Method     Method             `bson:"method"`       // 支持方式  1-钉钉
	HookID     string             `bson:"hook_id"`
	UpdatedAt  time.Time          `bson:"updated_at"`
}

func (AlarmRule) CollectionName() string {
	return CollectionNameAlarmRule
}

func (ar AlarmRule) Key() string {
	return fmt.Sprintf("%s_%s_%s", ar.ModuleName, ar.Short, ar.Level)
}

type HookURL struct {
	ID        primitive.ObjectID `bson:"_id,omitempty"`
	Name      string             `bson:"name"`
	URL       string             `bson:"url"`
	Method    Method             `bson:"method"`
	KeyWord   string             `bson:"key_word"`
	HideText  []string           `bson:"hide_text"`
	UpdatedAt time.Time          `bson:"updated_at"`
}

func (HookURL) CollectionName() string {
	return CollectionNameHookURL
}

type Method int32

func (m Method) Int32() int32 {
	return int32(m)
}
func (m Method) String() string {
	v := "UNKNOWN"
	switch m {
	case MethodDingDing:
		v = "DingDing"
	}
	return v
}

func AlarmRuleIndexMany() []mongo.Index {
	return []mongo.Index{{
		Collection: CollectionNameAlarmRule,
		Keys: bson.D{
			{
				Key: "module_name", Value: 1,
			},
			{
				Key: "short", Value: 1,
			},
			{
				Key: "level", Value: 1,
			},
		},
		Unique:     true,
		Background: true,
	}, {
		Collection: CollectionNameHookURL,
		Keys: bson.D{
			{
				Key: "name", Value: 1,
			},
		},
		Unique:     true,
		Background: true,
	}}
}
