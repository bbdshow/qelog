package qezap

import (
	"testing"
	"time"

	"go.uber.org/zap/zapcore"

	"go.uber.org/zap"
)

func TestOption(t *testing.T) {

	var testCases = []struct {
		Filename      string
		Addr          []string
		ModuleName    string
		ZapOpt        []zap.Option
		Transport     Transport
		Level         zapcore.Level
		Mode          Mode
		MaxSize       uint64
		MaxAge        time.Duration
		Compress      bool
		MaxConcurrent uint
		MaxPacketSize uint
		Timeout       time.Duration
	}{
		{
			Filename:      "./log/x.log",
			Addr:          []string{},
			ModuleName:    "",
			ZapOpt:        []zap.Option{zap.Development()},
			Transport:     TransportMock,
			Level:         zap.ErrorLevel,
			Mode:          ModeRelease,
			MaxSize:       400,
			MaxAge:        10 * time.Second,
			Compress:      true,
			MaxConcurrent: 0,
			MaxPacketSize: 0,
			Timeout:       0,
		},
		{
			Filename:      "./log/remote.log",
			Addr:          []string{"127.0.0.1:31082"},
			ModuleName:    "testing",
			ZapOpt:        []zap.Option{zap.Development()},
			Transport:     TransportMock,
			Level:         zap.ErrorLevel,
			Mode:          ModeRelease,
			MaxSize:       400,
			MaxAge:        10 * time.Second,
			Compress:      true,
			MaxConcurrent: 5,
			MaxPacketSize: 5000,
			Timeout:       3 * time.Second,
		},
	}

	for i, v := range testCases {
		opt := defaultOptions()

		opts := []Option{
			WithFilename(v.Filename),
			WithAddrsAndModuleName(v.Addr, v.ModuleName),
			WithZapOptions(v.ZapOpt...),
			WithTransport(v.Transport),
			WithLevel(v.Level),
			WithMode(v.Mode),
			WithRotateMaxSizeAge(v.MaxSize, v.MaxAge),
			WithGzipCompress(v.Compress),
			WithRemoteConcurrent(v.MaxConcurrent),
			WithRemotePacketSize(v.MaxPacketSize),
			WithRemoteWriteTimeout(v.Timeout),
		}

		for _, v := range opts {
			v.apply(opt)
		}

		if opt.Local.Filename != v.Filename {
			t.Fatalf("case %d: opt %v", i, opt.Local)
		}
		if v.Level != opt.Level {
			t.Fatalf("case %d: opt %v", i, opt.Level)
		}
		if v.Mode != opt.Mode {
			t.Fatalf("case %d: opt %v", i, opt.Mode)
		}
		if v.Compress != opt.Local.GzipCompress {
			t.Fatalf("case %d: opt %v", i, opt.Local)
		}
		if int64(v.MaxSize) != opt.Local.MaxSize {
			t.Fatalf("case %d: opt %v", i, opt.Local)
		}

		if v.MaxAge != opt.Local.MaxAge {
			t.Fatalf("case %d: opt %v", i, opt.Local)
		}

		if len(v.Addr) > 0 {
			// remote
			if len(v.Addr) != len(opt.Remote.Addrs) || v.ModuleName != opt.Remote.ModuleName {
				t.Fatalf("case %d: opt %v", i, opt.Remote)
			}
			if v.Transport != opt.Remote.Transport {
				t.Fatalf("case %d: opt %v", i, opt.Remote)
			}
			if int(v.MaxPacketSize) != opt.Remote.MaxPacketSize {
				if v.MaxPacketSize > 4<<10 {
					if opt.Remote.MaxPacketSize != 4000<<10 {
						t.Fatalf("case %d: MaxPacketSize limit 4000<<10 %v", i, opt.Remote)
					}
				}
			}

			if int(v.MaxConcurrent) != opt.Remote.MaxConcurrent {
				t.Fatalf("case %d: opt %v", i, opt.Remote)
			}
		}

	}
}
