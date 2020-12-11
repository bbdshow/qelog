package entity

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/huzhongqing/qelog/types/errors"
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
	RespStatusCodeWithError(c, http.StatusBadRequest, err)
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
	code := errors.CodeFailed
	message := ""

	if e, ok := err.(errors.Error); ok {
		// 如果是自定义错误，就重写 code
		code = e.Code
		message = e.Message
	}
	switch code {
	case errors.ErrCodeSystemException:
		e := errLog{
			Method:      c.Request.Method,
			Path:        c.Request.URL.RequestURI(),
			HandlerFunc: c.HandlerName(),
			Err:         err.Error(),
		}
		fmt.Println(e)
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
