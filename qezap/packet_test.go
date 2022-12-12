package qezap

import (
	"strconv"
	"testing"
)

func TestNewPacket(t *testing.T) {
	_ = testNewPacket(t, 1024)
}

func testNewPacket(t *testing.T, maxSize int) *Packet {
	module := "test"
	p := newPacket(module, maxSize)
	if p == nil {
		t.Fatal("packet nil")
	}
	if p.maxSize != maxSize {
		t.Fatal("packet maxSize ", p.maxSize)
	}
	return p
}

func TestPacket_Append(t *testing.T) {
	maxSize := 1024
	p := testNewPacket(t, maxSize)

	largeTxt := ""
	for i := 0; i < maxSize; i++ {
		largeTxt += strconv.Itoa(i)
	}
	var testCases = []struct {
		In      string
		Out     string
		CanPush bool
		IsPut   bool
	}{
		{
			In:      "hello",
			Out:     "hello",
			CanPush: false,
		},
		{
			In:      " world",
			Out:     "hello world",
			CanPush: false,
		},
		{
			In:      ", qelog",
			Out:     "hello world, qelog",
			CanPush: false,
		},
		{
			In:      "-qezap",
			Out:     "hello world, qelog-qezap",
			CanPush: false,
			IsPut:   true,
		},
		{
			In:      largeTxt,
			Out:     largeTxt,
			CanPush: true,
			IsPut:   true,
		},
	}

	for i, testCase := range testCases {
		data := p.Append([]byte(testCase.In))
		out := string(data.Data().GetData())
		if out != testCase.Out {
			t.Fatalf("case index %d, out %s", i, out)
		}
		if data.CanPush() != testCase.CanPush {
			t.Fatalf("case index %d, canPush %t", i, data.CanPush())
		}
		if testCase.IsPut {
			p.SwitchNextDataPacket()
			if cap(data.Data().Data) > 2*p.maxSize {
				p.PoolPutDataPacket(data)
				if len(data.Data().Data) <= 0 {
					t.Fatalf("data should be waiting GC")
				}
			} else {
				// clear data
				p.PoolPutDataPacket(data)
				if !data.IsEmpty() {
					t.Fatalf("data should be empty")
				}
			}
			// switch next, should empty
			if !p.DataPacket().IsEmpty() {
				t.Fatalf("switch data should be empty")
			}
		}
	}
}
