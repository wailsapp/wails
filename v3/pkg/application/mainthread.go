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
		// handlePanic recovers and returns normally, so releasing the caller
		// from the closure's last statement would skip it on a panic and leave
		// wg.Wait() blocked forever. Deferred first, LIFO still runs
		// handlePanic to completion before the caller is let go.
		defer wg.Done()
		defer handlePanic()
		fn()
	})
	wg.Wait()
}

func InvokeSyncWithResult[T any](fn func() T) (res T) {
	var wg sync.WaitGroup
	wg.Add(1)
	globalApplication.dispatchOnMainThread(func() {
		defer wg.Done()
		defer handlePanic()
		res = fn()
	})
	wg.Wait()
	return res
}

func InvokeSyncWithError(fn func() error) (err error) {
	var wg sync.WaitGroup
	wg.Add(1)
	globalApplication.dispatchOnMainThread(func() {
		defer wg.Done()
		defer handlePanic()
		err = fn()
	})
	wg.Wait()
	return
}

func InvokeSyncWithResultAndError[T any](fn func() (T, error)) (res T, err error) {
	var wg sync.WaitGroup
	wg.Add(1)
	globalApplication.dispatchOnMainThread(func() {
		defer wg.Done()
		defer handlePanic()
		res, err = fn()
	})
	wg.Wait()
	return res, err
}

func InvokeSyncWithResultAndOther[T any, U any](fn func() (T, U)) (res T, other U) {
	var wg sync.WaitGroup
	wg.Add(1)
	globalApplication.dispatchOnMainThread(func() {
		defer wg.Done()
		defer handlePanic()
		res, other = fn()
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
