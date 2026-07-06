package runtime

import (
	"context"
)

// ScreenSize mirrors the v2 runtime.ScreenSize type.
type ScreenSize struct {
	Width  int `json:"width"`
	Height int `json:"height"`
}

// Screen mirrors the v2 runtime.Screen type.
// v3 equivalent: application.Screen.
type Screen struct {
	IsCurrent    bool       `json:"isCurrent"`
	IsPrimary    bool       `json:"isPrimary"`
	Width        int        `json:"width"`
	Height       int        `json:"height"`
	Size         ScreenSize `json:"size"`
	PhysicalSize ScreenSize `json:"physicalSize"`
}

// ScreenGetAll mirrors the v2 runtime.ScreenGetAll function.
// v3 equivalent: app.Screen.GetAll.
func ScreenGetAll(_ context.Context) ([]Screen, error) {
	a := app()
	if a == nil {
		return nil, errNoApp
	}

	currentID := ""
	if w := currentWindow(); w != nil {
		if screen, err := w.GetScreen(); err == nil && screen != nil {
			currentID = screen.ID
		}
	}

	screens := a.Screen.GetAll()
	result := make([]Screen, 0, len(screens))
	for _, screen := range screens {
		result = append(result, Screen{
			IsCurrent: currentID != "" && screen.ID == currentID,
			IsPrimary: screen.IsPrimary,
			Width:     screen.Size.Width,
			Height:    screen.Size.Height,
			Size: ScreenSize{
				Width:  screen.Size.Width,
				Height: screen.Size.Height,
			},
			PhysicalSize: ScreenSize{
				Width:  screen.PhysicalBounds.Width,
				Height: screen.PhysicalBounds.Height,
			},
		})
	}
	return result, nil
}
