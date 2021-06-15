package httputil

import (
	"fmt"

	"github.com/bbdshow/qelog/infra/jwt"
	"github.com/bbdshow/qelog/infra/logs"

	"github.com/gin-gonic/gin"

	"go.uber.org/zap"
)

var Authorization = "X-Authorization"

// signingKey 自定义的 加密Key, 如果没有就使用全局默认的
func AuthAdmin(enable bool, signingKey ...string) gin.HandlerFunc {
	return func(c *gin.Context) {
		if !enable {
			c.Next()
			return
		}
		token := c.GetHeader(Authorization)
		if token == "" {
			logs.Qezap.Info(fmt.Sprintf("%s header required", Authorization), logs.Qezap.FieldTraceID(c.Request.Context()))
			RespError(c, ErrUnauthorized)
			c.Abort()
			return
		}

		ok, err := jwt.VerifyJWTToken(token, signingKey...)
		if err != nil || !ok {
			logs.Qezap.Info(fmt.Sprintf("%s token verify", Authorization), zap.Error(err), logs.Qezap.FieldTraceID(c.Request.Context()))
			RespError(c, ErrUnauthorized)
			c.Abort()
			return
		}

		if err := SetJWTClaims(c, token, signingKey...); err != nil {
			logs.Qezap.Info(fmt.Sprintf("%s token set claims", Authorization), zap.Error(err), logs.Qezap.FieldTraceID(c.Request.Context()))
			RespError(c, ErrUnauthorized)
			c.Abort()
			return
		}
		c.Next()
	}
}
