package config

import (
	"github.com/BurntSushi/toml"
)

var GlobalConfig *Config

func SetGlobalConfig(cfg *Config) {
	if cfg == nil {
		panic("config required")
	}
	GlobalConfig = cfg
}

type Config struct {
	Env          string
	ReceiverAddr string
	QueryAddr    string
	MongoDB      MongoDB
}

func InitConfig(filename string) *Config {
	cfg := &Config{}
	_, err := toml.DecodeFile(filename, cfg)
	if err != nil {
		panic("config init " + err.Error())
	}
	return cfg
}

func (c *Config) Release() bool {
	return c.Env == "release"
}

type MongoDB struct {
	DataBase string
	URI      string
}

type AdminUser struct {
	Username string
	Password string
}
