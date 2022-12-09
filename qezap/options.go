package qezap

import (
	"fmt"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"path"
	"time"
)

//type RemoteOption interface {
//	apply(*remoteOption)
//}

type remoteOption struct {
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

func defaultRemoteOption() *remoteOption {
	return &remoteOption{
		Transport:      TransportGRPC,
		MaxConcurrent:  50,
		MaxPacketSize:  32 << 10,
		WriteTimeout:   5 * time.Second,
		BackupFilename: backupFilename,
	}
}

type localOption struct {
	// local fs filename
	Filename string
	// single file max size, if 0 do not cut file. default: 500MB
	MaxSize int64
	// cut file max keep time. default: 0 keep forever
	MaxAge time.Duration
	// cut file enable Gzip compress. default: true
	GzipCompress bool
}

func defaultLocalOption() *localOption {
	return &localOption{
		Filename:     logFilename,
		MaxSize:      500 << 20,
		MaxAge:       0,
		GzipCompress: true,
	}
}

type options struct {
	// Uber-zap logger level
	Level zapcore.Level
	// Uber-zap Option
	Zap []zap.Option
	// when addrs not empty, enable remote transport
	EnableRemote bool

	Local  *localOption
	Remote *remoteOption
}

func defaultOptions() *options {
	return &options{
		Level:  zapcore.DebugLevel,
		Local:  defaultLocalOption(),
		Remote: defaultRemoteOption(),
	}
}

type Option interface {
	apply(*options)
}

type setOption struct {
	f func(*options)
}

func (s *setOption) apply(o *options) { s.f(o) }

func newSetOption(f func(*options)) *setOption { return &setOption{f: f} }

// WithFilename setting logger filename
func WithFilename(filename string) Option {
	return newSetOption(func(o *options) {
		dir, file := path.Split(filename)
		o.Local.Filename = file
		o.Remote.BackupFilename = path.Join(dir, "backup", fmt.Sprintf("bak.%s", file))
	})
}

func WithAddrs(addrs []string) Option {
	return newSetOption(func(o *options) {
		o.Remote.Addrs = addrs
	})
}

// WithZapOptions setting zap option
func WithZapOptions(opts ...zap.Option) Option {
	return newSetOption(func(o *options) {
		o.Zap = append(o.Zap, opts...)
	})
}
