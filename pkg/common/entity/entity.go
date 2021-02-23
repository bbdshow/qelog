package entity

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type PageReq struct {
	Page  int64 `json:"page" form:"page" default:"1" `
	Limit int64 `json:"limit" from:"limit" default:"20"`
}

func (v PageReq) SetPage(opt *options.FindOptions) {
	opt.SetSkip((v.Page - 1) * v.Limit)
	opt.SetLimit(v.Limit)
}

type ObjectIDReq struct {
	ID string `json:"id" form:"id" binding:"required,len=24"`
}
type OmitObjectIDReq struct {
	ID string `json:"id" form:"id" binding:"omitempty,len=24"`
}

func (v OmitObjectIDReq) ObjectID() (primitive.ObjectID, error) {
	return primitive.ObjectIDFromHex(v.ID)
}

func (v ObjectIDReq) ObjectID() (primitive.ObjectID, error) {
	return primitive.ObjectIDFromHex(v.ID)
}

type TimeReq struct {
	BeginTsSec int64 `json:"beginTsSec" form:"beginTsSec"`
	EndTsSec   int64 `json:"endTsSec" form:"endTsSec"`
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
	DBIndex int32  `json:"dbIndex" binding:"required,min=1,max=16"`
	Desc    string `json:"desc" binding:"omitempty,gte=1,lte=128"`
}

type FindModuleListReq struct {
	Name string `json:"name" form:"name"`
	PageReq
}

type FindModuleList struct {
	ID             string  `json:"id"`
	Name           string  `json:"name"`
	Desc           string  `json:"desc"`
	DBIndex        int32   `json:"dbIndex"`
	HistoryDBIndex []int32 `json:"historyDbIndex"`
	UpdatedTsSec   int64   `json:"updatedTsSec"`
}

type UpdateModuleReq struct {
	ObjectIDReq
	DBIndex int32  `json:"dbIndex" binding:"required,min=1,max=16"`
	Desc    string `json:"desc" binding:"required,gte=1,lte=128"`
}

type DeleteModuleReq struct {
	ObjectIDReq
	Name string `json:"name" binding:"required"`
}

type GetDBIndexResp struct {
	SuggestDBIndex int32          `json:"suggestDbIndex"`
	MaxDBIndex     int32          `json:"maxDbIndex"`
	UseState       []DBIndexState `json:"useState"`
}
type DBIndexState struct {
	Index int32 `json:"index"`
	Use   int   `json:"use"`
}

type FindLoggingListReq struct {
	DBIndex        int32  `json:"dbIndex" binding:"required,min=0"`
	ModuleName     string `json:"moduleName" binding:"required"`
	Short          string `json:"short"`
	Level          int32  `json:"level" binding:"omitempty,min=-2,max=5"`
	IP             string `json:"ip"`
	ConditionOne   string `json:"conditionOne"`
	ConditionTwo   string `json:"conditionTwo"`
	ConditionThree string `json:"conditionThree"`
	TimeReq
	PageReq
}

type FindLoggingByTraceIDReq struct {
	DBIndex    int32  `json:"dbIndex" binding:"required,min=0"`
	ModuleName string `json:"moduleName" binding:"required"`
	TraceID    string `json:"traceId" binding:"required,gte=19"`
}

type FindLoggingList struct {
	ID             string `json:"id"`
	TsMill         int64  `json:"tsMill"`
	Level          int32  `json:"level"`
	Short          string `json:"short"`
	Full           string `json:"full"`
	ConditionOne   string `json:"conditionOne"`
	ConditionTwo   string `json:"conditionTwo"`
	ConditionThree string `json:"conditionThree"`
	TraceID        string `json:"traceId"`
	IP             string `json:"ip"`
}

type FindAlarmRuleListReq struct {
	Enable     int    `json:"enable" binding:"omitempty,min=-1,max=1"`
	ModuleName string `json:"moduleName"`
	Short      string `json:"short"`
	PageReq
}

type FindAlarmRuleList struct {
	ID           string `json:"id"`
	Enable       bool   `json:"enable"`
	ModuleName   string `json:"moduleName"`
	Short        string `json:"short"`
	Level        int32  `json:"level"`
	Tag          string `json:"tag"`
	RateSec      int64  `json:"rateSec"`
	Method       int32  `json:"method"`
	HookID       string `json:"hookId"`
	UpdatedTsSec int64  `json:"updatedTsSec"`
}

type CreateAlarmRuleReq struct {
	ModuleName string `json:"moduleName" binding:"required"`
	Short      string `json:"short" binding:"required"`
	Level      int32  `json:"level" binding:"min=-1,max=8"`
	Tag        string `json:"tag" binding:"omitempty,gte=1,lte=128"`
	RateSec    int64  `json:"rateSec" binding:"min=0"`
	Method     int32  `json:"method" binding:"required,min=1"`
	HookID     string `json:"hookId" binding:"required,len=24"`
}

type UpdateAlarmRuleReq struct {
	ObjectIDReq
	Enable bool `json:"enable"`
	// module_name 不支持修改
	CreateAlarmRuleReq
}

type DeleteAlarmRuleReq struct {
	ObjectIDReq
}

type DropLoggingCollectionReq struct {
	Host string `json:"host" binding:"required"`
	Name string `json:"name" binding:"required"`
}

type FindHookURLListReq struct {
	OmitObjectIDReq
	Name    string `json:"name" form:"name"`
	KeyWord string `json:"keyWord" form:"keyWord"`
	Method  int32  `json:"method" form:"method"`
	PageReq
}
type FindHookURLList struct {
	ID           string `json:"id"`
	Name         string `json:"name"`
	URL          string `json:"url"`
	Method       int32  `json:"method"`
	KeyWord      string `json:"keyWord"`
	UpdatedTsSec int64  `json:"updatedTsSec"`
}

type CreateHookURLReq struct {
	Name    string `json:"name" binding:"required"`
	URL     string `json:"url" binding:"required"`
	Method  int32  `json:"method" binding:"required"`
	KeyWord string `json:"keyWord" binding:"omitempty,lte=24"`
}

type UpdateHookURLReq struct {
	ObjectIDReq
	CreateHookURLReq
}

type DelHookURLReq struct {
	ObjectIDReq
}

type PingHookURLReq struct {
	ObjectIDReq
}
