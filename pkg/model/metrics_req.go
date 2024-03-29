package model

import "sort"

type MetricsDBStatsResp struct {
	DBStats []DBStat `json:"dbStats"`
}

// DBStat database stat
type DBStat struct {
	Host         string `json:"host"`
	DBName       string `json:"dbName"`
	Collections  int32  `json:"collections"`
	DataSize     int64  `json:"dataSize"`
	StorageSize  int64  `json:"storageSize"`
	IndexSize    int64  `json:"indexSize"`
	Objects      int64  `json:"objects"`
	Indexs       int64  `json:"indexs"`
	UpdatedTsSec int64  `json:"updatedTsSec"`
}

type MetricsCollStatsReq struct {
	ModuleName string `json:"moduleName" form:"moduleName" binding:"required"`
}

type CollStat struct {
	ModuleName     string           `json:"moduleName"`
	Host           string           `json:"host"`
	DBName         string           `json:"dbName"`
	Name           string           `json:"name"`
	Size           int64            `json:"size"`
	Count          int64            `json:"count"`
	AvgObjSize     int64            `json:"avgObjSize"`
	StorageSize    int64            `json:"storageSize"`
	Capped         bool             `json:"capped"`
	TotalIndexSize int64            `json:"totalIndexSize"`
	IndexSizes     map[string]int64 `json:"indexSizes"`
	UpdatedTsSec   int64            `json:"updatedTsSec"`
	CreatedTsSec   int64            `json:"createdTsSec"`
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

type MetricsModuleTrendReq struct {
	ModuleName string `json:"moduleName" form:"moduleName" binding:"required"`
	LastDay    int    `json:"lastDay" form:"lastDay" binding:"required,min=1,max=7"`
}

type MetricsModuleTrendResp struct {
	// X coordinates
	XData       []string `json:"xData"`
	LegendData  []string `json:"legendData"`
	LevelSeries []Serie  `json:"levelSeries"`
	IPSeries    []Serie  `json:"ipSeries"`
	Number      int64    `json:"number"`
	Size        int64    `json:"size"`
}

type Serie struct {
	Index int32   `json:"index"`
	Name  string  `json:"name"`
	Type  string  `json:"type"`
	Color string  `json:"color"`
	Data  []int32 `json:"data"`
}
type AscSeries []Serie

func (v AscSeries) Len() int           { return len(v) }
func (v AscSeries) Swap(i, j int)      { v[i], v[j] = v[j], v[i] }
func (v AscSeries) Less(i, j int) bool { return v[i].Index < v[j].Index }

type DescSeries []Serie

func (v DescSeries) Len() int           { return len(v) }
func (v DescSeries) Swap(i, j int)      { v[i], v[j] = v[j], v[i] }
func (v DescSeries) Less(i, j int) bool { return v[i].Index > v[j].Index }
func SortSeries(series []Serie, order string) {
	if order == "ASC" {
		sort.Sort(AscSeries(series))
		return
	}
	if order == "DESC" {
		sort.Sort(DescSeries(series))
	}
}
