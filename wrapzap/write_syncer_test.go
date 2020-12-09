package wrapzap

import (
	"io/ioutil"
	"os"
	"strings"
	"testing"
	"time"
)

func TestWriteSync_Write(t *testing.T) {
	cfg := DefaultWriteConfig("logger.log")
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
	cfg := DefaultWriteConfig("logger.log")
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
