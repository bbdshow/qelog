package config

import (
	"errors"

	"github.com/huzhongqing/qelog/infra/defval"

	"github.com/BurntSushi/toml"
)

var Global *Config

func SetGlobalConfig(cfg *Config) {
	if cfg == nil {
		panic("config required")
	}
	Global = cfg
}

type Config struct {
	Env              string `default:"dev"`
	ReceiverAddr     string `default:"0.0.0.0:31081"`
	ReceiverGRPCAddr string `default:":31082"`
	ManagerAddr      string `default:"0.0.0.0:31080"`

	AuthEnable    bool `default:"true"`
	AlarmEnable   bool `default:"true"`
	MetricsEnable bool `default:"true"`

	AdminUser AdminUser
	// 日志配置，管理端产生的日志，也可以存储到远端
	Logging Logging

	// 储存管理配置的实例
	Main MongoDB
	// 存储日志内容的实例
	Sharding []MongoDBIndex

	// 不同的分片索引，可能存储在不同的数据库与集合里
	// 索引决定集合的命名
	MaxShardingIndex int32 `default:"8"`
	// 分片时间(天)跨度
	DaySpan int `default:"7"`
}

func InitConfig(filename string) *Config {
	cfg := &Config{}
	if err := defval.ParseDefaultVal(cfg); err != nil {
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

func (c *Config) Print() Config {
	cfg := *c
	// 脱敏
	if c.Release() {
		cfg.Main = MongoDB{}
		cfg.Sharding = []MongoDBIndex{}
	}
	return cfg
}

type MongoDBIndex struct {
	// 这个库需负责的存储序列
	Index    []int32 `default:"1,2,3,4,5,6,7,8"`
	DataBase string  `default:"sharding_qelog_db"`
	URI      string  `default:"mongodb://127.0.0.1:27017/admin"`
}

type MongoDB struct {
	DataBase string `default:"qelog"`
	URI      string `default:"mongodb://127.0.0.1:27017/admin"`
}

type AdminUser struct {
	Username string `default:"admin"`
	Password string `default:"111111"`
}

type Logging struct {
	Module   string   `default:"qelog"`
	Addr     []string `default:"127.0.0.1:31082"`
	Filename string   `default:"./log/logger.log"`
}
