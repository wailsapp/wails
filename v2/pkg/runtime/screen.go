package runtime

import (
	"context"

	"github.com/wailsapp/wails/v2/internal/frontend"
)

type Screen = frontend.Screen

// ScreenGetAll returns all screens
func ScreenGetAll(ctx context.Context) ([]Screen, error) {
	appFrontend := getFrontend(ctx)
	return appFrontend.ScreenGetAll()
}
