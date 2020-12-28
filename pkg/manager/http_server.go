package manager

import (
	"net/http"
	"os"
	"time"

	"github.com/huzhongqing/qelog/libs/logs"
	"go.uber.org/zap"

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
		gin.DefaultErrorWriter = logs.Qezap.Clone().SetWritePrefix("ginError").SetWriteLevel(zap.ErrorLevel)
		gin.DefaultWriter = logs.Qezap.Clone().SetWritePrefix("ginDebug").SetWriteLevel(zap.DebugLevel)
	} else {
		if err := srv.database.UpsertCollectionIndexMany(
			model.ModuleIndexMany()); err != nil {
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

	handler.POST("/login")

	v1 := handler.Group("/v1", middleware...)

	v1.GET("/module/list", srv.FindModuleList)
	//v1.GET("/module")
	v1.POST("/module", srv.CreateModule)
	v1.PUT("/module", srv.UpdateModule)
	v1.DELETE("/module", srv.DeleteModule)

	// 配置报警规则
	v1.GET("/alarm-rule/list", srv.FindAlarmRuleList)
	v1.POST("/alarm-rule", srv.CreateAlarmRule)
	v1.PUT("/alarm-rule", srv.UpdateAlarmRule)
	v1.DELETE("/alarm-rule", srv.DeleteAlarmRule)

	// 获取 db 信息
	v1.GET("/db-index", srv.GetDBIndex)

	// 搜索日志
	v1.POST("/logging/list", srv.FindLoggingList)

	// 报表
}
func (srv *HTTPService) FindModuleList(c *gin.Context) {
	in := &entity.FindModuleListReq{}
	if err := c.ShouldBind(in); err != nil {
		httputil.RespError(c, httputil.ErrArgsInvalid.MergeError(err))
		return
	}

	out := &entity.ListResp{}
	if err := srv.manager.FindModuleList(c.Request.Context(), in, out); err != nil {
		httputil.RespError(c, err)
		return
	}
	httputil.RespData(c, http.StatusOK, out)
}

func (srv *HTTPService) CreateModule(c *gin.Context) {
	in := &entity.CreateModuleReq{}
	if err := c.ShouldBind(in); err != nil {
		httputil.RespError(c, httputil.ErrArgsInvalid.MergeError(err))
		return
	}
	if err := srv.manager.CreateModule(c.Request.Context(), in); err != nil {
		httputil.RespError(c, err)
		return
	}
	httputil.RespSuccess(c)
}

func (srv *HTTPService) UpdateModule(c *gin.Context) {
	in := &entity.UpdateModuleReq{}
	if err := c.ShouldBind(in); err != nil {
		httputil.RespError(c, httputil.ErrArgsInvalid.MergeError(err))
		return
	}
	if err := srv.manager.UpdateModule(c.Request.Context(), in); err != nil {
		httputil.RespError(c, err)
		return
	}
	httputil.RespSuccess(c)
}

func (srv *HTTPService) DeleteModule(c *gin.Context) {
	in := &entity.DeleteModuleReq{}
	if err := c.ShouldBind(in); err != nil {
		httputil.RespError(c, httputil.ErrArgsInvalid.MergeError(err))
		return
	}
	if err := srv.manager.DeleteModule(c.Request.Context(), in); err != nil {
		httputil.RespError(c, err)
		return
	}
	httputil.RespSuccess(c)
}

func (srv *HTTPService) FindLoggingList(c *gin.Context) {
	in := &entity.FindLoggingListReq{}
	if err := c.ShouldBind(in); err != nil {
		httputil.RespError(c, httputil.ErrArgsInvalid.MergeError(err))
		return
	}
	out := &entity.ListResp{}
	if err := srv.manager.FindLoggingList(c.Request.Context(), in, out); err != nil {
		httputil.RespError(c, err)
		return
	}

	httputil.RespData(c, http.StatusOK, out)
}

func (srv *HTTPService) GetDBIndex(c *gin.Context) {
	out := &entity.GetDBIndexResp{}
	if err := srv.manager.GetDBIndex(c.Request.Context(), out); err != nil {
		httputil.RespError(c, err)
		return
	}
	httputil.RespData(c, http.StatusOK, out)
}

func (srv *HTTPService) FindAlarmRuleList(c *gin.Context) {
	in := &entity.FindAlarmRuleListReq{}
	if err := c.ShouldBind(in); err != nil {
		httputil.RespError(c, httputil.ErrArgsInvalid.MergeError(err))
		return
	}

	out := &entity.ListResp{}
	if err := srv.manager.FindAlarmRuleList(c.Request.Context(), in, out); err != nil {
		httputil.RespError(c, err)
		return
	}
	httputil.RespData(c, http.StatusOK, out)
}

func (srv *HTTPService) CreateAlarmRule(c *gin.Context) {
	in := &entity.CreateAlarmRuleReq{}
	if err := c.ShouldBind(in); err != nil {
		httputil.RespError(c, httputil.ErrArgsInvalid.MergeError(err))
		return
	}
	if err := srv.manager.CreateAlarmRule(c.Request.Context(), in); err != nil {
		httputil.RespError(c, err)
		return
	}
	httputil.RespSuccess(c)
}

func (srv *HTTPService) UpdateAlarmRule(c *gin.Context) {
	in := &entity.UpdateAlarmRuleReq{}
	if err := c.ShouldBind(in); err != nil {
		httputil.RespError(c, httputil.ErrArgsInvalid.MergeError(err))
		return
	}
	if err := srv.manager.UpdateAlarmRule(c.Request.Context(), in); err != nil {
		httputil.RespError(c, err)
		return
	}
	httputil.RespSuccess(c)
}

func (srv *HTTPService) DeleteAlarmRule(c *gin.Context) {
	in := &entity.DeleteAlarmRuleReq{}
	if err := c.ShouldBind(in); err != nil {
		httputil.RespError(c, httputil.ErrArgsInvalid.MergeError(err))
		return
	}
	if err := srv.manager.DeleteAlarmRule(c.Request.Context(), in); err != nil {
		httputil.RespError(c, err)
		return
	}
	httputil.RespSuccess(c)
}
