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

	"github.com/huzhongqing/qelog/pb"
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
	a := `{"_level":"INFO","_time":1609944069261.2573,"_caller":"example/main.go:39","_func":"main.loopWriteLogging","_short":"1","_traceid":"160994406926125730025244977","val":823537}`
	b := `{"_level":"INFO","_time":1609944069261.2573,"_caller":"example/main.go:39","_func":"main.loopWriteLogging","_short":"2","_traceid":"160994406926125730025244977","val":823537}`
	c := `{"_level":"INFO","_time":1612630357000.2573,"_caller":"example/main.go:39","_func":"main.loopWriteLogging","_short":"3","_traceid":"160994406926125730025244977","val":823537}`

	in := &pb.Packet{
		Id:     fmt.Sprintf("%d", time.Now().UnixNano()),
		Module: "example",
		Data:   []string{a, b, c},
	}
	resp, err := http.Post(fmt.Sprintf("%s/v1/receiver/packet", host), ContentTypeJSON, JSONReader(in))
	if err != nil {
		t.Fatal(err)
	}
	JSONOutput(resp, t)
}
