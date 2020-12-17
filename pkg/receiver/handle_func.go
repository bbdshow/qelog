package receiver

import (
	"github.com/gin-gonic/gin"
	"github.com/huzhongqing/qelog/pkg/common/entity"
	"github.com/huzhongqing/qelog/pkg/httputil"
)

func (rec *Receiver) ReceivePacket(c *gin.Context) {

	var arg entity.DataPacket
	if err := c.ShouldBind(&arg); err != nil {
		httputil.RespError(c, httputil.ErrArgsInvalid)
		return
	}

	if err := rec.srv.InsertPacket("", c.ClientIP(), arg); err != nil {
		httputil.RespError(c, httputil.ErrNotFound)
		return
	}
	httputil.RespSuccess(c)
}
