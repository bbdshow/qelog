package wrapzap

import (
	"testing"
	"time"
)

func TestWriteRemote_Write(t *testing.T) {
	wr := NewWriteRemote(NewWriteRemoteConfig("http://127.0.0.1:31081/v1/receive/packet", "test"))
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
