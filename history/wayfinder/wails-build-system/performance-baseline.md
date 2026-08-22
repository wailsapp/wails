# Wails manifest build-system performance baseline

Initial baseline measured 2026-08-16 on Linux/amd64 at `4db35ce43270` plus the
then-uncommitted manifest build-system implementation. The latest controlled
audit is dated 2026-08-22 below.

## Environment

- Bazzite/Fedora Linux, kernel 7.1.3
- AMD Ryzen 9 3950X, 16 cores / 32 hardware threads
- Go 1.26.5
- Node 22.22.2 and npm from the Codex desktop runtime
- Task 3.40.1-patched3
- Go build, Go module, npm, and downloaded-tool caches were warm unless a row
  explicitly says otherwise
- Wails action/Artifact caches lived in disposable project fixtures

The badge fixture exercises two services after adding a tiny benchmark service.
The dock example is the larger real-project fixture. Both remain inside the
Wails module so their imports and package graph are representative of the
examples they were copied from.

## Complete-command measurements

These historical baseline scenarios used one or two unmeasured warmups followed
by five measured CLI processes. The permanent release gate now uses two
warmups and seven measured samples as defined in `performance-acceptance.md`.
`Wall` is complete process latency. `Graph` is the duration reported
by Wails and approximates graph critical-path work. `CPU` is aggregate user and
system process time. Ranges are min–max across the five samples.

| Scenario | Median wall (range) | Median CPU | Median graph | Executed work |
| --- | ---: | ---: | ---: | --- |
| Manifest badge no-op build | 87.6ms (77.9–89.8) | 98.5ms | 9ms | 0 ran, 4 cached |
| Manifest dock no-op build | 74.2ms (72.5–83.1) | 86.0ms | 5ms | 0 ran, 4 cached |
| Taskfile Wake badge no-op | 91.2ms (80.4–98.8) | 103.9ms | 14ms | 0 ran, 6 cached |
| External Task badge no-op | 228.5ms (225.0–237.0) | 1,501.8ms | n/a | Task still invokes `go mod tidy` and `go build` |
| Manifest DEB no-op package | 92.0ms (81.7–96.2) | 104.1ms | 16ms | 0 ran, 6 cached |
| Manifest AppImage no-op package | 84.7ms (77.9–92.0) | 96.8ms | 6ms | 0 ran, 6 cached |
| Ordinary Go method-body edit | 904.6ms (869.9–921.0) | 1,407.4ms | 819ms | compile only; 1 ran, 3 cached |
| Binding-shape edit | 2,816.2ms (1,896.6–2,861.2) | 8,373.2ms | 2,700ms | typically bindings, frontend, compile; 3 ran, 1 cached |
| Frontend comment, identical bundle | 965.1ms (931.4–989.2) | 2,054.6ms | 884ms | frontend only; 1 ran, 3 cached |
| Frontend edit, changed bundle | 1,787.9ms (1,744.7–1,806.6) | 3,296.1ms | 1,700ms | frontend and compile; 2 ran, 2 cached |
| Missing binary Artifact restore | 85.3ms (79.6–89.6) | 99.4ms | 7ms | 0 ran, 4 cached; restored bytes matched |

Distinct binding shapes used `string`, `number`, `boolean`, array, and map
outputs. Two samples executed only bindings and compile because different Go
numeric types intentionally collapse to the same TypeScript `number` Artifact;
the unchanged frontend bundle was correctly restored instead of rebuilt.

The manifest no-op path is about 62% lower wall time than external Task on this
host and is within measurement noise of the legacy Taskfile-Wake experiment.
The important difference is architectural: the manifest path has four typed
semantic Nodes rather than interpreting six Taskfile tasks.

## Cold and lifecycle observations

These are deliberately not five-sample warm distributions:

| Scenario | Wall | CPU | Graph | Work |
| --- | ---: | ---: | ---: | --- |
| Dock action-cache cold build, tool caches warm | 3,184ms | 9,464ms | 3,100ms | 4 ran |
| First AppImage package in the new target workspace | 21,173ms | 52,718ms | 21,100ms | 3 ran, 3 cached; includes linuxdeploy/network/tool work |
| Live Dev method-body rebuild and backend replacement | 870ms | not captured | 870ms | compile only; 1 ran, 3 cached |

The Dev run kept Vite alive, rebuilt only Go compilation, started the
replacement backend after readiness stabilization, reported `Backend rebuilt
and restarted`, and left neither child processes nor port 9245 behind after
Ctrl-C.

## Planner and executor overhead

Five Go benchmark runs on the same host:

| Benchmark | Median |
| --- | ---: |
| Plan shared Linux/Windows/macOS graph | 687µs/op |
| Warm single-target no-op Executor | 2.26ms/op |
| Warm three-target no-op Executor | 3.28ms/op |

The Executor benchmarks use deterministic in-process handlers to isolate Wake
orchestration and cache lookup from native toolchains. Actual cross-platform
compilation and packaging remain host-gated acceptance work.

## Reproduction and gate

Build the helper and invoke it from a warmed disposable badge project:

```bash
go build -o /tmp/wails-build-benchmark ./scripts/benchmark-manifest-build.go

/tmp/wails-build-benchmark \
  -name badge-noop-linux-amd64 \
  -samples 7 \
  -warmups 2 \
  -baseline /path/to/wails/v3/internal/wake/benchmark/testdata/badge-noop-linux-amd64.json \
  -max-ms 100 \
  -max-regression 20 \
  -max-mad-percent 15 \
  -expect-ran 0 \
  -expect-cached 6 \
  -- /path/to/wails3 build --target linux/amd64
```

The checked-in baseline now records the audited six-Node graph at 99.092ms
median wall time. It passes the absolute 100ms budget; subsequent controlled
runs must also stay within the 20% relative budget. Raw checked-in samples remain in
`internal/wake/benchmark/testdata/badge-noop-linux-amd64.json`.

Run the in-process graph measurements with:

```bash
go test ./internal/wake/pipeline -run '^$' \
  -bench 'Benchmark(PlanMultiTarget|WarmNoopExecutor|WarmNoopMultiTargetExecutor)$' \
  -benchmem -count=5
```

Wall-clock numbers are host observations, not portable constants. The
permanent acceptance decision should gate a controlled runner by medians while
keeping semantic executed/cached-node assertions in ordinary tests.

## 2026-08-17 acceptance run

The completed acceptance harness ran against the rebuilt reference CLI:

- badge no-op: 84.049ms median wall, 5.55% MAD, 0 ran and 4 cached in all
  seven samples; the 100ms absolute and 20% checked-baseline gates passed;
- cold dock action cache: 2,787ms wall versus 2,700ms reported graph time,
  4 ran, and 3.22% orchestration overhead; the 5% gate passed;
- native Linux build and concurrent DEB/RPM/Arch packaging passed, followed by
  zero-work reruns at 91.0ms and 83.3ms respectively;
- native AppImage packaging completed in 19.9s and the complete four-format
  package Plan then reran at 89.9ms with 0 ran and 9 cached.

The package acceptance output showed RPM, Arch, and DEB adapter execution
overlapping after their shared asset dependency. Deterministic tests separately
prove package and multi-Target overlap when declared resource claims fit.

After the final concurrency-isolation review, a fresh badge fixture measured
93.503ms median wall with 8.19% MAD and 0 ran/4 cached in all seven samples. It
passed both latency budgets despite one 191ms scheduling outlier, which the
median/MAD policy correctly treated as diagnostic rather than decisive.

## 2026-08-22 full local audit

The production graph now has six Nodes for the default Linux binary path:
install, bindings, frontend, compile, publish, and final receipt collection.
The early four-Node measurements above remain historical evidence and must not
be used as current work-count expectations.

Two warmups and seven measured processes on the same host produced:

- 99.092ms median wall time, 5.05% MAD, 26ms median graph time, 204.036ms
  aggregate CPU, and 80.3MiB median peak RSS;
- 0 ran and 6 cached in every sample, with a stable 9,741,248-byte binary and
  stable SHA-256 digest;
- a modified published binary regenerated byte-identically in 65ms with only
  publication executing; and
- a missing published binary restored byte-identically in 49ms with all six
  Nodes reported cached.

The result passes the 100ms absolute median and 15% MAD gates. It is 20.1ms
faster than the pre-optimization 119.239ms audit run. The optimization caches a
final receipt only when every final output is reproducible, artifact-caches the
published output itself, and uses stable file identity metadata to avoid
rehashing unchanged regular-file artifacts on Linux and macOS. Signing and its
downstream publication/receipt path remain deliberately non-reusable.
