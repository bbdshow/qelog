package receiver

import (
	"github.com/huzhongqing/qelog/model/mongoclient"

	"github.com/huzhongqing/qelog/service/receiver"

	"github.com/gin-gonic/gin"
	"github.com/huzhongqing/qelog/types/entity"
	"github.com/huzhongqing/qelog/types/errors"
)

type HandleFunc struct {
	srv *receiver.Service
}

func NewHandleFunc(database *mongoclient.Database) *HandleFunc {
	h := &HandleFunc{
		srv: receiver.NewService(database),
	}

	return h
}

func (h *HandleFunc) ReceivePacket(c *gin.Context) {
	var arg entity.DataPacket
	if err := c.ShouldBind(&arg); err != nil {
		entity.RespError(c, errors.ErrArgsInvalid)
		return
	}

}
