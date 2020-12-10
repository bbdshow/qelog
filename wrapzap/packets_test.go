package wrapzap

import (
	"fmt"
	"io/ioutil"
	"testing"
	"time"
)

func TestPackets_WritePacket(t *testing.T) {
	p := NewPackets(1024)
	n, err := p.WritePacket(DataPacket{
		ID: "2",
		Data: []string{_jsonMessage{
			Time:  time.Now().String(),
			Level: "INFO",
			Field: "Bac\nkup",
		}.String()},
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
		dp := DataPacket{}
		ok, err := p.ReadPacket(&dp)
		if err != nil {
			t.Fatal(err)
		}
		if ok {
			fmt.Print(dp)
		}
	}
}

func TestPackets_AddPacket(t *testing.T) {
	p := NewPackets(1024)
	buf, flush := p.AddPacket([]byte(RandString(246)))
	if !flush {
		t.Fatal("flush should hash true")
	}
	if len(buf) != 1 {
		t.Fatal("buf should has value")
	}

	buf, flush = p.AddPacket([]byte(RandString(512)))
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
	fb := NewDataPacket(RandString(16), []string{"测试换行\n符的影响"})
	buf, flush = p.AddPacket(fb.Marshal())
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
	for _, b := range buf {
		fmt.Println(string(b))
	}
}
