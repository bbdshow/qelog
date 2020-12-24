package qezap

import (
	"fmt"
	"io/ioutil"
	"testing"
	"time"

	"github.com/huzhongqing/qelog/qezap/push"
)

func TestPackets_WritePacket(t *testing.T) {
	p := NewPackets(1024)
	n, err := p.WriteBakPacket(push.NewPacket("test", []string{_jsonMessage{
		Time:  time.Now().String(),
		Level: "INFO",
		Field: "Bac\nkup",
	}.String()}))

	if b, err := ioutil.ReadFile(p.bakFilename); err != nil {
		t.Fatal(err)
	} else {
		fmt.Print(string(b))
	}
	fmt.Println(n, err)
}

func TestPackets_FailedPacket(t *testing.T) {
	p := NewPackets(1024)
	for i := 0; i < 5; i++ {
		dp := push.Packet{}
		ok, err := p.ReadBakPacket(&dp)
		if err != nil {
			t.Fatal(err)
		}
		if ok {
			fmt.Print(dp)
		}
	}
}
