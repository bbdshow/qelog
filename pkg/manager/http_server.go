package manager

import (
	"net/http"
	"time"

	"github.com/huzhongqing/qelog/infra/httputil"
	"github.com/huzhongqing/qelog/pkg/config"

	"github.com/gin-gonic/gin"
)

type HTTPService struct {
	server *http.Server
}

func NewHTTPService() *HTTPService {
	srv := &HTTPService{}
	return srv
}

func (srv *HTTPService) Run(addr string) error {
	handler := gin.New()
	if config.Global.Release() {
		gin.SetMode(gin.ReleaseMode)
		handler.Use(httputil.GinLogger([]string{"/health", "/admin", "/static"}), httputil.GinRecoveryWithLogger())
	} else {
		handler.Use(gin.Logger(), gin.Recovery())
	}

	RegisterRouter(handler)

	srv.server = &http.Server{
		Addr:         addr,
		Handler:      handler,
		ReadTimeout:  120 * time.Second,
		WriteTimeout: 90 * time.Second,
	}
	return srv.server.ListenAndServe()
}

func (srv *HTTPService) Close() error {
	if srv.server != nil {
		_ = srv.server.Close()
	}
	return nil
}

func RegisterRouter(route *gin.Engine, midd ...gin.HandlerFunc) {
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
