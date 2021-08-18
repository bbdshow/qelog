package admin

//type Handler struct {
//	svc *Service
//}
//
//func NewHandler() *Handler {
//	h := &Handler{svc: NewService()}
//	return h
//}

//func (h *Handler) FindModuleList(c *gin.Context) {
//	in := &entity.FindModuleListReq{}
//	if err := c.ShouldBind(in); err != nil {
//		httputil.RespError(c, httputil.ErrArgsInvalid.MergeError(err))
//		return
//	}
//
//	out := &entity.ListResp{}
//	if err := h.svc.FindModuleList(c.Request.Context(), in, out); err != nil {
//		httputil.RespError(c, err)
//		return
//	}
//	httputil.RespData(c, http.StatusOK, out)
//}
//
//func (h *Handler) CreateModule(c *gin.Context) {
//	in := &entity.CreateModuleReq{}
//	if err := c.ShouldBind(in); err != nil {
//		httputil.RespError(c, httputil.ErrArgsInvalid.MergeError(err))
//		return
//	}
//	if err := h.svc.CreateModule(c.Request.Context(), in); err != nil {
//		httputil.RespError(c, err)
//		return
//	}
//	httputil.RespSuccess(c)
//}
//
//func (h *Handler) UpdateModule(c *gin.Context) {
//	in := &entity.UpdateModuleReq{}
//	if err := c.ShouldBind(in); err != nil {
//		httputil.RespError(c, httputil.ErrArgsInvalid.MergeError(err))
//		return
//	}
//	if err := h.svc.UpdateModule(c.Request.Context(), in); err != nil {
//		httputil.RespError(c, err)
//		return
//	}
//	httputil.RespSuccess(c)
//}
//
//func (h *Handler) DeleteModule(c *gin.Context) {
//	in := &entity.DeleteModuleReq{}
//	if err := c.ShouldBind(in); err != nil {
//		httputil.RespError(c, httputil.ErrArgsInvalid.MergeError(err))
//		return
//	}
//	if err := h.svc.DeleteModule(c.Request.Context(), in); err != nil {
//		httputil.RespError(c, err)
//		return
//	}
//	httputil.RespSuccess(c)
//}

//func (h *Handler) FindLoggingList(c *gin.Context) {
//	in := &entity.FindLoggingListReq{}
//	if err := c.ShouldBind(in); err != nil {
//		httputil.RespError(c, httputil.ErrArgsInvalid.MergeError(err))
//		return
//	}
//	out := &entity.ListResp{}
//	if err := h.svc.FindLoggingList(c.Request.Context(), in, out); err != nil {
//		httputil.RespError(c, err)
//		return
//	}
//
//	httputil.RespData(c, http.StatusOK, out)
//}
//
//func (h *Handler) FindLoggingByTraceID(c *gin.Context) {
//	in := &entity.FindLoggingByTraceIDReq{}
//	if err := c.ShouldBind(in); err != nil {
//		httputil.RespError(c, httputil.ErrArgsInvalid.MergeError(err))
//		return
//	}
//	out := &entity.ListResp{}
//	if err := h.svc.FindLoggingByTraceID(c.Request.Context(), in, out); err != nil {
//		httputil.RespError(c, err)
//		return
//	}
//
//	httputil.RespData(c, http.StatusOK, out)
//}
//
//func (h *Handler) GetShardingIndex(c *gin.Context) {
//	out := &entity.GetShardingIndexResp{}
//	if err := h.svc.GetShardingIndex(c.Request.Context(), out); err != nil {
//		httputil.RespError(c, err)
//		return
//	}
//	httputil.RespData(c, http.StatusOK, out)
//}

//func (h *Handler) FindAlarmRuleList(c *gin.Context) {
//	in := &entity.FindAlarmRuleListReq{}
//	if err := c.ShouldBind(in); err != nil {
//		httputil.RespError(c, httputil.ErrArgsInvalid.MergeError(err))
//		return
//	}
//
//	out := &entity.ListResp{}
//	if err := h.svc.FindAlarmRuleList(c.Request.Context(), in, out); err != nil {
//		httputil.RespError(c, err)
//		return
//	}
//	httputil.RespData(c, http.StatusOK, out)
//}
//
//func (h *Handler) CreateAlarmRule(c *gin.Context) {
//	in := &entity.CreateAlarmRuleReq{}
//	if err := c.ShouldBind(in); err != nil {
//		httputil.RespError(c, httputil.ErrArgsInvalid.MergeError(err))
//		return
//	}
//	if err := h.svc.CreateAlarmRule(c.Request.Context(), in); err != nil {
//		httputil.RespError(c, err)
//		return
//	}
//	httputil.RespSuccess(c)
//}
//
//func (h *Handler) UpdateAlarmRule(c *gin.Context) {
//	in := &entity.UpdateAlarmRuleReq{}
//	if err := c.ShouldBind(in); err != nil {
//		httputil.RespError(c, httputil.ErrArgsInvalid.MergeError(err))
//		return
//	}
//	if err := h.svc.UpdateAlarmRule(c.Request.Context(), in); err != nil {
//		httputil.RespError(c, err)
//		return
//	}
//	httputil.RespSuccess(c)
//}
//
//func (h *Handler) DeleteAlarmRule(c *gin.Context) {
//	in := &entity.DeleteAlarmRuleReq{}
//	if err := c.ShouldBind(in); err != nil {
//		httputil.RespError(c, httputil.ErrArgsInvalid.MergeError(err))
//		return
//	}
//	if err := h.svc.DeleteAlarmRule(c.Request.Context(), in); err != nil {
//		httputil.RespError(c, err)
//		return
//	}
//	httputil.RespSuccess(c)
//}

//func (h *Handler) MetricsDBStats(c *gin.Context) {
//	out := &entity.ListResp{}
//	if err := h.svc.MetricsDBStats(c.Request.Context(), out); err != nil {
//		httputil.RespError(c, err)
//		return
//	}
//	httputil.RespData(c, http.StatusOK, out)
//}
//
//func (h *Handler) MetricsCollStats(c *gin.Context) {
//	in := &entity.MetricsCollStatsReq{}
//	if err := c.ShouldBind(in); err != nil {
//		httputil.RespError(c, httputil.ErrArgsInvalid.MergeError(err))
//		return
//	}
//
//	out := &entity.ListResp{}
//	if err := h.svc.MetricsCollStats(c.Request.Context(), in, out); err != nil {
//		httputil.RespError(c, err)
//		return
//	}
//	httputil.RespData(c, http.StatusOK, out)
//}
//
//func (h *Handler) MetricsModuleList(c *gin.Context) {
//	in := &entity.MetricsModuleListReq{}
//	if err := c.ShouldBind(in); err != nil {
//		httputil.RespError(c, httputil.ErrArgsInvalid.MergeError(err))
//		return
//	}
//	out := &entity.ListResp{}
//
//	if err := h.svc.MetricsModuleList(c.Request.Context(), in, out); err != nil {
//		httputil.RespError(c, err)
//		return
//	}
//	httputil.RespData(c, http.StatusOK, out)
//}
//
//func (h *Handler) MetricsModuleTrend(c *gin.Context) {
//	in := &entity.MetricsModuleTrendReq{}
//	if err := c.ShouldBind(in); err != nil {
//		httputil.RespError(c, httputil.ErrArgsInvalid.MergeError(err))
//		return
//	}
//	out := &entity.MetricsModuleTrendResp{}
//
//	if err := h.svc.MetricsModuleTrend(c.Request.Context(), in, out); err != nil {
//		httputil.RespError(c, err)
//		return
//	}
//	httputil.RespData(c, http.StatusOK, out)
//}
//
//func (h *Handler) DropLoggingCollection(c *gin.Context) {
//	in := &entity.DropLoggingCollectionReq{}
//	if err := c.ShouldBind(in); err != nil {
//		httputil.RespError(c, httputil.ErrArgsInvalid.MergeError(err))
//		return
//	}
//	if err := h.svc.DropLoggingCollection(c.Request.Context(), in); err != nil {
//		httputil.RespError(c, err)
//		return
//	}
//	httputil.RespSuccess(c)
//}

//func (h *Handler) FindHookURLList(c *gin.Context) {
//	in := &entity.FindHookURLListReq{}
//	if err := c.ShouldBind(in); err != nil {
//		httputil.RespError(c, httputil.ErrArgsInvalid.MergeError(err))
//		return
//	}
//	out := &entity.ListResp{}
//
//	if err := h.svc.FindHookURLList(c.Request.Context(), in, out); err != nil {
//		httputil.RespError(c, err)
//		return
//	}
//	httputil.RespData(c, http.StatusOK, out)
//}
//
//func (h *Handler) CreateHookURL(c *gin.Context) {
//	in := &entity.CreateHookURLReq{}
//	if err := c.ShouldBind(in); err != nil {
//		httputil.RespError(c, httputil.ErrArgsInvalid.MergeError(err))
//		return
//	}
//	if err := h.svc.CreateHookURL(c.Request.Context(), in); err != nil {
//		httputil.RespError(c, err)
//		return
//	}
//	httputil.RespSuccess(c)
//}
//
//func (h *Handler) UpdateHookURL(c *gin.Context) {
//	in := &entity.UpdateHookURLReq{}
//	if err := c.ShouldBind(in); err != nil {
//		httputil.RespError(c, httputil.ErrArgsInvalid.MergeError(err))
//		return
//	}
//	if err := h.svc.UpdateHookURL(c.Request.Context(), in); err != nil {
//		httputil.RespError(c, err)
//		return
//	}
//	httputil.RespSuccess(c)
//}
//
//func (h *Handler) DelHookURL(c *gin.Context) {
//	in := &entity.DelHookURLReq{}
//	if err := c.ShouldBind(in); err != nil {
//		httputil.RespError(c, httputil.ErrArgsInvalid.MergeError(err))
//		return
//	}
//	if err := h.svc.DelHookURL(c.Request.Context(), in); err != nil {
//		httputil.RespError(c, err)
//		return
//	}
//	httputil.RespSuccess(c)
//}
//
//func (h *Handler) PingHookURL(c *gin.Context) {
//	in := &entity.PingHookURLReq{}
//	if err := c.ShouldBind(in); err != nil {
//		httputil.RespError(c, httputil.ErrArgsInvalid.MergeError(err))
//		return
//	}
//	if err := h.svc.PingHookURL(c.Request.Context(), in); err != nil {
//		httputil.RespError(c, err)
//		return
//	}
//	httputil.RespSuccess(c)
//}
