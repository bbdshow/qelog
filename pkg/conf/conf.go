package conf

import (
	"fmt"
	"math/rand"

	"github.com/bbdshow/bkit/conf"
	"github.com/bbdshow/bkit/db/mongo"
	"github.com/bbdshow/bkit/gen/defval"
	"github.com/bbdshow/bkit/logs"
)

var (
	Conf = &Config{}
)

// Config required setting
// if config.toml does not exist this filed, it is set 'defval' value
type Config struct {
	Env string `defval:"dev"`

	MongoGroup MongoGroup
	Mongo      mongo.Config
	Logging    *logs.Config

	Admin    Admin
	Receiver Receiver
}

func InitConf(path ...string) error {
	return conf.ReadConfig(conf.FlagConfigPath(path...), Conf)
}

func (c *Config) Validate() error {
	if c.Admin.Username == "" || c.Admin.Password == "" {
		return fmt.Errorf("admin username password required")
	}

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
	cp := new(Config)
	*cp = *c
	conns := make([]mongo.Conn, 0)
	for _, conn := range cp.Mongo.Conns {
		conns = append(conns, conn)
	}
	cp.Mongo.Conns = conns
	_ = defval.InitialNullVal(cp)
	return *cp
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

func (mg MongoGroup) Databases() []string {
	names := make([]string, 0)
	names = append(names, mg.AdminDatabase)
	names = append(names, mg.ReceiverDatabase...)
	return names
}

type Receiver struct {
	HttpListenAddr string `defval:"0.0.0.0:31081"` // if empty, disable http server
	RpcListenAddr  string `defval:":31082"`
	AlarmEnable    bool   `defval:"true"`
	MetricsEnable  bool   `defval:"true"`
}

type Admin struct {
	HttpListenAddr string `defval:"0.0.0.0:31080"`
	AuthEnable     bool   `defval:"true"`
	Username       string `defval:"admin"` // manager: username/passwd
	Password       string `defval:"111111"`
}
