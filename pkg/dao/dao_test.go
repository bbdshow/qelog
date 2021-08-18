package dao

import (
	"github.com/bbdshow/qelog/pkg/conf"
	"os"
	"testing"
)

var d *Dao

func TestMain(m *testing.M) {
	if err := conf.InitConf("../configs/config.toml"); err != nil {
		panic(err)
	}
	d = New(conf.Conf)
	os.Exit(m.Run())
}

func TestClose(t *testing.T) {
	d.Close()
}
