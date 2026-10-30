//go:build !windows

package cmds

import (
	"bytes"
	"errors"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
)

func TestRouteFrontendArguments(t *testing.T) {
	for _, tc := range []struct {
		command string
		want    []string
	}{
		{`npm install "./local package" --ignore-scripts`, []string{"install", "./local package", "--ignore-scripts"}},
		{"npm run build -- --mode production", []string{"run", "build", "--", "--mode", "production"}},
		{"npm run --workspace=web build", []string{"run", "--workspace=web", "build"}},
		{"bun install --frozen-lockfile --ignore-scripts", []string{"install", "--frozen-lockfile", "--ignore-scripts"}},
		{"pnpm install --offline --frozen-lockfile", []string{"install", "--offline", "--frozen-lockfile"}},
		{"yarn install --immutable", []string{"install", "--immutable"}},
	} {
		t.Run(tc.command, func(t *testing.T) {
			manager, _, _ := strings.Cut(tc.command, " ")
			base := installFakeCommand(t, manager, `printf '%s\n' "$@"`)
			output := captureRoutedOutput(t, tc.command, base)
			if want := strings.Join(tc.want, "\n") + "\n"; output != want {
				t.Errorf("forwarded arguments = %q, want %q", output, want)
			}
		})
	}
}

func TestNpmInstallIgnoreScripts(t *testing.T) {
	if _, err := exec.LookPath("npm"); err != nil {
		t.Skip("npm is needed for the offline install regression")
	}
	root := t.TempDir()
	for name, content := range map[string]string{
		"package.json": `{"name":"wake-ignore-scripts-check","version":"1.0.0","private":true,"scripts":{"postinstall":"node -e \"require('fs').writeFileSync('postinstall-ran','ran')\""}}`,
		"user.npmrc":   "",
		"global.npmrc": "",
	} {
		if err := os.WriteFile(filepath.Join(root, name), []byte(content), 0644); err != nil {
			t.Fatal(err)
		}
	}
	var env []string
	for _, value := range os.Environ() {
		if !strings.HasPrefix(strings.ToLower(value), "npm_config_") {
			env = append(env, value)
		}
	}
	env = append(env,
		"npm_config_offline=true", "npm_config_audit=false", "npm_config_fund=false",
		"npm_config_update_notifier=false", "npm_config_ignore_scripts=false",
		"npm_config_cache="+filepath.Join(root, "cache"),
		"npm_config_userconfig="+filepath.Join(root, "user.npmrc"),
		"npm_config_globalconfig="+filepath.Join(root, "global.npmrc"),
	)
	cmd := Route("npm install --ignore-scripts --offline --no-audit --no-fund", RouteOptions{Dir: root, Env: env})
	var output bytes.Buffer
	cmd.(OutputSetter).SetOutput(&output, &output)
	if err := cmd.Run(); err != nil {
		t.Fatalf("offline install failed: %v; output: %s", err, output.String())
	}
	if _, err := os.Stat(filepath.Join(root, "postinstall-ran")); !errors.Is(err, os.ErrNotExist) {
		t.Errorf("ignored postinstall hook ran: %v; output: %s", err, output.String())
	}
}
