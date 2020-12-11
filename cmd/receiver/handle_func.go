package receiver

import (
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/huzhongqing/qelog/types/entity"
	"github.com/huzhongqing/qelog/types/errors"
)

type HandleFunc struct {
}

func NewHandleFunc() *HandleFunc {
	hf := &HandleFunc{}

	return hf
}

func (hf *HandleFunc) ReceivePacket(c *gin.Context) {
	var arg entity.DataPacket
	if err := c.ShouldBind(&arg); err != nil {
		entity.RespError(c, errors.ErrArgsInvalid)
		return
	}
	for _, d := range arg.Data {
		fmt.Println(d)
	}
}
