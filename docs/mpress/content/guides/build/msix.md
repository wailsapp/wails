---
title: "MSIX Packaging"
description: "Packaging your Wails v3 application as an MSIX package"
slug: "guides/build/msix"
sourcePath: "guides/build/msix.md"
---

MSIX is the modern Windows application packaging format. Wails can generate an MSIX package as part of your Windows build.

Instructions for MSIX packaging are documented in the [Windows Packaging](/guides/build/windows/#msix-package) guide.

## Direct CLI packaging

Run the MSIX tool on Windows, including in CI. The default backend uses `MakeAppx.exe`; signing also needs `signtool.exe`. Both ship with the Windows SDK. The install helper opens the Microsoft Store and, if needed, the SDK download page; complete installation before packaging.

Set the identity in `build/config.yml`, then build and package your executable:

```yaml
info:
  companyName: "Example Corp"
  productName: "MyApp"
  productIdentifier: "com.example.myapp"
  description: "MyApp"
  version: "1.0.0"
```

```powershell
wails3 tool msix-install-tools
wails3 build GOOS=windows
wails3 tool msix --executable bin/myapp.exe --name myapp.exe
```

The direct command writes `MyApp.msix` in the current directory. The [Windows packaging task](/guides/build/windows/#msix-package) instead supplies its own output path.

## CLI options

| Option | Meaning |
| --- | --- |
| `--config` | Configuration file; defaults to `build/config.yml`. |
| `--executable`, `--name` | Existing executable and its filename inside the package; both are required. |
| `--out` | Output file; defaults to `<ProductName>.msix`. |
| `--arch` | Package architecture: `x64` (default), `x86`, `arm`, `arm64`, `x86a64`, or `neutral`. Go aliases `amd64` and `386` are accepted. Match the executable architecture. |
| `--publisher` | Publisher identity; defaults to `CN=<companyName>`. |
| `--cert`, `--cert-password` | PFX certificate path and password for signing. |
| `--use-makeappx` | Use the default Windows SDK packager. |
| `--use-msix-tool` | Opt into `MsixPackagingTool.exe`, which must be in `PATH`. |

## Signing and CI

For distribution outside the Store, sign with a certificate trusted on the target machine. The certificate Subject must match `--publisher` exactly. The MakeAppx backend invokes SignTool with SHA256 when `--cert` is supplied. See [Microsoft's signing guide](https://learn.microsoft.com/en-us/windows/msix/package/sign-msix-package-guide).

This Windows workflow step assumes Wails and the SDK are installed and a preceding step has securely provisioned the PFX file at `CERT_PATH`. A path secret alone does not upload a certificate:

```yaml
- name: MSIX
  if: runner.os == 'Windows'
  shell: pwsh
  run: |
    wails3 build GOOS=windows
    wails3 tool msix --executable bin/myapp.exe --name myapp.exe --publisher "$env:MSIX_PUBLISHER" --cert "$env:CERT_PATH" --cert-password "$env:CERT_PASSWORD"
  env:
    MSIX_PUBLISHER: ${{ vars.MSIX_PUBLISHER }}
    CERT_PATH: ${{ secrets.WINDOWS_CERT_PATH }}
    CERT_PASSWORD: ${{ secrets.WINDOWS_CERT_PASSWORD }}
```

## File associations and assets

Add extensions without a leading dot to `build/config.yml`; the generated manifest adds the dot:

```yaml
fileAssociations:
  - ext: myext
    name: MyApp Document
    description: MyApp Document
    iconName: fileicon
```

Handle runtime file opening as described in [File Associations](/guides/file-associations/). The MakeAppx backend currently copies only the executable and generates transparent placeholder images. It does not import project `Assets/` files or convert `iconName` icons. Use a custom packaging workflow for branded assets or additional DLLs.

| Generated asset | Size (pixels) |
| --- | --- |
| `Square150x150Logo.png` | 150×150 |
| `Square44x44Logo.png` | 44×44 |
| `Wide310x150Logo.png` | 310×150 |
| `StoreLogo.png` | 50×50 |
| `SplashScreen.png` | 620×300 |
| `FileIcon.png` | 44×44 |

`FileIcon.png` is generated only when file associations are configured. These files live in the package's `Assets/` directory.

## Store submission and troubleshooting

Reserve your application in [Partner Center](https://partner.microsoft.com/dashboard) and use its package identity and publisher when preparing the submission. The Store signs MSIX packages during submission; you do not need to buy a signing certificate for that path. See [Microsoft's package requirements](https://learn.microsoft.com/en-us/windows/apps/publish/publish-your-app/msix/app-package-requirements).

If `MakeAppx.exe` or `signtool.exe` cannot be found, install or repair the Windows SDK. Wails searches `PATH` and standard SDK locations. For signing failures, check the certificate Subject, validity and trust on the target machine; see [MSIX troubleshooting](https://learn.microsoft.com/en-us/windows/msix/msix-troubleshooting-guide).
