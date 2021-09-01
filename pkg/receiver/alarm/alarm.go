package alarm

import (
	"context"
	"fmt"
	"github.com/bbdshow/bkit/logs"
	"github.com/bbdshow/bkit/util/alert"
	"github.com/bbdshow/bkit/util/inet"
	"github.com/bbdshow/qelog/pkg/model"
	"go.uber.org/zap"
	"strings"
	"sync"
	"sync/atomic"
	"time"
)

var (
	ContentPrefix = "[QELOG]"
	machineIP, _  = inet.GetLocalIPV4()
)

type Alarm struct {
	mutex     sync.RWMutex
	ruleState map[string]*RuleState
	hooks     map[string]*model.HookURL
	modules   map[string]bool
	// 报警信息隐藏文字
	hideTexts []string
}

func NewAlarm() *Alarm {
	a := &Alarm{
		mutex:     sync.RWMutex{},
		ruleState: make(map[string]*RuleState, 0),
		hooks:     make(map[string]*model.HookURL, 0),
		modules:   make(map[string]bool),
		hideTexts: make([]string, 0),
	}
	return a
}

func (a *Alarm) AddHideText(txt []string) {
	for _, v := range txt {
		if v != "" {
			a.hideTexts = append(a.hideTexts, v)
		}
	}
}

// 如果模块没有设置报警，则不用判断具体的状态了
func (a *Alarm) ModuleIsEnable(name string) bool {
	a.mutex.RLock()
	enable, ok := a.modules[name]
	a.mutex.RUnlock()
	return ok && enable
}

func (a *Alarm) IsAlarm(docs []*model.Logging) {
	a.mutex.RLock()
	defer a.mutex.RUnlock()
	for _, v := range docs {
		state, ok := a.ruleState[v.Key()]
		if ok {
			state.Send(v)
		}
	}
}

func (a *Alarm) InitRuleState(rules []*model.AlarmRule, hooks []*model.HookURL) {
	modules := make(map[string]bool)
	ruleState := make(map[string]*RuleState, len(rules))
	hooksMap := make(map[string]*model.HookURL, len(hooks))
	for _, v := range hooks {
		hooksMap[v.ID.Hex()] = v
	}
	for _, rule := range rules {
		ruleState[rule.Key()] = new(RuleState).UpsertRule(rule, hooksMap[rule.HookID])
		modules[rule.ModuleName] = true
	}
	// 替换状态机
	a.mutex.Lock()
	defer a.mutex.Unlock()
	for _, state := range a.ruleState {
		v, ok := ruleState[state.Key()]
		if ok {
			ruleState[state.Key()] = state.UpsertRule(v.rule, v.hook)
		}
	}

	a.ruleState = ruleState
	a.modules = modules
	a.hooks = hooksMap
}

type RuleState struct {
	key            string
	hook           *model.HookURL
	rule           *model.AlarmRule
	count          int32
	latestSendTime int64
	method         alert.Alarm
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
		ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
		defer cancel()
		err := rs.method.Send(ctx, rs.parsingContent(v))
		if err != nil {
			logs.Qezap.Error("AlarmSend", zap.String(rs.method.Method(), err.Error()), zap.Any("content", rs.parsingContent(v)))
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

func (rs *RuleState) parsingContent(v *model.Logging) string {
	str := fmt.Sprintf(`%s
标签: %s
IP: %s
时间: %s
等级: %s
短消息: %s
详情: %s
频次: %d/%ds
报警节点: %s`, rs.KeyWord(), rs.rule.Tag, v.IP, time.Unix(v.TimeSec, 0).Format("2006-01-02 15:04:05"), v.Level.String(),
		v.Short, v.Full, atomic.LoadInt32(&rs.count), rs.rule.RateSec, machineIP)

	// 隐藏字段
	if rs.hook != nil {
		for _, hide := range rs.hook.HideText {
			str = strings.ReplaceAll(str, hide, "****")
		}
	}
	return str
}

func (rs *RuleState) Key() string {
	return rs.key
}
func (rs *RuleState) Rule() *model.AlarmRule {
	return rs.rule
}
func (rs *RuleState) KeyWord() string {
	if rs.hook != nil && rs.hook.KeyWord != "" {
		return rs.hook.KeyWord
	}
	return ContentPrefix
}

func (rs *RuleState) UpsertRule(new *model.AlarmRule, hook *model.HookURL) *RuleState {
	if rs.rule == nil || !rs.rule.UpdatedAt.Equal(new.UpdatedAt) {
		rs.rule = new
		rs.hook = hook
		rs.key = new.Key()
		rs.latestSendTime = 0
		switch rs.rule.Method {
		case model.MethodDingDing:
			rs.method = alert.NewDingDing()
			if rs.hook != nil {
				rs.method.SetHookURL(rs.hook.URL)
			}
		}
	}
	return rs
}
