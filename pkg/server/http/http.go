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
	adminSvc *admin.Service

	receiverSvc *receiver.Service
)

func NewAdminHttpServer(cfg *conf.Config, svc *admin.Service) runner.Server {
	adminSvc = svc

	midFlag := ginutil.MStd
	if cfg.Release() {
		midFlag = ginutil.MRelease | ginutil.MTraceId | ginutil.MRecoverLogger
	}
	httpHandler := ginutil.DefaultEngine(midFlag)
	registerAdminRouter(httpHandler, cfg)

	return runner.NewHttpServer(httpHandler)
}

func registerAdminRouter(e *gin.Engine, cfg *conf.Config) {
	e.POST("/v1/login", login)

	v1 := e.Group("/v1")
	{
		v1.GET("/module/list", findModuleList)
		v1.POST("/module", createModule)
		v1.PUT("/module", updateModule)
		v1.DELETE("/module", delModule)
	}

	// 配置报警规则
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

	// 单页应用
	e.StaticFile("/favicon.ico", "web/favicon.ico")
	e.Static("/static", "web/static")
	e.Static("/admin", "web")
}

func NewReceiverHttpServer(cfg *conf.Config, svc *receiver.Service) runner.Server {
	receiverSvc = svc

	midFlag := ginutil.MStd
	if cfg.Release() {
		midFlag = ginutil.MRelease | ginutil.MTraceId | ginutil.MRecoverLogger
	}
	httpHandler := ginutil.DefaultEngine(midFlag)
	registerReceiverRouter(httpHandler, cfg)

	return runner.NewHttpServer(httpHandler)
}

func registerReceiverRouter(e *gin.Engine, cfg *conf.Config) {
	e.POST("/v1/receiver/packet", receiverPacket)
}
