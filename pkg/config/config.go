package config

import (
	"errors"

	"github.com/huzhongqing/qelog/libs/defval"

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
	Env              string `default:"dev"`
	ReceiverAddr     string `default:"0.0.0.0:31081"`
	ReceiverGRPCAddr string `default:":31082"`
	ManagerAddr      string `default:"0.0.0.0:31080"`

	AuthEnable    bool `default:"true"`
	AlarmEnable   bool `default:"true"`
	MetricsEnable bool `default:"true"`

	// 不同的分片索引，可能存储在不同的数据库与集合里
	// 索引决定集合的命名
	MaxShardingIndex int32 `default:"8"`

	// 储存管理配置的实例
	Main MongoDB
	// 存储日志内容的实例
	Sharding []MongoDBIndex

	AdminUser AdminUser

	// 日志配置，管理端产生的日志，也可以存储到远端
	Logging Logging
}

func InitConfig(filename string) *Config {
	cfg := &Config{}
	if err := defval.ParseDefaultVal(cfg); err != nil {
		panic(err)
	}

	if err := defval.ParseDefaultVal(&cfg.AdminUser); err != nil {
		panic(err)
	}

	_, err := toml.DecodeFile(filename, cfg)
	if err != nil {
		panic("config init " + err.Error())
	}
	return cfg
}

func (c *Config) Release() bool {
	return c.Env == "release"
}

func (c *Config) Validate() error {
	if c.MaxShardingIndex <= 0 {
		c.MaxShardingIndex = 8
	}
	if c.Main.URI == "" {
		return errors.New("main.uri required")
	}

	if len(c.Sharding) <= 0 {
		return errors.New("sharding required")
	}
	indexExists := make(map[int32]struct{})
	for _, v := range c.Sharding {
		for _, i := range v.Index {
			_, ok := indexExists[i]
			if ok {
				return errors.New("sharding index dump key")
			}
			indexExists[i] = struct{}{}
		}
	}

	return nil
}

type MongoDBIndex struct {
	// 这个库需负责的存储序列
	Index    []int32
	DataBase string
	URI      string
}

type MongoDB struct {
	DataBase string
	URI      string
}

type AdminUser struct {
	Username string `default:"admin"`
	Password string `default:"111111"`
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

		AuthEnable:    false,
		AlarmEnable:   true,
		MetricsEnable: true,

		MaxShardingIndex: 8,

		Main: MongoDB{
			DataBase: "qelog",
			URI:      "mongodb://127.0.0.1:27017/admin",
		},
		Sharding: []MongoDBIndex{
			{
				Index:    []int32{1, 2, 3, 4},
				DataBase: "qelog",
				URI:      "mongodb://127.0.0.1:27017/admin",
			},
			{
				Index:    []int32{5, 6, 7, 8},
				DataBase: "qelog2",
				URI:      "mongodb://127.0.0.1:27017/admin",
			}},
		AdminUser: AdminUser{
			Username: "admin",
			Password: "111111",
		},
	}
}
