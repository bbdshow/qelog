package entity

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type PageReq struct {
	Page     int64 `json:"page" form:"page" default:"1" `
	PageSize int64 `json:"page_size" from:"page_size" default:"20"`
}

func (v PageReq) SetPage(opt *options.FindOptions) {
	opt.SetSkip((v.Page - 1) * v.PageSize)
	opt.SetLimit(v.PageSize)
}

type ObjectIDReq struct {
	ID string `json:"id" form:"id" binding:"required,eq=24"`
}

func (v ObjectIDReq) ObjectID() (primitive.ObjectID, error) {
	return primitive.ObjectIDFromHex(v.ID)
}

type TimeReq struct {
	BeginTsSec int64 `json:"begin_ts_sec" form:"begin_ts_sec"`
	EndTsSec   int64 `json:"end_ts_sec" form:"end_ts_sec"`
}

func (v TimeReq) BeginTime() time.Time {
	return time.Unix(v.BeginTsSec, 0)
}
func (v TimeReq) EndTime() time.Time {
	return time.Unix(v.EndTsSec, 0)
}

func (v TimeReq) InitTimeSection(t time.Duration) (b, e time.Time) {
	// 没有查询时间
	if v.BeginTsSec <= 0 && v.EndTsSec <= 0 {
		e = time.Now()
		b = e.Add(-t)
		return b, e
	}
	// 只有结束时间
	if v.BeginTsSec <= 0 && v.EndTsSec > 0 {
		e = time.Unix(v.EndTsSec, 0)
		b = e.Add(-t)
		return b, e
	}
	// 只有开始时间
	if v.BeginTsSec > 0 && v.EndTsSec <= 0 {
		b = time.Unix(v.BeginTsSec, 0)
		e = time.Now()
		return b, e
	}

	// 时间都存在
	return time.Unix(v.BeginTsSec, 0), time.Unix(v.EndTsSec, 0)
}

type ListResp struct {
	Count int64       `json:"count"`
	List  interface{} `json:"list"`
}

type CreateModuleReq struct {
	Name    string `json:"name" binding:"required,gte=2,lte=24,lowercase"`
	DBIndex int32  `json:"db_index" binding:"required,min=1,max=16"`
	Desc    string `json:"desc" binding:"required,gte=1,lte=128"`
}

type FindModuleListReq struct {
	Name string `json:"name"`
	PageReq
}

type FindModuleList struct {
	ID             string  `json:"id"`
	Name           string  `json:"name"`
	Desc           string  `json:"desc"`
	DBIndex        int32   `json:"db_index"`
	HistoryDBIndex []int32 `json:"history_db_index"`
	UpdatedTsSec   int64   `json:"updated_ts_sec"`
}

type UpdateModuleReq struct {
	ObjectIDReq
	DBIndex int32  `json:"db_index" binding:"required,min=1,max=16"`
	Desc    string `json:"desc" binding:"required,gte=1,lte=128"`
}

type DeleteModuleReq struct {
	ObjectIDReq
	Name string `json:"name" binding:"required"`
}

type GetDBIndexResp struct {
	SuggestDBIndex int32          `json:"suggest_db_index"`
	MaxDBIndex     int32          `json:"max_db_index"`
	UseState       []DBIndexState `json:"use_state"`
}
type DBIndexState struct {
	Index int32 `json:"index"`
	Use   int   `json:"use"`
}

type FindLoggingListReq struct {
	DBIndex        int32  `json:"db_index" binding:"required"`
	ModuleName     string `json:"module_name" binding:"required"`
	Short          string `json:"short"`
	Level          int32  `json:"level" binding:"omitempty,min=-1,max=8"`
	IP             string `json:"ip"`
	ConditionOne   string `json:"condition_one"`
	ConditionTwo   string `json:"condition_two"`
	ConditionThree string `json:"condition_three"`
	TimeReq
	PageReq
}

type FindLoggingList struct {
	ID             string `json:"id"`
	TsMill         int64  `json:"ts_mill"`
	Level          int32  `json:"level"`
	ShortMsg       string `json:"short_msg"`
	Full           string `json:"full"`
	ConditionOne   string `json:"condition_one"`
	ConditionTwo   string `json:"condition_two"`
	ConditionThree string `json:"condition_three"`
	IP             string `json:"ip"`
}

type FindAlarmRuleListReq struct {
	Enable     int    `json:"enable" binding:"omitempty,min=-1,max=1"`
	ModuleName string `json:"module_name"`
	Short      string `json:"short"`
	PageReq
}

type FindAlarmRuleList struct {
	ID           string `json:"id"`
	Enable       bool   `json:"enable"`
	ModuleName   string `json:"module_name"`
	Short        string `json:"short"`
	Level        int32  `json:"level"`
	Tag          string `json:"tag"`
	RateSec      int64  `json:"rate_sec"`
	Method       int32  `json:"method"`
	HookURL      string `json:"hook_url"`
	UpdatedTsSec int64  `json:"updated_ts_sec"`
}

type CreateAlarmRuleReq struct {
	ModuleName string `json:"module_name" binding:"required"`
	Short      string `json:"short" binding:"required"`
	Level      int32  `json:"level" binding:"required,min=0,max=8"`
	Tag        string `json:"tag" binding:"omitempty,gte=1,lte=128"`
	RateSec    int64  `json:"rate_sec" binding:"required,min=1"`
	Method     int32  `json:"method" binding:"required,min=1"`
	HookURL    string `json:"hook_url" binding:"required"`
}

type UpdateAlarmRuleReq struct {
	ObjectIDReq
	Enable int `json:"enable" binding:"required,min=0,max=1"`
	// module_name 不支持修改
	CreateAlarmRuleReq
}

type DeleteAlarmRuleReq struct {
	ObjectIDReq
}
