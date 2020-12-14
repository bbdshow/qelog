package receiver

import (
	"context"
	"net/http"
	"time"

	"github.com/huzhongqing/qelog/libs/mongoclient"

	"github.com/gin-gonic/gin"
	"github.com/huzhongqing/qelog/config"
)

type Receiver struct {
	cfg *config.Config

	database *mongoclient.Database
	server   *http.Server
}

func New(cfg *config.Config) (*Receiver, error) {
	r := &Receiver{
		cfg: cfg,
	}

	ctx, _ := context.WithTimeout(context.Background(), 30*time.Second)
	cli, err := mongoclient.NewMongoClientByURI(ctx, cfg.MongoDB.URI)
	if err != nil {
		return nil, err
	}
	r.database = &mongoclient.Database{Database: cli.Database(cfg.MongoDB.DataBase)}

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
	if r.database != nil {
		ctx, _ := context.WithTimeout(context.Background(), 5*time.Second)
		_ = r.database.Client().Disconnect(ctx)
	}
	return nil
}

func (r *Receiver) route(handler *gin.Engine, mids ...gin.HandlerFunc) {
	handler.HEAD("/", func(c *gin.Context) { c.Status(200) })

	handleFunc := NewHandleFunc(r.database)

	handler.Use(mids...)
	handler.POST("/v1/receive/packet", handleFunc.ReceivePacket)
}
