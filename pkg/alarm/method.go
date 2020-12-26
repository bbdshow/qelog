package alarm

import "context"

type Methoder interface {
	SetHookURL(string)
	Send(ctx context.Context, content string) error
}
