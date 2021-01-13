package qezap

import (
	"fmt"
	"math/rand"
	"sync"
	"time"

	"github.com/huzhongqing/qelog/pb"
)

var (
	_maxPacketSize = 4 << 10
	_module        = ""

	_packetPool = sync.Pool{New: func() interface{} {
		return &packet{
			p: &pb.Packet{Module: _module, Data: make([]byte, 0, 1024)},
		}
	}}
)

func setMaxPacketSizeAndModule(size int, module string) {
	_maxPacketSize = size
	_module = module
}

type packet struct {
	p      *pb.Packet
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
	p.isFree = true
	return p
}

func newPacket() *packet {
	p := _packetPool.Get().(*packet)
	return p
}

var incNum int64 = 0

func id() string {
	incNum++
	if incNum >= 10000 {
		incNum = 0
	}
	return fmt.Sprintf("%d_%02d_%04d", time.Now().UnixNano()/1e6, rand.Int31n(100), incNum)
}
