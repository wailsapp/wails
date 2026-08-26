# Actionable HCL configuration diagnostics

Type: task
Status: resolved
Blocked by: none
Label: ready-for-agent

## Question

Make invalid `wails.hcl` feedback consistently actionable. Every syntax,
schema, semantic, path, target, format, signing and host-capability error that
originates from configuration should identify the manifest, exact source range,
field, invalid value where safe, and a concise correction or supported-value
hint. Build and Dev must render these diagnostics before executing work and
must never fall back to Taskfiles.

Evaluate whether a separate `wails3 config check` materially improves editor,
CI or preflight workflows beyond the normal command and `--plan` validation
surface. Add it only if it provides a distinct user benefit; excellent
build-time diagnostics are mandatory either way. Cover text output,
machine-readable output if introduced, terminal capability differences,
goldens, source-range accuracy, unknown-field suggestions, multiple errors,
security redaction and diagnostic-path performance.

## Comments

### 2026-08-26 — Added to the resumed goal

The current implementation reports a source location for many parse/schema
errors, but presentation is a single compact CLI error. This ticket covers the
complete diagnostic experience rather than only adding another command.

## Answer

Implemented. Manifest validation errors render the exact file,
line and column, source line, caret range, field, detail, and a concise hint
for supported-value and unknown-field cases. The same renderer is used by
normal build and Dev command failures before pipeline execution.

`wails3 config check [profile]` was added because it has distinct CI/editor
value: without changing build state it validates the base manifest and either
one requested profile or every named profile, including structural Plan
selection. Explicit block-label and attribute origins now survive semantic
validation, Plan construction, signing checks and configuration-owned host
capability failures. Missing fields fall back to their owning block; failures
caused only by CLI selections remain ordinary command errors.

Focused tests cover source extraction, caret accuracy, actionable hints,
unknown fields, joined diagnostics, non-manifest fallback, terminal-neutral
text, golden output and environment-value redaction. Source files are loaded at
most once while rendering joined failures. On the Linux audit host, the
no-source renderer benchmark measured 1.262–1.345 microseconds per diagnostic
with 512 bytes and 10 allocations. Real CLI acceptance verified exact ranges
for an invalid profile format and a cross-build host/toolchain mismatch.

No additional JSON format was introduced: `config check` has a stable exit
status and human diagnostic stream, while the already-versioned
`build --plan --json` contract remains the machine-readable planning surface.
