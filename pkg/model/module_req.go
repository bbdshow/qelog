package model

type CreateModuleReq struct {
	Name     string `json:"name" binding:"required,gte=2,lte=64,lowercase"`
	Desc     string `json:"desc" binding:"omitempty,gte=1,lte=128"`
	DaySpan  int    `json:"daySpan" binding:"omitempty,gte=1,lte=31"`
	MaxMonth int    `json:"maxMonth" binding:"omitempty,gte=1"`
}

type FindModuleListReq struct {
	Name string `json:"name" form:"name"`
	PageReq
}

type FindModuleList struct {
	ID           string `json:"id"`
	Name         string `json:"name"`
	Desc         string `json:"desc"`
	Bucket       string `json:"bucket"`
	DaySpan      int    `json:"daySpan"`
	MaxMonth     int    `json:"maxMonth"`
	Database     string `json:"database"`
	Prefix       string `json:"prefix"`
	UpdatedTsSec int64  `json:"updatedTsSec"`
}

type UpdateModuleReq struct {
	ObjectIDReq
	Bucket   string `json:"bucket" binding:"omitempty,gte=1,lte=12"`
	DaySpan  int    `json:"daySpan" binding:"omitempty,gte=1,lte=31"`
	MaxMonth int    `json:"maxMonth" binding:"omitempty,gte=1"`
	Database string `json:"database" binding:"omitempty,gte=1"`
	Prefix   string `json:"prefix" binding:"omitempty,gte=1,lte=12"`
	Desc     string `json:"desc" binding:"omitempty,gte=1,lte=128"`
}

type DelModuleReq struct {
	ObjectIDReq
	Name string `json:"name" binding:"required"`
}
