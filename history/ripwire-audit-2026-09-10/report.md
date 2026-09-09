# Ripwire audit of the HCL build-system branch

Date: 10 September 2026. Audited `codex/hcl-build-system` at `6c556f0b9771a2c9cc6b35b74348a1202e2faa52`. The full branch comparison uses merge-base `f38ae268344700a3f48dd1e95b716f2059f7269e`; the incremental comparison starts at the previous review, `906a9aa27`. This report commit does not change application code.

**Result: three actionable findings.** One compiler-cache defect can ship stale native code; two additional issues affect mobile development and the HCL tutorial. The existing tests pass. Ripwire also reports substantial structural debt, which is a prioritization aid rather than a count of bugs or a security certification.

## Installation

Installed official [redhat-et/ripwire v0.5.0](https://github.com/redhat-et/ripwire/releases/tag/v0.5.0), verified against its published SHA-256. The anonymous installer hit GitHub's API rate limit; authenticated release discovery and direct checksum-verified download completed installation.

- Executable: `~/.local/bin/ripwire`, available as `ripwire` on PATH.
- Sixteen usage skills: `~/.agents/skills/ripwire-*`, linked to stable assets under `~/.local/share/ripwire/skills`.
- The skill-installer helper fetched the matching release, then Ripwire's own installer activated its canonical Codex manifest. Duplicate identical copies from the intermediate legacy Codex directory were removed.
- Ripwire's skill scanner: 23 Markdown files scanned, zero findings, zero unscannable files. Its doctor confirms binary, all 22 grammars and exact 16-skill parity.
- Optional advisory hooks and MCP configuration were not enabled. Doctor reports those optional integrations as absent; they are not needed to run the CLI or use the installed skills.
- The contributor-only skill for building Ripwire itself was not activated. The installed usage skills become available on the next Codex turn.

Release and executable hashes are recorded in [measurements.json](measurements.json).

## Standards

**[P1] Compiler cache inputs omit valid CGo source extensions.** Both the project input list in `v3/internal/wake/pipeline/planner.go:316` and the local dependency list at line 834 omit `.cxx` and other Go-supported source types, including Objective-C `.m` and Fortran/SWIG inputs. Editing `helper.cxx` does not change the actual compile-node input snapshots. A warm build can therefore reuse and publish an executable containing old native code. This violates the documented content-input/action-identity contract in `v3/internal/wake/AGENTS.md`.

Reproduction: `python3 history/ripwire-audit-2026-09-10/reproduce-cxx.py` injects an isolated Go test using an overlay and currently fails with `CXX source edit must invalidate compile inputs`. It does not modify package sources. A separate real CGo fixture printed `1`, then `2`, when only its `.cxx` implementation changed, confirming that Go consumes the file Wake overlooks. Align both lists with Go's accepted compilation inputs and add regression coverage before accepting the fix.

No additional code-smell heuristic is promoted to a correctness violation.

## Spec

**[P2] Mobile development replaces an unchanged backend.** `v3/internal/commands/manifest_dev_mobile.go:119` sets `BackendChanged` to unconditional `true`. The controller consequently stops, installs and launches the Android/iOS app after a successful rebuild even when its compile/package results are cached—for example, after a manifest edit that only changes `dev.debounce_ms`. This loses live app state and performs unnecessary deployment. The contract in `v3/internal/wake/AGENTS.md` says: “Use compile and after_build results to avoid replacing an unchanged backend.” Compare the installable artifact identity, including packaging changes. This finding is confirmed from control flow; this audit did not repeat a native mobile session.

**[P2] The primary HCL tutorial omits its required opt-in.** `docs/src/content/docs/experimental/hcl-builds.mdx:23–24` says root-manifest presence selects the native build system. The new gate also requires `WAILS_EXP_USE_WAKE` to be present. Without it, build/dev remain on Taskfiles and migration/config/eject commands are hidden. A real CLI reproduction showed `config check` displaying general help without validating the manifest when the variable was unset; with it enabled, the same command validated the file. The updated Wake instructions explicitly preserve legacy behavior when the variable is unset. Add the environment prerequisite before the tutorial's configuration and commands, and update the manual acceptance guide similarly.

Axis totals: Standards 1 (worst P1); Spec 2 (worst P2). Follow-ups are tracked in [the local audit map](../wayfinder/ripwire-audit-2026-09-10/map.md). They are bugs prioritized for nightly patches, outside the feature-only v3.1.0 milestone/project.

## Ripwire measurements

| Pass | Result | Interpretation |
| --- | --- | --- |
| Repository quality panel | 18,313 eligible functions/methods; 706 meet the two-family threshold | Risk candidates, not confirmed defects; all 706 result rows retained |
| Full branch quality delta | 1,451 rows: 1,198 new-symbol debt, 253 preexisting-worse; 247 gating rows, exit 2 | The default structural gate fails; multiple rows can describe one symbol |
| Since the previous review | 170 rows: 80 new-symbol debt, 90 preexisting-worse; 75 gating rows, exit 2 | Includes complexity/size growth in migration, template installation and packaging |
| Lint | 8,292 reported findings; all locator rows retained | Predominantly naming/style heuristics; naming-confusable capture is capped, so the total is a floor |
| Hotspots | 1,130 ranked files out of a 4,548-file indexed corpus | Recent churn multiplied by complexity; top 30 retained |
| Clones | 23.3% reported duplicated lines | Includes tests, generated code and cross-platform copies; not a removal target without review |
| Test seams | 4,765 edges not reached by statically resolved tests | Not equivalent to actual test coverage; interface dispatch and naming ambiguity limit inference |
| Dead-code candidates | 37 candidates | Verify native callbacks and generated/indirect references before removing anything |
| Branch review context | 219 eligible changed files represented | Complete file list; nested detail trimmed to the requested token budget |
| Commands/Wake panel | 306 candidates out of 2,143 eligible bodies | Separate subsystem scan; top 40 retained |

The strongest build-system maintenance priorities are the multi-platform handler `manifest_pipeline.go`, migration analysis, the planner, and the new development-session integration. Reported examples: `packageIOS` cognitive complexity 54, `analyseMigration` 123, and `planTarget` 156. These numbers suggest where smaller independently testable units would help; they do not justify refactoring blindly before the cache and lifecycle regressions are fixed.

Repository-wide high scores also include Windows window dispatch and generated JavaScript runtime bundles. Those must not be attributed solely to the HCL changes. The graph reports 3,621 ambiguous and 2,608 unresolved references, including dubious cross-version/cross-language matches; graph-based test/dead-code results are leads requiring confirmation.

## Scope, commands and evidence

Scans ran against the complete checked-out Wails repository, including v2 and v3. Standard exclusions were `.git`, `node_modules`, `history`, `dist`, `.wails` and `bin`. Generated assets outside those names remain included. The scan is structural triage, not a dataflow/taint proof or an exhaustive security audit. It does not execute native APIs, prove every cache dependency, or certify absence of defects.

Representative reproduction commands, from the repository root:

```sh
ripwire --version
ripwire . --doctor --agent=codex
ripwire . --quality-panel --limit=10000 --exclude=.git --exclude=node_modules --exclude=history --exclude=dist --exclude=.wails --exclude=bin
ripwire . --quality-delta=f38ae268344700a3f48dd1e95b716f2059f7269e..6c556f0b9
ripwire . --quality-delta=906a9aa27..6c556f0b9
ripwire . --pr-context=f38ae268344700a3f48dd1e95b716f2059f7269e --limit=10000 --max-tokens=100000
```

The same exclusions were applied to every repository-wide audit invocation. Additional passes were `--lint --limit=10000`, `--hotspots --limit=30`, `--deps`, `--clones --limit=50`, `--seams --limit=100`, and `--dead-code --limit=10000`. The subsystem pass was `ripwire v3/internal/commands v3/internal/wake --quality-panel --limit=40`.

[measurements.json](measurements.json) retains the exact counts, floors, truncation disclosures, hashes and top rows. Full raw reports and command logs remain locally in `raw/` (ignored to avoid committing megabytes of generated XML); the metadata and confirmed-bug reproduction are committed. Initial short panel/context runs were expanded before reporting the final results.

## Validation and remaining limits

Passed on Linux at the audited source commit:

```sh
# From v3/
go test ./internal/wake/... ./internal/commands ./cmd/wails3 ./internal/templates ./internal/setupwizard ./internal/features ./internal/report/...
go test -race ./internal/dev/... ./internal/devruntime/... ./internal/wake/... ./internal/commands ./cmd/wails3 ./internal/features ./internal/report/...
go vet ./internal/dev/... ./internal/devruntime/... ./internal/wake/... ./internal/commands ./cmd/wails3 ./internal/templates ./internal/setupwizard ./internal/features ./internal/report/...
```

The added `.cxx` audit reproduction intentionally fails and exposes a gap in the passing suite. Earlier test failures at `03ca64342` were superseded by the passing rerun at `6c556f0b9`; they are not reported as current defects.

Native Windows/macOS/iOS execution and physical Android device acceptance were not repeated here. Previous native acceptance predates the new mobile/desktop session rewrite and cannot establish its behavior. No correctness fixes were made during this installation/audit task.
