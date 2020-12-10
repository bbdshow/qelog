package wrapzap

import "time"

type WriteRemoteConfig struct {
	Endpoint   string
	Token      string
	ModuleName string

	MaxPool      int
	MaxPacket    int
	WriteTimeout time.Duration
}

type WriteRemote struct {
	pusher Pusher
}
