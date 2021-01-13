package qezap

import (
	"os"
	"testing"
)

func TestBackupWrite(t *testing.T) {
	bw := NewBackupWrite("./testing/backup.log")
	_, _ = bw.WriteBakPacket([]byte("hello"))
	_, _ = bw.WriteBakPacket([]byte("world"))

	hello, err := bw.ReadBakPacket()
	if err != nil {
		t.Fatal(err)
	}
	if hello == nil {
		t.Fatal("not exists hello")
	}
	if string(hello) != "hello" {
		t.Fatal("hello not eq ", string(hello))
	}

	world, err := bw.ReadBakPacket()
	if err != nil {
		t.Fatal(err)
	}
	if world == nil {
		t.Fatal("not exists world")
	}

	if string(world) != "world" {
		t.Fatal("world not eq ", string(world))
	}
	_ = bw.Close()

	os.Remove(bw.filename)
}
