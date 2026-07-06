package runtime

import (
	"context"
)

// EnvironmentInfo mirrors the v2 runtime.EnvironmentInfo type.
// v3 equivalent: application.EnvironmentInfo.
type EnvironmentInfo struct {
	BuildType string `json:"buildType"`
	Platform  string `json:"platform"`
	Arch      string `json:"arch"`
}

// Quit mirrors the v2 runtime.Quit function.
// v3 equivalent: app.Quit.
func Quit(_ context.Context) {
	if a := app(); a != nil {
		a.Quit()
	}
}

// Hide mirrors the v2 runtime.Hide function.
// v3 equivalent: app.Hide.
func Hide(_ context.Context) {
	if a := app(); a != nil {
		a.Hide()
	}
}

// Show mirrors the v2 runtime.Show function.
// v3 equivalent: app.Show.
func Show(_ context.Context) {
	if a := app(); a != nil {
		a.Show()
	}
}

// Environment mirrors the v2 runtime.Environment function.
// v3 equivalent: app.Env.Info.
func Environment(_ context.Context) EnvironmentInfo {
	a := app()
	if a == nil {
		return EnvironmentInfo{}
	}
	info := a.Env.Info()
	buildType := "production"
	if info.Debug {
		buildType = "dev"
	}
	return EnvironmentInfo{
		BuildType: buildType,
		Platform:  info.OS,
		Arch:      info.Arch,
	}
}
