package push

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"math/rand"
	"net/http"
	"strconv"
	"time"

	"google.golang.org/grpc/balancer/roundrobin"

	"google.golang.org/grpc/resolver"

	"google.golang.org/grpc"
)

func NewPacket(module string, data []string) *Packet {
	return &Packet{
		// 尽可能唯一ID, 后面随机3位，是防止多进程同一时刻
		Id:      strconv.FormatInt(time.Now().UnixNano()/1e6, 10) + "_" + strconv.FormatInt(rand.Int63n(1000), 10),
		Module:  module,
		Data:    data,
		IsRetry: false,
	}
}

type Pusher interface {
	PushPacket(ctx context.Context, in *Packet) error
	Concurrent() int
}

type GRRCPush struct {
	cli   PushClient
	conn  *grpc.ClientConn
	cChan chan struct{}
}

func NewGRPCPush(addrs []string, concurrent int) (*GRRCPush, error) {
	if len(addrs) == 0 {
		return nil, fmt.Errorf("addrs required")
	}
	if concurrent <= 0 {
		concurrent = 1
	}

	resolver.Register(NewLocalResolverBuilder(addrs))

	ctx, _ := context.WithTimeout(context.Background(), 5*time.Second)
	conn, err := grpc.DialContext(ctx, DialLocalServiceName, grpc.WithInsecure(), grpc.WithBlock(),
		grpc.WithDefaultServiceConfig(fmt.Sprintf(`{"LoadBalancingPolicy": "%s"}`, roundrobin.Name)))
	if err != nil {
		return nil, err
	}

	gp := &GRRCPush{
		cli:   NewPushClient(conn),
		conn:  conn,
		cChan: make(chan struct{}, concurrent),
	}

	return gp, nil
}

func (gp *GRRCPush) PushPacket(ctx context.Context, in *Packet) error {
	gp.cChan <- struct{}{}
	defer func() {
		<-gp.cChan
	}()
	if ctx == nil {
		ctx, _ = context.WithTimeout(context.Background(), 30*time.Second)
	}

	resp := make(chan error, 1)
	go func() {
		resp <- gp.push(ctx, in)
	}()

	select {
	case err := <-resp:
		return err
	case <-ctx.Done():
		return ctx.Err()
	}
}

func (gp *GRRCPush) push(ctx context.Context, in *Packet) error {
	resp, err := gp.cli.PushPacket(ctx, in)
	if err != nil {
		return err
	}
	if resp.Code != 0 {
		return fmt.Errorf("response error %s", resp.String())
	}
	return nil
}

func (gp *GRRCPush) Concurrent() int {
	return len(gp.cChan)
}

func (gp *GRRCPush) Close() error {
	if gp.conn != nil {
		return gp.conn.Close()
	}
	return nil
}

type HttpPush struct {
	addr   string
	client *http.Client

	cChan chan struct{}
}

func NewHttpPush(addr string, concurrent int) (*HttpPush, error) {
	if addr == "" {
		return nil, fmt.Errorf("addr required")
	}
	if concurrent <= 0 {
		concurrent = 1
	}
	hp := &HttpPush{
		addr:   addr,
		client: &http.Client{},
		cChan:  make(chan struct{}, concurrent),
	}

	return hp, nil
}

func (hp *HttpPush) PushPacket(ctx context.Context, in *Packet) error {
	hp.cChan <- struct{}{}
	defer func() {
		<-hp.cChan
	}()
	if ctx == nil {
		ctx = context.Background()
	}

	resp := make(chan error, 1)
	go func() {
		byt, err := json.Marshal(in)
		if err != nil {
			resp <- err
			return
		}
		resp <- hp.push(byt)
	}()

	select {
	case err := <-resp:
		return err
	case <-ctx.Done():
		return ctx.Err()
	}
}

func (hp *HttpPush) push(data []byte) error {
	contentType := "application/json"
	resp, err := hp.client.Post(hp.addr, contentType, bytes.NewReader(data))
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode == 200 {
		return nil
	}
	respBody, _ := ioutil.ReadAll(resp.Body)
	return fmt.Errorf("http status code %d, response body %s", resp.StatusCode, string(respBody))
}

func (hp *HttpPush) Concurrent() int {
	return len(hp.cChan)
}

func (hp *HttpPush) Close() error {
	if hp.client != nil {
		hp.client.CloseIdleConnections()
	}
	return nil
}
