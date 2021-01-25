package manager

import (
	"net/http"
	"time"

	"github.com/huzhongqing/qelog/libs/logs"
	"go.uber.org/zap"

	"github.com/huzhongqing/qelog/libs/jwt"

	"github.com/huzhongqing/qelog/pkg/config"

	"github.com/huzhongqing/qelog/pkg/common/entity"

	"github.com/huzhongqing/qelog/pkg/httputil"

	"github.com/huzhongqing/qelog/pkg/storage"

	"github.com/gin-gonic/gin"
)

type HTTPService struct {
	server  *http.Server
	manager *Service
}

func NewHTTPService(sharding *storage.Sharding) *HTTPService {
	srv := &HTTPService{
		manager: NewService(sharding),
	}
	return srv
}

func (srv *HTTPService) Run(addr string) error {
	handler := gin.New()
	if config.GlobalConfig.Release() {
		gin.SetMode(gin.ReleaseMode)
		gin.DefaultErrorWriter = logs.Qezap.Clone().SetWritePrefix("[GIN-Recovery]").SetWriteLevel(zap.ErrorLevel)
		handler.Use(httputil.GinLogger([]string{"/"}), gin.Recovery())
	} else {
		handler.Use(gin.Logger(), gin.Recovery())
	}

	srv.route(handler)

	srv.server = &http.Server{
		Addr:         addr,
		Handler:      handler,
		ReadTimeout:  90 * time.Second,
		WriteTimeout: 120 * time.Second,
	}
	return srv.server.ListenAndServe()
}

func (srv *HTTPService) Close() error {
	if srv.server != nil {
		_ = srv.server.Close()
	}
	return nil
}

func (srv *HTTPService) route(handler *gin.Engine) {
	handler.HEAD("/", func(c *gin.Context) { c.Status(200) })

	handler.POST("/v1/login", srv.Login)

	v1 := handler.Group("/v1", AuthVerify(config.GlobalConfig.AuthEnable), httputil.HandlerRegisterTraceID())
	module := v1.Group("/module", httputil.HandlerLogging(true))
	{
		module.GET("/list", srv.FindModuleList)
		module.POST("", srv.CreateModule)
		module.PUT("", srv.UpdateModule)
		module.DELETE("", srv.DeleteModule)
	}
	// 配置报警规则
	alarmRule := v1.Group("/alarm-rule", httputil.HandlerLogging(true))
	{
		alarmRule.GET("/list", srv.FindAlarmRuleList)
		alarmRule.POST("", srv.CreateAlarmRule)
		alarmRule.PUT("", srv.UpdateAlarmRule)
		alarmRule.DELETE("", srv.DeleteAlarmRule)
	}

	// 获取 db 信息
	v1.GET("/db-index", srv.GetDBIndex)

	// 搜索日志
	logging := v1.Group("/logging")
	{
		logging.POST("/list", srv.FindLoggingList)
		logging.POST("/traceid", srv.FindLoggingByTraceID)
		logging.DELETE("/collection", srv.DropLoggingCollection)
	}

	// 报表
	metrics := v1.Group("/metrics")
	{
		metrics.GET("/dbstats", srv.MetricsDBStats)
		metrics.GET("/collstats", srv.MetricsCollStats)
		metrics.GET("/module/list", srv.MetricsModuleList)
		metrics.GET("/module/trend", srv.MetricsModuleTrend)
	}

	// 单页应用
	handler.StaticFile("/favicon.ico", "web/favicon.ico")
	handler.Static("/static", "web/static")
	handler.Static("/admin", "web")
}

func (srv *HTTPService) Login(c *gin.Context) {
	in := &entity.LoginReq{}
	if err := c.ShouldBind(in); err != nil {
		httputil.RespError(c, httputil.ErrArgsInvalid.MergeError(err))
		return
	}
	if in.Username != config.GlobalConfig.AdminUser.Username ||
		in.Password != config.GlobalConfig.AdminUser.Password {
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

func (srv *HTTPService) FindModuleList(c *gin.Context) {
	in := &entity.FindModuleListReq{}
	if err := c.ShouldBind(in); err != nil {
		httputil.RespError(c, httputil.ErrArgsInvalid.MergeError(err))
		return
	}

	out := &entity.ListResp{}
	if err := srv.manager.FindModuleList(c.Request.Context(), in, out); err != nil {
		httputil.RespError(c, err)
		return
	}
	httputil.RespData(c, http.StatusOK, out)
}

func (srv *HTTPService) CreateModule(c *gin.Context) {
	in := &entity.CreateModuleReq{}
	if err := c.ShouldBind(in); err != nil {
		httputil.RespError(c, httputil.ErrArgsInvalid.MergeError(err))
		return
	}
	if err := srv.manager.CreateModule(c.Request.Context(), in); err != nil {
		httputil.RespError(c, err)
		return
	}
	httputil.RespSuccess(c)
}

func (srv *HTTPService) UpdateModule(c *gin.Context) {
	in := &entity.UpdateModuleReq{}
	if err := c.ShouldBind(in); err != nil {
		httputil.RespError(c, httputil.ErrArgsInvalid.MergeError(err))
		return
	}
	if err := srv.manager.UpdateModule(c.Request.Context(), in); err != nil {
		httputil.RespError(c, err)
		return
	}
	httputil.RespSuccess(c)
}

func (srv *HTTPService) DeleteModule(c *gin.Context) {
	in := &entity.DeleteModuleReq{}
	if err := c.ShouldBind(in); err != nil {
		httputil.RespError(c, httputil.ErrArgsInvalid.MergeError(err))
		return
	}
	if err := srv.manager.DeleteModule(c.Request.Context(), in); err != nil {
		httputil.RespError(c, err)
		return
	}
	httputil.RespSuccess(c)
}

func (srv *HTTPService) FindLoggingList(c *gin.Context) {
	in := &entity.FindLoggingListReq{}
	if err := c.ShouldBind(in); err != nil {
		httputil.RespError(c, httputil.ErrArgsInvalid.MergeError(err))
		return
	}
	out := &entity.ListResp{}
	if err := srv.manager.FindLoggingList(c.Request.Context(), in, out); err != nil {
		httputil.RespError(c, err)
		return
	}

	httputil.RespData(c, http.StatusOK, out)
}

func (srv *HTTPService) FindLoggingByTraceID(c *gin.Context) {
	in := &entity.FindLoggingByTraceIDReq{}
	if err := c.ShouldBind(in); err != nil {
		httputil.RespError(c, httputil.ErrArgsInvalid.MergeError(err))
		return
	}
	out := &entity.ListResp{}
	if err := srv.manager.FindLoggingByTraceID(c.Request.Context(), in, out); err != nil {
		httputil.RespError(c, err)
		return
	}

	httputil.RespData(c, http.StatusOK, out)
}

func (srv *HTTPService) GetDBIndex(c *gin.Context) {
	out := &entity.GetDBIndexResp{}
	if err := srv.manager.GetDBIndex(c.Request.Context(), out); err != nil {
		httputil.RespError(c, err)
		return
	}
	httputil.RespData(c, http.StatusOK, out)
}

func (srv *HTTPService) FindAlarmRuleList(c *gin.Context) {
	in := &entity.FindAlarmRuleListReq{}
	if err := c.ShouldBind(in); err != nil {
		httputil.RespError(c, httputil.ErrArgsInvalid.MergeError(err))
		return
	}

	out := &entity.ListResp{}
	if err := srv.manager.FindAlarmRuleList(c.Request.Context(), in, out); err != nil {
		httputil.RespError(c, err)
		return
	}
	httputil.RespData(c, http.StatusOK, out)
}

func (srv *HTTPService) CreateAlarmRule(c *gin.Context) {
	in := &entity.CreateAlarmRuleReq{}
	if err := c.ShouldBind(in); err != nil {
		httputil.RespError(c, httputil.ErrArgsInvalid.MergeError(err))
		return
	}
	if err := srv.manager.CreateAlarmRule(c.Request.Context(), in); err != nil {
		httputil.RespError(c, err)
		return
	}
	httputil.RespSuccess(c)
}

func (srv *HTTPService) UpdateAlarmRule(c *gin.Context) {
	in := &entity.UpdateAlarmRuleReq{}
	if err := c.ShouldBind(in); err != nil {
		httputil.RespError(c, httputil.ErrArgsInvalid.MergeError(err))
		return
	}
	if err := srv.manager.UpdateAlarmRule(c.Request.Context(), in); err != nil {
		httputil.RespError(c, err)
		return
	}
	httputil.RespSuccess(c)
}

func (srv *HTTPService) DeleteAlarmRule(c *gin.Context) {
	in := &entity.DeleteAlarmRuleReq{}
	if err := c.ShouldBind(in); err != nil {
		httputil.RespError(c, httputil.ErrArgsInvalid.MergeError(err))
		return
	}
	if err := srv.manager.DeleteAlarmRule(c.Request.Context(), in); err != nil {
		httputil.RespError(c, err)
		return
	}
	httputil.RespSuccess(c)
}

func (srv *HTTPService) MetricsDBStats(c *gin.Context) {
	out := &entity.ListResp{}
	if err := srv.manager.MetricsDBStats(c.Request.Context(), out); err != nil {
		httputil.RespError(c, err)
		return
	}
	httputil.RespData(c, http.StatusOK, out)
}

func (srv *HTTPService) MetricsCollStats(c *gin.Context) {
	in := &entity.MetricsCollStatsReq{}
	if err := c.ShouldBind(in); err != nil {
		httputil.RespError(c, httputil.ErrArgsInvalid.MergeError(err))
		return
	}

	out := &entity.ListResp{}
	if err := srv.manager.MetricsCollStats(c.Request.Context(), in, out); err != nil {
		httputil.RespError(c, err)
		return
	}
	httputil.RespData(c, http.StatusOK, out)
}

func (srv *HTTPService) MetricsModuleList(c *gin.Context) {
	in := &entity.MetricsModuleListReq{}
	if err := c.ShouldBind(in); err != nil {
		httputil.RespError(c, httputil.ErrArgsInvalid.MergeError(err))
		return
	}
	out := &entity.ListResp{}

	if err := srv.manager.MetricsModuleList(c.Request.Context(), in, out); err != nil {
		httputil.RespError(c, err)
		return
	}
	httputil.RespData(c, http.StatusOK, out)
}

func (srv *HTTPService) MetricsModuleTrend(c *gin.Context) {
	in := &entity.MetricsModuleTrendReq{}
	if err := c.ShouldBind(in); err != nil {
		httputil.RespError(c, httputil.ErrArgsInvalid.MergeError(err))
		return
	}
	out := &entity.MetricsModuleTrendResp{}

	if err := srv.manager.MetricsModuleTrend(c.Request.Context(), in, out); err != nil {
		httputil.RespError(c, err)
		return
	}
	httputil.RespData(c, http.StatusOK, out)
}

func (srv *HTTPService) DropLoggingCollection(c *gin.Context) {
	in := &entity.DropLoggingCollectionReq{}
	if err := c.ShouldBind(in); err != nil {
		httputil.RespError(c, httputil.ErrArgsInvalid.MergeError(err))
		return
	}
	if err := srv.manager.DropLoggingCollection(c.Request.Context(), in); err != nil {
		httputil.RespError(c, err)
		return
	}
	httputil.RespSuccess(c)
}
