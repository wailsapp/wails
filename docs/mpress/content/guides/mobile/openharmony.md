---
title: "OpenHarmony"
description: "Build Wails applications for OpenHarmony — Wake pipeline targets, DevEco SDK setup, HAP packaging and DevEco emulator verification"
slug: "guides/mobile/openharmony"
sourcePath: "guides/mobile/openharmony.md"
---

@note{type="caution" title="Experimental Feature"}
OpenHarmony support is experimental and has only been verified end-to-end on the DevEco Studio x86_64 emulator. Read [Verified status and boundaries](#verified-status-and-boundaries) before relying on it.
@end

Wails v3 applications can be built for OpenHarmony as a HAP using the same `main.go` and frontend you already write for desktop. The Go backend is compiled to a C shared library with an OpenHarmony-capable Go toolchain, and a Stage-model `UIAbility` hosts an ArkWeb view. Assets are served **in-process** from the Go asset server through a local `https://wails.localhost/` origin — no HTTP listener and no open ports. Service bindings, events and `@wailsio/runtime` work through the ArkWeb bridge.

OpenHarmony builds go through the experimental [Wake pipeline](/experimental/wake/) (a `wails.yaml` manifest), not the legacy Taskfile flow.

## Supported targets

| Target | Status |
| --- | --- |
| `openharmony/amd64` | End-to-end verified: builds a HAP and runs on the DevEco x86_64 emulator |
| `openharmony/arm64` | Build verified (`c-shared` strict link + HAP assembly for `arm64-v8a`); not installed or run on a physical device |

@note{type="note"}
The `hap` package format is registered for both targets. The `arm64` target produces an `arm64-v8a` HAP but has no device acceptance evidence, so treat it as build-only for now.
@end

## Requirements

- **DevEco Studio** with the OpenHarmony SDK. The pipeline locates the SDK through `DEVECO_SDK_HOME`.
- The **OpenHarmony-capable Go toolchain** on `PATH` (verified with `go1.27.4-hmos-devel`). Stock Go cannot build `-buildmode=c-shared` for `GOOS=openharmony`, so the custom toolchain's `bin` directory must come first on `PATH`.
- An active `wails.yaml` manifest and the Wake feature opt-in.
- npm, as for other Wails targets, unless the frontend `dist` directory is staged in advance.

### Environment variables

| Variable | Value | Notes |
| --- | --- | --- |
| `WAILS_EXP_USE_WAKE` | `1` | Presence-based opt-in for the Wake/YAML pipeline; set it in every terminal or CI job that builds |
| `DEVECO_SDK_HOME` | `<DevEco Studio>/sdk` | Required so hvigor can locate the OpenHarmony SDK (e.g. `C:/Program Files/Huawei/DevEco Studio/sdk` on Windows) |
| `PATH` | prepend `<OH Go toolchain>/bin` | Must resolve the custom OpenHarmony Go, not stock Go |
| `OHOS_BASE_SDK_HOME` | set by the pipeline | Derived from `DEVECO_SDK_HOME` for hvigor; you do not set it manually |

In PowerShell:

```powershell
$env:WAILS_EXP_USE_WAKE = '1'
$env:DEVECO_SDK_HOME = 'C:\Program Files\Huawei\DevEco Studio/sdk'
$env:PATH = 'C:\Users\<you>\sdk\go1.27.4-hmos\toolchains\<toolchain>\bin;' + $env:PATH
```

In Bash/Zsh:

```bash
export WAILS_EXP_USE_WAKE=1
export DEVECO_SDK_HOME="/path/to/DevEco Studio/sdk"
export PATH="/path/to/go1.27.4-hmos/toolchains/<toolchain>/bin:$PATH"
```

## Building

```bash
wails3 build -targets openharmony/amd64 -formats hap
```

The pipeline generates the OpenHarmony project from the `build_assets/openharmony/` template (an `OpenHarmonyOverlayGen` overlay injects the application's `main` into the Go package), compiles the backend as `libwails_go.so` with `GOOS=openharmony GOARCH=amd64 CGO_ENABLED=1`, and assembles the HAP with hvigor. The artifact is written to `bin/<name>.hap`. The end-to-end verification project produced `bin/ohe2e.hap` (about **11.2 MiB**), containing the x86_64 native libraries `libentry.so`, `libwails_go.so` and `libc++_shared.so`.

To build for arm64, use `-targets openharmony/arm64`; the resulting HAP contains `arm64-v8a` libraries.

## Verifying on the DevEco emulator

With the emulator connected over `hdc`, install and launch the HAP from the project directory:

```bash
hdc -t 127.0.0.1:5555 install bin/ohe2e.hap
hdc -t 127.0.0.1:5555 shell aa start -b com.wails.ohe2e -a EntryAbility
hdc -t 127.0.0.1:5555 shell hilog -x -T Wails
```

The generated template logs under the `Wails` hilog tag. The success signal is:

```
M3 PASS: local HTML/JS/CSS, Greet and Go event
```

The bundle identifier comes from your manifest — the end-to-end run used `com.wails.ohe2e`, and the native HelloWails verification project uses `com.wails.nativehello`. A passing run loads the local origin and its CSS/JS through the Go asset handler, calls `Greet` through the real Wails runtime, and returns Go events to the page, with no local HTTP listener.

## Signing

The pipeline includes an OpenHarmony signing stage (`v3/internal/commands/openharmony_sign.go`). Configure it under `openharmony.signing` in `wails.yaml`:

```yaml
openharmony:
  signing:
    certificate: signing/openharmony.p12           # .p12 keystore
    entitlements: signing/openharmony.cer          # .cer signing certificate
    provisioning_profile: signing/openharmony.p7b  # .p7b provisioning profile
    identity: wails                                # key alias inside the keystore
    credential: OHOS_SIGNING_PASSWORD              # name of the env var holding the keystore password, never the password itself
```

`key_alias` is accepted as an alternative to `identity`. Set the password in the environment (the variable name is the value of `credential`, not the password):

```bash
export OHOS_SIGNING_PASSWORD='<p12 password>'
```

The signing stage injects a `signingConfigs` section into a staging copy of the generated hvigor project and re-runs `hvigorw assembleHap` to produce a signed HAP. See the [YAML manifest fields →](/reference/wails-yaml-fields/) for all `openharmony.signing.*` fields.

@note{type="caution" title="Verification boundary"}
The signing stage has been verified only by the injection unit test (`TestInjectOpenHarmonySigningConfig` PASS). No DevEco signing credentials were available locally, so a real signed-HAP build and installation have **not** been verified end-to-end. The unsigned HAP path above is the one that produced the emulator M3 PASS.
@end

## Verified status and boundaries

@note{type="caution" title="What has actually been tested"}
Only the **DevEco Studio x86_64 emulator (OpenHarmony 7.0.0.107, API 26)** has been verified end-to-end. The first launch after installing a HAP can freeze in ArkWeb initialization (`THREAD_BLOCK_6S`); re-run once and it passes. This is a known emulator cold-start issue, not a bridge regression.
@end

| Area | Status |
| --- | --- |
| `openharmony/amd64` build → HAP → emulator M3 PASS | ✅ End-to-end verified |
| Local-origin assets, `Greet` binding, Go → page events | ✅ Verified, no local HTTP listener |
| `openharmony/arm64` build | ✅ Builds; ❌ not run on a physical device |
| `harmonyos/amd64` build → `hos_hap` (commercial `runtimeOS: HarmonyOS`) | ✅ Builds; ❌ not run on a device — no physical HarmonyOS device in the verification setup |
| Physical arm64 device / HarmonyOS 3.0 compatibility | ❌ Not verified |
| Multiple windows | ❌ Not verified — a single Web window only |
| Full mobile API surface (`application.Mobile` / `IOS` / `Android`) | ❌ Not verified on OpenHarmony |
| Native API surface on `openharmony/amd64` | ⚠️ Emulator probe dispatched 45 calls (`M3 FULLAPI PASS`); real effects and permission paths not asserted |
| Explicit unsupported paths | ✅ Return errors to callers or emit `ohos:unsupported` / `ohos:error`; no silent ACK for missing platform APIs |
| HMS-only capabilities (share, TTS, tray) | ❌ No public OpenHarmony API; unsupported on `openharmony` target, potential only on `harmonyos` target |
| ArkWeb document and host rebuild, network isolation | ✅ In the dedicated native test project; not re-run through the productized pipeline |
| Signed HAP (`openharmony.signing.*`) | ⚠️ `signingConfigs` injection unit test passes (`TestInjectOpenHarmonySigningConfig`); ❌ no real DevEco signing credentials locally — signed-HAP build/install not verified |

These limits are deliberate: a packaged HAP and a compiling arm64 target are not proof of device behaviour. Do not read "builds" as "supported on devices".

@note{type="info" title="Platform naming"}
There are now two platform identifiers in the productized pipeline: `openharmony` (Go `GOOS=openharmony`, `runtimeOS: OpenHarmony`, format `hap`) and `harmonyos` (same Go toolchain and `x86_64-linux-ohos` / `aarch64-linux-ohos` clang targets, but `runtimeOS: HarmonyOS` and format `hos_hap`). The commercial OS is called **HarmonyOS** (local SDK `HarmonyOS 26.0.0`, API 26) — not `HarmonyOS NEXT`, an early name around HarmonyOS 5.0.0 / API 12 in Huawei docs. `go-hmos` / `ohos-go` / `hmgo` are toolchain-directory and ecosystem project names, not extra GOOS values. `harmonyos` is opt-in: it is absent from default targets and must be requested explicitly with `-targets harmonyos/<arch> -formats hos_hap`; signing reuses the `openharmony.signing.*` settings. Verified locally: the `hos_hap` build produced a HAP with `compileSdkType: HarmonyOS` and three x86_64 native libraries; the `openharmony/hap` path still produces `compileSdkType: OpenHarmony`. No HarmonyOS device was available for install/run.
@end

## Next steps

- [Mobile Overview →](/guides/mobile/) — the shared application model across iOS, Android and OpenHarmony.
- [Wake builds →](/experimental/wake/) — the manifest-driven pipeline OpenHarmony builds use.
- [YAML manifest fields →](/reference/wails-yaml-fields/) — including the `openharmony.*` signing settings.
