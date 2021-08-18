package types

import (
	apitypes "github.com/bbdshow/qelog/api/types"
	jsoniterator "github.com/json-iterator/go"
	"time"
)

type Decoder map[string]interface{}

func (dec Decoder) Level() Level {
	interfaceV, ok := dec[apitypes.EncoderLevelKey]
	if ok {
		val, ok1 := interfaceV.(string)
		if ok1 {
			return String2Level(val)
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
