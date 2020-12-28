package manager

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"testing"
	"time"

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

func TestCreateModule(t *testing.T) {
	in := entity.CreateModuleReq{
		Name:    "example",
		DBIndex: 1,
		Desc:    "example 演示",
	}
	resp, err := http.Post(fmt.Sprintf("%s/v1/module", host), ContentTypeJSON, JSONReader(in))
	if err != nil {
		t.Fatal(err)
	}
	JSONOutput(resp, t)
}

func TestFindLoggingList(t *testing.T) {
	in := entity.FindLoggingListReq{
		DBIndex:        1,
		ModuleName:     "example",
		Short:          "",
		Level:          -1,
		IP:             "",
		ConditionOne:   "",
		ConditionTwo:   "",
		ConditionThree: "",
		TimeReq:        entity.TimeReq{BeginTsSec: time.Now().Add(-48 * time.Hour).Unix()},
		PageReq:        entity.PageReq{},
	}
	resp, err := http.Post(fmt.Sprintf("%s/v1/logging/list", host), ContentTypeJSON, JSONReader(in))
	if err != nil {
		t.Fatal(err)
	}
	JSONOutput(resp, t)

}

func TestGetDBIndex(t *testing.T) {
	resp, err := http.Get(fmt.Sprintf("%s/v1/db-index", host))
	if err != nil {
		t.Fatal(err)
	}
	JSONOutput(resp, t)
}
