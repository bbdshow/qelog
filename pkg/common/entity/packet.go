package entity

type JSONPacket struct {
	Id     string   `json:"id"`
	Module string   `json:"module"`
	Data   []string `json:"data"`
}
