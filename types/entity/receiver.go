package entity

type DataPacket struct {
	Name string   `json:"name"`
	ID   string   `json:"id"`
	Data []string `json:"data"`
}
