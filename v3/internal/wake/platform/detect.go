package platform

import "runtime"

func OS() string {
	return runtime.GOOS
}

func Arch() string {
	return runtime.GOARCH
}

func OSFamily() string {
	switch runtime.GOOS {
	case "linux", "darwin", "freebsd", "openbsd", "netbsd":
		return "unix"
	case "windows":
		return "windows"
	default:
		return runtime.GOOS
	}
}

func NumCPU() int {
	return runtime.NumCPU()
}

func ExeExt() string {
	if runtime.GOOS == "windows" {
		return ".exe"
	}
	return ""
}

func Matches(os string) bool {
	return runtime.GOOS == os
}

func Filter(platforms []string) bool {
	if len(platforms) == 0 {
		return true
	}
	for _, p := range platforms {
		if runtime.GOOS == p {
			return true
		}
	}
	return false
}
