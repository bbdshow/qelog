package manager

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"testing"

	"github.com/huzhongqing/qelog/pkg/common/entity"
)

const (
	host            = "http://127.0.0.1:31080"
	ContentTypeJSON = "application/json"
)

func JSONReader(v interface{}) io.Reader {
	byt, _ := json.Marshal(v)
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

func TestCreateModuleRegister(t *testing.T) {
	in := entity.CreateModuleRegisterReq{
		ModuleName: "example",
		DBIndex:    1,
		Desc:       "example 演示",
	}
	resp, err := http.Post(fmt.Sprintf("%s/v1/module", host), ContentTypeJSON, JSONReader(in))
	if err != nil {
		t.Fatal(err)
	}
	JSONOutput(resp, t)
}
