package alarm

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

type Methoder interface {
	SetHookURL(string)
	Send(ctx context.Context, content string) error
	Method() string
}

type DingDingMethod struct {
	hookURL string
	client  *http.Client
}

func NewDingDingMethod() *DingDingMethod {
	return &DingDingMethod{
		client: &http.Client{},
	}
}

func (ddm *DingDingMethod) SetHookURL(url string) {
	ddm.hookURL = url
}

func (ddm *DingDingMethod) Send(ctx context.Context, content string) error {
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

	req, err := http.NewRequestWithContext(ctx, "POST", ddm.hookURL, bytes.NewReader(byt))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json")
	resp, err := ddm.client.Do(req)
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
		return fmt.Errorf("%s method error %d %s", ddm.Method(), respBody.ErrCode, respBody.ErrMsg)
	}
	return nil
}

func (ddm *DingDingMethod) Method() string {
	return "DingDing"
}
