# Wake manifest cache MVP — throwaway prototype

This prototype answers one question: can a minimal `wails.toml` drive a real
bindings → frontend → Go binary build while making no-op and narrowly scoped
incremental builds fast enough to justify replacing generated Taskfiles?

It is deliberately limited to one host desktop Target, npm, TypeScript
interface bindings, and the production build mode. It is not a production API.

Build the prototype CLI, then run the normal command from a project containing
the minimal manifest:

```bash
cd v3
go build -o /tmp/wails3-wake-mvp ./cmd/wails3
cd examples/badge
WAILS_WAKE_MVP=1 /tmp/wails3-wake-mvp build
```

The full relevant state is visible after every run through Pulse: each Node is
shown as executed or cached, the total build time includes planning and cache
validation, and the resulting binary is listed as an Artifact. Prototype cache
metadata is written under the project's `.wails/`; Artifact bytes use the
operating system's user-cache directory under `wails/wake-mvp/artifacts`.

The code is intentionally marked `mvpprototype` and lives only on the local
throwaway prototype branch. Review the timings and invalidation behavior before
lifting any part into the production Wake engine.
