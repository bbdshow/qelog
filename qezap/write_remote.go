package qezap

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"sync"
	"time"

	"go.uber.org/multierr"

	"github.com/bbdshow/qelog/api/receiverpb"
)

type WriteRemote struct {
	mutex sync.Mutex
	cfg   *Config

	pusher Pusher
	packet *packet

	bw   *BackupWrite
	once sync.Once
}

func NewWriteRemote(cfg *Config) *WriteRemote {

	w := &WriteRemote{
		cfg: cfg,
		bw:  NewBackupWrite(cfg.BackupFilename),
	}
	setMaxPacketSizeAndModule(cfg.MaxPacketSize, cfg.ModuleName)

	w.once.Do(func() {
		go w.initPusher()
		go w.backgroundSendPacket()
		go w.backgroundRetrySendPacket()
	})
	return w
}

func (w *WriteRemote) Write(b []byte) (n int, err error) {
	w.mutex.Lock()
	if w.packet == nil {
		w.packet = newPacket()
	}
	p := w.packet.append(b)
	if p.isFree {
		w.push(p)
	}
	w.mutex.Unlock()

	return len(b), nil
}

func (w *WriteRemote) push(p *packet) {
	defer func() {
		w.packet = nil
	}()
	if len(p.p.Data) <= 0 {
		// 没有类容的包，直接丢掉
		return
	}
	if w.pusher == nil {
		_ = w.backup(p.p)
		p.free()
		return
	}
	// 如果发送者满负荷，则直接丢文件
	if w.pusher.Concurrent() >= w.cfg.MaxConcurrent {
		_ = w.backup(p.p)
		p.free()
		return
	}

	ctx := context.Background()
	if w.cfg.WriteTimeout > 0 {
		ctx, _ = context.WithTimeout(context.Background(), w.cfg.WriteTimeout)
	}

	go func() {
		if err := w.pusher.PushPacket(ctx, p.p); err != nil {
			if err == ErrUnavailable {
				// 只有当服务不可用时，放入错误备份文件里
				_ = w.backup(p.p)
			}
			log.Printf("write remote push packet %s\n", err.Error())
		}
		p.free()
	}()
}

func (w *WriteRemote) initPusher() {
	// 在发送的时候，才去链接， 如果链接不通，不能影响主进程
	tick := time.NewTicker(time.Second)
	for {
		if w.pusher == nil {
			if w.cfg.Transport == "http" {
				pusher, err := NewHttpPush(w.cfg.Addrs[0], w.cfg.MaxConcurrent)
				if err != nil {
					log.Printf("init http push error %s\n", err.Error())
					goto next
				}
				w.pusher = pusher
			} else {
				pusher, err := NewGRPCPush(w.cfg.Addrs, w.cfg.MaxConcurrent)
				if err != nil {
					log.Printf("init grpc push error %s\n", err.Error())
					goto next
				}
				w.pusher = pusher

			}
			log.Printf("init %s push success \n", w.cfg.Transport)
			tick.Stop()
			return
		}
	next:
		select {
		case <-tick.C:
		}
	}
}

type _jsonPacket struct {
	ID     string `json:"id"`
	Module string `json:"module"`
	Data   string `json:"data"`
}

func (w *WriteRemote) backup(in *receiverpb.Packet) error {
	jsonPacket := _jsonPacket{
		ID:     in.Id,
		Module: in.Module,
		Data:   string(in.Data),
	}

	byt, err := json.Marshal(jsonPacket)
	if err != nil {
		return err
	}
	_, err = w.bw.WriteBakPacket(byt)
	return err
}

// 当一定时间内，包容量没有达到，则也会默认发送已在缓存中的日志
func (w *WriteRemote) backgroundSendPacket() {
	tick := time.NewTicker(time.Second)
	for {
		select {
		case <-tick.C:
			w.mutex.Lock()
			if w.packet != nil {
				p := w.packet.flush()
				w.push(p)
			}
			w.mutex.Unlock()
		}
	}
}

func (w *WriteRemote) backgroundRetrySendPacket() {
	tick := time.NewTicker(200 * time.Millisecond)
	for {
		select {
		case <-tick.C:
			byt, err := w.bw.ReadBakPacket()
			if err != nil {
				fmt.Println("packets retry", err.Error())
			}
			if len(byt) > 0 {
				jsonPacket := &_jsonPacket{}
				if err := json.Unmarshal(byt, jsonPacket); err != nil {
					fmt.Println("packets retry", err.Error())
					break
				}
				v := &receiverpb.Packet{
					Id:     jsonPacket.ID,
					Module: jsonPacket.Module,
					Data:   []byte(jsonPacket.Data),
				}
			loop:
				if w.pusher != nil {
					ctx, _ := context.WithTimeout(context.Background(), w.cfg.WriteTimeout)
					if err := w.pusher.PushPacket(ctx, v); err == nil {
						break
					} else {
						log.Printf("write remote push packet %s\n", err.Error())
					}
				}
				time.Sleep(time.Second)
				goto loop
			}
		}
	}
}

func (w *WriteRemote) Sync() error {
	w.mutex.Lock()
	defer w.mutex.Unlock()
	if w.packet != nil {
		p := w.packet.flush()
		w.push(p)
	}
	sendEmpty := make(chan struct{}, 1)
	go func() {
		for {
			if w.pusher != nil && w.pusher.Concurrent() == 0 {
				time.Sleep(10 * time.Millisecond)
				sendEmpty <- struct{}{}
				return
			}
			time.Sleep(50 * time.Millisecond)
		}
	}()

	ctx, _ := context.WithTimeout(context.Background(), w.cfg.WriteTimeout)
	select {
	case <-ctx.Done():
		log.Println("sync ", ctx.Err())
	case <-sendEmpty:
	}

	return w.bw.Close()
}

func (w *WriteRemote) Close() error {
	var err error
	err = multierr.Append(err, w.Sync())
	if w.pusher != nil {
		err = multierr.Append(err, w.pusher.Close())
	}
	return err
}
