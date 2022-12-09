package qezap

import (
	"io/ioutil"
	"os"
	"path"
	"strconv"
	"strings"
	"testing"
	"time"
)

func TestNewWriteLocal(t *testing.T) {
	_ = testNewWriteLocal(t)
}

func testNewWriteLocal(t *testing.T) *WriteLocal {
	w := NewWriteLocal(defaultLocalOption())
	if w == nil {
		t.Fatal("write local nil")
	}
	if w.file != nil {
		t.Fatal("write local w file must nil")
	}
	return w
}

func TestWriteSync_Write(t *testing.T) {
	w := testNewWriteLocal(t)
	n, err := w.Write([]byte("hello write sync"))
	if err != nil {
		t.Fatal(err)
	}
	if n != 16 {
		t.Fatal("n", n)
	}
}

func TestWriteSync_isRotate(t *testing.T) {
	w := testNewWriteLocal(t)
	w.opt.MaxSize = 1024
	largeTxt := ""
	for i := 0; i < int(w.opt.MaxSize); i++ {
		largeTxt += strconv.Itoa(i)
	}
	var testCases = []struct {
		In       string
		IsRotate bool
	}{
		{
			In:       "hello write local",
			IsRotate: false,
		},
		{
			In:       largeTxt,
			IsRotate: true,
		},
	}

	for i, v := range testCases {
		n, err := w.Write([]byte(v.In))
		if err != nil {
			t.Fatal(err)
		}
		if n != len(v.In) {
			t.Fatalf("case %d :in len invalid", i)
		}
		if err := w.isRotate(n); err != nil {
			t.Fatalf("case %d :isRotate %v", i, err)
		}

		time.Sleep(time.Second)

		fs, err := ioutil.ReadDir(w.dir)
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
			if v.IsRotate {
				t.Fatalf("case %d: should rotate, but not found file .gz", i)
			}
		} else {
			if err := os.Remove(path.Join(w.dir, hitFilename)); err != nil {
				t.Fatal("remove", hitFilename, err.Error())
			}

			if !v.IsRotate {
				t.Fatalf("case %d:not should rotate, but found file %s", i, hitFilename)
			}
		}
	}
	w.Close()
	time.Sleep(time.Second)
	if err := os.RemoveAll(w.dir); err != nil {
		t.Fatalf(err.Error())
	}
}

//func TestWrite_delExpiredFile(t *testing.T) {
//	cfg := NewConfig(nil, "").
//		SetMaxSize(1024).
//		SetMaxAge(10 * time.Second)
//
//	ws := NewWriteSync(cfg)
//	go func() {
//		tick := time.NewTicker(30 * time.Second)
//		for {
//			select {
//			case <-tick.C:
//				tick.Stop()
//				return
//			default:
//				ws.Write([]byte("hello write sync"))
//				time.Sleep(10 * time.Millisecond)
//			}
//		}
//	}()
//	time.Sleep(90 * time.Second)
//}
