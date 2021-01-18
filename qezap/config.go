package qezap

import (
	"errors"
	"path"
	"strings"
	"time"
)

var (
	_logFilename    = "./log/logger.log"
	_backupFilename = path.Join(path.Dir(_logFilename), "backup", "backup.log")
)

type Config struct {
	// 本地文件
	Filename string
	// 超过大小 滚动 默认 0 不滚动
	MaxSize int64
	// 保留文件时间
	MaxAge time.Duration // 滚动日志文件最大时间， 默认 0 永久
	// Gzip 压缩 滚动日志是否Gzip压缩， 默认 false 不压缩
	GzipCompress bool

	// 是否开启远程传输
	EnableRemote bool
	// 远端传输协议 支持 HTTP、gRPC  默认 gRPC
	Transport string
	//  访问地址 HTTP["http://127.0.0.1:31081/v1/receiver/packet"] gRPC["127.0.0.1:31082"]
	//  HTTP 只取第一个地址， gRPC 取所有地址，然后轮询负载
	Addrs []string
	// 管理后台注册的模块名
	ModuleName string
	// 最大远端写入并发
	// 如果超出并发限制，直接写入备份文件，并间隔背景发送
	MaxConcurrent int // 默认 50 个并发
	// 最大数据包缓冲容量
	// 包的大小对写入效率有着比较重要的影响。 设置的相对大，有利于减少rpc调用次数，整体写入速度会更快。
	// 但是因为使用 sync.Pool 占用内存会更高一点。
	// 小对象对于GC与内存占用相对更加友好 (grpc 默认最大4MB一个包)
	MaxPacketSize int           // 默认 32kb
	WriteTimeout  time.Duration // 默认 5s
	// 远程数据包发送失败，备份到文件
	BackupFilename string
}

func NewConfig(addrs []string, moduleName string) *Config {
	cfg := &Config{
		Filename:     _logFilename,
		MaxSize:      500 << 20,
		MaxAge:       0,
		GzipCompress: true,

		EnableRemote:   false,
		Transport:      "grpc",
		Addrs:          addrs,
		ModuleName:     moduleName,
		MaxConcurrent:  50,
		MaxPacketSize:  32 << 10,
		WriteTimeout:   5 * time.Second,
		BackupFilename: _backupFilename,
	}
	if len(addrs) > 0 {
		cfg.EnableRemote = true
	}
	return cfg
}

func (cfg *Config) SetFilename(filename string) *Config {
	dir := path.Dir(filename)
	cfg.Filename = filename
	cfg.BackupFilename = path.Join(dir, "backup", "backup.log")
	return cfg
}

func (cfg *Config) SetEnableRemote(enable bool) *Config {
	cfg.EnableRemote = enable
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
