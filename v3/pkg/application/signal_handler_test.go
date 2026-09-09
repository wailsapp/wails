package application

import (
	"go/build"
	"testing"
)

// TestSignalHandlerBuildVariants keeps mobile lifecycle handling separate from
// the desktop/server Unix signal handler.
func TestSignalHandlerBuildVariants(t *testing.T) {
	for _, platform := range []string{"darwin", "linux", "windows", "ios", "android"} {
		for _, mode := range []string{"desktop", "server"} {
			t.Run(platform+"/"+mode, func(t *testing.T) {
				ctx := build.Default
				ctx.GOOS, ctx.GOARCH, ctx.CgoEnabled = platform, "arm64", true
				ctx.BuildTags = nil
				if mode == "server" {
					ctx.BuildTags = []string{"server"}
				}
				for file, want := range map[string]bool{
					"signal_handler_desktop.go": platform != "ios" && platform != "android",
					"signal_handler_ios.go":     platform == "ios",
					"signal_handler_android.go": platform == "android",
				} {
					got, err := ctx.MatchFile(".", file)
					if err != nil {
						t.Fatal(err)
					}
					if got != want {
						t.Errorf("%s selected = %v, want %v", file, got, want)
					}
				}
			})
		}
	}
}
