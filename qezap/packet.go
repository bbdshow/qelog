package qezap

import (
	"sync"

	"github.com/bbdshow/qelog/api/receiverpb"
	apitypes "github.com/bbdshow/qelog/api/types"
)

type Packet struct {
	maxSize int
	module  string
	pool    sync.Pool

	data *DataPacket
}

func newPacket(module string, maxSize ...int) *Packet {
	p := &Packet{
		maxSize: 4 << 10, // 4KB
		module:  module,
		pool: sync.Pool{
			New: func() interface{} {
				return &DataPacket{
					p: &receiverpb.Packet{Module: module, Data: make([]byte, 0, 1024)},
				}
			},
		},
	}
	if len(maxSize) > 0 {
		p.maxSize = maxSize[0]
	}
	return p
}

// Append warning: concurrent not safe
func (p *Packet) Append(b []byte) *DataPacket {
	d := p.DataPacket()
	d.p.Data = append(d.p.Data, b...)
	if len(d.p.Data) >= p.maxSize {
		p.SetCanPush(d)
	}
	return d
}

// SwitchNextDataPacket current data have been processed, so set nil, switch next data packet
func (p *Packet) SwitchNextDataPacket() {
	p.data = nil
}

func (p *Packet) SetCanPush(d *DataPacket) {
	d.p.Id = id()
	d.canPush = true
}

func (p *Packet) DataPacket() *DataPacket {
	if p.data == nil {
		p.data = p.pool.Get().(*DataPacket)
	}
	return p.data
}

func (p *Packet) PoolPutDataPacket(d *DataPacket) {
	if cap(d.p.Data) > 2*p.maxSize {
		// cap too large, waiting GC
		return
	}
	// clear data, keep cap
	d.p.Data = d.p.Data[:0]
	d.p.Id = ""
	d.canPush = false
	p.pool.Put(d)
}

type DataPacket struct {
	p *receiverpb.Packet
	// true: data full load or time of arrival, can push
	canPush bool
}

func (d *DataPacket) IsEmpty() bool {
	return len(d.p.Data) <= 0
}

func (d *DataPacket) Data() *receiverpb.Packet {
	return d.p
}

func (d *DataPacket) CanPush() bool {
	return d.canPush
}

func id() string {
	return apitypes.NewTraceID().Hex()
}
