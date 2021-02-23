package receiver

import (
	"net/http"
	"time"

	"github.com/huzhongqing/qelog/api"

	"github.com/gin-gonic/gin"
	"github.com/huzhongqing/qelog/infra/httputil"
	"github.com/huzhongqing/qelog/pkg/config"
	"github.com/huzhongqing/qelog/pkg/storage"
)

type HTTPService struct {
	server   *http.Server
	receiver *Service
}

func NewHTTPService() *HTTPService {
	srv := &HTTPService{
		receiver: NewService(storage.ShardingDB),
	}
	return srv
}

func (srv *HTTPService) Run(addr string) error {
	handler := gin.New()
	if config.Global.Release() {
		gin.SetMode(gin.ReleaseMode)
	}
	handler.Use(gin.Recovery())

	handler.HEAD("/", func(c *gin.Context) { c.Status(200) })
	handler.POST("/v1/receiver/packet", srv.ReceivePacket)

	srv.server = &http.Server{
		Addr:         addr,
		Handler:      handler,
		ReadTimeout:  120 * time.Second,
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

func (srv *HTTPService) ReceivePacket(c *gin.Context) {
	in := &api.JSONPacket{}
	if err := c.ShouldBind(in); err != nil {
		httputil.RespError(c, httputil.ErrArgsInvalid)
		return
	}

	if err := srv.receiver.InsertJSONPacket(c.Request.Context(), c.ClientIP(), in); err != nil {
		httputil.RespDataWithError(c, http.StatusBadRequest, nil, err)
		return
	}
	httputil.RespSuccess(c)
}
