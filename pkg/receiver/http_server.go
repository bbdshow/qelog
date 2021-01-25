package receiver

import (
	"net/http"
	"time"

	"github.com/huzhongqing/qelog/libs/logs"
	"go.uber.org/zap"

	"github.com/gin-gonic/gin"
	"github.com/huzhongqing/qelog/api"
	"github.com/huzhongqing/qelog/pkg/config"
	"github.com/huzhongqing/qelog/pkg/httputil"
	"github.com/huzhongqing/qelog/pkg/storage"
)

type HTTPService struct {
	server   *http.Server
	receiver *Service
}

func NewHTTPService(sharding *storage.Sharding) *HTTPService {
	srv := &HTTPService{
		receiver: NewService(sharding),
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
	srv.receiver.Sync()
	if srv.server != nil {
		_ = srv.server.Close()
	}
	return nil
}

func (srv *HTTPService) route(handler *gin.Engine) {
	handler.HEAD("/", func(c *gin.Context) { c.Status(200) })
	handler.POST("/v1/receiver/packet", srv.ReceivePacket)
}

func (srv *HTTPService) ReceivePacket(c *gin.Context) {
	in := &api.JSONPacket{}
	if err := c.ShouldBind(in); err != nil {
		httputil.RespError(c, httputil.ErrArgsInvalid)
		return
	}

	if err := srv.receiver.InsertJSONPacket(c.Request.Context(), c.ClientIP(), in); err != nil {
		httputil.RespStatusCodeWithError(c, http.StatusBadRequest, err)
		return
	}
	httputil.RespSuccess(c)
}
