package qzzap

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"os"
	"sync"
)

type Packets struct {
	mutex sync.Mutex

	maxSize  int
	dataSize int
	data     []string

	w           *WriteSync
	bakFilename string
	offset      int64
}

func NewPackets(maxSize int) *Packets {
	if maxSize <= 0 {
		maxSize = 0
	}
	p := &Packets{
		maxSize:     maxSize,
		dataSize:    0,
		data:        make([]string, 0, 1024),
		bakFilename: "./failed.bak/packets.bak",
		offset:      0,
	}
	p.initWrite()

	return p
}

func (p *Packets) initWrite() {
	p.w = NewWriteSync(WriteSyncConfig{
		Filename:     p.bakFilename,
		MaxSize:      0, // 不滚动
		TTL:          0, // 不切割
		GzipCompress: false,
	})
}

func (p *Packets) AddPacket(b []byte) (data []string, flush bool) {
	p.mutex.Lock()
	// 缓存起，超过一定时间/容量再发送
	p.data = append(p.data, string(b))
	p.dataSize += len(b)

	if p.dataSize >= p.maxSize {
		data = make([]string, len(p.data))
		copy(data, p.data)

		flush = true
		// reset
		p.data = p.data[:0]
		p.dataSize = 0
	}
	p.mutex.Unlock()

	return data, flush
}

func (p *Packets) PullPacket() (data []string, flush bool) {
	p.mutex.Lock()
	if p.dataSize > 0 {
		data = make([]string, len(p.data))
		copy(data, p.data)

		flush = true
		// reset
		p.data = p.data[:0]
		p.dataSize = 0
	}
	p.mutex.Unlock()
	return data, flush
}

func (p *Packets) WritePacket(v interface{}) (n int, err error) {
	p.mutex.Lock()
	b, _ := json.Marshal(v)
	if p.w == nil {
		p.initWrite()
	}
	n, err = fmt.Fprintln(p.w, string(b))
	p.mutex.Unlock()
	return
}

func (p *Packets) ReadPacket(v interface{}) (ok bool, err error) {
	p.mutex.Lock()
	defer p.mutex.Unlock()

	_, err = os.Stat(p.bakFilename)
	if os.IsNotExist(err) {
		return false, nil
	}
	f, err := os.Open(p.bakFilename)
	if err != nil {
		return false, err
	}
	if _, err := f.Seek(p.offset, io.SeekStart); err != nil {
		return false, err
	}
	buf := bufio.NewReader(f)
	for {
		b, err := buf.ReadBytes('\n')
		if err != nil {
			if err == io.EOF {
				// 文件读取完了，就删除了
				// 关闭 file io
				_ = f.Close()
				_ = p.w.Close()

				p.w = nil
				if err := os.Remove(p.bakFilename); err == nil {
					p.offset = 0
				} else {
					log.Printf("os remove %s error %s\n", p.bakFilename, err.Error())
				}
				break
			}
			return false, err
		}

		p.offset += int64(len(b))

		if err := json.Unmarshal(b, &v); err != nil {
			return false, err
		}
		_ = f.Close()
		break
	}

	return true, nil
}
