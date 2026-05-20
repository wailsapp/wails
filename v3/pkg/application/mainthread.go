package application

import (
	"sync"
)

var mainThreadFunctionStore = make(map[uint]func())
var mainThreadFunctionStoreLock sync.RWMutex

func generateFunctionStoreID() uint {
	startID := 0
	for {
		if _, ok := mainThreadFunctionStore[uint(startID)]; !ok {
			return uint(startID)
		}
		startID++
		if startID == 0 {
			Fatal("Too many functions have been dispatched to the main thread")
		}
	}
}

func InvokeSync(fn func()) {
	var wg sync.WaitGroup
	wg.Add(1)
	globalApplication.dispatchOnMainThread(func() {
		defer handlePanic()
		fn()
		wg.Done()
	})
	wg.Wait()
}

func InvokeSyncWithResult[T any](fn func() T) (res T) {
	var wg sync.WaitGroup
	wg.Add(1)
	globalApplication.dispatchOnMainThread(func() {
		defer handlePanic()
		res = fn()
		wg.Done()
	})
	wg.Wait()
	return res
}

func InvokeSyncWithError(fn func() error) (err error) {
	var wg sync.WaitGroup
	wg.Add(1)
	globalApplication.dispatchOnMainThread(func() {
		defer handlePanic()
		err = fn()
		wg.Done()
	})
	wg.Wait()
	return
}

func InvokeSyncWithResultAndError[T any](fn func() (T, error)) (res T, err error) {
	var wg sync.WaitGroup
	wg.Add(1)
	globalApplication.dispatchOnMainThread(func() {
		defer handlePanic()
		res, err = fn()
		wg.Done()
	})
	wg.Wait()
	return res, err
}

func InvokeSyncWithResultAndOther[T any, U any](fn func() (T, U)) (res T, other U) {
	var wg sync.WaitGroup
	wg.Add(1)
	globalApplication.dispatchOnMainThread(func() {
		defer handlePanic()
		res, other = fn()
		wg.Done()
	})
	wg.Wait()
	return res, other
}

func InvokeAsync(fn func()) {
	globalApplication.dispatchOnMainThread(func() {
		defer handlePanic()
		fn()
	})
}
