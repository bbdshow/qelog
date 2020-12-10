package wrapzap

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"math/rand"
	"os"
	"sync"
	"time"
)

const (
	defaultFlushInterval = 3 // time.Second
)

type Packets struct {
	mutex          sync.Mutex
	maxSize        int
	flushTimestamp int64 // 3秒一个周期

	bufferSize int
	buffers    bytes.Buffer

	w           *WriteSync
	bakFilename string
	offset      int64
}

type FailedPacket struct {
	ID   string
	Data string
}

func NewFailedPacket(data string) FailedPacket {
	return FailedPacket{
		ID:   RandString(16),
		Data: data,
	}
}

func NewPackets(maxSize int) *Packets {
	if maxSize <= 0 {
		maxSize = 0
	}
	p := &Packets{
		maxSize:        maxSize,
		flushTimestamp: time.Now().Unix(),
		bufferSize:     0,
		buffers:        bytes.Buffer{},
		bakFilename:    "./failed.bak/packets.bak",
		offset:         0,
	}

	p.w = NewWriteSync(WriteSyncConfig{
		Filename:     p.bakFilename,
		MaxSize:      0, // 不滚动
		TTL:          0, // 不切割
		GzipCompress: false,
	})

	return p
}

func (p *Packets) AddPacket(b []byte) (buffers []byte, flush bool) {
	p.mutex.Lock()

	// 缓存起，超过一定时间/容量再发送
	n, _ := p.buffers.Write(b)
	p.bufferSize += n

	if p.bufferSize >= p.maxSize || time.Now().Unix()-p.flushTimestamp > defaultFlushInterval {
		flush = true
		buffers = p.buffers.Bytes()

		p.buffers.Reset()
		p.bufferSize = 0
		p.flushTimestamp = time.Now().Unix()
	}

	p.mutex.Unlock()

	return buffers, flush
}

func (p *Packets) PutFailedPacket(v interface{}) (n int, err error) {
	p.mutex.Lock()
	b, _ := json.Marshal(v)
	n, err = fmt.Fprintln(p.w, string(b))
	p.mutex.Unlock()
	return
}

func (p *Packets) GetFailedPacket(v interface{}) (ok bool, err error) {
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
	defer func() {
		if f != nil {
			_ = f.Close()
		}
	}()
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
				if err := os.Remove(p.bakFilename); err == nil {
					p.offset = 0
				}
				break
			}
			return false, err
		}

		p.offset += int64(len(b))

		if err := json.Unmarshal(b, &v); err != nil {
			return false, err
		}

		break
	}

	return true, nil
}

func RandString(length int) string {
	baseChar := "0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
	byteChar := []byte(baseChar)
	str := ""
	for i := 0; i < length; i++ {
		rand.New(rand.NewSource(time.Now().UnixNano() + rand.Int63n(1000000)))
		str += string(byteChar[rand.Intn(len(byteChar))])
	}
	return str
}
