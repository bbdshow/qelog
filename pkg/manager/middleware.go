package manager

import (
	"bytes"
	"encoding/json"
	"io/ioutil"

	"github.com/huzhongqing/qelog/libs/logs"
	"go.uber.org/zap"

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

func HandlerRegisterTraceID() gin.HandlerFunc {
	return func(c *gin.Context) {
		httputil.WithTraceID(c)
		c.Next()
	}
}

type loggingRespWriter struct {
	gin.ResponseWriter
	body *bytes.Buffer
}

func (lrw *loggingRespWriter) Write(b []byte) (int, error) {
	lrw.body.Write(b)
	return lrw.ResponseWriter.Write(b)
}

func HandlerLogging(enable bool) gin.HandlerFunc {
	return func(c *gin.Context) {
		if !enable {
			c.Next()
			return
		}

		// 记录请求
		ip := c.ClientIP()
		path := c.Request.URL.Path
		uri := c.Request.URL.RequestURI()
		request := ""
		if c.Request.Body != nil {
			b, err := ioutil.ReadAll(c.Request.Body)
			if err != nil {
				httputil.RespError(c, err)
				c.Abort()
				return
			}
			// 这里已经读取第一次，就关掉
			_ = c.Request.Body.Close()
			body := bytes.NewBuffer(b)
			request = body.String()
			c.Request.Body = ioutil.NopCloser(body)
		}
		logs.Qezap.InfoWithCtx(c.Request.Context(), "Request", zap.String("reqBody", request), logs.Qezap.ConditionOne(ip), logs.Qezap.ConditionTwo(path), logs.Qezap.ConditionThree(uri))

		lrw := &loggingRespWriter{body: bytes.NewBuffer([]byte{}), ResponseWriter: c.Writer}
		c.Writer = lrw

		c.Next()
		baseResp := httputil.BaseResp{}
		if err := json.Unmarshal(lrw.body.Bytes(), &baseResp); err != nil || baseResp.Code != 0 {
			logs.Qezap.ErrorWithCtx(c.Request.Context(), "Response", zap.String("respBody", lrw.body.String()), logs.Qezap.ConditionOne(ip), logs.Qezap.ConditionTwo(path), logs.Qezap.ConditionThree(uri))
		}
	}
}
