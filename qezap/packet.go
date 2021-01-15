package qezap

import (
	"sync"

	"github.com/huzhongqing/qelog/api/receiverpb"
	apitypes "github.com/huzhongqing/qelog/api/types"
)

var (
	_maxPacketSize = 4 << 10
	_module        = ""

	_packetPool = sync.Pool{New: func() interface{} {
		return &packet{
			p: &receiverpb.Packet{Module: _module, Data: make([]byte, 0, 1024)},
		}
	}}
)

func setMaxPacketSizeAndModule(size int, module string) {
	_maxPacketSize = size
	_module = module
}

type packet struct {
	p      *receiverpb.Packet
	isFree bool
}

func (p *packet) free() {
	if cap(p.p.Data) > 2*_maxPacketSize {
		// 如果扩容的太大，就让 GC 回收
		return
	}

	p.p.Data = p.p.Data[:0]
	p.p.Id = ""
	p.isFree = false
	_packetPool.Put(p)
}

func (p *packet) append(b []byte) *packet {
	p.p.Data = append(p.p.Data, b...)
	if _maxPacketSize <= len(p.p.Data) {
		p.p.Id = id()
		p.isFree = true
	}
	return p
}

func (p *packet) flush() *packet {
	p.p.Id = id()
	p.isFree = true
	return p
}

func newPacket() *packet {
	p := _packetPool.Get().(*packet)
	p.p.Id = ""
	p.isFree = false
	return p
}

func id() string {
	return apitypes.NewTraceID().Hex()
}
