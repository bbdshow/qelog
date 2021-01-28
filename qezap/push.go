package qezap

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"time"

	"google.golang.org/grpc/balancer/roundrobin"
	"google.golang.org/grpc/resolver"

	"github.com/huzhongqing/qelog/api/receiverpb"
	"google.golang.org/grpc"
)

var (
	ErrUnavailable = errors.New("Push Unavailable")
)

type Pusher interface {
	PushPacket(ctx context.Context, in *receiverpb.Packet) error
	Concurrent() int
	Close() error
}

type GRRCPush struct {
	cli   receiverpb.ReceiverClient
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
	conn, err := grpc.DialContext(ctx, DialLocalServiceName, grpc.WithInsecure(),
		grpc.WithDefaultServiceConfig(fmt.Sprintf(`{"LoadBalancingPolicy": "%s"}`, roundrobin.Name)))
	if err != nil {
		return nil, err
	}

	gp := &GRRCPush{
		cli:   receiverpb.NewReceiverClient(conn),
		conn:  conn,
		cChan: make(chan struct{}, concurrent),
	}

	return gp, nil
}

func (gp *GRRCPush) PushPacket(ctx context.Context, in *receiverpb.Packet) error {
	gp.cChan <- struct{}{}
	defer func() {
		<-gp.cChan
	}()
	if ctx == nil {
		ctx, _ = context.WithTimeout(context.Background(), 10*time.Second)
	}

	if err := gp.push(ctx, in); err != nil {
		return err
	}

	return nil
}

func (gp *GRRCPush) push(ctx context.Context, in *receiverpb.Packet) error {
	resp, err := gp.cli.PushPacket(ctx, in)
	if err != nil {
		// 认为服务不可用
		log.Printf("grpc push %s\n", err)
		return ErrUnavailable
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

func (hp *HttpPush) PushPacket(ctx context.Context, in *receiverpb.Packet) error {
	hp.cChan <- struct{}{}
	defer func() {
		<-hp.cChan
	}()
	if ctx == nil {
		ctx, _ = context.WithTimeout(context.Background(), 10*time.Second)
	}
	v := struct {
		ID     string   `json:"id"`
		Module string   `json:"module"`
		Data   []string `json:"data"`
	}{ID: in.Id, Module: in.Module}
	byteItems := bytes.Split(in.Data, []byte{'\n'})
	for _, b := range byteItems {
		if b == nil || bytes.Equal(b, []byte{}) || bytes.Equal(b, []byte{'\n'}) {
			continue
		}
		v.Data = append(v.Data, string(b))
	}

	return hp.push(ctx, v)
}

func (hp *HttpPush) push(ctx context.Context, body interface{}) error {
	byt, err := json.Marshal(body)
	if err != nil {
		return err
	}

	contentType := "application/json"
	req, err := http.NewRequestWithContext(ctx, "POST", hp.addr, bytes.NewReader(byt))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", contentType)
	resp, err := hp.client.Do(req)
	if err != nil {
		// 认为服务不可用
		log.Printf("http push %s\n", err)
		return ErrUnavailable
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
