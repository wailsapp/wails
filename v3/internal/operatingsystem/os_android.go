//go:build android

package operatingsystem

import (
	"fmt"
	"runtime"
)

func platformInfo() (*OS, error) {
	return &OS{
		ID:       "android",
		Name:     "Android",
		Version:  fmt.Sprintf("Go %s", runtime.Version()),
		Branding: "Android",
	}, nil
}
