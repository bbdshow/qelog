package types

import (
	"fmt"
	"testing"

	"github.com/huzhongqing/qelog/pkg/common/model"
)

func BenchmarkDecoder(b *testing.B) {

	for i := 0; i < b.N; i++ {
		str := `{"_level":"DEBUG","_time":1607961003768.121,"_caller":"wrapzap/wrap_zap_test.go:21","_func":"github.com/huzhongqing/qelog/wrapzap.TestNewWrapZap","_short":"Debug","k":"v","l":"fd6iStvg5A0U9apaZS7K"}`
		dec := Decoder{}
		if err := Unmarshal([]byte(str), &dec); err != nil {
			b.Fatal(err)
		}
		dec.TimeMill()
		dec.Short()
		dec.Level()
		dec.Condition(1)
		dec.TraceIDHex()
		dec.Full()
	}
}

func TestDecoder_Full(t *testing.T) {
	str := `{"_level":"DEBUG","_time":1607961003768.121,"_caller":"wrapzap/wrap_zap_test.go:21","_func":"github.com/huzhongqing/qelog/wrapzap.TestNewWrapZap","_short":"Debug","k":"v","l":"fd6iStvg5A0U9apaZS7K"}`
	dec := Decoder{}
	if err := Unmarshal([]byte(str), &dec); err != nil {
		t.Fatal(err)
	}
	r := model.Logging{
		Module:     "",
		IP:         "",
		Level:      dec.Level(),
		Short:      dec.Short(),
		Full:       dec.Full(),
		Condition1: "",
		Condition2: "",
		Condition3: "",
		TimeMill:   dec.TimeMill(),
		TimeSec:    0,
	}
	fmt.Println(r)
}
