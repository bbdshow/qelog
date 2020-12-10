package wrapzap

import (
	"context"
	"fmt"
	"sync"
	"time"
)

type WriteRemoteConfig struct {
	URL        string
	ModuleName string

	MaxConcurrent int
	MaxPacket     int
	WriteTimeout  time.Duration
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
	cfg     WriteRemoteConfig
	pusher  Pusher
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
	})
	return wr
}

func (wr *WriteRemote) Write(b []byte) (n int, err error) {
	buffers, flush := wr.packets.AddPacket(b)
	if flush {
		data := make([]string, len(buffers))
		for i, buf := range buffers {
			data[i] = string(buf)
		}
		err = wr.push(RandString(16), data)
	}
	return len(b), err
}

func (wr *WriteRemote) push(id string, data []string) error {
	dp := DataPacket{
		Name: wr.cfg.ModuleName,
		ID:   id,
		Data: data,
	}
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
				_ = wr.push(v.ID, v.Data)
			}
		}
	}
}

func (wr *WriteRemote) Sync() error { return nil }
