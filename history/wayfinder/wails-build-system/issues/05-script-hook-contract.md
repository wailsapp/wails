# Script hook contract and safe extension boundary

Type: grilling
Status: open
Blocked by: 02

## Question

What phases may invoke user-owned scripts, what arguments and environment are
stable, how are target/profile/output values passed, how are failures reported,
and what input/output declaration is required for a hook to participate in the
cache? The cache contract must treat the referenced script content and relevant
executable metadata as implicit inputs. Confirm that inline shell remains
unsupported in the new manifest.

## Comments

### 2026-08-16 — Accepted implementation default

- Stable phases are `before_build`, `after_build`, `before_package`,
  `after_package`, `before_sign`, and `after_sign`, with one script per phase.
- Scripts receive no interpolated shell command. They are executed directly
  from the Manifest directory with the stable `WAILS_*` environment already
  specified in the proposal; the long form may override the directory.
- Hooks inherit the phase's Project, Target, or Package scope. A non-zero exit
  fails that scope and is reported as a normal Step failure.
- Hooks always run unless the long form explicitly sets `cache = true` and
  declares complete `inputs` and `outputs`. The script bytes, executable mode,
  phase, environment contract version, and resolved scope are implicit cache
  inputs. Inline shell remains unsupported.
