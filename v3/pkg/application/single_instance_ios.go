//go:build ios

package application

// setupSingleInstance sets up single instance on iOS
func (a *App) setupSingleInstance() error {
	// iOS apps are always single instance
	return nil
}

type iosLock struct {
	manager *singleInstanceManager
}

func newPlatformLock(manager *singleInstanceManager) (platformLock, error) {
	return &iosLock{
		manager: manager,
	}, nil
}

func (l *iosLock) acquire(uniqueID string) error {
	// iOS apps are always single instance
	return nil
}

func (l *iosLock) release() {
	// iOS apps are always single instance
}

func (l *iosLock) notify(data string) error {
	// iOS apps are always single instance
	return nil
}