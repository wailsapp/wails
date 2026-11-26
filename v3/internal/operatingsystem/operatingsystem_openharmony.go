//go:build openharmony

package operatingsystem

import (
	"fmt"
	"runtime"
)

func platformInfo() (*OS, error) {
	return &OS{
		ID:       "openharmony",
		Name:     "OpenHarmony",
		Version:  fmt.Sprintf("Go %s", runtime.Version()),
		Branding: "OpenHarmony",
	}, nil
}