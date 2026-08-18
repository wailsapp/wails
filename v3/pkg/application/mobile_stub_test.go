//go:build !ios && !android

package application

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestMobileOpenURLReportsUnsupportedPlatform(t *testing.T) {
	require.Error(t, Mobile.OpenURL("https://example.com"))
}
