package wrapzap

import "context"

type Pusher interface {
	Push(ctx context.Context, b []byte) error
	Close() error
}
