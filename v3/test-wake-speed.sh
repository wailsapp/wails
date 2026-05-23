#!/bin/bash
set -euo pipefail

SCRIPT_DIR="$(cd "$(dirname "$0")" && pwd)"
TEST_DIR=$(mktemp -d)
PROJECT_NAME="wakecompare"
RUNS=5

echo "=== Wake vs Task Build Speed Comparison ==="
echo "Test directory: $TEST_DIR"
echo "Runs per method: $RUNS"

echo ""
echo "[1/4] Building wails3 binary..."
cd "$SCRIPT_DIR"
go build -o "$TEST_DIR/wails3" ./cmd/wails3
echo "  wails3 built"

WA="$TEST_DIR/wails3"

create_project() {
    local dir=$1
    mkdir -p "$dir/bin" "$dir/build/darwin" "$dir/build/windows" "$dir/build/linux" "$dir/build/ios" "$dir/build/android" "$dir/frontend/dist"

    echo "module wakecompare" > "$dir/go.mod"
    echo "" >> "$dir/go.mod"
    echo "go 1.23" >> "$dir/go.mod"
    echo "" >> "$dir/go.mod"
    echo "require github.com/wailsapp/wails/v3 v3.0.0-00010101000000-000000000000" >> "$dir/go.mod"
    echo "" >> "$dir/go.mod"
    echo "replace github.com/wailsapp/wails/v3 => $SCRIPT_DIR" >> "$dir/go.mod"

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

    mkdir -p "$dir/frontend/dist"
    echo '<html><body>test</body></html>' > "$dir/frontend/dist/index.html"

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

    cat > "$dir/build/Taskfile.yml" << 'EOF'
version: '3'
tasks:
  go:mod:tidy:
    cmds:
      - go mod tidy
EOF

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
    go mod tidy 2>&1 >/dev/null || true
}

echo ""
echo "[2/4] Creating project A (task CLI)..."
create_project "$TEST_DIR/projectA"
PROJECT_A="$TEST_DIR/projectA"

echo ""
echo "[3/4] Creating project B (WAILS_USE_WAKE)..."
create_project "$TEST_DIR/projectB"
PROJECT_B="$TEST_DIR/projectB"

# Run benchmarks
echo ""
echo "[4/4] Running benchmarks ($RUNS runs each)..."
echo ""

# Warm up both
echo "  Warming up..."
cd "$PROJECT_A"
rm -rf bin
$WA build >/dev/null 2>&1 || true
cd "$PROJECT_B"
rm -rf bin
WAILS_USE_WAKE=true $WA build >/dev/null 2>&1 || true

# Benchmark task CLI
echo ""
echo "=== Task CLI ==="
declare -a TASK_TIMES
for i in $(seq 1 $RUNS); do
    cd "$PROJECT_A"
    rm -rf bin
    START_NS=$(date +%s%N)
    $WA build >/dev/null 2>&1 || true
    END_NS=$(date +%s%N)
    ELAPSED_MS=$(( (END_NS - START_NS) / 1000000 ))
    TASK_TIMES+=($ELAPSED_MS)
    echo "  Run $i: ${ELAPSED_MS}ms"
done

# Benchmark wake
echo ""
echo "=== WAKE ==="
declare -a WAKE_TIMES
for i in $(seq 1 $RUNS); do
    cd "$PROJECT_B"
    rm -rf bin
    START_NS=$(date +%s%N)
    WAILS_USE_WAKE=true $WA build >/dev/null 2>&1 || true
    END_NS=$(date +%s%N)
    ELAPSED_MS=$(( (END_NS - START_NS) / 1000000 ))
    WAKE_TIMES+=($ELAPSED_MS)
    echo "  Run $i: ${ELAPSED_MS}ms"
done

# Calculate averages
TASK_SUM=0
for t in "${TASK_TIMES[@]}"; do TASK_SUM=$((TASK_SUM + t)); done
TASK_AVG=$((TASK_SUM / RUNS))

WAKE_SUM=0
for t in "${WAKE_TIMES[@]}"; do WAKE_SUM=$((WAKE_SUM + t)); done
WAKE_AVG=$((WAKE_SUM / RUNS))

echo ""
echo "=== Results ==="
echo "Task CLI avg: ${TASK_AVG}ms"
echo "Wake avg:     ${WAKE_AVG}ms"

if [ "$TASK_AVG" -gt "$WAKE_AVG" ]; then
    DIFF=$((TASK_AVG - WAKE_AVG))
    PCT=$(( (DIFF * 100) / TASK_AVG ))
    echo "Wake is ${DIFF}ms faster (${PCT}% improvement)"
else
    DIFF=$((WAKE_AVG - TASK_AVG))
    PCT=$(( (DIFF * 100) / TASK_AVG ))
    echo "Task CLI is ${DIFF}ms faster (${PCT}% overhead for wake)"
fi

echo ""
echo "Artifacts at: $TEST_DIR"
