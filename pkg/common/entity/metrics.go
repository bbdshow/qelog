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

type MetricsModuleList struct {
	ModuleName   string       `json:"module_name"`
	Number       int64        `json:"number"`
	Size         string       `json:"size"`
	LevelsTop    []NumberData `json:"levels_top"`
	IPsTop       []NumberData `json:"i_ps_top"`
	CreatedTsSec int64        `json:"created_ts_sec"`
}

type NumberData struct {
	Name   string `json:"name"`
	Number int64  `json:"number"`
}

type MetricsModuleReq struct {
	ModuleName string `json:"module_name" binding:"required"`
	LastDay    int    `json:"last_day" binding:"required,min=1,max=7"`
}

type MetricsModule struct {
	ModuleName string
	Numbers    int64
	Size       string
	Levels     []string
	LevelData  []LineNumData
	IPs        []string
	IPData     []LineNumData
}

type LineNumData struct {
	Name string
	Data []int64
}

type MetricsIPUpTop struct {
}
