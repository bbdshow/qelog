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

func TestFindLoggingByTraceID(t *testing.T) {
	in := entity.FindLoggingByTraceIDReq{
		DBIndex:    7,
		ModuleName: "example",
		TraceID:    "1610201334536326300111441",
	}
	resp, err := http.Post(fmt.Sprintf("%s/v1/logging/traceid", host), ContentTypeJSON, JSONReader(in))
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

func TestCreateAlarmRule(t *testing.T) {
	in := &entity.CreateAlarmRuleReq{
		ModuleName: "example",
		Short:      "Alarm",
		Level:      3,
		Tag:        "[test]",
		RateSec:    30,
		Method:     1,
		HookURL:    "http://",
	}

	resp, err := http.Post(fmt.Sprintf("%s/v1/alarm-rule", host), ContentTypeJSON, JSONReader(in))
	if err != nil {
		t.Fatal(err)
	}
	JSONOutput(resp, t)
}

func TestUpdateAlarmRule(t *testing.T) {
	in := &entity.UpdateAlarmRuleReq{
		ObjectIDReq: entity.ObjectIDReq{ID: "5feac25b147eb46108e919a0"},
		Enable:      true,
		CreateAlarmRuleReq: entity.CreateAlarmRuleReq{
			ModuleName: "example",
			Short:      "Alarm",
			Level:      3,
			Tag:        "[test_ding]",
			RateSec:    30,
			Method:     1,
			HookURL:    "https://oapi.dingtalk.com/robot/send?access_token=00eca7373a1472267cc2a2a75ebab1ac476d3be37f3b7397d1f605b8d8e277b4",
		},
	}
	cli := &http.Client{}
	req, err := http.NewRequest("PUT", fmt.Sprintf("%s/v1/alarm-rule", host), JSONReader(in))
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Set("Content-Type", "application/json")
	resp, err := cli.Do(req)
	if err != nil {
		t.Fatal(err)
	}
	JSONOutput(resp, t)
}

func TestMetricsDBStats(t *testing.T) {
	resp, err := http.Get(fmt.Sprintf("%s/v1/metrics/dbstats", host))
	if err != nil {
		t.Fatal(err)
	}
	JSONOutput(resp, t)
}

func TestMetricsCollStats(t *testing.T) {
	resp, err := http.Get(fmt.Sprintf("%s/v1/metrics/collstats?host=%s&dbName=%s", host,
		"127.0.0.1:27017", "qelog"))
	if err != nil {
		t.Fatal(err)
	}
	JSONOutput(resp, t)
}

func TestMetricsModuleTrend(t *testing.T) {
	resp, err := http.Get(fmt.Sprintf("%s/v1/metrics/module/trend?moduleName=%s&lastDay=%d", host, "example", 7))
	if err != nil {
		t.Fatal(err)
	}
	JSONOutput(resp, t)
}
