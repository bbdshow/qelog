package qezap

import (
	"bufio"
	"bytes"
	"io"
	"log"
	"os"
	"sync"
)

type BackupWrite struct {
	mutex    sync.Mutex
	w        *WriteSync
	filename string
	offset   int64
}

func NewBackupWrite(filename string) *BackupWrite {
	bw := &BackupWrite{
		w:        nil,
		filename: filename,
		offset:   0,
	}
	bw.initWrite()
	return bw
}

func (bw *BackupWrite) initWrite() {
	bw.w = NewWriteSync(WriteSyncConfig{
		Filename:     bw.filename,
		MaxSize:      0, // 不滚动
		TTL:          0, // 不切割
		GzipCompress: false,
	})
}

func (bw *BackupWrite) WriteBakPacket(b []byte) (n int, err error) {
	bw.mutex.Lock()
	if bw.w == nil {
		bw.initWrite()
	}
	n, err = bw.w.Write(b)
	if err != nil {
		return n, err
	}
	n1, err := bw.w.Write([]byte{'\n'})
	if err != nil {
		return n, err
	}
	bw.mutex.Unlock()
	return n + n1, nil
}

func (bw *BackupWrite) ReadBakPacket() (b []byte, err error) {
	bw.mutex.Lock()
	defer bw.mutex.Unlock()

	_, err = os.Stat(bw.filename)
	if os.IsNotExist(err) {
		return nil, nil
	}
	f, err := os.Open(bw.filename)
	if err != nil {
		return nil, err
	}
	if _, err := f.Seek(bw.offset, io.SeekStart); err != nil {
		return nil, err
	}
	buf := bufio.NewReader(f)
	b, err = buf.ReadBytes('\n')
	if err != nil {
		if err == io.EOF {
			//  清空文件夹
			if _, err := f.Seek(0, io.SeekStart); err != nil {
				log.Printf("f.Seek offset zero %s\n", err.Error())
				return nil, nil
			}
			if err := os.Truncate(bw.filename, 0); err != nil {
				log.Printf(" os.Truncate  %s\n", err.Error())
				return nil, nil
			}
			_ = f.Close()
			bw.offset = 0
			return nil, nil
		}
		return nil, err
	}

	if len(b) > 0 {
		bw.offset += int64(len(b))
		b = bytes.TrimSuffix(b, []byte{'\n'})
		return b, nil
	}

	return nil, nil
}

func (bw *BackupWrite) Close() error {
	if bw.w != nil {
		err := bw.w.Close()
		bw.w = nil
		return err
	}
	return nil
}
