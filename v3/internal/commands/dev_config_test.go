package commands

import (
	"testing"

	"github.com/atterpac/refresh/process"
	"github.com/stretchr/testify/require"
	"gopkg.in/yaml.v3"
)

func TestDevConfigGatesApplicationOnFrontendReadiness(t *testing.T) {
	data, err := buildAssets.ReadFile("build_assets/config.yml")
	require.NoError(t, err)

	var config struct {
		DevMode struct {
			Executes []process.Execute `yaml:"executes"`
		} `yaml:"dev_mode"`
	}
	require.NoError(t, yaml.Unmarshal(data, &config))

	require.Len(t, config.DevMode.Executes, 4)
	require.Equal(t, process.Background, config.DevMode.Executes[1].Type)
	require.Equal(t, "wails3 task common:dev:wait", config.DevMode.Executes[2].Cmd)
	require.Equal(t, process.Once, config.DevMode.Executes[2].Type)
	require.Equal(t, process.Primary, config.DevMode.Executes[3].Type)
}

func TestFrontendDevServerTimeoutIsTaskLocal(t *testing.T) {
	data, err := buildAssets.ReadFile("build_assets/Taskfile.tmpl.yml")
	require.NoError(t, err)

	var taskfile struct {
		Tasks map[string]yaml.Node `yaml:"tasks"`
	}
	require.NoError(t, yaml.Unmarshal(data, &taskfile))

	waitTaskNode, exists := taskfile.Tasks["dev:wait"]
	require.True(t, exists)
	var waitTask struct {
		Vars map[string]int `yaml:"vars"`
		Cmds []string       `yaml:"cmds"`
	}
	require.NoError(t, waitTaskNode.Decode(&waitTask))
	require.Equal(t, 60, waitTask.Vars["TIMEOUT"])
	require.Equal(t, []string{"wails3 tool waitport --timeout {{.Opn}}.TIMEOUT{{.Cls}}"}, waitTask.Cmds)
}
