package commands

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"golang.org/x/sys/windows"
)

func TestProtectHookContextPathUsesCurrentUserOnlyDACL(t *testing.T) {
	path := filepath.Join(t.TempDir(), "context.json")
	require.NoError(t, os.WriteFile(path, []byte("{}\n"), 0o600))
	require.NoError(t, protectHookContextPath(path, false))

	descriptor, err := windows.GetNamedSecurityInfo(
		path,
		windows.SE_FILE_OBJECT,
		windows.DACL_SECURITY_INFORMATION,
	)
	require.NoError(t, err)
	control, _, err := descriptor.Control()
	require.NoError(t, err)
	assert.NotZero(t, control&windows.SE_DACL_PROTECTED)
	dacl, defaulted, err := descriptor.DACL()
	require.NoError(t, err)
	assert.False(t, defaulted)
	assert.Equal(t, uint16(1), dacl.AceCount)
	user, err := windows.GetCurrentProcessToken().GetTokenUser()
	require.NoError(t, err)
	assert.True(t, strings.Contains(descriptor.String(), user.User.Sid.String()))
}
