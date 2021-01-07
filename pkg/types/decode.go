package types

import (
	"fmt"
	"strings"
	"time"

	"github.com/huzhongqing/qelog/pkg/common/model"
	jsoniterator "github.com/json-iterator/go"
)

type Decoder struct {
	Val map[string]interface{}
}

func (dec Decoder) Level() model.Level {
	interfaceV, ok := dec.Val["_level"]
	if ok {
		val, ok1 := interfaceV.(string)
		if ok1 {
			return LevelStr2Int(val)
		}
	}
	return 0
}
func (dec Decoder) Time() int64 {
	interfaceV, ok := dec.Val["_time"]
	if ok {
		val, ok1 := interfaceV.(float64)
		if ok1 {
			return int64(val)
		}
	}
	return time.Now().UnixNano() / 1e6
}

func (dec Decoder) Short() string {
	interfaceV, ok := dec.Val["_short"]
	if ok {
		val, ok1 := interfaceV.(string)
		if ok1 {
			return val
		}
	}
	return ""
}

func (dec Decoder) Condition(num int) string {
	interfaceV, ok := dec.Val[fmt.Sprintf("_condition%d", num)]
	if ok {
		val, ok1 := interfaceV.(string)
		if ok1 {
			return val
		}
	}
	return ""
}

func (dec Decoder) TraceID() string {
	interfaceV, ok := dec.Val["_traceid"]
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
	delFields := []string{"_level", "_time", "_short",
		"_condition1", "_condition2", "_condition3", "_traceid"}
	for _, v := range delFields {
		delete(dec.Val, v)
	}
	b, _ := Marshal(dec.Val)
	return string(b)
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
