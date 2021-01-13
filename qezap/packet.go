package qezap

import (
	"fmt"
	"math/rand"
	"sync"
	"time"

	"github.com/huzhongqing/qelog/pb"
)

type Packet struct {
	mutex    sync.Mutex
	pool     sync.Pool
	packet   *pb.Packet
	maxSize  int
	capLimit int
}

func NewPacket(module string, maxSize int) *Packet {
	p := &Packet{
		packet:  nil,
		maxSize: maxSize,
		// 预留一定容量，避免最后一条信息写入，导致扩容
		capLimit: maxSize + 128,
	}
	p.pool = sync.Pool{New: func() interface{} {
		return &pb.Packet{Module: module, Data: make([]byte, 0, p.capLimit)}
	}}

	return p
}

func (p *Packet) initPacket() {
	v := p.pool.Get().(*pb.Packet)
	p.packet = v
}
func (p *Packet) FreePacket(v *pb.Packet) {
	if cap(v.Data) > p.capLimit {
		// 如果产生了扩容，则不放回池子里面，让GC回收
		return
	}

	v.Id = ""
	v.Data = v.Data[:0]

	p.pool.Put(v)
}

func (p *Packet) AppendData(b []byte) *pb.Packet {
	p.mutex.Lock()
	if p.packet == nil {
		p.initPacket()
	}
	p.packet.Data = append(p.packet.Data, b...)
	if p.maxSize <= len(p.packet.Data) {
		v := p.packet
		v.Id = p.packetID()
		p.packet = nil
		p.mutex.Unlock()
		return v
	}
	p.mutex.Unlock()
	return nil
}

func (p *Packet) FlushData() *pb.Packet {
	p.mutex.Lock()
	if p.packet == nil {
		p.mutex.Unlock()
		return nil
	}
	v := p.packet
	v.Id = p.packetID()

	p.packet = nil
	p.mutex.Unlock()
	return v
}

var incNum int64 = 0

func (p *Packet) packetID() string {
	incNum++
	if incNum >= 10000 {
		incNum = 0
	}
	return fmt.Sprintf("%d_%02d_%04d", time.Now().UnixNano()/1e6, rand.Int31n(100), incNum)
}
