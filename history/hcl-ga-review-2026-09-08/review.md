# HCL build-system GA code review

Follow-up: [fixes and verification](fixes.md) records subsequent repairs; the findings below describe the original reviewed commit.

Reviewed 8 September 2026. **Do not ship this branch as the v3 GA replacement build system yet.** There are 16 actionable findings: 7 P1 and 9 P2. Priorities describe release impact, not whether a native reproduction was possible.

Branch: `wailsapp/wails:codex/hcl-build-system`, HEAD `3452def9defa5ad7f711efe1aa14a34ca93cbe41`. Downloaded into `/var/home/lea/projects/wails-hcl-review` as an isolated detached worktree. Compared with merge base `f38ae268344700a3f48dd1e95b716f2059f7269e` against fetched master `8f4517fcd72430d815fe9d2fe8a37a23d1b495d5`: 221 files, 42,408 added lines and 767 removed lines. The original working directory and its changes were preserved.

This is a review, not a fixes branch. No implementation changes were committed or pushed. Reproduction tests are archived here outside the active Go packages. New local follow-up tickets are indexed in [the review ticket map](../wayfinder/hcl-ga-review-2026-09-08/map.md).

## Validation

Passed on this Linux host at the reviewed HEAD:

- `go test ./internal/wake/... ./internal/commands ./cmd/wails3 ./internal/templates ./internal/setupwizard ./internal/report/...`
- The same package scope with `go test -race`.
- `go test ./internal/generator/...` (passed with GTK deprecation warnings).
- `go vet` for the same package scope.

The existing passing suite does not cover the regressions described below. Additional focused reproductions establish input/action-key collisions, migration losses, dev failures, and mobile generated-file discrepancies. Details and source are archived alongside this report. Source-confirmed findings are explicitly distinguished from executed reproductions. Native Windows/macOS/iOS packaging, installation, signed releases and device acceptance were not run in this review.

## Standards

The direct-input and tool/environment-aware cache guarantees in `v3/internal/wake/AGENTS.md` are not met. No additional independently established coding-convention violations are claimed. Possible divergent change in the roughly 2,900-line platform handler is a maintenance judgement, not a GA blocker. The four correctness findings in this axis follow.

### F01 — [P1] Embedded resources do not invalidate compiled binaries

[Source](/var/home/lea/projects/wails-hcl-review/v3/internal/wake/pipeline/planner.go:316) · [Follow-up ticket](../wayfinder/hcl-ga-review-2026-09-08/issues/01.md)

Reproduced: editing data.json referenced by //go:embed leaves every input snapshot of the real compile node unchanged.

**Impact:** A normal warm build can silently ship old JSON, SQL migrations, templates, or other embedded resources.

**Cause:** The source snapshot uses a fixed extension allowlist and excludes build, dist and frontend directories. Frontend output dependencies do not cover arbitrary Go embed inputs.

**Required fix/acceptance:** Discover the effective Go compiler inputs, including embedded files and local dependencies. Require an embedded-resource-only edit to invalidate compilation and update the resulting executable.

### F02 — [P2] Binding cache ignores registrations inside methods

[Source](/var/home/lea/projects/wails-hcl-review/v3/internal/wake/cache/cache.go:609) · [Follow-up ticket](../wayfinder/hcl-ga-review-2026-09-08/issues/02.md)

Reproduced: changing `application.RegisterEvent[string]("first")` to `application.RegisterEvent[int]("other")` inside a receiver method produces the same SnapshotGoAPI digest.

**Impact:** Generated event names and payload types can remain stale; NewService calls inside methods have the same underlying risk.

**Cause:** Semantic hashing removes every receiver-method body, while generator event discovery examines registration calls anywhere in the AST/types information.

**Required fix/acceptance:** Retain binding-relevant registrations in semantic fingerprints. Test method-contained event/service registration changes through binding generation.

### F03 — [P2] Installed custom frontend commands are rejected

[Source](/var/home/lea/projects/wails-hcl-review/v3/internal/wake/pipeline/environment.go:125) · [Follow-up ticket](../wayfinder/hcl-ga-review-2026-09-08/issues/03.md)

Reproduced: /usr/bin/python3 exists, but frontend.install = ["python3", "install.py"] fails host planning as a missing tool.

**Impact:** Valid configured commands using Python, corepack, make, or project wrappers cannot run.

**Cause:** Host discovery probes only a fixed knownTools list; validation requires the configured executable to appear in that list.

**Required fix/acceptance:** Resolve and validate the actual configured executable, including supported project-relative commands, and include its identity in caching.

### F04 — [P1] CPU-feature configuration does not invalidate cached compilation

[Source](/var/home/lea/projects/wails-hcl-review/v3/internal/commands/manifest_pipeline.go:512) · [Follow-up ticket](../wayfinder/hcl-ga-review-2026-09-08/issues/04.md)

Reproduced: actual manifestHandler.Identity and cache.ActionKey are identical for GOAMD64=v3 and GOAMD64=v1. This checks cache identity; architecture-specific binaries were not compiled.

**Impact:** A warm rebuild requesting baseline CPU support can reuse an executable built for newer CPUs. GOARM64, GOEXPERIMENT and GOROOT are also absent from this allowlist.

**Cause:** Subprocesses inherit compiler environment values that relevantEnvironment omits.

**Required fix/acceptance:** Fingerprint effective Go toolchain/build configuration. Verify changing GOAMD64 invalidates the cache and inspect the rebuilt executable's build information.

## Spec

The reference is `v3/wep/proposals/manifest-build-system/proposal.md`, read with current HCL documentation and platform handoffs. The design is broadly implemented, but the following twelve behavioral requirements are incomplete or incorrect. Hooks, config checking and Android deployment are documented experiments, so they are not treated as unexplained scope expansion.

### F05 — [P1] Android release versions and SDK policy are silently ignored

[Source](/var/home/lea/projects/wails-hcl-review/v3/internal/commands/manifest_pipeline.go:1608) · [Follow-up ticket](../wayfinder/hcl-ga-review-2026-09-08/issues/05.md)

Reproduced with real HCL loading, planning and asset generation: version_code=42, version_name="2.4.1", minimum_sdk=28, target_sdk=35 generates Gradle values 1, "1.0.0", 21, 36.

**Impact:** Release updates can have the wrong version code/name, and the installed-app compatibility and target-SDK policy differ from the manifest.

**Cause:** Android-specific fields decode successfully but do not reach the effective target/assets adapter; Gradle substitutions use common project fields and never apply target_sdk.

**Required fix/acceptance:** Carry all four Android fields into Gradle generation with defined precedence. Verify APK/AAB metadata with native tools after a schema-to-generated-file regression test.

**Contract:** Platform blocks expose native identity, assets, SDK policy

### F06 — [P1] MSIX packaging signs within a cacheable package stage

[Source](/var/home/lea/projects/wails-hcl-review/v3/internal/commands/manifest_pipeline.go:2313) · [Follow-up ticket](../wayfinder/hcl-ga-review-2026-09-08/issues/06.md)

Source-confirmed call path: packageMSIX passes CertificatePath to ToolMSIX; msix.go:305 automatically runs signtool when a certificate exists. Native Windows signing was not executed.

**Impact:** Signing can occur with sign=false, before before_sign hooks, and its output can be cached. Password-protected PFX files can fail because CertificatePassword is not supplied here.

**Cause:** The legacy packaging helper performs signing implicitly while the new planner models packaging as CacheArtifact and signing as a separate stage.

**Required fix/acceptance:** Package unsigned and sign exclusively through SignArtifact. Test sign=false, hook ordering, protected PFX credentials, repeated signing and native signature verification.

**Contract:** Signing, notarization, credentials, and externally stateful operations are never reusable cache entries.

### F07 — [P2] iOS icons still come from legacy paths

[Source](/var/home/lea/projects/wails-hcl-review/v3/internal/commands/manifest_pipeline.go:1330) · [Follow-up ticket](../wayfinder/hcl-ga-review-2026-09-08/issues/07.md)

Source-confirmed: IOSXcodeGen reads build/appicon.png or the embedded Wails icon (ios_xcode_gen.go:248). The manifest-selected icon is staged elsewhere, and ios.icon is copied after generation to an unused path.

**Impact:** A project can ship the Wails/default or leftover legacy icon despite selecting a custom icon. Legacy icon contents also become an undeclared input.

**Cause:** The HCL adapter delegates icon generation without passing its resolved icon.

**Required fix/acceptance:** Pass the resolved platform/project icon to the asset-catalog generator; test distinct custom images and removal of legacy build assets.

**Contract:** Platform blocks expose native identity, assets, SDK policy

### F08 — [P2] iOS background modes never reach Info.plist

[Source](/var/home/lea/projects/wails-hcl-review/v3/internal/commands/manifest_pipeline.go:2783) · [Follow-up ticket](../wayfinder/hcl-ga-review-2026-09-08/issues/08.md)

Reproduced with real HCL loading, planning and asset generation: background_modes=["audio", "remote-notification"] produces a plist without UIBackgroundModes.

**Impact:** An app relying on background audio or notification behavior does not receive its requested bundle configuration.

**Cause:** writeGeneratedConfig emits identity, associations and protocols but omits the iOS backgroundModes configuration consumed by IOSXcodeGen.

**Required fix/acceptance:** Carry validated modes into generated iOS configuration and assert the final Info.plist contains them.

**Contract:** Platform blocks expose native identity, assets, SDK policy

### F09 — [P2] Complete replacement plists are rewritten

[Source](/var/home/lea/projects/wails-hcl-review/v3/internal/commands/manifest_pipeline.go:1578) · [Follow-up ticket](../wayfinder/hcl-ga-review-2026-09-08/issues/09.md)

Source-confirmed: applyUserPlatformInputs copies the supplied plist, then applyGeneratedTargetSettings rewrites CFBundleExecutable and can rewrite CFBundleVersion/minimum-version values.

**Impact:** The packaged replacement differs from the user-provided file, despite the complete-replacement contract. The original source file itself remains intact.

**Cause:** Generated-setting rewrites run after user input replacement without distinguishing generated from user-owned metadata.

**Required fix/acceptance:** Preserve replacement bytes; reject conflicting settings with a diagnostic. Test byte identity in the assembled artifact, including custom executable/build-number fields.

**Contract:** Wails copies it byte-for-byte into .wails/, preserving executable permissions, with no interpolation, merge, sanitization, or line-ending change.

### F10 — [P1] Development app assembly can overwrite production bundles

[Source](/var/home/lea/projects/wails-hcl-review/v3/internal/wake/pipeline/planner.go:397) · [Follow-up ticket](../wayfinder/hcl-ga-review-2026-09-08/issues/10.md)

Reproduced with actual development packageApp assembly: the selected bin/app.app output replaces a pre-existing production bundle and deletes its sentinel file. The planner coverage test also explicitly expects the production path.

**Impact:** On macOS, starting development after a release build can replace the release .app in build.output with a development bundle.

**Cause:** Only the compile binary is routed under .wails/dev; runnable-app assembly selects the production output when Development is true.

**Required fix/acceptance:** Place every development artifact under .wails/dev and prove a production bundle's digest remains unchanged across dev startup/rebuild.

**Contract:** Dev retains debug information, uses the frontend dev command, and stores transient binaries beneath .wails/dev/.

### F11 — [P1] Migration silently drops Taskfile environment configuration

[Source](/var/home/lea/projects/wails-hcl-review/v3/internal/commands/manifest_migrate.go:521) · [Follow-up ticket](../wayfinder/hcl-ga-review-2026-09-08/issues/11.md)

Reproduced: adding root env.CGO_CFLAGS='-DREVIEW_REQUIRED_FEATURE=1' to the stock badge Taskfile yields Complete=true, no diagnostics, and generated HCL without the setting.

**Impact:** Activation can silently change compilation behavior while declaring the migration safe.

**Cause:** Stock comparison focuses on task ASTs and a few root variables; Taskfile-level environment is neither translated nor treated as a blocker.

**Required fix/acceptance:** Audit file/task/include-level behavioral attributes. Translate supported environment settings and block any reachable behavior without a proven equivalent; add activation tests.

**Contract:** Reachable custom behavior that cannot be proved equivalent is a blocker.

### F12 — [P2] Migration reports success while dropping VITE_PORT

[Source](/var/home/lea/projects/wails-hcl-review/v3/internal/commands/manifest_migrate.go:639) · [Follow-up ticket](../wayfinder/hcl-ga-review-2026-09-08/issues/12.md)

Reproduced: VITE_PORT=9357 yields Complete=true and internal doc.Dev.Port=9357, but encoded HCL contains no 9357. hclDev has no port field; dev defaults to 9245 without a CLI override.

**Impact:** Cutover loses a project's configured development port, potentially causing conflicts or breaking frontend configuration.

**Cause:** The migration writes a field in the internal document that the HCL schema/encoder cannot represent.

**Required fix/acceptance:** Represent the migrated port in the supported contract or explicitly block/report its loss. Test encode/load/dev startup, not only the intermediate document.

**Contract:** It extracts only behavior representable by typed HCL and writes a fresh inactive draft.

### F13 — [P2] Development from a subdirectory uses the wrong project root

[Source](/var/home/lea/projects/wails-hcl-review/v3/internal/commands/manifest_dev.go:128) · [Follow-up ticket](../wayfinder/hcl-ga-review-2026-09-08/issues/13.md)

Reproduced via the dev adapter: manifest loading discovers the parent project, but startFrontend receives the original subdirectory as root.

**Impact:** Running wails3 dev from frontend/ or another subdirectory looks for frontend/frontend, resolves binaries incorrectly, and establishes the wrong watch root.

**Cause:** root remains the result of getwd instead of loaded.Config.Root after upward discovery.

**Required fix/acceptance:** Use the discovered root for all session effects and test real startup from a nested directory.

**Contract:** The directory containing the nearest manifest is the project root.

### F14 — [P2] Frontend development ignores declared environment

[Source](/var/home/lea/projects/wails-hcl-review/v3/internal/commands/manifest_dev.go:604) · [Follow-up ticket](../wayfinder/hcl-ga-review-2026-09-08/issues/14.md)

Reproduced with a real child process: a configured frontend environment variable is empty. Both startup paths only add the Wails port/URL variables.

**Impact:** Frontend dev commands lose declared configuration even though finite frontend build/install handlers receive it; environment-only manifest changes also do not restart the frontend.

**Cause:** startFrontendDev omits config.Frontend.Environment, and frontendSessionChanged omits environment comparison.

**Required fix/acceptance:** Merge declared frontend environment into child startup and make environment-only changes trigger a transactional frontend restart.

**Contract:** wails.hcl describes frontend commands and development policy; manifest reload is transactional.

### F15 — [P2] Development cleanup leaves descendants alive after parent exit

[Source](/var/home/lea/projects/wails-hcl-review/v3/internal/commands/manifest_dev.go:746) · [Follow-up ticket](../wayfinder/hcl-ga-review-2026-09-08/issues/15.md)

Reproduced on Linux: start a shell that launches sleep 60 and exits, wait for the shell, then call stop; the child remains alive. The reproduction explicitly kills its child afterward.

**Impact:** A failed/exited frontend wrapper can leave servers running and ports occupied after the development session terminates.

**Cause:** stop returns immediately when p.done is closed, treating parent exit as proof that the entire process group has exited. It also returns on parent exit after interrupt without checking descendants.

**Required fix/acceptance:** Track and clean the owned process group even after its leader exits, while avoiding PID-reuse hazards. Test early parent exit and interrupt-ignoring descendants.

**Contract:** Wake provides cancellation; wake/AGENTS.md requires terminating and reaping complete process groups on every exit path.

### F16 — [P1] Android production AAB generation silently falls back to debug signing

[Source](/var/home/lea/projects/wails-hcl-review/v3/internal/commands/manifest_pipeline.go:2353) · [Follow-up ticket](../wayfinder/hcl-ga-review-2026-09-08/issues/16.md)

Reproduced in the actual generated Gradle file: release signingConfig selects signingConfigs.debug when ANDROID_KEYSTORE_FILE is absent. packageAndroid invokes bundleRelease using that file. A native signed AAB was not built in this review.

**Impact:** A purported production/unsigned AAB is debug-signed, and packaging can consume undeclared ambient legacy signing credentials before the explicit non-cacheable signing stage.

**Cause:** The new adapter reuses the legacy Android Gradle signing fallback unchanged inside cacheable package generation.

**Required fix/acceptance:** Generate unsigned release bundles for the package stage and route explicit release signing through the named credential contract. Verify the final AAB signer, sign=false behavior, missing credentials, and repeated signing with native tools.

**Contract:** Android production never silently uses a debug keystore.


## Outstanding release work

The branch already has two explicitly open implementation tickets: [Matching-host release verification](/var/home/lea/projects/wails-hcl-review/history/wayfinder/wails-build-system/issues/12-matching-host-release-verification.md) and [Device/emulator deployment](/var/home/lea/projects/wails-hcl-review/history/wayfinder/wails-build-system/issues/14-device-emulator-deployment.md). Do not equate the many resolved implementation tickets with release acceptance.

| Platform | Existing evidence in branch | Outstanding before claiming GA support |
| --- | --- | --- |
| Linux | amd64 binaries and DEB/RPM/Arch/AppImage builds; arm64 cross-build/package evidence | Native arm64 matrix; installation/launch/binding-call smoke tests on target systems for package formats; real package signing |
| Windows | Implementation and cross-compilation evidence | Native amd64/arm64 builds, NSIS/MSIX install and launch, warm-cache correctness, Authenticode including invalid-credential failure paths |
| macOS | Implementation plus September 4 bug fixes | Exact-HEAD native amd64/arm64/universal app and DMG evidence, mount/launch/binding calls, signing and notarization; protect production outputs during dev |
| iOS | Implemented build/package path and structural tests | Xcode simulator build/install/launch; device IPA with provisioning and release credentials; validate requested assets/metadata. Deployment CLI is a separate scope decision |
| Android | AAB ABI matrix and x86_64 emulator acceptance | Physical-device arm64 install/launch/binding/logs/cancellation/disconnect cases; credentialed release AAB verification; verify corrected version/SDK metadata |

The status page is dated August 29 and refers to `acdd11858`; it predates `8fe525f6d` and this HEAD. Absence of recorded native evidence is not proof nobody has performed the tests, but GA needs accessible evidence tied to the actual release candidate.

For each supported platform/format, retain CLI commit, host/tool versions, JSON plan, cold/warm logs, artifact hashes and receipt, install/launch/binding-call evidence, and native signature verification. Reusable second builds must do zero work and preserve bytes. Signing must run again; failed credentials must not publish stale success artifacts.

Complete a schema-to-final-artifact coverage audit with non-default values, especially version/SDK fields, platform metadata, icons, entitlements, capabilities and background modes. Decoder and fake-handler tests alone have allowed accepted configuration to go unused.

Fix the confirmed regressions with permanent tests, integrate onto the intended release branch, then rerun the affected acceptance rows and no-op performance gate (currently 150 ms on the documented Linux benchmark host). The cache fixes affect the benchmark path, so previous timing evidence needs refreshing.

Decide explicitly whether the experimental hooks/config/deployment surfaces are GA-supported, gated, or deferred. iOS deployment implementation remains outstanding if included; arbitrary stages, inline shell/HCL expressions, remote caching and automatic translation of arbitrary Taskfiles are deliberate non-goals, not implicit GA blockers. The branch handoff explicitly says WEP approval is not required.

## Axis totals

Standards: 4 correctness findings (worst P1: stale compilation). Spec: 12 findings (worst P1: incorrect release artifacts, signing, migration and development-output ownership). One optional structural smell is recorded separately and excluded from the 16-bug count.
