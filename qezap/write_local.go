package qezap

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

// WriteLocal local fs writer impl
// support rotate, gzip, maxAge
type WriteLocal struct {
	mutex sync.Mutex
	once  sync.Once

	opt *LocalOption

	dir  string
	size int64

	// check compress
	compressing chan struct{}

	// local fs object
	file *os.File

	isExit int32
	exit   chan struct{}
}

func NewWriteLocal(opt *LocalOption) *WriteLocal {
	w := &WriteLocal{
		mutex: sync.Mutex{},
		opt:   opt,
		dir:   path.Dir(opt.Filename),
		size:  0,

		compressing: make(chan struct{}, 1),
		file:        nil,

		exit: make(chan struct{}),
	}
	w.once.Do(func() {
		go w.bgDelExpiredFile()
	})
	return w
}

// Write impl
func (w *WriteLocal) Write(b []byte) (n int, err error) {
	w.mutex.Lock()
	defer w.mutex.Unlock()

	if w.file == nil {
		err = w.openFile()
		if err != nil {
			return n, err
		}
	}

	n, err = w.file.Write(b)
	if err != nil {
		return n, err
	}

	if err := w.isRotate(n); err != nil {
		return n, errors.New("rotate " + err.Error())
	}

	return n, nil
}

// Sync empty impl
func (w *WriteLocal) Sync() error {
	return nil
}

// Close file handle
func (w *WriteLocal) Close() error {
	if atomic.LoadInt32(&w.isExit) == 1 {
		return nil
	}
	atomic.StoreInt32(&w.isExit, 1)
	close(w.exit)

	// if compressing, waiting
	w.compressing <- struct{}{}

	if w.file != nil {
		return w.file.Close()
	}
	return nil
}

func (w *WriteLocal) openFile() error {
	err := os.MkdirAll(w.dir, os.ModePerm|os.ModeDir)
	if err != nil {
		return err
	}

	// 查看文件信息
	info, err := os.Stat(w.opt.Filename)
	if err != nil {
		if os.IsNotExist(err) {
			// 不存在文件
			f, err := os.Create(w.opt.Filename)
			if err != nil {
				return err
			}
			w.file = f
			return nil
		}
		return err
	}

	// 存在文件
	f, err := os.OpenFile(w.opt.Filename, os.O_APPEND|os.O_WRONLY, os.ModeAppend)
	if err != nil {
		return err
	}
	w.file = f
	w.size = info.Size()

	return nil
}

// gzip compress file
func (w *WriteLocal) gzipCompress(filename string) error {
	if filename == "" || atomic.LoadInt32(&w.isExit) == 1 {
		return nil
	}
	w.compressing <- struct{}{}
	defer func() {
		<-w.compressing
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
	// delete src file
	_ = src.Close()
	return os.Remove(filename)
}

func (w *WriteLocal) isRotate(n int) error {
	if w.opt.MaxSize <= 0 || atomic.AddInt64(&w.size, int64(n)) < w.opt.MaxSize {
		return nil
	}
	// 滚动先关闭原文件
	if err := w.file.Close(); err != nil {
		fmt.Println("ws.file.Close()", err.Error())
	}
	w.file = nil
	// 滚动, 有任何操作失败的地方，都不滚动
	rotateFilename := w.rotateFilename()
	err := os.Rename(w.opt.Filename, rotateFilename)
	if err != nil {
		return err
	}
	// 新建文件
	f, err := os.Create(w.opt.Filename)
	if err != nil {
		return err
	}
	// 新建文件
	w.size = 0
	w.file = f

	if w.opt.GzipCompress {
		go func() {
			if err := w.gzipCompress(rotateFilename); err != nil {
				log.Println("gzip compress", err.Error())
			}
		}()
	}
	return nil
}

func (w *WriteLocal) rotateFilename() string {
	filename := strings.Replace(w.opt.Filename, ".log", fmt.Sprintf("%s.bak.log", time.Now().Format("20060102150405.00")), 1)
	return filename
}

// delete expired backup file
func (w *WriteLocal) bgDelExpiredFile() {
	if w.opt.MaxAge <= 0 {
		return
	}

	tick := time.NewTicker(30 * time.Second)
	for {
		select {
		case <-tick.C:
			expired := time.Now().Add(-w.opt.MaxAge)
			fs, err := ioutil.ReadDir(w.dir)
			if err == nil {
				for _, f := range fs {
					if !f.IsDir() {
						// delete suffix *.bak.log *.bak.log.gz
						if strings.HasSuffix(f.Name(), ".bak.log") ||
							strings.HasSuffix(f.Name(), ".bak.log.gz") {
							if f.ModTime().Before(expired) {
								w.mutex.Lock()
								_ = os.Remove(path.Join(w.dir, f.Name()))
								w.mutex.Unlock()
							}
						}
					}
				}
			}
		case <-w.exit:
			return
		}
	}
}
