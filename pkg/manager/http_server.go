package manager

import (
	"net/http"
	"os"
	"time"

	"github.com/huzhongqing/qelog/pkg/common/model"

	"github.com/huzhongqing/qelog/pkg/common/entity"

	"github.com/huzhongqing/qelog/pkg/httputil"

	"github.com/huzhongqing/qelog/pkg/storage"

	"github.com/gin-gonic/gin"
	"github.com/huzhongqing/qelog/libs/mongo"
)

type HTTPService struct {
	database *mongo.Database
	server   *http.Server
	manager  *Service
}

func NewHTTPService(database *mongo.Database) *HTTPService {
	srv := &HTTPService{
		database: database,
		manager:  NewService(storage.New(database)),
	}
	return srv
}

func (srv *HTTPService) Run(addr string) error {
	handler := gin.Default()
	if os.Getenv("ENV") == gin.ReleaseMode {
		gin.SetMode(gin.ReleaseMode)
		handler = gin.New()
	} else {
		if err := srv.database.UpsertCollectionIndexMany(model.ModuleRegisterIndexMany()); err != nil {
			return err
		}
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
	if srv.server != nil {
		_ = srv.server.Close()
	}
	return nil
}

func (srv *HTTPService) route(handler *gin.Engine, middleware ...gin.HandlerFunc) {
	handler.HEAD("/", func(c *gin.Context) { c.Status(200) })

	v1 := handler.Group("/v1", middleware...)
	v1.POST("/module", srv.CreateModuleRegister)
}

func (srv *HTTPService) CreateModuleRegister(c *gin.Context) {
	var arg entity.CreateModuleRegisterReq
	if err := c.ShouldBind(&arg); err != nil {
		httputil.RespError(c, httputil.ErrArgsInvalid.MergeError(err))
		return
	}

	if err := srv.manager.CreateModuleRegister(c.Request.Context(), &arg); err != nil {
		httputil.RespError(c, httputil.ErrSystemException.MergeError(err))
		return
	}
	httputil.RespSuccess(c)
}
