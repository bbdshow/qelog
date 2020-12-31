package qezap

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"os"
	"sync"
)

// data 数据最大容量
var _maxPacketSize = 1024

type Packets struct {
	data *BuffSliceString

	mutex       sync.Mutex
	backWrite   *WriteSync
	bakFilename string
	offset      int64
}

type BuffSliceString struct {
	mutex sync.RWMutex
	val   []string
	size  int
}

func NewBuffSliceString() *BuffSliceString {
	return &BuffSliceString{val: make([]string, 0, 1024)}
}

func (bss *BuffSliceString) Append(v string) {
	bss.mutex.Lock()
	bss.val = append(bss.val, v)
	bss.size += len(v)
	bss.mutex.Unlock()
}

func (bss *BuffSliceString) Reset() {
	bss.mutex.Lock()
	bss.val = bss.val[:0]
	bss.size = 0
	bss.mutex.Unlock()
}

func (bss *BuffSliceString) Len() int {
	bss.mutex.RLock()
	l := len(bss.val)
	bss.mutex.RUnlock()
	return l
}

func (bss *BuffSliceString) Size() int {
	bss.mutex.RLock()
	size := bss.size
	bss.mutex.RUnlock()
	return size
}

func (bss *BuffSliceString) Val() []string {
	bss.mutex.RLock()
	val := bss.val[0:]
	bss.mutex.RUnlock()
	return val
}

var dataFree = sync.Pool{
	New: func() interface{} { return NewBuffSliceString() },
}

func NewPackets(maxSize int, backup string) *Packets {
	if maxSize > 1 {
		_maxPacketSize = maxSize
	}
	data := dataFree.Get().(*BuffSliceString)
	p := &Packets{
		data:        data,
		bakFilename: backup,
		offset:      0,
	}
	p.initBackWrite()

	return p
}

func (p *Packets) initBackWrite() {
	p.backWrite = NewWriteSync(WriteSyncConfig{
		Filename:     p.bakFilename,
		MaxSize:      0, // 不滚动
		TTL:          0, // 不切割
		GzipCompress: false,
	})
}

func (p *Packets) AddPacket(b []byte) (*BuffSliceString, bool) {
	// 缓存起，超过一定时间/容量再发送
	p.data.Append(string(b))
	if p.data.Size() >= _maxPacketSize {
		return p.PullPacket()
	}
	return nil, false
}

func (p *Packets) Free(bss *BuffSliceString) {
	bss.Reset()
	dataFree.Put(bss)
}

func (p *Packets) PullPacket() (*BuffSliceString, bool) {
	bss := p.data
	p.data = dataFree.Get().(*BuffSliceString)
	return bss, true
}

func (p *Packets) WriteBakPacket(v interface{}) (n int, err error) {
	p.mutex.Lock()
	b, _ := json.Marshal(v)
	if p.backWrite == nil {
		p.initBackWrite()
	}
	n, err = fmt.Fprintln(p.backWrite, string(b))
	p.mutex.Unlock()
	return
}

func (p *Packets) ReadBakPacket(v interface{}) (ok bool, err error) {
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
	b, err := buf.ReadBytes('\n')
	if err != nil {
		if err == io.EOF {
			//  清空文件夹
			if _, err := f.Seek(0, io.SeekStart); err != nil {
				log.Printf("f.Seek offset zero %s\n", err.Error())
				return false, nil
			}
			if err := os.Truncate(p.bakFilename, 0); err != nil {
				log.Printf(" os.Truncate  %s\n", err.Error())
				return false, nil
			}
			_ = f.Close()
			p.offset = 0
			return false, nil
		}
		return false, err
	}

	if len(b) > 0 {
		p.offset += int64(len(b))
		if err := json.Unmarshal(b, &v); err != nil {
			return false, err
		}
		_ = f.Close()
		return true, nil
	}

	return false, nil
}

func (p *Packets) Close() error {
	if p.backWrite == nil {
		return nil
	}
	return p.backWrite.Close()
}
