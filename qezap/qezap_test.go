package qezap

import (
	"testing"

	"go.uber.org/zap/zapcore"
)

func testLoggerNew(t *testing.T, opts ...Option) *Logger {
	lg := New(opts...)
	if lg == nil {
		t.Fatal("logger nil")
	}
	if lg.local == nil {
		t.Fatalf("local writer nil")
	}

	if lg.opt.EnableRemote {
		if lg.remote == nil {
			t.Fatalf("remote write nil")
		}
	}

	return lg
}

func TestLogger_New(t *testing.T) {
	local := testLoggerNew(t)
	local.Info("info")
	_ = local.Close()
	multi := testLoggerNew(t,
		WithAddrsAndModuleName([]string{"127.0.0.1:31082"}, "testing"),
		WithTransport(TransportMock),
	)
	multi.Info("multi info")
	_ = multi.Close()
}

func TestLogger_LevelWriter(t *testing.T) {
	lg := testLoggerNew(t,
		WithAddrsAndModuleName([]string{"127.0.0.1:31082"}, "testing"),
		WithTransport(TransportMock),
		WithFilename("./log/lw.log"),
	)
	lg.Info("multi info")
	defer lg.Close()

	ioWriter := lg.NewLevelWriter(zapcore.InfoLevel, "IOWriter")
	str := []byte("hello, this write used for io.Writer")
	n, err := ioWriter.Write(str)
	if err != nil {
		t.Fatalf("NewLevelWriter %v", err)
	}
	if n != len(str) {
		t.Fatalf("NewLevelWriter write length %d exception", n)
	}

	lg.SetEnabledLevel(zapcore.ErrorLevel)

	n, err = ioWriter.Write(str)
	if err == nil {
		t.Fatalf("NewLevelWriter should error")
	}
	t.Logf("change log level, write should error: [%v] n %d", err, n)
}
