package event

import (
	"errors"
	"reflect"
	"sync"
)

// 监听器函数
type EventHandler func(data interface{})

// 事件基类
type Event struct {
	lock          sync.RWMutex
	eventHandlers map[string][]EventHandler
}

func NewEvent() *Event {
	return &Event{
		eventHandlers: make(map[string][]EventHandler),
	}
}

// 监听事件
func (e *Event) AddEventHandler(eventName string, eventHandler EventHandler) (err error) {
	e.lock.Lock()
	defer e.lock.Unlock()
	e.eventHandlers[eventName] = append(e.eventHandlers[eventName], eventHandler)
	return nil
}

// 执行2023年3月20日09:21:50
func (e *Event) Invoke(data interface{}) {
	e.lock.Lock()
	defer e.lock.Unlock()
	//顺序调度事件
	for _, handlers := range e.eventHandlers {
		for _, handler := range handlers {
			handler(data)
		}
	}
}

// 移除事件处理
func (e *Event) RemoveEventHandler(eventName string, eventHandler EventHandler) (err error) {
	e.lock.Lock()
	defer e.lock.Unlock()
	eventList, ok := e.eventHandlers[eventName]
	if !ok {
		return errors.New(eventName + " Not registered ")
	}
	for i, event := range eventList {
		if reflect.ValueOf(event).Pointer() == reflect.ValueOf(eventHandler).Pointer() {
			eventList = append(eventList[:i], eventList[i+1:]...)
			e.eventHandlers[eventName] = eventList
			break
		}
	}
	return nil
}

// 删除所有事件
func (e *Event) RemoveAllEvent() {
	e.lock.Lock()
	defer e.lock.Unlock()
	e.eventHandlers = make(map[string][]EventHandler)
}
