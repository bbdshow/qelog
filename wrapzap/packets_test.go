package wrapzap

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"testing"
	"time"
)

type jsonMessage struct {
	Time  string `json:"time"`
	Level string `json:"level"`
	Field string `json:"field"`
}

func (m jsonMessage) String() string {
	b, _ := json.Marshal(m)
	return string(b)
}

func TestPackets_PutFailedPacket(t *testing.T) {
	p := NewPackets(1024)
	n, err := p.PutFailedPacket(FailedPacket{
		ID: "2",
		Data: jsonMessage{
			Time:  time.Now().String(),
			Level: "INFO",
			Field: "Backup",
		}.String(),
	})

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
		fp := FailedPacket{}
		ok, err := p.GetFailedPacket(&fp)
		if err != nil {
			t.Fatal(err)
		}
		if ok {
			fmt.Println(fp)
		}
	}
}

func TestPackets_AddPacket(t *testing.T) {
	p := NewPackets(1024)
	buf, flush := p.AddPacket([]byte(RandString(512)))
	if flush {
		t.Fatal("flush should false")
	}
	if len(buf) > 0 {
		t.Fatal("buf should nil")
	}

	time.Sleep(4 * time.Second)

	buf, flush = p.AddPacket([]byte(RandString(246)))
	if !flush {
		t.Fatal("flush should hash true")
	}
	if len(buf) < 0 {
		t.Fatal("buf should has value")
	}
	buf, flush = p.AddPacket([]byte(RandString(256)))
	if flush {
		t.Fatal("flush should false")
	}
	if len(buf) > 0 {
		t.Fatal("buf should nil")
	}

	buf, flush = p.AddPacket([]byte(RandString(1024)))

	if !flush {
		t.Fatal("flush should hash true")
	}
	if len(buf) < 0 {
		t.Fatal("buf should has value")
	}
}
