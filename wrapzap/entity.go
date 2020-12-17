package wrapzap

import (
	"encoding/json"
	"strconv"
	"time"
)

type DataPacket struct {
	ID      string   `json:"id"`
	Module  string   `json:"module"`
	Data    []string `json:"data"`
	IsRetry bool     `json:"is_retry"`
}

func (dp DataPacket) Marshal() []byte {
	b, _ := json.Marshal(dp)
	return b
}

func NewDataPacket(module string, data []string) DataPacket {
	return DataPacket{
		Module:  module,
		ID:      strconv.FormatInt(time.Now().UnixNano(), 10),
		Data:    data,
		IsRetry: false,
	}
}
