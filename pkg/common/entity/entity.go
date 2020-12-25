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
	BeginUnix int64 `json:"begin_unix" form:"begin_unix"`
	EndUnix   int64 `json:"end_unix" form:"end_unix"`
}

func (v TimeReq) BeginAt() time.Time {
	return time.Unix(v.BeginUnix, 0)
}
func (v TimeReq) EndAt() time.Time {
	return time.Unix(v.EndUnix, 0)
}

func (v TimeReq) DefaultSection(t time.Duration) (b, e time.Time) {
	// 没有查询时间
	if v.BeginUnix <= 0 && v.EndUnix <= 0 {
		e = time.Now()
		b = e.Add(-t)
		return b, e
	}
	// 只有结束时间
	if v.BeginUnix <= 0 && v.EndUnix > 0 {
		e = time.Unix(v.EndUnix, 0)
		b = e.Add(-t)
		return b, e
	}
	// 只有开始时间
	if v.BeginUnix > 0 && v.EndUnix <= 0 {
		b = time.Unix(v.BeginUnix, 0)
		e = time.Now()
		return b, e
	}

	// 时间都存在
	return time.Unix(v.BeginUnix, 0), time.Unix(v.EndUnix, 0)
}

type ListResp struct {
	Count int64       `json:"count"`
	List  interface{} `json:"list"`
}

type CreateModuleReq struct {
	Name    string `json:"name" binding:"required,gte=2,lte=24"`
	DBIndex int32  `json:"db_index" binding:"required,min=1,max=16"`
	Desc    string `json:"desc" binding:"required,lte=128"`
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
	UpdatedAt      string  `json:"updated_at"`
}

type UpdateModuleReq struct {
	ObjectIDReq
	DBIndex int32 `json:"db_index" binding:"required,min=1,max=16"`
}

type DeleteModuleReq struct {
	ObjectIDReq
	Name string `json:"name" binding:"required"`
}

type GetDBIndexResp struct {
	SuggestDBIndex int32           `json:"suggest_db_index"`
	MaxDBIndex     int32           `json:"max_db_index"`
	UseState       map[int32]int32 `json:"use_state"`
}

type FindLoggingListReq struct {
	DBIndex        int32  `json:"db_index" binding:"required"`
	ModuleName     string `json:"module_name" binding:"required"`
	ShortMsg       string `json:"short_msg"`
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
	TimeUnixMill   int64  `json:"time_unix_mill"`
	Level          int    `json:"level"`
	ShortMsg       string `json:"short_msg"`
	Full           string `json:"full"`
	ConditionOne   string `json:"condition_one"`
	ConditionTwo   string `json:"condition_two"`
	ConditionThree string `json:"condition_three"`
	IP             string `json:"ip"`
}
