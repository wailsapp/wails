package frontend

type Events interface {
	On(eventName string, callback func(...interface{})) func()
	OnMultiple(eventName string, callback func(...interface{}), counter int) func()
	Once(eventName string, callback func(...interface{})) func()
	Emit(eventName string, data ...interface{})
	Off(eventName string)
	OffAll()
	Notify(sender Frontend, name string, data ...interface{})
}
