package types

import (
	"bytes"
	"fmt"
	"os"
	"strconv"
	"sync/atomic"
	"time"
)

var _pidString = func() string {
	pid := os.Getpid()
	return fmt.Sprintf("%05d", pid)
}()

var _incInt64 int64 = 0

type TraceID string

// [nsec:19]
func (tid TraceID) New() TraceID {
	buff := bytes.Buffer{}
	nsec := time.Now().UnixNano()
	nsecStr := strconv.FormatInt(nsec, 10)

	buff.WriteString(nsecStr)
	buff.WriteString(_pidString)
	buff.WriteString(strconv.FormatInt(atomic.AddInt64(&_incInt64, 1), 10))
	return TraceID(buff.String())
}

func (tid TraceID) Time() time.Time {
	if tid != "" && len(tid) >= 19 {
		nsec, _ := strconv.ParseInt(string(tid[:19]), 10, 64)
		return time.Unix(0, nsec)
	}
	return time.Unix(0, 0)
}

func (tid TraceID) String() string {
	return string(tid)
}
