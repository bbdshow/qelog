package cache

import (
	"fmt"
	"time"
)

type Cache interface {
	Get(key string) *Cmd
	Set(key string, value interface{}, ttl time.Duration) *StatusCmd

	// prefix - 前缀查询，"" 查询所有， 只返回当前有效的key
	Keys(prefix string) *SliceStringCmd

	Delete(key string) *StatusCmd

	// 删除所有 key
	FlushAll() *StatusCmd

	Save() *StatusCmd

	Close() error
}

type baseCmd struct {
	exists bool
	ttl    time.Duration
	err    error
}

func (cmd *baseCmd) Exists() bool {
	return cmd.exists
}

func (cmd *baseCmd) TTL() time.Duration {
	return cmd.ttl
}

func (cmd *baseCmd) Error() error {
	return cmd.err
}

type Cmd struct {
	baseCmd
	value interface{}
}

func (cmd *Cmd) ValString() string {
	switch cmd.value.(type) {
	case []byte:
		return string(cmd.value.([]byte))
	default:
		return fmt.Sprint(cmd.value)
	}
}

func (cmd *Cmd) Val() interface{} {
	return cmd.value
}

type SliceStringCmd struct {
	baseCmd
	value []string
}

func (cmd *SliceStringCmd) Val() []string {
	return cmd.value
}

const (
	StatusOK = "OK"
)

type StatusCmd struct {
	baseCmd
	value string
}

type BoolCmd struct {
	baseCmd
	value bool
}

func (cmd *BoolCmd) Val() bool {
	return cmd.value
}
