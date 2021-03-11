package httputil

import (
	"errors"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/huzhongqing/qelog/infra/logs"
	"go.uber.org/zap"
)

var HideValidatorErr = true

func SetHideValidatorErr(b bool) {
	HideValidatorErr = b
}

func ShouldBind(c *gin.Context, obj interface{}) error {
	if err := c.ShouldBind(obj); err != nil {
		if HideValidatorErr {
			logs.Qezap.Info("ShouldBind", zap.String("obj", err.Error()), logs.Qezap.FieldTraceID(c.Request.Context()))
			return errors.New("validator error")
		}
	}
	return nil
}

func ValidateStruct(obj interface{}) error {
	return binding.Validator.ValidateStruct(obj)
}

type Language string

func (l Language) IsChinese() bool {
	return l == CN
}
func (l Language) IsEnglish() bool {
	return l == EN
}

const (
	CN Language = "zh-cn"
	EN Language = "en"
)

func GetLanguage(c *gin.Context) Language {
	v := c.GetHeader("Language")
	switch strings.ToLower(v) {
	case string(EN):
		return EN
	default:
		return CN
	}
}
