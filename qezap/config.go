package qezap

import (
	"errors"
	"fmt"
	"path"
	"strings"
	"time"
)

type (
	Transport string
)

const (
	TransportGRPC Transport = "GRPC"
	TransportHTTP Transport = "HTTP"
	TransportMock Transport = "MOCK"
)

var (
	logFilename    = "./log/logger.log"
	backupFilename = path.Join(path.Dir(logFilename), "backup", "backup.log")
)

type Config struct {
	// local fs filename
	Filename string
	// single file max size, if 0 do not cut file. default: 500MB
	MaxSize int64
	// cut file max keep time. default: 0 keep forever
	MaxAge time.Duration
	// cut file enable Gzip compress. default: true
	GzipCompress bool
	// when addrs not empty, enable remote transport
	EnableRemote bool
	// remote transport protocol, support HTTP, GRPC, default: GRPC
	Transport Transport
	// remote endpoint address. eg HTTP:["http://xxx.com:31081/v1/receiver/packet"] GRPC:["192.168.10.1:31082","192.168.10.2:31082"]
	// HTTP suggest config load balancing address
	// GRPC qezap client impl local resolver, round_robin balancer
	Addrs []string
	// qelog admin register module name, equal to access token.
	ModuleName string
	// remote push max concurrent. if concurrent setting, data will be written backup file, bg retry send.
	// concurrent decision I/O max transfer. MAX=(MaxPacketSize*MaxConcurrent) default: 50
	MaxConcurrent int
	// send packet max size.  grpc client default body size max 4MB, but here default setting 32KB.
	// this setting 32KB + sync.Pool impl, can reduces memory overhead and friendly GC
	MaxPacketSize int
	// writeTimeout default 5s, if timeout, will be written backup file
	WriteTimeout time.Duration
	// back fs filename
	BackupFilename string
}

//func defaultConfig() *Config {
//	return &Config{
//		Filename:       logFilename,
//		MaxSize:        500 << 20,
//		MaxAge:         0,
//		GzipCompress:   true,
//		Transport:      TransportGRPC,
//		MaxConcurrent:  50,
//		MaxPacketSize:  32 << 10,
//		WriteTimeout:   5 * time.Second,
//		BackupFilename: backupFilename,
//	}
//}
//
//type OptionConfig interface {
//	apply(*Config)
//}
//
//type setOption struct {
//	f func(*Config)
//}
//
//func (s *setOption) apply(c *Config) {
//	s.f(c)
//}
//
//func newSetOption(f func(*Config)) *setOption {
//	return &setOption{
//		f: f,
//	}
//}
//
//// WithFilename setting
//logger filename
//func WithFilename(filename string) OptionConfig {
//	return newSetOption(func(c *Config) {
//		dir, file := path.Split(filename)
//		c.Filename = filename
//		c.BackupFilename = path.Join(dir, "backup", fmt.Sprintf("bak.%s", file))
//	})
//}

// Deprecated: use defaultConfig
func NewConfig(addrs []string, moduleName string) *Config {
	return &Config{
		Filename:       logFilename,
		MaxSize:        500 << 20,
		MaxAge:         0,
		GzipCompress:   true,
		EnableRemote:   len(addrs) > 0,
		Transport:      TransportGRPC,
		Addrs:          addrs,
		ModuleName:     moduleName,
		MaxConcurrent:  50,
		MaxPacketSize:  32 << 10,
		WriteTimeout:   5 * time.Second,
		BackupFilename: backupFilename,
	}
}

func (cfg *Config) SetFilename(filename string) *Config {
	dir, file := path.Split(filename)
	cfg.Filename = filename
	cfg.BackupFilename = path.Join(dir, "backup", fmt.Sprintf("bak.%s", file))
	return cfg
}

func (cfg *Config) SetEnableRemote(enable bool) *Config {
	cfg.EnableRemote = enable
	return cfg
}

func (cfg *Config) SetAddr(addr []string) *Config {
	cfg.Addrs = addr
	return cfg
}

func (cfg *Config) SetModule(module string) *Config {
	cfg.ModuleName = module
	return cfg
}

func (cfg *Config) SetHTTPTransport() *Config {
	cfg.Transport = "http"
	return cfg
}
func (cfg *Config) SetMaxSize(size int64) *Config {
	cfg.MaxSize = size
	return cfg
}

func (cfg *Config) SetMaxAge(t time.Duration) *Config {
	cfg.MaxAge = t
	return cfg
}

func (cfg *Config) SetWriteTimeout(t time.Duration) *Config {
	cfg.WriteTimeout = t
	return cfg
}

func (cfg *Config) SetMaxConcurrent(c int) *Config {
	cfg.MaxConcurrent = c
	return cfg
}

func (cfg *Config) SetMaxPacketSize(size int) *Config {
	cfg.MaxPacketSize = size
	return cfg
}

func (cfg *Config) Validate() error {
	if cfg.Filename == "" {
		return errors.New("filename required")
	}

	if cfg.MaxAge < 0 {
		cfg.MaxAge = 0
	}

	if cfg.EnableRemote {
		if cfg.ModuleName == "" {
			return errors.New("module name required")
		}
		if len(cfg.Addrs) == 0 {
			return errors.New("enable remote, addrs required")
		}
		for _, v := range cfg.Addrs {
			if cfg.Transport == "http" {
				if !strings.HasPrefix(v, "http://") {
					return errors.New("http addr invalid")
				}
			} else {
				ipPort := strings.Split(v, ":")
				if len(ipPort) < 2 || ipPort[1] == "" {
					return errors.New("gRPC addr invalid")
				}
			}
		}

		if cfg.MaxConcurrent <= 0 {
			cfg.MaxConcurrent = 1
		}
		if cfg.MaxPacketSize <= 0 {
			cfg.MaxPacketSize = 1 << 10
		}
		if cfg.WriteTimeout <= 0 {
			cfg.WriteTimeout = 5 * time.Second
		}
	}

	return nil
}
