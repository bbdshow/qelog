package wrapzap

import (
	"context"
	"fmt"
	"net/url"
	"sync"
	"time"
)

type WriteRemoteConfig struct {
	URL        string
	ModuleName string

	MaxConcurrent int           // 默认 1 个并发
	MaxPacket     int           // 默认不缓冲
	WriteTimeout  time.Duration // 默认不超时
}

func (cfg WriteRemoteConfig) Validate() error {
	_, err := url.Parse(cfg.URL)
	if err != nil {
		return err
	}
	if cfg.ModuleName == "" {
		return fmt.Errorf("moduleName required")
	}
	return nil
}

func NewWriteRemoteConfig(url, moduleName string) WriteRemoteConfig {
	return WriteRemoteConfig{
		URL:           url,
		ModuleName:    moduleName,
		MaxConcurrent: 20,
		MaxPacket:     500 * 1024,
		WriteTimeout:  5 * time.Second,
	}
	// 此配置当百兆网络畅通时, 理论，支持每秒10MB日志生成不卡顿， 如果超出，直接写入文件，缓慢背景发送
}

type WriteRemote struct {
	cfg    WriteRemoteConfig
	pusher Pusher

	packets *Packets

	once sync.Once
}

func NewWriteRemote(cfg WriteRemoteConfig) *WriteRemote {
	wr := &WriteRemote{
		cfg:     cfg,
		pusher:  NewHttpPush(cfg.URL, cfg.MaxConcurrent),
		packets: NewPackets(cfg.MaxPacket),
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
		err = wr.push(NewDataPacket(wr.cfg.ModuleName, data))
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
				_ = wr.push(NewDataPacket(wr.cfg.ModuleName, data))
			}
		}
	}
}

func (wr *WriteRemote) push(dp DataPacket) error {
	fmt.Println("data", dp)
	// 如果发送者满负荷，则直接丢文件
	if wr.pusher.Concurrent() >= wr.cfg.MaxConcurrent {
		_, _ = wr.packets.WritePacket(dp)
		return nil
	}

	ctx := context.Background()
	if wr.cfg.WriteTimeout > 0 {
		ctx, _ = context.WithTimeout(context.Background(), wr.cfg.WriteTimeout)
	}
	if err := wr.pusher.Push(ctx, dp.Marshal()); err != nil {
		fmt.Println("pusher", err.Error())
		// 放入错误备份文件里
		_, _ = wr.packets.WritePacket(dp)
		return err
	}
	return nil
}

func (wr *WriteRemote) backgroundRetry() {
	tick := time.NewTicker(time.Second)
	for {
		select {
		case <-tick.C:
			v := DataPacket{}
			ok, err := wr.packets.ReadPacket(&v)
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
		return wr.push(NewDataPacket(wr.cfg.ModuleName, data))
	}
	return nil
}
