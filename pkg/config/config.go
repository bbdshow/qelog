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

	// 存储对象分片索引容量，可能存储在不同的数据库与集合里
	// 索引决定集合的命名
	ShardingIndexSize int `default:"8"`
	// 数据分片时间(天)跨度
	DaySpan int `default:"7"`
	// 后台账号密码
	AdminUser AdminUser

	// 日志配置，管理端产生的日志，也可以存储到远端
	Logging Logging

	// 管理配置的存储对象
	Main MongoMainDB
	// 日志内容的分片配置存储对象
	Sharding []MongoShardingDB
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
	if err := cfg.Validate(); err != nil {
		panic("config validate " + err.Error())
	}
	return cfg
}

func (c *Config) Release() bool {
	return c.Env == "release"
}

func (c *Config) Validate() error {
	if c.Main.URI == "" {
		return errors.New("main.uri required")
	}

	if len(c.Sharding) <= 0 {
		return errors.New("sharding required")
	}
	indexExists := make(map[int]struct{})
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
		cfg.Main = MongoMainDB{}
		cfg.Sharding = []MongoShardingDB{}
	}
	return cfg
}

type MongoShardingDB struct {
	// 这个库需负责的分片索引
	Index    []int  `default:"1,2,3,4,5,6,7,8"`
	DataBase string `default:"sharding_qelog_db"`
	URI      string `default:"mongodb://127.0.0.1:27017/admin"`
}

type MongoMainDB struct {
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
