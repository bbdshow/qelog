package http

import (
	"github.com/bbdshow/bkit/ginutil"
	"github.com/bbdshow/qelog/api"
	"github.com/gin-gonic/gin"
)

func receiverPacket(c *gin.Context) {
	in := &api.JSONPacket{}
	if err := ginutil.ShouldBind(c, in); err != nil {
		ginutil.RespErr(c, err)
		return
	}

	if err := receiverSvc.JSONPacketToLogger(c.Request.Context(), c.ClientIP(), in); err != nil {
		ginutil.RespErr(c, err)
		return
	}
	ginutil.RespSuccess(c)
}
