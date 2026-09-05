package commands

import (
	"os"
	"runtime"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/wailsapp/wails/v3/internal/flags"
)

func TestMergeTags(t *testing.T) {
	assert.Equal(t, "mcp", mergeTags("", "mcp"))
	assert.Equal(t, "gtk4,mcp", mergeTags("gtk4", "mcp"))
	assert.Equal(t, "gtk4,server,mcp", mergeTags("gtk4,server", "mcp"))
	assert.Equal(t, "mcp", mergeTags("mcp", "mcp"))
	assert.Equal(t, "gtk4,mcp", mergeTags("gtk4,mcp", "mcp"))
	assert.Equal(t, "gtk4", mergeTags("gtk4"))
	assert.Equal(t, "", mergeTags(""))
}

func TestEnvTags(t *testing.T) {
	for _, value := range []string{"1", "true", "TRUE", "on", "yes"} {
		t.Setenv(mcpEnvVar, value)
		assert.Equal(t, []string{mcpBuildTag}, envTags(), "WAILS_MCP=%s", value)
	}
	for _, value := range []string{"", "0", "false", "off", "no", "nonsense"} {
		t.Setenv(mcpEnvVar, value)
		assert.Empty(t, envTags(), "WAILS_MCP=%s", value)
	}
}

func TestBuildCommandWithMCPEnvVar(t *testing.T) {
	currentOS := runtime.GOOS
	currentArch := runtime.GOARCH

	originalRunTask := runTaskFunc
	defer func() { runTaskFunc = originalRunTask }()

	var capturedOptions *RunTaskOptions
	var capturedOtherArgs []string
	runTaskFunc = func(options *RunTaskOptions, otherArgs []string) error {
		capturedOptions = options
		capturedOtherArgs = otherArgs
		return nil
	}

	originalArgs := os.Args
	defer func() { os.Args = originalArgs }()

	t.Setenv("GOOS", "")
	t.Setenv("GOARCH", "")
	os.Unsetenv("GOOS")
	os.Unsetenv("GOARCH")
	t.Setenv(mcpEnvVar, "1")

	// WAILS_MCP=1 alone adds the mcp tag. `build` is a root-dispatch verb, so it
	// targets the root "build" task and passes GOOS/ARCH as variables (#5615).
	err := Build(&flags.Build{}, nil)
	assert.NoError(t, err)
	assert.Equal(t, "build", capturedOptions.Name)
	assert.Equal(t, []string{"EXTRA_TAGS=mcp", "GOOS=" + currentOS, "ARCH=" + currentArch}, capturedOtherArgs)

	// WAILS_MCP=1 merges with user-supplied tags
	buildFlags := &flags.Build{}
	buildFlags.Tags = "gtk4"
	err = Build(buildFlags, nil)
	assert.NoError(t, err)
	assert.Equal(t, []string{"EXTRA_TAGS=gtk4,mcp", "GOOS=" + currentOS, "ARCH=" + currentArch}, capturedOtherArgs)

	// No duplicate tag when the user already passed it
	buildFlags = &flags.Build{}
	buildFlags.Tags = "mcp"
	err = Build(buildFlags, nil)
	assert.NoError(t, err)
	assert.Equal(t, []string{"EXTRA_TAGS=mcp", "GOOS=" + currentOS, "ARCH=" + currentArch}, capturedOtherArgs)
}
