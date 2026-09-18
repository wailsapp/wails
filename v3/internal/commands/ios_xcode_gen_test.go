package commands

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestIOSXcodeProjectRetainsObjectiveCClasses(t *testing.T) {
	output := filepath.Join(t.TempDir(), "project.pbxproj")
	err := renderTemplateTo(
		updatableBuildAssets,
		"updatable_build_assets/ios/project.pbxproj.tmpl",
		output,
		iOSProjectConfig{ProductName: "Test App", ProductIdentifier: "com.wails.test"},
	)
	require.NoError(t, err)

	data, err := os.ReadFile(output)
	require.NoError(t, err)
	project := string(data)

	for _, configuration := range []string{"Debug", "Release"} {
		t.Run(configuration, func(t *testing.T) {
			marker := "/* " + configuration + " */ = {"
			_, remainder, found := strings.Cut(project, marker)
			require.True(t, found, "%s build configuration should be present", configuration)
			buildSettings, _, found := strings.Cut(remainder, "\n\t\t};")
			require.True(t, found, "%s build configuration should be complete", configuration)

			assert.Equal(t, 1, strings.Count(buildSettings, "OTHER_LDFLAGS = ("),
				"build configuration should define linker flags once")
			assert.Equal(t, 1, strings.Count(buildSettings, `"$(inherited)",`),
				"build configuration should preserve inherited linker flags")
			assert.Equal(t, 1, strings.Count(buildSettings, `"-ObjC",`),
				"build configuration should retain Objective-C classes from the Go archive")
		})
	}
}
