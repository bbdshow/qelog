package defval

import (
	"fmt"
	"testing"
	"time"
)

type def struct {
	MyInt         int8              `default:"-1"`
	MyUint        uint              `default:"1"`
	MyString      string            `default:"hello"`
	MyBool        bool              `default:"true"`
	MyFloat       float32           `default:"66.6"`
	MySliceString []string          `default:"1,2,3"`
	MySliceFloat  []float64         `default:"66.6,77.7"`
	MySliceInt    []int8            `default:"-1,0,9"`
	MySliceUint   []uint16          `default:"0,2,4"`
	MyMap         map[string]string `default:"a=1,2"`
	MyDuration    time.Duration     `default:"30s"` // 30s
	MyStruct      MyStruct
	MyStruct2
}

type MyStruct struct {
	Key   string `default:"structKey"`
	Value MyStruct2
}

type MyStruct2 struct {
	Value int32 `default:"8"`
}

func TestParseDefaultVal(t *testing.T) {
	def := def{}
	if err := ParseDefaultVal(&def); err != nil {
		t.Fatal(err)
	}
	if !def.MyBool {
		t.Fatal("bool")
	}
	if def.MyDuration.String() != "30s" {
		fmt.Println(def.MyDuration.String())
		t.Fatal("duration")
	}
	if def.MyInt != -1 {
		t.Fatal("int")
	}
	if def.MyUint != 1 {
		t.Fatal("uint")
	}
	if def.MyFloat != 66.6 {
		t.Fatal("float")
	}
	if def.MyStruct.Value.Value != def.Value {
		t.Fatal("struct")
	}
	fmt.Printf("%#v \n", def)
}
