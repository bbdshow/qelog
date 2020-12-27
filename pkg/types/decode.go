package types

import (
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/huzhongqing/qelog/pkg/common/model"
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

// 删除一些不必要的字段，节约存储
func (dec Decoder) Full() string {
	delFields := []string{"_level", "_time", "_short", "_condition1", "_condition2", "_condition3"}
	for _, v := range delFields {
		delete(dec.Val, v)
	}
	b, _ := Marshal(dec.Val)
	return string(b)
}

func LevelStr2Int(lvl string) model.Level {
	switch strings.ToUpper(lvl) {
	case "DEBUG":
		return 0
	case "INFO":
		return 1
	case "WARN":
		return 2
	case "ERROR":
		return 3
	case "PANIC", "DPANIC":
		return 4
	case "FATAL":
		return 5
	}
	return -1
}

func Unmarshal(data []byte, v interface{}) error {
	return json.Unmarshal(data, v)
}

func Marshal(v interface{}) ([]byte, error) {
	return json.Marshal(v)
}
