package application

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestIOSOpenURLReturnsLaunchFailures(t *testing.T) {
	header, err := os.ReadFile("mobile_features_ios.h")
	require.NoError(t, err)
	implementation, err := os.ReadFile("mobile_features_ios.m")
	require.NoError(t, err)

	assert.Contains(t, string(header), "const char* ios_open_url(const char* url)")
	source := string(implementation)
	assert.Contains(t, source, "if (str == nil) return mfDup(@\"invalid URL encoding\")")
	assert.Contains(t, source, "completionHandler:^(BOOL success)")
	assert.Contains(t, source, "if (!success) errorMessage")
	assert.Contains(t, source, "timed out waiting for application to open URL")
	assert.NotContains(t, source, "canOpenURL:url")
}
