package metrics

import (
	"context"
	"github.com/bbdshow/bkit/logs"
	"github.com/bbdshow/qelog/common/types"
	"github.com/bbdshow/qelog/pkg/dao"
	"github.com/bbdshow/qelog/pkg/model"
	"sync"
	"time"

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
	mutex  sync.Mutex
	states map[string]*model.MetricsState
	d      *dao.Dao
}

func NewMetrics(d *dao.Dao) *Metrics {
	m := &Metrics{
		mutex:  sync.Mutex{},
		states: make(map[string]*model.MetricsState),
		d:      d,
	}
	return m
}

func initState(moduleName string) *model.MetricsState {
	now := time.Now()
	y, m, d := now.Date()
	date := time.Date(y, m, d, 0, 0, 0, 0, time.Local)
	section := time.Date(y, m, d, now.Hour(), 0, 0, 0, time.Local).Unix()

	return &model.MetricsState{
		Date:           date,
		Section:        section,
		ModuleName:     moduleName,
		Number:         0,
		Size:           0,
		Levels:         make(map[types.Level]int32),
		IPs:            make(map[string]int32),
		IncIntervalSec: incIntervalSec,
	}
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
	} else if state.IsIncr() {
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
func (m *Metrics) incr(s *model.MetricsState) {
	retry := 2
loop:
	if retry <= 0 {
		return
	}
	retry--
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	if err := m.d.IncrModuleMetrics(ctx, s); err != nil {
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
