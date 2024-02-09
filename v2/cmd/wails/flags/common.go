package flags

import (
	"fmt"
	"os"
	"runtime"
	"strings"

	"github.com/wailsapp/wails/v2/internal/system"
)

type Common struct {
	NoColour bool `description:"Disable colour in output"`
}

type Target struct {
	Platform string
	Arch     string
}

func defaultTarget() Target {
	defaultPlatform := os.Getenv("GOOS")
	if defaultPlatform == "" {
		defaultPlatform = runtime.GOOS
	}
	defaultArch := os.Getenv("GOARCH")
	if defaultArch == "" {
		if system.IsAppleSilicon {
			defaultArch = "arm64"
		} else {
			defaultArch = runtime.GOARCH
		}
	}

	return Target{
		Platform: defaultPlatform,
		Arch:     defaultArch,
	}
}

type TargetsCollection []Target

func (c TargetsCollection) MacTargetsCount() int {
	count := 0

	for _, t := range c {
		if strings.HasPrefix(t.Platform, "darwin") {
			count++
		}
	}

	return count
}

func (t Target) String() string {
	if t.Arch != "" {
		return fmt.Sprintf("%s/%s", t.Platform, t.Arch)
	}

	return t.Platform
}

func parseTargets(platforms string) TargetsCollection {
	platformList := strings.Split(platforms, ",")

	var targets []Target

	for _, platform := range platformList {
		parts := strings.Split(platform, "/")
		if len(parts) == 1 {
			architecture := defaultTarget().Arch
			targets = append(targets, Target{Platform: parts[0], Arch: architecture})
		} else if len(parts) == 2 {
			targets = append(targets, Target{Platform: parts[0], Arch: parts[1]})
		}
	}

	return targets
}
