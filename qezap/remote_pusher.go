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

	"github.com/bbdshow/qelog/api/receiverpb"
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

type gRRCPush struct {
	cli   receiverpb.ReceiverClient
	conn  *grpc.ClientConn
	cChan chan struct{}
}

func newGRPCPush(addrs []string, concurrent int) (*gRRCPush, error) {
	if len(addrs) == 0 {
		return nil, fmt.Errorf("addrs required")
	}
	if concurrent <= 0 {
		concurrent = 5
	}

	resolver.Register(NewLocalResolverBuilder(addrs))

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	// warning: disable permission verify
	conn, err := grpc.DialContext(ctx, DialLocalServiceName,
		grpc.WithDefaultServiceConfig(fmt.Sprintf(`{"LoadBalancingPolicy": "%s"}`, roundrobin.Name)),
		grpc.WithInsecure())
	if err != nil {
		return nil, err
	}

	gp := &gRRCPush{
		cli:   receiverpb.NewReceiverClient(conn),
		conn:  conn,
		cChan: make(chan struct{}, concurrent),
	}

	return gp, nil
}

func (gp *gRRCPush) PushPacket(ctx context.Context, in *receiverpb.Packet) error {
	gp.cChan <- struct{}{}
	defer func() {
		<-gp.cChan
	}()

	resp, err := gp.cli.PushPacket(ctx, in)
	if err != nil {
		// any error, the server is considered unavailable
		log.Printf("Pusher:grpc %s\n", err)
		return ErrUnavailable
	}

	if resp.Code != 0 {
		return fmt.Errorf("response error %s", resp.String())
	}
	return nil
}

func (gp *gRRCPush) Concurrent() int {
	return len(gp.cChan)
}

func (gp *gRRCPush) Close() error {
	if gp.conn != nil {
		return gp.conn.Close()
	}
	return nil
}

type httpPush struct {
	addr   string
	client *http.Client

	cChan chan struct{}
}

func newHttpPush(addr []string, concurrent int) (*httpPush, error) {
	if len(addr) == 0 {
		return nil, fmt.Errorf("addr required")
	}
	if concurrent <= 0 {
		concurrent = 5
	}
	hp := &httpPush{
		addr:   addr[0],
		client: &http.Client{},
		cChan:  make(chan struct{}, concurrent),
	}

	return hp, nil
}

func (hp *httpPush) PushPacket(ctx context.Context, in *receiverpb.Packet) error {
	hp.cChan <- struct{}{}
	defer func() {
		<-hp.cChan
	}()
	v := struct {
		ID     string   `json:"id"`
		Module string   `json:"module"`
		Data   []string `json:"data"`
	}{ID: in.Id, Module: in.Module}
	// data split multi message and filter
	byteItems := bytes.Split(in.Data, []byte{'\n'})
	for _, b := range byteItems {
		if b == nil || bytes.Equal(b, []byte{}) || bytes.Equal(b, []byte{'\n'}) {
			continue
		}
		v.Data = append(v.Data, string(b))
	}

	return hp.push(ctx, v)
}

func (hp *httpPush) push(ctx context.Context, body interface{}) error {
	if ctx == nil {
		var cancel context.CancelFunc
		ctx, cancel = context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()
	}

	byt, err := json.Marshal(body)
	if err != nil {
		return err
	}

	req, err := http.NewRequestWithContext(ctx, "POST", hp.addr, bytes.NewReader(byt))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json")
	resp, err := hp.client.Do(req)
	if err != nil {
		log.Printf("Pusher:http %s\n", err)
		return ErrUnavailable
	}
	defer resp.Body.Close()

	if resp.StatusCode == 200 {
		return nil
	}
	respBody, _ := ioutil.ReadAll(resp.Body)
	return fmt.Errorf("http status code %d, response body %s", resp.StatusCode, string(respBody))
}

func (hp *httpPush) Concurrent() int {
	return len(hp.cChan)
}

func (hp *httpPush) Close() error {
	if hp.client != nil {
		hp.client.CloseIdleConnections()
	}
	return nil
}

// mockPush impl mock pusher used to test
type mockPush struct {
	cChan chan struct{}
}

func newMockPush(addrs []string, concurrent int) (*mockPush, error) {
	if len(addrs) == 0 {
		return nil, fmt.Errorf("addrs required")
	}
	if concurrent <= 0 {
		concurrent = 5
	}

	mp := &mockPush{
		cChan: make(chan struct{}, concurrent),
	}
	return mp, nil
}

var mockErr = 0

func (mp *mockPush) PushPacket(_ context.Context, in *receiverpb.Packet) error {
	mp.cChan <- struct{}{}
	defer func() {
		<-mp.cChan
	}()
	if string(in.Data) == ErrUnavailable.Error() && mockErr == 0 {
		log.Printf("Waiting Retry Data: ID %s DATA %s", in.Id, string(in.Data))
		mockErr = 1
		return ErrUnavailable
	}
	log.Printf("ID %s DATA %s", in.Id, string(in.Data))
	return nil
}

func (mp *mockPush) Concurrent() int {
	return len(mp.cChan)
}

func (mp *mockPush) Close() error {
	return nil
}
