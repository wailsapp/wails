#!/bin/sh
set -eu

script_dir=$(CDPATH= cd -- "$(dirname -- "$0")" && pwd)
v3_dir=$(CDPATH= cd -- "$script_dir/../../../" && pwd)
output_dir=${1:-"$(mktemp -d "${TMPDIR:-/tmp}/native-editor-benchmark.XXXXXX")"}
launch_iterations=${LAUNCH_ITERATIONS:-20}
launch_document_sizes=${LAUNCH_DOCUMENT_SIZES:-"0 1048576 10485760"}
idle_seconds=${IDLE_SECONDS:-30}

mkdir -p "$output_dir"
wails_binary="$output_dir/native-notes-wails"
swift_binary="$output_dir/native-notes-swift"
sampler="$output_dir/rusage-sample"

(
    cd "$v3_dir"
    GOWORK=off go build -tags 'production wails_native' -trimpath -ldflags='-s -w' \
        -o "$wails_binary" ./examples/mac-native-editor
    swiftc -O -whole-module-optimization examples/mac-native-editor-swift/main.swift \
        -framework AppKit -o "$swift_binary"
)
cc -O2 "$script_dir/rusage_sample.c" -o "$sampler"

{
    printf 'implementation,binary_bytes\n'
    printf 'wails,%s\n' "$(stat -f %z "$wails_binary")"
    printf 'swift,%s\n' "$(stat -f %z "$swift_binary")"
} > "$output_dir/binaries.csv"
otool -L "$wails_binary" > "$output_dir/wails-linkage.txt"
otool -L "$swift_binary" > "$output_dir/swift-linkage.txt"

printf 'implementation,document_bytes,iteration,real_seconds,user_seconds,system_seconds,max_rss_bytes,peak_footprint_bytes\n' > "$output_dir/launches.csv"
measure_launches() {
    implementation=$1
    binary=$2
    document_bytes=$3
    iteration=1
    while [ "$iteration" -le "$launch_iterations" ]; do
        stats="$output_dir/time-$implementation-$iteration.txt"
        NATIVE_EDITOR_BENCH_AUTO_QUIT_MS=1000 \
            NATIVE_EDITOR_BENCH_DOCUMENT_BYTES="$document_bytes" \
            /usr/bin/time -l -o "$stats" "$binary" >/dev/null 2>&1
        awk -v implementation="$implementation" -v document_bytes="$document_bytes" -v iteration="$iteration" '
            / real .* user .* sys$/ { real=$1; user=$3; systime=$5 }
            /maximum resident set size$/ { rss=$1 }
            /peak memory footprint$/ { footprint=$1 }
            END { printf "%s,%s,%d,%s,%s,%s,%s,%s\n", implementation, document_bytes, iteration, real, user, systime, rss, footprint }
        ' "$stats" >> "$output_dir/launches.csv"
        rm "$stats"
        iteration=$((iteration + 1))
    done
}
for document_bytes in $launch_document_sizes; do
    measure_launches wails "$wails_binary" "$document_bytes"
    measure_launches swift "$swift_binary" "$document_bytes"
done

printf 'implementation,visibility,document_bytes,threads,ports,rusage\n' > "$output_dir/steady-state.csv"
measure_steady_state() {
    implementation=$1
    binary=$2
    visibility=$3
    document_bytes=$4
    hidden=0
    if [ "$visibility" = hidden ]; then hidden=1; fi
    ready="$output_dir/ready-$implementation-$visibility-$document_bytes"
    rm -f "$ready"
    NATIVE_EDITOR_BENCH_HIDDEN="$hidden" \
        NATIVE_EDITOR_BENCH_DOCUMENT_BYTES="$document_bytes" \
        NATIVE_EDITOR_BENCH_READY_FILE="$ready" \
        "$binary" >/dev/null 2>&1 &
    pid=$!
    attempts=0
    while [ ! -f "$ready" ] && kill -0 "$pid" 2>/dev/null; do
        sleep 0.05
        attempts=$((attempts + 1))
        if [ "$attempts" -gt 200 ]; then
            kill "$pid" 2>/dev/null || true
            wait "$pid" 2>/dev/null || true
            echo "timed out waiting for $implementation" >&2
            return 1
        fi
    done
    sleep 2
    process_stats=$(top -l 1 -pid "$pid" -stats threads,ports -ncols 2 | tail -1)
    threads=$(printf '%s\n' "$process_stats" | awk '{print $1}')
    ports=$(printf '%s\n' "$process_stats" | awk '{print $2}')
    rusage=$($sampler "$pid" "$idle_seconds")
    escaped_rusage=$(printf '%s' "$rusage" | sed 's/"/""/g')
    printf '%s,%s,%s,%s,%s,"%s"\n' \
        "$implementation" "$visibility" "$document_bytes" "$threads" "$ports" "$escaped_rusage" \
        >> "$output_dir/steady-state.csv"
    kill "$pid" 2>/dev/null || true
    wait "$pid" 2>/dev/null || true
}

if [ "${SKIP_STEADY_STATE:-0}" != 1 ]; then
    for document_bytes in 0 1048576 10485760; do
        measure_steady_state wails "$wails_binary" visible "$document_bytes"
        measure_steady_state swift "$swift_binary" visible "$document_bytes"
    done
    measure_steady_state wails "$wails_binary" hidden 0
    measure_steady_state swift "$swift_binary" hidden 0
fi

printf '%s\n' "$output_dir"
