package qezap

import (
	"os"
	"testing"
)

func TestWriteBackup_WriteRead(t *testing.T) {
	w := newWriteBackup("./log/backup.log")
	_, err := w.WriteBakPacket([]byte("hello"))
	if err != nil {
		t.Fatal(err)
	}
	_, err = w.WriteBakPacket([]byte("world"))
	if err != nil {
		t.Fatal(err)
	}

	hello, err := w.ReadBakPacket()
	if err != nil {
		t.Fatal(err)
	}
	if hello == nil {
		t.Fatal("not found hello")
	}
	if string(hello) != "hello" {
		t.Fatal("data not equal ", string(hello))
	}

	world, err := w.ReadBakPacket()
	if err != nil {
		t.Fatal(err)
	}
	if world == nil {
		t.Fatal("not found world")
	}

	if string(world) != "world" {
		t.Fatal("data not equal ", string(world))
	}
	_ = w.Close()

	_ = os.Remove(w.filename)
}
