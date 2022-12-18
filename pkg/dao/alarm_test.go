package dao

import (
	"context"
	"testing"

	"github.com/bbdshow/bkit/tests"
	"github.com/bbdshow/qelog/pkg/model"
)

func TestDao_FindAlarmRuleList(t *testing.T) {
	in := &model.FindAlarmRuleListReq{
		Enable:     -1,
		ModuleName: "",
		Short:      "a",
		PageReq:    model.PageReq{Page: 1, Limit: 20},
	}
	_, docs, err := d.FindAlarmRuleList(context.Background(), in)
	if err != nil {
		t.Fatal(err)
	}
	tests.PrintBeautifyJSON(docs)
}
