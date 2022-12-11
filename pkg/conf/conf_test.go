package conf

import (
	"testing"

	"github.com/bbdshow/bkit/conf"
)

func TestConfigToFile(t *testing.T) {
	cfg := &Config{}
	if err := conf.UnmarshalDefaultVal(cfg); err != nil {
		t.Fatal(err)
	}
	if err := conf.MarshalToFile(cfg, "../configs/config.toml"); err != nil {
		t.Fatal(err)
	}
}
