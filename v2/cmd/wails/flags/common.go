package flags

import (
	"fmt"
	"os"
	"runtime"
	"strings"

	"github.com/leaanthony/slicer"
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
		Platform: defaultPlatform + "/" + defaultArch,
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

func parseTargets(platform string) TargetsCollection {
	var result []Target
	var targets slicer.StringSlicer

	targets.AddSlice(strings.Split(platform, ","))
	targets.Deduplicate()

	targets.Each(func(platform string) {
		target := Target{
			Platform: "",
			Arch:     "",
		}

		platformSplit := strings.Split(platform, "/")

		target.Platform = platformSplit[0]

		if len(platformSplit) > 1 {
			target.Arch = platformSplit[1]
		} else {
			target.Arch = defaultTarget().Arch
		}

		result = append(result, target)
	})

	return result
}
