package types

import (
	"strings"
	"time"

	apitypes "github.com/bbdshow/qelog/api/types"

	"github.com/bbdshow/qelog/pkg/common/model"
	jsoniterator "github.com/json-iterator/go"
)

//type Decoder struct {
//	Val map[string]interface{}
//}

type Decoder map[string]interface{}

func (dec Decoder) Level() model.Level {
	interfaceV, ok := dec[apitypes.EncoderLevelKey]
	if ok {
		val, ok1 := interfaceV.(string)
		if ok1 {
			return LevelStr2Int(val)
		}
	}
	return 0
}
func (dec Decoder) TimeMill() int64 {
	interfaceV, ok := dec[apitypes.EncoderTimeKey]
	if ok {
		val, ok1 := interfaceV.(float64)
		if ok1 {
			return int64(val)
		}
	}
	return time.Now().UnixNano() / 1e6
}

func (dec Decoder) Short() string {
	interfaceV, ok := dec[apitypes.EncoderMessageKey]
	if ok {
		val, ok1 := interfaceV.(string)
		if ok1 {
			return val
		}
	}
	return ""
}

func (dec Decoder) Condition(num int) string {
	key := ""
	switch num {
	case 1:
		key = apitypes.EncoderConditionOneKey
	case 2:
		key = apitypes.EncoderConditionTwoKey
	case 3:
		key = apitypes.EncoderConditionThreeKey
	}
	interfaceV, ok := dec[key]
	if ok {
		val, ok1 := interfaceV.(string)
		if ok1 {
			return val
		}
	}
	return ""
}

func (dec Decoder) TraceIDHex() string {
	interfaceV, ok := dec[apitypes.EncoderTraceIDKey]
	if ok {
		val, ok1 := interfaceV.(string)
		if ok1 {
			return val
		}
	}
	return ""
}

// 删除一些不必要的字段，节约存储
func (dec Decoder) Full() string {
	delFields := []string{apitypes.EncoderLevelKey, apitypes.EncoderTimeKey, apitypes.EncoderMessageKey,
		apitypes.EncoderConditionOneKey, apitypes.EncoderConditionTwoKey, apitypes.EncoderConditionThreeKey, apitypes.EncoderTraceIDKey}
	for _, v := range delFields {
		delete(dec, v)
	}
	str, _ := MarshalToString(dec)
	return str
}

func LevelStr2Int(lvl string) model.Level {
	switch strings.ToUpper(lvl) {
	case "DEBUG":
		return -1
	case "INFO":
		return 0
	case "WARN":
		return 1
	case "ERROR":
		return 2
	case "DPANIC":
		return 3
	case "PANIC":
		return 4
	case "FATAL":
		return 5
	}
	return -2
}

// 频繁调用，快速解析
func Unmarshal(data []byte, v interface{}) error {
	return jsoniterator.Unmarshal(data, v)
}

func Marshal(v interface{}) ([]byte, error) {
	return jsoniterator.Marshal(v)
}

func MarshalToString(v interface{}) (string, error) {
	return jsoniterator.MarshalToString(v)
}
