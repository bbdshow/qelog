package alarm

import (
	"context"
	"fmt"
	"sync"
	"sync/atomic"
	"time"

	"github.com/huzhongqing/qelog/pkg/common/kit"

	"github.com/huzhongqing/qelog/libs/logs"
	"go.uber.org/zap"

	"github.com/huzhongqing/qelog/pkg/common/model"
)

var (
	ContentPrefix = "[QELOG]"
	machineIP, _  = kit.GetLocalIPV4()
)

type Alarm struct {
	mutex     sync.RWMutex
	ruleState map[string]*RuleState
	modules   map[string]bool
}

func NewAlarm() *Alarm {
	a := &Alarm{
		mutex:     sync.RWMutex{},
		ruleState: make(map[string]*RuleState, 0),
		modules:   make(map[string]bool),
	}
	return a
}

// 如果模块没有设置报警，则不用判断具体的状态了
func (a *Alarm) ModuleIsEnable(name string) bool {
	a.mutex.RLock()
	enable, ok := a.modules[name]
	a.mutex.RUnlock()
	return ok && enable
}

func (a *Alarm) AlarmIfHitRule(docs []*model.Logging) {
	a.mutex.RLock()
	defer a.mutex.RUnlock()
	for _, v := range docs {
		state, ok := a.ruleState[v.Key()]
		if ok {
			state.Send(v)
		}
	}
}

func (a *Alarm) InitRuleState(rules []*model.AlarmRule) {
	modules := make(map[string]bool)
	ruleState := make(map[string]*RuleState, len(rules))
	for _, rule := range rules {
		ruleState[rule.Key()] = new(RuleState).UpsertRule(rule)
		modules[rule.ModuleName] = true
	}
	a.mutex.RLock()
	for _, state := range a.ruleState {
		v, ok := ruleState[state.Key()]
		if ok {
			ruleState[state.Key()] = state.UpsertRule(v.rule)
		}
	}
	a.mutex.RUnlock()
	// 替换状态机
	a.mutex.Lock()
	a.ruleState = ruleState
	a.modules = modules
	a.mutex.Unlock()
}

type RuleState struct {
	key            string
	rule           *model.AlarmRule
	count          int32
	latestSendTime int64
	method         Methoder
}

func (rs *RuleState) Send(v *model.Logging) {
	if v == nil {
		return
	}
	atomic.AddInt32(&rs.count, 1)
	isSend := false
	if atomic.LoadInt64(&rs.latestSendTime) == 0 {
		//直接发送
		isSend = true
	} else if time.Now().Unix()-atomic.LoadInt64(&rs.latestSendTime) > rs.rule.RateSec {
		// 超出了间隔
		isSend = true
	}
	if isSend {
		ctx, _ := context.WithTimeout(context.Background(), 30*time.Second)
		err := rs.method.Send(ctx, rs.content(v))
		if err != nil {
			logs.Qezap.Error("AlarmSend", zap.String(rs.method.Method(), err.Error()))
		} else {
			atomic.StoreInt32(&rs.count, 0)
			// 如果间隔时间 <= 0  那么每次都直接发送
			latestSendTime := int64(0)
			if rs.rule.RateSec > 0 {
				latestSendTime = time.Now().Unix()
			}
			atomic.StoreInt64(&rs.latestSendTime, latestSendTime)
		}
	} else {
		atomic.AddInt32(&rs.count, 1)
	}
	return
}

func (rs *RuleState) content(v *model.Logging) string {
	return fmt.Sprintf(`%s
标签: %s
IP: %s
时间: %s
等级: %s
短消息: %s
详情: %s
频次: %d/%ds
报警节点: %s`, ContentPrefix, rs.rule.Tag, v.IP, time.Unix(v.TimeSec, 0), v.Level.String(),
		v.Short, v.Full, atomic.LoadInt32(&rs.count), rs.rule.RateSec, machineIP)
}

func (rs *RuleState) Key() string {
	return rs.key
}
func (rs *RuleState) Rule() *model.AlarmRule {
	return rs.rule
}

func (rs *RuleState) UpsertRule(new *model.AlarmRule) *RuleState {
	if rs.rule == nil || rs.rule.UpdatedAt != new.UpdatedAt {
		rs.rule = new
		rs.key = new.Key()
		rs.latestSendTime = 0
		switch rs.rule.Method {
		case model.MethodDingDing:
			rs.method = NewDingDingMethod()
			rs.method.SetHookURL(rs.rule.HookURL)
		}
	}
	return rs
}
