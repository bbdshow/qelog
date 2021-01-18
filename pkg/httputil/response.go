package httputil

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/url"

	apitypes "github.com/huzhongqing/qelog/api/types"

	"github.com/gin-gonic/gin"
)

func WithTraceID(c *gin.Context) {
	id := apitypes.NewTraceID()
	ctx := context.WithValue(c.Request.Context(), apitypes.EncoderTraceIDKey, id)
	c.Request = c.Request.WithContext(ctx)
}

type BaseResp struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	TraceID string `json:"traceId"`
}

func NewBaseResp(code int, msg string) *BaseResp {
	return &BaseResp{
		Code:    code,
		Message: msg,
		TraceID: "",
	}
}

func (br *BaseResp) WriteTraceID(c *gin.Context) *BaseResp {
	v := c.Request.Context().Value(apitypes.EncoderTraceIDKey)
	if v != nil {
		id, ok := v.(apitypes.TraceID)
		if ok {
			br.TraceID = id.Hex()
		}
	}
	return br
}

type DataResp struct {
	*BaseResp
	Data interface{} `json:"data"`
}

func (dr *DataResp) WriteTraceID(c *gin.Context) *DataResp {
	dr.BaseResp.WriteTraceID(c)
	return dr
}

func NewDataResp(code int, msg string, data interface{}) *DataResp {
	return &DataResp{
		BaseResp: NewBaseResp(code, msg),
		Data:     data,
	}
}

func RespError(c *gin.Context, err error) {
	// 请求算成功，取业务Code码
	RespStatusCodeWithError(c, http.StatusOK, err)
}

type ResponseErr struct {
	Method   string     `json:"method"`
	Path     string     `json:"path"`
	Form     url.Values `json:"form"`
	PostForm url.Values `json:"postForm"`
	Func     string     `json:"func"`
	Error    string     `json:"error"`
}

func (err ResponseErr) Marshal() ([]byte, error) {
	return json.Marshal(err)
}

func (err ResponseErr) String() string {
	byt, _ := err.Marshal()
	return string(byt)
}

func RespStatusCodeWithError(c *gin.Context, statusCode int, err error) {
	if err == nil {
		err = errors.New("nil error")
	}
	code := CodeFailed
	message := ""
	if e, ok := err.(Error); ok {
		// 如果是自定义错误，就重写 code
		code = e.Code
		message = e.Message
	}
	switch code {
	case ErrCodeSystemException:
		// 拦截响应中间件已经打日志
		//if gin.DefaultErrorWriter != nil {
		//	respErr := ResponseErr{
		//		Method:   c.Request.Method,
		//		Path:     c.Request.URL.RequestURI(),
		//		Form:     c.Request.Form,
		//		PostForm: c.Request.PostForm,
		//		Func:     c.HandlerName(),
		//		Error:    err.Error(),
		//	}
		//	byt, _ := respErr.Marshal()
		//	_, _ = gin.DefaultErrorWriter.Write(byt)
		//}
	}

	c.JSON(statusCode, NewBaseResp(code, message).WriteTraceID(c))
}

func RespData(c *gin.Context, httpCode int, ret interface{}) {
	if ret == nil {
		c.JSON(httpCode, NewBaseResp(0, "success").WriteTraceID(c))
		return
	}
	c.JSON(httpCode, NewDataResp(0, "success", ret).WriteTraceID(c))
}

func RespSuccess(c *gin.Context) {
	RespData(c, http.StatusOK, nil)
}
