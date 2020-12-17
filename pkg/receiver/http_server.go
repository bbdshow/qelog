package receiver

import (
	"net/http"
	"os"
	"time"

	"github.com/huzhongqing/qelog/pkg/storage"

	"github.com/gin-gonic/gin"
	"github.com/huzhongqing/qelog/libs/mongo"
)

type Receiver struct {
	database *mongo.Database
	server   *http.Server
	srv      *Service
}

func New(database *mongo.Database) *Receiver {
	r := &Receiver{
		database: database,
		srv:      NewService(storage.New(database)),
	}
	return r
}

func (rec *Receiver) Run(addr string) error {
	handler := gin.Default()
	if os.Getenv("ENV") == gin.ReleaseMode {
		gin.SetMode(gin.ReleaseMode)
		handler = gin.New()
	}

	rec.route(handler)

	rec.server = &http.Server{
		Addr:         addr,
		Handler:      handler,
		ReadTimeout:  90 * time.Second,
		WriteTimeout: 120 * time.Second,
	}

	return rec.server.ListenAndServe()
}

func (rec *Receiver) Close() error {
	if rec.server != nil {
		_ = rec.server.Close()
	}
	return nil
}

func (rec *Receiver) route(handler *gin.Engine, middleware ...gin.HandlerFunc) {
	handler.HEAD("/", func(c *gin.Context) { c.Status(200) })

	handler.Use(middleware...)
	handler.POST("/v1/receive/packet", rec.ReceivePacket)
}
