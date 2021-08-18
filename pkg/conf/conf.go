package conf

import (
	"github.com/bbdshow/bkit/conf"
	"github.com/bbdshow/bkit/db/mongo"
	"github.com/bbdshow/bkit/gen/defval"
	"github.com/bbdshow/bkit/logs"
	"math/rand"
)

var (
	Conf = &Config{}
)

type Config struct {
	Env string `defval:"dev"`

	MongoGroup MongoGroup
	Mongo      *mongo.Config
	Logging    *logs.Config

	Receiver Receiver
}

func InitConf(path ...string) error {
	return conf.ReadConfig(conf.FlagConfigPath(path...), Conf)
}

func (c *Config) Validate() error {
	return nil
}

func (c *Config) Release() bool {
	return c.Env == "release"
}

func (c *Config) EraseSensitive() Config {
	// 脱敏数据，可用于打印
	cc := *c
	_ = defval.InitialNullVal(&cc)
	return cc
}

type MongoGroup struct {
	AdminDatabase    string
	ReceiverDatabase []string
}

func (mg MongoGroup) IsExists(database string) bool {
	if database == mg.AdminDatabase {
		return true
	}
	for _, v := range mg.ReceiverDatabase {
		if v == database {
			return true
		}
	}
	return false
}

func (mg MongoGroup) RandReceiverDatabase() string {
	if len(mg.ReceiverDatabase) == 0 {
		return ""
	}
	i := rand.Intn(len(mg.ReceiverDatabase))
	return mg.ReceiverDatabase[i]
}

type Receiver struct {
	RpcListenAddr string `defval:":31082"`
	AlarmEnable   bool   `defval:"true"`
	MetricsEnable bool   `defval:"true"`
}
