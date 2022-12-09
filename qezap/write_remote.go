package qezap

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"sync"
	"sync/atomic"
	"time"

	"go.uber.org/multierr"

	"github.com/bbdshow/qelog/api/receiverpb"
)

type WriteRemote struct {
	once  sync.Once
	mutex sync.Mutex

	opt *remoteOption

	pusher Pusher

	packet *Packet

	// back up file writer
	bw *BackupWrite

	isClose int32
	exit    chan struct{}
}

// NewWriteRemote  impl remote io write, used of zap writer
func NewWriteRemote(opt *remoteOption) *WriteRemote {
	w := &WriteRemote{
		opt:  opt,
		exit: make(chan struct{}),
	}
	w.bw = NewBackupWrite(w.opt.BackupFilename)
	w.packet = NewPacket(w.opt.ModuleName, w.opt.MaxPacketSize)

	w.once.Do(func() {
		go w.initPusher()
		go w.bgTimeArrivalSendPacket()
		go w.bgRetrySendPacket()
	})
	return w
}

func (w *WriteRemote) Write(b []byte) (n int, err error) {
	if atomic.LoadInt32(&w.isClose) == 1 {
		return 0, fmt.Errorf("write remote is close")
	}
	w.mutex.Lock()
	data := w.packet.Append(b)
	if data.CanPush() {
		w.push(data)
	}
	w.mutex.Unlock()

	return len(b), nil
}

func (w *WriteRemote) push(data *DataPacket) {
	defer func() {
		// data have been process
		w.packet.SwitchNextDataPacket()
	}()
	if data.IsEmpty() {
		w.packet.PoolPutDataPacket(data)
		return
	}
	// if pusher exception, write to back up file
	if w.pusher == nil || w.pusher.Concurrent() >= w.opt.MaxConcurrent {
		_ = w.backup(data.Data())
		w.packet.PoolPutDataPacket(data)
		return
	}

	go func() {
		ctx := context.Background()
		if w.opt.WriteTimeout > 0 {
			var cancel context.CancelFunc
			ctx, cancel = context.WithTimeout(context.Background(), w.opt.WriteTimeout)
			defer cancel()
		}

		if err := w.pusher.PushPacket(ctx, data.Data()); err != nil {
			if err == ErrUnavailable {
				// 只有当服务不可用时，放入错误备份文件里
				_ = w.backup(data.Data())
			}
			log.Printf("Writer:remote push packet %s\n", err.Error())
		}
		w.packet.PoolPutDataPacket(data)
	}()
}

func (w *WriteRemote) initPusher() {
	initPusher := func(trans Transport, addrs []string, c int) (Pusher, error) {
		switch trans {
		case TransportHTTP:
			return NewHttpPush(addrs, c)
		case TransportGRPC:
			return NewGRPCPush(addrs, c)
		default:
			return nil, fmt.Errorf("init %s transport pusher invalid", trans)
		}
	}

	tick := time.NewTicker(time.Second)
	for {
		if w.pusher != nil {
			return
		}
		var err error
		w.pusher, err = initPusher(w.opt.Transport, w.opt.Addrs, w.opt.MaxConcurrent)
		if err != nil {
			log.Printf("Writer:init %s pusher %v\n", w.opt.Transport, err)
			<-tick.C
			continue
		}

		log.Printf("init %s pusher success \n", w.opt.Transport)
		tick.Stop()
		return
	}
}

type bakPacket struct {
	ID     string `json:"id"`
	Module string `json:"module"`
	Data   string `json:"data"`
}

func (w *WriteRemote) backup(in *receiverpb.Packet) error {
	bp := bakPacket{
		ID:     in.Id,
		Module: in.Module,
		Data:   string(in.Data),
	}

	byt, err := json.Marshal(bp)
	if err != nil {
		return err
	}
	_, err = w.bw.WriteBakPacket(byt)
	return err
}

// when interval time arrival, even if then packet not full, it also should send
func (w *WriteRemote) bgTimeArrivalSendPacket() {
	tick := time.NewTicker(time.Second)
	for {
		select {
		case <-tick.C:
			w.mutex.Lock()
			if d := w.packet.DataPacket(); d != nil {
				w.packet.SetCanPush(d)
				w.push(d)
			}
			w.mutex.Unlock()
		case <-w.exit:
			tick.Stop()
			return
		}
	}
}

// when packet send failed, interval retry data from back up file
func (w *WriteRemote) bgRetrySendPacket() {
	tick := time.NewTicker(200 * time.Millisecond)
	for {
		select {
		case <-tick.C:
			byt, err := w.bw.ReadBakPacket()
			if err != nil {
				log.Println("Writer:packet retry", err.Error())
			}
			if len(byt) > 0 {
				bp := &bakPacket{}
				if err := json.Unmarshal(byt, bp); err != nil {
					log.Println("Writer:packet retry Unmarshal", err.Error())
					break
				}
				v := &receiverpb.Packet{
					Id:     bp.ID,
					Module: bp.Module,
					Data:   []byte(bp.Data),
				}
			loop:
				// read data must send success, because back up file have been offset.
				// warning: if main process exception, back up file have content, packets are sent repeatedly, Packet.ID used of idempotent.
				if w.pusher != nil {
					if err := w.pusher.PushPacket(context.Background(), v); err == nil {
						break
					}
					log.Printf("Writer:packet retry push %v\n", err)
				}
				time.Sleep(time.Second)
				goto loop
			}
		case <-w.exit:
			return
		}
	}
}

// Sync write final content
func (w *WriteRemote) Sync() error {
	w.mutex.Lock()
	defer w.mutex.Unlock()

	// send last packet
	if d := w.packet.DataPacket(); d != nil {
		w.packet.SetCanPush(d)
		w.push(d)
	}
	// send pusher memory data
	wait := make(chan struct{}, 1)
	go func() {
		for {
			if w.pusher != nil && w.pusher.Concurrent() == 0 {
				time.Sleep(10 * time.Millisecond)
				wait <- struct{}{}
				return
			}
			time.Sleep(50 * time.Millisecond)
		}
	}()

	select {
	case <-wait:
	}

	return w.bw.Close()
}

func (w *WriteRemote) Close() error {
	atomic.StoreInt32(&w.isClose, 1)

	var err error
	err = multierr.Append(err, w.Sync())
	if w.pusher != nil {
		err = multierr.Append(err, w.pusher.Close())
	}

	close(w.exit)
	return err
}
