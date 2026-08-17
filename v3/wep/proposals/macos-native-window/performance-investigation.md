# Wails 3 performance investigation

## Purpose

This report records what the native macOS experiment and its Swift/AppKit
control teach us about ordinary Wails 3 applications. It separates four costs
that are easy to confuse:

1. executable size;
2. process and framework memory;
3. startup work and garbage collection;
4. steady-state CPU, threads, and wakeups.

The central correction is that WebKit being present in `otool -L` does **not**
mean WebKit is copied into the Wails executable. WebKit is a dynamically linked
system framework. The large Wails executable was primarily statically linked
Go code and metadata made reachable by the monolithic `application` package.
WebKit matters much more after launch, when it creates UI-process threads and
GPU, Networking, and WebContent XPC processes.

The results below were collected on 6 August 2026 from the
`feat/macos-window-chrome` worktree. They are engineering measurements, not a
cross-platform product benchmark.

## Executive findings

- A production, stripped normal Wails toolbar example is 8,790,130 bytes.
- The native editor built through the ordinary monolithic path is 8,657,586
  bytes and still links WebKit, despite never creating a WebView.
- The same editor built with the new `wails_native` boundary is 3,489,122
  bytes, a 59.7% reduction, and has no WebKit, HTTP, TLS, X.509, asset-server,
  JavaScript-transport, or built-in-updater dependency.
- The equivalent optimized Swift executable is 176,736 bytes. That comparison
  is not a like-for-like runtime comparison: macOS provides Swift and AppKit as
  dynamic system libraries, while Go places its runtime, stack maps, type data,
  and garbage collector in every executable.
- A stripped empty Go executable on the same toolchain is 1,145,394 bytes. A
  minimal Go executable that touches AppKit is 1,180,594 bytes. The 3.49 MB
  native Wails binary is therefore about 2.31 MB above the measured Go+AppKit
  floor, not 3.48 MB above a 176 KB Swift-equivalent floor.
- A normal-sized native editor did not trigger Go GC during startup or a
  20-second idle observation. Its sampled main thread was blocked in AppKit,
  and its Go runtime threads were parked.
- Loading 10 MiB of text triggered three very short Go collections in the first
  11 ms. The longest observed wall-clock collection was below 0.5 ms. The large
  document memory gap was caused by whole-value copies and ownership across
  Go, C, Foundation, and `NSTextView`, not by time spent collecting garbage.
- In one five-second idle sample, the normal WebView example had 35 threads in
  its main process and created three WebKit XPC processes. The native Wails
  editor had 14 sampled threads; Swift had 6. Every sampled process was mostly
  blocked rather than consuming CPU.
- Lowering `GOMAXPROCS` from 14 to 2 did not reduce the native editor's sampled
  thread count or physical footprint. It reduced one normal WebView sample
  from 35 to 27 main-process threads, but did not produce a reliable memory or
  CPU win. Wails should not silently override application concurrency.

The best normal-Wails opportunities are therefore package modularity, removal
of unconditional embedded/default assets, optional updater and HTTP boundaries,
lazy WebView creation, and fewer copies for large payloads. GC tuning is not a
high-priority general optimization based on this evidence.

## Test system

- Apple M4 Max, 36 GiB RAM
- macOS 26.4.1 (25E253), arm64
- Xcode 26.6
- Go 1.26.2
- Swift 6.3.3

All reported executable sizes use `-trimpath -ldflags='-s -w'`. Wails release
builds also use the `production` tag. None of the measured executables was
bundled, signed, compressed, or notarized.

## Reproduction commands

Normal Wails/WebKit:

```sh
GOWORK=off go build -tags production -trimpath -ldflags='-s -w' \
  -o mac-toolbar ./examples/mac-toolbar
```

Native application through the old monolithic package shape:

```sh
GOWORK=off go build -tags production -trimpath -ldflags='-s -w' \
  -o native-editor-monolithic ./examples/mac-native-editor
```

Native application with the compile-time boundary:

```sh
GOWORK=off go build -tags 'production wails_native' \
  -trimpath -ldflags='-s -w' \
  -o native-editor-native ./examples/mac-native-editor
```

Swift control:

```sh
swiftc -O -whole-module-optimization \
  examples/mac-native-editor-swift/main.swift \
  -framework AppKit -o native-editor-swift
```

The repeatable native/Swift launch and steady-state harness is
`examples/mac-native-editor/benchmark/run.sh`. It now builds Wails with both
`production` and `wails_native`, so it does not accidentally compare Swift
against a runtime-native but build-monolithic binary.

Useful inspection commands:

```sh
stat -f '%z' <binary>
otool -L <binary>
size -m <unstripped-binary>
go tool nm -size <unstripped-binary>
GOWORK=off go list -tags 'production wails_native' -deps \
  ./examples/mac-native-editor
GODEBUG=gctrace=1 <binary>
sample <pid> 5 10 -file process.sample.txt
```

Physical footprint was read with `proc_pid_rusage(RUSAGE_INFO_V4)` through
`examples/mac-native-editor/benchmark/rusage_sample.c`. RSS is shown only as a
diagnostic because it includes shared framework mappings.

## What is in each executable

### Stripped production size

| Executable | Bytes | MiB | Links WebKit |
|---|---:|---:|:---:|
| Normal Wails `mac-toolbar` | 8,790,130 | 8.38 | Yes |
| Native editor, monolithic build | 8,657,586 | 8.26 | Yes |
| Native editor, `wails_native` | 3,489,122 | 3.33 | No |
| Native editor, plus `wails_single_instance` | 3,656,178 | 3.49 | No |
| Swift/AppKit control | 176,736 | 0.17 | No |
| Empty Go control | 1,145,394 | 1.09 | No |
| Minimal Go+AppKit control | 1,180,594 | 1.13 | No |

The normal Wails and monolithic native sizes are close because both compile the
same broad `pkg/application` feature graph. Runtime `NativeOnly: true` avoids
starting frontend infrastructure, but a runtime branch cannot change which Go
packages and cgo translation units were compiled.

The production tag is itself significant. The same monolithic native editor
built without `production` was 10,467,602 bytes rather than 8,657,586 bytes, a
17.3% reduction. Wails' generated release tasks already use `production`,
`-trimpath`, and `-ldflags='-s -w'`; CI should keep asserting that contract.

### WebKit's actual contribution to file size

An intermediate diagnostic build excluded the Darwin WebView backend and its
WebKit framework linkage but left the frontend Go graph in place:

| Diagnostic build | Bytes |
|---|---:|
| Monolithic debug build | 10,467,586 |
| WebKit backend excluded only | 10,252,322 |
| Difference | 215,264 (2.1%) |

This is the cleanest proof that “the binary is large because it links WebKit”
was the wrong causal claim. The framework is supplied by macOS. The remaining
multi-megabyte reduction appeared only after the Go frontend graph was also
made unreachable at compile time.

### Package reachability

| Build | Dependency packages | Wails packages | `net`/`crypto` packages |
|---|---:|---:|---:|
| Monolithic production | 230 | 20 | 77 |
| `wails_native` production | 116 | 15 | 2 |
| Native plus `wails_single_instance` | 149 | 15 | 32 |

The boundary removed 81 packages, including:

- `internal/assetserver`;
- `internal/assetserver/bundledassets`;
- `internal/assetserver/webview`;
- `internal/runtime`;
- `pkg/updater`;
- `net/http`, TLS, X.509, HTTP/2 support, IDNA, and their cryptographic graph.

An earlier native build retained `crypto/rand` because the single-instance
feature was still compiled into `pkg/application`. On Go 1.26 that reached the
FIPS DRBG, whose symbol map contained four 8 MiB zero-filled BSS regions. Those
regions did not add 32 MiB to the executable and were not all resident merely
because they existed, but they illustrated the virtual-address and
initialization consequences of compiling optional features unconditionally.

The native mode now omits single-instance support unless the
`wails_single_instance` tag is supplied. That removed another 167,056 bytes
(4.6% of the previously lean native executable), reduced the dependency graph
from 149 to 116 packages, and removed AES, `crypto/rand`, and the FIPS DRBG from
the default native build. Normal Wails builds remain source- and
behaviour-compatible and continue to include the feature.

`math/rand/v2` is not a safe replacement here. The random value is the 96-bit
nonce for optional AES-256-GCM encryption of second-instance messages. Go's
`cipher.AEAD` contract requires that nonce to be unique for all time under a
given key, while `math/rand/v2` explicitly says its outputs may be predictable
and must not be used for security-sensitive work. A repeated GCM nonce can
break confidentiality and authentication.

Safe alternatives are:

- keep `crypto/rand` when encrypted single-instance support is selected;
- use the OS CSPRNG (`SecRandomCopyBytes`, `BCryptGenRandom`, or `getrandom`) to
  reduce Go crypto reachability on a specific platform;
- redesign encryption as an optional extension package so applications using
  plaintext single-instance messages do not compile AES or a secure RNG.

The OS-CSPRNG option preserves security but adds platform implementations and,
on macOS, Security.framework. It is less valuable than the feature boundary.
`math/rand/v2.ChaCha8` does not eliminate the need for a secure, non-repeating
seed and is therefore not a shortcut.

The repository-wide scan found two other production imports and several test
imports:

| Location | Purpose | Assessment |
|---|---|---|
| `pkg/application/single_instance.go` | AES-GCM nonce | Security-critical; keep a CSPRNG |
| `pkg/application/mcp_eval_enabled.go` | Correlates an HTTP/JavaScript evaluation response | Keep unpredictable because the identifier crosses a process/HTTP boundary; compiled only with `mcp` |
| `internal/commands/updater_tool.go` | Generates Ed25519 signing keys | Security-critical; keep `crypto/rand` |
| `internal/uuid/uuid.go` | Produces a v4 UUID for Windows notifications | Not authentication, but a conforming random UUID should remain CSPRNG-backed; not reachable in the macOS native editor |
| updater tests | Random test keys and payloads | Test-only; no shipped executable cost |

Consequently there is no production `crypto/rand` call that should simply be
changed to `math/rand/v2`. The useful optimization is to stop unrelated
applications from reaching the feature that needs it. On macOS, the ordinary
Wails example reached `crypto/rand` through single-instance support; the lean
native build now reaches none of the sources above.

### Text and metadata, not just function bodies

The unstripped normal toolbar executable had a 7,077,888-byte `__TEXT` segment;
the native build had 2,949,120 bytes. Approximate text-symbol attribution for
the normal build was:

| Namespace | Text bytes | Share of attributed text |
|---|---:|---:|
| Network and crypto | 762,736 | 22.3% |
| Go runtime | 493,444 | 14.4% |
| Wails application | 435,136 | 12.7% |
| Encoding | 165,696 | 4.8% |
| Wails updater | 44,672 | 1.3% |
| Wails asset server | 30,256 | 0.9% |
| Example application | 35,472 | 1.0% |
| Other standard/library code | 1,442,572 | 42.2% |

These figures count regular text symbols and are useful for direction, not
binary accounting. Removing a package also reduces strings, method tables,
interfaces, stack maps, and linker metadata. In the normal unstripped build,
notable metadata included:

- `runtime.pclntab`: 2,637,568 bytes;
- type/relative-type data: about 890,544 bytes;
- function metadata: about 482,280 bytes;
- Go string data: about 409,800 bytes.

In the first native symbol build (which still included optional
single-instance encryption), those corresponding values fell to roughly
1,127,568, 443,904, 222,024, and 62,264 bytes. After making single-instance
support optional, `runtime.pclntab` fell again to 1,067,536 bytes and Go string
data to 59,704 bytes. This is why package reachability can save far more than
the visible body of the Wails function being removed.

### Embedded assets

The following source assets are relevant:

| Asset | Raw bytes | Inclusion |
|---|---:|---|
| `pkg/application/assets/alpha/index.html` | 202,855 | Currently embedded by the core application package |
| Production `runtime.js` | 45,945 | Required by a normal frontend build |
| Debug `runtime.debug.js` | 420,561 | Non-production only |
| `mac-toolbar` demo HTML | 14,655 | Application-owned |

The alpha page is present in the normal production binary even when the app
provides its own asset handler; `strings` and the embed symbol table confirm
it. It should move out of the core package and into the alpha template/example
that actually uses it. This is a low-risk saving of roughly 200 KB raw, plus a
small amount of embed metadata, for every normal Wails application.

The 45.9 KB production runtime is not a meaningful first target. It is used by
normal Wails applications and is already far smaller than the debug runtime.

## Runtime behaviour

### Native Wails versus Swift

One settled three-second observation of the production controls recorded:

| Metric | Native Wails | Swift |
|---|---:|---:|
| Physical footprint at start | 25,740,320 bytes | 22,315,824 bytes |
| Sampled main-process threads | 14 | 6 |
| User CPU during observation | effectively zero | effectively zero |
| Interrupt wakeups | 1 | 1 |

The exact memory values vary across launches and OS state. The longer benchmark
in `performance.md` found a normal-document gap of roughly 4.7 MiB on this
machine. The durable conclusion is that the native Wails memory premium is a
few MiB, while idle CPU is negligible.

The five-second stack samples were more informative than a CPU percentage:

- both main threads spent essentially every sample waiting in
  `-[NSApplication run]` / `mach_msg2_trap`;
- native Wails' additional Go threads were parked in condition waits, `kevent`,
  or work queues;
- Swift had fewer runtime threads, but was waiting in the same AppKit loop;
- there was no active Go collector in the sample.

This is healthy idle behaviour. Reducing parked threads may improve process
shape and stack reservation, but it is not an urgent CPU optimization.

### Normal Wails/WebKit process shape

The normal toolbar demo uses an 1180×760 WebView with native toolbar, split
view, blur/accessory effects, and application content. It is intentionally not
an equivalent UI to the 980×680 text controls, so its memory is an anatomy
sample rather than a fair Wails-versus-Swift score.

After five seconds it had:

- 35 sampled threads in the Wails UI process;
- a WebKit GPU XPC process with 9 sampled threads;
- a WebKit Networking XPC process with 6 sampled threads;
- a WebKit WebContent XPC process with 11 sampled threads.

The main process settled around 55 MB physical footprint in one run. The XPC
snapshot reported roughly 38 MB for WebContent and 6.7 MB for Networking. GPU
footprint accounting was much larger and included purgeable/shared graphics
resources; it must not be naively added to application-private memory.

All four process samples were primarily blocked on kernel waits. The normal UI
process nevertheless produced materially more interrupts/wakeups than the
native controls in these short observations. That work is associated with
WebKit, JavaScriptCore, WebCore scrolling/compositing, and the animated window
surface, not with Go GC. A static hidden/visible benchmark should be added
before turning the observed wakeup count into a release claim.

### What the Go garbage collector did

With a sample note, `GODEBUG=gctrace=1` printed no collections during startup
or the 20-second idle window. The heap stayed below the trigger.

With a generated 10 MiB UTF-8 document, three collections occurred by 11 ms:

| Collection | Approximate heap transition | Wall time |
|---|---|---:|
| 1 | 11 MB → 10 MB live | <0.5 ms |
| 2 | 20 MB → 10 MB live | <0.3 ms |
| 3 | 20 MB → 10 MB live | <0.2 ms |

A forced fourth collection after the native UI was installed also completed in
well under 0.5 ms. It did not erase the large-document footprint. The important
cost was allocation and copying:

1. `os.ReadFile` owns a Go byte slice;
2. `string(data)` creates a Go string copy;
3. `C.CString` creates a C UTF-8 copy;
4. `NSString initWithUTF8String:` decodes into Foundation storage;
5. `NSTextView`/TextKit creates text and layout state.

The editor now drops its Go staging string once AppKit accepts it, which fixed
an avoidable permanent copy. It cannot remove the transient peak or the native
text-layout cost. Swift avoids the Go/C boundary and can hand Foundation-owned
data through fewer representations.

The same lesson applies to ordinary Wails bindings. A large value may exist as
Go data, encoded JSON, HTTP/custom-scheme response bytes, WebKit IPC data, a
JavaScript string/object, and rendered DOM state. Optimizing the collector will
not fix multiple simultaneous representations.

### `GOMAXPROCS`

The machine exposed 14 logical execution contexts to Go. Repeating the idle
sample with `GOMAXPROCS=2` produced:

- no change in the native editor's 14 sampled threads;
- no material native physical-footprint change;
- a reduction from 35 to 27 threads in one normal main-process sample;
- no repeatable CPU or memory improvement in the normal sample.

AppKit, WebKit, JavaScriptCore, cgo, dispatch queues, and Go all create or use
threads. `GOMAXPROCS` controls simultaneous Go execution, not the complete
process thread count. Wails should expose measurements and allow application
owners to tune it, but should not choose a small value globally.

## Where normal Wails 3 can improve

### P0: preserve release-build hygiene

Status: mostly already implemented.

Release tasks must continue to apply:

```text
-tags production -trimpath -buildvcs=false -ldflags="-w -s"
```

Why it matters: omitting `production` increased the measured monolithic editor
from 8.66 MB to 10.47 MB and embedded the 420.6 KB debug runtime instead of the
45.9 KB production runtime.

Action:

- add CI smoke builds that inspect build tags and executable size;
- fail if a release artifact contains `runtime.debug.js`, devtools sources, or
  local source paths;
- record size by platform and architecture to catch regressions.

### P1: remove unconditional alpha assets

Status: clear, low-risk normal-Wails win.

Move `pkg/application/assets/alpha` and `AlphaAssets` to a separate template or
explicit compatibility package. A normal app that supplies `Assets.Handler`
should not carry an unused 202,855-byte HTML page.

Acceptance test:

- build `examples/mac-toolbar` before and after;
- verify `assets/alpha/index.html` and `<title>Wails Alpha</title>` are absent;
- verify alpha scaffolding still works by explicitly importing/embedding it.

### P1: make the updater opt-in at build time

Status: architectural opportunity.

`App` currently exposes a concrete updater and initializes it eagerly. That
forces `pkg/application` to import `pkg/updater`, even for applications that do
not update themselves. The updater's own text is modest, but it reaches HTTP,
verification, archive, crypto, and WebView-window functionality; much of that
overlaps the normal frontend graph, so the exact incremental saving needs an
isolated build experiment.

Preferred v4 shape:

- updater is a separate optional module installed into an application;
- the core app exposes a small extension/host interface, not a concrete updater
  field;
- headless and WebView updater UIs are separate choices.

Practical v3 experiment:

- add a `wails_no_updater` tag or inverse `wails_updater` tag;
- compare `go list -deps`, text sections, metadata, and linkage;
- retain helper-mode support only when the updater is included.

### P1: split HTTP compatibility from the internal WebView transport

Status: largest plausible normal-binary opportunity, but API-sensitive.

The normal build attributed about 763 KB of text symbols to network and crypto,
before counting their strings and type metadata. Wails uses `net/http` types for
asset handlers, middleware, services, and transport even when the application
opens no network listener. Importing `net/http` reaches a broad standard-library
graph including TLS and X.509.

A production design should preserve today's HTTP-compatible API while allowing
a lean mode:

- define a small internal request/response/router abstraction for the
  `wails://` or equivalent custom-scheme path;
- put the `net/http` adapter, HTTP services, and server transport in an optional
  package/build feature;
- keep HTTP mode as the compatibility default in v3;
- allow v4 applications to choose `webview transport` versus `HTTP/server
  transport` explicitly.

The 763 KB text category is an upper bound, not a promised saving. Some crypto
and network code may remain reachable through updater, user code, remote URLs,
or platform services.

### P1: make platform managers and cgo frameworks optional

Status: normal application package remains broad.

Even the native tagged editor links ServiceManagement, UniformTypeIdentifiers,
and Carbon because autostart, sharing/type handling, and global shortcuts live
in the same application package. Go can dead-strip many methods, but cgo/ObjC
translation units, package initializers, framework flags, and public concrete
types still broaden the artifact.

Candidate modules:

- WebView;
- updater;
- autostart / ServiceManagement;
- global shortcuts / Carbon;
- sharing and type-identification support;
- server and HTTP transport;
- MCP/debug tooling.

Acceptance criterion: an application that does not select a feature should not
show its package in `go list -deps` or its framework in `otool -L`.

### P1: create WebViews lazily

Status: highest likely memory/process win for tray and occasionally-opened
normal Wails apps.

Linking WebKit is cheap; constructing `WKWebView` is not. The first WebView
creates the UI-process machinery and GPU, Networking, and WebContent XPC
processes. A tray application that may never open its window should be able to
define a window without constructing its WebView until `Show` or first load.

Requirements:

- preserve option mutation before realization;
- define whether preload is desired for latency-sensitive applications;
- make close/destroy semantics explicit so XPC resources can be released;
- benchmark first-show latency separately from application-launch latency.

Acceptance tests:

- a hidden, unrealized WebView window creates no WebKit XPC process;
- first `Show` creates exactly the expected process set;
- an opt-in preload mode preserves current immediate-load behaviour.

### P1: add byte/file-oriented large-data paths

Status: proven by the 10 MiB editor experiment.

For native text, add Foundation-owned file/data operations such as
`LoadFile`/`SaveFile` or a provider returning bytes with explicit ownership.
For normal Wails transport, add binary and streaming paths that avoid:

- `[]byte` to `string` conversions;
- base64 for binary payloads;
- complete JSON materialization when a stream or file is intended;
- repeated full-document callbacks on every edit;
- redundant Go and JavaScript copies of immutable data.

Benchmark at 1, 10, and 100 MiB and record peak physical footprint, lifetime
peak, bytes allocated, collections, and end-to-end latency. The current native
10 MiB footprint gap should be treated as a transfer-API problem.

### P2: lazy initialization inside `application.New`

Status: likely small individually, worthwhile after package boundaries.

`App.init` creates every manager, and normal initialization creates event,
message, binding, transport, asset-server, and updater state before the app
runs. Most constructors are cheap, and the no-GC startup result means this is
not currently a major pause source. Still, Wails can avoid work when:

- there are no services to bind;
- no updater is configured;
- no system tray, global shortcut, autostart, or dialog is used;
- a WebView has not yet been realized;
- a custom transport makes the default transport unnecessary.

Use lazy initialization to reduce retained objects and startup dependencies,
not as a substitute for compile-time modularity.

### P2: consolidate framework event readers

Status: process-shape improvement, low expected CPU return.

`App.Run` starts separate goroutines for application events, window events,
WebView requests, messages, keys, drops, toolbar actions, sharing, split panes,
sidebars, inspectors, text editors, native-window close events, and menus. Go
stacks begin small and these goroutines were parked, so the expected memory
win is modest.

A typed central dispatcher or demand-started readers could reduce goroutine
count and simplify shutdown. It should only be pursued with race, latency, and
backpressure tests; converting independent bounded channels into one global
queue can create head-of-line blocking.

### P2: profile the JavaScript/WebKit side independently

Status: required for normal-app CPU work.

The Go UI process was idle while JavaScriptCore and WebCore threads existed in
the same process and three XPC helpers existed outside it. A Go CPU profile
alone cannot explain normal Wails performance.

A standard Wails performance suite should capture:

- Go CPU, allocation, mutex, block, and goroutine profiles;
- Web Inspector timeline and JavaScript heap;
- samples of the UI, WebContent, GPU, and Networking processes;
- Core Animation frame pacing and layout/paint invalidations;
- idle wakeups with the window visible, occluded, hidden, and minimized.

The suite needs a static baseline page and a representative framework frontend
so application animation is not mistaken for Wails overhead.

### P3: do not make GC or `GOMAXPROCS` changes by default

Status: unsupported by current evidence.

The sample editor did no normal-document GC, large-document collections were
sub-millisecond, and lowering `GOMAXPROCS` did not improve the native process.
Global `GOGC`, `GOMEMLIMIT`, or `GOMAXPROCS` defaults could hurt compute-heavy
applications and hide copy/retention bugs.

Wails should instead:

- document these Go controls for application owners;
- expose optional runtime metrics in diagnostic builds;
- use allocation profiles to fix retained or repeated values;
- set no framework-wide concurrency policy without broad cross-platform data.

## Proposed normal-Wails benchmark matrix

Every optimization above should be tested against the same matrix:

| Scenario | Why |
|---|---|
| Empty production WebView, visible | Framework baseline |
| Same window hidden and minimized | Detect unnecessary rendering/wakeups |
| Tray app before first `Show` | Validate lazy WebView creation |
| One and five WebView windows | Process-pool and per-window scaling |
| No services vs 100 bound methods | Binding/reflection cost |
| 1, 10, 100 MiB text and binary payloads | Copying and GC scaling |
| Updater absent vs configured | Optional dependency and runtime cost |
| Default HTTP compatibility vs lean scheme transport | Network/crypto graph |
| `production` vs debug | Release hygiene guardrail |

For each row record:

- stripped executable and bundle size;
- package list and linked frameworks;
- launch wall/user/system time;
- UI-process and helper-process physical footprint;
- lifetime peak footprint;
- thread, goroutine, and XPC process count;
- Go allocations and GC pause/CPU;
- JavaScript heap and GC;
- idle and interrupt wakeups;
- first-window and first-interaction latency.

Results should be stored as CSV/JSON artifacts and summarized from medians and
percentiles over repeated runs. A single Activity Monitor number is not enough.

## Recommended order of work

1. Remove unconditional alpha assets and add production size/linkage CI.
2. Land and validate the `wails_native` boundary as the reference modularity
   experiment.
3. Isolate updater cost with an opt-out build and then redesign it as an
   optional extension.
4. Prototype a lean internal WebView request transport with an optional
   `net/http` adapter.
5. Add lazy WebView realization for hidden/tray-window workflows.
6. Add byte/file/stream transfer APIs and rerun the large-payload matrix.
7. Split platform managers/frameworks where measurements demonstrate a useful
   artifact or startup reduction.
8. Revisit event goroutines and runtime tuning only after the larger boundaries
   are complete.

## Bottom line

The native experiment did not show that Wails should abandon WebKit for normal
applications, nor that Go GC is the dominant cost. It showed that Wails 3's
current package shape makes runtime-optional features build-mandatory, and that
large values become expensive when copied across runtime boundaries.

For normal Wails 3, the immediate release gains are modest but concrete:
remove unused embedded assets, enforce production builds, make hidden WebViews
lazy, and begin isolating updater/HTTP/platform features. The larger v4 gain is
an application core whose dependency graph is assembled from selected window,
transport, updater, and native-control modules. That structure lets Go's linker
do the optimization it is already good at: removing code that was never made
reachable in the first place.
