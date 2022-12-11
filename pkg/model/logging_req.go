package model

type FindLoggingListReq struct {
	ModuleName     string `json:"moduleName" binding:"required"`
	Short          string `json:"short"`
	Level          int32  `json:"level" binding:"omitempty,min=-2,max=5"`
	IP             string `json:"ip"`
	ConditionOne   string `json:"conditionOne"`
	ConditionTwo   string `json:"conditionTwo"`
	ConditionThree string `json:"conditionThree"`
	// force collection name, used to when data is migrated shard
	ForceCollectionName string `json:"forceCollectionName"`
	ForceDatabase       string `json:"forceDatabase"`
	TimeReq
	PageReq
}

type FindLoggingByTraceIDReq struct {
	ModuleName          string `json:"moduleName" binding:"required"`
	TraceID             string `json:"traceId" binding:"required,gte=19"`
	ForceCollectionName string `json:"forceCollectionName"`
	ForceDatabase       string `json:"forceDatabase"`
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

type DropLoggingCollectionReq struct {
	ModuleName string `json:"moduleName" binding:"required"`
	Collection string `json:"collection" binding:"required"`
}
