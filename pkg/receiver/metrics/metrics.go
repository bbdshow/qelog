package metrics

import (
	"context"
	"fmt"
	"strings"
	"sync"
	"sync/atomic"
	"time"

	"github.com/huzhongqing/qelog/infra/logs"
	"github.com/huzhongqing/qelog/pkg/common/model"
	"github.com/huzhongqing/qelog/pkg/storage"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.uber.org/zap"
)

// 粗略的统计日志的写入情况， 因为当receiver异常推出的时候，可能不能及时入库
// 进行时统计，能够减少对Mongodb的资源消耗
// 如果在更新的时候因为网络问题导致失败，那么将直接忽略
var (
	incIntervalSec int64 = 60
)

func SetIncIntervalSec(sec int64) {
	incIntervalSec = sec
}

type Metrics struct {
	mutex              sync.Mutex
	states             map[string]*state
	moduleMetricsStore *storage.ModuleMetrics
}

func NewMetrics() *Metrics {
	m := &Metrics{
		mutex:              sync.Mutex{},
		states:             make(map[string]*state),
		moduleMetricsStore: storage.NewModuleMetrics(model.MainDB),
	}
	return m
}

type state struct {
	date       time.Time
	section    int64
	ModuleName string
	Number     int32
	Size       int32
	Levels     map[model.Level]int32
	IPs        map[string]int32
}

func initState(moduleName string) *state {
	now := time.Now()
	y, m, d := now.Date()
	date := time.Date(y, m, d, 0, 0, 0, 0, time.Local)
	section := time.Date(y, m, d, now.Hour(), 0, 0, 0, time.Local).Unix()

	return &state{
		date:       date,
		section:    section,
		ModuleName: moduleName,
		Number:     0,
		Size:       0,
		Levels:     make(map[model.Level]int32),
		IPs:        make(map[string]int32),
	}
}

func (s *state) IncrNumber(n int32) {
	atomic.AddInt32(&s.Number, n)
}
func (s *state) IncrSize(n int32) {
	atomic.AddInt32(&s.Size, n)
}

func (s *state) IncrLevel(lvl model.Level, n int32) {
	v, ok := s.Levels[lvl]
	if ok {
		s.Levels[lvl] = v + n
		return
	}
	s.Levels[lvl] = n
}

func (s *state) IncrIP(ip string, n int32) {
	if ip == "" {
		return
	}
	strs := strings.Split(ip, ".")
	if len(strs) <= 1 {
		// ipv6
		strs = strings.Split(ip, ":")
	}
	// 使用 _ 链接，便于mongodb更新
	ip = strings.Join(strs, "_")

	v, ok := s.IPs[ip]
	if ok {
		s.IPs[ip] = v + n
		return
	}
	s.IPs[ip] = n
}

func (s *state) isIncr() bool {
	// 超过一定时间，就可以写入了
	return time.Now().Unix()-s.section >= incIntervalSec
}

func (m *Metrics) Statistics(moduleName, ip string, docs []*model.Logging) {
	num := int32(len(docs))
	if num == 0 {
		return
	}

	m.mutex.Lock()
	defer m.mutex.Unlock()

	state, ok := m.states[moduleName]
	if !ok {
		state = initState(moduleName)
		m.states[moduleName] = state
	} else if state.isIncr() {
		// 先检查是否超过周期
		m.incr(state)
		state = initState(moduleName)
		m.states[moduleName] = state
	}

	state.IncrNumber(num)
	state.IncrIP(ip, num)
	for _, v := range docs {
		state.IncrSize(int32(v.Size))
		state.IncrLevel(v.Level, 1)
	}
}

// ignore update inc error
func (m *Metrics) incr(s *state) {
	filter := bson.M{
		"module_name":  s.ModuleName,
		"created_date": s.date,
	}

	opt := options.Update()
	opt.SetUpsert(true)

	fields := bson.M{
		"number": s.Number,
		"size":   s.Size,
		fmt.Sprintf("sections.%d.sum", s.section): s.Number,
	}
	for k, v := range s.Levels {
		fields[fmt.Sprintf("sections.%d.levels.%d", s.section, k.Int32())] = v
	}
	for k, v := range s.IPs {
		fields[fmt.Sprintf("sections.%d.ips.%s", s.section, k)] = v
	}

	update := bson.M{
		"$inc": fields,
	}

	retry := 2
loop:
	if retry <= 0 {
		return
	}
	retry--
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	if err := m.moduleMetricsStore.UpdateModuleMetrics(ctx, filter, update, opt); err != nil {
		logs.Qezap.Error("Metrics", zap.String("UpsertModuleMetrics", err.Error()))
		cancel()
		goto loop
	}
	cancel()
}

func (m *Metrics) Sync() {
	m.mutex.Lock()
	for moduleName, state := range m.states {
		m.incr(state)
		m.states[moduleName] = initState(moduleName)
	}
	m.mutex.Unlock()
}
