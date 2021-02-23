package alert

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

type Alarm interface {
	SetHookURL(string)
	Send(ctx context.Context, content string) error
	Method() string
}

type DingDing struct {
	hookURL string
	client  *http.Client
}

func NewDingDing() *DingDing {
	return &DingDing{
		client: &http.Client{},
	}
}

func (dd *DingDing) SetHookURL(url string) {
	dd.hookURL = url
}

func (dd *DingDing) Send(ctx context.Context, content string) error {
	reqBody := struct {
		MsgType string `json:"msgtype"`
		Text    struct {
			Content string `json:"content"`
		} `json:"text"`
	}{}
	respBody := struct {
		ErrCode int    `json:"errcode"`
		ErrMsg  string `json:"errmsg"`
	}{}
	reqBody.MsgType = "text"
	reqBody.Text.Content = content
	byt, err := json.Marshal(reqBody)
	if err != nil {
		return err
	}

	req, err := http.NewRequestWithContext(ctx, "POST", dd.hookURL, bytes.NewReader(byt))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json")
	resp, err := dd.client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	byt, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	if err := json.Unmarshal(byt, &respBody); err != nil {
		return err
	}

	if respBody.ErrCode != 0 {
		return fmt.Errorf("%s method error %d %s", dd.Method(), respBody.ErrCode, respBody.ErrMsg)
	}
	return nil
}

func (dd *DingDing) Method() string {
	return "DingDing"
}
