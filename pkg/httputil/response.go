package httputil

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"

	"github.com/gin-gonic/gin"
)

type BaseResp struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

func Success() BaseResp {
	return BaseResp{
		Code:    0,
		Message: "success",
	}
}

func SuccessWithData(data interface{}) DataResp {
	return DataResp{
		BaseResp: Success(),
		Data:     data,
	}
}

type DataResp struct {
	BaseResp
	Data interface{} `json:"data"`
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
		err = fmt.Errorf("error")
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
		if gin.DefaultErrorWriter != nil {
			respErr := ResponseErr{
				Method:   c.Request.Method,
				Path:     c.Request.URL.RequestURI(),
				Form:     c.Request.Form,
				PostForm: c.Request.PostForm,
				Func:     c.HandlerName(),
				Error:    err.Error(),
			}
			byt, _ := respErr.Marshal()
			_, _ = gin.DefaultErrorWriter.Write(byt)
		}
	}

	c.JSON(statusCode, BaseResp{
		Code:    code,
		Message: message,
	})
}

func RespData(c *gin.Context, httpCode int, ret interface{}) {
	if ret == nil {
		c.JSON(httpCode, Success())
		return
	}
	c.JSON(httpCode, SuccessWithData(ret))
}

func RespSuccess(c *gin.Context) {
	RespData(c, http.StatusOK, nil)
}
