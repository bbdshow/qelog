package entity

type MetricsIndexResp struct {
	// 数据库最新统计周期
	DBName      string `json:"db_name"`
	Collections int32  `json:"collections"`
	DataSize    string `json:"data_size"`
	StorageSize string `json:"storage_size"`
	IndexSize   string `json:"index_size"`
	Objects     int64  `json:"objects"`
	// 项目最新统计周期
	Modules     int64 `json:"modules"`
	Numbers     int64 `json:"numbers"`
	LoggingSize int64 `json:"logging_size"`
}

type MetricsModuleListReq struct {
	ModuleName string `json:"module_name"`
	PageReq
}

type MetricsModuleReq struct {
	ModuleName string `json:"module_name"`
	LastDay    int    `json:"last_day"`
}

type MetricsModule struct {
}

type MetricsIPUpTop struct {
}
