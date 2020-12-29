package qezap

import (
	"fmt"
	"testing"
)

func TestTraceID(t *testing.T) {
	var tid TraceID
	fmt.Println(tid.New().String(), tid.New().String())
	fmt.Println(tid.New().Time())
}
