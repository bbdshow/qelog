package cache

import (
	"os"
	"testing"
	"time"
)

func TestNewMemCache(t *testing.T) {
	memCache := NewMemCache(NewDefaultOptions())
	if setCmd := memCache.Set("2s", "true", 2*time.Second); setCmd.Error() != nil {
		t.Fatal(setCmd.Error())
	}
	getCmd := memCache.Get("2s")
	if getCmd.Error() != nil {
		t.Fatal(getCmd.Error())
	}
	if getCmd.ValString() != "true" {
		t.Fatal(getCmd.ValString())
	}

	time.Sleep(2 * time.Second)

	getCmd = memCache.Get("2s")
	if getCmd.Error() != nil {
		t.Fatal(getCmd.Error())
	}

	if getCmd.Exists() {
		t.Fatal("should del key")
	}
}

func TestMemCache_Size(t *testing.T) {
	opt := NewDefaultOptions()
	opt.Size = 32
	memCache := NewMemCache(opt)
	if setCmd := memCache.Set("2s", "true", 2*time.Second); setCmd.Error() != nil {
		t.Fatal(setCmd.Error())
	}

	if setCmd := memCache.Set("2s", "true", 2*time.Second); setCmd.Error() != nil {
		t.Fatal(setCmd.Error())
	}

	if setCmd := memCache.Set("30", "123456789012345678901234567890", 2*time.Second); setCmd.Error() == nil {
		t.Fatal("should over size")
	}
}

func TestMemCache_Save(t *testing.T) {
	opt := NewDefaultOptions()
	opt.Filename = "./cache.bak"
	memCache := NewMemCache(opt)
	if setCmd := memCache.Set("2s", "true", 2*time.Second); setCmd.Error() != nil {
		t.Fatal(setCmd.Error())
	}
	if err := memCache.Close(); err != nil {
		t.Fatal(err)
	}

	memCache1 := NewMemCache(opt)

	getCmd := memCache1.Get("2s")
	if getCmd.Error() != nil {
		t.Fatal(getCmd.Error())
	}
	if !getCmd.Exists() {
		t.Fatal("should 2s exists")
	}

	os.Remove("./cache.bak")
}
