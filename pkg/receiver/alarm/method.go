package alarm

import "context"

type Methoder interface {
	SetHookURL(string)
	Send(ctx context.Context, content string) error
	Method() string
}

type DingDingMethod struct {
}

func NewDingDingMethod() *DingDingMethod {
	return &DingDingMethod{}
}

func (ddm *DingDingMethod) SetHookURL(string) {}

func (ddm *DingDingMethod) Send(ctx context.Context, content string) error {
	return nil
}

func (ddm *DingDingMethod) Method() string {
	return "DingDing"
}
