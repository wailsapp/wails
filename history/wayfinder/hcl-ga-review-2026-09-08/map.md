# HCL build-system GA review follow-up

Review HEAD: `3452def9defa5ad7f711efe1aa14a34ca93cbe41`.

[Full review](../../hcl-ga-review-2026-09-08/review.md). These are newly reproduced or source-confirmed defects from the 8 September review. Existing native acceptance/deployment work remains in the original branch tickets rather than being duplicated here.

## Original findings (resolved; native acceptance remains separate)

- [F01: Embedded resources do not invalidate compiled binaries](issues/01.md) — P1.
- [F02: Binding cache ignores registrations inside methods](issues/02.md) — P2.
- [F03: Installed custom frontend commands are rejected](issues/03.md) — P2.
- [F04: CPU-feature configuration does not invalidate cached compilation](issues/04.md) — P1.
- [F05: Android release versions and SDK policy are silently ignored](issues/05.md) — P1.
- [F06: MSIX packaging signs within a cacheable package stage](issues/06.md) — P1.
- [F07: iOS icons still come from legacy paths](issues/07.md) — P2.
- [F08: iOS background modes never reach Info.plist](issues/08.md) — P2.
- [F09: Complete replacement plists are rewritten](issues/09.md) — P2.
- [F10: Development app assembly can overwrite production bundles](issues/10.md) — P1.
- [F11: Migration silently drops Taskfile environment configuration](issues/11.md) — P1.
- [F12: Migration reports success while dropping VITE_PORT](issues/12.md) — P2.
- [F13: Development from a subdirectory uses the wrong project root](issues/13.md) — P2.
- [F14: Frontend development ignores declared environment](issues/14.md) — P2.
- [F15: Development cleanup leaves descendants alive after parent exit](issues/15.md) — P2.

- [F16: Android production AAB generation silently falls back to debug signing](issues/16.md) — P1.

## Decisions so far

- [Embedded resources do not invalidate compiled binaries](issues/01.md): Compilation snapshots now discover go:embed resources, including resources in excluded build directories and local dependencies. Permanent cache and planner tests cover changes, removal and persistence; a real CLI smoke build changed executable output after an embedded JSON-only edit.
- [Binding cache ignores registrations inside methods](issues/02.md): Binding fingerprints retain receiver bodies containing event/service registration calls and use a new cache key version. Regression tests cover method registrations while ordinary method-body changes remain cache hits.
- [Installed custom frontend commands are rejected](issues/03.md): Explicit frontend command arrays are accepted by HCL validation, probed on the host and fingerprinted using their actual executable. HCL and host-probe tests pass; a real Python frontend builds successfully.
- [CPU-feature configuration does not invalidate cached compilation](issues/04.md): Compiler fingerprints include CPU-feature/compiler environment and the contents of the persistent GOENV file. Regression tests pass; real GOAMD64=v3 then v1 builds report the requested setting with go version -m.
- [Android release versions and SDK policy are silently ignored](issues/05.md): Android version name/code and minimum/target SDK now reach generated assets, with target-specific overrides retained. Real HCL-to-generated-Gradle regression passes. Final APK/AAB metadata verification remains in matching-host release acceptance.
- [MSIX packaging signs within a cacheable package stage](issues/06.md): MSIX packaging no longer passes certificate credentials to the packaging tool; signing belongs to the explicit signing stage. The packaging regression uses a signtool stub that fails if called. Native credentialed signature verification remains release acceptance work.
- [iOS icons still come from legacy paths](issues/07.md): HCL iOS generation explicitly passes the selected/staged icon to Xcode generation, preventing legacy icon fallback. The generated icon-pixel regression passes.
- [iOS background modes never reach Info.plist](issues/08.md): iOS background modes are carried through the plan into generated configuration and plist output. The real HCL-to-assets regression passes.
- [Complete replacement plists are rewritten](issues/09.md): Generated target settings are applied before complete user plist replacements. Regression verifies byte-for-byte custom plist preservation. Native semantic validation of user-provided plists remains part of release acceptance.
- [Development app assembly can overwrite production bundles](issues/10.md): Development app bundles use .wails/dev/<os>-<arch> rather than production output. Planner and actual app assembly regressions verify the production sentinel survives.
- [Migration silently drops Taskfile environment configuration](issues/11.md): Migration now blocks unsupported Taskfile environment/execution policy rather than activating a lossy conversion. Migration regression passes. Automatic translation of arbitrary Taskfile environments is intentionally not implemented.
- [Migration reports success while dropping VITE_PORT](issues/12.md): Migration now blocks non-default VITE_PORT with an explicit --port remediation instead of reporting successful lossy activation. Migration regression passes.
- [Development from a subdirectory uses the wrong project root](issues/13.md): Development uses the loaded manifest project root after discovery. Nested-working-directory regression passes.
- [Frontend development ignores declared environment](issues/14.md): Frontend startup receives the declared environment, and environment-only changes restart the session. Actual shell-process regression verifies both behaviors.
- [Development cleanup leaves descendants alive after parent exit](issues/15.md): Unix parent exit immediately cleans the owned process group. Linux regressions exercise surviving descendants and port release. Windows cleanup after parent exit remains unresolved and is tracked in ticket 17; native Windows release acceptance is still required.
- [Android production AAB generation silently falls back to debug signing](issues/16.md): Generated Android release builds no longer fall back to debug signing or consume legacy keystore signing in the cacheable package stage. Generated Gradle regression verifies unsigned release configuration; credentialed AAB signer verification remains release acceptance work.


## Remaining implementation work

Open follow-ups are recorded in `issues/`, including Windows descendant cleanup. Existing matching-host acceptance and device deployment tickets remain open in the original build-system tracker.

## Additional platform follow-up

- [Windows signing credential](issues/18-windows-signing-credential.md): fixed and regression-tested; native PFX verification remains.
- [Android launch identity](issues/19-android-launch-identity.md): fixed and verified on the emulator, including existing-APK launch and log cancellation.
- [Windows process ownership](issues/17-windows-process-cleanup.md): implemented; awaiting native Windows process tests.

See the [platform verification and manual handoff](../../hcl-ga-review-2026-09-08/platform-fixes.md).
