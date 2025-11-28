//go:build android

package application

// setupSingleInstance sets up single instance on Android
func (a *App) setupSingleInstance() error {
	// Android apps handle single instance via launch mode in manifest
	return nil
}

type androidLock struct {
	manager *singleInstanceManager
}

func newPlatformLock(manager *singleInstanceManager) (platformLock, error) {
	return &androidLock{
		manager: manager,
	}, nil
}

func (l *androidLock) acquire(uniqueID string) error {
	// Android apps handle single instance via launch mode in manifest
	return nil
}

func (l *androidLock) release() {
	// Android apps handle single instance via launch mode in manifest
}

func (l *androidLock) notify(data string) error {
	// Android apps handle single instance via launch mode in manifest
	return nil
}
