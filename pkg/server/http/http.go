package http

import (
	"github.com/bbdshow/bkit/ginutil"
	"github.com/bbdshow/bkit/runner"
	"github.com/bbdshow/qelog/pkg/admin"
	"github.com/bbdshow/qelog/pkg/conf"
	"github.com/bbdshow/qelog/pkg/receiver"
	"github.com/gin-gonic/gin"
)

var (
	adminSvc    *admin.Service
	receiverSvc *receiver.Service
	cfg         *conf.Config
)

func NewAdminHttpServer(c *conf.Config, svc *admin.Service) runner.Server {
	adminSvc = svc
	cfg = c

	midFlag := ginutil.MStd
	if cfg.Release() {
		midFlag = ginutil.MRelease | ginutil.MTraceId | ginutil.MRecoverLogger
	}
	// skip static file log
	ginutil.AddSkipPaths("/static/*filepath", "/admin/*filepath")

	httpHandler := ginutil.DefaultEngine(midFlag)
	registerAdminRouter(httpHandler)

	return runner.NewHttpServer(httpHandler)
}

func registerAdminRouter(e *gin.Engine) {
	e.POST("/v1/login", login)

	v1 := e.Group("/v1").Use(ginutil.JWTAuthVerify(cfg.Admin.AuthEnable))
	{
		v1.GET("/module/list", findModuleList)
		v1.POST("/module", createModule)
		v1.PUT("/module", updateModule)
		v1.DELETE("/module", delModule)
	}

	// alarm rule set
	{
		v1.GET("/alarmRule/list", findAlarmRuleList)
		v1.POST("/alarmRule", createAlarmRule)
		v1.PUT("/alarmRule", updateAlarmRule)
		v1.DELETE("/alarmRule", delAlarmRule)
		v1.GET("/alarmRule/hook/list", findHookURLList)
		v1.POST("/alarmRule/hook", createHookURL)
		v1.PUT("/alarmRule/hook", updateHookURL)
		v1.DELETE("/alarmRule/hook", delHookURL)
		v1.GET("/alarmRule/hook/ping", pingHookURL)
	}

	// log query
	{
		// Why use POST as query request method? The query condition text may be truncated because it is too large
		v1.POST("/logging/list", findLoggingList)
		v1.POST("/logging/traceid", findLoggingByTraceId)
		v1.DELETE("/logging/collection", dropLoggingCollection)
	}
	// log metrics
	{
		v1.GET("/metrics/dbStats", metricsDBStats)
		v1.GET("/metrics/collStats", metricsCollStats)
		v1.GET("/metrics/module/list", metricsModuleList)
		v1.GET("/metrics/module/trend", metricsModuleTrend)
	}

	// web static server
	e.StaticFile("/favicon.ico", "web/favicon.ico")
	e.Static("/static", "web/static")
	e.Static("/admin", "web")
}

func NewReceiverHttpServer(c *conf.Config, svc *receiver.Service) runner.Server {
	receiverSvc = svc
	cfg = c

	midFlag := ginutil.MStd
	if cfg.Release() {
		midFlag = ginutil.MRelease | ginutil.MTraceId | ginutil.MRecoverLogger
	}
	httpHandler := ginutil.DefaultEngine(midFlag)
	registerReceiverRouter(httpHandler)

	return runner.NewHttpServer(httpHandler)
}

func registerReceiverRouter(e *gin.Engine) {
	e.POST("/v1/receiver/packet", receiverPacket)
}
