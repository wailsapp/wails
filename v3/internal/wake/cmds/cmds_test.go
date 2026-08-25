package cmds

import (
	"os"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestRouteGoBuild(t *testing.T) {
	tests := []struct {
		name     string
		cmd      string
		wantType string
	}{
		{"simple", "go build", "*cmds.GoBuildCmd"},
		{"with output", "go build -o bin/app", "*cmds.GoBuildCmd"},
		{"with tags", "go build -tags production", "*cmds.GoBuildCmd"},
		{"with ldflags", `go build -ldflags "-s -w"`, "*cmds.GoBuildCmd"},
		{"with trimpath", "go build -trimpath", "*cmds.GoBuildCmd"},
		{"with package", "go build ./cmd/app", "*cmds.GoBuildCmd"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ex := Route(tt.cmd, RouteOptions{})
			require.IsType(t, &GoBuildCmd{}, ex)
		})
	}
}

func TestRouteGoRun(t *testing.T) {
	ex := Route("go run ./cmd/app", RouteOptions{})
	require.IsType(t, &GoRunCmd{}, ex)
}

func TestRouteGoTest(t *testing.T) {
	ex := Route("go test ./...", RouteOptions{})
	require.IsType(t, &GoTestCmd{}, ex)
}

func TestParseGoTestCount(t *testing.T) {
	// Regression: parseGoTestArgs used to read only the first digit of the
	// -count value via args[i+1][0]-'0', so `-count 10` became 1 and
	// `-count 100` became 1 too. Pin the full-integer parse in place.
	for _, in := range []string{"1", "10", "100", "1000"} {
		got := parseGoTestArgs([]string{"go", "test", "-count", in, "./..."})
		want := map[string]int{"1": 1, "10": 10, "100": 100, "1000": 1000}[in]
		require.Equal(t, want, got.Count, "-count %s parsed as %d", in, got.Count)
	}
}

func TestRouteGoModTidy(t *testing.T) {
	ex := Route("go mod tidy", RouteOptions{})
	require.IsType(t, &GoModTidyCmd{}, ex)
}

func TestRouteGoMod(t *testing.T) {
	ex := Route("go mod download", RouteOptions{})
	require.IsType(t, &GoModCmd{}, ex)
}

func TestRouteGoInstall(t *testing.T) {
	ex := Route("go install", RouteOptions{})
	require.IsType(t, &GoInstallCmd{}, ex)
}

func TestRouteGoVet(t *testing.T) {
	ex := Route("go vet", RouteOptions{})
	require.IsType(t, &GoVetCmd{}, ex)
}

func TestRouteGoFmt(t *testing.T) {
	ex := Route("go fmt", RouteOptions{})
	require.IsType(t, &GoFmtCmd{}, ex)
}

func TestRouteNpmInstall(t *testing.T) {
	ex := Route("npm install", RouteOptions{})
	require.IsType(t, &NpmInstallCmd{}, ex)
}

func TestRouteNpmRun(t *testing.T) {
	ex := Route("npm run build", RouteOptions{})
	require.IsType(t, &NpmRunCmd{}, ex)
}

func TestRouteBunInstall(t *testing.T) {
	ex := Route("bun install", RouteOptions{})
	require.IsType(t, &BunInstallCmd{}, ex)
}

func TestRoutePnpmInstall(t *testing.T) {
	ex := Route("pnpm install", RouteOptions{})
	require.IsType(t, &PnpmInstallCmd{}, ex)
}

func TestRouteYarnInstall(t *testing.T) {
	ex := Route("yarn install", RouteOptions{})
	require.IsType(t, &YarnInstallCmd{}, ex)
}

func TestRouteShellFallback(t *testing.T) {
	ex := Route("echo hello world", RouteOptions{})
	require.IsType(t, &ShellCmd{}, ex)
}

func TestParseGoBuildArgs(t *testing.T) {
	tests := []struct {
		name     string
		args     []string
		expected GoBuildOptions
	}{
		{
			name:     "simple build",
			args:     []string{"go", "build"},
			expected: GoBuildOptions{},
		},
		{
			name:     "with output",
			args:     []string{"go", "build", "-o", "bin/app"},
			expected: GoBuildOptions{Output: "bin/app"},
		},
		{
			name:     "with tags",
			args:     []string{"go", "build", "-tags", "production"},
			expected: GoBuildOptions{Tags: []string{"production"}},
		},
		{
			name:     "with multiple tags",
			args:     []string{"go", "build", "-tags", "production server"},
			expected: GoBuildOptions{Tags: []string{"production", "server"}},
		},
		{
			name:     "with ldflags",
			args:     []string{"go", "build", "-ldflags", "-s -w"},
			expected: GoBuildOptions{Ldflags: "-s -w"},
		},
		{
			name:     "with trimpath",
			args:     []string{"go", "build", "-trimpath"},
			expected: GoBuildOptions{Trimpath: true},
		},
		{
			name:     "with package",
			args:     []string{"go", "build", "./cmd/app"},
			expected: GoBuildOptions{Package: "./cmd/app"},
		},
		{
			name: "complex build",
			args: []string{"go", "build", "-tags", "production", "-ldflags", "-s -w", "-trimpath", "-o", "bin/app", "./cmd/app"},
			expected: GoBuildOptions{
				Tags:     []string{"production"},
				Ldflags:  "-s -w",
				Trimpath: true,
				Output:   "bin/app",
				Package:  "./cmd/app",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			opts := parseGoBuildArgs(tt.args)
			require.Equal(t, tt.expected.Output, opts.Output)
			require.Equal(t, tt.expected.Tags, opts.Tags)
			require.Equal(t, tt.expected.Ldflags, opts.Ldflags)
			require.Equal(t, tt.expected.Trimpath, opts.Trimpath)
			require.Equal(t, tt.expected.Package, opts.Package)
		})
	}
}

func TestParseGoTestArgs(t *testing.T) {
	tests := []struct {
		name     string
		args     []string
		expected GoTestOptions
	}{
		{
			name:     "simple test",
			args:     []string{"go", "test"},
			expected: GoTestOptions{},
		},
		{
			name:     "verbose test",
			args:     []string{"go", "test", "-v"},
			expected: GoTestOptions{Verbose: true},
		},
		{
			name:     "test with package",
			args:     []string{"go", "test", "./..."},
			expected: GoTestOptions{Package: "./..."},
		},
		{
			name:     "test with run",
			args:     []string{"go", "test", "-run", "TestFoo"},
			expected: GoTestOptions{Run: "TestFoo"},
		},
		{
			name:     "test with race",
			args:     []string{"go", "test", "-race"},
			expected: GoTestOptions{Race: true},
		},
		{
			name:     "test with cover",
			args:     []string{"go", "test", "-cover"},
			expected: GoTestOptions{Cover: true},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			opts := parseGoTestArgs(tt.args)
			require.Equal(t, tt.expected.Verbose, opts.Verbose)
			require.Equal(t, tt.expected.Package, opts.Package)
			require.Equal(t, tt.expected.Run, opts.Run)
			require.Equal(t, tt.expected.Race, opts.Race)
			require.Equal(t, tt.expected.Cover, opts.Cover)
		})
	}
}

func TestRouteSetsDirAndEnv(t *testing.T) {
	opts := RouteOptions{
		Dir: "/tmp/test",
		Env: []string{"FOO=bar"},
	}

	t.Run("go build sets dir", func(t *testing.T) {
		ex := Route("go build", opts)
		cmd := ex.(*GoBuildCmd)
		require.Equal(t, "/tmp/test", cmd.Dir)
		require.Equal(t, []string{"FOO=bar"}, cmd.Env)
	})

	t.Run("shell fallback sets dir", func(t *testing.T) {
		ex := Route("echo hello", opts)
		cmd := ex.(*ShellCmd)
		require.Equal(t, "/tmp/test", cmd.Dir)
		require.Equal(t, []string{"FOO=bar"}, cmd.Env)
	})

	t.Run("npm install sets dir", func(t *testing.T) {
		ex := Route("npm install", opts)
		cmd := ex.(*NpmInstallCmd)
		require.Equal(t, "/tmp/test", cmd.Dir)
		require.Equal(t, []string{"FOO=bar"}, cmd.Env)
	})
}

func TestDetectPackageManager(t *testing.T) {
	dir := t.TempDir()

	t.Run("npm", func(t *testing.T) {
		createFile(t, dir, "package-lock.json")
		require.Equal(t, "npm", DetectPackageManager(dir))
	})

	t.Run("bun", func(t *testing.T) {
		dir := t.TempDir()
		createFile(t, dir, "bun.lock")
		require.Equal(t, "bun", DetectPackageManager(dir))
	})

	t.Run("pnpm", func(t *testing.T) {
		dir := t.TempDir()
		createFile(t, dir, "pnpm-lock.yaml")
		require.Equal(t, "pnpm", DetectPackageManager(dir))
	})

	t.Run("yarn", func(t *testing.T) {
		dir := t.TempDir()
		createFile(t, dir, "yarn.lock")
		require.Equal(t, "yarn", DetectPackageManager(dir))
	})

	t.Run("default npm", func(t *testing.T) {
		require.Equal(t, "npm", DetectPackageManager(t.TempDir()))
	})
}

func createFile(t *testing.T, dir, name string) {
	f, err := os.Create(dir + "/" + name)
	require.NoError(t, err)
	f.Close()
}
