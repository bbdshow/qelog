package types

import (
	"bytes"
	"crypto/rand"
	"encoding/binary"
	"encoding/hex"
	"errors"
	"fmt"
	"io"
	"time"
)

var ErrInvalidHex = errors.New("the provided hex string is not a valid TraceID")

var processUnique = func() [4]byte {
	var b [4]byte
	_, err := io.ReadFull(rand.Reader, b[:])
	if err != nil {
		panic(fmt.Sprintf("connot init processUnique: %v", err))
	}
	return b
}()

var NilTraceID TraceID

type TraceID [12]byte

func NewTraceID() TraceID {
	var b [12]byte
	binary.BigEndian.PutUint64(b[0:8], uint64(time.Now().UnixNano()))
	copy(b[8:12], processUnique[:])
	return b
}

func (id TraceID) Time() time.Time {
	nsec := binary.BigEndian.Uint64(id[0:8])
	return time.Unix(0, int64(nsec))
}

func (id TraceID) Hex() string {
	return hex.EncodeToString(id[:])
}

func (id TraceID) IsZero() bool {
	return bytes.Equal(id[:], NilTraceID[:])
}

func TraceIDFromHex(s string) (TraceID, error) {
	b, err := hex.DecodeString(s)
	if err != nil {
		return NilTraceID, err
	}
	if len(b) != 12 {
		return NilTraceID, ErrInvalidHex
	}
	var id [12]byte
	copy(id[:], b[:])
	return id, nil
}
