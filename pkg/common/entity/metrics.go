package entity

type MetricsIndexResp struct {
	// 数据库最新统计周期
	DBName      string `json:"dbName"`
	Collections int32  `json:"collections"`
	DataSize    string `json:"dataSize"`
	StorageSize string `json:"storageSize"`
	IndexSize   string `json:"indexSize"`
	Objects     int64  `json:"objects"`
	// 项目最新统计周期
	Modules     int64 `json:"modules"`
	Numbers     int64 `json:"numbers"`
	LoggingSize int64 `json:"loggingSize"`
}

type MetricsModuleListReq struct {
	ModuleName string
	PageReq
}

type MetricsModuleList struct {
	ModuleName   string `json:"moduleName"`
	Number       int64  `json:"number"`
	Size         string `json:"size"`
	LevelsTop    []NumberData
	IPsTop       []NumberData
	CreatedTsSec int64
}

type NumberData struct {
	Name   string `json:"name"`
	Number int64  `json:"number"`
}

type MetricsModuleReq struct {
	ModuleName string ` binding:"required"`
	LastDay    int    ` binding:"required,min=1,max=7"`
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
