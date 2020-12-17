package sharding

import (
	"context"
	"fmt"
	"sync"
	"time"
)

type Format string

const (
	FormatYear  Format = "2006"
	FormatMonth        = "200601"
	FormatDay          = "20060102"
)

type Sharding struct {
	mutex  sync.Mutex
	names  map[string]struct{}
	format string
	prefix string
}
type QueryNames func(ctx context.Context) ([]string, error)

func NewSharding(format Format, prefix ...string) *Sharding {
	s := &Sharding{
		names:  make(map[string]struct{}),
		format: string(format),
		prefix: "sharding",
	}
	if len(prefix) > 0 && prefix[0] != "" {
		s.prefix = prefix[0]
	}

	return s
}

func (s *Sharding) GenerateName(bucket string, unix int64) string {
	name := fmt.Sprintf("%s_%s_%s",
		s.prefix, bucket, time.Unix(unix, 0).Format(s.format))
	return name
}

func (s *Sharding) NameExists(ctx context.Context, name string, queryNames QueryNames) (bool, error) {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	_, ok := s.names[name]
	if ok {
		return true, nil
	}
	names, err := queryNames(ctx)
	if err != nil {
		return false, err
	}
	exists := false
	for _, v := range names {
		if v == name {
			exists = true
		}
		s.names[v] = struct{}{}
	}

	return exists, nil
}
