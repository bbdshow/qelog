package receiver

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"testing"
	"time"

	"github.com/huzhongqing/qelog/api"
)

const (
	host            = "http://127.0.0.1:31081"
	ContentTypeJSON = "application/json"
)

func JSONReader(v interface{}) io.Reader {
	byt, _ := json.Marshal(v)
	fmt.Println(string(byt))
	return bytes.NewReader(byt)
}

func JSONOutput(resp *http.Response, t *testing.T) {
	defer resp.Body.Close()
	byt, _ := ioutil.ReadAll(resp.Body)
	val := make(map[string]interface{})
	if err := json.Unmarshal(byt, &val); err != nil {
		t.Fatal(err)
	}
	fmt.Println(string(byt))
	code, ok := val["code"]
	if ok {
		v, ok1 := code.(int)
		if ok1 && v != 0 {
			t.Fatal(code, val["message"])
		}
	}
}

func TestHTTPService_ReceivePacket(t *testing.T) {
	docs := []string{}
	count := 0
	ts := "20210221 09:00:00"
	ti, _ := time.ParseInLocation("20060102 15:04:05", ts, time.Local)
	for count < 100 {
		count++
		mill := ti.UnixNano()/1e6 + 1
		str := fmt.Sprintf(`{"_level":"INFO","_time":%d.2573,"_caller":"example/main.go:39","_func":"main.loopWriteLogging","_short":"%d","val":%d}`, mill, count, count)
		docs = append(docs, str)
	}

	in := &api.JSONPacket{
		Id:     fmt.Sprintf("%d", time.Now().UnixNano()),
		Module: "example",
		Data:   docs,
	}
	resp, err := http.Post(fmt.Sprintf("%s/v1/receiver/packet", host), ContentTypeJSON, JSONReader(in))
	if err != nil {
		t.Fatal(err)
	}
	JSONOutput(resp, t)
}
