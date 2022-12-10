package qezap

import (
	"strconv"
	"testing"
	"time"
)

func testNewWriteRemote(t *testing.T) *WriteRemote {
	opt := defaultRemoteOption()
	opt.Addrs = []string{"127.0.0.1:31082"}
	opt.Transport = TransportMock
	opt.MaxPacketSize = 1024
	w := newWriteRemote(opt)
	if w == nil {
		t.Fatal("write remote nil")
	}
	time.Sleep(time.Second)
	return w
}

func TestNewWriteRemote(t *testing.T) {
	w := testNewWriteRemote(t)
	if w == nil {
		t.Fatal("write remote nil")
	}
	time.Sleep(time.Second)
}

func TestWriter_Write(t *testing.T) {
	w := testNewWriteRemote(t)
	largeTxt := ""
	for i := 0; i < int(w.opt.MaxPacketSize); i++ {
		largeTxt += strconv.Itoa(i)
	}
	var testCases = []struct {
		In    string
		Data  string
		Sleep time.Duration
	}{
		{
			In:    "hello write remote",
			Data:  "hello write remote",
			Sleep: 300 * time.Millisecond,
		},
		{
			In:    "hello write local",
			Data:  "",
			Sleep: 1000 * time.Millisecond,
		},
		{
			In:    largeTxt,
			Data:  "",
			Sleep: 1 * time.Millisecond,
		},
		{
			In:    "next packet",
			Data:  "next packet",
			Sleep: 1 * time.Millisecond,
		},
	}

	for i, v := range testCases {
		n, err := w.Write([]byte(v.In))
		if err != nil {
			t.Fatalf(err.Error())
		}
		if n != len(v.In) {
			t.Fatalf("case %d: input n invalid %d", i, n)
		}

		time.Sleep(v.Sleep)

		if string(w.packet.DataPacket().Data().Data) != v.Data {
			t.Fatalf("case %d: data not equal %s", i, v.Data)
		}
	}
	time.Sleep(time.Second)
}

func TestWriter_Close(t *testing.T) {
	w := testNewWriteRemote(t)
	largeTxt := ""
	for i := 0; i < int(w.opt.MaxPacketSize); i++ {
		largeTxt += strconv.Itoa(i)
	}
	var testCases = []struct {
		In      string
		Data    string
		Sleep   time.Duration
		IsClose bool
		IsErr   bool
	}{
		{
			In:      "hello write remote",
			Data:    "",
			Sleep:   100 * time.Millisecond,
			IsClose: true,
		},
		{
			In:    "hello write local",
			Data:  "",
			Sleep: 1000 * time.Millisecond,
			IsErr: true,
		},
	}

	for i, v := range testCases {
		n, err := w.Write([]byte(v.In))
		if err != nil {
			if !v.IsErr {
				t.Fatalf(err.Error())
			}
			t.Log(err)
			break
		}
		if n != len(v.In) {
			t.Fatalf("case %d: input n invalid %d", i, n)
		}
		if v.IsClose {
			w.Close()
		}

		time.Sleep(v.Sleep)

		if string(w.packet.DataPacket().Data().Data) != v.Data {
			t.Fatalf("case %d: data not equal %s", i, v.Data)
		}
	}
	time.Sleep(time.Second)
}

func TestWriter_RetrySendPacket(t *testing.T) {
	w := testNewWriteRemote(t)

	var testCases = []struct {
		In    string
		Data  string
		Sleep time.Duration
	}{
		{
			In:    ErrUnavailable.Error(),
			Data:  "",
			Sleep: 2000 * time.Millisecond,
		},
		{
			In:    "hello write local",
			Data:  "hello write local",
			Sleep: 300 * time.Millisecond,
		},
	}

	for i, v := range testCases {
		n, err := w.Write([]byte(v.In))
		if err != nil {
			t.Fatalf(err.Error())
		}
		if n != len(v.In) {
			t.Fatalf("case %d: input n invalid %d", i, n)
		}

		time.Sleep(v.Sleep)

		if string(w.packet.DataPacket().Data().Data) != v.Data {
			t.Fatalf("case %d: data not equal %s", i, v.Data)
		}
	}
	time.Sleep(time.Second)
}
