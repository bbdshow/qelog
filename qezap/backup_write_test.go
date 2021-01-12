package qezap

import (
	"fmt"
	"testing"
)

func TestBackupWrite_BakPacket(t *testing.T) {
	bw := NewBackupWrite("./testing/backup")
	bw.WriteBakPacket([]byte("hello"))
	bw.WriteBakPacket([]byte("world"))

	hello, err := bw.ReadBakPacket()
	if err != nil {
		t.Fatal(err)
	}
	if hello == nil {
		t.Fatal("not exists hello")
	}
	if string(hello) != "hello" {
		fmt.Println(string(hello))
		t.Fatal("hello not eq")
	}

	world, err := bw.ReadBakPacket()
	if err != nil {
		t.Fatal(err)
	}
	if world == nil {
		t.Fatal("not exists world")
	}

	if string(world) != "world" {
		t.Fatal("world not eq")
	}
}
