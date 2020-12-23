package qzzap

import (
	"encoding/json"
	"math/rand"
	"time"
)

// 测试消息
type _jsonMessage struct {
	Time  string `json:"time"`
	Level string `json:"level"`
	Field string `json:"field"`
}

func (m _jsonMessage) Marshal() []byte {
	b, _ := json.Marshal(m)
	return b
}
func (m _jsonMessage) String() string {
	return string(m.Marshal())
}

func _RandString(length int) string {
	baseChar := "0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
	byteChar := []byte(baseChar)
	str := ""
	for i := 0; i < length; i++ {
		rand.New(rand.NewSource(time.Now().UnixNano() + rand.Int63n(1000000)))
		str += string(byteChar[rand.Intn(len(byteChar))])
	}
	return str
}
