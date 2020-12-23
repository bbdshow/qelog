package entity

type CreateModuleRegisterReq struct {
	ModuleName string `json:"module_name" binding:"required,gte=2,lte=24"`
	DBIndex    int32  `json:"db_index" binding:"required,min=1,max=16"`
	Desc       string `json:"desc" binding:"required"`
}
