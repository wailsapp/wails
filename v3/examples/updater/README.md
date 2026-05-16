# Updater example

A minimal Wails v3 application that ships in-app updates from a GitHub
release. Demonstrates:

- Configuring `app.Updater` with a single `Provider` (GitHub Releases)
- Wiring a menu item to `app.Updater.CheckAndInstall`
- Subscribing to `updater:*` events from both Go and JavaScript

## Run

```sh
go run .
```

Out of the box the example reports itself as `v1.0.0` and points at
[`wailsapp/updater-demo`](https://github.com/wailsapp/updater-demo),
which publishes a `v2.0.0` release for darwin/arm64, linux/amd64, and
windows/amd64. Click **App → Check for Updates…** and the flow runs
end-to-end: download, SHA256 verify, atomic swap, restart into the
new binary.

Override the defaults via env to point at your own repo:

```sh
APP_VERSION=1.0.0 GH_REPOSITORY=your-org/your-repo go run .
```

## Customising the update window

The framework opens a default window with release notes, a progress bar,
and Install / Skip / Remind / Cancel buttons. Override:

- CSS only: `Window: &updater.BuiltinWindow{CSS: …}`
- Full template: `Window: &updater.BuiltinWindow{HTML: …}`
- Window chrome: `Window: &updater.BuiltinWindow{Options: …}`
- Headless: `Window: updater.WindowNone`

## Verification

Out of the box the example uses **digest verification**: the GitHub
provider is configured with `ChecksumAsset: "SHA256SUMS"`, the
demo's release ships that sidecar, and the framework refuses to
install an artifact whose SHA-256 doesn't match what the sidecar
declares.

For stronger trust, layer on **signature verification** by setting
`publicKey` in `main.go` to your project's Ed25519 / Ed25519ph /
ECDSA-P256 public key. The framework refuses to install a signed
release when no key is configured — by design. The demo doesn't
sign its releases, so leave `publicKey` empty when running against
`wailsapp/updater-demo`.
