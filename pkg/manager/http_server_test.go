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
		fmt.Println(string(byt))
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
		Level:          0,
		IP:             "",
		ConditionOne:   "",
		ConditionTwo:   "",
		ConditionThree: "",
		TimeReq:        entity.TimeReq{BeginTsSec: time.Now().AddDate(0, 0, -13).Unix()},
		PageReq: entity.PageReq{
			Page:  1,
			Limit: 20,
		},
		ForceCollectionName: "logging_1_202102_03",
	}
	resp, err := http.Post(fmt.Sprintf("%s/v1/logging/list", host), ContentTypeJSON, JSONReader(in))
	if err != nil {
		t.Fatal(err)
	}
	JSONOutput(resp, t)

}

func TestFindLoggingByTraceID(t *testing.T) {
	in := entity.FindLoggingByTraceIDReq{
		DBIndex:    1,
		ModuleName: "example",
		TraceID:    "1666dcdd45c308587d4933fe",
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
