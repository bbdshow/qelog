package wrapzap

import (
	"bytes"
	"context"
	"fmt"
	"io/ioutil"
	"net/http"
)

type Pusher interface {
	Push(ctx context.Context, data []byte) error
	Concurrent() int
}

type HttpPush struct {
	url    string
	client http.Client

	cChan chan struct{}
}

func NewHttpPush(url string, concurrent int) *HttpPush {
	if url == "" {
		panic("url required")
	}
	if concurrent <= 0 {
		concurrent = 1
	}
	hp := &HttpPush{
		url:    url,
		client: http.Client{},
		cChan:  make(chan struct{}, concurrent),
	}

	return hp
}

func (hp *HttpPush) Push(ctx context.Context, data []byte) error {
	hp.cChan <- struct{}{}
	defer func() {
		<-hp.cChan
	}()
	if ctx == nil {
		ctx = context.Background()
	}

	resp := make(chan error, 1)
	go func() {
		resp <- hp.push(data)
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
	resp, err := hp.client.Post(hp.url, contentType, bytes.NewReader(data))
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
