# Pulse — the wake build TUI

## What this is

`pulse` is a new renderer for the `report.Reporter` contract: a calm,
information-dense terminal UI for wake builds with a pinned live region above
the scrolling log of completed steps. It replaces nothing — `report/termui` keeps
working. Pulse is designed to be the *premium* render path: enabled when stdout
is a TTY and the user hasn't asked for plain output. The non-TTY fallback
degrades to the same one-line-per-step structure as today.

It has **zero third-party rendering dependencies**. Colour, cursor control,
synchronized output, and width measurement are implemented in
`pulse/ansi` (≈80 lines) and `pulse/color.go` (the palette + profile detection).
The only non-stdlib import is `golang.org/x/term` for `IsTerminal` / `GetSize`,
the standard Go ecosystem escape hatch for TTY operations.

## Design philosophy

1. **Calm by default, dense on request.** The serial path looks like today —
   one spinner row, one progress strip, one summary block. Verbose mode adds
   commands and streamed output. Debug mode adds DAG/dep/var trace lines.
   Nothing surprising fires without a flag.

2. **Pinned live region, scrolling history above.** The Buck2 Superconsole
   model: the active steps and the progress bar sit at the bottom and are
   redrawn in place; completed steps print *above* the region and roll off
   into scrollback as the terminal scrolls. The scrollback you see at the end
   of a build is plain, copyable text — no escape soup, no orphaned cursor
   moves.

3. **Atomic frame updates.** Every paint cycle is wrapped in DEC mode 2026
   (synchronized output). Modern terminals (Ghostty, Kitty, WezTerm, iTerm2)
   buffer the writes and flip the affected region atomically; older terminals
   ignore the brackets and behave as today. Sync brackets are *only* emitted
   when the pinned region is involved — pure log writes in non-TTY mode stay
   clean.

4. **One signature accent.** Electric blue is the *only* saturated brand
   colour. Success, failure, warning, cached, and dim are all muted so the
   accent always pulls the eye to "what's happening right now" (the spinner,
   the progress head, the critical-step marker).

5. **No emoji as status glyphs.** Unicode geometric shapes (`✓ ✗ ⊙ · ◆ ◇ ◀`)
   are single-cell, monospace-safe, and ship in every terminal font. Emoji
   are variable-width and render inconsistently — they break the alignment
   the design depends on.

## Mockups

### Serial build, in flight (animated)

```
  build › darwin:build (8 steps)
  ━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━

  ⊙ [1/8] go:mod:tidy                                                    cached
  ✓ [2/8] generate:bindings                                              480ms
  ✓ [3/8] build:frontend                                                  1.3s
  ⠦ [4/8] build:amd64  • linking…                                         1.9s
  ━━━━━━━━━━━━━━━━━━▋──────────────────────────────────────  47%  ETA 2.4s
```

The accent rule under the header is wake's *signature stroke* — a single
anchor for the build's identity. The bottom two lines are the pinned region;
they redraw on an 80 ms tick. Everything else is plain scrollback.

### Parallel build, four workers

```
  build › darwin:release (12 steps)

  ⊙ [1/12] go:mod:tidy                                                   cached
  ✓ [2/12] generate:bindings                                             460ms
  ✓ [3/12] build:frontend                                                 1.2s

    Active
  ⠹ [4/12] build:amd64           • compiling 248 packages                 2.1s
  ⠹ [5/12] build:arm64           • compiling 248 packages                 2.0s
  ⠼ [6/12] test:unit             • running 184 tests                      1.7s
  ⠴ [7/12] test:integration      • running 24 tests                       1.4s

  ━━━━━━━━━━━━━━━╾────────────────────────────────────  31%  ETA 4.6s
  throughput  ▁▁▂▃▄▆▇█▇▆▄▃▂▁
```

The "Active" panel is unbounded (the worker pool sets the practical cap); each
spinner ticks in unison so the eye sees one heartbeat rather than four. The
throughput sparkline shows completions over a rolling ~10 s window.

### End of build, success

```
  ━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━

  ✓  build succeeded    3.6s  11.6s cpu · 3.2× speedup

  11 ran  •  1 cached

  Where time went  ━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
                   compile 5.5s  test 4.3s  package 1.2s  prepare 460ms

  Slowest
  build:amd64        ━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━    2.1s  ◀ critical
  build:arm64        ━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━─    2.0s
  build:frontend     ━━━━━━━━━━━━━━━━━━━━━━─────────────    1.3s
  codesign           ━━━━━━━━━━━━━━━────────────────────   920ms
  generate:bindings  ━━━━━━━━───────────────────────────   480ms
```

The summary tells the build's story: verdict, headline counts, where time
went, and which step was the bottleneck. The verdict line shows wall-clock
duration plus, when meaningfully different, cumulative CPU time and the
parallel speedup ratio — `3.2×` immediately tells you the worker pool earned
its keep. The bar chart is sorted descending and capped at 5 rows. For a
serial build the slowest step *is* the critical path; for a parallel build
it's the bottleneck — both are the most actionable single answer to "what
should I look at next?".

Zero-count entries are hidden ("0 cached" is noise), but "ran" always shows,
because zero ran-and-passed steps is itself informative.

### End of build, single failure

```
  ✗  build failed    after 2.2s

  1 ran  •  1 cached  •  1 failed

  Slowest
  test:unit          ━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━    2.2s  ◀ critical
  generate:bindings  ━━━━━━─────────────────────────────   380ms

  ╭─ test:unit ────────────────────────────────────────────────────────────────╮
  │ command   go test -count=1 ./internal/generator/...                        │
  │ status    exited 1                                                         │
  │                                                                            │
  │ === RUN   TestAnalyser/Service13                                           │
  │     analyser_test.go:142: missing service "Service13" in generated bindings│
  │     analyser_test.go:143:   want: ["Service11" "Service12" "Service13"]    │
  │     analyser_test.go:144:   got:  ["Service11" "Service12"]                │
  │ --- FAIL: TestAnalyser (0.42s)                                             │
  │     --- FAIL: TestAnalyser/Service13 (0.18s)                               │
  │ FAIL                                                                       │
  ╰────────────────────────────────────────────────────────────────────────────╯
```

The panel is right-sized to its widest content line; never edge-to-edge, so
it reads as a *panel* rather than a banner. Border is dim red, header is
bold red, body content is unstyled.

### End of build, multiple failures

```
  ✗  build failed    after 3.6s

  8 ran  •  1 cached  •  3 failed

  Slowest                                          (omitted for brevity)

  ✗  3 tasks failed

  ╭─ darwin:test:unit ─────────────────────────────────────────────────────────╮
  │ ...                                                                        │
  ╰────────────────────────────────────────────────────────────────────────────╯

  ╭─ darwin:vet ───────────────────────────────────────────────────────────────╮
  │ ...                                                                        │
  ╰────────────────────────────────────────────────────────────────────────────╯

  ╭─ darwin:package:dmg ───────────────────────────────────────────────────────╮
  │ ...                                                                        │
  ╰────────────────────────────────────────────────────────────────────────────╯
```

Multi-failure builds get an aggregation banner *before* the panel sequence so
the reader's first scroll-up shows the headline count, then one panel per
failure in order of occurrence.

## Visual vocabulary

### Palette

| Role     | TrueColor | 256-color | 16-color | Used for                                |
| -------- | --------- | --------- | -------- | --------------------------------------- |
| Accent   | `#7d9eff` | 111       | 94       | verb, spinner, progress fill, ◀ critical |
| Success  | `#7ec682` | 114       | 92       | ✓, "build succeeded", "N ran"           |
| Failure  | `#e06c75` | 167       | 91       | ✗, panel borders, "N failed"            |
| Warning  | `#e5c07b` | 179       | 93       | (reserved; debug `var` badge)           |
| Cached   | `#8ec5d8` | 109       | 96       | ⊙, "cached", "N cached"                 |
| Dim      | `#6e7785` | 244       | 90       | counters, separators, dependency arrows |
| Subtle   | `#9aa5b8` | 248       | 37       | secondary text (durations)              |

`DetectProfile` picks the richest tier the terminal advertises:

- `NO_COLOR` set → `None` (no SGR emitted at all)
- `COLORTERM=truecolor|24bit` → 24-bit
- `TERM` contains `256color` → 256
- `TERM` contains `kitty|ghostty|alacritty|wezterm` → 24-bit
- otherwise → 16-color
- `TERM=dumb` or stdout is not a TTY → None

### Glyphs

| Glyph | Meaning                       |
| ----- | ----------------------------- |
| `✓`   | step succeeded                |
| `✗`   | step failed                   |
| `⊙`   | cache hit                     |
| `·`   | skipped (up-to-date, platform mismatch) |
| `›`   | header arrow (`verb › target`) |
| `•`   | inline detail bullet          |
| `◆`   | active step (reserved)        |
| `◇`   | queued step (reserved)        |
| `◀`   | critical / bottleneck marker  |
| `━`   | progress bar filled           |
| `╾`   | progress bar head (transition glyph) |
| `─`   | progress bar empty            |
| `▁..█`| sparkline tier (8 levels)     |

Spinner: Wake's **heartbeat** — a vertical block that swells from a low
bar to full height and back down (`▁ ▂ ▃ ▄ ▅ ▆ ▇ █ ▇ ▆ ▅ ▄ ▃ ▂`),
80 ms per frame, ~1.1 s per beat (roughly a resting human pulse). This is
what earns the renderer its name. We deliberately don't use the generic
braille loading-spinner — that signifies "loading"; this signifies *wake*.
All in-flight steps share a single frame counter so the build pulses in
unison, like a chorus rather than a crowd.

Progress bar: thick horizontal stroke with **eighths-block sub-cell head**
(`━━━━▋───`). At each tick the head advances in 1/8-cell increments rather
than whole-cell jumps, so a slow build's bar still feels alive between
integer percent crossings.

### Box-drawing

Used in exactly one place: the failure panel border (`╭ ╮ ╰ ╯ ─ │`). Everywhere
else, structure comes from whitespace + indentation. The intent is that the
panel reads as a *capsule* — a discrete, frameable thing — while everything
else flows.

## Animation timing

| Cadence | What it drives                                       |
| ------- | ---------------------------------------------------- |
| 80 ms   | Spinner frame advance; pinned region full repaint    |
| 250 ms  | Throughput window roll-over (one bar per 250 ms slice) |
| Event   | StepEnd / StepFailed promote the completed line and force-repaint |

The repaint goroutine acquires the same `sync.Mutex` as the public interface
methods. A single ticker drives both serial and parallel modes.

## Interface fit and extensions

The Pulse Reporter satisfies `report.Reporter` exactly as documented:
serial-only, one in-flight step at a time, with the canonical lifecycle of
`StepStart → (StepInfo|StepCommand|StepOutput)* → (StepEnd|StepFailed)`.

For **parallel builds**, the Reporter contract is single-threaded by design.
Pulse adds an additive concrete-type API (no interface change):

```go
type StepID uint64
func (r *Reporter) ParallelStepStart(name, label string) StepID
func (r *Reporter) ParallelStepInfo(id StepID, msg string)
func (r *Reporter) ParallelStepEnd(id StepID, status report.Status, dur time.Duration)
func (r *Reporter) ParallelStepFailed(id StepID, f report.Failure)
```

Callers running multiple workers hand each worker the `StepID` from its own
`ParallelStepStart` and use that to address subsequent updates. Mixing the
serial and parallel paths in one build is supported (the serial path is just
an implicit single-worker case).

If we promote this into the `report.Reporter` interface later, the natural
shape is to make `StepStart` return a `StepID` and have every subsequent
step-scoped method take it — the serial caller can ignore the return value
and the parallel caller threads it through.

### Counter semantics in parallel mode

The `[N/M]` counter is assigned at `StepStart` time, so in a parallel run the
counters in the scrolling log are in **launch order, not completion order**:
the line `[6/12] test:unit (1.7s)` may appear above `[4/12] build:amd64 (2.1s)`
if test:unit started later but finished first. This is intentional — the
counter is "which step number is this?", which is stable for the lifetime of
the step. The visible interleaving is the build pipeline showing its shape.

### What's *not* in scope

These were considered and deliberately deferred:

- **Per-task PTY-emulated log panes (Turborepo-style split).** Requires an
  interactive Bubbletea-style event loop, mouse/keyboard handling, and a
  stripped vt100 parser. Wake's output is line-oriented; a one-shot animated
  inline TUI gives 90% of the value at 5% of the surface.

- **Critical-path-by-dependencies.** The current `report.Reporter` doesn't
  receive the DAG shape, only `totalSteps`. The "Slowest" panel uses wall-time
  ordering as a proxy. For a true critical path we'd need to push DAG edges
  into the reporter; this is the right next interface extension if we want it.

- **OSC 8 hyperlinks on step names.** The wake reporter doesn't currently know
  the source file for a task definition. Worth wiring in once the resolver
  exposes that; a one-line addition through the `styler.link` helper.

- **Inline-image bundle visualisations** (Kitty/iTerm graphics protocol).
  Tempting, brittle in tmux / SSH / VSCode terminal, off-brand for "calm
  default". Belongs behind a future `--graphics` flag if at all.

## Implementation map

| File                        | Lines | Purpose                                            |
| --------------------------- | ----- | -------------------------------------------------- |
| `pulse/ansi/ansi.go`        | ~95   | SGR / cursor / sync / terminal-detect helpers      |
| `pulse/color.go`            | ~120  | Profile detection, palette, hex→RGB                |
| `pulse/style.go`            | ~140  | `styler` + width-aware padding/truncation          |
| `pulse/progress.go`         | ~100  | Glyph constants, progress bar, sparkline, barchart |
| `pulse/region.go`           | ~75   | Pinned-region engine (erase+repaint, atomic update)|
| `pulse/pulse.go`            | ~330  | `Reporter` implementation, serial + parallel API   |
| `pulse/render.go`           | ~140  | Per-line renderers, layout math                    |
| `pulse/summary.go`          | ~195  | Verdict + counts + slowest + failure panels        |
| **Total**                   | **~1200** | Pure-stdlib + `golang.org/x/term`              |

## How to inspect the design

The fastest way is to run a real wake build with the pulse renderer on a
generated project:

```bash
wails3 init -n demo -t vanilla && cd demo
WAILS_USE_WAKE=true wails3 build       # cold build — full pulse run
WAILS_USE_WAKE=true wails3 build       # warm — most steps cached, fast verdict
WAILS_USE_WAKE=true WAKE_FORCE=true wails3 build   # forced clean — every step runs
WAILS_USE_WAKE=true WAKE_DEBUG=true wails3 build   # adds resolver internals
```

Failure-panel rendering can be exercised by introducing a transient
syntax error into a Go file (the build's `build:native` step fails and
the panel surfaces the error with OSC 8 file:line links).
