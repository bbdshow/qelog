package alarm

import "github.com/huzhongqing/qelog/pkg/common/model"

type Alarm struct {
	ruleState map[string]State
}

func (a *Alarm) ModuleIsEnable(name string) bool {
	return false
}
func (a *Alarm) AlarmIfHitRule([]*model.Logging) error {
	return nil
}

type State struct {
	rule          model.AlarmRule
	latestLogging *model.Logging
	nextSendTime  int64
	method        Methoder
}

func (s *State) IsSend(v *model.Logging) bool {
	return nil
}
