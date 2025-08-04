//go:build darwin

package application

func (app *darwinApp) CreateJumpList() *JumpList {
	return &JumpList{
		app:        app,
		categories: []JumpListCategory{},
	}
}

func (j *JumpList) applyPlatform() error {
	// Stub implementation for macOS
	// Jump lists are Windows-specific, so this is a no-op on macOS
	return nil
}