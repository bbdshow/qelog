package entity

type DataPacket struct {
	ID      string   `json:"id"`
	Module  string   `json:"module"`
	Data    []string `json:"data"`
	IsRetry bool     `json:"is_retry"`
}
