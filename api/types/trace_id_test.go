package types

import (
	"fmt"
	"testing"
)

func TestTraceID(t *testing.T) {
	id := NewTraceID()
	fmt.Println(id.Hex(), id.Time().String())

	id2, err := TraceIDFromHex(id.Hex())
	if err != nil {
		t.Fatal(err)
	}

	if id.Hex() != id2.Hex() {
		t.Fatal("hex unequal")
	}

	if id.Time() != id2.Time() {
		t.Fatal("time unequal")
	}
}
