package action

import (
	"errors"
)

// Action action的根数据 - 保存数据结构
type Action struct {
	actionName   string  // Action的名称
	rootData     string  // 初始化的数据 - json数据
	startEventID int64   // 起始版本ID
	events       []Event // 指令版本记录 DATA:Event数据
	eventLen     int64   // 记录长度
}

func NewAction(actionName string, rootData string, startEventID int64) *Action {
	return &Action{
		actionName:   actionName,
		rootData:     rootData,
		startEventID: startEventID,
		events:       make([]Event, 0),
		eventLen:     0,
	}
}

// AddEvent 记录一个Event
func (a *Action) AddEvent(event Event) (int64, error) {
	a.events = append(a.events, event)
	a.eventLen += 1
	return a.startEventID + a.eventLen, nil
}

// DeleteEvent 删除从EventID开始到最后的Event信息
func (a *Action) DeleteEvent(eventID int64) error {
	if eventID-a.startEventID < 0 || eventID-a.startEventID > a.eventLen {
		return errors.New("invalid event id")
	}
	a.events = append(a.events[:(eventID - a.startEventID - 1)])
	return nil
}

// GetMaxEventID 获取当前Action中最大EventID的数值
func (a *Action) GetMaxEventID() (int64, error) {
	return a.startEventID + a.eventLen, nil
}
