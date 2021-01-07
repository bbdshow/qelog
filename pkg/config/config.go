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
	Env              string
	ReceiverAddr     string
	ReceiverGRPCAddr string
	ManagerAddr      string

	AlarmEnable   bool
	MetricsEnable bool

	// 不同的模块存储不同的集合前缀里 (类似 redis db0-15 ...)
	DBSize int32

	MongoDB   MongoDB
	AdminUser AdminUser

	Logging Logging
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

type Logging struct {
	Module string
	Addr   []string
}

func MockDevConfig() *Config {
	return &Config{
		Env:              "dev",
		ReceiverAddr:     "0.0.0.0:31081",
		ReceiverGRPCAddr: ":31082",
		ManagerAddr:      "0.0.0.0:31080",

		AlarmEnable:   true,
		MetricsEnable: true,

		DBSize: 16,

		MongoDB: MongoDB{
			DataBase: "qelog",
			URI:      "mongodb://127.0.0.1:27017/admin",
		},
		AdminUser: AdminUser{
			Username: "admin",
			Password: "111111",
		},
	}
}
