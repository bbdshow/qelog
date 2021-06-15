package receiver

import (
	"github.com/bbdshow/qelog/api"
	"github.com/bbdshow/qelog/infra/httputil"
	"github.com/gin-gonic/gin"
	"net/http"
)

type HttpHandler struct {
	receiver *Service
}

func NewHttpHandler() *HttpHandler {
	srv := &HttpHandler{
		receiver: NewService(),
	}
	return srv
}

func (h *HttpHandler) ReceivePacket(c *gin.Context) {
	in := &api.JSONPacket{}
	if err := c.ShouldBind(in); err != nil {
		httputil.RespError(c, httputil.ErrArgsInvalid)
		return
	}

	if err := h.receiver.InsertJSONPacket(c.Request.Context(), c.ClientIP(), in); err != nil {
		httputil.RespDataWithError(c, http.StatusBadRequest, nil, err)
		return
	}
	httputil.RespSuccess(c)
}

func RegisterRouter(route *gin.Engine) {
	h := NewHttpHandler()
	route.POST("/v1/receiver/packet", h.ReceivePacket)
}
