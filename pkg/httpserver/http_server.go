package httpserver

import (
	"github.com/huzhongqing/qelog/infra/httputil"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"

)

type HTTPServer struct {
	env     string
	server  *http.Server
	handler *gin.Engine
}

func NewHTTPServer(env string) *HTTPServer {
	handler := gin.New()
	if env == gin.ReleaseMode {
		gin.SetMode(gin.ReleaseMode)
		handler.Use(httputil.GinLogger([]string{"/health","/admin", "/static"}), httputil.GinRecoveryWithLogger())
	} else {
		handler.Use(gin.Logger(), httputil.GinRecoveryWithLogger())
	}

	handler.HEAD("/health", func(c *gin.Context) { c.Status(200) })

	return &HTTPServer{
		server:  nil,
		env:     env,
		handler: handler,
	}
}

// 统一使用此 handler 注册路由
func (srv *HTTPServer) Engine() *gin.Engine {
	// 注册TraceID
	srv.handler.Use(httputil.HandlerRegisterTraceID())

	if srv.env != gin.ReleaseMode {
		// 测试环境，记录请求返回日志
		srv.handler.Use(httputil.HandlerLogging(true, []string{"/docs/", "/health"}...))
	}
	return srv.handler
}

func (srv *HTTPServer) Run(addr string) error {

	srv.server = &http.Server{
		Addr:         addr,
		Handler:      srv.handler,
		ReadTimeout:  90 * time.Second,
		WriteTimeout: 120 * time.Second,
	}
	return srv.server.ListenAndServe()
}

func (srv *HTTPServer) Close() error {
	if srv.server != nil {
		_ = srv.server.Close()
	}
	return nil
}
