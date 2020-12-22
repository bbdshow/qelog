package auth

import "github.com/gin-gonic/gin"

type Receiver struct {
}

func (rec *Receiver) Verify(c *gin.Context) {}
