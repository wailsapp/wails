# Script hook contract and safe extension boundary

Type: grilling
Status: resolved
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
  from the Manifest directory. One generated JSON context file is exposed as
  `WAILS_HOOK_CONTEXT_FILE`; the long form may override the directory.
- Hooks inherit the phase's Project, Target, or Package scope. A non-zero exit
  fails that scope and is reported as a normal Step failure.
- Hooks always run unless the long form explicitly sets `cache = true` and
  declares complete `inputs` and `outputs`. The script bytes, executable mode,
  phase, context contract version, and resolved scope are implicit cache
  inputs. Inline shell remains unsupported.

## Answer

The six hook phases are implemented as typed `RunHook` Nodes. `before_build`
runs once for the shared Project with empty target/output values;
`after_build` runs once per Target with its binary; package and signing hooks
form one before/after barrier per Target around the complete requested format
set, receiving the package path for one format or the common package parent for
several. Independent package Nodes remain concurrent inside those barriers.

Scripts are project-relative files, not inline shell. They run from the Project
or an explicit project-contained directory. The only Wails-owned context
variable is `WAILS_HOOK_CONTEXT_FILE`, which points to an atomically written,
owner-only JSON file under `.wails/hooks/context/`. Schema version 1 contains
the phase, command, scope, Project and working directories, profile, optional
Target, scoped output, and declared outputs. It contains no inherited
environment or secrets. Success and cancellation remove the file; a genuine
failure retains it and reports its location. Unix uses the script's
executable/shebang; Windows passes `.cmd`, `.bat`, and `.ps1` paths as one
interpreter argument. Non-zero exits retain their code and output as normal
Step failures, and cancellation terminates the script's process group.

Hooks run every time unless `cache = true` and complete inputs and outputs are
declared. Their cache identity implicitly includes script bytes/mode, phase,
context-contract version, and resolved scope. The random generated context path
is not part of cache identity. Output restoration is
content-addressed; output roots cannot contain scripts or inputs. Script,
directory, input, and output paths are constrained to the Project after
symlink resolution. Migration translates only conservative, argument-free
script-file tasks; inline shell remains a manual diagnostic.

Focused tests cover all phase/barrier shapes, multi-Target project sharing,
cache hit/restore/invalidation, default non-caching, context isolation/lifecycle,
path escapes, process-tree cancellation, failure propagation, and migration.
Manifest, pipeline, cache, and command tests pass normally and under the race
detector; vet and a Windows cross-compile also pass.
