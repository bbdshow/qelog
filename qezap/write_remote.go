package qezap

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"path"
	"sync"
	"time"

	"github.com/huzhongqing/qelog/pb"
)

var _backupFilename = path.Join(path.Dir(_logFilename), "backup", "backup.log")

type WriteRemoteConfig struct {
	Transport  string // 支持 http || grpc 默认grpc
	Addrs      []string
	ModuleName string

	MaxConcurrent      int           // 默认 1 个并发
	MaxPacket          int           // 默认不缓冲
	WriteTimeout       time.Duration // 默认不超时
	RemoteFailedBackup string        // 远程发送失败，备份文件
}

func (cfg WriteRemoteConfig) Validate() error {
	if len(cfg.Addrs) == 0 {
		return fmt.Errorf("address required, grpc [ip:port]  http[url]")
	}
	if cfg.ModuleName == "" {
		return fmt.Errorf("moduleName required")
	}
	return nil
}

func NewWriteRemoteConfig(addrs []string, moduleName string) WriteRemoteConfig {
	return WriteRemoteConfig{
		Transport:     "grpc",
		Addrs:         addrs,
		ModuleName:    moduleName,
		MaxConcurrent: 50,
		// 包的大小对写入效率有着比较重要的影响。 当设置 1MB时，会快于 64KB
		// 但是小对象对于GC相对更加友好 (grpc 默认最大4MB一个包)
		MaxPacket:          62 << 10,
		WriteTimeout:       5 * time.Second,
		RemoteFailedBackup: _backupFilename,
	}
	// 如果超出并发限制，直接写入文件，缓慢背景发送
}

type WriteRemote struct {
	cfg    WriteRemoteConfig
	pusher Pusher

	//packets *Packets

	packet *Packet

	bw *BackupWrite

	once sync.Once
}

func NewWriteRemote(cfg WriteRemoteConfig) *WriteRemote {
	if err := cfg.Validate(); err != nil {
		panic("config validate error " + err.Error())
	}
	wr := &WriteRemote{
		cfg: cfg,
		//packets: NewPackets(cfg.MaxPacket, cfg.RemoteFailedBackup),
		packet: NewPacket(cfg.ModuleName, cfg.MaxPacket),
		bw:     NewBackupWrite(cfg.RemoteFailedBackup),
	}

	wr.once.Do(func() {
		go wr.initPusher()
		go wr.backgroundRetry()
		go wr.pullPacket()
	})
	return wr
}

func (wr *WriteRemote) Write(b []byte) (n int, err error) {
	//data, flush := wr.packets.AddPacket(b)
	//if flush {
	//	wr.push(NewDataPacket(wr.cfg.ModuleName, data))
	//}
	sendPacket := wr.packet.AppendData(b)
	if sendPacket != nil {
		wr.push(sendPacket)
	}
	return len(b), nil
}

func (wr *WriteRemote) pullPacket() {
	tick := time.NewTicker(2 * time.Second)
	for {
		select {
		case <-tick.C:
			sendPacket := wr.packet.FlushData()
			if sendPacket != nil {
				wr.push(sendPacket)
			}
		}
	}
}

func (wr *WriteRemote) push(in *pb.Packet) {
	if in.Data == nil || len(in.Data) <= 0 {
		// 没有类容的包，直接丢掉
		return
	}
	if wr.pusher == nil {
		_ = wr.backup(in)
		wr.packet.FreePacket(in)
		return
	}
	// 如果发送者满负荷，则直接丢文件
	if wr.pusher.Concurrent() >= wr.cfg.MaxConcurrent {
		_ = wr.backup(in)
		wr.packet.FreePacket(in)
		return
	}

	ctx := context.Background()
	if wr.cfg.WriteTimeout > 0 {
		ctx, _ = context.WithTimeout(context.Background(), wr.cfg.WriteTimeout)
	}

	go func() {
		if err := wr.pusher.PushPacket(ctx, in); err != nil {
			if err == ErrUnavailable {
				// 只有当服务不可用时，放入错误备份文件里
				_ = wr.backup(in)
			}
			log.Printf("write remote push packet %s\n", err.Error())
		}
		wr.packet.FreePacket(in)
	}()
}

func (wr *WriteRemote) initPusher() {
	// 在发送的时候，才去链接， 如果链接不通，不能影响主进程
	tick := time.NewTicker(time.Second)
	for {
		if wr.pusher == nil {
			if wr.cfg.Transport == "http" {
				pusher, err := NewHttpPush(wr.cfg.Addrs[0], wr.cfg.MaxConcurrent)
				if err != nil {
					log.Printf("init http push error %s\n", err.Error())
					goto next
				}
				wr.pusher = pusher
			} else {
				pusher, err := NewGRPCPush(wr.cfg.Addrs, wr.cfg.MaxConcurrent)
				if err != nil {
					log.Printf("init grpc push error %s\n", err.Error())
					goto next
				}
				wr.pusher = pusher

			}
			log.Printf("init %s push success \n", wr.cfg.Transport)
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

func (wr *WriteRemote) backup(in *pb.Packet) error {
	jsonPacket := _jsonPacket{
		ID:     in.Id,
		Module: in.Module,
		Data:   string(in.Data),
	}

	byt, err := json.Marshal(jsonPacket)
	if err != nil {
		return err
	}
	_, err = wr.bw.WriteBakPacket(byt)
	return err
}

func (wr *WriteRemote) backgroundRetry() {
	tick := time.NewTicker(200 * time.Millisecond)
	for {
		select {
		case <-tick.C:
			byt, err := wr.bw.ReadBakPacket()
			if err != nil {
				fmt.Println("packets retry", err.Error())
			}
			if len(byt) > 0 {
				jsonPacket := &_jsonPacket{}
				if err := json.Unmarshal(byt, jsonPacket); err != nil {
					fmt.Println("packets retry", err.Error())
					break
				}
				v := &pb.Packet{
					Id:     jsonPacket.ID,
					Module: jsonPacket.Module,
					Data:   []byte(jsonPacket.Data),
				}
			loop:
				if wr.pusher != nil {
					ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
					if err := wr.pusher.PushPacket(ctx, v); err == nil {
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

func (wr *WriteRemote) Sync() error {
	sendPacket := wr.packet.FlushData()
	if sendPacket != nil {
		wr.push(sendPacket)
	}
	sendEmpty := make(chan struct{}, 1)
	go func() {
		for {
			if wr.pusher != nil && wr.pusher.Concurrent() == 0 {
				time.Sleep(10 * time.Millisecond)
				sendEmpty <- struct{}{}
				return
			}
			time.Sleep(50 * time.Millisecond)
		}
	}()

	ctx, _ := context.WithTimeout(context.Background(), 5*time.Second)
	select {
	case <-ctx.Done():
		log.Println("sync ", ctx.Err())
	case <-sendEmpty:
	}

	return wr.bw.Close()
}
