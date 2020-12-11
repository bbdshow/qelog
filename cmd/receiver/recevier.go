package receiver

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/huzhongqing/qelog/config"
)

type Receiver struct {
	cfg *config.Config

	server *http.Server
}

func New(cfg *config.Config) (*Receiver, error) {
	r := &Receiver{
		cfg: cfg,
	}
	return r, nil
}

func (r *Receiver) Run() error {
	handler := gin.Default()
	if r.cfg.Release() {
		gin.SetMode(gin.ReleaseMode)
		handler = gin.New()
	}

	r.route(handler)

	r.server = &http.Server{
		Addr:         r.cfg.ReceiverAddr,
		Handler:      handler,
		ReadTimeout:  90 * time.Second,
		WriteTimeout: 120 * time.Second,
	}

	return r.server.ListenAndServe()
}

func (r *Receiver) Close() error {
	if r.server != nil {
		_ = r.server.Close()
	}
	return nil
}

func (r *Receiver) route(handler *gin.Engine, mids ...gin.HandlerFunc) {
	handler.HEAD("/", func(c *gin.Context) { c.Status(200) })

	handleFunc := NewHandleFunc()

	handler.Use(mids...)
	handler.POST("/v1/receive/packet", handleFunc.ReceivePacket)
}
