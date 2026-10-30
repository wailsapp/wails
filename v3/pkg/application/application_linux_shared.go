//go:build linux && !android && !server

package application

func (a *linuxApp) GetFlags(options Options) map[string]any {
	if options.Flags == nil {
		options.Flags = make(map[string]any)
	}
	return options.Flags
}

func (a *linuxApp) destroy() {
	if !globalApplication.shouldQuit() {
		return
	}
	globalApplication.cleanup()
	appDestroy(a.application)
}

func (a *linuxApp) getAccentColor() string {
	return "rgb(0,122,255)"
}

func (a *linuxApp) hide() {
	a.hideAllWindows()
}

func (a *linuxApp) isOnMainThread() bool {
	return isOnMainThread()
}

func (a *linuxApp) name() string {
	return appName()
}

func (a *linuxApp) on(eventID uint) {
}

func (a *linuxApp) registerWindow(window pointer, id uint) {
	a.windowMapLock.Lock()
	a.windowMap[windowPointer(window)] = id
	a.windowMapLock.Unlock()
}

func (a *linuxApp) show() {
	a.showAllWindows()
}

func fatalHandler(errFunc func(error)) {
	// Stub for windows function
	return
}

func getNativeApplication() *linuxApp {
	return globalApplication.impl.(*linuxApp)
}

func writeInt(buf []byte, n int) int {
	if n < 0 {
		buf[0] = '-'
		return 1 + writeInt(buf[1:], -n)
	}
	if n == 0 {
		buf[0] = '0'
		return 1
	}
	tmp := n
	digits := 0
	for tmp > 0 {
		digits++
		tmp /= 10
	}
	for i := digits - 1; i >= 0; i-- {
		buf[i] = byte('0' + n%10)
		n /= 10
	}
	return digits
}
