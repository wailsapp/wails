//go:build linux

package application

func (app *linuxApp) CreateJumpList() *JumpList {
	return &JumpList{
		app:        app,
		categories: []JumpListCategory{},
	}
}

func (j *JumpList) applyPlatform() error {
	// Stub implementation for Linux
	// Jump lists are Windows-specific, so this is a no-op on Linux
	return nil
}