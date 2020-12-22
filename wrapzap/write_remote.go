package wrapzap

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/huzhongqing/wrapzap/push"
)

type WriteRemoteConfig struct {
	Transport  string // 支持 http grpc 默认 grpc
	Addrs      []string
	ModuleName string

	MaxConcurrent int           // 默认 1 个并发
	MaxPacket     int           // 默认不缓冲
	WriteTimeout  time.Duration // 默认不超时
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

func (cfg *WriteRemoteConfig) SetHTTPTransport() *WriteRemoteConfig {
	cfg.Transport = "http"
	return cfg
}

func NewWriteRemoteConfig(addrs []string, moduleName string) WriteRemoteConfig {
	return WriteRemoteConfig{
		Transport:     "grpc",
		Addrs:         addrs,
		ModuleName:    moduleName,
		MaxConcurrent: 20,
		MaxPacket:     500 * 1024,
		WriteTimeout:  5 * time.Second,
	}
	// 此配置当百兆网络畅通时, 理论，支持每秒10MB日志生成不卡顿， 如果超出，直接写入文件，缓慢背景发送
}

type WriteRemote struct {
	cfg    WriteRemoteConfig
	pusher push.Pusher

	packets *Packets

	once sync.Once
}

func NewWriteRemote(cfg WriteRemoteConfig) *WriteRemote {
	if err := cfg.Validate(); err != nil {
		panic("config validate error " + err.Error())
	}
	wr := &WriteRemote{
		cfg:     cfg,
		packets: NewPackets(cfg.MaxPacket),
	}

	if cfg.Transport == "http" {
		pusher, err := push.NewHttpPush(cfg.Addrs[0], cfg.MaxConcurrent)
		if err != nil {
			panic("init http push error " + err.Error())
		}
		wr.pusher = pusher
	} else {
		pusher, err := push.NewGRPCPush(cfg.Addrs, cfg.MaxConcurrent)
		if err != nil {
			panic("init grpc push error " + err.Error())
		}
		wr.pusher = pusher
	}

	wr.once.Do(func() {
		go wr.backgroundRetry()
		go wr.pullPacket()
	})
	return wr
}

func (wr *WriteRemote) Write(b []byte) (n int, err error) {
	data, flush := wr.packets.AddPacket(b)
	if flush && len(data) > 0 {
		err = wr.push(push.NewPacket(wr.cfg.ModuleName, data))
	}
	return len(b), err
}

func (wr *WriteRemote) pullPacket() {
	tick := time.NewTicker(2 * time.Second)
	for {
		select {
		case <-tick.C:
			data, flush := wr.packets.PullPacket()
			if flush && len(data) > 0 {
				_ = wr.push(push.NewPacket(wr.cfg.ModuleName, data))
			}
		}
	}
}

func (wr *WriteRemote) push(in *push.Packet) error {
	// 如果发送者满负荷，则直接丢文件
	if wr.pusher.Concurrent() >= wr.cfg.MaxConcurrent {
		_, _ = wr.packets.WritePacket(in)
		return nil
	}

	ctx := context.Background()
	if wr.cfg.WriteTimeout > 0 {
		ctx, _ = context.WithTimeout(context.Background(), wr.cfg.WriteTimeout)
	}

	if err := wr.pusher.PushPacket(ctx, in); err != nil {
		// 放入错误备份文件里
		_, _ = wr.packets.WritePacket(in)
		return err
	}
	return nil
}

func (wr *WriteRemote) backgroundRetry() {
	tick := time.NewTicker(time.Second)
	for {
		select {
		case <-tick.C:
			v := &push.Packet{}
			ok, err := wr.packets.ReadPacket(v)
			if err != nil {
				fmt.Println("packets retry", err.Error())
			}
			if ok {
				v.IsRetry = true
				_ = wr.push(v)
			}
		}
	}
}

func (wr *WriteRemote) Sync() error {
	data, flush := wr.packets.PullPacket()
	if flush && len(data) > 0 {
		return wr.push(push.NewPacket(wr.cfg.ModuleName, data))
	}
	return nil
}
