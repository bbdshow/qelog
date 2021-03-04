package entity

type FindLoggingListReq struct {
	ShardingIndex  int    `json:"shardingIndex"binding:"required,min=0"`
	ModuleName     string `json:"moduleName" binding:"required"`
	Short          string `json:"short"`
	Level          int32  `json:"level" binding:"omitempty,min=-2,max=5"`
	IP             string `json:"ip"`
	ConditionOne   string `json:"conditionOne"`
	ConditionTwo   string `json:"conditionTwo"`
	ConditionThree string `json:"conditionThree"`
	// 指定查询集合
	ForceCollectionName string `json:"forceCollectionName"`
	TimeReq
	PageReq
}

type FindLoggingByTraceIDReq struct {
	ShardingIndex int    `json:"shardingIndex" binding:"required,min=0"`
	ModuleName    string `json:"moduleName" binding:"required"`
	TraceID       string `json:"traceId" binding:"required,gte=19"`
	// 指定查询集合
	ForceCollectionName string `json:"forceCollectionName"`
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
