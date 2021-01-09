package manager

import (
	"github.com/gin-gonic/gin"
	"github.com/huzhongqing/qelog/libs/jwt"
	"github.com/huzhongqing/qelog/pkg/httputil"
)

func init() {
	jwt.SetSigningKey("qelog_jwt_signg_key")
}

func AuthVerify(enable bool) gin.HandlerFunc {
	return func(c *gin.Context) {
		if !enable {
			c.Next()
			return
		}
		token := c.GetHeader("X-QELOG-Token")
		if token == "" {
			httputil.RespError(c, httputil.ErrUnauthorized)
			c.Abort()
			return
		}

		ok, err := jwt.VerifyJWTToken(token)
		if err != nil || !ok {
			httputil.RespError(c, httputil.ErrUnauthorized)
			c.Abort()
			return
		}

		c.Next()
	}
}
