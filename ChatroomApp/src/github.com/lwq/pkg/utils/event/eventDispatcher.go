package event

type IEventDispatcher interface {
	//注册事件
	AddEventHandler(eventName string, eventHandler EventHandler) (err error)
	//移除事件
	RemoveEventHandler(eventName string) (err error)
	//事件派发
	Invoke(o interface{})
	//删除所有事件
	RemoveAllEvent()
}
