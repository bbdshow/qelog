package conf

import (
	"fmt"
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

	// 后台账号密码
	AdminUser AdminUser

	Admin    Admin
	Receiver Receiver
}

func InitConf(path ...string) error {
	return conf.ReadConfig(conf.FlagConfigPath(path...), Conf)
}

func (c *Config) Validate() error {

	mongoConn := map[string]struct{}{}
	for _, v := range c.Mongo.Conns {
		mongoConn[v.Database] = struct{}{}
	}
	if len(mongoConn) != len(c.Mongo.Conns) {
		return fmt.Errorf("mongo conns database must be different")
	}

	for db := range mongoConn {
		if !c.MongoGroup.IsExists(db) {
			return fmt.Errorf("mongo conns databse must be in the mongo group database")
		}
	}

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

func (mg MongoGroup) IsReceiverDatabase(database string) bool {
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
	HttpListenAddr string `defval:"0.0.0.0:31081"` // 如果 "" 则不开启 http 服务
	RpcListenAddr  string `defval:":31082"`
	AlarmEnable    bool   `defval:"true"`
	MetricsEnable  bool   `defval:"true"`
}

type Admin struct {
	HttpListenAddr string `defval:"0.0.0.0:31080"`
	AuthEnable     bool   `defval:"true"`
}

type AdminUser struct {
	Username string `default:"admin"`
	Password string `default:"111111"`
}
