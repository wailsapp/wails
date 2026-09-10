# HCL run review record

Date: 2026-09-10.

The required `coderabbit --plain` invocation was attempted. The installed CLI
rejects that obsolete option; `coderabbit review` uses plain output by default.
The initial uncommitted review included the implementation and new example files.
A repeat exceeded the free 150-file limit after the launch suite regenerated
frontend outputs, so the follow-up was scoped to `v3/internal` and `v3/cmd`.

## Findings addressed

- Duplicate application identifiers in legacy example metadata would make
  installed examples replace each other. New manifests now use unique IDs;
  the independent example inventory test checks uniqueness.
- A planning timeout stopped the example sweep. The runner now records it and
  continues; a two-example timeout fixture verified both rows and failure status.
- Architecture-specific compiler environments discarded OS-wide values. Target
  environments now merge by key, with architecture values taking precedence.
  The HCL resolution/ejection test checks this behaviour.
- APK rejection diagnostics omitted the new run path. The error now names both
  run and development deployment as supported APK flows.
- Disabled frontends still excluded Go packages under the frontend directory
  from compile snapshots. A regression test failed before the correction and
  passed afterwards. Enabled frontends retain their existing input selectors.

## Additional implementation review fixes

The real CLI parser could turn absent application arguments into an empty slice,
clearing configured defaults. The CLI normalises absent arguments to nil and
preserves an explicit empty `--`. A child-process CLI regression covers defaults,
replacement flags and empty overrides.

macOS registers the built app bundle with Launch Services before executing its
binary. The protocol examples preserve their URL schemes and existing icons.
Physical iOS environment variables use devicectl's JSON option rather than an
unverified environment prefix. Unit tests check bundle paths, registration errors,
JSON encoding and argument preservation. Native platform acceptance remains open.

The final CLI-scoped CodeRabbit review reported no findings. The internal review
reported the disabled-frontend cache issue above; its regression and the pipeline,
commands and CLI tests passed after the correction.
