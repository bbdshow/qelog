package wrapzap

import (
	"compress/gzip"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
	"path"
	"strings"
	"sync"
	"sync/atomic"
	"time"
)

type WriteSync struct {
	mutex sync.Mutex

	dir      string
	filename string

	size    int64
	maxSize int64

	ttl time.Duration // 0 永久保存

	compress    bool
	compressing chan struct{} // 正在压缩

	// 文件对象
	file *os.File
}

type WriteConfig struct {
	Filename string
	// 单文件最大 bytes
	MaxSize int64
	// 保留文件时间
	TTL time.Duration
	// Gzip 压缩
	GzipCompress bool
}

func DefaultWriteConfig(filename string) WriteConfig {
	if filename == "" {
		panic("write sync filename required")
	}
	return WriteConfig{
		Filename:     filename,
		MaxSize:      200 << 20, // 200MB
		TTL:          0,         // 保存永久
		GzipCompress: true,
	}
}

func NewWriteSync(cfg WriteConfig) *WriteSync {
	return &WriteSync{
		mutex:       sync.Mutex{},
		dir:         path.Dir(cfg.Filename),
		filename:    cfg.Filename,
		size:        0,
		maxSize:     cfg.MaxSize,
		ttl:         cfg.TTL,
		compress:    cfg.GzipCompress,
		compressing: make(chan struct{}, 1),
		file:        nil,
	}
}

func (ws *WriteSync) Write(b []byte) (n int, err error) {
	ws.mutex.Lock()
	defer ws.mutex.Unlock()

	if ws.file == nil {
		err = ws.openFile()
		if err != nil {
			return n, err
		}
	}

	n, err = ws.file.Write(b)
	if err != nil {
		return n, err
	}
	if err := ws.isRotate(n); err != nil {
		return n, errors.New("rotate " + err.Error())
	}

	return n, nil
}

func (ws *WriteSync) Sync() error {
	return nil
}

func (ws *WriteSync) Close() error {
	// 如果正在压缩，等压缩完再退出
	ws.compressing <- struct{}{}
	if ws.file != nil {
		return ws.file.Close()
	}
	return nil
}

// 打开文件
func (ws *WriteSync) openFile() error {
	err := os.MkdirAll(ws.dir, os.ModePerm|os.ModeDir)
	if err != nil {
		return err
	}

	// 查看文件信息
	info, err := os.Stat(ws.filename)
	if err != nil {
		if os.IsNotExist(err) {
			// 不存在文件
			f, err := os.Create(ws.filename)
			if err != nil {
				return err
			}
			ws.file = f
			return nil
		}
		return err
	}

	// 存在文件
	f, err := os.OpenFile(ws.filename, os.O_APPEND|os.O_WRONLY, os.ModeAppend)
	if err != nil {
		return err
	}
	ws.file = f
	ws.size = info.Size()

	return nil
}

// 压缩文件
func (ws *WriteSync) gzipCompress(filename string) error {
	if filename == "" {
		return nil
	}
	ws.compressing <- struct{}{}
	defer func() {
		<-ws.compressing
	}()
	destFilename := filename + ".gz"
	dest, err := os.Create(destFilename)
	if err != nil {
		return err
	}
	defer dest.Close()

	gz := gzip.NewWriter(dest)
	defer gz.Close()

	src, err := os.Open(filename)
	if err != nil {
		return err
	}

	srcInfo, _ := src.Stat()
	gz.Name = destFilename
	gz.ModTime = srcInfo.ModTime()

	if _, err := io.Copy(gz, src); err != nil {
		return err
	}
	if err := gz.Flush(); err != nil {
		return err
	}
	// 删除原文件
	_ = src.Close()
	return os.Remove(filename)
}
func (ws *WriteSync) isRotate(n int) error {
	if ws.maxSize <= 0 || atomic.AddInt64(&ws.size, int64(n)) < ws.maxSize {
		return nil
	}

	_ = ws.file.Close()
	ws.file = nil
	// 滚动, 有任何操作失败的地方，都不滚动
	rotateFilename := ws.rotateFilename()
	err := os.Rename(ws.filename, rotateFilename)
	if err != nil {
		return err
	}
	// 新建文件
	f, err := os.Create(ws.filename)
	if err != nil {
		return err
	}
	// 新建文件
	ws.size = 0
	ws.file = f

	if ws.compress {
		go func() {
			if err := ws.gzipCompress(rotateFilename); err != nil {
				log.Println("gzip compress", err.Error())
			}
		}()
	}
	//  检查一下过期文件
	go ws.delExpiredFile()
	return nil
}

func (ws *WriteSync) rotateFilename() string {
	filename := strings.Replace(ws.filename, ".log", fmt.Sprintf("%s.log", time.Now().Format("20060102150405.00")), 1)
	return filename
}
func (ws *WriteSync) delExpiredFile() {
	if ws.ttl > 0 {
		expired := time.Now().Add(-ws.ttl)
		fs, err := ioutil.ReadDir(ws.dir)
		if err == nil {
			for _, f := range fs {
				if !f.IsDir() {
					// 只删除 .log 或者 .log.gz
					if strings.HasSuffix(f.Name(), ".log") || strings.HasSuffix(f.Name(), ".log.gz") {
						if f.ModTime().Before(expired) {
							_ = os.Remove(path.Join(ws.dir, f.Name()))
						}
					}
				}
			}
		}
	}
}
