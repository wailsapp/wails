# Wails v3 release documentation audit

An audit of the checked-in v3 documentation found **31 actionable findings: 3 P1, 25 P2 and 3 P3**. The highest priorities are unsafe security examples and a misleading claim that binding directives prevent frontend access. The site builds and its links validate, but several documented examples and workflows do not match the implementation.

## Scope and baseline

- Audit date: 25 October 2026. Authoritative tracker: [Wails #6337](https://github.com/wailsapp/wails/issues/6337), following the completed #5430 and #6086 accuracy passes. This report records evidence; follow-up work and status belong on GitHub.
- Audited source: `wailsapp/wails`, branch `releases/v3.0.0`, commit `3e03d19bdc5599c7eac71185f89c6e0596dec1de`; audit branch `codex/v3-release-docs-audit`.
- Source of truth: `docs/mpress/content/`, the v3 module, CLI parser, binding generator, runtime, standard templates, build assets and repository workflows. The Docusaurus website is the v2 site and was excluded.
- Inventory: 204 English pages, 52,259 English source lines, nine translated sets (1,827 translated pages; the new release-policy page has no locale copies), and the supporting contributor/build references. Historical blog posts and changelog entries were inventoried but were not treated as current API promises.
- The English page inventory and site were checked on this release branch. Identifier/source sweeps from the initial audit were reconciled with the changed release-branch files. Semantic review concentrated on onboarding, API references, binding generation, events/lifecycle, build/platform instructions, security, testing and contribution guidance. Streams, updater, mobile, menus, dialogs and environment material received targeted source checks, not exhaustive platform execution.
- This is an audit report, not a correction patch. Production documentation and implementation were left unchanged. Findings were verified against the release checkout and its code; the live-site deployment was not audited.

## Severity and evidence

P1 means copied guidance can compromise application security. P2 means a documented workflow, API contract or validation gate is materially wrong or blocked. P3 means a narrower inconsistency or misleading destination. Counts are distinct remediation items; repeated instances and translation diagnostics are not counted as separate defects.

Validation labels distinguish executed reproductions from source inspection. Native macOS/Windows UI, Linux distribution installs, code signing and mobile device workflows were not executed.

## Validation results

| Check | Result |
| --- | --- |
| M-Press version | 1.0.19, matching the release contributor instructions |
| `mpress build --strict --no-purge-css` | Passed; 2,031 pages, 2,192 files |
| `mpress check` | Passed; no broken links or assets |
| `check_site.py docs/mpress/site` | Failed, nine missing release-policy language routes; all other checks passed |
| Documentation helper unit tests | Passed, 19 tests |
| `check_translations.py --mpress …` | Failed, nine missing translation files; 651 additional non-blocking warnings |
| Security examples | Reproduced nonce reuse and two directory-escape cases |
| Binding registry and server manager probe | Reproduced exposed annotated methods and contradictory return behaviour |
| Built CLI help and project generation | Reproduced invalid flag, signing command and source-file mismatches |
| Background panic | Reproduced process termination, exit status 2 |
| `sanity:gtk4` via Task dry-run | Reproduced missing-target error |

## Findings

### F01 P1 Encryption example reuses the GCM nonce

**Documentation:** `docs/mpress/content/guides/security.md:113`.

The example allocates a zero-filled nonce for every call and immediately passes it to `gcm.Seal`. It never fills the nonce with randomness. Executing the example twice with the same key and plaintext produced identical ciphertext and the nonce `000000000000000000000000`.

**Impact:** Applications copying this security example reuse a GCM nonce under the same key, compromising the security properties of authenticated encryption. This is a defect in the example; this audit does not establish that Wails itself uses this helper.

**Correction:** Replace it with a complete standard-library example using a fresh nonce from `crypto/rand`, or an appropriate random-nonce AEAD API. Explain the nonce prefix and matching decryption. [Go cipher reference](https://pkg.go.dev/crypto/cipher#NewGCM) documents the nonce requirement.

**Validation:** Executed standalone Go reproduction.

### F02 P1 File access example does not confine reads to the allowed directory

**Documentation:** `docs/mpress/content/guides/security.md:169`.

`strings.HasPrefix(absPath, allowedDir)` accepts a sibling directory sharing the prefix. The same code follows symlinks out of the allowed directory. A temporary fixture returned an outside file through both `allowed-outside/secret` and `allowed/link`, with no error.

**Impact:** Readers may expose files outside the intended directory through a bound method copied from the security guide.

**Correction:** Use a filesystem API that confines resolution beneath a root, such as `os.OpenRoot` and root-relative operations, and document supported Go versions. A lexical relative-path check alone does not address symlink escape.

**Validation:** Executed both sibling-prefix and symlink reproductions.

### F03 P1 Binding directives are presented as preventing frontend exposure

**Documentation:** `docs/mpress/content/features/bindings/advanced.md:47`; `docs/mpress/content/features/bindings/advanced.md:59`; `docs/mpress/content/contributing/architecture/bindings.md:88`.

`internal/generator/collect/service.go:239` makes `ignore` omit a generated method and makes `internal` suppress its JS export. The generated directives fixture still contains a callable `InternalMethod`. Runtime registration in `pkg/application/bindings.go:134` reflects exported Go methods and does not read source comments. A server-tagged probe registered both annotated methods and retrieved both by their fully qualified binding names.

**Impact:** A developer can mistake removal from the generated SDK for an access-control boundary and leave privileged exported methods callable.

**Correction:** Describe both directives as code-generation controls. State that exported methods remain registered in normal builds and that private operations must be unexported or moved off the bound service; enforce authorisation inside exposed operations. Qualify any separate obfuscation behaviour rather than promising universal hiding.

**Validation:** Executed binding registry probe with both directives.

### F04 P2 Typed event factories and signatures disagree with generated string constants

**Documentation:** `docs/mpress/content/reference/frontend-runtime.md:130`; `docs/mpress/content/reference/frontend-runtime.md:161`; `docs/mpress/content/reference/frontend-runtime.md:225`.

The reference imports `UserUpdated`/`UserLogin` from `./bindings/events`, invokes them as factories, and describes `EventType<T>` overloads. This release branch **does generate** `events/index.ts`: `internal/generator/render/eventnames.go` and `eventnames.js.tmpl` export string constants. Consequently `Events.On(UserLogin, handler)` and the import path can be valid, but `UserUpdated(payload)` is not callable. The runtime uses `CustomEvents` augmentation and event-name strings.

**Impact:** Factory-style calls and documented overload signatures remain wrong; claiming that the events module is absent would itself be wrong on this release branch. Existing report [#5645](https://github.com/wailsapp/wails/issues/5645) covers the original typed-event discovery problem.

**Correction:** Retain the generated constants and use `Events.Emit(UserUpdated, payload)`, or the string-literal equivalent. Document the runtime's actual name-based generics rather than EventType factories, and validate the example with this generator.

**Validation:** Compared release runtime signatures and generator templates; ran the release event-name renderer test and generated an events fixture.

### F05 P2 Typed events plugin example omits its required bindings path

**Documentation:** `docs/mpress/content/reference/frontend-runtime.md:93`.

The reference configures `plugins: [wails()]`. `src/plugins/vite.ts:25` requires `bindingsRoot: string`, then resolves `${bindingsRoot}/github.com/wailsapp/wails/v3/internal/eventcreate`. Without the argument it attempts a path starting with `undefined/`. The built-in Vanilla config correctly passes `wails("./bindings")`.

**Impact:** The documented plugin configuration fails type checking or event-module resolution.

**Correction:** Pass the bindings directory and explain that it is relative to the Vite project. Align with `guides/dev/frontend-frameworks.md:84` and the built-in template.

**Validation:** Source comparison.

### F06 P2 Time mappings contradict one another and omit the controlling option

**Documentation:** `docs/mpress/content/features/bindings/models.md:319`; `docs/mpress/content/features/bindings/models.md:468`; `docs/mpress/content/features/bindings/models.md:710`; `docs/mpress/content/features/bindings/methods.md:334`; `docs/mpress/content/concepts/bridge.md:131`.

The tables say `time.Time` becomes `any`; other sections promise `Date` and call `.getTime()`. `internal/flags/bindings.go:16` exposes `-time-type` with default `string`. The generated build task passes `-time-type=Date` for class/JS templates, but omits it for TypeScript interface templates. Committed fixtures demonstrate `Time: string` in the interface configuration and `Time: Date` with `DateFromTime` conversion in the class configuration.

**Impact:** Date methods fail for the default interface/string path. Switching generation modes silently changes the expected data shape if the reader follows these pages.

**Correction:** Document string versus Date explicitly, including `-time-type`, template defaults, JSON wire representation, and construction behaviour. Generate examples for both modes; do not replace every mapping with one universal type.

**Validation:** Source comparison.

### F07 P2 Model quick start reads a property excluded by its own JSON tag

**Documentation:** `docs/mpress/content/features/bindings/models.md:21`; `docs/mpress/content/features/bindings/models.md:42`.

The Go example declares `Name string` with `json:"name"`, then JavaScript reads `user.Name`. The later generated-model example on the same page uses `name` correctly.

**Impact:** The very first example returns `undefined` rather than the user name.

**Correction:** Read `user.name` and use one generated fixture to keep the Go definition and frontend property accesses aligned.

**Validation:** Source comparison.

### F08 P2 Counter service example does not compile

**Documentation:** `docs/mpress/content/features/bindings/methods.md:100`.

The stateful counter declares `mu sync.Mutex`, then `GetCount` calls `RLock` and `RUnlock`, methods of `sync.RWMutex`, not `sync.Mutex`.

**Impact:** Copying this complete service fails compilation.

**Correction:** Use `sync.RWMutex`, as the services page already does, or use `Lock`/`Unlock` consistently.

**Validation:** Extracted the complete CounterService example and compiled it with Go; RLock and RUnlock were both rejected on sync.Mutex.

### F09 P2 Custom service names are said to allow multiple instances of one type

**Documentation:** `docs/mpress/content/features/bindings/services.md:292`; `docs/mpress/content/features/bindings/services.md:307`.

The listed use cases for `ServiceOptions.Name` include multiple instances of the same service type. `pkg/application/services.go:22` limits the name to logging/debugging. `Bindings.Add` rejects a duplicate fully qualified method regardless of that display name, explicitly saying to wrap additional instances in dedicated structs.

**Impact:** The suggested registration pattern fails rather than exposing two independently named service instances.

**Correction:** Remove that use case and demonstrate distinct wrapper types. Explain that `Name` does not change a binding namespace.

**Validation:** Source comparison.

### F10 P2 Window creation callback is incorrectly described as runtime ready

**Documentation:** `docs/mpress/content/features/windows/events.md:52`.

`WindowManager.addAndRun` calls the OnCreate callbacks before `runOrDeferToAppRun`. Before `app.Run`, no window runtime has started. The page separately documents `WindowRuntimeReady` as the safe point for JavaScript.

**Impact:** Readers may emit events or call JavaScript too early and lose startup messages.

**Correction:** Describe OnCreate as a configuration/registration callback. Use `WindowRuntimeReady` for operations requiring the frontend runtime.

**Validation:** Traced the exact manager call order.

### F11 P2 Server mode reference gives incorrect return behaviour

**Documentation:** `docs/mpress/content/guides/server-build.md:202`; `docs/mpress/content/guides/server-build.md:212`.

The page says `NewWithOptions` warns and returns nil, and `app.Screen.GetPrimary()` returns an error. The shared manager creates a non-nil `WebviewWindow`; GetPrimary has signature `*Screen` and returns nil when there is no primary screen. A server-tagged probe observed `window nil=false` and `primary nil=true`.

**Impact:** Server adaptations based on the documented nil/error checks will be wrong; blindly dereferencing the screen can panic.

**Correction:** Separate native-host availability from Go window objects. Document exact Go signatures, nil behaviour, and the browser-side errors separately.

**Validation:** Executed a server-tagged manager probe.

### F12 P2 Build guides disagree about whether an ordinary build is debug or production

**Documentation:** `docs/mpress/content/getting-started/your-first-app.md:83`; `docs/mpress/content/getting-started/your-first-app.md:112`; `docs/mpress/content/guides/cli.md:83`; `docs/mpress/content/guides/build/cross-platform.md:21`.

The first-app and CLI guide promise a debug executable and console logging from `wails3 build`. The cross-platform guide correctly calls it production by default. Current Linux/macOS/Windows build assets choose production flags unless `DEV=true` is supplied.

**Impact:** Users cannot follow the promised debugging workflow and may assume inspector/logging behaviour that was compiled out.

**Correction:** Document production as the generated default. Explain `wails3 dev` and any supported explicit DEV override separately, including platform inspector constraints.

**Validation:** Source comparison.

### F13 P2 Several build examples pass variables the generated tasks do not use

**Documentation:** `docs/mpress/content/guides/cli.md:98`; `docs/mpress/content/guides/cli.md:163`; `docs/mpress/content/guides/cli.md:213`; `docs/mpress/content/guides/build/customization.md:125`.

Examples use `PLATFORM=linux CONFIG=production` and `PLATFORM=windows ARCH=amd64`. The wrapper selects `GOOS`; the generated root task dispatches on `.GOOS`; platform tasks select development flags using `.DEV`. PLATFORM and CONFIG do not control the stock tasks.

**Impact:** A command can succeed while building the host platform or an unintended mode.

**Correction:** Use `GOOS=...` and the supported mode variable, and label invented variables as examples for custom Taskfiles rather than built-in controls.

**Validation:** Source comparison.

### F14 P2 macOS deployment promise disagrees with the generated build target

**Documentation:** `docs/mpress/content/getting-started/installation.md:11`.

Installation advertises Intel deployment to 10.13 and Apple Silicon 11.0. `build_assets/darwin/Taskfile.yml:54` sets C flags, linker flags and `MACOSX_DEPLOYMENT_TARGET` to 12.0. Some package-level Cgo directives still name 10.13, so even the source has multiple apparent minima.

**Impact:** The documented support promise does not describe binaries produced by the standard build workflow.

**Correction:** Establish one tested deployment contract, distinguish development-host requirements from binary runtime minimums, and align native and Docker builds. Do not infer a new framework-wide minimum solely from one flag.

**Validation:** Static build-setting contradiction; older macOS execution not tested.

### F15 P2 Linux missing-dependency remedy installs the legacy stack for the default build

**Documentation:** `docs/mpress/content/guides/build/linux.md:184`.

The guide correctly says GTK4/WebKitGTK 6.0 is the default, then its unqualified missing-WebKit remedy installs WebKit2GTK 4.1 runtime packages. Those packages match the legacy `gtk3` build, not the default stack.

**Impact:** A user with a default binary follows the troubleshooting instructions and still lacks its runtime library.

**Correction:** Provide separate GTK4/WebKitGTK 6.0 and legacy GTK3/WebKit2GTK 4.1 remedies, tied to the actual build tag and target distribution. Qualify AppImage portability claims by bundled/external runtime requirements.

**Validation:** Source comparison.

### F16 P2 Custom transport walkthrough omits a required interface method

**Documentation:** `docs/mpress/content/guides/custom-transport.md:43`; `docs/mpress/content/guides/custom-transport.md:73`; `docs/mpress/content/guides/custom-transport.md:120`.

The walkthrough implements Start and Stop and assigns `&MyCustomTransport{}` to Options.Transport. It never implements `JSClient() []byte`, although the reference later on the same page and `pkg/application/transport.go:18` both require it.

**Impact:** The walkthrough fails to compile before a reader can add their transport implementation.

**Correction:** Include JSClient in the first implementation and a complete request/response example. Clearly mark incomplete pseudocode such as `return resp` in the frontend skeleton.

**Validation:** Compiled the documented method set against application.Transport; Go rejected the assignment because JSClient was missing.

### F17 P2 End to end guide points a backend workflow at a frontend-only server

**Documentation:** `docs/mpress/content/guides/e2e-testing.md:33`; `docs/mpress/content/guides/e2e-testing.md:136`.

The guide uses Playwright against port 9245, the Vite UI preview. Wails native invocation is not provided by an ordinary browser tab there. The CI example starts neither Vite nor a Wails server and supplies no mocks or `webServer` configuration. `guides/build/customization.md:239` explicitly warns that bindings do not work at that browser-preview URL.

**Impact:** The claimed complete workflow has no backend connection locally and no running app in CI.

**Correction:** Choose and document a real harness: explicit mocks for UI-only tests, server-tag integration tests for bindings, or a supported native-webview automation path. Start and wait for that harness in CI; state which native capabilities it excludes.

**Validation:** Source comparison.

### F18 P2 End to end examples confuse browser upload and viewport APIs with native UI

**Documentation:** `docs/mpress/content/guides/e2e-testing.md:80`; `docs/mpress/content/guides/e2e-testing.md:100`; `docs/mpress/content/guides/e2e-testing.md:63`.

The save-dialog example handles a Playwright `filechooser` and calls `setFiles`, which is an HTML file-upload operation, not a native save dialog. The window test checks a newly created browser viewport without setting it to 800x600, not a Wails window. The delete test registers its browser dialog handler after the triggering click.

**Impact:** Tests either fail, dismiss the confirmation before handling it, or assert properties unrelated to the native behaviour they claim to validate.

**Correction:** Use uploads only for HTML inputs, register browser-dialog handlers before the action, and distinguish browser viewport assertions from native window checks. [Playwright upload documentation](https://playwright.dev/docs/input#upload-files) and [dialog documentation](https://playwright.dev/docs/dialogs) define the supported operations.

**Validation:** Compared examples with Playwright primary documentation; no GUI automation executed.

### F19 P2 Windows installer output is documented in the wrong directory

**Documentation:** `docs/mpress/content/guides/build/windows.md:22`; `docs/mpress/content/reference/cli.md:108`.

Windows packaging says the installer lands in `build/windows/nsis/<AppName>-installer.exe`; the CLI reference says packaged outputs land in bin. `build_assets/windows/nsis/project.nsi.tmpl:75` sets OutFile to `../../../bin/${INFO_PROJECTNAME}-${ARCH}-installer.exe`.

**Impact:** Users and CI scripts look in the wrong directory and omit the architecture in the expected filename.

**Correction:** Document the bin path and architecture-bearing filename; distinguish NSIS source files from installer output.

**Validation:** Source comparison.

### F20 P2 Setup guide points readers at the wrong defaults file

**Documentation:** `docs/mpress/content/getting-started/setup.md:38`.

The page says `~/.config/wails/config.yaml`. `internal/defaults/defaults.go:112` returns `defaults.yaml` under `$XDG_CONFIG_HOME/wails`, falling back to `~/.config/wails`.

**Impact:** Users inspect or edit a file the setup wizard does not read.

**Correction:** Use defaults.yaml and describe the XDG override. Keep this distinct from per-project build/config.yml.

**Validation:** Source comparison.

### F21 P2 CLI reference invents a global flag spelling and scope

**Documentation:** `docs/mpress/content/reference/cli.md:111`; `docs/mpress/content/guides/cli.md:660`.

The reference claims `--no-colour` works on all commands. `flags.Common` declares NoColour, exposed by clir as `-nocolour` on commands embedding Common. The audited binary rejects `init --no-colour`; accepts `init -nocolour --help`; and rejects even `version -nocolour`, because version does not embed Common.

**Impact:** Copying the global reference produces a flag error.

**Correction:** Document the parser spelling and actual per-command scope, or change the CLI deliberately and verify that all commands really inherit it.

**Validation:** Executed the built CLI help and rejected-flag cases.

### F22 P2 CLI reference denies a signing command that the binary exposes

**Documentation:** `docs/mpress/content/reference/cli.md:71`.

The page explicitly says there is no `wails3 signing` command. The audited binary lists `signing` and `entitlements` at top level, and `wails3 signing --help` prints signing-specific flags. Conversely, `setup signing --help` prints setup help. This occurs despite the apparent nested construction in main.go.

**Impact:** The reference gives an absolute command-availability claim contradicted by the actual parser. Signing setup instructions need an executable check, not just inspection of builder calls.

**Correction:** Reconcile intended and actual command nesting, then document and test the chosen public path. Do not silently replace signing setup instructions based only on how NewSubCommand calls look.

**Validation:** Executed help only; no signing settings or credentials were changed.

### F23 P3 Template contributor guide mixes legacy metadata with the current layout

**Documentation:** `docs/mpress/content/contributing/template-system.md:35`; `docs/mpress/content/contributing/template-system.md:172`; `docs/mpress/content/contributing/template-system.md:220`.

The guide says framework folders contain template.json, while later showing template.yaml and its TypeScript flag. It describes base as the Go layer every template starts from, but base contains frontend and NEXTSTEPS; normal template installation overlays _common and the selected template. The loader supports template.json as a legacy fallback, not the current generated format.

**Impact:** Contributors receive contradictory instructions on where metadata and backend boilerplate belong.

**Correction:** Make YAML the current authoring format, explicitly label the JSON fallback, and describe _common versus the base skeleton using the actual extraction path.

**Validation:** Source comparison.

### F24 P2 Release documentation is missing nine required translations

**Documentation:** `docs/mpress/content/contributing/release-policy.md:1`; `docs/mpress/CONTRIBUTING.md` and `docs/mpress/scripts/check_translations.py`.

The release branch adds an English release-policy page, marked as a draft in its prose, without any of the nine locale copies. It has no build metadata excluding it from publication. The coverage checker reports nine missing translations. The site build includes the English route; its rendered-site validator requires corresponding translated routes. Stale segment findings are warnings unless --strict is supplied; missing translation files still fail the ordinary gate.

**Impact:** Warning-only stale auditing does not waive missing files. The release validation gate remains blocked by this specific page. The standalone M-Press build succeeds; `docs/mpress/scripts/build.sh` additionally runs the failing translation coverage check before invoking that build.

**Correction:** Translate the page when the owner approves its policy, or explicitly exclude an unpublished draft from both publication and coverage under one consistent rule. Do not silently accept missing published routes or treat draft policy as an approved release decision.

**Validation:** Ran the release coverage checker, which reports exactly nine missing files; rebuilt the site with M-Press 1.0.19. Full audit reports nine errors and 651 non-blocking warnings, with the same nine missing routes in the rendered-site check.

### F25 P2 Panic guide says a goroutine panic affects only that goroutine

**Documentation:** `docs/mpress/content/guides/panic-handling.md:15`.

The guide says only the panicking goroutine is affected. A standalone unrecovered goroutine panic terminates the process; the probe exited with status 2. Recovery must occur in that same goroutine. Wails recovery around particular dispatch paths is a separate behaviour.

**Impact:** Readers may assume a background worker crash is isolated and omit necessary recovery at its boundary.

**Correction:** Describe stack unwinding and process termination accurately, then identify precisely which Wails callbacks are recovered. Preserve the existing advice to handle background goroutines explicitly.

**Validation:** Executed standalone background-panic reproduction.

### F26 P3 Database guide links to a persistence example that the destination does not contain

**Documentation:** `docs/mpress/content/guides/patterns/database.md:15`.

The page promises that the TODO tutorial demonstrates SQLite persistence. The tutorial keeps an in-memory slice and instead recommends adding persistence as a future extension at lines 685 and 738.

**Impact:** The linked route exists, so link checks pass, but the promised worked example is absent.

**Correction:** Provide a real persistence example or link to a tutorial with actual persistence, describing its storage format accurately. Label the current database, REST, workflow, WML and Electron migration placeholders as incomplete in navigation.

**Validation:** Source comparison.

### F27 P2 Testing contributor reference names missing targets and overstates its test harness

**Documentation:** `docs/mpress/content/contributing/testing-ci.md:88`; `docs/mpress/content/contributing/testing-ci.md:96`; `docs/mpress/content/contributing/testing-ci.md:117`; `docs/mpress/content/contributing/testing-ci.md:142`.

`sanity:gtk4` does not exist; the default sanity path and `sanity:gtk3` are present. The workflow table also claims Go 1.25, whereas the checked-in matrix uses Go 1.26 and 1.27. The page sends readers to TESTING.md for benign-race guidance and suppression rationale, which that cross-platform example-build guide does not contain. It describes build-and-launch checks where Taskfile example targets run go build. Its workflow summary describes the GTK4 variant as BUILD_TAGS=gtk4, while the workflow now installs GTK4 by default and uses BUILD_TAGS=gtk3 for the legacy variant.

**Impact:** Maintainers cannot run the promised command and may overestimate what a passing example build verifies.

**Correction:** Generate or check target names against the Taskfile, describe build checks separately from launch/runtime tests, and link race guidance only when it exists.

**Validation:** Executed the missing target with Task dry-run; inspected Taskfile and workflow.

### F28 P3 Runtime generation reference gives the wrong debug filename

**Documentation:** `docs/mpress/content/reference/frontend-runtime.md:43`.

The page promises runtime.debug.js. `internal/commands/runtime.go:28` copies the debug source to a destination named runtime-debug.js.

**Impact:** A script or script tag based on the reference points at a nonexistent generated file.

**Correction:** Document runtime-debug.js for command output; distinguish it from the embedded source asset runtime.debug.js.

**Validation:** Source comparison.

### F29 P2 Window options reference contains the wrong callback type and default name

**Documentation:** `docs/mpress/content/features/windows/options.md:59`; `docs/mpress/content/features/windows/options.md:87`.

The supposedly complete options structure declares KeyBindings as `map[string]func(window *WebviewWindow)`. The actual field is `map[string]func(window Window)` at webview_window_options.go:228. The Name description promises a UUID; NewWindow instead uses `window-<numeric ID>`. The displayed structure also omits options such as FullscreenButtonState and OpenInspectorOnStartup.

**Impact:** A callback copied from the type reference fails assignment, and the default identity contract is misstated.

**Correction:** Use the exact public field types and actual default. Generate the structure from source or label it as a selected-field summary.

**Validation:** Source comparison.

### F30 P2 Onboarding switches between TypeScript and JavaScript without changing templates

**Documentation:** `docs/mpress/content/quick-start/first-app.md:30`; `docs/mpress/content/quick-start/first-app.md:126`; `docs/mpress/content/getting-started/your-first-app.md:30`; `docs/mpress/content/tutorials/02-todo-vanilla.md:181`; `docs/mpress/content/tutorials/03-notes-vanilla.md:32`.

The default is correctly described as TypeScript, but the next instructions tell the reader to open frontend/src/main.js. An actual `init -t vanilla` probe generated main.ts, not main.js. The generated index.html itself still names /src/main.js, so part of this is also a template defect, not merely stale prose. TODO/Notes assume JavaScript while selecting the default/vanilla template.

**Impact:** A new reader cannot open the named source file and may create a second entry point or encounter the template mismatch before reaching the lesson.

**Correction:** Choose one consistent lesson variant: explicitly select vanilla-js for JavaScript lessons or use TypeScript files and examples. Independently repair and verify the generated template entrypoint; changing only the prose does not fix that defect.

**Validation:** Executed CLI project generation into a temporary directory.

### F31 P2 Installation routes disagree about the required Go toolchain

**Documentation:** `docs/mpress/content/getting-started/installation.md:25`; `docs/mpress/content/quick-start/installation.md:25`; `docs/mpress/content/status.md:24`.

The public release checkout has conflicting Go requirements: `getting-started/installation.md:25` says Go 1.24, while `quick-start/installation.md:25` and `status.md:24` say Go 1.25. The release `v3/go.mod:3` selects Go 1.26.0 and the common starter go.mod still declares 1.24. Go can select/download a newer toolchain automatically, so this is not proof that every installation using an older Go launcher fails immediately.

**Impact:** Users cannot tell which toolchain a release build selects or whether it requires an automatic download.

**Correction:** Align the stated effective build requirement with v3/go.mod, and distinguish the installed Go launcher from the selected toolchain and the starter module declaration.

**Validation:** Compared both installation routes, release status, v3/go.mod and the generated common module template.

## Structural recommendations

The main maintenance problem is duplicated reference material drifting independently. There are two installation/first-app routes, two CLI references, several event and menu references, and repeated Go-to-JavaScript mapping tables. The CLI still opens the older getting-started route. Consolidate factual reference contracts into one maintained source, and let tutorials link to those facts.

Under Diataxis, keep the first-app and QR/TODO/Notes pages as tutorials with one reproducible path and exact generated filenames. Keep field types, default values, CLI flags and conversion rules in reference pages derived from executable fixtures. Keep how-to guides such as custom transport, signing and end-to-end testing complete enough to achieve their stated outcome. Explanation pages can discuss alternatives without duplicating normative API tables.

Fix the three P1 items first, then the first-app/event/type examples and build contracts, then testing and contribution gates. Avoid a broad rewrite until examples can be regenerated and checked. For upstream triage, these corrections are bugs: prioritise by severity and keep fixes outside the feature-only v3.1.0 milestone/project, following the repository release policy. The findings are tracked in Wails #6337; no PR or push was made. Existing reports were searched before filing; the completed #5430/#6086 and typed-event #5645 context is retained.

## Reproduction evidence

Executed probes used temporary files and directories, the checked-in v3 module and Go 1.26.5. CLI probes were built with the server tag to avoid requiring a desktop display; native signing and packaging actions were not invoked.

```text
Repeated ciphertext: true; nonce: 000000000000000000000000
Sibling bypass: "outside allowed directory" error=<nil>
Symlink bypass: "outside allowed directory" error=<nil>
main.AuditService.InternalMethod registered=true
main.AuditService.IgnoredMethod registered=true
server NewWithOptions nil=false; primary screen nil=true
flag provided but not defined: -no-colour
flag provided but not defined: -nocolour  [version command]
wails signing - Configure code signing  [signing --help]
Task "sanity:gtk4" does not exist
panic: background panic; exit status 2
```

The binding probe used `application.NewBindings(nil, nil)`, added `application.NewService(&AuditService{})` after creating an App, and queried `Bindings.Get(&application.CallOptions{MethodName: name})` for both fully qualified method names. AuditService declared one exported method with each directive. The probe did not start a listening server or invoke privileged methods.

The release translation audit with M-Press 1.0.19 reports nine missing-file errors and 651 warnings: 252 missing-segment, 216 protected-content, 73 untranslated, 54 glossary, 36 structure, nine missing-file, nine invalid-target and two requirement-language findings. These are checker diagnostics, not 651 independently confirmed translation defects.

## Follow up boundaries

This report applies to the checked-in `releases/v3.0.0` baseline above. It makes no claim about the live website or a later branch tip. Native platform, device and signing workflows remain unexecuted where noted. The report contains evidence and suggested corrections; GitHub #6337 is the work tracker, and no documentation or implementation fixes are included in this commit.

## Correction pass — 26 October 2026

The follow-up on `release-prep/v3.0.0` uses `ca6147d6e` as its source baseline.
It preserves the newer English release policy and release post. GitHub #6337
remains the authoritative tracker; this section records verification evidence.

The correction pass addresses the documented security, binding/time/event,
window/server, build/signing, onboarding, CLI, transport, testing and contributor
inconsistencies. Both first-app routes and the TODO/Notes tutorials explicitly
select the JavaScript template used by their examples. Go prerequisites are
aligned with the v3 module's effective Go 1.26 minimum; automatic toolchain
selection remains documented. M-Press setup uses the latest stable release.
Corrected pages have matching translations in all nine published locales, and
the release post's translations are stored under `docs/mpress/drafts/` until
publication. English release policy decisions remain the owner's existing draft.

Executed checks include a newly generated app, generated bindings, frontend
build, a real server-mode greeting request and browser-visible greeting toast;
compiled counter and HTTPS examples; authenticated-encryption round trips and
tamper rejection; rooted file access with sibling, parent and symlink escape
attempts; and runtime probes for service registration and server window/screen
contracts. The strict multilingual site build, local links, documentation helper
checks and mobile/desktop browser overflow checks also run against the changes.

Credits now serve the existing sponsor/contributor SVGs from the v3 site's own
static directory, and the daily image workflow maintains those copies. M-Press
v1.0.23 fixes intrinsic SVG object sizing and preserved bold code beside CJK text.
A separately reproduced translation bug skipped code-only source updates when
no prose was pending; v1.0.24 reconstructs those source-owned blocks while
retaining manual translated prose and avoiding unnecessary provider requests.

Remaining limits: native macOS/Windows UI, platform package installation, code
signing and mobile hardware were not exercised on this Linux host. The underlying
default TypeScript starter's HTML entry-point mismatch (F30) is framework work;
the corrected tutorials avoid that mismatch by selecting the JavaScript template.
The original audit is not evidence that every historical page or every existing
translation has been re-reviewed. Findings outside the corrected pages retain
the repository's previous translation-warning baseline. No live Wails deployment
or release announcement is part of this correction pass.

Final qualification used the published M-Press v1.0.24 Linux archive, verified
against its release checksum. The strict translation checker reports no errors
or warnings across all nine locales, using individually reviewed, file-hash-bound
exceptions for unchanged technical names and two French requirement heuristics;
no structural or protected-content finding is waived. All 477 scoped page/locale
audits and all nine unpublished release-post translations pass. The production
build emits 2,040 pages, the draft preview emits 2,050, local links/assets pass,
the site checker reports 204 pages per language and no errors, and all 19 helper
tests pass. Browser checks pass at 390px and 1440px, including the contributor-grid
proportions. The API overview's four empty bold spans were repaired without
changing their text.

CodeRabbit was run before committing. Its recommendation to use the visible
button text as the accessible name was rejected after checking the current
starter's explicit `aria-label="greet-btn"` with a real browser. The documented
selector matches that label; an explanatory code comment preserves the reason.
The CLI's actual `signing --help` and `entitlements --help` were also executed to
qualify the documented top-level commands against this branch.
