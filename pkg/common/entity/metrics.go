package entity

type MetricsCountResp struct {
	// 数据库最新统计周期
	DBCount DBCount `json:"dbCount"`
	// 项目最新统计周期
	//ModuleCount ModuleCount `json:"moduleCount"`
}

type DBCount struct {
	DBName      string `json:"dbName"`
	Collections int32  `json:"collections"`
	DataSize    int64  `json:"dataSize"`
	StorageSize int64  `json:"storageSize"`
	IndexSize   int64  `json:"indexSize"`
	Objects     int64  `json:"objects"`
	Indexs      int64  `json:"indexs"`
}

type ModuleCount struct {
	Modules     int64 `json:"modules"`
	Numbers     int64 `json:"numbers"`
	LoggingSize int64 `json:"loggingSize"`
}

type MetricsModuleListReq struct {
	ModuleName string `json:"moduleName" form:"moduleName"`
	DateTsSec  int64  `json:"dateTsSec" form:"dateTsSec" binding:"required,min=1"`
	PageReq
}

type MetricsModuleList struct {
	ModuleName   string `json:"moduleName"`
	Number       int64  `json:"number"`
	Size         int64  `json:"size"`
	CreatedTsSec int64  `json:"createdTsSec"`
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
