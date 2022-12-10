package qezap

import (
	"go.uber.org/multierr"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func newOneEncoderMultiWriterCore(enc zapcore.Encoder, level *zap.AtomicLevel, multiWriter []zapcore.WriteSyncer) *oneEncoderMultiWriter {
	return &oneEncoderMultiWriter{
		AtomicLevel: level,
		enc:         enc,
		multiWriter: multiWriter,
	}
}

// impl zapcore.Core,  one encoder, multi writer
type oneEncoderMultiWriter struct {
	*zap.AtomicLevel
	enc         zapcore.Encoder
	multiWriter []zapcore.WriteSyncer
}

func (mw *oneEncoderMultiWriter) With(fields []zap.Field) zapcore.Core {
	clone := mw.clone()
	for i := range fields {
		fields[i].AddTo(clone.enc)
	}
	return clone
}

func (mw *oneEncoderMultiWriter) Check(ent zapcore.Entry, ce *zapcore.CheckedEntry) *zapcore.CheckedEntry {
	if mw.Enabled(ent.Level) {
		return ce.AddCore(ent, mw)
	}
	return ce
}

func (mw *oneEncoderMultiWriter) Write(ent zapcore.Entry, fields []zap.Field) error {
	buf, err := mw.enc.EncodeEntry(ent, fields)
	if err != nil {
		return err
	}

	for _, w := range mw.multiWriter {
		_, err = w.Write(buf.Bytes())
		if err != nil {
			err = multierr.Append(err, err)
			continue
		}
	}
	buf.Free()

	if ent.Level > zapcore.ErrorLevel {
		// Since we may be crashing the program, sync the output. Ignore Sync
		// errors, pending a clean solution to issue #370.
		_ = mw.Sync()
	}
	return err
}

func (mw *oneEncoderMultiWriter) Sync() error {
	var err error
	for i := range mw.multiWriter {
		err = multierr.Append(err, mw.multiWriter[i].Sync())
	}
	return err
}

func (mw *oneEncoderMultiWriter) clone() *oneEncoderMultiWriter {
	return &oneEncoderMultiWriter{
		AtomicLevel: mw.AtomicLevel,
		enc:         mw.enc.Clone(),
		multiWriter: mw.multiWriter,
	}
}
