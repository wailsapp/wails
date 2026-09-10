package main

import (
	"encoding/json"
	"os"
	"os/exec"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestRunSeparatesApplicationFlagsBeforeCLIParsing(t *testing.T) {
	cli, app := detachApplicationArguments([]string{"wails3", "run", "--tags", "devtools", "--", "--help", "space value", "--"})
	require.Equal(t, []string{"wails3", "run", "--tags", "devtools"}, cli)
	require.Equal(t, []string{"--help", "space value", "--"}, app)
	_, app = detachApplicationArguments([]string{"wails3", "run", "--"})
	require.NotNil(t, app)
	require.Empty(t, app)
	_, app = detachApplicationArguments([]string{"wails3", "run"})
	require.Nil(t, app)
	cli, app = detachApplicationArguments([]string{"wails3", "build", "--", "--help"})
	require.Len(t, cli, 4)
	require.Nil(t, app)
}

// Exercise the real CLI parser/action in a child process: a direct call to
// RunApplication cannot detect an empty slice introduced by flag parsing.
func TestRunCLIProcess(t *testing.T) {
	if os.Getenv("WAILS_RUN_CLI_TEST_CHILD") != "1" {
		return
	}
	for i, arg := range os.Args {
		if arg == "--" {
			os.Args = append([]string{"wails3"}, os.Args[i+1:]...)
			break
		}
	}
	main()
	os.Exit(0)
}

func TestRunCLIInheritsAndReplacesManifestArguments(t *testing.T) {
	root := t.TempDir()
	require.NoError(t, os.WriteFile(filepath.Join(root, "go.mod"), []byte("module example.com/runcli\n\ngo 1.25\n"), 0600))
	require.NoError(t, os.WriteFile(filepath.Join(root, "main.go"), []byte(`package main
import ("os"; "encoding/json")
func main() { data,_:=json.Marshal(os.Args[1:]); if err:=os.WriteFile("args.json",data,0600);err!=nil{panic(err)} }
`), 0600))
	require.NoError(t, os.WriteFile(filepath.Join(root, "wails.hcl"), []byte(`version = 3
project {
 name = "runcli"
 product_name = "Run CLI"
 identifier = "com.example.runcli"
 version = "1.0.0"
}
frontend { disabled = true }
run { args = ["default", "space value"] }
`), 0600))
	for _, tc := range []struct {
		name        string
		flags, want []string
	}{
		{"defaults", nil, []string{"default", "space value"}},
		{"replacement", []string{"--", "--help", "other value"}, []string{"--help", "other value"}},
		{"clear", []string{"--"}, []string{}},
		{"compatibility", []string{"--appargs", `--help "space value"`}, []string{"--help", "space value"}},
		{"compatibility clear", []string{"--appargs="}, []string{}},
	} {
		t.Run(tc.name, func(t *testing.T) {
			args := append([]string{"-test.run=^TestRunCLIProcess$", "--", "run"}, tc.flags...)
			cmd := exec.Command(os.Args[0], args...)
			cmd.Dir = root
			cmd.Env = append(os.Environ(), "WAILS_RUN_CLI_TEST_CHILD=1", "WAILS_EXP_USE_WAKE=1", "GOWORK=off")
			output, err := cmd.CombinedOutput()
			require.NoError(t, err, string(output))
			data, err := os.ReadFile(filepath.Join(root, "args.json"))
			require.NoError(t, err)
			var got []string
			require.NoError(t, json.Unmarshal(data, &got))
			require.Equal(t, tc.want, got)
		})
	}
}

func TestDevCLIApplicationArguments(t *testing.T) {
	root := t.TempDir()
	require.NoError(t, os.WriteFile(filepath.Join(root, "go.mod"), []byte("module example.com/devcli\n\ngo 1.25\n"), 0600))
	require.NoError(t, os.WriteFile(filepath.Join(root, "main.go"), []byte("package main\nfunc main(){}\n"), 0600))
	require.NoError(t, os.WriteFile(filepath.Join(root, "wails.hcl"), []byte(`version = 3
project {
 name = "devcli"
 product_name = "Dev CLI"
 identifier = "com.example.devcli"
 version = "1.0.0"
}
frontend { disabled = true }
run { args = ["production"] }
dev { args = ["--config-path", "testing.yaml"] }
`), 0600))
	for _, tc := range []struct {
		name        string
		flags, want []string
		failure     string
	}{
		{name: "defaults", want: []string{"--config-path", "testing.yaml"}},
		{name: "vector", flags: []string{"--", "--help", `C:\Users\Example User\testing.yaml`, ""}, want: []string{"--help", `C:\Users\Example User\testing.yaml`, ""}},
		{name: "clear vector", flags: []string{"--"}, want: []string{}},
		{name: "alias", flags: []string{"--appargs", `--config-path 'space value.yaml'`}, want: []string{"--config-path", "space value.yaml"}},
		{name: "single dash alias", flags: []string{"-appargs=--help"}, want: []string{"--help"}},
		{name: "clear alias", flags: []string{"--appargs="}, want: []string{}},
		{name: "mixed", flags: []string{"--appargs=", "--"}, failure: "either --appargs"},
		{name: "repeated", flags: []string{"--appargs=one", "--appargs=two"}, failure: "only once"},
		{name: "malformed", flags: []string{`--appargs="unterminated`}, failure: "quoting"},
	} {
		t.Run(tc.name, func(t *testing.T) {
			args := append([]string{"-test.run=^TestRunCLIProcess$", "--", "dev", "--plan"}, tc.flags...)
			cmd := exec.Command(os.Args[0], args...)
			cmd.Dir = root
			cmd.Env = append(os.Environ(), "WAILS_RUN_CLI_TEST_CHILD=1", "WAILS_EXP_USE_WAKE=1", "GOWORK=off")
			output, err := cmd.CombinedOutput()
			if tc.failure != "" {
				require.Error(t, err)
				require.Contains(t, string(output), tc.failure)
				return
			}
			require.NoError(t, err, string(output))
			encoded, err := json.Marshal(tc.want)
			require.NoError(t, err)
			require.Contains(t, string(output), "Application arguments: "+string(encoded))
		})
	}
}
