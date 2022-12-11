package http

import (
	"time"

	"github.com/bbdshow/bkit/auth/jwt"
	"github.com/bbdshow/bkit/errc"
	"github.com/bbdshow/bkit/ginutil"
	"github.com/bbdshow/qelog/pkg/model"
	"github.com/gin-gonic/gin"
)

func login(c *gin.Context) {
	in := &model.LoginReq{}
	if err := ginutil.ShouldBind(c, in); err != nil {
		ginutil.RespErr(c, err)
		return
	}
	if in.Username != cfg.Admin.Username ||
		in.Password != cfg.Admin.Password {
		ginutil.RespErr(c, errc.ErrAuthInvalid.MultiMsg("user or passwd invalid"))
		return
	}

	claims := jwt.NewCustomClaims("", 72*time.Hour)
	token, err := jwt.GenerateJWTToken(claims)
	if err != nil {
		ginutil.RespErr(c, errc.ErrAuthInternalErr.MultiMsg("system exception,try again"))
		return
	}
	out := &model.LoginResp{
		Token: token,
	}
	ginutil.RespData(c, out)
}

func findModuleList(c *gin.Context) {
	in := &model.FindModuleListReq{}
	if err := ginutil.ShouldBind(c, in); err != nil {
		ginutil.RespErr(c, err)
		return
	}
	out := &model.ListResp{}
	if err := adminSvc.FindModuleList(c.Request.Context(), in, out); err != nil {
		ginutil.RespErr(c, err)
		return
	}
	ginutil.RespData(c, out)
}

func createModule(c *gin.Context) {
	in := &model.CreateModuleReq{}
	if err := ginutil.ShouldBind(c, in); err != nil {
		ginutil.RespErr(c, err)
		return
	}
	if err := adminSvc.CreateModule(c.Request.Context(), in); err != nil {
		ginutil.RespErr(c, err)
		return
	}
	ginutil.RespSuccess(c)
}

func updateModule(c *gin.Context) {
	in := &model.UpdateModuleReq{}
	if err := ginutil.ShouldBind(c, in); err != nil {
		ginutil.RespErr(c, err)
		return
	}
	if err := adminSvc.UpdateModule(c.Request.Context(), in); err != nil {
		ginutil.RespErr(c, err)
		return
	}
	ginutil.RespSuccess(c)
}

func delModule(c *gin.Context) {
	in := &model.DelModuleReq{}
	if err := ginutil.ShouldBind(c, in); err != nil {
		ginutil.RespErr(c, err)
		return
	}
	if err := adminSvc.DelModule(c.Request.Context(), in); err != nil {
		ginutil.RespErr(c, err)
		return
	}
	ginutil.RespSuccess(c)
}

func findAlarmRuleList(c *gin.Context) {
	in := &model.FindAlarmRuleListReq{}
	if err := ginutil.ShouldBind(c, in); err != nil {
		ginutil.RespErr(c, err)
		return
	}
	out := &model.ListResp{}
	if err := adminSvc.FindAlarmRuleList(c.Request.Context(), in, out); err != nil {
		ginutil.RespErr(c, err)
		return
	}
	ginutil.RespData(c, out)
}

func createAlarmRule(c *gin.Context) {
	in := &model.CreateAlarmRuleReq{}
	if err := ginutil.ShouldBind(c, in); err != nil {
		ginutil.RespErr(c, err)
		return
	}
	if err := adminSvc.CreateAlarmRule(c.Request.Context(), in); err != nil {
		ginutil.RespErr(c, err)
		return
	}
	ginutil.RespSuccess(c)
}

func updateAlarmRule(c *gin.Context) {
	in := &model.UpdateAlarmRuleReq{}
	if err := ginutil.ShouldBind(c, in); err != nil {
		ginutil.RespErr(c, err)
		return
	}
	if err := adminSvc.UpdateAlarmRule(c.Request.Context(), in); err != nil {
		ginutil.RespErr(c, err)
		return
	}
	ginutil.RespSuccess(c)
}

func delAlarmRule(c *gin.Context) {
	in := &model.DelAlarmRuleReq{}
	if err := ginutil.ShouldBind(c, in); err != nil {
		ginutil.RespErr(c, err)
		return
	}
	if err := adminSvc.DelAlarmRule(c.Request.Context(), in); err != nil {
		ginutil.RespErr(c, err)
		return
	}
	ginutil.RespSuccess(c)
}

func findHookURLList(c *gin.Context) {
	in := &model.FindHookURLListReq{}
	if err := ginutil.ShouldBind(c, in); err != nil {
		ginutil.RespErr(c, err)
		return
	}
	out := &model.ListResp{}

	if err := adminSvc.FindHookURLList(c.Request.Context(), in, out); err != nil {
		ginutil.RespErr(c, err)
		return
	}
	ginutil.RespData(c, out)
}

func createHookURL(c *gin.Context) {
	in := &model.CreateHookURLReq{}
	if err := ginutil.ShouldBind(c, in); err != nil {
		ginutil.RespErr(c, err)
		return
	}
	if err := adminSvc.CreateHookURL(c.Request.Context(), in); err != nil {
		ginutil.RespErr(c, err)
		return
	}
	ginutil.RespSuccess(c)
}

func updateHookURL(c *gin.Context) {
	in := &model.UpdateHookURLReq{}
	if err := ginutil.ShouldBind(c, in); err != nil {
		ginutil.RespErr(c, err)
		return
	}
	if err := adminSvc.UpdateHookURL(c.Request.Context(), in); err != nil {
		ginutil.RespErr(c, err)
		return
	}
	ginutil.RespSuccess(c)
}

func delHookURL(c *gin.Context) {
	in := &model.DelHookURLReq{}
	if err := ginutil.ShouldBind(c, in); err != nil {
		ginutil.RespErr(c, err)
		return
	}
	if err := adminSvc.DelHookURL(c.Request.Context(), in); err != nil {
		ginutil.RespErr(c, err)
		return
	}
	ginutil.RespSuccess(c)
}

func pingHookURL(c *gin.Context) {
	in := &model.PingHookURLReq{}
	if err := ginutil.ShouldBind(c, in); err != nil {
		ginutil.RespErr(c, err)
		return
	}
	if err := adminSvc.PingHookURL(c.Request.Context(), in); err != nil {
		ginutil.RespErr(c, err)
		return
	}
	ginutil.RespSuccess(c)
}

func findLoggingList(c *gin.Context) {
	in := &model.FindLoggingListReq{}
	if err := ginutil.ShouldBind(c, in); err != nil {
		ginutil.RespErr(c, err)
		return
	}
	out := &model.ListResp{}
	if err := adminSvc.FindLoggingList(c.Request.Context(), in, out); err != nil {
		ginutil.RespErr(c, err)
		return
	}
	ginutil.RespData(c, out)
}

func findLoggingByTraceId(c *gin.Context) {
	in := &model.FindLoggingByTraceIDReq{}
	if err := ginutil.ShouldBind(c, in); err != nil {
		ginutil.RespErr(c, err)
		return
	}
	out := &model.ListResp{}
	if err := adminSvc.FindLoggingByTraceID(c.Request.Context(), in, out); err != nil {
		ginutil.RespErr(c, err)
		return
	}
	ginutil.RespData(c, out)
}

func metricsDBStats(c *gin.Context) {
	out := &model.ListResp{}
	if err := adminSvc.MetricsDBStats(c.Request.Context(), out); err != nil {
		ginutil.RespErr(c, err)
		return
	}
	ginutil.RespData(c, out)
}

func metricsCollStats(c *gin.Context) {
	in := &model.MetricsCollStatsReq{}
	if err := ginutil.ShouldBind(c, in); err != nil {
		ginutil.RespErr(c, err)
		return
	}
	out := &model.ListResp{}
	if err := adminSvc.MetricsCollStats(c.Request.Context(), in, out); err != nil {
		ginutil.RespErr(c, err)
		return
	}
	ginutil.RespData(c, out)
}

func metricsModuleList(c *gin.Context) {
	in := &model.MetricsModuleListReq{}
	if err := ginutil.ShouldBind(c, in); err != nil {
		ginutil.RespErr(c, err)
		return
	}
	out := &model.ListResp{}
	if err := adminSvc.MetricsModuleList(c.Request.Context(), in, out); err != nil {
		ginutil.RespErr(c, err)
		return
	}
	ginutil.RespData(c, out)
}

func metricsModuleTrend(c *gin.Context) {
	in := &model.MetricsModuleTrendReq{}
	if err := ginutil.ShouldBind(c, in); err != nil {
		ginutil.RespErr(c, err)
		return
	}
	out := &model.MetricsModuleTrendResp{}
	if err := adminSvc.MetricsModuleTrend(c.Request.Context(), in, out); err != nil {
		ginutil.RespErr(c, err)
		return
	}
	ginutil.RespData(c, out)
}

func dropLoggingCollection(c *gin.Context) {
	in := &model.DropLoggingCollectionReq{}
	if err := ginutil.ShouldBind(c, in); err != nil {
		ginutil.RespErr(c, err)
		return
	}
	if err := adminSvc.DropLoggingCollection(c.Request.Context(), in); err != nil {
		ginutil.RespErr(c, err)
		return
	}
	ginutil.RespSuccess(c)
}
