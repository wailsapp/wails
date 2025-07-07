package runtime

import (
	"context"
)

// OpenDevTools opens the developer tools window if runtime devtools support is enabled.
// This function only works when the application is built with the -runtimedevtools flag.
// On platforms where devtools cannot be opened programmatically (e.g., macOS in production),
// this function will do nothing.
func OpenDevTools(ctx context.Context) {
	appFrontend := getFrontend(ctx)
	if appFrontend != nil {
		appFrontend.OpenDevTools()
	}
}