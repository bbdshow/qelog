package kit

import "encoding/json"

// 用于测试打印
func BeautifyToJSON(v interface{}) string {
	byt, _ := json.MarshalIndent(v, "", "  ")
	return string(byt)
}
