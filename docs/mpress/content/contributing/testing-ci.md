---
title: "Testing \u0026 Continuous Integration"
description: "How Wails v3 ensures quality with unit tests, integration suites, race detection, and GitHub Actions CI."
slug: "contributing/testing-ci"
sourcePath: "contributing/testing-ci.md"
---

Robust desktop frameworks demand rock-solid testing. Wails v3 employs a **layered strategy**:

| Layer | Goal | Tooling |
| --- | --- | --- |
| Unit tests | Fast feedback for isolated functions | `go test ./...` |
| Generator / CLI tests | Validate `wails3 generate bindings` and the CLI plumbing | `task test:generator`, `task test:cli` |
| Template tests | Make sure every shipped template still builds | `task test:templates` |
| Race detection | Catch data races in the runtime & bridge | `go test -race ./...` |
| CI matrix | Cross-OS confidence on every PR | GitHub Actions |

This document explains **where tests live**, **how to run them**, and **what the Taskfile orchestrates**.

`v3/TESTING.md` documents cross-platform example builds; it is not a race-suppression catalogue.

---

## 1. Directory Conventions

```
v3/
├── internal/.../_test.go     # Unit tests for internal packages
├── pkg/.../_test.go          # Public API tests
├── tasks/events/generate.go  # Code generator for event constants (NOT a test harness)
├── tests/                    # Top-level integration test harness
└── TESTING.md                # Cross-platform example-build guidance
```

Guidelines:

- **Keep unit tests next to code** (`foo.go` ↔ `foo_test.go`).
- Use the **black-box style** for `pkg/` packages (`package application_test`) where it improves API hygiene.
- Shared fixtures live where they're used (there is no central `internal/testutil/` package in this tree — wire helpers per package instead).

---

## 2. Unit Tests

### Writing Tests

```go
func TestEventConstants(t *testing.T) {
    assert.NotEmpty(t, events.Common.WindowFocus)
}
```

Recommendations:

- Use [`stretchr/testify`](https://github.com/stretchr/testify) — already in `go.mod`.
- Prefer **table-driven** tests whenever several inputs, edge cases, or expected outcomes exercise the same behaviour. Give every case a descriptive name.
- Stub platform quirks behind build tags (`foo_windows_test.go`, `foo_darwin_test.go`, …) when needed.

### Coverage expectation

New and changed logic is expected to have 100% Go statement coverage. Measure the package you changed rather than relying on a repository-wide percentage:

```bash
go test ./pkg/application -coverprofile=coverage.out
go tool cover -func=coverage.out
```

Some paths are reasonably untestable in a normal test environment—for example, platform-only failures, hardware-dependent behaviour, or a defensive fallback that cannot be induced safely. Keep those exceptions narrow and explain each uncovered path in the PR description.

### Running Locally

```bash
cd v3
go test ./... -cover
```

Or via Taskfile (real targets — there is no `task test` shortcut):

```
task test:cli            # CLI plumbing tests
task test:generator      # bindings generator round-trip tests
task test:templates      # build every shipped template
task test:infrastructure # supporting helpers
task test:examples       # exercise the example matrix (downloads as needed)
task test:all            # everything above
task sanity              # default GTK4 smoke check (legacy: sanity:gtk3)
task precommit           # what you should run before pushing
```

---

## 3. Integration Tests

`v3/tests/` hosts the cross-package integration harness. The Taskfile `test:example:*` and `test:examples:*` targets drive compilation checks across darwin / windows / linux (including Docker-based GTK3 / GTK4 matrices on Linux).

> Runnable examples live in `v3/examples/`. The test targets select and build
>
> the examples appropriate to the host platform or CI matrix.

Run the host-platform smoke suite with:

```
task test:examples       # host
task test:examples:all   # full matrix (slow)
```

---

## 4. Race Detection

Data races are fatal in GUI runtimes.

### Investigating a race

Inspect both conflicting accesses and reproduce on the affected platform. Do not treat a report as benign or suppress it without evidence. `v3/TESTING.md` describes example compilation; passing those builds does not establish race freedom or native runtime behaviour.

### Local Race Suite

```
go test -race ./...
```

> `wails3 dev` has no `-race` flag — its CLI flags are `--config`, `--port`,
>
> and `-s` (enable HTTPS). To exercise the runtime under the race detector,
>
> build a test app with `go build -race` and run it directly.

---

## 5. GitHub Actions Workflows

Real workflow files under `.github/workflows/` (verified against the tree):

| File | Purpose |
| --- | --- |
| `build-and-test-v3.yml` | Go 1.26/1.27 matrix, runtime and generator checks, template/example compilation, and application build checks. Linux uses GTK4/WebKitGTK 6.0 by default; the legacy variant sets `BUILD_TAGS=gtk3`. Native launch and interaction must be checked separately. |
| `cross-compile-test-v3.yml` | Cross-compile sanity checks |
| `auto-changelog-v3.yml`, `changelog-v3.yml` | Changelog automation |
| `nightly-release-v3.yml` | Nightly v3 release artefacts |
| `bump-webview2-v3.yml`, `release-webview2.yml` | WebView2 dependency / release management |
| `build-cross-image.yml` | Builds the cross-compiler container image |
| `publish-npm.yml` | Publishes the embedded `@wailsio/runtime` JS runtime to npm |
| `pr-master.yml` | PR-side checks against the `master` branch |
| `semgrep.yml` | Semgrep static analysis |
| `stale-issues.yml`, `generate-sponsor-image.yml`, `build-and-test.yml`, `weekly-release-v2.yml` | Repository housekeeping and v2 workflows; the old Crowdin sync/upload workflows are retired. |

The CI steps mirror the Taskfile targets above (`task test:cli`, `task test:generator`, `task test:templates`, `task test:examples`, …) so you can reproduce CI locally one-to-one. The smoke `wails3 build` step in `build-and-test-v3.yml` is invoked with **no extra flags** — there is no `-skip-package` flag on `wails3 build`.

---

## 6. Local CI Parity

There is no single `task ci` umbrella target. Reproduce CI by chaining the real targets:

```
task precommit
task test:cli
task test:generator
task test:templates
task test:examples
```

---

## 7. Troubleshooting Failing Tests

| Symptom | Likely Cause | Fix |
| --- | --- | --- |
| **Race in `webview_window_darwin.go`** | Mutating window state off the main thread | Marshal through `application.InvokeAsync` / `Invoke` so the call runs on the main thread |
| **Linux test hangs in headless CI** | GTK needs a display | Run under `xvfb-run`, e.g. `xvfb-run task test:examples:linux` |
| **Template build fails** | Frontend lockfile out of date | Re-run `wails3 init` against a clean directory to refresh the template |
| **Coverpkg errors** | Integration test importing `main` | Switch to build tag `//go:build integration` and gate the import |

---

## 8. Adding New Tests

1. **Unit** — create `*_test.go`, run `go test ./...`
2. **Generator / CLI** — extend the cases under `internal/generator/testcases/` or `internal/commands/*_test.go` and re-run `task test:generator` / `task test:cli`
3. **Templates / examples** — make sure the shipped templates still build with `task test:templates`

---

## 9. Key File Map

| What | Path |
| --- | --- |
| Generator round-trip test | `internal/generator/generate_test.go` |
| Build-assets test | `internal/commands/build-assets_test.go` |
| Cross-platform example builds | `v3/TESTING.md` |
| Taskfile test targets | `v3/Taskfile.yaml` |
| Event-constants generator | `v3/tasks/events/generate.go` |
| CI workflow | `.github/workflows/build-and-test-v3.yml` (Go 1.26/1.27 via `actions/setup-go@v5`) |
| Static analysis | `.github/workflows/semgrep.yml` |
| Runtime npm publish | `.github/workflows/publish-npm.yml` |

---

Quality isn't an afterthought in Wails v3. With unit tests, generator/template suites, race detection, and a cross-platform CI matrix you can contribute confidently, knowing your changes run green on every OS we support. Happy testing!
