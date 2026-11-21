//go:build openharmony && !server

package application

// setupSingleInstance sets up single instance on OHOS.
func (a *App) setupSingleInstance() error {
	// OHOS apps handle single instance via the module.json5 launch type.
	return nil
}

type openharmonyLock struct {
	manager *singleInstanceManager
}

func newPlatformLock(manager *singleInstanceManager) (platformLock, error) {
	return &openharmonyLock{manager: manager}, nil
}

func (l *openharmonyLock) acquire(uniqueID string) error { return nil }
func (l *openharmonyLock) release()                      {}
func (l *openharmonyLock) notify(data string) error      { return nil }
