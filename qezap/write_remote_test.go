package qezap

import (
	"strconv"
	"testing"
	"time"
)

func TestWriteRemote_Write(t *testing.T) {
	addrs := []string{"http://127.0.0.1:31081/v1/receive/packet"}
	wr := NewWriteRemote(NewWriteRemoteConfig(addrs, "test"))
	msg := _jsonMessage{
		Time:  time.Now().String(),
		Level: "INFO",
		Field: "哈哈\n嘿嘿",
	}
	_, err := wr.Write(msg.Marshal())
	if err != nil {
		t.Fatal(err)
	}
	time.Sleep(time.Minute)
}

func TestWriteRemote_Grpc(t *testing.T) {
	addrs := []string{"127.0.0.1:31082", "127.0.0.1:31182"}
	cfg := NewWriteRemoteConfig(addrs, "example")
	wr := NewWriteRemote(cfg)
	count := 100
	for count >= 0 {
		count--
		msg := _jsonMessage{
			Time:  time.Now().String(),
			Level: "INFO",
			Field: strconv.Itoa(count),
		}
		_, err := wr.Write(msg.Marshal())
		if err != nil {
			t.Fatal(err)
		}
		time.Sleep(time.Second)
	}

	time.Sleep(time.Minute)
}
