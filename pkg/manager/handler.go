package manager

import (
	"net/http"
	"time"

	"github.com/huzhongqing/qelog/pkg/storage"

	"github.com/gin-gonic/gin"
	"github.com/huzhongqing/qelog/infra/httputil"
	"github.com/huzhongqing/qelog/infra/jwt"
	"github.com/huzhongqing/qelog/pkg/common/entity"
	"github.com/huzhongqing/qelog/pkg/config"
)

type Handler struct {
	srv *Service
}

func NewHandler() *Handler {
	h := &Handler{srv: NewService(storage.ShardingDB)}
	return h
}

func (h *Handler) Login(c *gin.Context) {
	in := &entity.LoginReq{}
	if err := c.ShouldBind(in); err != nil {
		httputil.RespError(c, httputil.ErrArgsInvalid.MergeError(err))
		return
	}
	if in.Username != config.Global.AdminUser.Username ||
		in.Password != config.Global.AdminUser.Password {
		httputil.RespError(c, httputil.NewError(httputil.ErrCodeUnauthorized, "账户或密码错误"))
		return
	}

	claims := jwt.NewCustomClaims(nil, 72*time.Hour)
	token, err := jwt.GenerateJWTToken(claims)
	if err != nil {
		httputil.RespError(c, httputil.NewError(httputil.ErrCodeUnauthorized, "系统异常，联系管理员"))
		return
	}

	out := &entity.LoginResp{
		Token: token,
	}
	httputil.RespData(c, http.StatusOK, out)
}

func (h *Handler) FindModuleList(c *gin.Context) {
	in := &entity.FindModuleListReq{}
	if err := c.ShouldBind(in); err != nil {
		httputil.RespError(c, httputil.ErrArgsInvalid.MergeError(err))
		return
	}

	out := &entity.ListResp{}
	if err := h.srv.FindModuleList(c.Request.Context(), in, out); err != nil {
		httputil.RespError(c, err)
		return
	}
	httputil.RespData(c, http.StatusOK, out)
}

func (h *Handler) CreateModule(c *gin.Context) {
	in := &entity.CreateModuleReq{}
	if err := c.ShouldBind(in); err != nil {
		httputil.RespError(c, httputil.ErrArgsInvalid.MergeError(err))
		return
	}
	if err := h.srv.CreateModule(c.Request.Context(), in); err != nil {
		httputil.RespError(c, err)
		return
	}
	httputil.RespSuccess(c)
}

func (h *Handler) UpdateModule(c *gin.Context) {
	in := &entity.UpdateModuleReq{}
	if err := c.ShouldBind(in); err != nil {
		httputil.RespError(c, httputil.ErrArgsInvalid.MergeError(err))
		return
	}
	if err := h.srv.UpdateModule(c.Request.Context(), in); err != nil {
		httputil.RespError(c, err)
		return
	}
	httputil.RespSuccess(c)
}

func (h *Handler) DeleteModule(c *gin.Context) {
	in := &entity.DeleteModuleReq{}
	if err := c.ShouldBind(in); err != nil {
		httputil.RespError(c, httputil.ErrArgsInvalid.MergeError(err))
		return
	}
	if err := h.srv.DeleteModule(c.Request.Context(), in); err != nil {
		httputil.RespError(c, err)
		return
	}
	httputil.RespSuccess(c)
}

func (h *Handler) FindLoggingList(c *gin.Context) {
	in := &entity.FindLoggingListReq{}
	if err := c.ShouldBind(in); err != nil {
		httputil.RespError(c, httputil.ErrArgsInvalid.MergeError(err))
		return
	}
	out := &entity.ListResp{}
	if err := h.srv.FindLoggingList(c.Request.Context(), in, out); err != nil {
		httputil.RespError(c, err)
		return
	}

	httputil.RespData(c, http.StatusOK, out)
}

func (h *Handler) FindLoggingByTraceID(c *gin.Context) {
	in := &entity.FindLoggingByTraceIDReq{}
	if err := c.ShouldBind(in); err != nil {
		httputil.RespError(c, httputil.ErrArgsInvalid.MergeError(err))
		return
	}
	out := &entity.ListResp{}
	if err := h.srv.FindLoggingByTraceID(c.Request.Context(), in, out); err != nil {
		httputil.RespError(c, err)
		return
	}

	httputil.RespData(c, http.StatusOK, out)
}

func (h *Handler) GetShardingIndex(c *gin.Context) {
	out := &entity.GetShardingIndexResp{}
	if err := h.srv.GetShardingIndex(c.Request.Context(), out); err != nil {
		httputil.RespError(c, err)
		return
	}
	httputil.RespData(c, http.StatusOK, out)
}

func (h *Handler) FindAlarmRuleList(c *gin.Context) {
	in := &entity.FindAlarmRuleListReq{}
	if err := c.ShouldBind(in); err != nil {
		httputil.RespError(c, httputil.ErrArgsInvalid.MergeError(err))
		return
	}

	out := &entity.ListResp{}
	if err := h.srv.FindAlarmRuleList(c.Request.Context(), in, out); err != nil {
		httputil.RespError(c, err)
		return
	}
	httputil.RespData(c, http.StatusOK, out)
}

func (h *Handler) CreateAlarmRule(c *gin.Context) {
	in := &entity.CreateAlarmRuleReq{}
	if err := c.ShouldBind(in); err != nil {
		httputil.RespError(c, httputil.ErrArgsInvalid.MergeError(err))
		return
	}
	if err := h.srv.CreateAlarmRule(c.Request.Context(), in); err != nil {
		httputil.RespError(c, err)
		return
	}
	httputil.RespSuccess(c)
}

func (h *Handler) UpdateAlarmRule(c *gin.Context) {
	in := &entity.UpdateAlarmRuleReq{}
	if err := c.ShouldBind(in); err != nil {
		httputil.RespError(c, httputil.ErrArgsInvalid.MergeError(err))
		return
	}
	if err := h.srv.UpdateAlarmRule(c.Request.Context(), in); err != nil {
		httputil.RespError(c, err)
		return
	}
	httputil.RespSuccess(c)
}

func (h *Handler) DeleteAlarmRule(c *gin.Context) {
	in := &entity.DeleteAlarmRuleReq{}
	if err := c.ShouldBind(in); err != nil {
		httputil.RespError(c, httputil.ErrArgsInvalid.MergeError(err))
		return
	}
	if err := h.srv.DeleteAlarmRule(c.Request.Context(), in); err != nil {
		httputil.RespError(c, err)
		return
	}
	httputil.RespSuccess(c)
}

func (h *Handler) MetricsDBStats(c *gin.Context) {
	out := &entity.ListResp{}
	if err := h.srv.MetricsDBStats(c.Request.Context(), out); err != nil {
		httputil.RespError(c, err)
		return
	}
	httputil.RespData(c, http.StatusOK, out)
}

func (h *Handler) MetricsCollStats(c *gin.Context) {
	in := &entity.MetricsCollStatsReq{}
	if err := c.ShouldBind(in); err != nil {
		httputil.RespError(c, httputil.ErrArgsInvalid.MergeError(err))
		return
	}

	out := &entity.ListResp{}
	if err := h.srv.MetricsCollStats(c.Request.Context(), in, out); err != nil {
		httputil.RespError(c, err)
		return
	}
	httputil.RespData(c, http.StatusOK, out)
}

func (h *Handler) MetricsModuleList(c *gin.Context) {
	in := &entity.MetricsModuleListReq{}
	if err := c.ShouldBind(in); err != nil {
		httputil.RespError(c, httputil.ErrArgsInvalid.MergeError(err))
		return
	}
	out := &entity.ListResp{}

	if err := h.srv.MetricsModuleList(c.Request.Context(), in, out); err != nil {
		httputil.RespError(c, err)
		return
	}
	httputil.RespData(c, http.StatusOK, out)
}

func (h *Handler) MetricsModuleTrend(c *gin.Context) {
	in := &entity.MetricsModuleTrendReq{}
	if err := c.ShouldBind(in); err != nil {
		httputil.RespError(c, httputil.ErrArgsInvalid.MergeError(err))
		return
	}
	out := &entity.MetricsModuleTrendResp{}

	if err := h.srv.MetricsModuleTrend(c.Request.Context(), in, out); err != nil {
		httputil.RespError(c, err)
		return
	}
	httputil.RespData(c, http.StatusOK, out)
}

func (h *Handler) DropLoggingCollection(c *gin.Context) {
	in := &entity.DropLoggingCollectionReq{}
	if err := c.ShouldBind(in); err != nil {
		httputil.RespError(c, httputil.ErrArgsInvalid.MergeError(err))
		return
	}
	if err := h.srv.DropLoggingCollection(c.Request.Context(), in); err != nil {
		httputil.RespError(c, err)
		return
	}
	httputil.RespSuccess(c)
}

func (h *Handler) FindHookURLList(c *gin.Context) {
	in := &entity.FindHookURLListReq{}
	if err := c.ShouldBind(in); err != nil {
		httputil.RespError(c, httputil.ErrArgsInvalid.MergeError(err))
		return
	}
	out := &entity.ListResp{}

	if err := h.srv.FindHookURLList(c.Request.Context(), in, out); err != nil {
		httputil.RespError(c, err)
		return
	}
	httputil.RespData(c, http.StatusOK, out)
}

func (h *Handler) CreateHookURL(c *gin.Context) {
	in := &entity.CreateHookURLReq{}
	if err := c.ShouldBind(in); err != nil {
		httputil.RespError(c, httputil.ErrArgsInvalid.MergeError(err))
		return
	}
	if err := h.srv.CreateHookURL(c.Request.Context(), in); err != nil {
		httputil.RespError(c, err)
		return
	}
	httputil.RespSuccess(c)
}

func (h *Handler) UpdateHookURL(c *gin.Context) {
	in := &entity.UpdateHookURLReq{}
	if err := c.ShouldBind(in); err != nil {
		httputil.RespError(c, httputil.ErrArgsInvalid.MergeError(err))
		return
	}
	if err := h.srv.UpdateHookURL(c.Request.Context(), in); err != nil {
		httputil.RespError(c, err)
		return
	}
	httputil.RespSuccess(c)
}

func (h *Handler) DelHookURL(c *gin.Context) {
	in := &entity.DelHookURLReq{}
	if err := c.ShouldBind(in); err != nil {
		httputil.RespError(c, httputil.ErrArgsInvalid.MergeError(err))
		return
	}
	if err := h.srv.DelHookURL(c.Request.Context(), in); err != nil {
		httputil.RespError(c, err)
		return
	}
	httputil.RespSuccess(c)
}

func (h *Handler) PingHookURL(c *gin.Context) {
	in := &entity.PingHookURLReq{}
	if err := c.ShouldBind(in); err != nil {
		httputil.RespError(c, httputil.ErrArgsInvalid.MergeError(err))
		return
	}
	if err := h.srv.PingHookURL(c.Request.Context(), in); err != nil {
		httputil.RespError(c, err)
		return
	}
	httputil.RespSuccess(c)
}
