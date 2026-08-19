package templates

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestCommonTemplateLeavesBuildConfigurationToTheManifestWriter(t *testing.T) {
	_, err := templates.ReadFile("_common/wails.tmpl.toml")
	require.Error(t, err)
	_, err = templates.ReadFile("_common/wails.tmpl.hcl")
	require.Error(t, err)
	_, err = templates.ReadFile("_common/Taskfile.tmpl.yml")
	require.Error(t, err)
}
