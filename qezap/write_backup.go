package qezap

import (
	"bufio"
	"bytes"
	"io"
	"log"
	"os"
	"sync"
)

// writeBackup offset read backup file,used to retry remote packet
type writeBackup struct {
	mutex    sync.Mutex
	local    *WriteLocal
	filename string
	offset   int64
}

func newWriteBackup(filename string) *writeBackup {
	w := &writeBackup{
		local:    nil,
		filename: filename,
		offset:   0,
	}
	w.initLocalWrite()
	return w
}

func (w *writeBackup) initLocalWrite() {
	opt := DefaultLocalOption()
	opt.Filename = w.filename
	opt.MaxAge = 0
	opt.MaxSize = 0
	opt.GzipCompress = false
	w.local = NewWriteLocal(opt)
}

// WriteBakPacket write need retry packet. use '\n' end
func (w *writeBackup) WriteBakPacket(b []byte) (n int, err error) {
	w.mutex.Lock()
	defer w.mutex.Unlock()

	if w.local == nil {
		w.initLocalWrite()
	}
	n, err = w.local.Write(b)
	if err != nil {
		return n, err
	}
	// \n used to split
	_, _ = w.local.Write([]byte{'\n'})
	return n, nil
}

// ReadBakPacket read need retry packet.
func (w *writeBackup) ReadBakPacket() (b []byte, err error) {
	w.mutex.Lock()
	defer w.mutex.Unlock()

	_, err = os.Stat(w.filename)
	if os.IsNotExist(err) {
		return nil, nil
	}
	f, err := os.Open(w.filename)
	if err != nil {
		return nil, err
	}
	if _, err := f.Seek(w.offset, io.SeekStart); err != nil {
		return nil, err
	}
	buf := bufio.NewReader(f)
	b, err = buf.ReadBytes('\n')
	if err != nil {
		if err == io.EOF {
			//  all data retry over, reset offset
			if _, err := f.Seek(0, io.SeekStart); err != nil {
				log.Printf("f.Seek offset zero %s\n", err.Error())
				return nil, nil
			}
			// truncate file
			if err := os.Truncate(w.filename, 0); err != nil {
				log.Printf(" os.Truncate  %s\n", err.Error())
				return nil, nil
			}
			_ = f.Close()
			w.offset = 0
			return nil, nil
		}
		return nil, err
	}

	if len(b) > 0 {
		w.offset += int64(len(b))
		b = bytes.TrimSuffix(b, []byte{'\n'})
		return b, nil
	}

	return nil, nil
}

// Close release file handle
func (w *writeBackup) Close() error {
	if w.local != nil {
		err := w.local.Close()
		w.local = nil
		return err
	}
	return nil
}
