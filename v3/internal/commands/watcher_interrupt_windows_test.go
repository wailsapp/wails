//go:build windows

package commands

import (
	"fmt"
	"os/exec"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestIsInterruptError(t *testing.T) {
	interruptErr := exec.Command("cmd", "/c", "exit /b -1073741510").Run()
	require.Error(t, interruptErr)
	assert.True(t, isInterruptError(fmt.Errorf("process exited: %w", interruptErr)))

	exitErr := exec.Command("cmd", "/c", "exit /b 1").Run()
	require.Error(t, exitErr)
	assert.False(t, isInterruptError(exitErr))
}
