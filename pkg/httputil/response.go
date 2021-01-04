package httputil

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

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

type errLog struct {
	Method      string `json:"method"`
	Path        string `json:"path"`
	HandlerFunc string `json:"handler_func"`
	Err         string `json:"err"`
}

func (e errLog) Marshal() []byte {
	b, _ := json.Marshal(e)
	return b
}

func (e errLog) String() string {
	return string(e.Marshal())
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
		e := errLog{
			Method:      c.Request.Method,
			Path:        c.Request.URL.RequestURI(),
			HandlerFunc: c.HandlerName(),
			Err:         err.Error(),
		}

		if gin.DefaultErrorWriter != nil {
			l := log.New(gin.DefaultErrorWriter, "", log.LstdFlags)
			l.Println(e)
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
