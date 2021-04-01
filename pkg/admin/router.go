package admin

import (
	"github.com/gin-gonic/gin"
	"github.com/huzhongqing/qelog/infra/httputil"
	"github.com/huzhongqing/qelog/pkg/config"
)

func RegisterRouter(route *gin.Engine) {
	h := NewHandler()

	route.HEAD("/", func(c *gin.Context) { c.Status(200) })

	route.POST("/v1/login", h.Login)

	v1 := route.Group("/v1", httputil.AuthAdmin(config.Global.AuthEnable), httputil.HandlerRegisterTraceID())
	module := v1.Group("/module", httputil.HandlerLogging(true))
	{
		module.GET("/list", h.FindModuleList)
		module.POST("", h.CreateModule)
		module.PUT("", h.UpdateModule)
		module.DELETE("", h.DeleteModule)
	}
	// 配置报警规则
	alarmRule := v1.Group("/alarmRule", httputil.HandlerLogging(true))
	{
		alarmRule.GET("/list", h.FindAlarmRuleList)
		alarmRule.POST("", h.CreateAlarmRule)
		alarmRule.PUT("", h.UpdateAlarmRule)
		alarmRule.DELETE("", h.DeleteAlarmRule)
		alarmRule.GET("/hook/list", h.FindHookURLList)
		alarmRule.POST("/hook", h.CreateHookURL)
		alarmRule.PUT("/hook", h.UpdateHookURL)
		alarmRule.DELETE("/hook", h.DelHookURL)
		alarmRule.GET("/hook/ping", h.PingHookURL)
	}

	// 获取分片使用信息
	v1.GET("/shardingIndex", h.GetShardingIndex)

	// 搜索日志
	logging := v1.Group("/logging")
	{
		logging.POST("/list", h.FindLoggingList)
		logging.POST("/traceid", h.FindLoggingByTraceID)
		logging.DELETE("/collection", h.DropLoggingCollection)
	}

	// 报表
	metrics := v1.Group("/metrics")
	{
		metrics.GET("/dbStats", h.MetricsDBStats)
		metrics.GET("/collStats", h.MetricsCollStats)
		metrics.GET("/module/list", h.MetricsModuleList)
		metrics.GET("/module/trend", h.MetricsModuleTrend)
	}

	// 单页应用
	route.StaticFile("/favicon.ico", "web/favicon.ico")
	route.Static("/static", "web/static")
	route.Static("/admin", "web")

}
