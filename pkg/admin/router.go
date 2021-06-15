package admin

import (
	"github.com/bbdshow/qelog/infra/httputil"
	"github.com/bbdshow/qelog/pkg/config"
	"github.com/gin-gonic/gin"
)

func RegisterRouter(route *gin.Engine) {
	h := NewHandler()
	route.POST("/v1/login", h.Login)

	v1 := route.Group("/v1", httputil.AuthAdmin(config.Global.AuthEnable), httputil.HandlerLogging(true))
	{
		v1.GET("/module/list", h.FindModuleList)
		v1.POST("/module", h.CreateModule)
		v1.PUT("/module", h.UpdateModule)
		v1.DELETE("/module", h.DeleteModule)
	}
	// 配置报警规则
	{
		v1.GET("/alarmRule/list", h.FindAlarmRuleList)
		v1.POST("/alarmRule", h.CreateAlarmRule)
		v1.PUT("/alarmRule", h.UpdateAlarmRule)
		v1.DELETE("/alarmRule", h.DeleteAlarmRule)
		v1.GET("/alarmRule/hook/list", h.FindHookURLList)
		v1.POST("/alarmRule/hook", h.CreateHookURL)
		v1.PUT("/alarmRule/hook", h.UpdateHookURL)
		v1.DELETE("/alarmRule/hook", h.DelHookURL)
		v1.GET("/alarmRule/hook/ping", h.PingHookURL)
	}

	// 获取分片使用信息
	v1.GET("/shardingIndex", h.GetShardingIndex)

	// 搜索日志
	{
		v1.POST("/logging/list", h.FindLoggingList)
		v1.POST("/logging/traceid", h.FindLoggingByTraceID)
		v1.DELETE("/logging/collection", h.DropLoggingCollection)
	}

	// 报表
	{
		v1.GET("/metrics/dbStats", h.MetricsDBStats)
		v1.GET("/metrics/collStats", h.MetricsCollStats)
		v1.GET("/metrics/module/list", h.MetricsModuleList)
		v1.GET("/metrics/module/trend", h.MetricsModuleTrend)
	}

	// 单页应用
	route.StaticFile("/favicon.ico", "web/favicon.ico")
	route.Static("/static", "web/static")
	route.Static("/admin", "web")
}
