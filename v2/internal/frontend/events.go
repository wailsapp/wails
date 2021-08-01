package frontend

type Events interface {
	On(eventName string, callback func(...interface{}))
	OnMultiple(eventName string, callback func(...interface{}), counter int)
	Once(eventName string, callback func(...interface{}))
	Emit(eventName string, data ...interface{})
	Off(eventName string)
	Notify(name string, data ...interface{})
}
