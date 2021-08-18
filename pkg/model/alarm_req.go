package model

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

type FindHookURLListReq struct {
	OmitObjectIDReq
	Name    string `json:"name" form:"name"`
	KeyWord string `json:"keyWord" form:"keyWord"`
	Method  int32  `json:"method" form:"method"`
	PageReq
}

type FindHookURLList struct {
	ID           string   `json:"id"`
	Name         string   `json:"name"`
	URL          string   `json:"url"`
	Method       int32    `json:"method"`
	KeyWord      string   `json:"keyWord"`
	HideText     []string `json:"hideText"`
	UpdatedTsSec int64    `json:"updatedTsSec"`
}

type CreateHookURLReq struct {
	Name     string   `json:"name" binding:"required"`
	URL      string   `json:"url" binding:"required"`
	Method   int32    `json:"method" binding:"required"`
	KeyWord  string   `json:"keyWord" binding:"omitempty,lte=24"`
	HideText []string `json:"hideText"`
}

type UpdateHookURLReq struct {
	ObjectIDReq
	CreateHookURLReq
}

type DelHookURLReq struct {
	ObjectIDReq
}

type DelAlarmRuleReq struct {
	ObjectIDReq
}

type PingHookURLReq struct {
	ObjectIDReq
}
