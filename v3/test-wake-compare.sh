#!/bin/bash
set -euo pipefail

SCRIPT_DIR="$(cd "$(dirname "$0")" && pwd)"
TEST_DIR=$(mktemp -d)
PROJECT_NAME="wakecompare"

echo "=== Wake Binary Comparison Test ==="
echo "Test directory: $TEST_DIR"

echo ""
echo "[1/5] Building wails3 binary..."
cd "$SCRIPT_DIR"
go build -o "$TEST_DIR/wails3" ./cmd/wails3
echo "  wails3 built"

WA="$TEST_DIR/wails3"

create_project() {
    local dir=$1
    mkdir -p "$dir/bin" "$dir/build/darwin" "$dir/build/windows" "$dir/build/linux" "$dir/build/ios" "$dir/build/android" "$dir/frontend/dist"

    # go.mod
    echo "module wakecompare" > "$dir/go.mod"
    echo "" >> "$dir/go.mod"
    echo "go 1.23" >> "$dir/go.mod"
    echo "" >> "$dir/go.mod"
    echo "require github.com/wailsapp/wails/v3 v3.0.0-00010101000000-000000000000" >> "$dir/go.mod"
    echo "" >> "$dir/go.mod"
    echo "replace github.com/wailsapp/wails/v3 => $SCRIPT_DIR" >> "$dir/go.mod"

    # main.go
    cat > "$dir/main.go" << 'EOF'
package main

import (
	"embed"
	"github.com/wailsapp/wails/v3/pkg/application"
)

//go:embed all:frontend/dist
var assets embed.FS

func main() {
	app := application.New(application.Options{
		Name:        "WakeCompare",
		Description: "Test app",
		Assets: application.AssetOptions{
			Handler: application.AssetFileServerFS(assets),
		},
	})
	app.Window.NewWithOptions(application.WebviewWindowOptions{
		Title: "WakeCompare",
		URL:   "/",
	})
	app.Run()
}
EOF

    # frontend/dist
    echo '<html><body>test</body></html>' > "$dir/frontend/dist/index.html"

    # build/config.yml
    cat > "$dir/build/config.yml" << 'EOF'
dev_mode:
  root_path: .
  main: "main.go"
  pre_build:
    - command: 'echo pre-build'
    - command: 'go mod tidy'
  build:
    - command: 'wails3 build'
  watch:
    - "go"
    - "frontend"
  ignore:
    directory:
      - frontend/node_modules
      - bin
EOF

    # Taskfile.yml
    cat > "$dir/Taskfile.yml" << 'EOF'
version: '3'
vars:
  APP_NAME: "wakecompare"
  BIN_DIR: "bin"
includes:
  common: ./build/Taskfile.yml
  darwin: ./build/darwin/Taskfile.yml
tasks:
  build:
    cmds:
      - task: darwin:build
EOF

    # build/Taskfile.yml
    cat > "$dir/build/Taskfile.yml" << 'EOF'
version: '3'
tasks:
  go:mod:tidy:
    cmds:
      - go mod tidy
EOF

    # build/darwin/Taskfile.yml
    cat > "$dir/build/darwin/Taskfile.yml" << 'EOF'
version: '3'
includes:
  common: ../Taskfile.yml
tasks:
  build:
    deps:
      - task: common:go:mod:tidy
    cmds:
      - go build -o bin/wakecompare
    env:
      GOOS: darwin
      CGO_ENABLED: 1
EOF

    cd "$dir"
    go mod tidy 2>&1 | head -5 || true
}

echo ""
echo "[2/5] Creating project A (task CLI)..."
create_project "$TEST_DIR/projectA"
PROJECT_A="$TEST_DIR/projectA"

echo ""
echo "[3/5] Creating project B (WAILS_USE_WAKE)..."
create_project "$TEST_DIR/projectB"
PROJECT_B="$TEST_DIR/projectB"

echo ""
echo "[4/5] Building project A (task CLI)..."
cd "$PROJECT_A"
unset WAILS_USE_WAKE
echo "  Running: $WA build"
$WA build 2>&1 || true

echo ""
echo "[5/5] Building project B (WAILS_USE_WAKE=true)..."
cd "$PROJECT_B"
export WAILS_USE_WAKE=true
echo "  Running: WAILS_USE_WAKE=true $WA build"
$WA build 2>&1 || true

echo ""
echo "=== Results ==="

BIN_A="$PROJECT_A/bin/$PROJECT_NAME"
BIN_B="$PROJECT_B/bin/$PROJECT_NAME"

if [ -f "$BIN_A" ]; then
    SIZE_A=$(stat -f%z "$BIN_A")
    echo "Binary A (task):  $BIN_A ($SIZE_A bytes)"
    file "$BIN_A"
else
    echo "Binary A NOT FOUND at $BIN_A"
    ls -la "$PROJECT_A/bin/" 2>/dev/null || echo "  (no bin dir)"
fi

if [ -f "$BIN_B" ]; then
    SIZE_B=$(stat -f%z "$BIN_B")
    echo "Binary B (wake):  $BIN_B ($SIZE_B bytes)"
    file "$BIN_B"
else
    echo "Binary B NOT FOUND at $BIN_B"
    ls -la "$PROJECT_B/bin/" 2>/dev/null || echo "  (no bin dir)"
fi

if [ -f "$BIN_A" ] && [ -f "$BIN_B" ]; then
    echo ""
    if [ "$SIZE_A" = "$SIZE_B" ]; then
        echo "SIZE MATCH: $SIZE_A bytes"
    else
        echo "SIZE DIFF: A=$SIZE_A B=$SIZE_B (delta: $(($SIZE_B - $SIZE_A)))"
        echo "  (Expected if build timestamps differ)"
    fi

    echo ""
    echo "Symbol table comparison:"
    SYM_DIFF=$(diff <(nm "$BIN_A" | sort) <(nm "$BIN_B" | sort) | grep -c "^[<>]" || true)
    if [ "$SYM_DIFF" -eq 0 ]; then
        echo "Symbol tables match"
    else
        echo "$SYM_DIFF symbol differences"
        diff <(nm "$BIN_A" | sort) <(nm "$BIN_B" | sort) | grep "^[<>]" | head -5
    fi
fi

echo ""
echo "Artifacts at: $TEST_DIR"
echo "Done."
