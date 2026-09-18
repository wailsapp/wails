---
title: "Aset Rilis GitHub untuk Pemutakhir"
description: "Cara pemutakhir Wails memilih artefak aplikasi dari GitHub Releases dan menghindari paket penginstal."
slug: "guides/updater-github-release-assets"
sourcePath: "guides/updater-github-release-assets.md"
---

Penyedia GitHub Releases memilih aset rilis menggunakan `AssetMatcher` yang dikonfigurasi. Jika `AssetMatcher` adalah `nil`, penyedia tersebut menggunakan `github.DefaultAssetMatcher`.

## Pencocokan default

Pencocok default mencari platform dan arsitektur saat ini dalam nama file setiap aset. Pencocok ini mengenali alias arsitektur yang umum, termasuk:

- `amd64`, `x86_64`, dan `x64`
- `arm64` dan `aarch64`
- `386`, `i386`, `x86`, dan `ia32`

File sidecar seperti tanda tangan dan checksum diabaikan.

## Aset penginstal

Sebuah Rilis GitHub dapat memuat biner aplikasi yang digunakan oleh pemutakhir sekaligus penginstal konvensional yang ditujukan untuk instalasi pertama kali. Pencocok default mengabaikan aset yang nama filenya dalam huruf kecil:

- memuat `-installer.`
- memuat `_installer.`
- sama persis dengan `installer.exe`

Misalnya, dengan aset Windows berikut:

```text
myapp-windows-amd64.exe
myapp-windows-amd64-installer.exe
```

`DefaultAssetMatcher` memilih `myapp-windows-amd64.exe` dan mengabaikan penginstal. Hal ini mencegah pemutakhir mengganti aplikasi yang sedang berjalan dengan file executable penginstal NSIS atau penginstal yang dikemas dengan cara serupa.

Pemeriksaan ini sengaja dibuat terbatas. Nama aplikasi yang sekadar memuat kata `installer` tetap valid, termasuk:

```text
myinstaller.exe
installer-tool-windows-amd64.exe
myinstaller-windows-amd64.zip
```

## Skema penamaan khusus

Konfigurasikan `AssetMatcher` jika aset rilis Anda tidak mengikuti konvensi penamaan platform dan arsitektur, atau jika Anda memerlukan pemfilteran penginstal yang berbeda:

```go
import (
    "strings"

    "github.com/wailsapp/wails/v3/pkg/updater"
    "github.com/wailsapp/wails/v3/pkg/updater/providers/github"
)

gh, err := github.New(github.Config{
    Repository: "myorg/myapp",
    AssetMatcher: func(req updater.CheckRequest, assets []github.ReleaseAsset) int {
        for i, asset := range assets {
            name := strings.ToLower(asset.Name)
            if strings.Contains(name, req.Platform) &&
                strings.Contains(name, req.Arch) &&
                !strings.Contains(name, "-setup.") {
                return i
            }
        }
        return -1
    },
})
```

Pencocok khusus sepenuhnya menggantikan `DefaultAssetMatcher`, sehingga pencocok tersebut bertanggung jawab mengecualikan tanda tangan, checksum, penginstal, dan aset lain yang tidak boleh diinstal sebagai pembaruan aplikasi.
