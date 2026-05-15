# Updater example

A minimal Wails v3 application that ships in-app updates from a GitHub
release. Demonstrates:

- Configuring `app.Updater` with a single `Provider` (GitHub Releases)
- Wiring a menu item to `app.Updater.CheckAndInstall`
- Subscribing to `updater:*` events from both Go and JavaScript

## Run

```sh
APP_VERSION=1.0.0 GH_REPOSITORY=wailsapp/wails go run .
```

`APP_VERSION` is the version the application claims to currently run.
`GH_REPOSITORY` is the `owner/repo` pair the updater checks. Click
**App → Check for Updates…** to trigger the flow.

## Customising the update window

The framework opens a default window with release notes, a progress bar,
and Install / Skip / Remind / Cancel buttons. Override:

- CSS only: `Window: &updater.BuiltinWindow{CSS: …}`
- Full template: `Window: &updater.BuiltinWindow{HTML: …}`
- Window chrome: `Window: &updater.BuiltinWindow{Options: …}`
- Headless: `Window: updater.WindowNone`

## Verification

Set `publicKey` in `main.go` to your project's Ed25519 / Ed25519ph /
ECDSA-P256 public key. The framework refuses to install a signed release
when no key is configured — by design.
