package wrapzap

import (
	"io/ioutil"
	"os"
	"strings"
	"testing"
	"time"
)

func TestWriteSync_Write(t *testing.T) {
	cfg := DefaultWriteSyncConfig("logger.log")
	ws := NewWriteSync(cfg)

	n, err := ws.Write([]byte("hello write sync"))
	if err != nil {
		t.Fatal(err)
	}
	if n != 16 {
		t.Fatal("n", n)
	}
}

func TestWriteSync_isRotate(t *testing.T) {
	cfg := DefaultWriteSyncConfig("logger.log")
	cfg.MaxSize = 1024
	ws := NewWriteSync(cfg)
	ws.Write([]byte("hello write sync"))
	if err := ws.isRotate(1024); err != nil {
		t.Fatal(err)
	}
	time.Sleep(time.Second)
	fs, err := ioutil.ReadDir("./")
	if err != nil {
		t.Fatal(err)
	}
	hit := false
	hitFilename := ""
	for _, f := range fs {
		if strings.HasSuffix(f.Name(), ".gz") {
			hit = true
			hitFilename = f.Name()
			break
		}
	}
	if !hit {
		t.Fatal("not hit .gz")
	}
	if err := os.Remove(hitFilename); err != nil {
		t.Fatal("remove", hitFilename)
	}
}

func TestWrite_delExpiredFile(t *testing.T) {
	cfg := DefaultWriteSyncConfig("logger.log")
	cfg.MaxSize = 1024
	cfg.TTL = 10 * time.Second
	ws := NewWriteSync(cfg)
	go func() {
		tick := time.NewTicker(30 * time.Second)
		for {
			select {
			case <-tick.C:
				tick.Stop()
				return
			default:
				ws.Write([]byte("hello write sync"))
				time.Sleep(10 * time.Millisecond)
			}
		}
	}()
	time.Sleep(90 * time.Second)
}
