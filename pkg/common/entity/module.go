package entity

type CreateModuleReq struct {
	Name          string `json:"name" binding:"required,gte=2,lte=24,lowercase"`
	ShardingIndex int    `json:"shardingIndex" binding:"required,min=1,max=16"`
	Desc          string `json:"desc" binding:"omitempty,gte=1,lte=128"`
}

type FindModuleListReq struct {
	Name string `json:"name" form:"name"`
	PageReq
}

type FindModuleList struct {
	ID                   string `json:"id"`
	Name                 string `json:"name"`
	Desc                 string `json:"desc"`
	ShardingIndex        int    `json:"shardingIndex"`
	HistoryShardingIndex []int  `json:"historyShardingIndex"`
	UpdatedTsSec         int64  `json:"updatedTsSec"`
}

type UpdateModuleReq struct {
	ObjectIDReq
	ShardingIndex int    `json:"shardingIndex" binding:"required,min=1,max=16"`
	Desc          string `json:"desc" binding:"required,gte=1,lte=128"`
}

type DeleteModuleReq struct {
	ObjectIDReq
	Name string `json:"name" binding:"required"`
}

type GetShardingIndexResp struct {
	SuggestIndex      int                  `json:"suggestIndex"`
	ShardingIndexSize int                  `json:"shardingIndexSize"`
	UseState          []ShardingIndexState `json:"useState"`
}

type ShardingIndexState struct {
	Index int `json:"index"`
	Use   int `json:"use"`
}
